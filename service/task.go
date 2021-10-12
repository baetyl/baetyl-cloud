package service

import (
	"github.com/baetyl/baetyl-go/v2/task"

	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

//go:generate mockgen -destination=../mock/service/task.go -package=service github.com/baetyl/baetyl-cloud/v2/service TaskService

type TaskService interface {
	task.TaskProducer
	task.TaskWorker
}

// NewTaskService NewTaskService
func NewTaskService(config *config.CloudConfig) (TaskService, error) {
	taskService, err := plugin.GetPlugin(config.Plugin.Task)
	if err != nil {
		return nil, err
	}

	return taskService.(plugin.Task), nil
}
