package v3

import (
	"github.com/forbole/juno/v4/database/postgresql"
	"github.com/jmoiron/sqlx"

	"github.com/forbole/juno/v4/database"
)

var _ database.Migrator = &Migrator{}

// Migrator represents the database migrator that should be used to migrate from v2 of the database to v3
type Migrator struct {
	SQL *sqlx.DB
}

func NewMigrator(db *postgresql.Database) *Migrator {
	return &Migrator{
		//TODO adapt migrator
		//SQL: db.SQL,
	}
}
