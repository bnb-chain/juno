package pruning

import (
	"github.com/forbole/juno/v4/database"
	"github.com/forbole/juno/v4/modules"
	"github.com/forbole/juno/v4/types/config"
)

var (
	_ modules.Module                     = &Module{}
	_ modules.BlockModule                = &Module{}
	_ modules.AdditionalOperationsModule = &Module{}
)

// Module represents the pruning module allowing to clean the database periodically
type Module struct {
	cfg *Config
	db  database.Database
}

// NewModule builds a new Module instance
func NewModule(cfg config.Config, db database.Database) *Module {
	bz, err := cfg.GetBytes()
	if err != nil {
		panic(err)
	}

	pruningCfg, err := ParseConfig(bz)
	if err != nil {
		panic(err)
	}

	return &Module{
		cfg: pruningCfg,
		db:  db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "pruning"
}

// RunAdditionalOperations implements
func (m *Module) RunAdditionalOperations() error {
	return RunAdditionalOperations(m.cfg)
}
