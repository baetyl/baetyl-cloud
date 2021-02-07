package entities

import (
	"encoding/json"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
	"github.com/jinzhu/copier"
	"time"
)

type Task struct {
	Id           int64     `json:"id,omitempty" db:"id"`
	TaskName     string    `json:"taskName" db:"task_name"`
	Namespace    string    `json:"namespace,omitempty" db:"namespace"`
	ResourceName string    `json:"resourceName,omitempty" db:"resource_name"`
	ResourceType string    `json:"resourceType,omitempty" db:"resource_type"`
	Version      int64     `json:"version,omitempty" db:"version"`
	ExpireTime   int64     `json:"expireTime,omitempty" db:"expire_time"`
	Status       int       `json:"status,omitempty" db:"status"`
	Content      string    `json:"content,omitempty" db:"content"`
	CreatedTime  time.Time `json:"createdTime" db:"created_time"`
	UpdatedTime  time.Time `json:"updatedTime" db:"updated_time"`
}

func FromTaskModel(task *models.Task) (*Task, error) {
	t := new(Task)
	err := copier.Copy(t, task)

	if err != nil {
		return nil, err
	}

	content, err := json.Marshal(task.ProcessorsStatus)
	if err != nil {
		return nil, err
	}
	t.Content = string(content)

	return t, nil
}

func ToTaskModel(task *Task) (*models.Task, error) {
	t := new(models.Task)
	err := copier.Copy(t, task)

	if err != nil {
		return nil, err
	}

	processorsStatus := map[string]plugin.TaskStatus{}
	if task.Content != "" {
		err := json.Unmarshal([]byte(task.Content), processorsStatus)

		if err != nil {
			return nil, err
		}
	}
	t.ProcessorsStatus = processorsStatus

	return t, nil
}
