package task

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

func TestTaskRegister_AddTask(t *testing.T) {
	tasks := []taskData{
		{
			TaskName:      "task01",
			ProcessorName: "add_processor",
			processor:     mockAddProcessor,
		},
		{
			TaskName:      "task01",
			ProcessorName: "delete_processor",
			processor:     mockDeleteProcessor,
		},
	}

	for _, task := range tasks {
		res := TaskRegister.Register(task.TaskName, task.ProcessorName, task.processor)
		assert.True(t, res)
	}

	processors := TaskRegister.GetProcessorsByTask("task01")
	assert.Equal(t, len(tasks), len(processors))

	res := TaskRegister.Unregister("task01", "delete_processor")
	assert.True(t, res)
	processors = TaskRegister.GetProcessorsByTask("task01")
	assert.Equal(t, 1, len(processors))
	assert.Nil(t, processors["delete_processor"])

	TaskRegister.UnregisterProcessors("task01")
	processors = TaskRegister.GetProcessorsByTask("task01")
	assert.Equal(t, 0, len(processors))

}

type taskData struct {
	TaskName      string
	ProcessorName string
	processor     ProcessFunc
	Expected      bool
}

func mockDeleteProcessor(task *models.Task) error {
	return nil
}

func mockAddProcessor(task *models.Task) error {
	return nil
}