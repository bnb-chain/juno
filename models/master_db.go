package models

// MasterDB stores current master DB
type MasterDB struct {
	OneRowId bool `gorm:"one_row_id;not null;default:true;primaryKey"`
	// IsMaster defines if current DB is master DB
	IsMaster bool `gorm:"column:is_master;not null;default:true;"`
}

// TableName is used to set Master table name in database
func (m *MasterDB) TableName() string {
	return "master_db"
}
