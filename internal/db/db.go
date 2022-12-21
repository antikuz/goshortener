package db

import (
	"database/sql"
	"errors"
	"log"
	"os"

	"github.com/antikuz/goshortener/internal/models"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

type Storage struct {
	db *sql.DB
	logger *zap.Logger
}

func NewStorage(dbPath string, logger *zap.Logger) *Storage {
	db, err := connectDB(dbPath)
	if err != nil {
		logger.Sugar().Fatalf("Cannot create connect to DB, due to error: %v", err)
	}

	return &Storage{
		db: db,
		logger: logger,
	}
}

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
			secret_key TEXT UNIQUE,
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
			return nil, err
		}

	}
	return db, nil
}

func (s *Storage) NewURL(url models.ShortenURL) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("insert into urls values(?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(url.Key, url.Secret_key, url.Target_url, url.Is_active, url.Clicks)
	if err != nil {
		return err
	}
	
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) Close() {
	if err := s.db.Close(); err != nil {
		s.logger.Sugar().Errorf("Cant close database, due to err: %v", err)
	}
}