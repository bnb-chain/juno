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
	Height          uint64         `gorm:"column:height;uniqueIndex:idx_height"`
	ProposerAddress common.Address `gorm:"column:proposer_address;type:BINARY(20);index:idx_proposer_address"`
	Timestamp       uint64         `gorm:"column:timestamp"`
}

type Block struct {
	ID uint64 `gorm:"column:id;primaryKey"`

	BlockID
	Header

	NumTxs   int    `gorm:"column:num_txs"`
	TotalGas uint64 `gorm:"column:total_gas"`
}
