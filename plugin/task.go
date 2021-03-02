package plugin

import (
	"io"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

//go:generate mockgen -destination=../mock/plugin/task.go -package=plugin github.com/baetyl/baetyl-cloud/v2/plugin Task

// Task interface of Task
type Task interface {
	CreateTask(task *models.Task) (bool, error)
	GetTask(name string) (*models.Task, error)
	AcquireTaskLock(task *models.Task) (bool, error)
	GetNeedProcessTask(number int32, expiredSeconds int32) ([]*models.Task, error)
	UpdateTask(task *models.Task) (bool, error)
	DeleteTask(taskName string) (bool, error)

	io.Closer
}
