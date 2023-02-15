package mysql

import (
	"github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/forbole/juno/v4/database"
	"github.com/forbole/juno/v4/logging"
	"github.com/forbole/juno/v4/types"
	"gorm.io/gorm"
)

// Builder creates a database connection with the given database connection info
// from config. It returns a database connection handle or an error if the
// connection fails.
func Builder(ctx *database.Context) (database.Database, error) {

	return &Database{}, nil
}

// Database defines a wrapper around a SQL database and implements functionality
// for data aggregation and exporting.
type Database struct {
	db             *gorm.DB
	EncodingConfig *params.EncodingConfig
	Logger         logging.Logger
}

// HasBlock implements database.Database
func (db *Database) HasBlock(height int64) (bool, error) {
	var res bool
	return res, nil
}

// GetLastBlockHeight returns the last block height stored inside the database
func (db *Database) GetLastBlockHeight() (int64, error) {
	var height int64
	return height, nil
}

// GetMissingHeights returns a slice of missing block heights between startHeight and endHeight
func (db *Database) GetMissingHeights(startHeight, endHeight int64) []int64 {
	var result []int64
	if len(result) == 0 {
		return nil
	}

	return result
}

// SaveBlock implements database.Database
func (db *Database) SaveBlock(block *types.Block) error {

	return nil
}

// GetTotalBlocks implements database.Database
func (db *Database) GetTotalBlocks() int64 {
	var blockCount int64

	return blockCount
}

// SaveTx implements database.Database
func (db *Database) SaveTx(tx *types.Tx) error {
	return nil
}

// HasValidator implements database.Database
func (db *Database) HasValidator(addr string) (bool, error) {
	var res bool

	return res, nil
}

// SaveValidators implements database.Database
func (db *Database) SaveValidators(validators []*types.Validator) error {
	if len(validators) == 0 {
		return nil
	}

	return nil
}

// SaveCommitSignatures implements database.Database
func (db *Database) SaveCommitSignatures(signatures []*types.CommitSig) error {
	if len(signatures) == 0 {
		return nil
	}

	return nil
}

// SaveMessage implements database.Database
func (db *Database) SaveMessage(msg *types.Message) error {

	return nil
}

// Close implements database.Database
func (db *Database) Close() {

}
