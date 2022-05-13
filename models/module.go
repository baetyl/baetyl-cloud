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
	Flag              int               `json:"flag"`
	IsLatest          bool              `json:"isLatest,omitempty"`
	Description       string            `json:"description,omitempty"`
	CreationTimestamp time.Time         `json:"createTime,omitempty"`
	UpdateTimestamp   time.Time         `json:"updateTime,omitempty"`
}

type InitCMD struct {
	CMD    string `json:"cmd,omitempty"`
	APK    string `json:"apk,omitempty"`
	APKSys string `json:"apk_sys,omitempty"`
}
