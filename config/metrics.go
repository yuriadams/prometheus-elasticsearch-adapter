package config

import "github.com/prometheus/client_golang/prometheus"

// Apps's metrics
var (
	ReceivedSamples = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "received_samples_total",
			Help: "Total number of received samples.",
		},
	)
	SentSamples = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "es_sent_samples_total",
			Help: "Total number of processed samples sent to remote storage.",
		},
		[]string{"remote"},
	)
	FailedSamples = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "es_failed_samples_total",
			Help: "Total number of processed samples which failed on send to remote storage.",
		},
		[]string{"remote"},
	)
	WriteSamples = prometheus.NewSummary(prometheus.SummaryOpts{
		Name: "es_adapter_write_timeseries_samples",
		Help: "How many samples each written timeseries has.",
	})
	ReadDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "es_adapter_read_latency_seconds",
		Help: "How long it took us to respond to read requests.",
	})
	ReadErrors = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "es_adapter_read_failed_total",
		Help: "How many selects from Elasticsearch failed.",
	})
	SentBatchDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "es_sent_batch_duration_seconds",
			Help:    "Duration of sample batch send calls to the remote storage.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"remote"},
	)
)
