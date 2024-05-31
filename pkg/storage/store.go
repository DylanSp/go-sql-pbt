package storage

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // load Postgres driver
)

type DBConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	DBName   string
}

// wrapper around a sqlx.DB, used for performing database operations
type Store struct {
	db *sqlx.DB
}

func NewStore(cfg DBConfig) (*Store, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return &Store{
		db: db,
	}, nil
}
