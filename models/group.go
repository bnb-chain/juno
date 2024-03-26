package models

import (
	"github.com/forbole/juno/v4/common"
	"gorm.io/datatypes"
)

type Group struct {
	ID         uint64         `gorm:"column:id;primaryKey"`
	Owner      common.Address `gorm:"column:owner;type:BINARY(20);index:idx_owner"`
	GroupID    common.Hash    `gorm:"column:group_id;type:BINARY(32);index:idx_group_id;uniqueIndex:idx_account_group,priority:2"`
	GroupName  string         `gorm:"column:group_name;type:varchar(63);index:idx_group_name"`
	SourceType string         `gorm:"column:source_type;type:varchar(63)"`
	Extra      string         `gorm:"column:extra;type:varchar(512)"`

	AccountID      common.Address `gorm:"column:account_id;type:BINARY(20);uniqueIndex:idx_account_group,priority:1"`
	Operator       common.Address `gorm:"column:operator;type:BINARY(20)"`
	ExpirationTime int64          `gorm:"column:expiration_time"`

	CreateAt   int64 `gorm:"column:create_at"`
	CreateTime int64 `gorm:"column:create_time"`
	UpdateAt   int64 `gorm:"column:update_at"`
	UpdateTime int64 `gorm:"column:update_time"`
	Removed    bool  `gorm:"column:removed;default:false"`

	Tags datatypes.JSON `gorm:"column:tags;TYPE:json"` // tags

}

func (*Group) TableName() string {
	return "groups"
}
