package models

import "github.com/forbole/juno/v4/common"

type Group struct {
	Owner      common.Address `gorm:"owner;type:BINARY(20);index:idx_owner"`
	GroupID    int64          `gorm:"group_id;type:int;primaryKey"`
	GroupName  string         `gorm:"group_name;type:varchar(63)"`
	SourceType int            `gorm:"source_type;type:int"`
}

func (*Group) TableName() string {
	return "groups"
}
