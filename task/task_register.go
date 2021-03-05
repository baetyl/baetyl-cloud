package task

import (
	"sync"

	"github.com/baetyl/baetyl-go/v2/errors"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

var (
	ErrTaskNotExist    = errors.New("failed to get task, task not exist.")
	ErrProcessConflict = errors.New("failed to add processor, the processor already exist.")
	ErrProcessNotExist = errors.New("failed to unregister processor, the processor not exist.")
)

type ProcessorFunc func(task *models.Task) error

var TaskRegister = taskRegister{
	tasks: &sync.Map{},
}

type processor struct {
	name     string
	function ProcessorFunc
}

type ProcessorList []processor

type taskRegister struct {
	tasks *sync.Map
	sync.RWMutex
}

func (m *taskRegister) Register(taskName, processorName string, function ProcessorFunc) error {
	mapValue, exist := m.tasks.LoadOrStore(taskName, ProcessorList{processor{processorName, function}})
	if !exist {
		return nil
	}
	m.Lock()
	defer m.Unlock()
	psList := mapValue.(ProcessorList)
	for _, ps := range psList {
		if ps.name == processorName {
			return ErrProcessConflict
		}
	}

	psList = append(psList, processor{processorName, function})
	m.tasks.Store(taskName, psList)
	return nil
}

func (m *taskRegister) Unregister(taskName, processorName string) error {
	psList := m.GetProcessorListByTask(taskName)
	if psList == nil {
		return ErrProcessNotExist
	}

	m.Lock()
	defer m.Unlock()

	for i := 0; i < len(psList); {
		if psList[i].name == processorName {
			psList = append(psList[:i], psList[i+1:]...)
			m.tasks.Store(taskName, psList)
			return nil
		}
		i++
	}

	return ErrTaskNotExist
}

func (m *taskRegister) GetProcessorListByTask(taskName string) ProcessorList {
	if ps, exist := m.tasks.Load(taskName); exist {
		if processes, ok := ps.(ProcessorList); ok {
			return processes
		}
	}

	return nil
}

func (m *taskRegister) UnregisterTask(taskName string) {
	m.tasks.Delete(taskName)
}
