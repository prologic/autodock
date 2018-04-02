package collector

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	etypes "github.com/docker/docker/api/types/events"
	dockerclient "github.com/docker/docker/client"
	msgbusclient "github.com/prologic/msgbus/client"
	log "github.com/sirupsen/logrus"

	"github.com/docker/docker/api/types"
	"github.com/prologic/msgbus"

	"github.com/prologic/autodock/config"
	"github.com/prologic/autodock/events"
)

var (
	errChan      chan (error)
	eventChan    chan *events.Message
	eventErrChan chan (error)
	restartChan  chan (bool)
	recoverChan  chan (bool)
)

// Publisher ...
type Publisher interface {
	Publish(topic string, payload []byte) error
}

// MessageBusLocalPublisher ...
type MessageBusLocalPublisher struct {
	msgbus *msgbus.MessageBus
}

// NewMessageBusLocalPublisher ...
func NewMessageBusLocalPublisher(msgbus *msgbus.MessageBus) *MessageBusLocalPublisher {
	return &MessageBusLocalPublisher{msgbus}
}

// Publish ...
func (p *MessageBusLocalPublisher) Publish(topic string, payload []byte) error {
	message := p.msgbus.NewMessage(p.msgbus.NewTopic(topic), payload)
	p.msgbus.Put(message)
	return nil
}

// MessageBusRemotePublisher ...
type MessageBusRemotePublisher struct {
	client *msgbusclient.Client
}

// NewMessageBusRemotePublisher ...
func NewMessageBusRemotePublisher(url string) *MessageBusRemotePublisher {
	client := msgbusclient.NewClient(url, nil)
	return &MessageBusRemotePublisher{client}
}

// Publish ...
func (p *MessageBusRemotePublisher) Publish(topic string, payload []byte) error {
	return p.client.Publish(topic, string(payload))
}

// Collector ...
type Collector struct {
	cfg       *config.Config
	client    *dockerclient.Client
	publisher Publisher
}

// NewCollector ...
func NewCollector(cfg *config.Config, publisher Publisher) (*Collector, error) {
	c := &Collector{
		cfg:       cfg,
		publisher: publisher,
	}

	client, err := c.getDockerClient()
	if err != nil {
		return nil, err
	}
	c.client = client

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
			c.waitForSwarm()
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

				c.waitForSwarm()

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
				err := c.publisher.Publish(topic, payload)
				if err != nil {
					log.Errorf("error publishing event %s: %s", topic, err)
				}
			}
		}
	}()

	// start event handler
	restartChan <- true

	return c, nil
}

func (c *Collector) waitForSwarm() {
	log.Debug("waiting for event stream to become ready")

	for {
		if _, err := c.client.Info(context.Background()); err == nil {
			log.Debug("event stream appears to have recovered; restarting handler")
			return
		}

		log.Warn("event stream not yet ready; retrying")

		time.Sleep(time.Second * 1)
	}
}
