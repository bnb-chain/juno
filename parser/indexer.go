package parser

import (
	"context"
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/forbole/juno/v4/database"
	"github.com/forbole/juno/v4/modules"
	"github.com/forbole/juno/v4/node"
	"github.com/forbole/juno/v4/types"
)

type Indexer interface {
	// ProcessBlock fetches  a block for a given height and associated metadata and export it to a database.
	// It returns an error if any export process fails.
	ProcessBlock(height uint64) error

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
	HandleEvent(ctx context.Context, block *tmctypes.ResultBlock, index int, event sdk.Event)

	// ExportTxs accepts a slice of transactions and persists then inside the database.
	// An error is returned if write fails.
	ExportTxs(txs []*types.Tx) error

	// ExportAccounts accepts a slice of transactions and persists accounts inside the database.
	// An error is returned if write fails.
	ExportAccounts(txs []*types.Tx) error

	// ExportEvents accepts a slice of transactions and get events in order to save in database.
	ExportEvents(ctx context.Context, block *tmctypes.ResultBlock, txs []*abci.ResponseDeliverTx) error

	// ChainID returns the chain id of the node which indexer owns
	ChainID() (string, error)

	// HasBlock returns if the given block height exist
	HasBlock(height uint64) (bool, error)
}

func NewCommonIndexer(ctx *Context) CommonIndexer {
	return CommonIndexer{
		//Codec:   ctx.EncodingConfig.Codec,
		Modules: ctx.Modules,
		Node:    ctx.Node,
		DB:      ctx.Database,
	}
}

type CommonIndexer struct {
	Codec   codec.Codec
	Modules []modules.Module
	Node    node.Node
	DB      database.Database
}

// ChainID returns the chain id of the node which indexer owns
func (idx *CommonIndexer) ChainID() (string, error) {
	return idx.Node.ChainID()
}

// HasBlock returns if the given block height exist
func (idx *CommonIndexer) HasBlock(height uint64) (bool, error) {
	return idx.DB.HasBlock(context.TODO(), height)
}
