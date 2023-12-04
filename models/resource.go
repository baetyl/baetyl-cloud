// Package models 模型定义
package models

type ResourceList struct {
	Total        int `json:"total"`
	*ListOptions `json:",inline"`
	Items        []string `json:"items"`
	SysItems     []string `json:"sysItems"`
}
