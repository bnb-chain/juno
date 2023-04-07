package models

import "github.com/forbole/juno/v4/common"

type Bucket struct {
	ID uint64 `gorm:"column:id;primaryKey"`

	BucketID         common.Hash    `gorm:"column:bucket_id;type:BINARY(32);uniqueIndex:idx_bucket_id"`
	BucketName       string         `gorm:"column:bucket_name;type:varchar(64);uniqueIndex:idx_bucket_name"` // BucketName length between 3 and 63
	OwnerAddress     common.Address `gorm:"column:owner_address;type:BINARY(20);index:idx_owner"`            // OwnerAddress bucket creator
	PaymentAddress   common.Address `gorm:"column:payment_address;type:BINARY(20)"`
	PrimarySpAddress common.Address `gorm:"column:primary_sp_address;type:BINARY(20)"`
	OperatorAddress  common.Address `gorm:"column:operator_address;type:BINARY(20)"`
	SourceType       string         `gorm:"column:source_type;type:VARCHAR(50)"`
	ChargedReadQuota uint64         `gorm:"column:charged_read_quota"`
	Visibility       string         `gorm:"column:visibility;type:VARCHAR(50)"`

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
