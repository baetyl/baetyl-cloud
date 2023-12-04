package database

import (
	"github.com/baetyl/baetyl-go/v2/log"
	_ "github.com/mattn/go-sqlite3"
)

func MockNewDB() (*BaetylCloudDB, error) {
	var cfg CloudConfig
	cfg.Database.Type = "sqlite3"
	cfg.Database.URL = ":memory:"
	cfg.Database.ConnMaxLifetime = 150
	cfg.Database.MaxConns = 20
	cfg.Database.MaxIdleConns = 5
	db, err := NewDB(cfg)
	if err != nil {
		return nil, err
	}
	return &BaetylCloudDB{
		DB:  *db,
		log: log.L().With(log.Any("plugin", "test")),
	}, nil
}
