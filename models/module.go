package models

import (
	"time"
)

type Module struct {
	Name              string            `json:"name,omitempty"`
	Version           string            `json:"version,omitempty"`
	Image             string            `json:"image,omitempty"`
	Programs          map[string]string `json:"programs,omitempty"`
	Type              string            `json:"type,omitempty"`
	IsLatest          bool              `json:"isLatest,omitempty"`
	Description       string            `json:"description,omitempty"`
	CreationTimestamp time.Time         `json:"createTime,omitempty"`
	UpdateTimestamp   time.Time         `json:"updateTime,omitempty"`
}

type ModuleView struct {
	Module Module `json:"module,omitempty"`
}
