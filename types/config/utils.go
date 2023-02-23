package config

import (
	"path"
	"time"
)

var (
	HomePath = ""
)

const (
	TomlConfigType        = "toml"
	YamlConfigType        = "yaml"
	BlockSyncerWorkerType = "blocksyncer"
	ExplorerWorkerType    = "explorer"
)

// GetConfigFilePath returns the path to the configuration file given the executable name
func GetConfigFilePath(configType string) string {
	switch configType {
	case TomlConfigType:
		return path.Join(HomePath, "config.toml")
	case YamlConfigType:
		return path.Join(HomePath, "config.yaml")
	}
	return path.Join(HomePath, "config.yaml")
}

// GetAvgBlockTime returns the average_block_time in the configuration file or
// returns 3 seconds if it is not configured
func GetAvgBlockTime() time.Duration {
	if Cfg.Parser.AvgBlockTime == nil {
		return 3 * time.Second
	}
	return *Cfg.Parser.AvgBlockTime
}
