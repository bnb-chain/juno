package group

import (
	"github.com/forbole/juno/v4/modules"
	"github.com/forbole/juno/v4/types/config"
)

const (
	ModuleName = "group"
)

var (
	_ modules.Module = &Module{}
)

// Module represents the telemetry module
type Module struct {
}

// NewModule returns a new Module implementation
func NewModule(cfg config.Config) *Module {
	//bz, err := cfg.GetBytes()
	//if err != nil {
	//	panic(err)
	//}
	//
	//telemetryCfg, err := ParseConfig(bz)
	//if err != nil {
	//	panic(err)
	//}
	//
	//return &Module{
	//	cfg: telemetryCfg,
	//}
	return nil
}

// Name implements modules.Module
func (module *Module) Name() string {
	return ModuleName
}

func (module *Module) PrepareTables() {

}
