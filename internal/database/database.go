package database

import (
	"database/sql"

	"github.com/davidtaing/go-webhook-server/internal/logger"
	_ "github.com/mattn/go-sqlite3"
)

func Open(path string, logger *logger.Logger) *sql.DB {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		logger.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		logger.Fatal(err)
	}

	return db
}
