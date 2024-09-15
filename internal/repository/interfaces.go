package repository

import (
	"database/sql"
	_ "github.com/lib/pq"
	"tender-manager/config"
)

type PostgresDB struct {
	db *sql.DB
}

func New(cfg config.Config) (*PostgresDB, error) {
	db, err := sql.Open("postgres", cfg.PostgresConn)
	if err != nil {
		return &PostgresDB{}, err
	}

	return &PostgresDB{db: db}, nil
}
