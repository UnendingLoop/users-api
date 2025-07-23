package config

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func ConnectDB(path string) *sqlx.DB {
	db, err := sqlx.Open("sqlite3", path)
	if err != nil {
		log.Fatalf("Cannot open db: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Cannot connect to db: %v", err)
	}

	return db
}
