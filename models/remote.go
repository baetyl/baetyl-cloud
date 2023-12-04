// Package models 模型定义
package models

type Remote struct {
	System    bool   `json:"system"`
	Instance  string `json:"instance"`
	Container string `json:"container"`
	Mode      string `json:"mode"`
}

type RemoteDescribe struct {
	System       bool   `json:"system"`
	Instance     string `json:"instance"`
	Mode         string `json:"mode"`
	ResourceType string `json:"resourceType"`
}
