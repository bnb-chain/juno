package models

import (
	"github.com/forbole/juno/v4/common"
)

type Tx struct {
	ID uint64 `gorm:"column:id;primaryKey"`

	Hash      common.Hash `gorm:"column:hash;type:BINARY(32);uniqueIndex:idx_hash"`
	Height    uint64      `gorm:"height;uniqueIndex:idx_height_tx_index,priority:1"`
	BlockHash common.Hash `gorm:"column:block_hash;type:BINARY(32)"`
	TxIndex   uint32      `gorm:"column:tx_index;uniqueIndex:idx_height_tx_index,priority:2"`

	Success     bool   `gorm:"success"`
	Messages    string `gorm:"messages;type:jsonb;not null;default:'[]'"`
	Memo        string `gorm:"memo"`
	Signatures  string `gorm:"signatures"`
	SignerInfos string `gorm:"signer_infos;type:jsonb;not null;default:'[]'"`
	Fee         string `gorm:"fee;type:jsonb;not null;default:'[]'"`

	GasWanted uint64 `gorm:"gas_wanted"`
	GasUsed   uint64 `gorm:"gas_used"`
	RawLog    string `gorm:"raw_log"`
	Logs      string `gorm:"logs;type:jsonb"`

	// Timestamp uint64 `gorm:"column:timestamp"` ?
}
