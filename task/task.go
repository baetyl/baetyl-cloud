package task

import (
	"time"

	"github.com/baetyl/baetyl-go/v2/log"

	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
	"github.com/baetyl/baetyl-cloud/v2/service"
)

const (
	WarningNum = 10
)

type TaskManager struct {
	taskService  service.TaskService
	tasks        chan *models.Task
	lock         plugin.Locker
	concurrency  chan int
	closeSignal  chan int
	scheduleTime time.Duration
	config       *config.CloudConfig
}

func NewTaskManager(cfg *config.CloudConfig) (*TaskManager, error) {
	locker, err := plugin.GetPlugin(cfg.Plugin.Locker)
	if err != nil {
		log.L().Error("init task manager failed", log.Error(err))
		return nil, err
	}

	taskService, err := service.NewTaskService(cfg)
	if err != nil {
		log.L().Error("init task service error", log.Error(err))
		return nil, err
	}
	return &TaskManager{
		taskService:  taskService,
		lock:         locker.(plugin.Locker),
		tasks:        make(chan *models.Task, cfg.Task.QueueLength),
		concurrency:  make(chan int, cfg.Task.ConcurrentNum),
		closeSignal:  make(chan int, 1),
		scheduleTime: time.Duration(cfg.Task.ScheduleTime) * time.Second,
		config:       cfg,
	}, nil
}

func (m *TaskManager) Start() {
	log.L().Debug("task start")
	go m.FetchTask()
	go m.RunTasks()
}

func (m *TaskManager) Close() {
	m.closeSignal <- 1
}

func (m *TaskManager) FetchTask() {
	timer := time.NewTimer(m.scheduleTime)

	for {
		select {
		case <-m.closeSignal:
			break
		case <-timer.C:
			tasks, err := m.taskService.GetNeedProcessTasks()

			if err != nil {
				log.L().Error("get run tasks error", log.Error(err))
			}

			for _, t := range tasks {
				m.tasks <- t
			}

			timer.Reset(m.scheduleTime)
		}
	}

	close(m.tasks)
	timer.Stop()
}

func (m *TaskManager) RunTasks() {
	for task := range m.tasks {
		m.concurrency <- 1
		go m.runTask(task)
	}
}

func (m *TaskManager) runTask(task *models.Task) {
	defer func() { <-m.concurrency }()
	lock, err := m.lock.LockWithExpireTime(task.Name, m.config.Lock.ExpireTime)
	if err != nil || !lock {
		log.L().Error("get lock error",
			log.Any("name", task.Name),
			log.Any("namespace", task.Namespace),
			log.Any("resourceType", task.ResourceType),
			log.Any("resourceName", task.ResourceName),
			log.Error(err))
		return
	}

	processors := TaskRegister.GetProcessorsByTask(task.RegistrationName)

	if task.ProcessorsStatus == nil {
		task.ProcessorsStatus = map[string]models.TaskStatus{}
	}

	task.Status = int(models.TaskFinished)

	for pName, processFunc := range processors {
		if isNeedRunTask(pName, task.ProcessorsStatus) {
			err := processFunc(task)
			if err != nil {
				if task.Version > WarningNum {
					log.L().Error("run process error", log.Any("name", task.Name),
						log.Any("registrationName", task.RegistrationName),
						log.Any("processorName", pName), log.Any("namespace", task.Namespace),
						log.Any("resourceType", task.ResourceType),
						log.Any("resourceName", task.ResourceName), log.Error(err))
				} else {
					log.L().Warn("run process error", log.Any("name", task.Name),
						log.Any("registrationName", task.RegistrationName),
						log.Any("processorName", pName), log.Any("namespace", task.Namespace),
						log.Any("resourceType", task.ResourceType),
						log.Any("resourceName", task.ResourceName), log.Error(err))
				}

				task.ProcessorsStatus[pName] = models.TaskNeedRetry

				// set to need retry
				task.Status = int(models.TaskNeedRetry)
			} else {
				task.ProcessorsStatus[pName] = models.TaskFinished
			}
		}
	}

	task.Version = task.Version + 1
	_, err = m.taskService.UpdateTask(task)
	if err != nil {
		log.L().Error("update task error",
			log.Any("name", task.Name),
			log.Any("namespace", task.Namespace),
			log.Any("resourceType", task.ResourceType),
			log.Any("resourceName", task.ResourceName),
			log.Error(err))
	}
}

func isNeedRunTask(processorName string, processorsStatus map[string]models.TaskStatus) bool {
	pStatus, ok := processorsStatus[processorName]
	if !ok || pStatus < models.TaskFinished {
		return true
	}
	return false
}
