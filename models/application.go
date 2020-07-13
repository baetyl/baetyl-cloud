package models

import (
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"time"
)

type ApplicationView struct {
	specV1.Application `json:",inline"`
	Registries         []RegistryView `json:"registries,omitempty"`
}

type AppItem struct {
	Name              string            `json:"name,omitempty" validate:"omitempty,resourceName"`
	Type              string            `json:"type,omitempty" default:"container"`
	Labels            map[string]string `json:"labels,omitempty"`
	Selector          string            `json:"selector"`
	Version           string            `json:"version,omitempty"`
	Namespace         string            `json:"namespace,omitempty"`
	CreationTimestamp time.Time         `json:"createTime,omitempty"`
	Description       string            `json:"description,omitempty"`
	System            bool              `json:"system,omitempty"`
}

// ApplicationList app List
type ApplicationList struct {
	Total       int          `json:"total"`
	ListOptions *ListOptions `json:"listOptions"`
	Items       []AppItem    `json:"items"`
}

type ServiceFunction struct {
	Functions []specV1.ServiceFunction `json:"functions,omitempty"`
}
