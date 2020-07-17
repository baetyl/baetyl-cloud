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
	CreateSystemConfig(sysConfig *models.SystemConfig) (sql.Result, error)
	DeleteSystemConfig(key string) (sql.Result, error)
	GetSystemConfig(key string) (*models.SystemConfig, error)
	ListSystemConfig(key string, page, size int) ([]models.SystemConfig, error)
	CountSystemConfig(key string) (int, error)
	UpdateSystemConfig(sysConfig *models.SystemConfig) (sql.Result, error)

	CreateSystemConfigTx(tx *sqlx.Tx, sysConfig *models.SystemConfig) (sql.Result, error)
	DeleteSystemConfigTx(tx *sqlx.Tx, key string) (sql.Result, error)
	GetSystemConfigTx(tx *sqlx.Tx, key string) (*models.SystemConfig, error)
	ListSystemConfigTx(tx *sqlx.Tx, key string, page, size int) ([]models.SystemConfig, error)
	CountSystemConfigTx(tx *sqlx.Tx, key string) (int, error)
	UpdateSystemConfigTx(tx *sqlx.Tx, sysConfig *models.SystemConfig) (sql.Result, error)

	io.Closer
}