package server

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	eventsProcessed = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "autodock",
			Subsystem: "totals",
			Name:      "events_processed",
			Help:      "Total number of events processed",
		},
	)

	uptime = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "autodock",
			Subsystem: "totals",
			Name:      "uptime",
			Help:      "Uptime in seconds",
		},
	)
)

type Metrics struct {
	EventsProcessed prometheus.Counter
	Uptime          prometheus.Counter
}

func NewMetrics() *Metrics {
	prometheus.MustRegister(eventsProcessed)
	prometheus.MustRegister(uptime)

	return &Metrics{
		EventsProcessed: eventsProcessed,
		Uptime:          uptime,
	}
}
