package plugin

import (
	"database/sql"
	"io"

	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/jmoiron/sqlx"
)

//go:generate mockgen -destination=../mock/plugin/task.go -package=plugin github.com/baetyl/baetyl-cloud/v2/plugin Task

// Task interface of Task
type Task interface {
	CreateTask(task *models.Task) (sql.Result, error)
	UpdateTask(task *models.Task) (sql.Result, error)
	GetTask(traceId string) (*models.Task, error)
	DeleteTask(traceId string) (sql.Result, error)
	CountTask(task *models.Task) (int, error)

	GetTaskTx(tx *sqlx.Tx, traceId string) (*models.Task, error)
	CreateTaskTx(tx *sqlx.Tx, task *models.Task) (sql.Result, error)
	UpdateTaskTx(tx *sqlx.Tx, task *models.Task) (sql.Result, error)
	DeleteTaskTx(tx *sqlx.Tx, traceId string) (sql.Result, error)
	io.Closer
}
