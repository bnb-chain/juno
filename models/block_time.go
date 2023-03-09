package models

type AverageBlockTimePerMinute struct {
	OneRowId    bool    `gorm:"one_row_id;not null;default:true;primaryKey"`
	AverageTime float64 `gorm:"average_time;not null"`
	Height      uint64  `gorm:"height;not null;index:idx_height"`
}

func (*AverageBlockTimePerMinute) TableName() string {
	return "average_block_time_per_minute"
}

type AverageBlockTimePerHour struct {
	OneRowId    bool    `gorm:"one_row_id;not null;default:true;primaryKey"`
	AverageTime float64 `gorm:"average_time;not null"`
	Height      uint64  `gorm:"height;not null;index:idx_height"`
}

func (*AverageBlockTimePerHour) TableName() string {
	return "average_block_time_per_hour"
}

type AverageBlockTimePerDay struct {
	OneRowId    bool    `gorm:"one_row_id;not null;default:true;primaryKey"`
	AverageTime float64 `gorm:"average_time;not null"`
	Height      uint64  `gorm:"height;not null;index:idx_height"`
}

func (*AverageBlockTimePerDay) TableName() string {
	return "average_block_time_per_day"
}

type AverageBlockTimeFromGenesis struct {
	OneRowId    bool    `gorm:"one_row_id;not null;default:true;primaryKey"`
	AverageTime float64 `gorm:"average_time;not null"`
	Height      uint64  `gorm:"height;not null;index:idx_height"`
}

func (*AverageBlockTimeFromGenesis) TableName() string {
	return "average_block_time_from_genesis"
}
