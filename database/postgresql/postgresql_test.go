package postgresql_test

import (
	"io/ioutil"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/stretchr/testify/suite"

	"github.com/forbole/juno/v4/database"
	databaseconfig "github.com/forbole/juno/v4/database/config"
	postgres "github.com/forbole/juno/v4/database/postgresql"
)

func TestDatabaseTestSuite(t *testing.T) {
	suite.Run(t, new(DbTestSuite))
}

type DbTestSuite struct {
	suite.Suite

	database *postgres.Database
}

func (suite *DbTestSuite) SetupTest() {
	// Create the codec
	codec := simapp.MakeTestEncodingConfig()

	// Build the database
	dbCfg := databaseconfig.NewDatabaseConfig(
		"postgres://bdjuno:password@localhost:6433/bdjuno?sslmode=disable&search_path=public",
		-1,
		-1,
		100000,
		100,
	)
	db, err := postgres.Builder(database.NewContext(dbCfg, &codec))
	suite.Require().NoError(err)

	bigDipperDb, ok := (db).(*postgres.Database)
	suite.Require().True(ok)

	// Delete the public schema
	err = bigDipperDb.Db.Exec(`DROP SCHEMA public CASCADE;`).Error
	suite.Require().NoError(err)

	// Re-create the schema
	err = bigDipperDb.Db.Exec(`CREATE SCHEMA public;`).Error
	suite.Require().NoError(err)

	dirPath := path.Join(".")
	dir, err := ioutil.ReadDir(dirPath)
	suite.Require().NoError(err)

	for _, fileInfo := range dir {
		if !strings.Contains(fileInfo.Name(), ".sql") {
			continue
		}

		file, err := ioutil.ReadFile(filepath.Join(dirPath, fileInfo.Name()))
		suite.Require().NoError(err)

		commentsRegExp := regexp.MustCompile(`/\*.*\*/`)
		requests := strings.Split(string(file), ";")
		for _, request := range requests {
			err := bigDipperDb.Db.Exec(commentsRegExp.ReplaceAllString(request, "")).Error
			suite.Require().NoError(err)
		}
	}

	suite.database = bigDipperDb
}
