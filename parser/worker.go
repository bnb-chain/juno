package parser

import (
	"fmt"
	"github.com/forbole/juno/v4/log"
	"github.com/forbole/juno/v4/types"
	"github.com/forbole/juno/v4/types/config"
	"time"
)

// NewWorker allows to create a new Worker implementation.
func NewWorker(indexer Indexer, queue types.HeightQueue, workerIndex int, concurrentSync bool, workerType string) Worker {
	return Worker{
		Indexer:        indexer,
		WorkerIndex:    workerIndex,
		Queue:          queue,
		ConcurrentSync: concurrentSync,
		WorkerType:     workerType,
	}
}

// Worker defines a job consumer that is responsible for getting and
// aggregating block and associated data and exporting it to a database.
type Worker struct {
	Indexer        Indexer
	WorkerIndex    int
	Queue          types.HeightQueue
	ConcurrentSync bool
	WorkerType     string
}

// Start starts a worker by listening for new jobs (block heights) from the
// given worker queue. Any failed job is logged and re-enqueued.
func (w *Worker) Start() {
	log.WorkerCount.Inc()
	chainID, err := w.Indexer.ChainID()
	if err != nil {
		log.Errorw("error while getting chain ID from the node ", "err", err)
	}

	for i := range w.Queue {
		if err := w.ProcessIfNotExists(i); err != nil {
			if w.ConcurrentSync {
				// re-enqueue any failed job after average block time
				// TODO: Implement exponential backoff or max retries for a block height.
				go func() {
					log.Errorw("re-enqueueing failed block", "height", i, "err", err)
					w.Queue <- i
				}()
				continue
			}

			for err != nil {
				time.Sleep(config.GetAvgBlockTime())
				err = w.ProcessIfNotExists(i)
			}
		}

		log.WorkerHeight.WithLabelValues(fmt.Sprintf("%d", w.WorkerIndex), chainID).Set(float64(i))
	}
}

// ProcessIfNotExists defines the job consumer workflow. It will fetch a block for a given
// height and associated metadata and export it to a database if it does not exist yet. It returns an
// error if any export process fails.

func (w *Worker) ProcessIfNotExists(height uint64) error {
	exists, err := w.Indexer.HasBlock(height)
	if err != nil {
		return fmt.Errorf("error while searching for block: %s", err)
	}

	if exists {
		log.Debugw("skipping already exported block", "height", height)
		return nil
	}

	return w.Indexer.ProcessBlock(height)
}
