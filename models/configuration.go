package models

import (
	"reflect"
	"time"

	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
)

// ConfigurationList Configuration List
type ConfigurationList struct {
	Total       int                    `json:"total"`
	ListOptions *ListOptions           `json:"listOptions"`
	Items       []specV1.Configuration `json:"items"`
}

type ConfigurationView struct {
	Name              string            `json:"name,omitempty" validate:"resourceName,nonBaetyl"`
	Namespace         string            `json:"namespace,omitempty"`
	Labels            map[string]string `json:"labels,omitempty"`
	Data              []ConfigDataItem  `json:"data,omitempty" default:"[]" validate:"dive"`
	CreationTimestamp time.Time         `json:"createTime,omitempty"`
	UpdateTimestamp   time.Time         `json:"updateTime,omitempty"`
	Description       string            `json:"description,omitempty"`
	Version           string            `json:"version,omitempty"`
	System            bool              `json:"system,omitempty"`
}

type ConfigDataItem struct {
	Key   string            `json:"key,omitempty" validate:"required,validConfigKeys"`
	Value map[string]string `json:"value,omitempty"`
}

type ConfigFunctionItem struct {
	ConfigObjectItem `json:",inline"`
	Function         string `json:"function,omitempty"`
	Version          string `json:"version,omitempty"`
	Runtime          string `json:"runtime,omitempty"`
	Handler          string `json:"handler,omitempty"`
}

type ConfigObjectItem struct {
	Source   string `json:"source,omitempty"`
	Account  string `json:"account,omitempty"`
	Endpoint string `json:"endpoint,omitempty"`
	Bucket   string `json:"bucket,omitempty"`
	Object   string `json:"object,omitempty"`
	Unpack   string `json:"unpack,omitempty"`
	MD5      string `json:"md5,omitempty"`
	Ak       string `json:"ak,omitempty"`
	Sk       string `json:"sk,omitempty"`
}

func EqualConfig(config1, config2 *specV1.Configuration) bool {
	return reflect.DeepEqual(config1.Labels, config2.Labels) &&
		reflect.DeepEqual(config1.Data, config2.Data) &&
		reflect.DeepEqual(config1.Description, config2.Description)
}
