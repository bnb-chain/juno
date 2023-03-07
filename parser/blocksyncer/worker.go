package blocksyncer

import (
	"errors"
	"fmt"

	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/forbole/juno/v4/log"
	"github.com/forbole/juno/v4/models"
	"github.com/forbole/juno/v4/parser"
)

// Worker defines a job consumer that is responsible for getting and
// aggregating block and associated data and exporting it to a database.
type Worker struct {
	parser.CommonWorker
}

// ExportBlock for blocksyncer worker only exports basic block data and module related events
func (w *Worker) ExportBlock(args ...interface{}) error {
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
	err := w.DB.SaveBlockLight(w.Ctx, models.NewBlockFromTmBlock(b, 0))
	if err != nil {
		return fmt.Errorf("failed to persist block: %s", err)
	}

	// Call the event handlers
	return w.ExportEvents(b, r.TxsResults)

}

// Process for blocksyncer worker only process basic block data and events
func (w *Worker) Process(height uint64) error {
	log.Debugw("processing block", "height", height)

	block, err := w.Node.Block(int64(height))
	if err != nil {
		return fmt.Errorf("failed to get block from node: %s", err)
	}

	events, err := w.Node.BlockResults(int64(height))
	if err != nil {
		return fmt.Errorf("failed to get block results from node: %s", err)
	}

	return w.ExportBlock(block, events)
}
