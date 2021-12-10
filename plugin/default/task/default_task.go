package task

import (
	"github.com/baetyl/baetyl-go/v2/task"

	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

func init() {
	plugin.RegisterFactory("defaulttask", New)
}

type defaultTask struct {
	task.TaskProducer
	task.TaskWorker
	broker  task.TaskBroker
	backend task.TaskBackend
}

func New() (plugin.Plugin, error) {
	broker := task.NewChannelBroker(10)
	backend := task.NewMapBackend()

	return &defaultTask{
		task.NewTaskProducer(broker, backend),
		task.NewTaskWorker(broker, backend),
		broker,
		backend,
	}, nil
}

func (dt *defaultTask) Close() error {
	return dt.broker.Close()
}
