package models

import (
	"github.com/lib/pq"

	"github.com/forbole/juno/v4/common"
)

type Object struct {
	ID uint64 `gorm:"column:id;primaryKey"`

	BucketID   common.Hash `gorm:"column:bucket_id;type:BINARY(32);index:idx_bucket_id"`
	BucketName string      `gorm:"column:bucket_name;type:varchar(64);index:idx_bucket_name_object_name,priority:1"`
	ObjectID   common.Hash `gorm:"column:object_id;type:BINARY(32);uniqueIndex:idx_object_id"`
	ObjectName string      `gorm:"column:object_name;type:varchar(1024);index:idx_bucket_name_object_name,length:512,priority:2"`

	Creator             common.Address `gorm:"column:creator;type:BINARY(20)"`
	Owner               common.Address `gorm:"column:owner;type:BINARY(20);index:idx_owner"`
	LocalVirtualGroupId uint32         `gorm:"column:local_virtual_group_id;index:idx_lvg_id"`
	Operator            common.Address `gorm:"column:operator;type:BINARY(20)"`
	PayloadSize         uint64         `gorm:"column:payload_size"`
	Visibility          string         `gorm:"column:visibility;type:VARCHAR(50)"`
	ContentType         string         `gorm:"column:content_type"`
	Status              string         `gorm:"column:status;type:VARCHAR(50)"`
	RedundancyType      string         `gorm:"column:redundancy_type;type:VARCHAR(50)"`
	SourceType          string         `gorm:"column:source_type;type:VARCHAR(50)"`
	CheckSums           pq.ByteaArray  `gorm:"column:checksums;type:text"`
	DeleteAt            int64          `gorm:"column:delete_at"`
	DeleteReason        string         `gorm:"column:delete_reason;type:varchar(256);"`

	CreateAt     int64       `gorm:"column:create_at"`
	CreateTxHash common.Hash `gorm:"column:create_tx_hash;type:BINARY(32);not null"`
	CreateTime   int64       `gorm:"column:create_time"` // seconds
	UpdateAt     int64       `gorm:"column:update_at;index:idx_update_at"`
	UpdateTxHash common.Hash `gorm:"column:update_tx_hash;type:BINARY(32);not null"`
	SealedTxHash common.Hash `gorm:"column:sealed_tx_hash;type:BINARY(32)"`
	UpdateTime   int64       `gorm:"column:update_time"` // seconds
	Removed      bool        `gorm:"column:removed;default:false"`

	Tags string `gorm:"column:tags;TYPE:json"` // tags

}

func (*Object) TableName() string {
	return "objects"
}
