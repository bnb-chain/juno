package tomlconfig

import (
	databaseconfig "github.com/forbole/juno/v4/database/config"
	loggingconfig "github.com/forbole/juno/v4/log/config"
	"github.com/forbole/juno/v4/node/remote"
	parserconfig "github.com/forbole/juno/v4/parser/config"
	"github.com/forbole/juno/v4/types/config"
)

type TomlConfig struct {
	Chain          config.ChainConfig
	Node           NodeConfig
	Parser         parserconfig.Config
	Database       databaseconfig.Config
	Logging        loggingconfig.Config
	RecreateTables bool
	Backup         bool
	DsnBackup      string
}

type NodeConfig struct {
	Type string
	RPC  *remote.RPCConfig
	GRPC *remote.GRPCConfig
}
