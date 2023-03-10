package models

import "github.com/ethereum/go-ethereum/common"

type Epoch struct {
	ID          int         `gorm:"id;type:tinyint(1);uniqueIndex:uniq_id"`
	BlockHeight int64       `gorm:"block_height;type:bigint(64)"`
	BlockHash   common.Hash `gorm:"block_hash;type:BINARY(32)"`
	UpdateTime  int64       `gorm:"update_time;type:bigint(64)"`
}

func (*Epoch) TableName() string {
	return "epoch"
}
