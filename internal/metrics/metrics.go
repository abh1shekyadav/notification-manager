package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	NotificationsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "notifications_total",
			Help: "Total number of notifications processed, labeled by status.",
		},
		[]string{"type", "status"},
	)

	ProcessingDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "notification_processing_duration_seconds",
			Help:    "Duration taken to process notifications, labeled by type.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"type"},
	)
)

func Init() {
	prometheus.MustRegister(NotificationsTotal, ProcessingDuration)
}
