package service

import (
	"testing"

	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/stretchr/testify/assert"
)

func TestTaskService(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	taskService, err := NewTaskService(mockObject.conf)
	assert.NoError(t, err)
	task := gentTask()

	mockObject.task.EXPECT().CreateTask(task).Return(true, nil)

	err = taskService.AddTask(task)
	assert.NoError(t, err)

	mockObject.conf.Task.LockExpiredTime = 1
	mockObject.conf.Task.BatchNum = 10
	mockObject.task.EXPECT().GetNeedProcessTask(mockObject.conf.Task.BatchNum, mockObject.conf.Task.LockExpiredTime).Return([]*models.Task{
		task,
	}, nil)
	tasks, err := taskService.GetNeedProcessTasks()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(tasks))
	mockObject.task.EXPECT().UpdateTask(task).Return(true, nil)
	res, err := taskService.UpdateTask(task)
	assert.NoError(t, err)
	assert.True(t, res)

}

func gentTask() *models.Task {
	return &models.Task{
		Name:             "task01",
		Namespace:        "default",
		RegistrationName: "test_task",
		ResourceType:     "namespace",
		ResourceName:     "namespace01",
	}
}
