package entities

import (
	"encoding/json"
	"time"

	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/jinzhu/copier"
)

type Task struct {
	Id               int64     `json:"id,omitempty" db:"id"`
	Name             string    `json:"name,omitempty" db:"name"`
	RegistrationName string    `json:"registrationName" db:"registration_name"`
	Namespace        string    `json:"namespace,omitempty" db:"namespace"`
	ResourceName     string    `json:"resourceName,omitempty" db:"resource_name"`
	ResourceType     string    `json:"resourceType,omitempty" db:"resource_type"`
	Version          int64     `json:"version,omitempty" db:"version"`
	ExpireTime       int64     `json:"expireTime,omitempty" db:"expire_time"`
	Status           int       `json:"status,omitempty" db:"status"`
	Content          string    `json:"content,omitempty" db:"content"`
	CreateTime       time.Time `json:"createTime" db:"create_time"`
	UpdateTime       time.Time `json:"updateTime" db:"update_time"`
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

	processorsStatus := map[string]models.TaskStatus{}
	if task.Content != "" {
		err := json.Unmarshal([]byte(task.Content), &processorsStatus)

		if err != nil {
			return nil, err
		}
	}
	t.ProcessorsStatus = processorsStatus

	return t, nil
}
