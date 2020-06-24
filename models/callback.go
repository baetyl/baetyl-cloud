package models

import (
	"time"
)

type Callback struct {
	Name        string            `json:"name,omitempty" validate:"omitempty,resourceName"`
	Namespace   string            `json:"namespace,omitempty"`
	Method      string            `json:"method,omitempty" binding:"required"`
	Url         string            `json:"url,omitempty" binding:"required"`
	Params      map[string]string `json:"params,omitempty"`
	Header      map[string]string `json:"header,omitempty"`
	Body        map[string]string `json:"body,omitempty"`
	Description string            `json:"description,omitempty"`
	CreateTime  time.Time         `json:"createTime,omitempty"`
	UpdateTime  time.Time         `json:"updateTime,omitempty"`
}
