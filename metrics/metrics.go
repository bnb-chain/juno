package metrics

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/forbole/juno/v4/log"
)

type Metrics interface {
	RecordWorkerCount()
	RecordWorkerHeight(WorkerIdx int, chainID string, height uint64)
	RecordWorkerLatencyHist(blockTime time.Time)
	RecordDBBlockCount(totalBlocks int64)
	RecordDBLatestHeight(dbLatestHeight uint64)
	RecordDBLatencyHist(blockTime time.Time)
	TestRecorder()
}

const Namespace = "juno"

type Impl struct {
	WorkerCount       prometheus.Counter
	WorkerHeight      *prometheus.GaugeVec
	WorkerLatencyHist prometheus.Histogram
	DBBlockCount      prometheus.Gauge
	DBLatestHeight    prometheus.Gauge
	DBLatencyHist     prometheus.Histogram
}

func DefaultMetrics() Metrics {
	return &Impl{
		WorkerCount:       DefaultWorkerCount(),
		WorkerHeight:      DefaultWorkerHeight(),
		WorkerLatencyHist: DefaultWorkerLatencyHist(),
		DBBlockCount:      DefaultDBBlockCount(),
		DBLatestHeight:    DefaultDBLatestHeight(),
		DBLatencyHist:     DefaultDBLatencyHist(),
	}
}

func DefaultWorkerCount() prometheus.Counter {
	return prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: Namespace,
			Subsystem: "worker",
			Name:      "count",
			Help:      "Count of active workers.",
		},
	)
}

func DefaultWorkerHeight() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: Namespace,
			Subsystem: "worker",
			Name:      "height",
			Help:      "Height of the last indexed block.",
		},
		[]string{"worker_index", "chain_id"},
	)
}

func DefaultWorkerLatencyHist() prometheus.Histogram {
	return prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: Namespace,
			Subsystem: "worker",
			Name:      "latency",
			Buckets:   prometheus.ExponentialBuckets(0.01, 3, 15),
		},
	)
}

func DefaultDBBlockCount() prometheus.Gauge {
	return prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: Namespace,
			Subsystem: "db",
			Name:      "total_blocks",
			Help:      "Total number of blocks in database.",
		},
	)
}

func DefaultDBLatestHeight() prometheus.Gauge {
	return prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: Namespace,
			Subsystem: "db",
			Name:      "latest_height",
			Help:      "Latest block height in the database.",
		},
	)
}

func DefaultDBLatencyHist() prometheus.Histogram {
	return prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: Namespace,
			Subsystem: "db",
			Name:      "latency",
			Buckets:   prometheus.ExponentialBuckets(0.01, 3, 15),
		},
	)
}

func (i *Impl) RecordWorkerCount() {
	i.WorkerCount.Inc()
}
func (i *Impl) RecordWorkerHeight(workerIdx int, chainID string, height uint64) {
	i.WorkerHeight.WithLabelValues(fmt.Sprintf("%d", workerIdx), chainID).Set(float64(height))
}
func (i *Impl) RecordWorkerLatencyHist(blockTime time.Time) {
	i.WorkerLatencyHist.Observe(float64(time.Since(blockTime).Milliseconds()))
}
func (i *Impl) RecordDBBlockCount(totalBlocks int64) {
	i.DBBlockCount.Set(float64(totalBlocks))
}
func (i *Impl) RecordDBLatestHeight(dbLatestHeight uint64) {
	i.DBLatestHeight.Set(float64(dbLatestHeight))
}
func (i *Impl) RecordDBLatencyHist(blockTime time.Time) {
	i.DBLatencyHist.Observe(float64(time.Since(blockTime).Milliseconds()))
}

func (i *Impl) TestRecorder() {
	log.Warnw("default output from juno", "Error", "Error")
}

func (i *Impl) RegisterMetrics() {
	prometheus.Register(i.WorkerCount)
	prometheus.Register(i.WorkerHeight)
	prometheus.Register(i.WorkerLatencyHist)
	prometheus.Register(i.DBBlockCount)
	prometheus.Register(i.DBLatencyHist)
	prometheus.Register(i.DBLatestHeight)
}
