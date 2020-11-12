package plugin

import (
	"database/sql"
	"io"

	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/jmoiron/sqlx"
)

//go:generate mockgen -destination=../mock/plugin/batch.go -package=plugin github.com/baetyl/baetyl-cloud/v2/plugin Batch

// Batch interface of Batch
type Batch interface {
	GetBatch(name, ns string) (*models.Batch, error)
	ListBatch(ns string, filter *models.Filter) ([]models.Batch, error)
	CreateBatch(batch *models.Batch) (sql.Result, error)
	UpdateBatch(batch *models.Batch) (sql.Result, error)
	DeleteBatch(name, ns string) (sql.Result, error)
	CountBatch(ns, name string) (int, error)
	CountBatchByCallback(callbackName, ns string) (int, error)
	GetBatchTx(tx *sqlx.Tx, name, ns string) (*models.Batch, error)
	ListBatchTx(tx *sqlx.Tx, ns string, filter *models.Filter) ([]models.Batch, error)
	CreateBatchTx(tx *sqlx.Tx, batch *models.Batch) (sql.Result, error)
	UpdateBatchTx(tx *sqlx.Tx, batch *models.Batch) (sql.Result, error)
	DeleteBatchTx(tx *sqlx.Tx, name, ns string) (sql.Result, error)
	CountBatchTx(tx *sqlx.Tx, ns, name string) (int, error)
	CountBatchByCallbackTx(tx *sqlx.Tx, callbackName, ns string) (int, error)
	io.Closer
}
