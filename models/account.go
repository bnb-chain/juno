package models

import "github.com/forbole/juno/v4/common"

type AccountType string

const (
	GeneralAccount AccountType = "general"
	PaymentAccount AccountType = "payment"
)

type Account struct {
	ID uint64 `gorm:"column:id;primaryKey"`

	Address             common.Address `gorm:"column:address;type:BINARY(20);uniqueIndex:idx_address"`
	Type                AccountType    `gorm:"column:type;not null;default:'general'"`
	Balance             *common.Big    `gorm:"column:balance" json:"-"` // changed from events of coin_spent and coin_receive
	BalanceU64          uint64         `gorm:"-" json:"balance"`
	TxCount             uint64         `gorm:"column:tx_count;not null;default:0"`
	LastActiveTimestamp uint64         `gorm:"last_active_timestamp"`
	Refundable          bool           `gorm:"column:refundable;not null;default:true"`
}

func (*Account) TableName() string {
	return "accounts"
}

type AccountGroup struct {
	ID uint64 `gorm:"column:id;primaryKey"`

	Address common.Address `gorm:"column:address;type:BINARY(20);index:idx_address"`
	GroupID uint64         `gorm:"column:group_id;index:idx_group_id"`
}

func (*AccountGroup) TableName() string {
	return "account_groups"
}

// AccountRelation currently only need by explorer
type AccountRelation struct {
	ID uint64 `gorm:"column:id;primaryKey"`

	Owner   common.Address `gorm:"column:owner;type:BINARY(20);not null;index:idx_owner;uniqueIndex:idx_owner_payment,priority:1"`
	Payment common.Address `gorm:"column:payment;type:BINARY(20);not null;index:idx_payment;uniqueIndex:idx_owner_payment,priority:2"`
	Index   uint64         `gorm:"column:index"`
}

func (*AccountRelation) TableName() string {
	return "account_relations"
}
