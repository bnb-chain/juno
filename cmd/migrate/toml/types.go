package tomlconfig

import (
	databaseconfig "github.com/forbole/juno/v4/database/config"
	loggingconfig "github.com/forbole/juno/v4/log/config"
	"github.com/forbole/juno/v4/node/remote"
	parserconfig "github.com/forbole/juno/v4/parser/config"
	"github.com/forbole/juno/v4/types/config"
)

type TomlConfig struct {
	//Chain    ChainConfig
	//Node     NodeConfig
	//Parser   ParserConfig
	//Database DBConfig
	//Logging  LogConfig
	//// The following are there to support modules which config are present if they are enabled
	//
	//Telemetry TelemetryConfig
	//Pruning   PruningConfig
	//PriceFeed *pricefeedconfig.Config

	Chain    config.ChainConfig
	Node     NodeConfig
	Parser   parserconfig.Config
	Database databaseconfig.Config
	Logging  loggingconfig.Config
}

type NodeConfig struct {
	Type string
	RPC  *remote.RPCConfig
	GRPC *remote.GRPCConfig
}

//type RPCConfig struct {
//	ClientName     string
//	Address        string
//	MaxConnections int
//}
//
//type GRPCConfig struct {
//	Address  string
//	Insecure bool
//}

//type ChainConfig struct {
//	Bech32Prefix string
//	Modules      []string
//}

//type ParserConfig struct {
//	GenesisFilePath string
//	Workers         int64
//	StartHeight     int64
//	AvgBlockTime    *time.Duration
//	ParseNewBlocks  bool
//	ParseOldBlocks  bool
//	ParseGenesis    bool
//	FastSync        bool
//}

//type DBConfig struct {
//	Type               string
//	DSN                string
//	Secrets            *Params
//	SlowThreshold      time.Duration
//	MaxOpenConnections int
//	MaxIdleConnections int
//	ConnMaxIdleTime    time.Duration
//	ConnMaxLifetime    time.Duration
//	PartitionSize      int64
//	PartitionBatchSize int64
//}
//
//type LogConfig struct {
//	ServiceName string
//	RootDir     string
//	Level       string
//}
//
//type TelemetryConfig struct {
//	Port uint
//}
//
//type PruningConfig struct {
//	KeepRecent int64
//	KeepEvery  int64
//	Interval   int64
//}
//
//type Params struct {
//	SecretId string `yaml:"SecretId"`
//	Region   string `yaml:"Region"`
//}
