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
	"github.com/forbole/juno/v4/node"
	"github.com/forbole/juno/v4/types"
)

type Indexer interface {
	// Process fetches a block for a given height and associated metadata and export it to a database.
	// It returns an error if any export process fails.
	Process(height uint64) error

	// Processed tells whether the current Indexer has already processed the given height of Block
	// An error is returned if the operation fails.
	Processed(ctx context.Context, height uint64) (bool, error)

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

	// HandleTx accepts the transaction and calls the tx handlers.
	HandleTx(tx *types.Tx)

	// HandleMessage accepts the transaction and handles messages contained
	// inside the transaction.
	HandleMessage(block *tmctypes.ResultBlock, index int, msg sdk.Msg, tx *types.Tx)

	// HandleEvent accepts the transaction and handles events contained inside the transaction.
	HandleEvent(ctx context.Context, block *tmctypes.ResultBlock, txHash common.Hash, event sdk.Event) error

	// ExportEpoch accepts a finalized block height and block hash then inside the database.
	ExportEpoch(ctx context.Context, block *tmctypes.ResultBlock) error

	// GetBlockRecordNum returns total number of blocks stored in database.
	GetBlockRecordNum(ctx context.Context) int64

	// GetLastBlockRecordHeight returns the last block height stored inside the database
	GetLastBlockRecordHeight(ctx context.Context) (uint64, error)
}

func DefaultIndexer(codec codec.Codec, proxy node.Node, db database.Database, modules []modules.Module) Indexer {
	return &Impl{
		Ctx:     context.TODO(),
		codec:   codec,
		Node:    proxy,
		DB:      db,
		Modules: modules,
	}
}

type Impl struct {
	Ctx context.Context

	Modules []modules.Module

	codec codec.Codec

	Node node.Node
	DB   database.Database
}

func (i *Impl) ExportEpoch(ctx context.Context, block *tmctypes.ResultBlock) error {
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

func (i *Impl) HandleTx(tx *types.Tx) {
	// Call the tx handlers
	for _, module := range i.Modules {
		if transactionModule, ok := module.(modules.TransactionModule); ok {
			err := transactionModule.HandleTx(tx)
			if err != nil {
				log.Errorw("error while handling transaction", "module", module.Name(), "height", tx.Height,
					"txHash", tx.TxHash, "err", err)
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

	// If it's a MsgExecute, we need to make sure the included messages are handled as well
	if msgExec, ok := msg.(*authz.MsgExec); ok {
		for authzIndex, msgAny := range msgExec.Msgs {
			var executedMsg sdk.Msg
			err := i.codec.UnpackAny(msgAny, &executedMsg)
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
func (i *Impl) HandleEvent(ctx context.Context, block *tmctypes.ResultBlock, txHash common.Hash, event sdk.Event) error {
	for _, module := range i.Modules {
		if eventModule, ok := module.(modules.EventModule); ok {
			err := eventModule.HandleEvent(ctx, block, txHash, event)
			if err != nil {
				log.Errorw("failed to handle event", "module", module.Name(), "event", event, "error", err)
				return err
			}
		}
	}
	return nil
}

// Process fetches a block for a given height and associated metadata and export it to a database.
// It returns an error if any export process fails.
func (i *Impl) Process(height uint64) error {
	log.Debugw("processing block", "height", height)

	block, err := i.Node.Block(int64(height))
	if err != nil {
		return fmt.Errorf("failed to get block from node: %s", err)
	}

	log.WorkerLatencyHist.Observe(float64(time.Since(block.Block.Time).Milliseconds()))

	blockResults, err := i.Node.BlockResults(int64(height))
	if err != nil {
		return fmt.Errorf("failed to get block results from node: %s", err)
	}

	txs, err := i.Node.Txs(block)
	if err != nil {
		return fmt.Errorf("failed to get transactions for block: %s", err)
	}

	err = i.ExportBlock(block, blockResults, txs, nil)
	if err != nil {
		return err
	}

	err = i.ExportTxs(block, txs)
	if err != nil {
		return err
	}

	err = i.ExportEventsByTxs(i.Ctx, block, txs)
	if err != nil {
		return err
	}

	log.DBLatencyHist.Observe(float64(time.Since(block.Block.Time).Milliseconds()))

	return nil
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

	i.HandleBlock(block, events, txs, vals)

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
	for ind, tx := range txs {
		// save the transaction
		err := i.DB.SaveTx(context.TODO(), uint64(block.Block.Time.UTC().UnixNano()), ind, tx)
		if err != nil {
			return fmt.Errorf("error while storing tx with hash %s, %s", tx.TxHash, err)
		}

		// call the tx handlers
		i.HandleTx(tx)

		// handle all messages contained inside the transaction
		sdkMsgs := make([]sdk.Msg, len(tx.Body.Messages))
		for ind, msg := range tx.Body.Messages {
			var stdMsg sdk.Msg
			err := i.codec.UnpackAny(msg, &stdMsg)
			if err != nil {
				return err
			}
			sdkMsgs[ind] = stdMsg
		}

		// call the msg handlers
		for ind, sdkMsg := range sdkMsgs {
			i.HandleMessage(block, ind, sdkMsg, tx)
		}
	}

	return nil
}

func (i *Impl) ExportEvents(ctx context.Context, block *tmctypes.ResultBlock, blockResults *tmctypes.ResultBlockResults) error {
	txsResults := blockResults.TxsResults

	for _, tx := range txsResults {
		for _, event := range tx.Events {
			if err := i.HandleEvent(ctx, block, common.Hash{}, sdk.Event(event)); err != nil {
				return err
			}
		}
	}
	return nil
}

func (i *Impl) ExportEventsByTxs(ctx context.Context, block *tmctypes.ResultBlock, txs []*types.Tx) error {
	for _, tx := range txs {
		txHash := common.HexToHash(tx.TxHash)
		for _, event := range tx.Events {
			if err := i.HandleEvent(ctx, block, txHash, sdk.Event(event)); err != nil {
				return err
			}
		}
	}
	return nil
}

// Processed tells whether the current Indexer has already processed the given height of Block
// An error is returned if the operation fails.
func (i *Impl) Processed(ctx context.Context, height uint64) (bool, error) {
	return i.DB.HasBlock(ctx, height)
}

// GetBlockRecordNum returns total number of blocks stored in database.
func (i *Impl) GetBlockRecordNum(ctx context.Context) int64 {
	return i.DB.GetTotalBlocks(ctx)
}

// GetLastBlockRecordHeight returns the last block height stored inside the database
func (i *Impl) GetLastBlockRecordHeight(ctx context.Context) (uint64, error) {
	return i.DB.GetLastBlockHeight(ctx)
}
