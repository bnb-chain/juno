package models

import (
	"github.com/forbole/juno/v4/common"
)

type Tx struct {
	ID uint64 `gorm:"column:id;primaryKey"`

	Hash    common.Hash `gorm:"column:hash;type:BINARY(32);not null;uniqueIndex:idx_hash"`
	Height  uint64      `gorm:"column:height;not null;uniqueIndex:idx_height_tx_index,priority:1"`
	TxIndex uint32      `gorm:"column:tx_index;not null;uniqueIndex:idx_height_tx_index,priority:2"`

	Success     bool   `gorm:"column:success"`
	Messages    string `gorm:"column:messages;type:json;not null;default:(JSON_ARRAY())"`
	Memo        string `gorm:"column:memo"`
	Signatures  string `gorm:"column:signatures"`
	SignerInfos string `gorm:"column:signer_infos;type:json;not null;default:(JSON_ARRAY())"`
	Fee         string `gorm:"column:fee;type:json;not null;default:(JSON_ARRAY())"`

	GasWanted uint64 `gorm:"column:gas_wanted"`
	GasUsed   uint64 `gorm:"column:gas_used"`
	RawLog    string `gorm:"column:raw_log"`
	Logs      string `gorm:"column:logs;type:json;not null;default:(JSON_ARRAY())"`

	// Timestamp uint64 `gorm:"column:timestamp"` ?
}

func (*Tx) TableName() string {
	return "txs"
}
