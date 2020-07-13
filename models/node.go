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
