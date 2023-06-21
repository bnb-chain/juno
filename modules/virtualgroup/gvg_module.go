package virtualgroup

import (
	"context"

	"gorm.io/gorm/schema"

	"github.com/forbole/juno/v4/database"
	"github.com/forbole/juno/v4/models"
	"github.com/forbole/juno/v4/modules"
)

const (
	GVGModuleName = "global_virtual_group"
)

var (
	_ modules.Module              = &GVGModule{}
	_ modules.PrepareTablesModule = &GVGModule{}
)

// GVGModule represents the payment module
type GVGModule struct {
	db database.Database
}

// NewGVGModule builds a new Module instance
func NewGVGModule(db database.Database) *GVGModule {
	return &GVGModule{
		db: db,
	}
}

// Name implements modules.Module
func (m *GVGModule) Name() string {
	return GVGModuleName
}

// PrepareTables implements
func (m *GVGModule) PrepareTables() error {
	return m.db.PrepareTables(context.TODO(), []schema.Tabler{&models.GlobalVirtualGroup{}})
}

// AutoMigrate implements
func (m *GVGModule) AutoMigrate() error {
	return m.db.AutoMigrate(context.TODO(), []schema.Tabler{&models.GlobalVirtualGroup{}})
}
