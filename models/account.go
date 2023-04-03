package models

import "github.com/forbole/juno/v4/common"

type AccountType string

const (
	GeneralAccountType AccountType = "general"
	PaymentAccountType AccountType = "payment"
)

type Account struct {
	ID uint64 `gorm:"column:id;primaryKey" json:"-"`

	Address             common.Address `gorm:"column:address;type:BINARY(20);uniqueIndex:idx_address"`
	Type                AccountType    `gorm:"column:type;not null;default:'general'"`
	Balance             *common.Big    `gorm:"column:balance" json:"-"`
	TxCount             uint64         `gorm:"column:tx_count;not null;default:0"`
	LastActiveTimestamp uint64         `gorm:"last_active_timestamp"`

	BalanceString string `gorm:"-" json:"balance"`
}

func (*Account) TableName() string {
	return "accounts"
}

type AccountGroup struct {
	ID uint64 `gorm:"column:id;primaryKey" json:"-"`

	Address common.Address `gorm:"column:address;type:BINARY(20);index:idx_address"`
	GroupID uint64         `gorm:"column:group_id;index:idx_group_id"`
}

func (*AccountGroup) TableName() string {
	return "account_groups"
}
