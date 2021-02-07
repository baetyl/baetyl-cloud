package models

import (
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

type Task struct {
	Id               int64                        `json:"id,omitempty"`
	TaskName         string                       `json:"taskName"`
	Namespace        string                       `json:"namespace,omitempty"`
	ResourceName     string                       `json:"resourceName,omitempty"`
	ResourceType     string                       `json:"resourceType,omitempty"`
	Version          int64                        `json:"version,omitempty"`
	ExpireTime       int64                        `json:"expireTime,omitempty"`
	Status           int                          `json:"status,omitempty"`
	ProcessorsStatus map[string]plugin.TaskStatus `json:"processorsStatus,omitempty"`
}
