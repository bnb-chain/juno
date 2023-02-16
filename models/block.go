package models

import (
	"github.com/forbole/juno/v4/common"
)

type PartSetHeader struct {
	Total uint32      `gorm:"-"`
	Hash  common.Hash `gorm:"-"`
}

type BlockID struct {
	Hash          common.Hash   `gorm:"column:hash;type:BINARY(32);uniqueIndex:idx_hash"`
	PartSetHeader PartSetHeader `gorm:"-"`
}

type Header struct {
	Height             uint64      `gorm:"column:height;uniqueIndex:idx_height"`
	LastCommitHash     common.Hash `gorm:"last_commit_hash"`     // commit from validators from the last block
	DataHash           common.Hash `gorm:"data_hash"`            // transactions
	ValidatorsHash     common.Hash `gorm:"validators_hash"`      // validators for the current block
	NextValidatorsHash common.Hash `gorm:"next_validators_hash"` // validators for the next block
	ConsensusHash      common.Hash `gorm:"consensus_hash"`       // consensus params for current block
	AppHash            common.Hash `gorm:"app_hash"`             // state after txs from the previous block
	// root hash of all results from the txs from the previous block
	// see `deterministicResponseDeliverTx` to understand which parts of a tx is hashed into here
	LastResultsHash common.Hash `gorm:"last_results_hash"`
	// consensus info
	EvidenceHash    common.Hash    `gorm:"evidence_hash"` // evidence included in the block
	ProposerAddress common.Address `gorm:"column:proposer_address;type:BINARY(20);index:idx_proposer_address"`

	Timestamp uint64 `gorm:"column:timestamp"`
}

type Block struct {
	ID uint64 `gorm:"column:id;primaryKey"`

	BlockID
	Header

	NumTxs   int    `gorm:"column:num_txs"`
	TotalGas uint64 `gorm:"column:total_gas"`
}

func (*Block) TableName() string {
	return "blocks"
}

type Genesis struct {
	OneRowId      bool   `gorm:"one_row_id;primaryKey;default:true"`
	ChainID       string `gorm:"chain_id"`
	Timestamp     uint64 `gorm:"timestamp"`
	InitialHeight uint64 `gorm:"initial_height"`
}

func (*Genesis) TableName() string {
	return "geneses"
}
