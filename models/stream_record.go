package models

import (
	"github.com/forbole/juno/v4/common"
)

type StreamRecord struct {
	ID uint64 `gorm:"column:id;primaryKey" json:"-"`

	Account         common.Address `gorm:"column:account;type:BINARY(20);uniqueIndex:idx_account"`
	CrudTimestamp   int64          `gorm:"column:crud_timestamp"`
	NetflowRate     *common.Big    `gorm:"column:netflow_rate"`
	StaticBalance   *common.Big    `gorm:"column:static_balance"`
	BufferBalance   *common.Big    `gorm:"column:buffer_balance"`
	LockBalance     *common.Big    `gorm:"column:lock_balance"`
	Status          string         `gorm:"column:status"`
	SettleTimestamp int64          `gorm:"column:settle_timestamp"`
}

func (*StreamRecord) TableName() string {
	return "stream_records"
}
