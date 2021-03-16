package lock

import (
	"fmt"
	mockPlugin "github.com/baetyl/baetyl-cloud/v2/mock/plugin"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTaskLocker_Lock(t *testing.T) {
	mockTask := initMockTask(t)
	locker := taskLocker{
		task: mockTask,
	}

	task := &models.Task{
		Name:             "task01",
		Namespace:        "default",
		RegistrationName: "delete",
		ResourceName:     "namespace01",
		ResourceType:     "namespace",
	}

	mockTask.EXPECT().AcquireTaskLock(task).Return(true, nil)
	mockTask.EXPECT().GetTask(task.Name).Return(task, nil)

	err := locker.Lock(task.Name)
	assert.NoError(t, err)
}

func TestTaskLocker_LockWithExpireTime(t *testing.T) {
	mockTask := initMockTask(t)
	locker := taskLocker{
		task: mockTask,
	}

	task := &models.Task{
		Name:             "task01",
		Namespace:        "default",
		RegistrationName: "delete",
		ResourceName:     "namespace01",
		ResourceType:     "namespace",
	}

	mockTask.EXPECT().GetTask(task.Name).Return(task, nil)
	mockTask.EXPECT().AcquireTaskLock(task).Return(false, fmt.Errorf("lock error"))
	err := locker.LockWithExpireTime(task.Name, 10)
	assert.Error(t, err)
	assert.Equal(t, "lock error", err.Error())

	mockTask.EXPECT().GetTask(task.Name).Return(task, fmt.Errorf("get task error"))
	err = locker.LockWithExpireTime(task.Name, 10)
	assert.Error(t, err)
	assert.Equal(t, "get task error", err.Error())

	mockTask.EXPECT().GetTask(task.Name).Return(nil, nil)
	err = locker.LockWithExpireTime(task.Name, 10)
	assert.NoError(t, err)
}

func TestTaskLocker_ReleaseLock(t *testing.T) {
	mockTask := initMockTask(t)
	locker := taskLocker{
		task: mockTask,
	}

	task := &models.Task{
		Name:             "task01",
		Namespace:        "default",
		RegistrationName: "delete",
		ResourceName:     "namespace01",
		ResourceType:     "namespace",
	}

	mockTask.EXPECT().GetTask(task.Name).Return(task, nil)
	mockTask.EXPECT().AcquireTaskLock(task).Return(false, fmt.Errorf("failed"))
	err := locker.Unlock(task.Name)
	assert.Error(t, err)
	assert.Equal(t, "failed", err.Error())
}

func TestTaskLocker_Close(t *testing.T) {
	mockTask := initMockTask(t)
	locker := taskLocker{
		task: mockTask,
	}
	mockTask.EXPECT().Close().Return(nil)

	err := locker.Close()
	assert.NoError(t, err)
}

func initMockTask(t *testing.T) *mockPlugin.MockTask {
	ctl := gomock.NewController(t)
	task := mockPlugin.NewMockTask(ctl)

	return task
}
