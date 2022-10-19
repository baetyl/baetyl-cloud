package models

import (
	"time"

	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
)

type ApplicationView struct {
	Name              string                `json:"name,omitempty" binding:"res_name"`
	Mode              string                `json:"mode,omitempty" default:"kube"`
	Type              string                `json:"type,omitempty" default:"container"`
	Labels            map[string]string     `json:"labels,omitempty"`
	Namespace         string                `json:"namespace,omitempty"`
	CreationTimestamp time.Time             `json:"createTime,omitempty"`
	Version           string                `json:"version,omitempty"`
	Selector          string                `json:"selector,omitempty"`
	NodeSelector      string                `json:"nodeSelector,omitempty"`
	InitServices      []ServiceView         `json:"initServices,omitempty" binding:"dive"`
	Services          []ServiceView         `json:"services,omitempty" binding:"dive"`
	Volumes           []VolumeView          `json:"volumes,omitempty" binding:"dive"`
	Description       string                `json:"description,omitempty"`
	System            bool                  `json:"system,omitempty"`
	Registries        []RegistryView        `json:"registries,omitempty"`
	CronStatus        specV1.CronStatusCode `json:"cronStatus" default:"0"`
	CronTime          time.Time             `json:"cronTime,omitempty"`
	HostNetwork       bool                  `json:"hostNetwork,omitempty"` // specifies host network mode of service
	Replica           int                   `json:"replica"`
	Workload          string                `json:"workload,omitempty"` // deployment | daemonset | statefulset | job
	JobConfig         *specV1.AppJobConfig  `json:"jobConfig,omitempty"`
	Ota               specV1.OtaInfo        `json:"ota,omitempty"`
	AutoScaleCfg      *specV1.AutoScaleCfg  `json:"autoScaleCfg,omitempty"`
}

// VolumeView volume view
type VolumeView struct {
	// specified name of the volume
	Name        string                       `json:"name,omitempty" binding:"required,res_name"`
	HostPath    *specV1.HostPathVolumeSource `json:"hostPath,omitempty"`
	Config      *specV1.ObjectReference      `json:"config,omitempty"`
	Secret      *specV1.ObjectReference      `json:"secret,omitempty"`
	EmptyDir    *specV1.EmptyDirVolumeSource `json:"emptyDir,omitempty"`
	Certificate *specV1.ObjectReference      `json:"certificate,omitempty"`
}

type AppItem struct {
	Name              string                `json:"name,omitempty" binding:"omitempty,res_name"`
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
	HostNetwork       bool                  `json:"hostNetwork,omitempty" yaml:"hostNetwork,omitempty"` // specifies host network mode of service
	Replica           int                   `json:"replica,omitempty" yaml:"replica,omitempty" binding:"required" default:"1"`
	Workload          string                `json:"workload,omitempty" yaml:"workload,omitempty"` // deployment | daemonset | statefulset | job
	JobConfig         *specV1.AppJobConfig  `json:"jobConfig,omitempty" yaml:"jobConfig,omitempty"`
	Ota               specV1.OtaInfo        `json:"ota,omitempty"`
	AutoScaleCfg      *specV1.AutoScaleCfg  `json:"autoScaleCfg,omitempty"`
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
