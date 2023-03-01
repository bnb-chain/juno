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

	abci "github.com/tendermint/tendermint/abci/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/forbole/juno/v4/common"
	"github.com/forbole/juno/v4/database"
	"github.com/forbole/juno/v4/log"
	"github.com/forbole/juno/v4/models"
	"github.com/forbole/juno/v4/modules"
	"github.com/forbole/juno/v4/node"
	"github.com/forbole/juno/v4/types"
	"github.com/forbole/juno/v4/types/config"
	eventutil "github.com/forbole/juno/v4/types/event"
	"github.com/forbole/juno/v4/utils/syncutils"
)

// Worker defines a job consumer that is responsible for getting and aggregating block and associated data and exporting it to a database.
type Worker interface {
	// Start starts a worker by listening for new jobs (block heights) from the given worker queue. Any failed job is logged and re-enqueued.
	Start(typedWorker Worker)

	// ProcessIfNotExists defines the job consumer workflow. It will fetch a block for a given  height and associated metadata and export it to a database if it does not exist yet. It returns an error if any export process fails.
	ProcessIfNotExists(typedWorker Worker, height uint64) error

	// Process fetches  a block for a given height and associated metadata and export it to a database.
	// It returns an error if any export process fails.
	Process(height uint64) error

	// ProcessTransactions fetches transactions for a given height and stores them into the database.
	// It returns an error if the export process fails.
	ProcessTransactions(height int64) error

	// HandleGenesis accepts a GenesisDoc and calls all the registered genesis handlers in the order in which they have been registered.
	HandleGenesis(genesisDoc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error

	// SaveValidators persists a list of Tendermint validators with an address and a consensus public key. An error is returned if the public key cannot be Bech32 encoded or if the DB write fails.
	SaveValidators(vals []*tmtypes.Validator) error

	// ExportBlock accepts a finalized block and a corresponding set of transactions
	// and persists them to the database along with attributable metadata. An error
	// is returned if write fails.
	ExportBlock(args ...interface{}) error

	// ExportCommit accepts a block commitment and a corresponding set of
	// validators for the commitment and persists them to the database. An error is
	// returned if any write fails or if there is any missed aggregated data.
	ExportCommit(commit *tmtypes.Commit, vals *tmctypes.ResultValidators) error

	// SaveTx accepts the transaction and persists it inside the database.
	// An error is returned if write fails.
	SaveTx(tx *types.Tx) error

	// HandleTx accepts the transaction and calls the tx handlers.
	HandleTx(tx *types.Tx)

	// HandleMessage accepts the transaction and handles messages contained
	// inside the transaction.
	HandleMessage(index int, msg sdk.Msg, tx *types.Tx)

	// HandleEvent accepts the transaction and handles events contained inside the transaction.
	HandleEvent(ctx context.Context, index int, event sdk.Event)

	// ExportTxs accepts a slice of transactions and persists then inside the database.
	// An error is returned if write fails.
	ExportTxs(txs []*types.Tx) error

	// ExportAccounts accepts a slice of transactions and persists accounts inside the database.
	// An error is returned if write fails.
	ExportAccounts(txs []*types.Tx) error

	// ExportEvents accepts a slice of transactions and get events in order to save in database.
	ExportEvents(txs []*abci.ResponseDeliverTx) error
}

type CommonWorker struct {
	Ctx            context.Context
	Index          int
	Queue          types.HeightQueue
	Codec          codec.Codec
	Modules        []modules.Module
	Node           node.Node
	DB             database.Database
	ConcurrentSync bool
	WorkerType     string
}

// NewWorker allows to create a new Worker implementation.
func NewWorker(ctx *Context, queue types.HeightQueue, index int, concurrentSync bool, workerType string) CommonWorker {
	return CommonWorker{
		Index:          index,
		Codec:          ctx.EncodingConfig.Codec,
		Node:           ctx.Node,
		Queue:          queue,
		DB:             ctx.Database,
		Modules:        ctx.Modules,
		ConcurrentSync: concurrentSync,
		WorkerType:     workerType,
	}
}

// Start starts a worker by listening for new jobs (block heights) from the
// given worker queue. Any failed job is logged and re-enqueued.
func (w *CommonWorker) Start(typedWorker Worker) {
	log.WorkerCount.Inc()
	chainID, err := w.Node.ChainID()
	if err != nil {
		log.Errorw("error while getting chain ID from the node ", "err", err)
	}

	for i := range w.Queue {
		if err := w.ProcessIfNotExists(typedWorker, i); err != nil {
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
				err = w.ProcessIfNotExists(typedWorker, i)
			}
		}

		log.WorkerHeight.WithLabelValues(fmt.Sprintf("%d", w.Index), chainID).Set(float64(i))
	}
}

// ProcessIfNotExists defines the job consumer workflow. It will fetch a block for a given
// height and associated metadata and export it to a database if it does not exist yet. It returns an
// error if any export process fails.
func (w *CommonWorker) ProcessIfNotExists(worker Worker, height uint64) error {
	exists, err := w.DB.HasBlock(w.Ctx, height)
	if err != nil {
		return fmt.Errorf("error while searching for block: %s", err)
	}

	if exists {
		log.Debugw("skipping already exported block", "height", height)
		return nil
	}

	return worker.Process(height)
}

// ProcessTransactions fetches transactions for a given height and stores them into the database.
// It returns an error if the export process fails.
func (w *CommonWorker) ProcessTransactions(height int64) error {
	block, err := w.Node.Block(height)
	if err != nil {
		return fmt.Errorf("failed to get block from node: %s", err)
	}

	txs, err := w.Node.Txs(block)
	if err != nil {
		return fmt.Errorf("failed to get transactions for block: %s", err)
	}

	return syncutils.BatchRun(
		func() error {
			return w.ExportTxs(txs)
		},
		func() error {
			return w.ExportAccounts(txs)
		},
	)
}

// HandleGenesis accepts a GenesisDoc and calls all the registered genesis handlers
// in the order in which they have been registered.
func (w *CommonWorker) HandleGenesis(genesisDoc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	// Call the genesis handlers
	for _, module := range w.Modules {
		if genesisModule, ok := module.(modules.GenesisModule); ok {
			if err := genesisModule.HandleGenesis(genesisDoc, appState); err != nil {
				log.Errorw("error while handling genesis", "module", module, "err", err)
			}
		}
	}

	return nil
}

// SaveValidators persists a list of Tendermint validators with an address and a
// consensus public key. An error is returned if the public key cannot be Bech32
// encoded or if the DB write fails.
func (w *CommonWorker) SaveValidators(vals []*tmtypes.Validator) error {
	var validators = make([]*models.Validator, len(vals))
	for index, val := range vals {
		consAddr := sdk.ConsAddress(val.Address).String()

		validators[index] = models.NewValidator(common.HexToAddress(consAddr), models.BytesToPubkey(val.PubKey.Bytes()))
	}

	err := w.DB.SaveValidators(w.Ctx, validators)
	if err != nil {
		return fmt.Errorf("error while saving validators: %s", err)
	}

	return nil
}

// ExportCommit accepts a block commitment and a corresponding set of
// validators for the commitment and persists them to the database. An error is
// returned if any write fails or if there is any missed aggregated data.
func (w *CommonWorker) ExportCommit(commit *tmtypes.Commit, vals *tmctypes.ResultValidators) error {
	var signatures []*types.CommitSig
	for _, commitSig := range commit.Signatures {
		// Avoid empty commits
		if commitSig.Signature == nil {
			continue
		}

		valAddr := sdk.ConsAddress(commitSig.ValidatorAddress)
		val := FindValidatorByAddr(valAddr.String(), vals)
		if val == nil {
			return fmt.Errorf("failed to find validator by commit validator address %s", valAddr.String())
		}

		signatures = append(signatures, types.NewCommitSig(
			types.ConvertValidatorAddressToBech32String(commitSig.ValidatorAddress),
			val.VotingPower,
			val.ProposerPriority,
			commit.Height,
			commitSig.Timestamp,
		))
	}

	err := w.DB.SaveCommitSignatures(w.Ctx, signatures)
	if err != nil {
		return fmt.Errorf("error while saving commit signatures: %s", err)
	}

	return nil
}

// SaveTx accepts the transaction and persists it inside the database.
// An error is returned if write fails.
func (w *CommonWorker) SaveTx(tx *types.Tx) error {
	err := w.DB.SaveTx(w.Ctx, tx)
	if err != nil {
		return fmt.Errorf("failed to handle transaction with hash %s: %s", tx.TxHash, err)
	}
	return nil
}

// HandleTx accepts the transaction and calls the tx handlers.
func (w *CommonWorker) HandleTx(tx *types.Tx) {
	// Call the tx handlers
	for _, module := range w.Modules {
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
func (w *CommonWorker) HandleMessage(index int, msg sdk.Msg, tx *types.Tx) {
	// Allow modules to handle the message
	for _, module := range w.Modules {
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
			err := w.Codec.UnpackAny(msgAny, &executedMsg)
			if err != nil {
				log.Errorw("unable to unpack MsgExec inner message", "index", authzIndex, "error", err)
			}

			for _, module := range w.Modules {
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
func (w *CommonWorker) HandleEvent(ctx context.Context, index int, event sdk.Event) {
	// Allow modules to handle the message
	for _, module := range w.Modules {
		if eventModule, ok := module.(modules.EventModule); ok {
			err := eventModule.HandleEvent(ctx, index, event)
			if err != nil {
				log.Errorw("error while handling event", "module", module, "event", event, "err", err)
			}
		}
	}

}

// ExportTxs accepts a slice of transactions and persists then inside the database.
// An error is returned if write fails.
func (w *CommonWorker) ExportTxs(txs []*types.Tx) error {
	// handle all transactions inside the block
	for _, tx := range txs {
		// save the transaction
		err := w.SaveTx(tx)
		if err != nil {
			return fmt.Errorf("error while storing txs: %s", err)
		}

		// call the tx handlers
		w.HandleTx(tx)

		// handle all messages contained inside the transaction
		sdkMsgs := make([]sdk.Msg, len(tx.Body.Messages))
		for i, msg := range tx.Body.Messages {
			var stdMsg sdk.Msg
			err := w.Codec.UnpackAny(msg, &stdMsg)
			if err != nil {
				return err
			}
			sdkMsgs[i] = stdMsg
		}

		// call the msg handlers
		for i, sdkMsg := range sdkMsgs {
			w.HandleMessage(i, sdkMsg, tx)
		}
	}

	totalBlocks := w.DB.GetTotalBlocks(w.Ctx)
	log.DbBlockCount.WithLabelValues("total_blocks_in_db").Set(float64(totalBlocks))

	dbLatestHeight, err := w.DB.GetLastBlockHeight(w.Ctx)
	if err != nil {
		return err
	}
	log.DbLatestHeight.WithLabelValues("db_latest_height").Set(float64(dbLatestHeight))

	return nil
}

// ExportAccounts accepts a slice of transactions and persists accounts inside the database.
// An error is returned if write fails.
func (w *CommonWorker) ExportAccounts(txs []*types.Tx) error {
	// save account
	for _, tx := range txs {
		for _, l := range tx.Logs {
			for _, event := range l.Events {
				for _, attr := range event.Attributes {
					if common.IsHexAddress(attr.Value) {
						account := &models.Account{
							Address: common.HexToAddress(attr.Value),
							TxCount: 1,
						}
						err := w.DB.SaveAccount(context.TODO(), account)
						if err != nil {
							return fmt.Errorf("error while storing account: %s", err)
						}
					}
				}
			}
		}
	}
	return nil
}

// ExportEvents accepts a slice of transactions and get events in order to save in database.
func (w *CommonWorker) ExportEvents(txs []*abci.ResponseDeliverTx) error {
	// get all events in order from the txs within the block
	for _, tx := range txs {
		// handle all events contained inside the transaction
		events := filterEventsType(tx)
		// call the event handlers
		for i, event := range events {
			w.HandleEvent(w.Ctx, i, event)
		}
	}
	return nil
}

func filterEventsType(tx *abci.ResponseDeliverTx) []sdk.Event {
	filteredEvents := make([]sdk.Event, 0)
	for _, event := range tx.Events {
		if _, ok := eventutil.EventProcessedMap[event.Type]; ok {
			filteredEvents = append(filteredEvents, sdk.Event(event))
		}
	}
	return filteredEvents
}
