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

var DbBlockCount = promauto.NewGauge(
	prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: "db",
		Name:      "total_blocks",
		Help:      "Total number of blocks in database.",
	},
)

// DbLatestHeight represents the Telemetry counter used to track the last indexed height in the database
var DbLatestHeight = promauto.NewGauge(
	prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: "db",
		Name:      "latest_height",
		Help:      "Latest block height in the database.",
	},
)
