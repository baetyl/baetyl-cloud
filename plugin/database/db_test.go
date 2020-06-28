package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func MockNewDB() (*dbStorage, error) {
	var cfg CloudConfig
	cfg.Database.Type = "sqlite3"
	cfg.Database.URL = ":memory:"
	db, err := sqlx.Open(cfg.Database.Type, cfg.Database.URL)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &dbStorage{db: db, cfg: cfg}, nil
}
