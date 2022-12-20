package db

import (
	"database/sql"
	"errors"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func connectDB(dbPath string) (*sql.DB, error) {
	var db *sql.DB
	if _, err := os.Stat(dbPath); errors.Is(err, os.ErrNotExist) {
		db, err = sql.Open("sqlite3", dbPath)
		if err != nil {
			return nil, err
		}
	
		sqlStmt := `
		CREATE TABLE urls (
			key TEXT NOT NULL UNIQUE,
			secret_key TEXT NOT NULL UNIQUE,
			target_url TEXT NOT NULL,
			is_active BOOLEAN NOT NULL,
			clicks INTEGER NOT NULL
		);
		CREATE UNIQUE INDEX idx_urls_key ON urls (key);
		CREATE UNIQUE INDEX idx_urls_secret_key ON urls (secret_key);
		`
		_, err := db.Exec(sqlStmt)
		if err != nil {
			log.Printf("%q: %s\n", err, sqlStmt)
			return nil, err
		}
	} else {
		db, err = sql.Open("sqlite3", dbPath)
		if err != nil {
			log.Fatal(err)
		}

	}
	return db, nil
}
func NewDB(dbPath string) *sql.DB {
	db, err := connectDB(dbPath)
	if err != nil {
		log.Fatal(db)
	}

	return db
}