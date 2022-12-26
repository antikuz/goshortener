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
	db     *sql.DB
	logger *zap.Logger
}

func NewStorage(dbPath string, logger *zap.Logger) *Storage {
	db, err := connectDB(dbPath)
	if err != nil {
		logger.Sugar().Fatalf("Cannot create connect to DB, due to error: %v", err)
	}

	return &Storage{
		db:     db,
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
			secret_key TEXT,
			target_url TEXT NOT NULL,
			is_active BOOLEAN NOT NULL,
			clicks INTEGER NOT NULL,
			expires DATE
		);
		CREATE UNIQUE INDEX idx_urls_key ON urls (key);
		CREATE INDEX idx_urls_secret_key ON urls (secret_key);
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

func (s *Storage) AddShortURL(url models.ShortenURL) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("insert into urls values(?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(url.Key, url.Secret_key, url.Target_url, url.Is_active, url.Clicks, url.Expires)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetURL(urlHash string) (models.ShortenURL, error) {
	row := s.db.QueryRow("SELECT * FROM urls WHERE key = ?", urlHash)
	urlModel := models.ShortenURL{}
	if err := row.Scan(
		&urlModel.Key,
		&urlModel.Secret_key,
		&urlModel.Target_url,
		&urlModel.Is_active,
		&urlModel.Clicks,
		&urlModel.Expires,
	); err == sql.ErrNoRows {
		return models.ShortenURL{}, err
	}

	return urlModel, nil
}

func (s *Storage) DeleteURL(urlHash string) error {
	_, err := s.db.Exec("DELETE FROM urls WHERE key = ?", urlHash)
	return err
}

func (s *Storage) GetExpiredURL() ([]models.ShortenURL, error) {
	rows, err := s.db.Query("SELECT * FROM urls WHERE expires < date('now');")
	if err != nil {
		return []models.ShortenURL{}, err
	}
	defer rows.Close()

	urlModels := []models.ShortenURL{}
	for rows.Next() {
		m := models.ShortenURL{}
		if err = rows.Scan(
			&m.Key,
			&m.Secret_key,
			&m.Target_url,
			&m.Is_active,
			&m.Clicks,
			&m.Expires,
		); err != nil {
			return []models.ShortenURL{}, err
		}
		urlModels = append(urlModels, m)
	}

	return urlModels, nil
}

func (s *Storage) Close() {
	if err := s.db.Close(); err != nil {
		s.logger.Sugar().Errorf("Cant close database, due to err: %v", err)
	}
}