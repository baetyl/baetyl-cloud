package models

import "time"

type Task struct {
	TraceId    string    `json:"trace_id,omitempty" db:"trace_id"`
	Namespace  string    `json:"namespace,omitempty" db:"namespace"`
	Node       string    `json:"node,omitempty" db:"node"`
	Type       string    `json:"type,omitempty" db:"type"`
	State      string    `json:"state,omitempty" db:"state"`
	Step       string    `json:"step,omitempty" db:"step"`
	OldVersion string    `json:"old_version,omitempty" db:"old_version"`
	NewVersion string    `json:"new_version,omitempty" db:"new_version"`
	CreateTime time.Time `json:"createTime,omitempty" db:"create_time"`
	UpdateTime time.Time `json:"updateTime,omitempty" db:"update_time"`
}
