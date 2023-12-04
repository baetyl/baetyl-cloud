package models

import "time"

type NodeConfiguration struct {
	NodeName          string    `json:"nodeName" db:"node_name"`
	Namespace         string    `json:"namespace" db:"namespace"`
	Data              string    `json:"data" db:"data"`
	ConfigurationType string    `json:"configurationType" db:"type"`
	CreateTime        time.Time `json:"createTime" db:"create_time"`
	UpdateTime        time.Time `json:"updateTime" db:"update_time"`
}
