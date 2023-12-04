// Package models 模型定义
package models

import (
	"time"
)

const (
	ErrDmpResourceNotFound = "ResourceNotFoundException"
)

type DMPLinkInfo struct {
	Entrypoint   string
	ProductKey   string
	DeviceName   string
	DeviceSecret string
	InstanceID   string
}

type DMPInstance struct {
	Name       string `json:"name,omitempty"`
	State      string `json:"state,omitempty"`
	InstanceID string `json:"instanceId,omitempty"`
}

type DMPNamespaceList struct {
	Items        []DMPNamespace `json:"items"`
	*ListOptions `json:",inline"`
	Total        int `json:"total"`
}

type DMPNamespace struct {
	Name       string    `json:"name"`
	Namespace  string    `json:"namespace"`
	InstanceID string    `json:"instanceId"`
	Config     DMPConfig `json:"dmpConfig"`
	CreateTime time.Time `json:"createTime"`
	UpdateTime time.Time `json:"updateTime"`
}

type DMPNamespaceViewList struct {
	Items        []DMPNamespaceView `json:"items"`
	*ListOptions `json:",inline"`
	Total        int `json:"total"`
}

type DMPNamespaceView struct {
	Name       string    `json:"name"`
	Namespace  string    `json:"namespace"`
	InstanceID string    `json:"instanceId"`
	CreateTime time.Time `json:"createTime"`
	UpdateTime time.Time `json:"updateTime"`
}

type DMPConfig struct {
	ManagerEndpoint string `json:"managerEndpoint,omitempty"`
	DeviceEndpoint  string `json:"deviceEndpoint,omitempty"`
	Ak              string `json:"ak,omitempty"`
	Sk              string `json:"sk,omitempty"`
}

func (dn *DMPNamespace) ToView() *DMPNamespaceView {
	return &DMPNamespaceView{
		Name:       dn.Name,
		Namespace:  dn.Namespace,
		InstanceID: dn.InstanceID,
		CreateTime: dn.CreateTime,
		UpdateTime: dn.UpdateTime,
	}
}

func (dnl *DMPNamespaceList) ToViewList() *DMPNamespaceViewList {
	res := &DMPNamespaceViewList{
		ListOptions: dnl.ListOptions,
		Total:       dnl.Total,
	}
	for _, dn := range dnl.Items {
		res.Items = append(res.Items, *dn.ToView())
	}
	return res
}
