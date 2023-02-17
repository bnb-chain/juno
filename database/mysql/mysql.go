package mysql

import (
	"context"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/forbole/juno/v4/database"
	"github.com/forbole/juno/v4/log"
	"github.com/forbole/juno/v4/logging"
	"github.com/forbole/juno/v4/models"
	"github.com/forbole/juno/v4/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Builder creates a database connection with the given database connection info
// from config. It returns a database connection handle or an error if the
// connection fails.
func Builder(ctx *database.Context) (database.Database, error) {
	db, err := New(&ctx.Cfg)
	if err != nil {
		return nil, err
	}
	return &Database{
		db:             db,
		EncodingConfig: ctx.EncodingConfig,
		Logger:         ctx.Logger,
	}, nil
}

// Database defines a wrapper around a SQL database and implements functionality
// for data aggregation and exporting.
type Database struct {
	db             *gorm.DB
	EncodingConfig *params.EncodingConfig
	Logger         logging.Logger
}

func (db *Database) PrepareTables(ctx context.Context) error {

	//TODO load by modules

	db.db.Migrator().AutoMigrate(&models.Account{})
	db.db.Migrator().AutoMigrate(&models.Block{})
	db.db.Migrator().AutoMigrate(&models.Tx{})

	//block_syncer tables
	db.db.Migrator().AutoMigrate(&models.Bucket{})
	db.db.Migrator().AutoMigrate(&models.Group{})
	db.db.Migrator().AutoMigrate(&models.Object{})

	//validator tables
	db.db.Migrator().AutoMigrate(&models.Validator{})
	db.db.Migrator().AutoMigrate(&models.ValidatorInfo{})
	db.db.Migrator().AutoMigrate(&models.ValidatorDescription{})
	db.db.Migrator().AutoMigrate(&models.ValidatorCommission{})
	db.db.Migrator().AutoMigrate(&models.ValidatorVotingPower{})
	db.db.Migrator().AutoMigrate(&models.ValidatorStatus{})
	db.db.Migrator().AutoMigrate(&models.ValidatorSigningInfo{})

	//modules

	return nil
}

// HasBlock implements database.Database
func (db *Database) HasBlock(height int64) (bool, error) {
	var res bool

	var block models.Block
	if err := db.db.Table("blocks").Where("`height` = ?", height).
		First(&block).
		Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		log.Errorf("Other DB error: %s", err.Error())
		return false, err
	}
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

	for _, val := range validators {
		valDo := &models.Validator{
			ConsensusAddress: val.ConsAddr,
			ConsensusPubkey:  val.ConsPubKey,
		}
		db.inTx(valDo, true)
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

func (db *Database) inTx(inst interface{}, overwrite bool) {

	if inst == nil {
		log.Warnw("nil instance diggerSaveAndDontCare")
		return
	}

	tx := db.db.Begin()

	if !overwrite {
		if err := tx.Clauses(clause.Insert{Modifier: "IGNORE"}).CreateInBatches(inst, 2000).Error; err != nil {
			log.Warnw("insert save error", "error", err)
		}
	}

	if err := tx.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Error; err != nil {
		log.Warnw("insert save error", "error", err)
	}

	if err := tx.Commit().Error; err != nil {
		log.Errorw("insert commit error", "error", err)
		tx.Rollback()
	}
}
