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

type Postgres struct {
	Conn    *sql.DB
	Queries *Queries
}

func NewPostgres(config PostgresConfig) (*Postgres, error) {
	connStr := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable", config.User, config.Password, config.Host, config.Db)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return &Postgres{
		Conn:    db,
		Queries: New(db),
	}, nil
}
