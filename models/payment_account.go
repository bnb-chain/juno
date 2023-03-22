package models

import "github.com/forbole/juno/v4/common"

type PaymentAccount struct {
	ID uint64 `gorm:"column:id;primaryKey" json:"-"`

	Addr       common.Address `gorm:"column:addr;type:BINARY(20);not null;uniqueIndex:idx_addr"`
	Owner      common.Address `gorm:"column:owner;type:BINARY(20);not null;index:idx_owner"`
	Refundable bool           `gorm:"column:refundable;not null;default:true"`

	UpdateAt   int64 `gorm:"column:update_at;type:bigint(64)"`
	UpdateTime int64 `gorm:"column:update_time;type:bigint(64)"`
}

func (*PaymentAccount) TableName() string {
	return "payment_accounts"
}
