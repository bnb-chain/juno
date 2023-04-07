package log

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const Namespace = "juno"

// WorkerCount represents the Telemetry counter used to track the worker count
var WorkerCount = promauto.NewCounter(
	prometheus.CounterOpts{
		Namespace: Namespace,
		Subsystem: "worker",
		Name:      "count",
		Help:      "Count of active workers.",
	},
)

// WorkerHeight represents the Telemetry counter used to track the last indexed height for each worker
var WorkerHeight = promauto.NewGaugeVec(
	prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: "worker",
		Name:      "height",
		Help:      "Height of the last indexed block.",
	},
	[]string{"worker_index", "chain_id"},
)

var WorkerLatencyHist = promauto.NewHistogram(
	prometheus.HistogramOpts{
		Namespace: Namespace,
		Subsystem: "worker",
		Name:      "latency",
		Buckets:   prometheus.ExponentialBuckets(0.01, 3, 15),
	},
)

var DBBlockCount = promauto.NewGauge(
	prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: "db",
		Name:      "total_blocks",
		Help:      "Total number of blocks in database.",
	},
)

// DBLatestHeight represents the Telemetry counter used to track the last indexed height in the database
var DBLatestHeight = promauto.NewGauge(
	prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: "db",
		Name:      "latest_height",
		Help:      "Latest block height in the database.",
	},
)

var DBLatencyHist = promauto.NewHistogram(
	prometheus.HistogramOpts{
		Namespace: Namespace,
		Subsystem: "db",
		Name:      "latency",
		Buckets:   prometheus.ExponentialBuckets(0.01, 3, 15),
	},
)
