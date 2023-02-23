package postgresql

import (
	"context"

	"github.com/forbole/juno/v4/database"
	"github.com/forbole/juno/v4/database/sqlclient"
)

// Builder creates a database connection with the given database connection info
// from config. It returns a database connection handle or an error if the
// connection fails.
func Builder(ctx *database.Context) (database.Database, error) {
	db, err := sqlclient.New(&ctx.Cfg)
	if err != nil {
		return nil, err
	}
	return &Database{
		Impl: database.Impl{
			Db:             db,
			EncodingConfig: ctx.EncodingConfig,
		},
	}, nil
}

// type check to ensure interface is properly implemented
var _ database.Database = &Database{}

// Database defines a wrapper around a SQL database and implements functionality
// for data aggregation and exporting.
type Database struct {
	database.Impl
}

// GetMissingHeights returns a slice of missing block heights between startHeight and endHeight
func (db *Database) GetMissingHeights(ctx context.Context, startHeight, endHeight uint64) []uint64 {
	var result []uint64
	stmt := `SELECT generate_series($1::int,$2::int) EXCEPT SELECT height FROM blocks ORDER BY 1;`
	err := db.Db.Select(&result, stmt, startHeight, endHeight)
	if err != nil {
		return nil
	}

	if len(result) == 0 {
		return nil
	}

	return result
}
