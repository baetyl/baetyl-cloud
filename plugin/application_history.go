package plugin

import (
	"database/sql"
	"io"

	"github.com/jmoiron/sqlx"

	v1 "github.com/baetyl/baetyl-go/v2/spec/v1"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

//go:generate mockgen -destination=../mock/plugin/application_history.go -package=plugin github.com/baetyl/baetyl-cloud/v2/plugin AppHistory

type AppHistory interface {
	CreateApplicationHis(app *v1.Application) (sql.Result, error)
	UpdateApplicationHis(app *v1.Application, oldVersion string) (sql.Result, error)
	DeleteApplicationHis(name, namespace, version string) (sql.Result, error)
	GetApplicationHis(name, namespace, version string) (*v1.Application, error)
	ListApplicationHis(namespace string, filter *models.Filter) ([]v1.Application, error)
	CreateApplicationHisWithTx(tx *sqlx.Tx, app *v1.Application) (sql.Result, error)
	UpdateApplicationHisWithTx(tx *sqlx.Tx, app *v1.Application, oldVersion string) (sql.Result, error)
	DeleteApplicationHisWithTx(tx *sqlx.Tx, name, namespace, version string) (sql.Result, error)
	CountApplicationHis(tx *sqlx.Tx, name, namespace string) (int, error)
	io.Closer
}
