package server

import (
	"net/http"
	"time"

	"github.com/prologic/msgbus"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/prologic/autodock/config"
	"github.com/prologic/autodock/metrics"
	"github.com/prologic/autodock/proxy"
)

// Server ...
type Server struct {
	cfg     *config.Config
	msgbus  *msgbus.MessageBus
	proxy   *proxy.Proxy
	metrics *metrics.Metrics
}

// NewServer ...
func NewServer(cfg *config.Config) (*Server, error) {
	s := &Server{
		cfg:     cfg,
		msgbus:  msgbus.NewMessageBus(&msgbus.Options{}),
		metrics: metrics.NewMetrics(),
	}

	// uptime ticker
	t := time.NewTicker(time.Second * 1)
	go func() {
		for range t.C {
			s.metrics.Uptime.Inc()
		}
	}()

	return s, nil
}

// EnableMessageBus ...
func (s *Server) EnableMessageBus() error {
	http.Handle("/events/", http.StripPrefix("/events/", s.msgbus))
	return nil
}

// EnableProxy ...
func (s *Server) EnableProxy() error {
	proxy, err := s.getDockerProxy()
	if err != nil {
		return err
	}

	http.Handle("/proxy/", http.StripPrefix("/proxy/", proxy))

	return nil
}

// Run ...
func (s *Server) Run() error {
	http.Handle("/metrics", prometheus.Handler())

	return http.ListenAndServe(s.cfg.Bind, nil)
}
