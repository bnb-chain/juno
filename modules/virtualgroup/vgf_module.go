package virtualgroup

import (
	"context"

	"gorm.io/gorm/schema"

	"github.com/forbole/juno/v4/database"
	"github.com/forbole/juno/v4/models"
	"github.com/forbole/juno/v4/modules"
)

const (
	VGFModuleName = "global_virtual_group_family"
)

var (
	_ modules.Module              = &VGFModule{}
	_ modules.PrepareTablesModule = &VGFModule{}
)

// VGFModule represents the payment module
type VGFModule struct {
	db database.Database
}

// NewVGFModule builds a new Module instance
func NewVGFModule(db database.Database) *VGFModule {
	return &VGFModule{
		db: db,
	}
}

// Name implements modules.Module
func (m *VGFModule) Name() string {
	return VGFModuleName
}

// PrepareTables implements
func (m *VGFModule) PrepareTables() error {
	return m.db.PrepareTables(context.TODO(), []schema.Tabler{&models.GlobalVirtualGroupFamily{}})
}

// AutoMigrate implements
func (m *VGFModule) AutoMigrate() error {
	return m.db.AutoMigrate(context.TODO(), []schema.Tabler{&models.GlobalVirtualGroupFamily{}})
}
