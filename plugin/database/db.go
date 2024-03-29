package database

import (
	"context"
	"database/sql"
	"io"
	"strings"
	"time"

	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/baetyl/baetyl-go/v2/log"
	"github.com/jmoiron/sqlx"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

var (
	ErrDatabaseTx = errors.New("failed to transfer interface to sqlx.Tx")
)

// DBStorage
type DBStorage interface {
	Transact(func(*sqlx.Tx) error) error
	Exec(tx *sqlx.Tx, sql string, args ...interface{}) (sql.Result, error)
	Query(tx *sqlx.Tx, sql string, data interface{}, args ...interface{}) error
	BeginTx() (*sqlx.Tx, error)
	Commit(tx *sqlx.Tx)
	Rollback(tx *sqlx.Tx)

	io.Closer
}

// DBStorage
type DB struct {
	db  *sqlx.DB
	cfg CloudConfig
	Log *log.Logger
}

func init() {
	plugin.RegisterFactory("database", NewExtDB)
}

var (
	HookSQL = func(s string) string {
		return s
	}
)

// New New
func New() (plugin.Plugin, error) {
	var cfg CloudConfig
	if err := common.LoadConfig(&cfg); err != nil {
		return nil, errors.Trace(err)
	}
	return NewDB(cfg)
}

func NewDB(cfg CloudConfig) (*DB, error) {
	if cfg.Database.Decryption {
		decryptedURL, err := genDecryptedURL(cfg.Database.URL)
		if decryptedURL == "" || err != nil {
			return nil, errors.Trace(err)
		}
		cfg.Database.URL = decryptedURL
	}

	db, err := sqlx.Open(cfg.Database.Type, cfg.Database.URL)
	if err != nil {
		return nil, errors.Trace(err)
	}
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetMaxOpenConns(cfg.Database.MaxConns)
	db.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime) * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return &DB{
		db:  db,
		cfg: cfg,
		Log: log.With(log.Any("plugin", "database")),
	}, nil
}

func genDecryptedURL(originURL string) (string, error) {
	var decryptedURL string
	decrypt, err := plugin.GetPlugin("decryption")
	if err != nil {
		return "", errors.Trace(err)
	}
	dec, ok := decrypt.(plugin.Decrypt)
	if !ok {
		return "", errors.Trace(errors.New("plugin type conversion error"))
	}
	oriPassword := originURL[strings.Index(originURL, ":")+1 : strings.LastIndex(originURL, "@")]
	newPassword, err := dec.Decrypt(oriPassword)
	if newPassword == "" || err != nil {
		return "", errors.Trace(err)
	}
	decryptedURL = strings.Replace(originURL, oriPassword, newPassword, -1)

	return decryptedURL, nil
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
	sql = HookSQL(sql)
	if tx == nil {
		res, err = d.db.Exec(sql, args...)
	} else {
		res, err = tx.Exec(sql, args...)
	}
	err = errors.Trace(err)
	return
}

func (d *DB) Query(tx *sqlx.Tx, sql string, data interface{}, args ...interface{}) (err error) {
	sql = HookSQL(sql)
	if tx == nil {
		err = d.db.Select(data, sql, args...)
	} else {
		err = tx.Select(data, sql, args...)
	}
	err = errors.Trace(err)
	return
}

func (d *DB) BeginTx() (*sqlx.Tx, error) {
	return d.db.Beginx()
}

func (d *DB) Commit(tx *sqlx.Tx) {
	if tx == nil {
		return
	}
	err := tx.Commit()
	if err != nil {
		log.Error(err)
	}
}

func (d *DB) Rollback(tx *sqlx.Tx) {
	if tx == nil {
		return
	}
	err := tx.Rollback()
	if err != nil {
		log.Error(err)
	}
}

type BaetylCloudDB struct {
	DB
	log *log.Logger
}

func NewExtDB() (plugin.Plugin, error) {
	db, err := New()
	if err != nil {
		return nil, err
	}
	dbExt := &BaetylCloudDB{
		log: log.L().With(log.Any("plugin", "databaseext")),
	}
	dbExt.DB = *db.(*DB)
	return dbExt, nil
}

func (d *BaetylCloudDB) InterfaceToTx(tx interface{}) (*sqlx.Tx, error) {
	if tx == nil {
		return nil, nil
	}
	transaction, ok := tx.(*sqlx.Tx)
	if !ok {
		return nil, common.Error(common.ErrConvertConflict, common.Field("error", ErrDatabaseTx))
	}
	return transaction, nil
}
