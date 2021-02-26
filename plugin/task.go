package plugin

import (
	"io"
	"sync"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

//go:generate mockgen -destination=../mock/plugin/task.go -package=plugin github.com/baetyl/baetyl-cloud/v2/plugin Task

// Task interface of Task
type Task interface {
	CreateTask(task *models.Task) (bool, error)
	GetTask(name string) (*models.Task, error)
	AcquireTaskLock(task *models.Task) (bool, error)
	GetNeedProcessTask(number int32) ([]*models.Task, error)
	UpdateTask(task *models.Task) (bool, error)
	DeleteTask(taskName string) (bool, error)

	io.Closer
}

type ProcessFunc func(task *models.Task) error

var TaskRegister taskRegister = taskRegister{
	tasks: &sync.Map{},
}

type taskRegister struct {
	tasks *sync.Map
	sync.RWMutex
}

func (m *taskRegister) AddTask(taskName, processorName string, processor ProcessFunc) bool {
	ps := map[string]ProcessFunc{processorName: processor}
	processors, exist := m.tasks.LoadOrStore(taskName, ps)
	if !exist {
		return true
	}

	m.Lock()
	defer m.Unlock()

	ps, ok := processors.(map[string]ProcessFunc)
	if !ok {
		ps = map[string]ProcessFunc{}
		m.tasks.Store(taskName, ps)
	}

	ps[processorName] = processor
	return true
}

func (m *taskRegister) DeleteTaskProcessor(taskName, processorName string) bool {
	processors := m.GetTasksByName(taskName)

	m.Lock()
	defer m.Unlock()

	delete(processors, processorName)
	return true
}

func (m *taskRegister) GetTasksByName(taskName string) map[string]ProcessFunc {
	if ps, exist := m.tasks.Load(taskName); exist {
		if processes, ok := ps.(map[string]ProcessFunc); ok {
			return processes
		}
	}

	return map[string]ProcessFunc{}
}

func (m *taskRegister) DeleteTasksByName(taskName string) {
	m.tasks.Delete(taskName)
}
