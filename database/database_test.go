package database

import (
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/forbole/juno/v4/common"
	"github.com/forbole/juno/v4/models"

	//"gorm.io/driver/postgres"
)

type DBType struct {
	dsn    string
	dbType string
	config gorm.Config
}

type Block struct {
	ID uint64 `gorm:"column:id;primaryKey" json:"-"`

	Hash     common.Hash `gorm:"column:hash;uniqueIndex:idx_hash"`
	Height   uint64      `gorm:"column:height;uniqueIndex:idx_height"`
	NumTxs   int         `gorm:"column:num_txs"`
	TotalGas uint64      `gorm:"column:total_gas"`
}

func (*Block) TableName() string {
	return "block"
}

type store struct {
	db *gorm.DB
}

func SetUpGorm(dbType *DBType) (*store, error) {

	var db *gorm.DB

	switch dbType.dbType {
	case "mysql":
		db, _ = gorm.Open(mysql.Open(dbType.dsn), &gorm.Config{})
	case "postgres":
		db, _ = gorm.Open(postgres.Open(dbType.dsn), &gorm.Config{})
	}

	db.Migrator().AutoMigrate(&Block{})

	return &store{db: db}, nil

	//create table

}

func TestMySQL(t *testing.T) {
	dbType := &DBType{
		dsn:    "root:Zxl_212901@tcp(localhost:3306)/test_gorm?parseTime=true&multiStatements=true&loc=Local",
		dbType: "mysql",
		config: gorm.Config{},
	}
	dbMySql, err := SetUpGorm(dbType)
	if err != nil {
		t.Log("init failed")
	}
	res, err := dbMySql.HasBlock(0)
	t.Log(res)

}

func TestPostgreSQL(t *testing.T) {
	dbType := &DBType{
		dsn:    "postgres://postgres:2129zxl2129@localhost:5432/test-gorm?sslmode=disable&search_path=public",
		dbType: "postgres",
		config: gorm.Config{},
	}
	dbMySql, err := SetUpGorm(dbType)
	if err != nil {
		t.Log("init failed")
	}
	res, err := dbMySql.HasBlock(0)
	t.Log(res)
}

func (storeDB *store) HasBlock(height int64) (bool, error) {
	var res bool

	var block models.Block
	if err := storeDB.db.Table("block").Where("height = ?", height).
		First(&block).
		Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return res, nil
}
