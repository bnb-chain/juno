package models

import "github.com/forbole/juno/v4/common"

type Object struct {
	Id                   uint64         `gorm:"id;type:bigint(64);primaryKey"`
	Creator              common.Address `gorm:"creator_address;type:BINARY(20)"`
	Owner                common.Address `gorm:"owner;type:BINARY(20);index:idx_owner"`
	BucketID             int64          `gorm:"bucket_id;type:int;index:idx_bucket_id"`
	BucketName           string         `gorm:"bucket_name;type:varchar(63)"`                              // BucketName length between 3 and 63
	ObjectName           string         `gorm:"object_name;type:varchar(63);uniqueIndex:uniq_object_name"` // BucketName length between 3 and 63
	ObjectID             int64          `gorm:"object_id;type:int;uniqueIndex:uniq_object_id"`
	PayloadSize          int64          `gorm:"payload_size;type:int"`
	IsPublic             bool           `gorm:"is_public;type:tinyint(1)"`
	ContentType          string         `gorm:"content_type;type:varchar(20)"`
	CreateAt             int64          `gorm:"create_at;type:bigint(64)"`
	ObjectStatus         string         `gorm:"object_status;type:varchar(64)"`
	RedundancyType       string         `gorm:"redundancy_type;type:varchar(64)"`
	SourceType           string         `gorm:"source_type;type:varchar(64)"`
	CheckSums            string         `gorm:"checksums;type:text"`
	SecondarySpAddresses string         `gorm:"secondary_sp_addresses;type:text"`
	PrimarySpAddress     common.Address `gorm:"primary_sp_address;type:BINARY(20)"`
	LockedBalance        common.Hash    `gorm:"locked_balance;type:BINARY(32)"`
	Removed              bool           `gorm:"removed"`
}

func (*Object) TableName() string {
	return "objects"
}
