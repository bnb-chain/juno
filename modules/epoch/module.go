package epoch

import (
	"context"
	"github.com/forbole/juno/v4/database"
	"github.com/forbole/juno/v4/models"
	"github.com/forbole/juno/v4/modules"
	"gorm.io/gorm/schema"
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
