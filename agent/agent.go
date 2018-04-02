package agent

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	etypes "github.com/docker/docker/api/types/events"
	dockerclient "github.com/docker/docker/client"
	msgbusclient "github.com/prologic/msgbus/client"
	log "github.com/sirupsen/logrus"

	"github.com/docker/docker/api/types"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/prologic/autodock/config"
	"github.com/prologic/autodock/events"
	"github.com/prologic/autodock/metrics"
)

// Agent ...
type Agent struct {
	cfg     *config.Config
	client  *dockerclient.Client
	msgbus  *msgbusclient.Client
	metrics *metrics.Metrics
}

var (
	errChan      chan (error)
	eventChan    chan *events.Message
	eventErrChan chan (error)
	restartChan  chan (bool)
	recoverChan  chan (bool)
)

// NewAgent ...
func NewAgent(cfg *config.Config) (*Agent, error) {
	s := &Agent{
		cfg:     cfg,
		msgbus:  msgbusclient.NewClient(cfg.MsgBusURL, nil),
		metrics: metrics.NewMetrics(),
	}

	client, err := s.getDockerClient()
	if err != nil {
		return nil, err
	}
	s.client = client

	// channel setup
	errChan = make(chan error)
	eventErrChan = make(chan error)
	restartChan = make(chan bool)
	recoverChan = make(chan bool)
	eventChan = make(chan *events.Message)

	// eventErrChan handler
	// this handles event stream errors
	go func() {
		for range eventErrChan {
			// error from swarm event stream; attempt to restart
			log.Error("event stream fail; attempting to reconnect")
			s.waitForSwarm()
			restartChan <- true
		}
	}()

	// errChan handler
	// this is a general error handling channel
	go func() {
		for err := range errChan {
			log.Error(err)
			// HACK: check for errors from swarm and restart
			// events.  an example is "No primary manager elected"
			// before the event handler is created and thus
			// won't send the error there
			if strings.Index(err.Error(), "500 Internal Server Error") > -1 {
				log.Error("swarm error detected")

				s.waitForSwarm()

				restartChan <- true
			}
		}
	}()

	// restartChan handler
	go func() {
		for range restartChan {
			log.Debug("starting event handling")

			log.Debug("using event stream")
			ctx, cancel := context.WithCancel(context.Background())
			evtChan, evtErrChan := client.Events(ctx, types.EventsOptions{})
			defer cancel()

			go func(ch <-chan error) {
				for {
					err := <-ch
					eventErrChan <- err
				}
			}(evtErrChan)

			// since the event stream channel is receive
			// only we wrap it to be able to send
			// autodock events on the autodock chan
			go func(ch <-chan etypes.Message) {
				for {
					msg := <-ch
					m := &events.Message{msg}

					eventChan <- m
				}
			}(evtChan)

			// trigger initial load
			eventChan <- &events.Message{
				etypes.Message{
					ID:     "0",
					Status: "autodock-start",
				},
			}

		}
	}()

	go func() {
		for e := range eventChan {
			log.Debugf(
				"event received: id=%s, type=%s status=%s action=%s",
				e.ID, e.Type, e.Status, e.Action,
			)

			if e.ID == "" && e.Type == "" {
				continue
			}

			topic := string(e.Type)
			payload, err := json.Marshal(e)
			if err != nil {
				log.Errorf("error encoding event: %s", err)
			} else {
				// FIXME: We're doing a lot of copies here :/ string -> []byte
				err := s.msgbus.Publish(topic, string(payload))
				if err != nil {
					log.Errorf("error publishing event %s: %s", topic, err)
				} else {
					s.metrics.EventsProcessed.Inc()
				}
			}
		}
	}()

	// uptime ticker
	t := time.NewTicker(time.Second * 1)
	go func() {
		for range t.C {
			s.metrics.Uptime.Inc()
		}
	}()

	// start event handler
	restartChan <- true

	return s, nil
}

func (s *Agent) waitForSwarm() {
	log.Debug("waiting for event stream to become ready")

	for {
		if _, err := s.client.Info(context.Background()); err == nil {
			log.Debug("event stream appears to have recovered; restarting handler")
			return
		}

		log.Warn("event stream not yet ready; retrying")

		time.Sleep(time.Second * 1)
	}
}

// Run ...
func (s *Agent) Run() error {
	http.Handle(
		"/metrics",
		prometheus.Handler(),
	)

	if err := http.ListenAndServe(s.cfg.Bind, nil); err != nil {
		return err
	}

	return nil
}
