package models

import "github.com/forbole/juno/v4/common"

type Object struct {
	ID uint64 `gorm:"column:id;primaryKey" json:"-"`

	BucketID   common.Hash `gorm:"column:bucket_id;type:BINARY(32);index:idx_bucket_id"`
	BucketName string      `gorm:"column:bucket_name;type:varchar(63)"`
	ObjectID   common.Hash `gorm:"column:object_id;type:BINARY(32);uniqueIndex:idx_object_id"`
	ObjectName string      `gorm:"column:object_name;type:varchar(63)"`

	CreatorAddress       common.Address   `gorm:"column:creator_address;type:BINARY(20)"`
	OwnerAddress         common.Address   `gorm:"column:owner_address;type:BINARY(20);index:idx_owner"`
	PrimarySpAddress     common.Address   `gorm:"column:primary_sp_address;type:BINARY(20)"`
	OperatorAddress      common.Address   `gorm:"column:operator_address;type:BINARY(20)"`
	SecondarySpAddresses []common.Address `gorm:"secondary_sp_addresses;type:BINARY(80)"`
	PayloadSize          uint64           `gorm:"column:payload_size"`
	IsPublic             bool             `gorm:"column:is_public"`
	ContentType          string           `gorm:"column:content_type"`
	Status               string           `gorm:"column:status;type:VARCHAR(50)"`
	RedundancyType       string           `gorm:"column:redundancy_type;type:VARCHAR(50)"`
	SourceType           string           `gorm:"column:source_type;type:VARCHAR(50)"`
	CheckSums            [][]byte         `gorm:"column:checksums;type:blob"`

	CreateAt   int64 `gorm:"column:create_at"`
	CreateTime int64 `gorm:"column:create_time"`
	UpdateAt   int64 `gorm:"column:update_at"`
	UpdateTime int64 `gorm:"column:update_time"`
	Removed    bool  `gorm:"column:removed;default:false"`
}

func (*Object) TableName() string {
	return "objects"
}
