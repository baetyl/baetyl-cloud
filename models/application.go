package models

import (
	"reflect"
	"time"

	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
)

type ApplicationView struct {
	Name              string            `json:"name,omitempty" validate:"resourceName"`
	Type              string            `json:"type,omitempty" default:"container"`
	Labels            map[string]string `json:"labels,omitempty"`
	Namespace         string            `json:"namespace,omitempty"`
	CreationTimestamp time.Time         `json:"createTime,omitempty"`
	Version           string            `json:"version,omitempty"`
	Selector          string            `json:"selector,omitempty"`
	Services          []specV1.Service  `json:"services,omitempty" validate:"dive"`
	Volumes           []VolumeView      `json:"volumes,omitempty" validate:"dive"`
	Description       string            `json:"description,omitempty"`
	System            bool              `json:"system,omitempty"`
	Registries        []RegistryView    `json:"registries,omitempty"`
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

func EqualApp(app1, app2 *specV1.Application) bool {
	if len(app1.Volumes) != len(app2.Volumes) {
		return false
	}

	if len(app1.Volumes) != 0 {
		for i := range app1.Volumes {
			v1 := app1.Volumes[i]
			v2 := app2.Volumes[i]
			flag := (v1.Name == v2.Name) &&
				((v1.Secret != nil && v2.Secret != nil && v1.Secret.Name == v2.Secret.Name) || (v1.Secret == nil && v2.Secret == nil)) &&
				((v1.Config != nil && v2.Config != nil && v1.Config.Name == v2.Config.Name) || (v1.Config == nil && v2.Config == nil)) &&
				((v1.HostPath != nil && v2.HostPath != nil && v1.HostPath.Path == v2.HostPath.Path) || (v1.HostPath == nil && v2.HostPath == nil))
			if !flag {
				return false
			}
		}
	}

	if len(app1.Services) != len(app2.Services) {
		return false
	}

	if len(app1.Services) != 0 {
		for j := range app1.Services {
			s1 := app1.Services[j]
			s2 := app2.Services[j]
			if flag := (s1.Name == s2.Name) && (s1.Hostname == s2.Hostname) && (s1.Image == s2.Image) && (s1.Replica == s2.Replica) && (s1.HostNetwork == s2.HostNetwork) &&
				(s1.Runtime == s2.Runtime) && reflect.DeepEqual(s1.SecurityContext, s2.SecurityContext) && reflect.DeepEqual(s1.FunctionConfig, s2.FunctionConfig); !flag {
				return false
			}

			if len(s1.VolumeMounts) != len(s2.VolumeMounts) || (len(s1.VolumeMounts) != 0 && !reflect.DeepEqual(s1.VolumeMounts, s2.VolumeMounts)) {
				return false
			}

			if len(s1.Ports) != len(s2.Ports) || (len(s1.Ports) != 0 && !reflect.DeepEqual(s1.Ports, s2.Ports)) {
				return false
			}

			if len(s1.Devices) != len(s2.Devices) || (len(s1.Devices) != 0 && !reflect.DeepEqual(s1.Devices, s2.Devices)) {
				return false
			}

			if len(s1.Args) != len(s2.Args) || (len(s1.Args) != 0 && !reflect.DeepEqual(s1.Args, s2.Args)) {
				return false
			}

			if len(s1.Env) != len(s2.Env) || (len(s1.Env) != 0 && !reflect.DeepEqual(s1.Env, s2.Env)) {
				return false
			}

			if len(s1.Labels) != len(s2.Labels) || (len(s1.Labels) != 0 && !reflect.DeepEqual(s1.Labels, s2.Labels)) {
				return false
			}

			if len(s1.Functions) != len(s2.Functions) || (len(s1.Functions) != 0 && !reflect.DeepEqual(s1.Functions, s2.Functions)) {
				return false
			}
			if s1.Resources != nil && s2.Resources != nil {
				if len(s1.Resources.Limits) != len(s2.Resources.Limits) || (len(s1.Resources.Limits) != 0 && !reflect.DeepEqual(s1.Resources.Limits, s2.Resources.Limits)) {
					return false
				}

				if len(s1.Resources.Requests) != len(s2.Resources.Requests) || (len(s1.Resources.Requests) != 0 && !reflect.DeepEqual(s1.Resources.Requests, s2.Resources.Requests)) {
					return false
				}
			} else if s1.Resources == nil && s2.Resources == nil {

			} else {
				return false
			}
		}
	}

	if len(app1.Labels) != len(app2.Labels) || (len(app1.Labels) != 0 && !reflect.DeepEqual(app1.Labels, app2.Labels)) {
		return false
	}

	return reflect.DeepEqual(app1.Selector, app2.Selector) &&
		reflect.DeepEqual(app1.Description, app2.Description)
}
