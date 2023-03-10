package models

import (
	"github.com/shopspring/decimal"

	"github.com/forbole/juno/v4/common"
)

type StreamRecord struct {
	Id              uint64          `gorm:"id;type:bigint(64);primaryKey"`
	Account         common.Address  `gorm:"account;type:BINARY(20);not null;index:idx_account"`
	CrudTimestamp   int64           `gorm:"update_time;type:bigint(64)"`
	NetflowRate     decimal.Decimal `gorm:"netflow_rate"`
	StaticBalance   decimal.Decimal `gorm:"static_balance"`
	BufferBalance   decimal.Decimal `gorm:"buffer_balance"`
	LockBalance     decimal.Decimal `gorm:"lock_balance"`
	Status          string          `gorm:"status"`
	SettleTimestamp int64           `gorm:"settle_time;type:bigint(64)"`
	OutFlows        string          `gorm:"out_flows;type:json"`
}

func (*StreamRecord) TableName() string {
	return "stream_records"
}
