package service

import (
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
	"sync"
	"time"
)

type TaskService struct {
	task plugin.Task
}

func (t *TaskService) AddTask(task *models.Task) error {
	_, err := t.task.CreateTask(task)
	return err
}

func (t *TaskService) GetNeedProcessTasks() ([]*models.Task, error) {
	return t.task.GetNeedProcessTask(100, 300)
}

func (t *TaskService) UpdateTask(task *models.Task) (bool, error) {
	return t.task.UpdateTask(task)
}

type TaskProcessor func(task *models.Task) error

type taskRegister struct {
	tasks *sync.Map
	sync.RWMutex
}

var TaskRegister taskRegister = taskRegister{}

func (m *taskRegister) AddTask(taskName, processorName string, processor TaskProcessor) bool {
	ps := map[string]TaskProcessor{processorName: processor}
	processors, exist := m.tasks.LoadOrStore(taskName, ps)
	if !exist {
		return true
	}

	m.Lock()
	defer m.Unlock()

	ps, ok := processors.(map[string]TaskProcessor)
	if !ok {
		ps = map[string]TaskProcessor{}
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

func (m *taskRegister) GetTasksByName(taskName string) map[string]TaskProcessor {
	if ps, exist := m.tasks.Load(taskName); exist {
		if processes, ok := ps.(map[string]TaskProcessor); ok {
			return processes
		}
	}

	return map[string]TaskProcessor{}
}

func (m *taskRegister) DeleteTasksByName(taskName string) {
	m.tasks.Delete(taskName)
}

type TaskManager struct {
	taskService TaskService
	tasks       chan *models.Task
	lock        plugin.Locker
	concurrency chan int
	closeSignal chan int
}

func NewTaskManager(service TaskService, locker plugin.Locker, taskQueueNum, concuurencyNum int) *TaskManager {
	return &TaskManager{
		taskService: service,
		lock:        locker,
		tasks:       make(chan *models.Task, taskQueueNum),
		concurrency: make(chan int, concuurencyNum),
		closeSignal: make(chan int, 1),
	}
}

func (m *TaskManager) Start() {
	go m.FetchTask()
	go m.RunTasks()
}

func (m *TaskManager) Close() {
	m.closeSignal <- 1
}

func (m *TaskManager) FetchTask() {
	ticker := time.NewTicker(time.Second)

	for {
		select {
		case <-ticker.C:
			tasks, err := m.taskService.GetNeedProcessTasks()
			if err != nil {
			}

			for _, t := range tasks {
				m.tasks <- t
			}
		case <-m.closeSignal:
			break
		}
	}
	close(m.tasks)
}

func (m *TaskManager) RunTasks() {
	for task := range m.tasks {
		m.concurrency <- 1
		go m.runTask(task)
	}
}

func (m *TaskManager) runTask(task *models.Task) {
	defer func() { <-m.concurrency }()
	lock, err := m.lock.LockWithExpireTime(task.TaskName, 100)
	if err != nil || !lock {
		return
	}

	processors := TaskRegister.GetTasksByName(task.TaskName)

	if task.ProcessorsStatus == nil {
		task.ProcessorsStatus = map[string]plugin.TaskStatus{}
	}

	task.Status = int(plugin.TaskFinished)

	for pName, processFunc := range processors {
		if needRunTask(pName, task.ProcessorsStatus) {
			err := processFunc(task)
			if err != nil {
				task.ProcessorsStatus[pName] = plugin.TaskNeedRetry

				// set to need retry
				task.Status = int(plugin.TaskNeedRetry)
			} else {
				task.ProcessorsStatus[pName] = plugin.TaskFinished
			}
		}
	}

	_, err = m.taskService.UpdateTask(task)
	if err != nil {
	}
}

func needRunTask(processorName string, processorsStatus map[string]plugin.TaskStatus) bool {
	pStatus, ok := processorsStatus[processorName]
	if !ok || pStatus < plugin.TaskFinished {
		return true
	}
	return false
}
