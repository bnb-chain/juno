package models

import "github.com/forbole/juno/v4/common"

type Group struct {
	ID         uint64         `gorm:"column:id;primaryKey"`
	Owner      common.Address `gorm:"column:owner;type:BINARY(20);index:idx_owner"`
	GroupID    common.Hash    `gorm:"column:group_id;type:BINARY(32);index:idx_group_id;unique_index:idx_account_group,priority:2"`
	GroupName  string         `gorm:"column:group_name;type:varchar(63)"`
	SourceType string         `gorm:"column:source_type;type:varchar(63)"`

	AccountID       common.Hash    `gorm:"column:account_id;type:BINARY(32);unique_index:idx_account_group,priority:1"`
	OperatorAddress common.Address `gorm:"column:operator_address;type:BINARY(20)"`

	CreateAt   int64 `gorm:"column:create_at"`
	CreateTime int64 `gorm:"column:create_time"`
	UpdateAt   int64 `gorm:"column:update_at"`
	UpdateTime int64 `gorm:"column:update_time"`
	Removed    bool  `gorm:"column:removed;default:false"`
}

func (*Group) TableName() string {
	return "groups"
}
