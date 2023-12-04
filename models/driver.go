// Package models 模型定义
package models

import (
	"reflect"
	"time"

	v1 "github.com/baetyl/baetyl-go/v2/spec/v1"
)

const IpcImgPath = "/var/lib/baetyl/image"

type DriverList struct {
	Items        []Driver `json:"items"`
	*ListOptions `json:",inline"`
	Total        int `json:"total"`
}

func EqualDriver(old, new *Driver) bool {
	if old.Name != new.Name || old.Type != new.Type || old.Protocol != new.Protocol ||
		old.Architecture != new.Architecture || old.Description != new.Description ||
		old.DefaultConfig != new.DefaultConfig || old.ProgramConfig != new.ProgramConfig {
		return false
	}
	if !equalService(old.Service, new.Service) {
		return false
	}
	if len(old.Volumes) != len(new.Volumes) || !reflect.DeepEqual(old.Volumes, new.Volumes) {
		return false
	}
	if len(old.Registries) != len(new.Registries) || !reflect.DeepEqual(old.Registries, new.Registries) {
		return false
	}
	if len(old.Labels) != len(new.Labels) || !reflect.DeepEqual(old.Labels, new.Labels) {
		return false
	}
	return true
}

func equalService(oldSvc, newSvc *Service) bool {
	if oldSvc.Image != newSvc.Image || oldSvc.HostNetwork != newSvc.HostNetwork {
		return false
	}
	if !reflect.DeepEqual(oldSvc.SecurityContext, newSvc.SecurityContext) {
		return false
	}
	if len(oldSvc.VolumeMounts) != len(newSvc.VolumeMounts) || (len(newSvc.VolumeMounts) > 0 && !reflect.DeepEqual(oldSvc.VolumeMounts, newSvc.VolumeMounts)) {
		return false
	}
	if len(oldSvc.Args) != len(newSvc.Args) || (len(newSvc.Args) > 0 && !reflect.DeepEqual(oldSvc.Args, newSvc.Args)) {
		return false
	}
	if len(oldSvc.Ports) != len(newSvc.Ports) || (len(newSvc.Ports) > 0 && !reflect.DeepEqual(oldSvc.Ports, newSvc.Ports)) {
		return false
	}
	if len(oldSvc.Env) != len(newSvc.Env) || (len(newSvc.Env) > 0 && !reflect.DeepEqual(oldSvc.Env, newSvc.Env)) {
		return false
	}
	if !equalResource(oldSvc.Resources, newSvc.Resources) {
		return false
	}
	return true
}

func equalResource(oldResources, newResources *v1.Resources) bool {
	// since old and new resources might be {}, {requests:{}}, {limits:{}}, {request:{}, limits:{}} or null
	// all of them equals to {}
	newRes := new(v1.Resources)
	if newResources != nil {
		if len(newResources.Limits) > 0 {
			newRes.Limits = newResources.Limits
		}
		if len(newResources.Requests) > 0 {
			newRes.Requests = newResources.Requests
		}
	}
	oldRes := new(v1.Resources)
	if oldResources != nil {
		if len(oldResources.Limits) > 0 {
			oldRes.Limits = oldResources.Limits
		}
		if len(oldResources.Requests) > 0 {
			oldRes.Requests = oldResources.Requests
		}
	}
	return reflect.DeepEqual(oldRes, newRes)
}

type Driver struct {
	Name          string            `json:"name,omitempty" binding:"omitempty,res_name"`
	Namespace     string            `json:"namespace,omitempty"`
	Version       string            `json:"version,omitempty"`
	Type          byte              `json:"type,omitempty"`
	Mode          string            `json:"mode,omitempty" default:"kube"`
	Labels        map[string]string `json:"labels,omitempty"`
	Protocol      string            `json:"protocol,omitempty"`
	Architecture  string            `json:"arch,omitempty"`
	Description   string            `json:"description,omitempty"`
	DefaultConfig string            `json:"defaultConfig,omitempty"`
	*Service      `json:",inline,omitempty"`
	Volumes       []v1.Volume    `json:"volumes,omitempty"`
	Registries    []RegistryView `json:"registries,omitempty"`
	CreateTime    time.Time      `json:"createTime,omitempty"`
	UpdateTime    time.Time      `json:"updateTime,omitempty"`
	ProgramConfig string         `json:"programConfig,omitempty"`
}

type Service struct {
	Image           string              `json:"image,omitempty"`
	Resources       *v1.Resources       `json:"resources,omitempty"`
	Ports           []v1.ContainerPort  `json:"ports,omitempty"`
	Env             []v1.Environment    `json:"env,omitempty"`
	SecurityContext *v1.SecurityContext `json:"security,omitempty"`
	HostNetwork     bool                `json:"hostNetwork,omitempty"`
	Args            []string            `json:"args,omitempty"`
	VolumeMounts    []v1.VolumeMount    `json:"volumeMounts,omitempty"`
}
