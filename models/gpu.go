// Package models 模型定义
package models

type GPUInfoInput struct {
	GPUMetrics bool     `json:"gpuMetrics"`
	Subs       []string `json:"subs"`
}

type GPUInfo struct {
	GPUMetrics bool             `json:"gpuMetrics"`
	Subs       []SubNodeGPUInfo `json:"subs,omitempty"`
}

type SubNodeGPUInfo struct {
	NodeName   string            `json:"nodeName"`
	GPUMetrics bool              `json:"gpuMetrics"`
	Ready      bool              `json:"ready"`
	Labels     map[string]string `json:"labels"`
}
