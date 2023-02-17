package models

type Account struct {
	Address string `gorm:"column:address;primaryKey"`
	Type    string `gorm:"column:type"`
	Balance uint64 `gorm:"column:balance"` // changed from events of coin_spent and coin_receive
	TxCount uint64 `gorm:"column:tx_count"`
}

func (*Account) TableName() string {
	return "accounts"
}
