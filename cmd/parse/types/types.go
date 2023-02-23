package types

import (
	"github.com/cosmos/cosmos-sdk/simapp"

	tomlconfig "github.com/forbole/juno/v4/cmd/migrate/toml"
	"github.com/forbole/juno/v4/database"
	"github.com/forbole/juno/v4/database/builder"
	"github.com/forbole/juno/v4/modules/registrar"
	"github.com/forbole/juno/v4/types/config"
)

// Config contains all the configuration for the "parse" command
type Config struct {
	registrar             registrar.Registrar
	configParser          config.Parser
	encodingConfigBuilder EncodingConfigBuilder
	setupCfg              SdkConfigSetup
	buildDb               database.Builder
	fileType              string
	tomlConfig            *tomlconfig.TomlConfig
}

// NewConfig allows to build a new Config instance
func NewConfig() *Config {
	return &Config{}
}

// WithRegistrar sets the modules registrar to be used
func (cfg *Config) WithRegistrar(r registrar.Registrar) *Config {
	cfg.registrar = r
	return cfg
}

// WithFileType sets the config type to be used
func (cfg *Config) WithFileType(fileType string) *Config {
	cfg.fileType = fileType
	return cfg
}

// GetRegistrar returns the modules registrar to be used
func (cfg *Config) GetRegistrar() registrar.Registrar {
	if cfg.registrar == nil {
		return &registrar.EmptyRegistrar{}
	}
	return cfg.registrar
}

// WithConfigParser sets the configuration parser to be used
func (cfg *Config) WithConfigParser(p config.Parser) *Config {
	cfg.configParser = p
	return cfg
}

// WithTomlConfig sets the tomlConfig
func (cfg *Config) WithTomlConfig(tomlConfig *tomlconfig.TomlConfig) *Config {
	cfg.tomlConfig = tomlConfig
	return cfg
}

// GetConfigParser returns the configuration parser to be used
func (cfg *Config) GetConfigParser(fileType string) config.Parser {
	if cfg.configParser == nil {
		switch fileType {
		case config.YamlConfigType:
			return config.DefaultConfigParser
		case config.TomlConfigType:
			return config.TomlConfigParser
		}
		return config.DefaultConfigParser
	}
	return cfg.configParser
}

// WithEncodingConfigBuilder sets the configurations builder to be used
func (cfg *Config) WithEncodingConfigBuilder(b EncodingConfigBuilder) *Config {
	cfg.encodingConfigBuilder = b
	return cfg
}

// GetEncodingConfigBuilder returns the encoding config builder to be used
func (cfg *Config) GetEncodingConfigBuilder() EncodingConfigBuilder {
	if cfg.encodingConfigBuilder == nil {
		return simapp.MakeTestEncodingConfig
	}
	return cfg.encodingConfigBuilder
}

// WithSetupConfig sets the SDK setup configurator to be used
func (cfg *Config) WithSetupConfig(s SdkConfigSetup) *Config {
	cfg.setupCfg = s
	return cfg
}

// GetSetupConfig returns the SDK configuration builder to use
func (cfg *Config) GetSetupConfig() SdkConfigSetup {
	if cfg.setupCfg == nil {
		return DefaultConfigSetup
	}
	return cfg.setupCfg
}

// WithDBBuilder sets the database builder to be used
func (cfg *Config) WithDBBuilder(b database.Builder) *Config {
	cfg.buildDb = b
	return cfg
}

// GetDBBuilder returns the database builder to be used
func (cfg *Config) GetDBBuilder() database.Builder {
	if cfg.buildDb == nil {
		return builder.Builder
	}
	return cfg.buildDb
}
