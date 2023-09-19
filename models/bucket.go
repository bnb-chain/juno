package models

import (
	"github.com/forbole/juno/v4/common"
	"github.com/shopspring/decimal"
)

type Bucket struct {
	ID uint64 `gorm:"column:id;primaryKey"`

	BucketID                   common.Hash    `gorm:"column:bucket_id;type:BINARY(32);uniqueIndex:idx_bucket_id"`
	BucketName                 string         `gorm:"column:bucket_name;type:varchar(64);uniqueIndex:idx_bucket_name"` // BucketName length between 3 and 63
	Owner                      common.Address `gorm:"column:owner;type:BINARY(20);index:idx_owner"`
	PaymentAddress             common.Address `gorm:"column:payment_address;type:BINARY(20)"`
	GlobalVirtualGroupFamilyId uint32         `gorm:"column:global_virtual_group_family_id;index:idx_vgf_id"`
	Operator                   common.Address `gorm:"column:operator;type:BINARY(20)"`
	SourceType                 string         `gorm:"column:source_type;type:VARCHAR(50)"`
	ChargedReadQuota           uint64         `gorm:"column:charged_read_quota"`
	Visibility                 string         `gorm:"column:visibility;type:VARCHAR(50)"`
	Status                     string         `gorm:"column:status;type:varchar(64);"`
	DeleteAt                   int64          `gorm:"column:delete_at"`
	DeleteReason               string         `gorm:"column:delete_reason;type:varchar(256);"`

	StorageSize decimal.Decimal `gorm:"column:storage_size;type:DECIMAL(65, 0);not null"`
	ChargeSize  decimal.Decimal `gorm:"column:charge_size;type:DECIMAL(65, 0);not null"`

	CreateAt     int64       `gorm:"column:create_at"`
	CreateTxHash common.Hash `gorm:"column:create_tx_hash;type:BINARY(32);not null"`
	CreateTime   int64       `gorm:"column:create_time"` // seconds
	UpdateAt     int64       `gorm:"column:update_at"`
	UpdateTxHash common.Hash `gorm:"column:update_tx_hash;type:BINARY(32);not null"`
	UpdateTime   int64       `gorm:"column:update_time"` // seconds
	Removed      bool        `gorm:"column:removed;default:false"`
}

func (*Bucket) TableName() string {
	return "buckets"
}
