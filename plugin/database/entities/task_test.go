package entities

import (
	"encoding/json"
	"testing"

	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/stretchr/testify/assert"
)

func TestFromTaskModel(t *testing.T) {
	mTask := &models.Task{
		Id:               1,
		Name:             "task_test",
		Namespace:        "default",
		RegistrationName: "test_delete",
		ProcessorsStatus: map[string]models.TaskStatus{
			"app_delete": models.TaskFinished,
		},
	}

	task, err := FromTaskModel(mTask)
	assert.NoError(t, err)
	processorStatus, _ := json.Marshal(&mTask.ProcessorsStatus)

	assert.Equal(t, string(processorStatus), task.Content)
	assert.Equal(t, mTask.Id, task.Id)
	assert.Equal(t, mTask.Name, task.Name)
	assert.Equal(t, mTask.Namespace, task.Namespace)
	assert.Equal(t, mTask.RegistrationName, task.RegistrationName)
}

func TestToTaskModel(t *testing.T) {
	task := &Task{
		Id:               1,
		Name:             "task_test",
		Namespace:        "default",
		RegistrationName: "test_delete",
		Content:          "{\"app_delete\":3}",
	}

	mTask, err := ToTaskModel(task)
	assert.NoError(t, err)

	processorStatus := map[string]models.TaskStatus{
		"app_delete": models.TaskFinished,
	}

	assert.Equal(t, processorStatus, mTask.ProcessorsStatus)
	assert.Equal(t, mTask.Id, task.Id)
	assert.Equal(t, mTask.Name, task.Name)
	assert.Equal(t, mTask.Namespace, task.Namespace)
	assert.Equal(t, mTask.RegistrationName, task.RegistrationName)
}
