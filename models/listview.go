// Package models 模型定义
package models

type ListView struct {
	Items        interface{} `json:"items"`
	*ListOptions `json:",inline"`
	Total        int `json:"total"`
}

type MenuItem struct {
	Text  string `json:"text"`
	Key   string `json:"key"`
	Unit  string `json:"unit"`
	Extra string `json:"extra"`
}
