package parser

import (
	"context"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/forbole/juno/v4/database"
	"github.com/forbole/juno/v4/log"
	"github.com/forbole/juno/v4/modules"
	"github.com/forbole/juno/v4/node"
	"github.com/forbole/juno/v4/types"
	"github.com/forbole/juno/v4/types/config"
	"github.com/forbole/juno/v4/types/utils"
)

// Worker defines a job consumer that is responsible for getting and
// aggregating block and associated data and exporting it to a database.
type Worker struct {
	ctx context.Context

	index int

	queue   types.HeightQueue
	codec   codec.Codec
	modules []modules.Module

	node    node.Node
	db      database.Database
	indexer Indexer

	concurrentSync bool
}

// NewWorker allows to create a new Worker implementation.
func NewWorker(ctx *Context, queue types.HeightQueue, index int, concurrentSync bool) *Worker {
	return &Worker{
		index:          index,
		codec:          ctx.EncodingConfig.Marshaler,
		node:           ctx.Node,
		queue:          queue,
		db:             ctx.Database,
		indexer:        DefaultIndexer(ctx.EncodingConfig.Marshaler, ctx.Node, ctx.Database, ctx.Modules),
		modules:        ctx.Modules,
		concurrentSync: concurrentSync,
	}
}

func (w *Worker) SetIndexer(indexer Indexer) {
	w.indexer = indexer
}

// Start starts a worker by listening for new jobs (block heights) from the
// given worker queue. Any failed job is logged and re-enqueued.
func (w *Worker) Start() {
	log.WorkerCount.Inc()
	chainID, err := w.node.ChainID()
	if err != nil {
		log.Errorw("error while getting chain ID from the node ", "err", err)
	}

	for i := range w.queue {
		if err := w.ProcessIfNotExists(i); err != nil {
			if w.concurrentSync {
				// re-enqueue any failed job after average block time
				// TODO: Implement exponential backoff or max retries for a block height.
				go func() {
					log.Errorw("re-enqueueing failed block", "height", i, "err", err)
					w.queue <- i
				}()
				continue
			}

			for err != nil {
				log.Errorw("error while process block", "height", i, "err", err)
				time.Sleep(config.GetAvgBlockTime())
				err = w.ProcessIfNotExists(i)
			}
		} else {
			log.WorkerHeight.WithLabelValues(fmt.Sprintf("%d", w.index), chainID).Set(float64(i))
		}
	}
}

// ProcessIfNotExists defines the job consumer workflow. It will fetch a block for a given
// height and associated metadata and export it to a database if it does not exist yet. It returns an
// error if any export process fails.
func (w *Worker) ProcessIfNotExists(height uint64) error {
	exists, err := w.indexer.Processed(w.ctx, height)
	if err != nil {
		return fmt.Errorf("error while searching for block: %s", err)
	}

	if exists {
		log.Infow("skipping already exported block", "height", height)
		return nil
	}

	return w.Process(height)
}

// Process fetches  a block for a given height and associated metadata and export it to a database.
// It returns an error if any export process fails.
func (w *Worker) Process(height uint64) error {
	log.Infow("processing block", "height", height)

	if height == 0 {
		cfg := config.Cfg.Parser

		genesisDoc, genesisState, err := utils.GetGenesisDocAndState(cfg.GenesisFilePath, w.node)
		if err != nil {
			return fmt.Errorf("failed to get genesis: %s", err)
		}

		return w.indexer.HandleGenesis(genesisDoc, genesisState)
	}

	err := w.indexer.Process(height)

	if err == nil {
		log.Infow("processed block", "height", height)

		totalBlocks := w.indexer.GetBlockRecordNum(context.TODO())
		log.DBBlockCount.Set(float64(totalBlocks))

		dbLatestHeight, err := w.indexer.GetLastBlockRecordHeight(context.TODO())
		if err != nil {
			return err
		}
		log.DBLatestHeight.Set(float64(dbLatestHeight))
	}

	return err
}

// ProcessTransactions fetches transactions for a given height and stores them into the database.
// It returns an error if the export process fails.
func (w *Worker) ProcessTransactions(height int64) error {
	block, err := w.node.Block(height)
	if err != nil {
		return fmt.Errorf("failed to get block from node: %s", err)
	}

	txs, err := w.node.Txs(block)
	if err != nil {
		return fmt.Errorf("failed to get transactions for block: %s", err)
	}

	return w.indexer.ExportTxs(block, txs)
}

// ProcessEvents fetches events for a given height and stores them into the database.
// It returns an error if the export process fails.
func (w *Worker) ProcessEvents(height int64) error {
	blockResults, err := w.node.BlockResults(height)
	block, err := w.node.Block(height)

	if err != nil {
		return fmt.Errorf("failed to get block results from node: %s", err)
	}

	return w.indexer.ExportEvents(w.ctx, block, blockResults)
}
