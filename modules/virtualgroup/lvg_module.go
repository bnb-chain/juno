package virtualgroup

import (
	"context"

	"gorm.io/gorm/schema"

	"github.com/forbole/juno/v4/database"
	"github.com/forbole/juno/v4/models"
	"github.com/forbole/juno/v4/modules"
)

const (
	LVGModuleName = "local_virtual_group"
)

var (
	_ modules.Module              = &LVGModule{}
	_ modules.PrepareTablesModule = &LVGModule{}
)

// LVGModule represents the payment module
type LVGModule struct {
	db database.Database
}

// NewLVGModule builds a new Module instance
func NewLVGModule(db database.Database) *LVGModule {
	return &LVGModule{
		db: db,
	}
}

// Name implements modules.Module
func (m *LVGModule) Name() string {
	return LVGModuleName
}

// PrepareTables implements
func (m *LVGModule) PrepareTables() error {
	return m.db.PrepareTables(context.TODO(), []schema.Tabler{&models.LocalVirtualGroup{}})
}

// AutoMigrate implements
func (m *LVGModule) AutoMigrate() error {
	return m.db.AutoMigrate(context.TODO(), []schema.Tabler{&models.LocalVirtualGroup{}})
}
