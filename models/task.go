package models

type TaskStatus int

const (
	TaskNew TaskStatus = iota
	TaskProcessing
	TaskNeedRetry
	TaskFinished
	TaskFailed
)

type Task struct {
	Id               int64                 `json:"id,omitempty"`
	Name             string                `json:"name,omitempty"`
	RegistrationName string                `json:"registrationName,omitempty"`
	Namespace        string                `json:"namespace,omitempty"`
	ResourceName     string                `json:"resourceName,omitempty"`
	ResourceType     string                `json:"resourceType,omitempty"`
	Version          int64                 `json:"version,omitempty"`
	ExpireTime       int64                 `json:"expireTime,omitempty"`
	Status           int                   `json:"status,omitempty"`
	ProcessorsStatus map[string]TaskStatus `json:"processorsStatus,omitempty"`
}
