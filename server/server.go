package server

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/prologic/autodock/config"
	"github.com/prologic/autodock/events"
	"github.com/prologic/autodock/proxy"

	"github.com/docker/docker/api/types"
	etypes "github.com/docker/docker/api/types/events"
	"github.com/docker/docker/client"
	"github.com/prologic/msgbus"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

const (
	defaultPollInterval = time.Millisecond * 2000
)

type Server struct {
	cfg           *config.Config
	client        *client.Client
	msgbus        *msgbus.MessageBus
	proxy         *proxy.Proxy
	metrics       *Metrics
	containerHash string
}

var (
	errChan      chan (error)
	eventChan    chan *events.Message
	eventErrChan chan (error)
	restartChan  chan (bool)
	recoverChan  chan (bool)
)

// NewServer ...
func NewServer(cfg *config.Config) (*Server, error) {
	s := &Server{
		cfg:           cfg,
		msgbus:        msgbus.NewMessageBus(&msgbus.Options{}),
		metrics:       NewMetrics(),
		containerHash: "",
	}

	client, err := s.getDockerClient()
	if err != nil {
		return nil, err
	}
	s.client = client

	proxy, err := s.getDockerProxy()
	if err != nil {
		return nil, err
	}
	s.proxy = proxy

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
					m := &events.Message{
						msg,
					}

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

			payload, err := json.Marshal(e)
			if err != nil {
				log.Errorf("error encoding event: %s", err)
			} else {
				topic := s.msgbus.NewTopic(e.Type)
				s.msgbus.Put(s.msgbus.NewMessage(topic, payload))
			}

			// counter
			s.metrics.EventsProcessed.Inc()
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

func (s *Server) waitForSwarm() {
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
func (s *Server) Run() error {
	http.Handle(
		"/events/",
		http.StripPrefix("/events/", s.msgbus),
	)

	http.Handle(
		"/metrics",
		prometheus.Handler(),
	)

	http.Handle(
		"/proxy/",
		http.StripPrefix("/proxy/", s.proxy),
	)

	if err := http.ListenAndServe(s.cfg.Bind, nil); err != nil {
		return err
	}

	return nil
}
