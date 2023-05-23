package database

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func MockNewDB() (*DB, error) {
	var cfg CloudConfig
	cfg.Database.Type = "sqlite3"
	cfg.Database.URL = ":memory:"
	db, err := sqlx.Open(cfg.Database.Type, cfg.Database.URL)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return &DB{db: db, Cfg: cfg}, nil
}
