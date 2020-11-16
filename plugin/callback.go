package plugin

import (
	"database/sql"
	"io"

	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/jmoiron/sqlx"
)

//go:generate mockgen -destination=../mock/plugin/callback.go -package=plugin github.com/baetyl/baetyl-cloud/v2/plugin Callback

// Callback interface of Callback
type Callback interface {
	GetCallback(name, namespace string) (*models.Callback, error)
	CreateCallback(callback *models.Callback) (sql.Result, error)
	UpdateCallback(callback *models.Callback) (sql.Result, error)
	DeleteCallback(name, ns string) (sql.Result, error)
	GetCallbackTx(tx *sqlx.Tx, name, namespace string) (*models.Callback, error)
	CreateCallbackTx(tx *sqlx.Tx, callback *models.Callback) (sql.Result, error)
	UpdateCallbackTx(tx *sqlx.Tx, callback *models.Callback) (sql.Result, error)
	DeleteCallbackTx(tx *sqlx.Tx, name, ns string) (sql.Result, error)
	io.Closer
}
