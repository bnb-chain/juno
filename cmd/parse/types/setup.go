package types

import (
	"fmt"
	"os"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/forbole/juno/v4/database"
	"github.com/forbole/juno/v4/log"
	modsregistrar "github.com/forbole/juno/v4/modules/registrar"
	nodebuilder "github.com/forbole/juno/v4/node/builder"
	"github.com/forbole/juno/v4/parser"
	"github.com/forbole/juno/v4/types/config"
	"github.com/forbole/juno/v4/types/env"
)

// GetParserContext setups all the things that can be used to later parse the chain state
func GetParserContext(cfg config.Config, parseConfig *Config) (*parser.Context, error) {
	// Build the codec
	encodingConfig := parseConfig.GetEncodingConfigBuilder()()

	// Set up the SDK configuration
	sdkConfig, sealed := getConfig()
	if !sealed {
		parseConfig.GetSetupConfig()(cfg, sdkConfig)
		sdkConfig.Seal()
	}

	// Get the db context
	databaseCtx := database.NewContext(cfg.Database, &encodingConfig)

	// Try to get dsn from ENV first
	if dsn, errOfEnv := getDBConfigFromEnv(env.BlockSyncerDsn); errOfEnv != nil {
		log.Warn("load block syncer db config from ENV failed, try to use config from file")
	} else {
		databaseCtx.Cfg.DSN = dsn
		log.Infof("Using DB config from ENV")
	}

	// Init the db
	db, err := parseConfig.GetDBBuilder()(databaseCtx)
	if err != nil {
		return nil, err
	}

	// Init the client
	cp, err := nodebuilder.BuildNode(cfg.Node, &encodingConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to start client: %s", err)
	}

	// Setup the logging
	lvl, _ := log.ParseLevel(cfg.Logging.Level)
	log.Init(lvl, log.StandardizePath(cfg.Logging.RootDir, cfg.Logging.ServiceName))

	// Get the modules
	context := modsregistrar.NewContext(cfg, sdkConfig, &encodingConfig, db, cp)
	mods := parseConfig.GetRegistrar().BuildModules(context)
	registeredModules := modsregistrar.GetModules(mods, cfg.Chain.Modules)

	return parser.NewContext(&encodingConfig, cp, db, registeredModules), nil
}

// getConfig returns the SDK Config instance as well as if it's sealed or not
func getConfig() (config *sdk.Config, sealed bool) {
	sdkConfig := sdk.GetConfig()
	fv := reflect.ValueOf(sdkConfig).Elem().FieldByName("sealed")
	return sdkConfig, fv.Bool()
}

func getDBConfigFromEnv(dsn string) (string, error) {
	dsnVal, ok := os.LookupEnv(dsn)
	if !ok {
		return "", fmt.Errorf("dsn %s config is not set in environment", dsnVal)
	}
	return dsnVal, nil
}
