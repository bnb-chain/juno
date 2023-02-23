package parser

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/forbole/juno/v4/common"
	"github.com/forbole/juno/v4/database"
	"github.com/forbole/juno/v4/log"
	"github.com/forbole/juno/v4/models"
	"github.com/forbole/juno/v4/modules"
	"github.com/forbole/juno/v4/node"
	"github.com/forbole/juno/v4/types"
	"github.com/forbole/juno/v4/types/config"
	"github.com/forbole/juno/v4/types/utils"
	"github.com/forbole/juno/v4/utils/syncutils"
	"github.com/gogo/protobuf/proto"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

// Worker defines a job consumer that is responsible for getting and
// aggregating block and associated data and exporting it to a database.
type Worker struct {
	ctx context.Context

	index int

	queue   types.HeightQueue
	codec   codec.Codec
	modules []modules.Module

	node node.Node
	db   database.Database
}

// NewWorker allows to create a new Worker implementation.
func NewWorker(ctx *Context, queue types.HeightQueue, index int) Worker {
	return Worker{
		index:   index,
		codec:   ctx.EncodingConfig.Codec,
		node:    ctx.Node,
		queue:   queue,
		db:      ctx.Database,
		modules: ctx.Modules,
	}
}

// Start starts a worker by listening for new jobs (block heights) from the
// given worker queue. Any failed job is logged and re-enqueued.
func (w Worker) Start() {
	log.WorkerCount.Inc()
	chainID, err := w.node.ChainID()
	if err != nil {
		log.Errorw("error while getting chain ID from the node ", "err", err)
	}

	for i := range w.queue {
		if err := w.ProcessIfNotExists(i); err != nil {
			// re-enqueue any failed job after average block time
			time.Sleep(config.GetAvgBlockTime())

			// TODO: Implement exponential backoff or max retries for a block height.
			go func() {
				log.Errorw("re-enqueueing failed block", "height", i, "err", err)
				w.queue <- i
			}()
		}

		log.WorkerHeight.WithLabelValues(fmt.Sprintf("%d", w.index), chainID).Set(float64(i))
	}
}

// ProcessIfNotExists defines the job consumer workflow. It will fetch a block for a given
// height and associated metadata and export it to a database if it does not exist yet. It returns an
// error if any export process fails.
func (w Worker) ProcessIfNotExists(height int64) error {
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
func (w Worker) Process(height int64) error {
	if height == 0 {
		cfg := config.Cfg.Parser

		genesisDoc, genesisState, err := utils.GetGenesisDocAndState(cfg.GenesisFilePath, w.node)
		if err != nil {
			return fmt.Errorf("failed to get genesis: %s", err)
		}

		return w.HandleGenesis(genesisDoc, genesisState)
	}

	log.Debugw("processing block", "height", height)

	block, err := w.node.Block(height)
	if err != nil {
		return fmt.Errorf("failed to get block from node: %s", err)
	}

	events, err := w.node.BlockResults(height)
	if err != nil {
		return fmt.Errorf("failed to get block results from node: %s", err)
	}

	txs, err := w.node.Txs(block)
	if err != nil {
		return fmt.Errorf("failed to get transactions for block: %s", err)
	}

	vals, err := w.node.Validators(height)
	if err != nil {
		return fmt.Errorf("failed to get validators for block: %s", err)
	}

	return w.ExportBlock(block, events, txs, vals)
}

// ProcessTransactions fetches transactions for a given height and stores them into the database.
// It returns an error if the export process fails.
func (w Worker) ProcessTransactions(height int64) error {
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
			return w.ExportTxs(txs)
		},
		func() error {
			return w.ExportAccounts(txs)
		},
	)
}

// HandleGenesis accepts a GenesisDoc and calls all the registered genesis handlers
// in the order in which they have been registered.
func (w Worker) HandleGenesis(genesisDoc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
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

// SaveValidators persists a list of Tendermint validators with an address and a
// consensus public key. An error is returned if the public key cannot be Bech32
// encoded or if the DB write fails.
func (w Worker) SaveValidators(vals []*tmtypes.Validator) error {
	var validators = make([]*models.Validator, len(vals))
	for index, val := range vals {
		consAddr := sdk.ConsAddress(val.Address).String()

		validators[index] = models.NewValidator(common.HexToAddress(consAddr), models.BytesToPubkey(val.PubKey.Bytes()))
	}

	err := w.db.SaveValidators(w.ctx, validators)
	if err != nil {
		return fmt.Errorf("error while saving validators: %s", err)
	}

	return nil
}

// ExportBlock accepts a finalized block and a corresponding set of transactions
// and persists them to the database along with attributable metadata. An error
// is returned if write fails.
func (w Worker) ExportBlock(
	b *tmctypes.ResultBlock, r *tmctypes.ResultBlockResults, txs []*types.Tx, vals *tmctypes.ResultValidators,
) error {
	// Save all validators
	err := w.SaveValidators(vals.Validators)
	if err != nil {
		return err
	}

	// Make sure the proposer exists
	proposerAddr := sdk.ConsAddress(b.Block.ProposerAddress)
	val := findValidatorByAddr(proposerAddr.String(), vals)
	if val == nil {
		return fmt.Errorf("failed to find validator by proposer address %s: %s", proposerAddr.String(), err)
	}

	// Save the block
	err = w.db.SaveBlock(w.ctx, models.NewBlockFromTmBlock(b, sumGasTxs(txs)))
	if err != nil {
		return fmt.Errorf("failed to persist block: %s", err)
	}

	//currently no need
	// Save the commits
	//err = w.ExportCommit(b.Block.LastCommit, vals)
	//if err != nil {
	//	return err
	//}

	// Call the block handlers
	for _, module := range w.modules {
		if blockModule, ok := module.(modules.BlockModule); ok {
			err = blockModule.HandleBlock(b, r, txs, vals)
			if err != nil {
				log.Errorw("error while handling block", "module", module, "height", b, "err", err)
			}
		}
	}

	// Export the transactions and accounts
	return syncutils.BatchRun(
		func() error {
			return w.ExportTxs(txs)
		},
		func() error {
			return w.ExportAccounts(txs)
		},
	)
}

// ExportCommit accepts a block commitment and a corresponding set of
// validators for the commitment and persists them to the database. An error is
// returned if any write fails or if there is any missed aggregated data.
func (w Worker) ExportCommit(commit *tmtypes.Commit, vals *tmctypes.ResultValidators) error {
	var signatures []*types.CommitSig
	for _, commitSig := range commit.Signatures {
		// Avoid empty commits
		if commitSig.Signature == nil {
			continue
		}

		valAddr := sdk.ConsAddress(commitSig.ValidatorAddress)
		val := findValidatorByAddr(valAddr.String(), vals)
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

	err := w.db.SaveCommitSignatures(w.ctx, signatures)
	if err != nil {
		return fmt.Errorf("error while saving commit signatures: %s", err)
	}

	return nil
}

// saveTx accepts the transaction and persists it inside the database.
// An error is returned if write fails.
func (w Worker) saveTx(tx *types.Tx) error {
	err := w.db.SaveTx(w.ctx, tx)
	if err != nil {
		return fmt.Errorf("failed to handle transaction with hash %s: %s", tx.TxHash, err)
	}
	return nil
}

// handleTx accepts the transaction and calls the tx handlers.
func (w Worker) handleTx(tx *types.Tx) {
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

// handleMessage accepts the transaction and handles messages contained
// inside the transaction.
func (w Worker) handleMessage(index int, msg sdk.Msg, tx *types.Tx) {
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

// ExportTxs accepts a slice of transactions and persists then inside the database.
// An error is returned if write fails.
func (w Worker) ExportTxs(txs []*types.Tx) error {
	// handle all transactions inside the block
	for _, tx := range txs {
		// save the transaction
		err := w.saveTx(tx)
		if err != nil {
			return fmt.Errorf("error while storing txs: %s", err)
		}

		// call the tx handlers
		w.handleTx(tx)

		// handle all messages contained inside the transaction
		sdkMsgs := make([]sdk.Msg, len(tx.Body.Messages))
		for i, msg := range tx.Body.Messages {
			var stdMsg sdk.Msg
			err := w.codec.UnpackAny(msg, &stdMsg)
			if err != nil {
				return err
			}
			sdkMsgs[i] = stdMsg
		}

		// call the msg handlers
		for i, sdkMsg := range sdkMsgs {
			w.handleMessage(i, sdkMsg, tx)
		}
	}

	totalBlocks := w.db.GetTotalBlocks(w.ctx)
	log.DbBlockCount.WithLabelValues("total_blocks_in_db").Set(float64(totalBlocks))

	dbLatestHeight, err := w.db.GetLastBlockHeight(w.ctx)
	if err != nil {
		return err
	}
	log.DbLatestHeight.WithLabelValues("db_latest_height").Set(float64(dbLatestHeight))

	return nil
}

// ExportAccounts accepts a slice of transactions and persists accounts inside the database.
// An error is returned if write fails.
func (w Worker) ExportAccounts(txs []*types.Tx) error {
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
						err := w.db.SaveAccount(context.TODO(), account)
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
