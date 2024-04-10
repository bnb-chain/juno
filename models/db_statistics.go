package models

type DataStat struct {
	OneRowId         bool   `gorm:"one_row_id;not null;default:true;primaryKey"`
	BlockHeight      int64  `gorm:"column:block_height;type:bigint(64)"`
	ObjectTotalCount string `gorm:"column:object_total_count;type:VARCHAR(2048)"`
	ObjectSealCount  string `gorm:"column:object_seal_count;type:VARCHAR(2048)"`
	ObjectDelCount   string `gorm:"column:object_del_count;type:VARCHAR(2048)"`
	UpdateTime       int64  `gorm:"update_time;type:bigint(64)"`
}

func (*DataStat) TableName() string {
	return "data_stat"
}

type BlockResult struct {
	BlockHeight uint64 `gorm:"primaryKey"`
	Result      string `gorm:"column:result;type:mediumtext"`
}

func (*BlockResult) TableName() string {
	return "block_result"
}
