package models

import "github.com/forbole/juno/v4/common"

type Object struct {
	Owner                common.Address `gorm:"owner;type:BINARY(20);index:idx_owner"`
	BucketID             int64          `gorm:"bucket_id;type:int;index:idx_bucket_id"`
	BucketName           string         `gorm:"bucket_name;type:varchar(63)"` // BucketName length between 3 and 63
	ObjectName           string         `gorm:"object_name;type:varchar(63)"` // BucketName length between 3 and 63
	ObjectID             int64          `gorm:"object_id;type:int;primaryKey"`
	PayloadSize          string         `gorm:"payload_size;type:varchar(20)"`
	IsPublic             bool           `gorm:"is_public;type:tinyint(1)"`
	ContentType          string         `gorm:"content_type;type:varchar(20)"`
	CreateAt             int64          `gorm:"create_at;type:bigint(64)"`
	ObjectStatus         int            `gorm:"object_status;type:int"`
	RedundancyType       int            `gorm:"redundancy_type;type:int"`
	SourceType           int            `gorm:"source_type;type:int"`
	CheckSums            string         `gorm:"checksums;type:text"`
	SecondarySpAddresses string         `gorm:"secondary_sp_addresses;type:text"`
	LockedBalance        common.Hash    `gorm:"locked_balance;type:BINARY(32)"`
}

func (*Object) TableName() string {
	return "objects"
}
