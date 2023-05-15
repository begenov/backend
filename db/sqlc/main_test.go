package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	testDriver = "postgres"
	testDSN    = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries

var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(testDriver, testDSN)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}
	testQueries = New(testDB)
	os.Exit(m.Run())
}
