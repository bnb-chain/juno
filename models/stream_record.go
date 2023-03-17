package models

import (
	"github.com/shopspring/decimal"

	"github.com/forbole/juno/v4/common"
)

type StreamRecord struct {
	Id              uint64          `gorm:"column:id;type:bigint(64);primaryKey"`
	Account         common.Address  `gorm:"column:account;type:BINARY(20);not null;index:idx_account"`
	CrudTimestamp   int64           `gorm:"column:update_time;type:bigint(64)"`
	NetflowRate     decimal.Decimal `gorm:"column:netflow_rate"`
	StaticBalance   decimal.Decimal `gorm:"column:static_balance"`
	BufferBalance   decimal.Decimal `gorm:"column:buffer_balance"`
	LockBalance     decimal.Decimal `gorm:"column:lock_balance"`
	Status          string          `gorm:"column:status"`
	SettleTimestamp int64           `gorm:"column:settle_time;type:bigint(64)"`
	OutFlows        string          `gorm:"column:out_flows;type:json"`
}

func (*StreamRecord) TableName() string {
	return "stream_records"
}
