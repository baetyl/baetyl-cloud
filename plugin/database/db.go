package database

import (
	"database/sql"
	"io"
	"time"

	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/baetyl/baetyl-go/v2/log"
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
	Log *log.Logger
}

func init() {
	plugin.RegisterFactory("database", New)
}

// New New
func New() (plugin.Plugin, error) {
	var cfg CloudConfig
	if err := common.LoadConfig(&cfg); err != nil {
		return nil, errors.Trace(err)
	}
	return NewDB(cfg)
}

func NewDB(cfg CloudConfig) (*DB, error) {
	db, err := sqlx.Open(cfg.Database.Type, cfg.Database.URL)
	if err != nil {
		return nil, errors.Trace(err)
	}
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetMaxOpenConns(cfg.Database.MaxConns)
	db.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime) * time.Second)
	err = db.Ping()
	if err != nil {
		return nil, errors.Trace(err)
	}
	return &DB{
		db:  db,
		cfg: cfg,
		Log: log.With(log.Any("plugin", "database")),
	}, nil
}

// Close Close
func (d *DB) Close() (err error) {
	err = d.db.Close()
	err = errors.Trace(err)
	return
}

func (d *DB) Transact(handler func(*sqlx.Tx) error) (err error) {
	tx, err := d.db.Beginx()
	if err != nil {
		err = errors.Trace(err)
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
	err = errors.Trace(err)
	return
}

func (d *DB) Exec(tx *sqlx.Tx, sql string, args ...interface{}) (res sql.Result, err error) {
	if tx == nil {
		res, err = d.db.Exec(sql, args...)
	} else {
		res, err = tx.Exec(sql, args...)
	}
	err = errors.Trace(err)
	return
}

func (d *DB) Query(tx *sqlx.Tx, sql string, data interface{}, args ...interface{}) (err error) {
	if tx == nil {
		err = d.db.Select(data, sql, args...)
	} else {
		err = tx.Select(data, sql, args...)
	}
	err = errors.Trace(err)
	return

}
