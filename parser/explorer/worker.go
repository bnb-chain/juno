package explorer

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/forbole/juno/v4/log"
	"github.com/forbole/juno/v4/models"
	"github.com/forbole/juno/v4/modules"
	"github.com/forbole/juno/v4/parser"
	"github.com/forbole/juno/v4/types"
	"github.com/forbole/juno/v4/types/config"
	"github.com/forbole/juno/v4/types/utils"
	"github.com/forbole/juno/v4/utils/syncutils"
)

// Worker defines a job consumer that is responsible for getting and
// aggregating block and associated data and exporting it to a database.
type Worker struct {
	parser.CommonWorker
}

// Process fetches a block for a given height and associated metadata and export it to a database.
// It returns an error if any export process fails.
func (w *Worker) Process(height uint64) error {
	log.Debugw("processing block", "height", height)

	if height == 0 {
		cfg := config.Cfg.Parser

		genesisDoc, genesisState, err := utils.GetGenesisDocAndState(cfg.GenesisFilePath, w.Node)
		if err != nil {
			return fmt.Errorf("failed to get genesis: %s", err)
		}

		return w.HandleGenesis(genesisDoc, genesisState)
	}

	block, err := w.Node.Block(int64(height))
	if err != nil {
		return fmt.Errorf("failed to get block from node: %s", err)
	}

	events, err := w.Node.BlockResults(int64(height))
	if err != nil {
		return fmt.Errorf("failed to get block results from node: %s", err)
	}

	txs, err := w.Node.Txs(block)
	if err != nil {
		return fmt.Errorf("failed to get transactions for block: %s", err)
	}

	vals, err := w.Node.Validators(int64(height))
	if err != nil {
		return fmt.Errorf("failed to get validators for block: %s", err)
	}

	return w.ExportBlock(block, events, txs, vals)
}

// ExportBlock accepts a finalized block and a corresponding set of transactions and persists them to the database along with attributable metadata. An error is returned if write fails.
func (w *Worker) ExportBlock(args ...interface{}) error {
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
	err := w.SaveValidators(vals.Validators)
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
	err = w.DB.SaveBlock(w.Ctx, models.NewBlockFromTmBlock(b, parser.SumGasTxs(txs)))
	if err != nil {
		return fmt.Errorf("failed to persist block: %s", err)
	}

	//currently no need
	// Save the commits
	err = w.ExportCommit(b.Block.LastCommit, vals)
	if err != nil {
		return err
	}

	// Call the block handlers
	for _, module := range w.Modules {
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
