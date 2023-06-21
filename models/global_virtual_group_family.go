package models

import (
	"github.com/lib/pq"

	"github.com/forbole/juno/v4/common"
)

type GlobalVirtualGroupFamily struct {
	ID                         uint64         `gorm:"column:id;primaryKey"`
	GlobalVirtualGroupFamilyId uint32         `gorm:"column:global_virtual_group_family_id;index:idx_vgf_id"`
	GlobalVirtualGroupIds      pq.StringArray `gorm:"column:global_virtual_group_ids;type:MEDIUMTEXT"`
	VirtualPaymentAddress      common.Address `gorm:"column:virtual_payment_address;type:BINARY(20)"`

	CreateAt     int64       `gorm:"column:create_at"`
	CreateTxHash common.Hash `gorm:"column:create_tx_hash;type:BINARY(32);not null"`
	CreateTime   int64       `gorm:"column:create_time"` // seconds
	UpdateAt     int64       `gorm:"column:update_at"`
	UpdateTxHash common.Hash `gorm:"column:update_tx_hash;type:BINARY(32);not null"`
	UpdateTime   int64       `gorm:"column:update_time"` // seconds
	Removed      bool        `gorm:"column:removed;default:false"`
}

func (*GlobalVirtualGroupFamily) TableName() string {
	return "global_virtual_group_families"
}
