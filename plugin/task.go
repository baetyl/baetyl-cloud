package plugin

import (
	"io"

	"github.com/baetyl/baetyl-go/v2/task"
)

//go:generate mockgen -destination=../mock/plugin/task.go -package=plugin github.com/baetyl/baetyl-cloud/v2/plugin Task

// Task interface of Task
type Task interface {
	task.TaskProducer
	task.TaskWorker

	io.Closer
}
