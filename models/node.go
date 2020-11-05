package models

import (
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
)

// NodeViewList node view list
type NodeViewList struct {
	Total       int               `json:"total"`
	ListOptions *ListOptions      `json:"listOptions"`
	Items       []specV1.NodeView `json:"items"`
}

// NodeList node list
type NodeList struct {
	Total       int           `json:"total"`
	ListOptions *ListOptions  `json:"listOptions"`
	Items       []specV1.Node `json:"items"`
}

type ListOptions struct {
	LabelSelector string `json:"selector,omitempty"`
	FieldSelector string `json:"fieldSelector,omitempty"`
	Limit         int64  `json:"limit,omitempty"`
	Continue      string `json:"continue,omitempty"`
}

type NodeNames struct {
	Names []string `json:"names,"validate:"maxLength=20"`
}

type NodeProperties struct {
	Properties []NodeProperty `yaml:"properties,omitempty" json:"properties,omitempty"`
}

type NodeProperty struct {
	Name    string        `yaml:"name,omitempty" json:"name,omitempty"`
	Type    string        `yaml:"type,omitempty" json:"type,omitempty"`
	Current PropertyValue `yaml:"current,omitempty" json:"current,omitempty"`
	Expect  PropertyValue `yaml:"expect,omitempty" json:"expect,omitempty"`
}

type PropertyValue struct {
	Value      string `yaml:"value,omitempty" json:"value,omitempty"`
	UpdateTime string `yaml:"updateTime,omitempty" json:"updateTime,omitempty"`
}

type NodeMode struct {
	Mode string `yaml:"mode,omitempty" json:"mode,omitempty"`
}

type NodePropertiesMeta struct {
	ReportMeta map[string]interface{}
	DesireMeta map[string]interface{}
}
