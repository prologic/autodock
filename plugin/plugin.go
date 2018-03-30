package plugin

import (
	"flag"
	"fmt"
	"os"
	"time"

	dockerclient "github.com/docker/docker/client"
	msgbusclient "github.com/prologic/msgbus/client"
	log "github.com/sirupsen/logrus"

	"github.com/prologic/msgbus"
)

const dockerAPIVersion = "1.30"

// RunFunc ...
type RunFunc func(ctx Context) error

// HandlerFunc ...
type HandlerFunc func(id uint64, payload []byte, created time.Time) error

// Context ...
type Context interface {
	On(event string, handler HandlerFunc)
	MessageBus() *msgbusclient.Client
	Docker() *dockerclient.Client
}

type pluginContext struct {
	host   string
	port   int
	msgbus *msgbusclient.Client
	docker *dockerclient.Client
	topics map[string]*msgbusclient.Subscriber
}

// On ...
func (ctx *pluginContext) On(event string, handler HandlerFunc) {
	subscriber := ctx.msgbus.Subscribe(
		event,
		func(msg *msgbus.Message) error {
			return handler(msg.ID, msg.Payload, msg.Created)
		},
	)

	ctx.topics[event] = subscriber

	go subscriber.Run()
}

// MessageBus ...
func (ctx *pluginContext) MessageBus() *msgbusclient.Client {
	return ctx.msgbus
}

// Docker ...
func (ctx *pluginContext) Docker() *dockerclient.Client {
	return ctx.docker
}

// Plugin ...
type Plugin struct {
	ctx         Context
	Name        string
	Version     string
	Description string

	Run RunFunc
}

func (p *Plugin) init() error {
	var (
		version bool
		debug   bool
		host    string
		port    int
	)

	flag.BoolVar(&version, "v", false, "display version information")
	flag.BoolVar(&debug, "d", false, "enable debug logging")

	flag.StringVar(&host, "host", "localhost", "autodock host to connect to")
	flag.IntVar(&port, "port", 8000, "autodock port to connect to")

	flag.Parse()

	if version {
		fmt.Printf("%s v%s", p.Name, p.Version)
		os.Exit(0)
	}

	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	msgbus := msgbusclient.NewClient(
		fmt.Sprintf("http://%s:%d/events/", host, port),
		nil,
	)

	docker, err := dockerclient.NewClientWithOpts(
		dockerclient.WithHost(
			fmt.Sprintf("tcp://%s:%d/proxy/", host, port),
		),
		dockerclient.WithVersion(dockerAPIVersion),
	)
	if err != nil {
		return err
	}

	p.ctx = &pluginContext{
		msgbus: msgbus,
		docker: docker,
		topics: make(map[string]*msgbusclient.Subscriber),
	}

	return nil
}

// Execute ..
func (p *Plugin) Execute() error {
	p.init()
	return p.Run(p.ctx)
}
