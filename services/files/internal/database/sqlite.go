package database

import (
	"database/sql"
	"fmt"

	_ "github.com/glebarez/sqlite"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/config"
)

func NewSQLite(config *config.Config) (*sql.DB, error) {
	db, err := sql.Open("sqlite", config.RootPath+"/metadata.db")

	if err != nil {
		return nil, fmt.Errorf("Failed to connect to sqlite database, details: %v", err)
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	statement := `
		CREATE TABLE IF NOT EXISTS files (
			id TEXT NOT NULL PRIMARY KEY,
			filename TEXT NOT NULL,
			location TEXT NOT NULL,
			size INTEGER NOT NULL,
			user_id TEXT NOT NULL
		)
	`

	_, err = db.Exec(statement)

	if err != nil {
		return nil, fmt.Errorf("Failed to create tables, details: %v", err)
	}

	return db, nil
}
