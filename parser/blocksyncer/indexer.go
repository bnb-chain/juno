package blocksyncer

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"

	"github.com/forbole/juno/v4/log"
	"github.com/forbole/juno/v4/models"
	"github.com/forbole/juno/v4/parser"
	eventutil "github.com/forbole/juno/v4/types/event"
)

// Indexer TODO krish move this file to SP repo
type Indexer struct {
	parser.Impl
}

// Process only process basic block data and events
func (idx *Indexer) Process(height uint64) error {
	log.Debugw("processing block", "height", height)

	block, err := idx.Node.Block(int64(height))
	if err != nil {
		return fmt.Errorf("failed to get block from node: %s", err)
	}

	events, err := idx.Node.BlockResults(int64(height))
	if err != nil {
		return fmt.Errorf("failed to get block results from node: %s", err)
	}

	err = idx.ExportBlock(block, nil, nil, nil)
	if err != nil {
		return err
	}

	err = idx.ExportEvents(events)
	if err != nil {
		return err
	}

	return nil
}

// ExportBlock accepts a finalized block and persists then inside the database.
// An error is returned if write fails.
func (idx *Indexer) ExportBlock(
	block *tmctypes.ResultBlock, _ *tmctypes.ResultBlockResults, _ []*types.Tx, _ *tmctypes.ResultValidators,
) error {
	// Save the block
	err := idx.DB.SaveBlock(idx.Ctx, models.NewBlockFromTmBlock(block, 0))
	if err != nil {
		return fmt.Errorf("failed to persist block: %s", err)
	}

	return nil
}

// ExportEvents accepts a slice of transactions and get events in order to save in database.
func (idx *Indexer) ExportEvents(events *tmctypes.ResultBlockResults) error {
	if idx.Worker == nil {
		return nil
	}

	log.Debugw("puppet worker exists, handle events...")

	txsResults := events.TxsResults
	// get all events in order from the txsResults within the block
	for _, tx := range txsResults {
		// handle all events contained inside the transaction
		events := filterEventsType(tx)
		// call the event handlers
		for i, event := range events {
			idx.Worker.HandleEvent(i, event)
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
