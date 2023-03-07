package parser

import (
	"github.com/bnb-chain/greenfield/app/params"
	"github.com/forbole/juno/v4/database"
	"github.com/forbole/juno/v4/modules"
	"github.com/forbole/juno/v4/node"
)

// Context represents the context that is shared among different workers
type Context struct {
	EncodingConfig *params.EncodingConfig
	Node           node.Node
	Database       database.Database
	Modules        []modules.Module
}

// NewContext builds a new Context instance
func NewContext(
	encodingConfig *params.EncodingConfig, proxy node.Node, db database.Database, modules []modules.Module,
) *Context {
	return &Context{
		EncodingConfig: encodingConfig,
		Node:           proxy,
		Database:       db,
		Modules:        modules,
	}
}
