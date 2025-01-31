package metrics

import (
	"time"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/models"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// How often we poll for updates.
const metricsPollingInterval = 1 * time.Minute

// CollectedMetrics stores different collected + timestamped values.
type CollectedMetrics struct {
	CPUUtilizations  []timestampedValue `json:"cpu"`
	RAMUtilizations  []timestampedValue `json:"memory"`
	DiskUtilizations []timestampedValue `json:"disk"`
}

// Metrics is the shared Metrics instance.
var Metrics *CollectedMetrics

// Start will begin the metrics collection and alerting.
func Start(getStatus func() models.Status) {
	host := data.GetServerURL()
	if host == "" {
		host = "unknown"
	}
	labels = map[string]string{
		"version": config.VersionNumber,
		"host":    host,
	}

	activeViewerCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name:        "owncast_instance_active_viewer_count",
		Help:        "The number of viewers.",
		ConstLabels: labels,
	})

	activeChatClientCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name:        "owncast_instance_active_chat_client_count",
		Help:        "The number of connected chat clients.",
		ConstLabels: labels,
	})

	chatUserCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name:        "owncast_instance_total_chat_users",
		Help:        "The total number of chat users on this Owncast instance.",
		ConstLabels: labels,
	})

	currentChatMessageCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name:        "owncast_instance_current_chat_message_count",
		Help:        "The number of chat messages currently saved before cleanup.",
		ConstLabels: labels,
	})

	cpuUsage = promauto.NewGauge(prometheus.GaugeOpts{
		Name:        "owncast_instance_cpu_use_pct",
		Help:        "CPU percentage used as seen within Owncast",
		ConstLabels: labels,
	})

	Metrics = new(CollectedMetrics)
	go startViewerCollectionMetrics()

	for range time.Tick(metricsPollingInterval) {
		handlePolling()
	}
}

func handlePolling() {
	// Collect hardware stats
	collectCPUUtilization()
	collectRAMUtilization()
	collectDiskUtilization()

	// Alerting
	handleAlerting()
}
