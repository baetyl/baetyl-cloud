// Package models 模型定义
package models

import (
	"time"

	"github.com/baetyl/baetyl-cloud/v2/common"
)

type Batch struct {
	Name            string            `json:"name,omitempty" binding:"omitempty,res_name"`
	Namespace       string            `json:"namespace,omitempty"`
	Description     string            `json:"description,omitempty"`
	Accelerator     string            `json:"accelerator,omitempty"`
	SysApps         []string          `json:"sysApps,omitempty"`
	QuotaNum        string            `json:"quotaNum,omitempty"`
	EnableWhitelist int               `json:"enableWhitelist,omitempty"`
	Cluster         int               `json:"cluster,omitempty"`
	SecurityType    common.Security   `json:"securityType,omitempty"`
	SecurityKey     string            `json:"securityKey,omitempty"`
	CallbackName    string            `json:"callbackName,omitempty"`
	Labels          map[string]string `json:"labels,omitempty" binding:"omitempty,label"`
	Fingerprint     Fingerprint       `json:"fingerprint,omitempty"`
	CreateTime      time.Time         `json:"createTime,omitempty"`
	UpdateTime      time.Time         `json:"updateTime,omitempty"`
}

type Fingerprint struct {
	Type       int    `json:"type,omitempty" validate:"oneof=1 2 4 8 16 32"`
	SnPath     string `json:"snPath,omitempty"`
	InputField string `json:"inputField,omitempty"`
}

type BatchList struct {
	Total        int `json:"total"`
	*ListOptions `json:",inline"`
	Items        []Batch `json:"items"`
}
