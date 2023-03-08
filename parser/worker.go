package parser

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/gogo/protobuf/proto"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/forbole/juno/v4/database"
	"github.com/forbole/juno/v4/log"
	"github.com/forbole/juno/v4/modules"
	"github.com/forbole/juno/v4/node"
	"github.com/forbole/juno/v4/types"
	"github.com/forbole/juno/v4/types/config"
	"github.com/forbole/juno/v4/types/utils"
	"github.com/forbole/juno/v4/utils/syncutils"
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

// NewPuppetWorker create a puppet Worker
func NewPuppetWorker(modules []modules.Module) *Worker {
	return &Worker{
		index:          -1,
		modules:        modules,
		concurrentSync: false,
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
		}

		log.Infow("processed block", "height", i)
		log.WorkerHeight.WithLabelValues(fmt.Sprintf("%d", w.index), chainID).Set(float64(i))
	}
}

// ProcessIfNotExists defines the job consumer workflow. It will fetch a block for a given
// height and associated metadata and export it to a database if it does not exist yet. It returns an
// error if any export process fails.
func (w *Worker) ProcessIfNotExists(height uint64) error {
	exists, err := w.db.HasBlock(w.ctx, height)
	if err != nil {
		return fmt.Errorf("error while searching for block: %s", err)
	}

	if exists {
		log.Debugw("skipping already exported block", "height", height)
		return nil
	}

	return w.Process(height)
}

// Process fetches  a block for a given height and associated metadata and export it to a database.
// It returns an error if any export process fails.
func (w *Worker) Process(height uint64) error {
	log.Debugw("processing block", "height", height)

	if height == 0 {
		cfg := config.Cfg.Parser

		genesisDoc, genesisState, err := utils.GetGenesisDocAndState(cfg.GenesisFilePath, w.node)
		if err != nil {
			return fmt.Errorf("failed to get genesis: %s", err)
		}

		return w.HandleGenesis(genesisDoc, genesisState)
	}

	return w.indexer.Process(height)
}

// HandleGenesis accepts a GenesisDoc and calls all the registered genesis handlers
// in the order in which they have been registered.
func (w *Worker) HandleGenesis(genesisDoc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	// Call the genesis handlers
	for _, module := range w.modules {
		if genesisModule, ok := module.(modules.GenesisModule); ok {
			if err := genesisModule.HandleGenesis(genesisDoc, appState); err != nil {
				log.Errorw("error while handling genesis", "module", module, "err", err)
			}
		}
	}

	return nil
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

	return syncutils.BatchRun(
		func() error {
			return w.indexer.ExportTxs(txs)
		},
		func() error {
			return w.indexer.ExportAccounts(block, txs)
		},
	)
}

func (w *Worker) HandleBlock(block *tmctypes.ResultBlock, events *tmctypes.ResultBlockResults, txs []*types.Tx, vals *tmctypes.ResultValidators) {
	for _, module := range w.modules {
		if blockModule, ok := module.(modules.BlockModule); ok {
			err := blockModule.HandleBlock(block, events, txs, vals)
			if err != nil {
				log.Errorw("error while handling block", "module", module.Name(), "height", block.Block.Height, "err", err)
			}
		}
	}
}

// HandleTx accepts the transaction and calls the tx handlers.
func (w *Worker) HandleTx(tx *types.Tx) {
	// Call the tx handlers
	for _, module := range w.modules {
		if transactionModule, ok := module.(modules.TransactionModule); ok {
			err := transactionModule.HandleTx(tx)
			if err != nil {
				log.Errorw("error while handling transaction", "module", module.Name(), "height", tx.Height,
					"txHash", tx.TxHash, "err", err)
			}
		}
	}
}

// HandleMessage accepts the transaction and handles messages contained
// inside the transaction.
func (w *Worker) HandleMessage(index int, msg sdk.Msg, tx *types.Tx) {
	// Allow modules to handle the message
	for _, module := range w.modules {
		if messageModule, ok := module.(modules.MessageModule); ok {
			err := messageModule.HandleMsg(index, msg, tx)
			if err != nil {
				log.Errorw("error while handling message", "module", module, "height", tx.Height,
					"txHash", tx.TxHash, "msg", proto.MessageName(msg), "err", err)
			}
		}
	}

	// If it's a MsgExecute, we need to make sure the included messages are handled as well
	if msgExec, ok := msg.(*authz.MsgExec); ok {
		for authzIndex, msgAny := range msgExec.Msgs {
			var executedMsg sdk.Msg
			err := w.codec.UnpackAny(msgAny, &executedMsg)
			if err != nil {
				log.Errorw("unable to unpack MsgExec inner message", "index", authzIndex, "error", err)
			}

			for _, module := range w.modules {
				if messageModule, ok := module.(modules.AuthzMessageModule); ok {
					err = messageModule.HandleMsgExec(index, msgExec, authzIndex, executedMsg, tx)
					if err != nil {
						log.Errorw("error while handling message", "module", module, "height", tx.Height,
							"txHash", tx.TxHash, "msg", proto.MessageName(executedMsg), "err", err)
					}
				}
			}
		}
	}
}

// HandleEvent accepts the transaction and handles events contained inside the transaction.
func (w *Worker) HandleEvent(index int, event sdk.Event) {
	// Allow modules to handle the message
	for _, module := range w.modules {
		if eventModule, ok := module.(modules.EventModule); ok {
			err := eventModule.HandleEvent(index, event)
			if err != nil {
				log.Errorw("error while handling event", "module", module, "event", event, "err", err)
			}
		}
	}
}
