package models

import "github.com/forbole/juno/v4/common"

type Bucket struct {
	ID uint64 `gorm:"column:id;primaryKey" json:"-"`

	BucketID         common.Hash    `gorm:"column:bucket_id;type:BINARY(32);uniqueIndex:idx_bucket_id"`
	BucketName       string         `gorm:"column:bucket_name;type:varchar(63);uniqueIndex:idx_bucket_name"` // BucketName length between 3 and 63
	OwnerAddress     common.Address `gorm:"column:owner_address;type:BINARY(20);index:idx_owner"`            // OwnerAddress bucket creator
	PaymentAddress   common.Address `gorm:"column:payment_address;type:BINARY(20)"`
	PrimarySpAddress common.Address `gorm:"column:primary_sp_address;type:BINARY(20)"`
	OperatorAddress  common.Address `gorm:"column:operator_address;type:BINARY(20)"`
	SourceType       string         `gorm:"column:source_type;type:VARCHAR(50)"`
	ReadQuota        uint64         `gorm:"column:read_quota"`
	Visibility       int32          `gorm:"column:visibility"`

	CreateAt   int64 `gorm:"column:create_at"`
	CreateTime int64 `gorm:"column:create_time"`
	UpdateAt   int64 `gorm:"column:update_at"`
	UpdateTime int64 `gorm:"column:update_time"`
	Removed    bool  `gorm:"column:removed;default:false"`
}

func (*Bucket) TableName() string {
	return "buckets"
}
