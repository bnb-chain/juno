package explorer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/gogo/protobuf/proto"
	abci "github.com/tendermint/tendermint/abci/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/forbole/juno/v4/common"
	"github.com/forbole/juno/v4/log"
	"github.com/forbole/juno/v4/models"
	"github.com/forbole/juno/v4/modules"
	"github.com/forbole/juno/v4/parser"
	"github.com/forbole/juno/v4/types"
	"github.com/forbole/juno/v4/types/config"
	"github.com/forbole/juno/v4/types/utils"
	"github.com/forbole/juno/v4/utils/syncutils"
)

type Indexer struct {
	parser.CommonIndexer
}

// ProcessBlock fetches a block for a given height and associated metadata and export it to a database.
// It returns an error if any export process fails.
func (idx *Indexer) ProcessBlock(height uint64) error {
	log.Debugw("processing block", "height", height)

	if height == 0 {
		cfg := config.Cfg.Parser

		genesisDoc, genesisState, err := utils.GetGenesisDocAndState(cfg.GenesisFilePath, idx.Node)
		if err != nil {
			return fmt.Errorf("failed to get genesis: %s", err)
		}

		return idx.HandleGenesis(genesisDoc, genesisState)
	}

	block, err := idx.Node.Block(int64(height))
	if err != nil {
		return fmt.Errorf("failed to get block from node: %s", err)
	}

	events, err := idx.Node.BlockResults(int64(height))
	if err != nil {
		return fmt.Errorf("failed to get block results from node: %s", err)
	}

	txs, err := idx.Node.Txs(block)
	if err != nil {
		return fmt.Errorf("failed to get transactions for block: %s", err)
	}

	vals, err := idx.Node.Validators(int64(height))
	if err != nil {
		return fmt.Errorf("failed to get validators for block: %s", err)
	}

	return idx.ExportBlock(block, events, txs, vals)
}

// ExportBlock accepts a finalized block and a corresponding set of transactions and persists them to the database along with attributable metadata. An error is returned if write fails.
func (idx *Indexer) ExportBlock(args ...interface{}) error {
	var (
		b    *tmctypes.ResultBlock
		r    *tmctypes.ResultBlockResults
		txs  []*types.Tx
		vals *tmctypes.ResultValidators
	)

	for _, arg := range args {
		switch arg.(type) {
		case *tmctypes.ResultBlock:
			b = arg.(*tmctypes.ResultBlock)
		case *tmctypes.ResultBlockResults:
			r = arg.(*tmctypes.ResultBlockResults)
		case []*types.Tx:
			txs = arg.([]*types.Tx)
		case *tmctypes.ResultValidators:
			vals = arg.(*tmctypes.ResultValidators)
		default:
			return errors.New("block result type not supported")
		}
	}

	// Save all validators
	err := idx.SaveValidators(vals.Validators)
	if err != nil {
		return err
	}

	// Make sure the proposer exists
	proposerAddr := sdk.ConsAddress(b.Block.ProposerAddress)
	val := parser.FindValidatorByAddr(proposerAddr.String(), vals)
	if val == nil {
		return fmt.Errorf("failed to find validator by proposer address %s: %s", proposerAddr.String(), err)
	}

	// Save the block
	err = idx.DB.SaveBlock(context.TODO(), models.NewBlockFromTmBlock(b, parser.SumGasTxs(txs)))
	if err != nil {
		return fmt.Errorf("failed to persist block: %s", err)
	}

	//currently no need
	// Save the commits
	err = idx.ExportCommit(b.Block.LastCommit, vals)
	if err != nil {
		return err
	}

	// Call the block handlers
	for _, module := range idx.Modules {
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
			return idx.ExportTxs(txs)
		},
		func() error {
			return idx.ExportAccounts(txs)
		},
	)
}

// ProcessTransactions fetches transactions for a given height and stores them into the database.
// It returns an error if the export process fails.
func (idx *Indexer) ProcessTransactions(height int64) error {
	block, err := idx.Node.Block(height)
	if err != nil {
		return fmt.Errorf("failed to get block from node: %s", err)
	}

	txs, err := idx.Node.Txs(block)
	if err != nil {
		return fmt.Errorf("failed to get transactions for block: %s", err)
	}

	return syncutils.BatchRun(
		func() error {
			return idx.ExportTxs(txs)
		},
		func() error {
			return idx.ExportAccounts(txs)
		},
	)
}

// HandleGenesis accepts a GenesisDoc and calls all the registered genesis handlers
// in the order in which they have been registered.
func (idx *Indexer) HandleGenesis(genesisDoc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	// Call the genesis handlers
	for _, module := range idx.Modules {
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
func (idx *Indexer) SaveValidators(vals []*tmtypes.Validator) error {
	var validators = make([]*models.Validator, len(vals))
	for index, val := range vals {
		consAddr := sdk.ConsAddress(val.Address).String()

		validators[index] = models.NewValidator(common.HexToAddress(consAddr), models.BytesToPubkey(val.PubKey.Bytes()))
	}

	err := idx.DB.SaveValidators(context.TODO(), validators)
	if err != nil {
		return fmt.Errorf("error while saving validators: %s", err)
	}

	return nil
}

// ExportCommit accepts a block commitment and a corresponding set of
// validators for the commitment and persists them to the database. An error is
// returned if any write fails or if there is any missed aggregated data.
func (idx *Indexer) ExportCommit(commit *tmtypes.Commit, vals *tmctypes.ResultValidators) error {
	var signatures []*types.CommitSig
	for _, commitSig := range commit.Signatures {
		// Avoid empty commits
		if commitSig.Signature == nil {
			continue
		}

		valAddr := sdk.ConsAddress(commitSig.ValidatorAddress)
		val := parser.FindValidatorByAddr(valAddr.String(), vals)
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

	err := idx.DB.SaveCommitSignatures(context.TODO(), signatures)
	if err != nil {
		return fmt.Errorf("error while saving commit signatures: %s", err)
	}

	return nil
}

// SaveTx accepts the transaction and persists it inside the database.
// An error is returned if write fails.
func (idx *Indexer) SaveTx(tx *types.Tx) error {
	err := idx.DB.SaveTx(context.TODO(), tx)
	if err != nil {
		return fmt.Errorf("failed to handle transaction with hash %s: %s", tx.TxHash, err)
	}
	return nil
}

// HandleTx accepts the transaction and calls the tx handlers.
func (idx *Indexer) HandleTx(tx *types.Tx) {
	// Call the tx handlers
	for _, module := range idx.Modules {
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
func (idx *Indexer) HandleMessage(index int, msg sdk.Msg, tx *types.Tx) {
	// Allow modules to handle the message
	for _, module := range idx.Modules {
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
			err := idx.Codec.UnpackAny(msgAny, &executedMsg)
			if err != nil {
				log.Errorw("unable to unpack MsgExec inner message", "index", authzIndex, "error", err)
			}

			for _, module := range idx.Modules {
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
func (idx *Indexer) HandleEvent(ctx context.Context, block *tmctypes.ResultBlock, index int, event sdk.Event) {
	// Allow modules to handle the message
	for _, module := range idx.Modules {
		if eventModule, ok := module.(modules.EventModule); ok {
			err := eventModule.HandleEvent(ctx, block, index, event)
			if err != nil {
				log.Errorw("error while handling event", "module", module, "event", event, "err", err)
			}
		}
	}

}

// ExportTxs accepts a slice of transactions and persists then inside the database.
// An error is returned if write fails.
func (idx *Indexer) ExportTxs(txs []*types.Tx) error {
	// handle all transactions inside the block
	for _, tx := range txs {
		// save the transaction
		err := idx.SaveTx(tx)
		if err != nil {
			return fmt.Errorf("error while storing txs: %s", err)
		}

		// call the tx handlers
		idx.HandleTx(tx)

		// handle all messages contained inside the transaction
		sdkMsgs := make([]sdk.Msg, len(tx.Body.Messages))
		for i, msg := range tx.Body.Messages {
			var stdMsg sdk.Msg
			err := idx.Codec.UnpackAny(msg, &stdMsg)
			if err != nil {
				return err
			}
			sdkMsgs[i] = stdMsg
		}

		// call the msg handlers
		for i, sdkMsg := range sdkMsgs {
			idx.HandleMessage(i, sdkMsg, tx)
		}
	}

	totalBlocks := idx.DB.GetTotalBlocks(context.TODO())
	log.DbBlockCount.WithLabelValues("total_blocks_in_db").Set(float64(totalBlocks))

	dbLatestHeight, err := idx.DB.GetLastBlockHeight(context.TODO())
	if err != nil {
		return err
	}
	log.DbLatestHeight.WithLabelValues("db_latest_height").Set(float64(dbLatestHeight))

	return nil
}

// ExportAccounts accepts a slice of transactions and persists accounts inside the database.
// An error is returned if write fails.
func (idx *Indexer) ExportAccounts(txs []*types.Tx) error {
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
						err := idx.DB.SaveAccount(context.TODO(), account)
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
func (idx *Indexer) ExportEvents(ctx context.Context, block *tmctypes.ResultBlock, txs []*abci.ResponseDeliverTx) error {
	return nil
}
