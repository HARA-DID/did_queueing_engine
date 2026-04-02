package pkg

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics holds all Prometheus counters and histograms for the worker.
type Metrics struct {
	EventsReceived  prometheus.Counter
	EventsProcessed *prometheus.CounterVec 
	EventsRetried   prometheus.Counter
	EventsDLQ       prometheus.Counter
	ProcessDuration prometheus.Histogram
}

func NewMetrics() *Metrics {
	return &Metrics{
		EventsReceived: promauto.NewCounter(prometheus.CounterOpts{
			Name: "worker_events_received_total",
			Help: "Total number of events received from Redis stream.",
		}),
		EventsProcessed: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "worker_events_processed_total",
			Help: "Total number of events processed, partitioned by status.",
		}, []string{"status"}),
		EventsRetried: promauto.NewCounter(prometheus.CounterOpts{
			Name: "worker_events_retried_total",
			Help: "Total number of event retry attempts.",
		}),
		EventsDLQ: promauto.NewCounter(prometheus.CounterOpts{
			Name: "worker_events_dlq_total",
			Help: "Total number of events pushed to the dead-letter queue.",
		}),
		ProcessDuration: promauto.NewHistogram(prometheus.HistogramOpts{
			Name:    "worker_event_process_duration_seconds",
			Help:    "Histogram of event processing durations.",
			Buckets: prometheus.DefBuckets,
		}),
	}
}
