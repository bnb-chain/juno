package epoch

import (
	"context"

	"gorm.io/gorm/schema"

	"github.com/forbole/juno/v4/database"
	"github.com/forbole/juno/v4/models"
	"github.com/forbole/juno/v4/modules"
)

const (
	ModuleName = "epoch"
)

var (
	_ modules.Module              = &Module{}
	_ modules.PrepareTablesModule = &Module{}
)

// Module represents the telemetry module
type Module struct {
	db database.Database
}

// NewModule builds a new Module instance
func NewModule(db database.Database) *Module {
	return &Module{
		db: db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return ModuleName
}

// PrepareTables implements
func (m *Module) PrepareTables() error {
	return m.db.PrepareTables(context.TODO(), []schema.Tabler{&models.Epoch{}})
}

// RecreateTables implements
func (m *Module) RecreateTables() error {
	return m.db.RecreateTables(context.TODO(), []schema.Tabler{&models.Epoch{}})
}
