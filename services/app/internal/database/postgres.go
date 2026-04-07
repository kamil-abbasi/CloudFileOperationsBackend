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

	createUsersTableStatement := `
		CREATE TABLE IF NOT EXISTS users (
			id UUID NOT NULL PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL
		)
	`

	createDirectoriesTableStatement := `
		CREATE TABLE IF NOT EXISTS directories (
			id UUID NOT NULL PRIMARY KEY,
			user_id UUID NOT NULL,
			parent_id UUID,
			name TEXT NOT NULL,
			location TEXT NOT NULL,
			FOREIGN KEY (parent_id)
				REFERENCES directories(id)
				ON DELETE CASCADE
		)
	`

	createFilesTableStatement := `
		CREATE TABLE IF NOT EXISTS files (
			id UUID NOT NULL PRIMARY KEY,
			user_id UUID NOT NULL,
			directory_id UUID,
			name TEXT NOT NULL,
			size INTEGER NOT NULL,
			location TEXT NOT NULL,
			FOREIGN KEY(directory_id)
				REFERENCES directories(id)
				ON DELETE CASCADE
		)
	`

	createItemTypeEnumStatement := `
		CREATE TYPE item_type AS ENUM('file', 'directory')
	`

	createDirectoryItemsTableStatement := `
		CREATE TABLE IF NOT EXISTS directory_items (
			id UUID NOT NULL PRIMARY KEY,
			user_id UUID NOT NULL,
			type item_type NOT NULL
		)
	`

	_, err = db.Exec(createUsersTableStatement)
	_, err = db.Exec(createDirectoriesTableStatement)
	_, err = db.Exec(createFilesTableStatement)
	_, err = db.Exec(createItemTypeEnumStatement)
	_, err = db.Exec(createDirectoryItemsTableStatement)

	if err != nil {
		return nil, fmt.Errorf("Failed to create tables, details: %v", err)
	}

	return db, nil
}
