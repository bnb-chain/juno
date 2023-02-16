package models

type Account struct {
	Address string `gorm:"address;primaryKey"`
	Type    string `gorm:"type"`
	Balance uint64 `gorm:"balance"` // changed from events of coin_spent and coin_receive
	TxCount uint64 `gorm:"tx_count"`
}

func (*Account) TableName() string {
	return "accounts"
}
