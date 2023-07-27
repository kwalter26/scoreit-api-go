package db

import (
	"database/sql"
	"github.com/kwalter26/scoreit-api-go/util"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB
var testStore Store

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..", true)
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)

	testStore = NewStore(testDB)

	os.Exit(m.Run())
}
