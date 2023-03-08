package parser

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

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

	// ExportBlock accepts a finalized block and persists then inside the database.
	// An error is returned if write fails.
	ExportBlock(block *tmctypes.ResultBlock, events *tmctypes.ResultBlockResults, txs []*types.Tx, vals *tmctypes.ResultValidators) error

	// ExportTxs accepts a slice of transactions and persists then inside the database.
	// An error is returned if write fails.
	ExportTxs(txs []*types.Tx) error

	// ExportValidators accepts ResultValidators and persists validators inside the database.
	// An error is returned if write fails.
	ExportValidators(block *tmctypes.ResultBlock, vals *tmctypes.ResultValidators) error

	// ExportCommit accepts ResultValidators and persists validator commit signatures inside the database.
	// An error is returned if write fails.
	ExportCommit(block *tmctypes.ResultBlock, vals *tmctypes.ResultValidators) error

	// ExportAccounts accepts a slice of transactions and persists accounts inside the database.
	// An error is returned if write fails.
	ExportAccounts(block *tmctypes.ResultBlock, txs []*types.Tx) error

	// ExportEvents accepts a slice of transactions and get events in order to save in database.
	ExportEvents(events *tmctypes.ResultBlockResults) error
}

func DefaultIndexer(codec codec.Codec, proxy node.Node, db database.Database, modules []modules.Module) Indexer {
	return &Impl{
		codec:  codec,
		Node:   proxy,
		DB:     db,
		Worker: NewPuppetWorker(modules),
	}
}

type Impl struct {
	Ctx context.Context

	Worker *Worker

	codec codec.Codec

	Node node.Node
	DB   database.Database
}

// Process fetches a block for a given height and associated metadata and export it to a database.
// It returns an error if any export process fails.
func (i *Impl) Process(height uint64) error {
	log.Debugw("processing block", "height", height)

	block, err := i.Node.Block(int64(height))
	if err != nil {
		return fmt.Errorf("failed to get block from node: %s", err)
	}

	events, err := i.Node.BlockResults(int64(height))
	if err != nil {
		return fmt.Errorf("failed to get block results from node: %s", err)
	}

	txs, err := i.Node.Txs(block)
	if err != nil {
		return fmt.Errorf("failed to get transactions for block: %s", err)
	}

	vals, err := i.Node.Validators(int64(height))
	if err != nil {
		return fmt.Errorf("failed to get validators for block: %s", err)
	}

	err = i.ExportBlock(block, events, txs, vals)
	if err != nil {
		return err
	}

	err = i.ExportValidators(block, vals)
	if err != nil {
		return err
	}

	err = i.ExportTxs(txs)
	if err != nil {
		return err
	}

	err = i.ExportAccounts(block, txs)
	if err != nil {
		return err
	}

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

	// Call the block handlers
	if i.Worker != nil {
		log.Debugw("puppet worker exists, handle block...")
		i.Worker.HandleBlock(block, events, txs, vals)
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
func (i *Impl) ExportTxs(txs []*types.Tx) error {
	// handle all transactions inside the block
	for _, tx := range txs {
		// save the transaction
		err := i.DB.SaveTx(context.TODO(), tx)
		if err != nil {
			return fmt.Errorf("error while storing tx with hash %s, %s", tx.TxHash, err)
		}

		// call the tx handlers
		if i.Worker != nil {
			log.Debugw("puppet worker exists, handle tx...")
			i.Worker.HandleTx(tx)
		}

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
		if i.Worker != nil {
			log.Debugw("puppet worker exists, handle msg...")
			for ind, sdkMsg := range sdkMsgs {
				i.Worker.HandleMessage(ind, sdkMsg, tx)
			}
		}
	}

	totalBlocks := i.DB.GetTotalBlocks(context.TODO())
	log.DbBlockCount.WithLabelValues("total_blocks_in_db").Set(float64(totalBlocks))

	dbLatestHeight, err := i.DB.GetLastBlockHeight(context.TODO())
	if err != nil {
		return err
	}
	log.DbLatestHeight.WithLabelValues("db_latest_height").Set(float64(dbLatestHeight))

	return nil
}

// ExportAccounts accepts a slice of transactions and persists accounts inside the database.
// An error is returned if write fails.
func (i *Impl) ExportAccounts(block *tmctypes.ResultBlock, txs []*types.Tx) error {
	// save account
	for _, tx := range txs {
		for _, l := range tx.Logs {
			for _, event := range l.Events {
				for _, attr := range event.Attributes {
					if common.IsHexAddress(attr.Value) {
						account := &models.Account{
							Address:             common.HexToAddress(attr.Value),
							TxCount:             1,
							LastActiveTimestamp: uint64(block.Block.Time.Unix()),
						}
						err := i.DB.SaveAccount(context.TODO(), account)
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

func (i *Impl) ExportEvents(events *tmctypes.ResultBlockResults) error {
	return nil
}
