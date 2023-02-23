package validator

import (
	"context"

	"github.com/forbole/juno/v4/database"
	"github.com/forbole/juno/v4/models"
	"github.com/forbole/juno/v4/modules"
	"gorm.io/gorm/schema"
)

var (
	_ modules.Module              = &Module{}
	_ modules.PrepareTablesModule = &Module{}
)

// Module represents the basic module which is required by both explorer and storage-provider
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
	return "validator"
}

// PrepareTables implements
func (m *Module) PrepareTables() error {
	return m.db.PrepareTables(context.TODO(), []schema.Tabler{
		&models.Validator{},
		&models.ValidatorInfo{},
		&models.ValidatorDescription{},
		&models.ValidatorCommission{},
		&models.ValidatorVotingPower{},
		&models.ValidatorStatus{},
		&models.ValidatorSigningInfo{}})
}
