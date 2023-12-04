package entities

import "time"

type NodeConfiguration struct {
	ID                int64     `json:"ID" db:"id"`
	NodeName          string    `json:"nodeName" db:"node_name"`
	Namespace         string    `json:"namespace" db:"namespace"`
	Data              string    `json:"data" db:"data"`
	ConfigurationType string    `json:"configurationType" db:"type"`
	CreateTime        time.Time `json:"createTime" db:"create_time"`
	UpdateTime        time.Time `json:"updateTime" db:"update_time"`
}
