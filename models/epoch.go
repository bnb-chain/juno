package models

import "github.com/forbole/juno/v4/common"

type Epoch struct {
	OneRowId    bool        `gorm:"one_row_id;not null;default:true;primaryKey"`
	BlockHeight int64       `gorm:"block_height;type:bigint(64)"`
	BlockHash   common.Hash `gorm:"block_hash;type:BINARY(32)"`
	UpdateTime  int64       `gorm:"update_time;type:bigint(64)"`
}

func (*Epoch) TableName() string {
	return "epoch"
}
