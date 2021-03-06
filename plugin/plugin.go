package plugin

import (
	"fmt"
	"net/http"
	"os"
	"time"

	dockerclient "github.com/docker/docker/client"
	"github.com/prologic/msgbus"
	msgbusclient "github.com/prologic/msgbus/client"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

const (
	apiVersion = "1.39"
)

// RunFunc ...
type RunFunc func(ctx Context) error

// HandlerFunc ...
type HandlerFunc func(id uint64, payload []byte, created time.Time) error

// Context ...
type Context interface {
	On(event string, handler HandlerFunc)
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

	subscriber.Start()
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

	flag.BoolVarP(&version, "version", "v", false, "display version information")
	flag.BoolVarP(&debug, "debug", "d", false, "enable debug logging")

	flag.StringVarP(&host, "host", "h", "localhost", "autodock host to connect to")
	flag.IntVarP(&port, "port", "p", 8000, "autodock port to connect to")

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

	var httpClient *http.Client

	dockerURL := fmt.Sprintf("tcp://%s:%d/proxy", host, port)

	defaultHeaders := map[string]string{
		"User-Agent": fmt.Sprintf("autodock-%s", p.Version),
	}

	docker, err := dockerclient.NewClient(
		dockerURL,
		apiVersion,
		httpClient,
		defaultHeaders,
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
