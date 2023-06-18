package repository

import (
	"database/sql"
	"log"
	"testing"

	"github.com/begenov/backend/internal/config"
)

var db *sql.DB

func TestMain(m *testing.M) {
	var err error

	cfg, err := config.Init("../..")
	if err != nil {
		log.Fatalln(err)
	}

	db, err = sql.Open(cfg.Postgres.Driver, cfg.Postgres.DSN)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}
	transferRepo = NewTransferRepo(db)
	repo = New(db)
	entryRepo = NewEntryRepo(db)
	userRepo = NewUserRepo(db)

}
