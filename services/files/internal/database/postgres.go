package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type PostgresConfig struct {
	Db       string
	User     string
	Password string
	Host     string
}

func NewPostgres(config PostgresConfig) (*sql.DB, error) {
	connStr := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable", config.User, config.Password, config.Host, config.Db)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

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
