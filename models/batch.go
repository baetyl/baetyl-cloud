package models

import (
	"time"

	"github.com/baetyl/baetyl-cloud/v2/common"
)

type Batch struct {
	Name            string            `json:"name,omitempty" validate:"omitempty,resourceName"`
	Namespace       string            `json:"namespace,omitempty"`
	Description     string            `json:"description,omitempty"`
	QuotaNum        int               `json:"quotaNum,omitempty"`
	EnableWhitelist int               `json:"enableWhitelist,omitempty"`
	SecurityType    common.Security   `json:"securityType,omitempty"`
	SecurityKey     string            `json:"securityKey,omitempty"`
	CallbackName    string            `json:"callbackName,omitempty"`
	Labels          map[string]string `json:"labels,omitempty" validate:"omitempty,validLabels"`
	Fingerprint     Fingerprint       `json:"fingerprint,omitempty"`
	CreateTime      time.Time         `json:"createTime,omitempty"`
	UpdateTime      time.Time         `json:"updateTime,omitempty"`
}

type Fingerprint struct {
	Type       int    `json:"type,omitempty"`
	SnPath     string `json:"snPath,omitempty"`
	InputField string `json:"inputField,omitempty"`
}
