package object

import (
	"context"

	"gorm.io/gorm/schema"

	"github.com/forbole/juno/v4/database"
	"github.com/forbole/juno/v4/models"
	"github.com/forbole/juno/v4/modules"
)

const (
	ModuleName = "object"
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
func (o *Module) Name() string {
	return ModuleName
}

// PrepareTables implements
func (o *Module) PrepareTables() error {
	return o.db.PrepareTables(context.TODO(), []schema.Tabler{&models.Object{}})
}
