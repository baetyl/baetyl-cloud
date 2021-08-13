package models

import (
	"time"

	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
)

type ApplicationView struct {
	Name              string                `json:"name,omitempty" validate:"resourceName"`
	Mode              string                `json:"mode,omitempty" default:"kube"`
	Type              string                `json:"type,omitempty" default:"container"`
	Labels            map[string]string     `json:"labels,omitempty"`
	Namespace         string                `json:"namespace,omitempty"`
	CreationTimestamp time.Time             `json:"createTime,omitempty"`
	Version           string                `json:"version,omitempty"`
	Selector          string                `json:"selector,omitempty"`
	NodeSelector      string                `json:"nodeSelector,omitempty"`
	Services          []ServiceView         `json:"services,omitempty" validate:"dive"`
	Volumes           []VolumeView          `json:"volumes,omitempty" validate:"dive"`
	Description       string                `json:"description,omitempty"`
	System            bool                  `json:"system,omitempty"`
	Registries        []RegistryView        `json:"registries,omitempty"`
	CronStatus        specV1.CronStatusCode `json:"cronStatus" default:"0"`
	CronTime          time.Time             `json:"cronTime,omitempty"`
}

// VolumeView volume view
type VolumeView struct {
	// specified name of the volume
	Name        string                       `json:"name,omitempty" binding:"required" validate:"omitempty,resourceName"`
	HostPath    *specV1.HostPathVolumeSource `json:"hostPath,omitempty"`
	Config      *specV1.ObjectReference      `json:"config,omitempty"`
	Secret      *specV1.ObjectReference      `json:"secret,omitempty"`
	Certificate *specV1.ObjectReference      `json:"certificate,omitempty"`
}

type AppItem struct {
	Name              string                `json:"name,omitempty" validate:"omitempty,resourceName"`
	Mode              string                `json:"mode,omitempty" default:"kube"`
	Type              string                `json:"type,omitempty" default:"container"`
	Labels            map[string]string     `json:"labels,omitempty"`
	Selector          string                `json:"selector"`
	NodeSelector      string                `json:"nodeSelector"`
	Version           string                `json:"version,omitempty"`
	Namespace         string                `json:"namespace,omitempty"`
	CreationTimestamp time.Time             `json:"createTime,omitempty"`
	Description       string                `json:"description,omitempty"`
	System            bool                  `json:"system,omitempty"`
	CronStatus        specV1.CronStatusCode `json:"cronStatus,omitempty" default:"0"`
	CronTime          time.Time             `json:"cronTime,omitempty"`
}

// ApplicationList app List
type ApplicationList struct {
	Total        int `json:"total"`
	*ListOptions `json:",inline"`
	Items        []AppItem `json:"items"`
}

type ServiceFunction struct {
	Functions []specV1.ServiceFunction `json:"functions,omitempty"`
}

type ServiceView struct {
	specV1.Service `json:",inline"`
	ProgramConfig  string `json:"programConfig,omitempty"`
}
