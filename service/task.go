package service

import (
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

//go:generate mockgen -destination=../mock/service/task.go -package=service github.com/baetyl/baetyl-cloud/v2/service TaskService

type TaskService interface {
	AddTask(task *models.Task) error
	GetNeedProcessTasks() ([]*models.Task, error)
	UpdateTask(task *models.Task) (bool, error)
}

// NewTaskService NewTaskService
func NewTaskService(config *config.CloudConfig) (TaskService, error) {
	task, err := plugin.GetPlugin(config.Plugin.Task)
	if err != nil {
		return nil, err
	}

	return &taskService{
		task: task.(plugin.Task),
		cfg:  config,
	}, nil
}

type taskService struct {
	task plugin.Task
	cfg  *config.CloudConfig
}

func (t *taskService) AddTask(task *models.Task) error {
	_, err := t.task.CreateTask(task)
	return err
}

func (t *taskService) GetNeedProcessTasks() ([]*models.Task, error) {
	return t.task.GetNeedProcessTask(t.cfg.Task.BatchNum)
}

func (t *taskService) UpdateTask(task *models.Task) (bool, error) {
	return t.task.UpdateTask(task)
}
