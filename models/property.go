package models

import "time"

type Property struct {
	Key        string    `yaml:"key,omitempty" json:"key,omitempty" db:"key"`
	Value      string    `yaml:"value,omitempty" json:"value,omitempty" db:"value"`
	CreateTime time.Time `yaml:"createTime,omitempty" json:"createTime,omitempty" db:"create_time"`
	UpdateTime time.Time `yaml:"updateTime,omitempty" json:"updateTime,omitempty" db:"update_time"`
}
