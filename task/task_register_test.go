package task

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

func TestTaskRegister(t *testing.T) {
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
		err := TaskRegister.Register(task.TaskName, task.ProcessorName, task.processor)
		assert.Nil(t, err)
	}

	err := TaskRegister.Register(tasks[0].TaskName, tasks[0].ProcessorName, tasks[0].processor)
	assert.Equal(t, err, ErrProcessConflict)

	psList := TaskRegister.GetProcessorListByTask("task01")
	assert.Equal(t, len(tasks), len(psList))

	psList = TaskRegister.GetProcessorListByTask("task03")
	assert.Nil(t, psList)

	err = TaskRegister.Unregister("task01", "delete_processor")
	assert.Nil(t, err)
	psList = TaskRegister.GetProcessorListByTask("task01")
	assert.Equal(t, 1, len(psList))
	assert.Equal(t, psList[0].name, "add_processor")

	TaskRegister.UnregisterTask("task01")
	psList = TaskRegister.GetProcessorListByTask("task01")
	assert.Equal(t, 0, len(psList))

}

type taskData struct {
	TaskName      string
	ProcessorName string
	processor     ProcessorFunc
	Expected      bool
}

func mockDeleteProcessor(task *models.Task) error {
	return nil
}

func mockAddProcessor(task *models.Task) error {
	return nil
}