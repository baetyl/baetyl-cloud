package plugin

import (
	"database/sql"
	"io"

	"github.com/jmoiron/sqlx"

	"github.com/baetyl/baetyl-cloud/v2/models"
	v1 "github.com/baetyl/baetyl-go/v2/spec/v1"
)

//go:generate mockgen -destination=../mock/plugin/application_history.go -package=plugin github.com/baetyl/baetyl-cloud/v2/plugin AppHistory

type AppHistory interface {
	CreateApplication(app *v1.Application) (sql.Result, error)
	UpdateApplication(app *v1.Application, oldVersion string) (sql.Result, error)
	DeleteApplication(name, namespace, version string) (sql.Result, error)
	GetApplication(name, namespace, version string) (*v1.Application, error)
	ListApplication(namespace string, filter *models.Filter) ([]v1.Application, error)
	CreateApplicationWithTx(tx *sqlx.Tx, app *v1.Application) (sql.Result, error)
	UpdateApplicationWithTx(tx *sqlx.Tx, app *v1.Application, oldVersion string) (sql.Result, error)
	DeleteApplicationWithTx(tx *sqlx.Tx, name, namespace, version string) (sql.Result, error)
	CountApplication(tx *sqlx.Tx, name, namespace string) (int, error)
	io.Closer
}
