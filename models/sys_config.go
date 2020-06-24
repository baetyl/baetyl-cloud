package models

import "time"

type SysConfig struct {
	Type       string    `yaml:"type,omitempty" json:"type,omitempty" db:"type"`
	Key        string    `yaml:"key,omitempty" json:"key,omitempty" db:"name"`
	Value      string    `yaml:"value,omitempty" json:"value,omitempty" db:"value"`
	CreateTime time.Time `yaml:"createTime,omitempty" json:"createTime,omitempty" db:"create_time"`
	UpdateTime time.Time `yaml:"updateTime,omitempty" json:"updateTime,omitempty" db:"update_time"`
}

type SysConfigView struct {
	SysConfigs []SysConfig `json:"sysconfigs"`
}
