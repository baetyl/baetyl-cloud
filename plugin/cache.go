package plugin

import (
	"database/sql"
	"github.com/baetyl/baetyl-cloud/models"
	"github.com/jmoiron/sqlx"
	"io"
)

//go:generate mockgen -destination=../mock/plugin/cache.go -package=plugin github.com/baetyl/baetyl-cloud/plugin CacheStorage

type CacheStorage interface {
	Transact(func(*sqlx.Tx) error) error
	// new system config
	AddCache(key,value string) (sql.Result, error)
	DeleteCache(key string) (sql.Result, error)
	GetCache(key string) (*models.Cache, error)
	ListCache(key string, page, size int) ([]models.Cache, error)
	CountCache(key string) (int, error)
	ReplaceCache(key ,value string) (sql.Result, error)

	AddCacheTx(tx *sqlx.Tx, key,value string) (sql.Result, error)
	DeleteCacheTx(tx *sqlx.Tx, key string) (sql.Result, error)
	GetCacheTx(tx *sqlx.Tx, key string) (*models.Cache, error)
	ListCacheTx(tx *sqlx.Tx, key string, page, size int) ([]models.Cache, error)
	CountCacheTx(tx *sqlx.Tx, key string) (int, error)
	ReplaceCacheTx(tx *sqlx.Tx, key ,value string) (sql.Result, error)

	io.Closer
}