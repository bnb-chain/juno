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

	"github.com/forbole/juno/v4/common"
	"github.com/forbole/juno/v4/database"
	"github.com/forbole/juno/v4/log"
	"github.com/forbole/juno/v4/models"
	"github.com/forbole/juno/v4/modules"
	blockmodule "github.com/forbole/juno/v4/modules/block"
	"github.com/forbole/juno/v4/node"
	"github.com/forbole/juno/v4/types"
	"github.com/forbole/juno/v4/utils/syncutils"
)

type Indexer interface {
	// Process fetches a block for a given height and associated metadata and export it to a database.
	// It returns an error if any export process fails.
	Process(height uint64) error

	// Processed tells whether the current Indexer has already processed the given height of Block
	// An error is returned if the operation fails.
	Processed(ctx context.Context, height uint64) (bool, error)

	GetLatestHeight(ctx context.Context) uint64

	// ExportBlock accepts a finalized block and persists then inside the database.
	// An error is returned if write fails.
	ExportBlock(block *tmctypes.ResultBlock, events *tmctypes.ResultBlockResults, txs []*types.Tx, vals *tmctypes.ResultValidators) error

	// ExportTxs accepts a slice of transactions and persists then inside the database.
	// An error is returned if write fails.
	ExportTxs(block *tmctypes.ResultBlock, txs []*types.Tx) error

	// ExportValidators accepts ResultValidators and persists validators inside the database.
	// An error is returned if write fails.
	ExportValidators(block *tmctypes.ResultBlock, vals *tmctypes.ResultValidators) error

	// ExportCommit accepts ResultValidators and persists validator commit signatures inside the database.
	// An error is returned if write fails.
	ExportCommit(block *tmctypes.ResultBlock, vals *tmctypes.ResultValidators) error

	// ExportEvents accepts a slice of transactions and get events in order to save in database.
	ExportEvents(ctx context.Context, block *tmctypes.ResultBlock, events *tmctypes.ResultBlockResults) error

	// HandleGenesis accepts a GenesisDoc and calls all the registered genesis handlers
	// in the order in which they have been registered.
	HandleGenesis(genesisDoc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error

	HandleBlock(block *tmctypes.ResultBlock, events *tmctypes.ResultBlockResults, txs []*types.Tx, vals *tmctypes.ResultValidators)

	// HandleMessage accepts the transaction and handles messages contained
	// inside the transaction.
	HandleMessage(block *tmctypes.ResultBlock, index int, msg sdk.Msg, tx *types.Tx)

	// HandleEvent accepts the transaction and handles events contained inside the transaction.
	HandleEvent(ctx context.Context, block *tmctypes.ResultBlock, txHash common.Hash, event sdk.Event)

	// ExportEpoch accepts a finalized block height and block hash then inside the database.
	ExportEpoch(block *tmctypes.ResultBlock) error
}

func DefaultIndexer(codec codec.Codec, proxy node.Node, db database.Database, modules []modules.Module) Indexer {
	return &Impl{
		Ctx:     context.TODO(),
		Codec:   codec,
		Node:    proxy,
		DB:      db,
		Modules: modules,
	}
}

type Impl struct {
	Ctx context.Context

	Modules []modules.Module

	Codec codec.Codec

	Node node.Node
	DB   database.Database
}

func (i *Impl) ExportEpoch(block *tmctypes.ResultBlock) error {
	return nil
}

func (i *Impl) HandleGenesis(genesisDoc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	// Call the genesis handlers
	for _, module := range i.Modules {
		if genesisModule, ok := module.(modules.GenesisModule); ok {
			if err := genesisModule.HandleGenesis(genesisDoc, appState); err != nil {
				log.Errorw("error while handling genesis", "module", module, "err", err)
			}
		}
	}

	return nil
}

func (i *Impl) HandleBlock(block *tmctypes.ResultBlock, events *tmctypes.ResultBlockResults, txs []*types.Tx, vals *tmctypes.ResultValidators) {
	for _, module := range i.Modules {
		if blockModule, ok := module.(modules.BlockModule); ok {
			err := blockModule.HandleBlock(block, events, txs, vals)
			if err != nil {
				log.Errorw("error while handling block", "module", module.Name(), "height", block.Block.Height, "err", err)
			}
		}
	}
}

func (i *Impl) HandleMessage(block *tmctypes.ResultBlock, index int, msg sdk.Msg, tx *types.Tx) {
	// Allow modules to handle the message
	for _, module := range i.Modules {
		if messageModule, ok := module.(modules.MessageModule); ok {
			err := messageModule.HandleMsg(block, index, msg, tx)
			if err != nil {
				log.Errorw("error while handling message", "module", module, "height", tx.Height,
					"txHash", tx.TxHash, "msg", proto.MessageName(msg), "err", err)
			}
		}
	}
}

func (i *Impl) HandleMsgExec(block *tmctypes.ResultBlock, index int, msg sdk.Msg, tx *types.Tx) {
	//If it's a MsgExecute, we need to make sure the included messages are handled as well
	if msgExec, ok := msg.(*authz.MsgExec); ok {
		for authzIndex, msgAny := range msgExec.Msgs {
			var executedMsg sdk.Msg
			err := i.Codec.UnpackAny(msgAny, &executedMsg)
			if err != nil {
				log.Errorw("unable to unpack MsgExec inner message", "index", authzIndex, "error", err)
			}

			for _, module := range i.Modules {
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
func (i *Impl) HandleEvent(ctx context.Context, block *tmctypes.ResultBlock, txHash common.Hash, event sdk.Event) {
	br := syncutils.NewBatchRunner()

	for _, module := range i.Modules {
		module := module
		br.AddTasks(func() error {
			if eventModule, ok := module.(modules.EventModule); ok {
				err := eventModule.HandleEvent(ctx, block, txHash, event)
				if err != nil {
					log.Errorw("failed to handle event", "module", module.Name(), "event", event, "error", err)
				}
			}
			return nil
		})
	}

	err := br.Exec()
	if err != nil {
		log.Errorw("failed to handle event in batch runner", "error", err)
	}
}

// Process fetches a block for a given height and associated metadata and export it to a database.
// It returns an error if any export process fails.
func (i *Impl) Process(height uint64) error {
	log.Debugw("processing block", "height", height)

	start := time.Now()

	block, err := i.Node.Block(int64(height))
	if err != nil {
		return fmt.Errorf("failed to get block from node: %s", err)
	}

	log.WorkerLatencyHist.Observe(float64(time.Since(block.Block.Time).Milliseconds()))

	txs, err := i.Node.Txs(block)
	if err != nil {
		return fmt.Errorf("failed to get transactions for block: %s", err)
	}

	log.IndexerLatencyHist.WithLabelValues("rpc").Observe(float64(time.Since(start).Milliseconds()))

	err = syncutils.BatchRun(
		func() error {
			return i.ProcessBlock(block, txs)
		},
		func() error {
			return i.ProcessTxs(block, txs)
		},
		func() error {
			return i.ProcessEvents(block, txs)
		},
	)

	log.DBLatencyHist.Observe(float64(time.Since(block.Block.Time).Milliseconds()))
	log.IndexerLatencyHist.WithLabelValues("processing").Observe(float64(time.Since(start).Milliseconds()))

	if err != nil {
		return err
	}

	return i.DB.SaveEpoch(context.TODO(), &models.Epoch{BlockHeight: int64(height), BlockHash: common.HexToHash(block.Block.Hash().String())})
}

func (i *Impl) ProcessBlock(block *tmctypes.ResultBlock, txs []*types.Tx) error {
	todo := false
	for _, module := range i.Modules {
		if module.Name() == (&blockmodule.Module{}).Name() {
			todo = true
			break
		}
	}
	if !todo {
		return nil
	}

	defer func(start time.Time) {
		log.IndexerLatencyHist.WithLabelValues("export_block").Observe(float64(time.Since(start).Milliseconds()))
	}(time.Now())

	return i.ExportBlock(block, nil, txs, nil)
}

func (i *Impl) ProcessTxs(block *tmctypes.ResultBlock, txs []*types.Tx) error {
	todo := false
	for _, module := range i.Modules {
		if module.Name() == (&blockmodule.Module{}).Name() {
			todo = true
			break
		}
	}
	if !todo {
		return nil
	}

	defer func(start time.Time) {
		log.IndexerLatencyHist.WithLabelValues("export_txs").Observe(float64(time.Since(start).Milliseconds()))
	}(time.Now())

	return i.ExportTxs(block, txs)
}

func (i *Impl) ProcessEvents(block *tmctypes.ResultBlock, txs []*types.Tx) error {
	defer func(start time.Time) {
		log.IndexerLatencyHist.WithLabelValues("export_events").Observe(float64(time.Since(start).Milliseconds()))
	}(time.Now())

	return i.ExportEventsByTxs(i.Ctx, block, txs)
}

// ExportBlock accepts a finalized block and persists then inside the database.
// An error is returned if write fails.
func (i *Impl) ExportBlock(
	block *tmctypes.ResultBlock, events *tmctypes.ResultBlockResults, txs []*types.Tx, vals *tmctypes.ResultValidators,
) error {
	// Save the block
	err := i.DB.SaveBlock(i.Ctx, models.NewBlockFromTmBlock(block, SumGasTxs(txs)))
	if err != nil {
		return fmt.Errorf("failed to persist block: %s", err)
	}

	return nil
}

func (i *Impl) ExportValidators(block *tmctypes.ResultBlock, vals *tmctypes.ResultValidators) error {
	var validators = make([]*models.Validator, len(vals.Validators))
	for index, val := range vals.Validators {
		consAddr := sdk.ConsAddress(val.Address).String()

		validators[index] = models.NewValidator(common.HexToAddress(consAddr), models.BytesToPubkey(val.PubKey.Bytes()))
	}

	err := i.DB.SaveValidators(context.TODO(), validators)
	if err != nil {
		return fmt.Errorf("error while saving validators: %s", err)
	}

	// Make sure the proposer exists
	proposerAddr := sdk.ConsAddress(block.Block.ProposerAddress)
	val := FindValidatorByAddr(proposerAddr.String(), vals)
	if val == nil {
		return fmt.Errorf("failed to find validator by proposer address %s: %s", proposerAddr.String(), err)
	}

	return nil
}

// ExportCommit accepts a block commitment and a corresponding set of
// validators for the commitment and persists them to the database. An error is
// returned if any write fails or if there is any missed aggregated data.
func (i *Impl) ExportCommit(block *tmctypes.ResultBlock, vals *tmctypes.ResultValidators) error {
	commit := block.Block.LastCommit

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

	err := i.DB.SaveCommitSignatures(context.TODO(), signatures)
	if err != nil {
		return fmt.Errorf("error while saving commit signatures: %s", err)
	}

	return nil
}

// ExportTxs accepts a slice of transactions and persists then inside the database.
// An error is returned if write fails.
func (i *Impl) ExportTxs(block *tmctypes.ResultBlock, txs []*types.Tx) error {
	// handle all transactions inside the block
	br := syncutils.NewBatchRunner()

	for ind, tx := range txs {
		ind, tx := ind, tx
		br.AddTasks(func() error {
			// save the transaction
			err := i.DB.SaveTx(context.TODO(), uint64(block.Block.Time.UTC().Unix()), ind, tx)
			if err != nil {
				return fmt.Errorf("error while storing todo with hash %s, %s", tx.TxHash, err)
			}

			// handle all messages contained inside the transaction
			sdkMsgs := make([]sdk.Msg, len(tx.Body.Messages))
			for ind, msg := range tx.Body.Messages {
				var stdMsg sdk.Msg
				err := i.Codec.UnpackAny(msg, &stdMsg)
				if err != nil {
					return err
				}
				sdkMsgs[ind] = stdMsg
			}

			// call the msg handlers
			for ind, sdkMsg := range sdkMsgs {
				i.HandleMessage(block, ind, sdkMsg, tx)
			}
			return nil
		})
	}

	err := br.Exec()
	if err != nil {
		return err
	}

	return nil
}

func (i *Impl) ExportEvents(ctx context.Context, block *tmctypes.ResultBlock, blockResults *tmctypes.ResultBlockResults) error {
	txsResults := blockResults.TxsResults

	for _, tx := range txsResults {
		for _, event := range tx.Events {
			i.HandleEvent(ctx, block, common.Hash{}, sdk.Event(event))
		}
	}
	return nil
}

func (i *Impl) ExportEventsByTxs(ctx context.Context, block *tmctypes.ResultBlock, txs []*types.Tx) error {
	for _, tx := range txs {
		if tx.Successful() {
			for _, event := range tx.Events {
				i.HandleEvent(ctx, block, common.HexToHash(tx.TxHash), sdk.Event(event))
			}
		}
	}
	return nil
}

// Processed tells whether the current Indexer has already processed the given height of Block
// An error is returned if the operation fails.
func (i *Impl) Processed(ctx context.Context, height uint64) (bool, error) {
	epoch, err := i.DB.GetEpoch(ctx)
	if err != nil {
		return false, err
	}

	if epoch.BlockHeight == 0 && epoch.BlockHash == common.EmptyHash {
		return false, nil
	}

	return uint64(epoch.BlockHeight) >= height, nil
}

func (i *Impl) GetLatestHeight(ctx context.Context) uint64 {
	ep, err := i.DB.GetEpoch(context.Background())
	if err != nil {
		return 0
	}
	return uint64(ep.BlockHeight)
}
