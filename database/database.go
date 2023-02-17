package database

import (
	"context"

	"github.com/cosmos/cosmos-sdk/simapp/params"
	databaseconfig "github.com/forbole/juno/v4/database/config"
	"github.com/forbole/juno/v4/types"
)

// Database represents an abstract database that can be used to save data inside it
type Database interface {
	// PrepareTables create tables
	PrepareTables(ctx context.Context) error

	// HasBlock tells whether the database has already stored the block having the given height.
	// An error is returned if the operation fails.
	HasBlock(height int64) (bool, error)

	// GetLastBlockHeight returns the last block height stored in database..
	// An error is returned if the operation fails.
	GetLastBlockHeight() (int64, error)

	// GetMissingHeights returns a slice of missing block heights between startHeight and endHeight
	GetMissingHeights(startHeight, endHeight int64) []int64

	// SaveBlock will be called when a new block is parsed, passing the block itself
	// and the transactions contained inside that block.
	// An error is returned if the operation fails.
	// NOTE. For each transaction inside txs, SaveTx will be called as well.
	SaveBlock(block *types.Block) error

	// GetTotalBlocks returns total number of blocks stored in database.
	GetTotalBlocks() int64

	// SaveTx will be called to save each transaction contained inside a block.
	// An error is returned if the operation fails.
	SaveTx(tx *types.Tx) error

	// HasValidator returns true if a given validator by consensus address exists.
	// An error is returned if the operation fails.
	HasValidator(address string) (bool, error)

	// SaveValidators stores a list of validators if they do not already exist.
	// An error is returned if the operation fails.
	SaveValidators(validators []*types.Validator) error

	// SaveCommitSignatures stores a  slice of validator commit signatures.
	// An error is returned if the operation fails.
	SaveCommitSignatures(signatures []*types.CommitSig) error

	// SaveMessage stores a single message.
	// An error is returned if the operation fails.
	SaveMessage(msg *types.Message) error

	// Close closes the connection to the database
	Close()
}

// PruningDb represents a database that supports pruning properly
type PruningDb interface {
	// Prune prunes the data for the given height, returning any error
	Prune(height int64) error

	// StoreLastPruned saves the last height at which the database was pruned
	StoreLastPruned(height int64) error

	// GetLastPruned returns the last height at which the database was pruned
	GetLastPruned() (int64, error)
}

// Context contains the data that might be used to build a Database instance
type Context struct {
	Cfg            databaseconfig.Config
	EncodingConfig *params.EncodingConfig
}

// NewContext allows to build a new Context instance
func NewContext(cfg databaseconfig.Config, encodingConfig *params.EncodingConfig) *Context {
	return &Context{
		Cfg:            cfg,
		EncodingConfig: encodingConfig,
	}
}

// Builder represents a method that allows to build any database from a given codec and configuration
type Builder func(ctx *Context) (Database, error)
