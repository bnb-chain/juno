package registrar

import (
	"github.com/bnb-chain/greenfield/app/params"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/forbole/juno/v4/database"
	"github.com/forbole/juno/v4/log"
	"github.com/forbole/juno/v4/modules"
	"github.com/forbole/juno/v4/modules/block"
	"github.com/forbole/juno/v4/modules/bucket"
	"github.com/forbole/juno/v4/modules/epoch"
	"github.com/forbole/juno/v4/modules/group"
	"github.com/forbole/juno/v4/modules/messages"
	"github.com/forbole/juno/v4/modules/object"
	"github.com/forbole/juno/v4/modules/payment"
	"github.com/forbole/juno/v4/modules/permission"
	"github.com/forbole/juno/v4/modules/pruning"
	"github.com/forbole/juno/v4/modules/telemetry"
	"github.com/forbole/juno/v4/modules/validator"
	"github.com/forbole/juno/v4/node"
	"github.com/forbole/juno/v4/types/config"
)

// Context represents the context of the modules registrar
type Context struct {
	JunoConfig     config.Config
	SDKConfig      *sdk.Config
	EncodingConfig *params.EncodingConfig
	Database       database.Database
	Proxy          node.Node
}

// NewContext allows to build a new Context instance
func NewContext(
	parsingConfig config.Config, sdkConfig *sdk.Config, encodingConfig *params.EncodingConfig,
	database database.Database, proxy node.Node,
) Context {
	return Context{
		JunoConfig:     parsingConfig,
		SDKConfig:      sdkConfig,
		EncodingConfig: encodingConfig,
		Database:       database,
		Proxy:          proxy,
	}
}

// Registrar represents a module registrar. This allows to build a list of modules that can later be used by
// specifying their names inside the TOML configuration file.
type Registrar interface {
	BuildModules(context Context) modules.Modules
}

// ------------------------------------------------------------------------------------------------------------------

var (
	_ Registrar = &EmptyRegistrar{}
)

// EmptyRegistrar represents a Registrar which does not register any custom module
type EmptyRegistrar struct{}

// BuildModules implements Registrar
func (*EmptyRegistrar) BuildModules(_ Context) modules.Modules {
	return nil
}

// ------------------------------------------------------------------------------------------------------------------

var (
	_ Registrar = &DefaultRegistrar{}
)

// DefaultRegistrar represents a registrar that allows to handle the default Juno modules
type DefaultRegistrar struct {
	parser messages.MessageAddressesParser
}

// NewDefaultRegistrar builds a new DefaultRegistrar
func NewDefaultRegistrar(parser messages.MessageAddressesParser) *DefaultRegistrar {
	return &DefaultRegistrar{
		parser: parser,
	}
}

// BuildModules implements Registrar
func (r *DefaultRegistrar) BuildModules(ctx Context) modules.Modules {
	return modules.Modules{
		block.NewModule(ctx.Database),
		validator.NewModule(ctx.Database),
		bucket.NewModule(ctx.Database),
		group.NewModule(ctx.Database),
		object.NewModule(ctx.Database),
		pruning.NewModule(ctx.JunoConfig, ctx.Database),
		telemetry.NewModule(ctx.JunoConfig),
		epoch.NewModule(ctx.Database),
		payment.NewModule(ctx.Database),
		permission.NewModule(ctx.Database),
		group.NewModule(ctx.Database),
	}
}

// ------------------------------------------------------------------------------------------------------------------

// GetModules returns the list of module implementations based on the given module names.
// For each module name that is specified but not found, a warning log is printed.
func GetModules(mods modules.Modules, names []string) []modules.Module {
	var modulesImpls []modules.Module
	for _, name := range names {
		module, found := mods.FindByName(name)
		if found {
			modulesImpls = append(modulesImpls, module)
		} else {
			log.Errorw("Module is required but not registered. Be sure to register it using registrar.RegisterModule", "module", name)
		}
	}
	return modulesImpls
}
