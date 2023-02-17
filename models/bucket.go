package models

import "github.com/forbole/juno/v4/common"

type Bucket struct {
	Owner            common.Address `gorm:"owner;type:BINARY(20);index:idx_owner"`
	BucketID         int64          `gorm:"bucket_id;type:int;primaryKey"`
	BucketName       string         `gorm:"bucket_name;type:varchar(63)"` // BucketName length between 3 and 63
	CreateAt         int64          `gorm:"create_at;type:bigint(64)"`    // create at block height
	IsPublic         bool           `gorm:"is_public;type:tinyint(1)"`
	SourceType       int            `gorm:"source_type;type:int"`
	PaymentAddress   common.Address `gorm:"payment_address";type:BINARY(20)`
	PrimarySpAddress common.Address `gorm:"primary_sp_address";type:BINARY(20)`
	ReadQuota        int32          `gorm:"read_quota;type:bigint(64)"`
	PaymentPriceTime int64          `gorm:"payment_price_time";type:bigint(64)`
	SpAddress        common.Address `gorm:"sp_address;type:BINARY(20)"`
	Rate             int64          `gorm:"rate"`
}

func (*Bucket) TableName() string {
	return "bucket"
}
