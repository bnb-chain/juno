package models

import "github.com/forbole/juno/v4/common"

type Bucket struct {
	ID               uint64         `gorm:"id;type:bigint(64);primaryKey"`
	BucketID         int64          `gorm:"bucket_id;type:int;index:idx_bucket_id"`
	Owner            common.Address `gorm:"owner;type:BINARY(20);index:idx_owner"`
	BucketName       string         `gorm:"bucket_name;type:varchar(63)"` // BucketName length between 3 and 63
	CreateAt         int64          `gorm:"create_at;type:bigint(64)"`    // create at block height
	IsPublic         bool           `gorm:"is_public;type:tinyint(1)"`
	SourceType       string         `gorm:"source_type;type:varchar(63)"`
	PaymentAddress   common.Address `gorm:"payment_address;type:BINARY(20)"`
	PrimarySpAddress common.Address `gorm:"primary_sp_address;type:BINARY(20)"`
	ReadQuota        string         `gorm:"read_quota;type:varchar(256)"`
	PaymentPriceTime int64          `gorm:"payment_price_time;type:bigint(64)"`
	SpAddress        common.Address `gorm:"sp_address;type:BINARY(20)"`
	Rate             int64          `gorm:"rate"`
	Removed          bool           `gorm:"removed"`
	OperatorAddress  common.Address `gorm:"operator_address;type:BINARY(20)"`
}

func (*Bucket) TableName() string {
	return "buckets"
}
