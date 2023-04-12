package log

import (
	"errors"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const Namespace = "juno"

type MetricsType int

const (
	WorkerCountType MetricsType = iota
	WorkerHeightType
	WorkerLatencyHistType
	DBBlockCountType
	DBLatestHeightType
	DBLatencyHistType
)

type WorkerHeightLabels struct {
	WorkerIdx   int
	ChainID     string
	BlockHeight uint64
}

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

func RecordWorkerCount(_ interface{}) error {
	WorkerCount.Inc()
	return nil
}

func RecordWorkerHeight(workerHeightLabels interface{}) error {

	if heightMetrics, ok := workerHeightLabels.(WorkerHeightLabels); ok {
		WorkerHeight.WithLabelValues(fmt.Sprintf("%d", heightMetrics.WorkerIdx), heightMetrics.ChainID).Set(float64(heightMetrics.BlockHeight))
		return nil
	}
	return errors.New("type error")
}

func RecordDBLatestHeight(dbLatestHeight interface{}) error {
	if dbLatestHeightMetrics, ok := dbLatestHeight.(uint64); ok {
		DBLatestHeight.Set(float64(dbLatestHeightMetrics))
	}
	return errors.New("type error")
}

func RecordDBLatencyHist(dbLatencyHist interface{}) error {
	if dbLatencyHistMetrics, ok := dbLatencyHist.(time.Time); ok {
		DBLatencyHist.Observe(float64(time.Since(dbLatencyHistMetrics).Milliseconds()))
	}
	return errors.New("type error")
}

func RecordWorkerLatencyHist(workerLatencyHist interface{}) error {
	if workerLatencyHistMetrics, ok := workerLatencyHist.(time.Time); ok {
		WorkerLatencyHist.Observe(float64(time.Since(workerLatencyHistMetrics).Milliseconds()))
	}
	return errors.New("type error")
}

func RecordDBBlockCount(dbBlockCount interface{}) error {
	if dbBlockCountMetrics, ok := dbBlockCount.(int64); ok {
		DBBlockCount.Set(float64(dbBlockCountMetrics))
	}
	return errors.New("type error")
}
