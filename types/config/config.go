package config

import (
	"strings"

	databaseconfig "github.com/forbole/juno/v4/database/config"
	loggingconfig "github.com/forbole/juno/v4/log/config"
	nodeconfig "github.com/forbole/juno/v4/node/config"
	"github.com/forbole/juno/v4/node/remote"
	parserconfig "github.com/forbole/juno/v4/parser/config"
)

var (
	// Cfg represents the configuration to be used during the execution
	Cfg Config
)

// Config defines all necessary juno configuration parameters.
type Config struct {
	bytes []byte

	Chain    ChainConfig           `yaml:"chain"`
	Node     nodeconfig.Config     `yaml:"node"`
	Parser   parserconfig.Config   `yaml:"parsing"`
	Database databaseconfig.Config `yaml:"database"`
	Logging  loggingconfig.Config  `yaml:"logging"`
}

// NewConfig builds a new Config instance
func NewConfig(
	nodeCfg nodeconfig.Config,
	chainCfg ChainConfig, dbConfig databaseconfig.Config,
	parserConfig parserconfig.Config, loggingConfig loggingconfig.Config,
) Config {
	return Config{
		Node:     nodeCfg,
		Chain:    chainCfg,
		Database: dbConfig,
		Parser:   parserConfig,
		Logging:  loggingConfig,
	}
}

func DefaultConfig() Config {
	return NewConfig(
		nodeconfig.DefaultConfig(),
		DefaultChainConfig(), databaseconfig.DefaultDatabaseConfig(),
		parserconfig.DefaultParsingConfig(), loggingconfig.DefaultLogConfig(),
	)
}

func (c Config) GetBytes() ([]byte, error) {
	return c.bytes, nil
}

// ---------------------------------------------------------------------------------------------------------------------

type ChainConfig struct {
	Bech32Prefix string   `yaml:"bech32_prefix"`
	Modules      []string `yaml:"modules"`
}

// NewChainConfig returns a new ChainConfig instance
func NewChainConfig(bech32Prefix string, modules []string) ChainConfig {
	return ChainConfig{
		Bech32Prefix: bech32Prefix,
		Modules:      modules,
	}
}

// DefaultChainConfig returns the default instance of ChainConfig
func DefaultChainConfig() ChainConfig {
	return NewChainConfig("cosmos", nil)
}

func (cfg ChainConfig) IsModuleEnabled(moduleName string) bool {
	for _, module := range cfg.Modules {
		if strings.EqualFold(module, moduleName) {
			return true
		}
	}

	return false
}

type TomlConfig struct {
	Chain          ChainConfig
	Node           NodeConfig
	Parser         parserconfig.Config
	Database       databaseconfig.Config
	Logging        loggingconfig.Config
	RecreateTables bool
	EnableDualDB   bool
	DsnSwitched    string
}

type NodeConfig struct {
	Type string
	RPC  *remote.RPCConfig
	GRPC *remote.GRPCConfig
}
