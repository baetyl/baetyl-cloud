package database

import (
	"database/sql"
	"io"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

// DBStorage
type DBStorage interface {
	Transact(func(*sqlx.Tx) error) error
	Exec(tx *sqlx.Tx, sql string, args ...interface{}) (sql.Result, error)
	Query(tx *sqlx.Tx, sql string, data interface{}, args ...interface{}) error

	io.Closer
}

// DBStorage
type DB struct {
	db  *sqlx.DB
	cfg CloudConfig
}

func init() {
	plugin.RegisterFactory("database", New)
}

// New New
func New() (plugin.Plugin, error) {
	var cfg CloudConfig
	if err := common.LoadConfig(&cfg); err != nil {
		return nil, err
	}
	return NewDB(cfg)
}

func NewDB(cfg CloudConfig) (*DB, error) {
	db, err := sqlx.Open(cfg.Database.Type, cfg.Database.URL)
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetMaxOpenConns(cfg.Database.MaxConns)
	db.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime) * time.Second)
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &DB{
		db:  db,
		cfg: cfg,
	}, nil
}

// Close Close
func (d *DB) Close() error {
	return d.db.Close()
}

func (d *DB) Transact(handler func(*sqlx.Tx) error) (err error) {
	tx, err := d.db.Beginx()
	if err != nil {
		return
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	err = handler(tx)
	return
}

func (d *DB) Exec(tx *sqlx.Tx, sql string, args ...interface{}) (sql.Result, error) {
	if tx == nil {
		return d.db.Exec(sql, args...)
	}
	return tx.Exec(sql, args...)
}

func (d *DB) Query(tx *sqlx.Tx, sql string, data interface{}, args ...interface{}) error {
	if tx == nil {
		return d.db.Select(data, sql, args...)
	}
	return tx.Select(data, sql, args...)
}
