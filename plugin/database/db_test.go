package database

import (
	"github.com/baetyl/baetyl-cloud/common"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"testing"
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

func TestNew(t *testing.T) {
	p, err := New()
	assert.Error(t, err)
	assert.Nil(t, p)
	assert.EqualError(t, err, "open etc/baetyl/cloud.yml: no such file or directory")
	common.SetConfFile("../../scripts/config/cloud.yml")
	p, err = New()
	assert.Nil(t, err)
	assert.NotNil(t, p)
}
