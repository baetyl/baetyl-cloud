// Package models 模型定义
package models

import (
	"time"
)

type Record struct {
	Name             string    `json:"name,omitempty" db:"name"`
	Namespace        string    `json:"namespace,omitempty" db:"namespace"`
	BatchName        string    `json:"batchName,omitempty" db:"batch_name"`
	FingerprintValue string    `json:"fingerprintValue,omitempty" db:"fingerprint_value" binding:"omitempty,fingerprint"`
	Active           int       `json:"active,omitempty" db:"active"`
	NodeName         string    `json:"nodeName,omitempty" db:"node_name" binding:"omitempty,res_name"`
	ActiveIP         string    `json:"activeIP,omitempty" db:"active_ip"`
	ActiveTime       time.Time `json:"activeTime,omitempty" db:"active_time"`
	CreateTime       time.Time `json:"createTime,omitempty" db:"create_time"`
	UpdateTime       time.Time `json:"updateTime,omitempty" db:"update_time"`
}

type RecordList struct {
	Total   int `json:"total"`
	*Filter `json:",inline"`
	Items   []Record `json:"items"`
}
