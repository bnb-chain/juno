package types

import (
	"fmt"
	"os"

	tomlconfig "github.com/forbole/juno/v4/cmd/migrate/toml"
	nodeconfig "github.com/forbole/juno/v4/node/config"
	"github.com/forbole/juno/v4/node/remote"
	"github.com/forbole/juno/v4/types/config"

	"github.com/spf13/cobra"

	"github.com/forbole/juno/v4/types"
)

// ReadConfigPreRunE represents a Cobra cmd function allowing to read the config before executing the command itself
func ReadConfigPreRunE(cfg *Config) types.CobraCmdFunc {
	return func(_ *cobra.Command, _ []string) error {
		return UpdatedGlobalCfg(cfg)
	}
}

func NewParseConfigFromToml(tomlConfig *tomlconfig.TomlConfig) config.Config {
	config := config.Config{}
	config.Chain = tomlConfig.Chain
	config.Parser = tomlConfig.Parser
	config.Database = tomlConfig.Database
	config.Logging = tomlConfig.Logging

	nodeConfig := &nodeconfig.Config{}
	nodeConfig.Type = tomlConfig.Node.Type

	var nodeDetails = new(remote.Details)
	nodeDetails.RPC = tomlConfig.Node.RPC
	nodeDetails.GRPC = tomlConfig.Node.GRPC
	nodeConfig.Details = nodeDetails
	config.Node = *nodeConfig

	return config

}

// ReadConfig allows to read the configuration using the provided cfg
func ReadConfig(cfg *Config) (config.Config, error) {
	//if tomlConfig exist
	if cfg.tomlConfig != nil {
		return NewParseConfigFromToml(cfg.tomlConfig), nil
	}

	file := config.GetConfigFilePath(cfg.fileType)

	// Make sure the path exists
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return config.Config{}, fmt.Errorf("config file does not exist (%s). Make sure you have run the init command", file)
	}

	// Read the config
	return config.Read(file, cfg.GetConfigParser(cfg.fileType))
}

// UpdatedGlobalCfg parses the configuration file using the provided configuration and sets the
// parsed config as the global one
func UpdatedGlobalCfg(cfg *Config) error {
	junoCfg, err := ReadConfig(cfg)
	if err != nil {
		return err
	}

	// Set the global configuration
	config.Cfg = junoCfg
	return nil
}
