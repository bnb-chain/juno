package models

import (
	"github.com/forbole/juno/v4/common"

	"github.com/lib/pq"
)

type Object struct {
	ID uint64 `gorm:"column:id;primaryKey"`

	BucketID   common.Hash `gorm:"column:bucket_id;type:BINARY(32);index:idx_bucket_id"`
	BucketName string      `gorm:"column:bucket_name;type:varchar(63)"`
	ObjectID   common.Hash `gorm:"column:object_id;type:BINARY(32);uniqueIndex:idx_object_id"`
	ObjectName string      `gorm:"column:object_name;type:varchar"`

	CreatorAddress       common.Address `gorm:"column:creator_address;type:BINARY(20)"`
	OwnerAddress         common.Address `gorm:"column:owner_address;type:BINARY(20);index:idx_owner"`
	PrimarySpAddress     common.Address `gorm:"column:primary_sp_address;type:BINARY(20)"`
	OperatorAddress      common.Address `gorm:"column:operator_address;type:BINARY(20)"`
	SecondarySpAddresses pq.StringArray `gorm:"column:secondary_sp_addresses;type:BINARY(80)"`
	PayloadSize          uint64         `gorm:"column:payload_size"`
	Visibility           string         `gorm:"column:visibility;type:VARCHAR(50)"`
	ContentType          string         `gorm:"column:content_type"`
	Status               string         `gorm:"column:status;type:VARCHAR(50)"`
	RedundancyType       string         `gorm:"column:redundancy_type;type:VARCHAR(50)"`
	SourceType           string         `gorm:"column:source_type;type:VARCHAR(50)"`
	CheckSums            pq.ByteaArray  `gorm:"column:checksums;type:blob"`

	CreateAt     int64       `gorm:"column:create_at"`
	CreateTxHash common.Hash `gorm:"column:create_tx_hash;type:BINARY(32);not null"`
	CreateTime   int64       `gorm:"column:create_time"` // seconds
	UpdateAt     int64       `gorm:"column:update_at"`
	UpdateTxHash common.Hash `gorm:"column:update_tx_hash;type:BINARY(32);not null"`
	UpdateTime   int64       `gorm:"column:update_time"` // seconds
	Removed      bool        `gorm:"column:removed;default:false"`
}

func (*Object) TableName() string {
	return "objects"
}
