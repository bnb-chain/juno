package database

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/forbole/juno/v4/log"
	"github.com/forbole/juno/v4/models"
	"github.com/forbole/juno/v4/types/config"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"strings"

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
	HasBlock(ctx context.Context, height int64) (bool, error)

	// GetLastBlockHeight returns the last block height stored in database..
	// An error is returned if the operation fails.
	GetLastBlockHeight(ctx context.Context) (int64, error)

	// GetMissingHeights returns a slice of missing block heights between startHeight and endHeight
	GetMissingHeights(ctx context.Context, startHeight, endHeight int64) []int64

	// SaveBlock will be called when a new block is parsed, passing the block itself
	// and the transactions contained inside that block.
	// An error is returned if the operation fails.
	// NOTE. For each transaction inside txs, SaveTx will be called as well.
	SaveBlock(ctx context.Context, block *types.Block) error

	// GetTotalBlocks returns total number of blocks stored in database.
	GetTotalBlocks(ctx context.Context) int64

	// SaveTx will be called to save each transaction contained inside a block.
	// An error is returned if the operation fails.
	SaveTx(ctx context.Context, tx *types.Tx) error

	// HasValidator returns true if a given validator by consensus address exists.
	// An error is returned if the operation fails.
	HasValidator(ctx context.Context, address string) (bool, error)

	// SaveValidators stores a list of validators if they do not already exist.
	// An error is returned if the operation fails.
	SaveValidators(ctx context.Context, validators []*types.Validator) error

	// SaveCommitSignatures stores a  slice of validator commit signatures.
	// An error is returned if the operation fails.
	SaveCommitSignatures(ctx context.Context, signatures []*types.CommitSig) error

	// SaveMessage stores a single message.
	// An error is returned if the operation fails.
	SaveMessage(ctx context.Context, msg *types.Message) error

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

type Impl struct {
	Db             *gorm.DB
	EncodingConfig *params.EncodingConfig
}

// createPartitionIfNotExists creates a new partition having the given partition id if not existing
func (db *Impl) createPartitionIfNotExists(table string, partitionID int64) error {
	partitionTable := fmt.Sprintf("%s_%d", table, partitionID)

	stmt := fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s PARTITION OF %s FOR VALUES IN (%d)",
		partitionTable,
		table,
		partitionID,
	)
	err := db.Db.Exec(stmt).Error

	if err != nil {
		return err
	}

	return nil
}

// -------------------------------------------------------------------------------------------------------------------

func (db *Impl) PrepareTables(ctx context.Context) error {
	return nil
}

// HasBlock implements database.Database
func (db *Impl) HasBlock(ctx context.Context, height int64) (bool, error) {
	var res bool

	var block models.Block
	if err := db.Db.Table("blocks").Where("height = ?", height).
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
func (db *Impl) GetLastBlockHeight(ctx context.Context) (int64, error) {
	stmt := `SELECT height FROM blocks ORDER BY height DESC LIMIT 1;`

	var height int64
	err := db.Db.Raw(stmt).Scan(&height).Error
	if err != nil {
		if errIsNotFound(err) {
			// If no rows stored in block table, return 0 as height
			return 0, nil
		}
		return 0, fmt.Errorf("error while getting last block height, error: %s", err)
	}

	return height, nil
}

// GetMissingHeights returns a slice of missing block heights between startHeight and endHeight
func (db *Impl) GetMissingHeights(ctx context.Context, startHeight, endHeight int64) []int64 {
	var result []int64
	stmt := `SELECT generate_series($1::int,$2::int) EXCEPT SELECT height FROM blocks ORDER BY 1;`
	err := db.Db.Select(&result, stmt, startHeight, endHeight)
	if err != nil {
		return nil
	}

	if len(result) == 0 {
		return nil
	}

	return result
}

// SaveBlock implements database.Database
func (db *Impl) SaveBlock(ctx context.Context, block *types.Block) error {
	sqlStatement := `
INSERT INTO block (height, hash, num_txs, total_gas, proposer_address, timestamp)
VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT DO NOTHING`

	proposerAddress := sql.NullString{Valid: len(block.ProposerAddress) != 0, String: block.ProposerAddress}
	err := db.Db.Exec(sqlStatement,
		block.Height, block.Hash, block.TxNum, block.TotalGas, proposerAddress, block.Timestamp,
	).Error
	return err
}

// GetTotalBlocks implements database.Database
func (db *Impl) GetTotalBlocks(ctx context.Context) int64 {
	var blockCount int64
	err := db.Db.Raw(`SELECT count(*) FROM blocks;`).Scan(&blockCount)
	if err != nil {
		return 0
	}

	return blockCount
}

// SaveTx implements database.Database
func (db *Impl) SaveTx(ctx context.Context, tx *types.Tx) error {
	var partitionID int64

	partitionSize := config.Cfg.Database.PartitionSize
	if partitionSize > 0 {
		partitionID = tx.Height / partitionSize
		err := db.createPartitionIfNotExists("transaction", partitionID)
		if err != nil {
			return err
		}
	}

	return db.saveTxInsidePartition(tx, partitionID)
}

// saveTxInsidePartition stores the given transaction inside the partition having the given id
func (db *Impl) saveTxInsidePartition(tx *types.Tx, partitionID int64) error {
	sqlStatement := `
INSERT INTO transaction 
(hash, height, success, messages, memo, signatures, signer_infos, fee, gas_wanted, gas_used, raw_log, logs, partition_id) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) 
ON CONFLICT (hash, partition_id) DO UPDATE 
	SET height = excluded.height, 
		success = excluded.success, 
		messages = excluded.messages,
		memo = excluded.memo, 
		signatures = excluded.signatures, 
		signer_infos = excluded.signer_infos,
		fee = excluded.fee, 
		gas_wanted = excluded.gas_wanted, 
		gas_used = excluded.gas_used,
		raw_log = excluded.raw_log, 
		logs = excluded.logs`

	var sigs = make([]string, len(tx.Signatures))
	for index, sig := range tx.Signatures {
		sigs[index] = base64.StdEncoding.EncodeToString(sig)
	}

	var msgs = make([]string, len(tx.Body.Messages))
	for index, msg := range tx.Body.Messages {
		bz, err := db.EncodingConfig.Codec.MarshalJSON(msg)
		if err != nil {
			return err
		}
		msgs[index] = string(bz)
	}
	msgsBz := fmt.Sprintf("[%s]", strings.Join(msgs, ","))

	feeBz, err := db.EncodingConfig.Codec.MarshalJSON(tx.AuthInfo.Fee)
	if err != nil {
		return fmt.Errorf("failed to JSON encode tx fee: %s", err)
	}

	var sigInfos = make([]string, len(tx.AuthInfo.SignerInfos))
	for index, info := range tx.AuthInfo.SignerInfos {
		bz, err := db.EncodingConfig.Codec.MarshalJSON(info)
		if err != nil {
			return err
		}
		sigInfos[index] = string(bz)
	}
	sigInfoBz := fmt.Sprintf("[%s]", strings.Join(sigInfos, ","))

	logsBz, err := db.EncodingConfig.Amino.MarshalJSON(tx.Logs)
	if err != nil {
		return err
	}

	err = db.Db.Exec(sqlStatement,
		tx.TxHash, tx.Height, tx.Successful(),
		msgsBz, tx.Body.Memo, pq.Array(sigs),
		sigInfoBz, string(feeBz),
		tx.GasWanted, tx.GasUsed, tx.RawLog, string(logsBz),
		partitionID,
	).Error
	return err
}

// HasValidator implements database.Database
func (db *Impl) HasValidator(ctx context.Context, addr string) (bool, error) {
	var res bool
	stmt := `SELECT EXISTS(SELECT 1 FROM validator WHERE consensus_address = $1);`
	err := db.Db.Raw(stmt, addr).Scan(&res).Error
	return res, err
}

// SaveValidators implements database.Database
func (db *Impl) SaveValidators(ctx context.Context, validators []*types.Validator) error {
	if len(validators) == 0 {
		return nil
	}

	stmt := `INSERT INTO validators (consensus_address, consensus_pubkey) VALUES `

	var vparams []interface{}
	for _, val := range validators {
		stmt += fmt.Sprintf("(?, ?),")
		vparams = append(vparams, val.ConsAddr, val.ConsPubKey)
	}

	stmt = stmt[:len(stmt)-1] // Remove trailing ,
	stmt += " ON CONFLICT DO NOTHING"
	err := db.Db.Exec(stmt, vparams...).Error
	return err
}

// SaveCommitSignatures implements database.Database
func (db *Impl) SaveCommitSignatures(ctx context.Context, signatures []*types.CommitSig) error {
	if len(signatures) == 0 {
		return nil
	}

	stmt := `INSERT INTO pre_commit (validator_address, height, timestamp, voting_power, proposer_priority) VALUES `

	var sparams []interface{}
	for i, sig := range signatures {
		si := i * 5

		stmt += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d),", si+1, si+2, si+3, si+4, si+5)
		sparams = append(sparams, sig.ValidatorAddress, sig.Height, sig.Timestamp, sig.VotingPower, sig.ProposerPriority)
	}

	stmt = stmt[:len(stmt)-1]
	stmt += " ON CONFLICT (validator_address, timestamp) DO NOTHING"
	err := db.Db.Exec(stmt, sparams...).Error
	return err
}

// SaveMessage implements database.Database
func (db *Impl) SaveMessage(ctx context.Context, msg *types.Message) error {
	var partitionID int64
	partitionSize := config.Cfg.Database.PartitionSize
	if partitionSize > 0 {
		partitionID = msg.Height / partitionSize
		err := db.createPartitionIfNotExists("message", partitionID)
		if err != nil {
			return err
		}
	}

	return db.saveMessageInsidePartition(msg, partitionID)
}

// saveMessageInsidePartition stores the given message inside the partition having the provided id
func (db *Impl) saveMessageInsidePartition(msg *types.Message, partitionID int64) error {
	stmt := `
INSERT INTO message(transaction_hash, index, type, value, involved_accounts_addresses, height, partition_id) 
VALUES ($1, $2, $3, $4, $5, $6, $7) 
ON CONFLICT (transaction_hash, index, partition_id) DO UPDATE 
	SET height = excluded.height, 
		type = excluded.type,
		value = excluded.value,
		involved_accounts_addresses = excluded.involved_accounts_addresses`

	err := db.Db.Exec(stmt, msg.TxHash, msg.Index, msg.Type, msg.Value, pq.Array(msg.Addresses), msg.Height, partitionID).Error
	return err
}

// Close implements database.Database
func (db *Impl) Close() {
	var err error
	if err != nil {
		log.Errorw("error while closing connection", "err", err)
	}
}

// -------------------------------------------------------------------------------------------------------------------

// GetLastPruned implements database.PruningDb
func (db *Impl) GetLastPruned() (int64, error) {
	var lastPrunedHeight int64
	err := db.Db.Raw(`SELECT coalesce(MAX(last_pruned_height),0) FROM pruning LIMIT 1;`).Scan(&lastPrunedHeight).Error
	return lastPrunedHeight, err
}

// StoreLastPruned implements database.PruningDb
func (db *Impl) StoreLastPruned(height int64) error {
	err := db.Db.Exec(`DELETE FROM pruning`).Error
	if err != nil {
		return err
	}

	err = db.Db.Exec(`INSERT INTO pruning (last_pruned_height) VALUES ($1)`, height).Error
	return err
}

// Prune implements database.PruningDb
func (db *Impl) Prune(height int64) error {
	err := db.Db.Exec(`DELETE FROM pre_commit WHERE height = $1`, height).Error
	if err != nil {
		return err
	}

	err = db.Db.Exec(`
DELETE FROM message 
USING transaction 
WHERE message.transaction_hash = transaction.hash AND transaction.height = $1
`, height).Error
	return err
}

func errIsNotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows) || errors.Is(err, gorm.ErrRecordNotFound)
}
