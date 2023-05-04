package models

import (
	"github.com/forbole/juno/v4/common"
)

type Bucket struct {
	ID uint64 `gorm:"column:id;primaryKey"`

	BucketID         common.Hash    `gorm:"column:bucket_id;type:BINARY(32);uniqueIndex:idx_bucket_id"`
	BucketName       string         `gorm:"column:bucket_name;type:varchar(64);uniqueIndex:idx_bucket_name"`
	Owner            common.Address `gorm:"column:owner;type:BINARY(20);index:idx_owner"` // Owner bucket creator
	PaymentAddress   common.Address `gorm:"column:payment_address;type:BINARY(20)"`
	PrimarySpAddress common.Address `gorm:"column:primary_sp_address;type:BINARY(20)"`
	Operator         common.Address `gorm:"column:operator;type:BINARY(20)"`
	SourceType       string         `gorm:"column:source_type;type:VARCHAR(50)"`
	ChargedReadQuota uint64         `gorm:"column:charged_read_quota"`
	Visibility       string         `gorm:"column:visibility;type:VARCHAR(50)"`
	Status           string         `gorm:"column:status;type:VARCHAR(64)"`
	Removed          bool           `gorm:"column:removed;default:false"`

	CreateHeight int64       `gorm:"column:create_height"`
	CreateTxHash common.Hash `gorm:"column:create_tx_hash;type:BINARY(32);not null"`
	CreateTime   int64       `gorm:"column:create_time"`
	UpdateHeight int64       `gorm:"column:update_height"`
	UpdateTxHash common.Hash `gorm:"column:update_tx_hash;type:BINARY(32);not null"`
	UpdateTime   int64       `gorm:"column:update_time"`
	DeleteTime   int64       `gorm:"column:delete_time"`
	DeleteReason string      `gorm:"column:delete_reason;type:VARCHAR(256)"`
}

func (*Bucket) TableName() string {
	return "buckets"
}
