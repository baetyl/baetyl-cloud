// Package models 模型定义
package models

type Lock struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
	TTL     int64  `json:"ttl,omitempty"`
}
