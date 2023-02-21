package config

import (
	"path"
	"time"
)

var (
	HomePath = ""
)

// GetConfigFilePath returns the path to the configuration file given the executable name
func GetConfigFilePath(fileType string) string {
	switch fileType {
	case "toml":
		return path.Join(HomePath, "config.toml")
	case "yaml":
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
