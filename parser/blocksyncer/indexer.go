package blocksyncer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/forbole/juno/v4/log"
	"github.com/forbole/juno/v4/models"
	"github.com/forbole/juno/v4/modules"
	"github.com/forbole/juno/v4/parser"
	"github.com/forbole/juno/v4/types"
	eventutil "github.com/forbole/juno/v4/types/event"
)

type Indexer struct {
	Ctx context.Context
	parser.CommonIndexer
}

// ProcessBlock for blocksyncer worker only process basic block data and events
func (idx *Indexer) ProcessBlock(height uint64) error {
	log.Debugw("processing block", "height", height)

	block, err := idx.Node.Block(int64(height))
	if err != nil {
		return fmt.Errorf("failed to get block from node: %s", err)
	}

	events, err := idx.Node.BlockResults(int64(height))
	if err != nil {
		return fmt.Errorf("failed to get block results from node: %s", err)
	}

	return idx.ExportBlock(block, events)
}

// ExportBlock for blocksyncer worker only exports basic block data and module related events
func (idx *Indexer) ExportBlock(args ...interface{}) error {
	var (
		b *tmctypes.ResultBlock
		r *tmctypes.ResultBlockResults
	)

	for _, arg := range args {
		switch arg.(type) {
		case *tmctypes.ResultBlock:
			b = arg.(*tmctypes.ResultBlock)
		case *tmctypes.ResultBlockResults:
			r = arg.(*tmctypes.ResultBlockResults)
		default:
			return errors.New("block result type not supported")
		}
	}

	// Save the block simple
	err := idx.DB.SaveBlockLight(context.TODO(), models.NewBlockFromTmBlock(b, 0))
	if err != nil {
		return fmt.Errorf("failed to persist block: %s", err)
	}

	// Call the event handlers
	return idx.ExportEvents(idx.Ctx, b, r.TxsResults)

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

// ExportEvents accepts a slice of transactions and get events in order to save in database.
func (idx *Indexer) ExportEvents(ctx context.Context, block *tmctypes.ResultBlock, txs []*abci.ResponseDeliverTx) error {
	// get all events in order from the txs within the block
	for _, tx := range txs {
		// handle all events contained inside the transaction
		events := filterEventsType(tx)
		// call the event handlers
		for i, event := range events {
			idx.HandleEvent(ctx, block, i, event)
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

// ExportTxs accepts a slice of transactions and persists then inside the database.
// An error is returned if write fails.
func (idx *Indexer) ExportTxs(txs []*types.Tx) error {
	return nil
}

// ExportAccounts accepts a slice of transactions and persists accounts inside the database.
// An error is returned if write fails.
func (idx *Indexer) ExportAccounts(txs []*types.Tx) error {
	return nil
}

// ProcessTransactions fetches transactions for a given height and stores them into the database.
// It returns an error if the export process fails.
func (idx *Indexer) ProcessTransactions(height int64) error {
	return nil
}

// HandleGenesis accepts a GenesisDoc and calls all the registered genesis handlers
// in the order in which they have been registered.
func (idx *Indexer) HandleGenesis(genesisDoc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	return nil
}

// SaveValidators persists a list of Tendermint validators with an address and a
// consensus public key. An error is returned if the public key cannot be Bech32
// encoded or if the DB write fails.
func (idx *Indexer) SaveValidators(vals []*tmtypes.Validator) error {
	return nil
}

// ExportCommit accepts a block commitment and a corresponding set of
// validators for the commitment and persists them to the database. An error is
// returned if any write fails or if there is any missed aggregated data.
func (idx *Indexer) ExportCommit(commit *tmtypes.Commit, vals *tmctypes.ResultValidators) error {
	return nil
}

// SaveTx accepts the transaction and persists it inside the database.
// An error is returned if write fails.
func (idx *Indexer) SaveTx(tx *types.Tx) error {
	return nil
}

// HandleTx accepts the transaction and calls the tx handlers.
func (idx *Indexer) HandleTx(tx *types.Tx) {
}

// HandleMessage accepts the transaction and handles messages contained
// inside the transaction.
func (idx *Indexer) HandleMessage(index int, msg sdk.Msg, tx *types.Tx) {
}
