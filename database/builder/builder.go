package builder

import (
	"errors"

	"github.com/forbole/juno/v4/database"
	databaseconfig "github.com/forbole/juno/v4/database/config"
	"github.com/forbole/juno/v4/database/mysql"
	"github.com/forbole/juno/v4/database/postgresql"
)

// Builder represents a generic Builder implementation that build the proper database
// instance based on the configuration the user has specified
func Builder(ctx *database.Context) (database.Database, error) {
	switch ctx.Cfg.Type {
	case databaseconfig.PostgreSQL:
		return postgresql.Builder(ctx)
	case databaseconfig.MySQL:
		return mysql.Builder(ctx)
	default:
		return nil, errors.New("unsupported database type")
	}
}
