package database

import (
	"database/sql"
	"log"
	"sync"

	_ "github.com/glebarez/sqlite"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/config"
)

type Database struct {
	instance *sql.DB
	config   *config.Config
	once     sync.Once
}

var dbInstance *Database

func New(config *config.Config) *Database {
	if dbInstance == nil {
		dbInstance = &Database{
			config: config,
		}
	}

	return dbInstance
}

func (database *Database) Init() {
	database.once.Do(func() {
		db, err := sql.Open("sqlite", database.config.RootPath+"/metadata.db")

		if err != nil {
			log.Fatalf("Failed to connect to sqlite database, details: %v", err)
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
			log.Fatalf("Failed to create tables, details: %v", err)
		}

		database.instance = db
	})
}

func GetInstance() *sql.DB {

	if dbInstance.instance == nil {
		log.Fatal("Database not initialized")
	}

	return dbInstance.instance
}
