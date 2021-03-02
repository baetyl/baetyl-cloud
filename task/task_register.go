package task

import (
	"sync"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

type ProcessFunc func(task *models.Task) error

var TaskRegister = taskRegister {
	tasks: &sync.Map{},
}

type taskRegister struct {
	tasks *sync.Map
	sync.RWMutex
}

func (m *taskRegister) Register(taskName, processorName string, processor ProcessFunc) bool {
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

func (m *taskRegister) Unregister(taskName, processorName string) bool {
	processors := m.GetProcessorsByTask(taskName)

	m.Lock()
	defer m.Unlock()

	delete(processors, processorName)
	return true
}

func (m *taskRegister) GetProcessorsByTask(taskName string) map[string]ProcessFunc {
	if ps, exist := m.tasks.Load(taskName); exist {
		if processes, ok := ps.(map[string]ProcessFunc); ok {
			return processes
		}
	}

	return map[string]ProcessFunc{}
}

func (m *taskRegister) UnregisterProcessors(taskName string) {
	m.tasks.Delete(taskName)
}