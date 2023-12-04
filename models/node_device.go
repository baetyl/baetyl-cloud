// Package models 模型定义
package models

import (
	"reflect"
	"time"

	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/jinzhu/copier"
)

type NodeDevice struct {
	Name           string        `json:"name,omitempty"`
	Namespace      string        `json:"namespace,omitempty"`
	Version        string        `json:"version,omitempty"`
	DeviceModel    string        `json:"deviceModel,omitempty"`
	AccessTemplate string        `json:"accessTemplate,omitempty"`
	NodeName       string        `json:"nodeName,omitempty"`
	DriverName     string        `json:"driverName,omitempty"`
	DriverInstName string        `json:"driverInstName,omitempty"`
	Config         *DeviceConfig `json:"config,omitempty"`
	CreateTime     time.Time     `json:"createTime,omitempty"`
	UpdateTime     time.Time     `json:"updateTime,omitempty"`
}

type NodeBindDeviceView struct {
	Name           string            `json:"name,omitempty"`
	DeviceModel    string            `json:"deviceModel,omitempty"`
	AccessTemplate string            `json:"accessTemplate,omitempty"`
	NodeName       string            `json:"nodeName,omitempty"`
	DriverName     string            `json:"driverName,omitempty"`
	DriverInstName string            `json:"driverInstName,omitempty"`
	Status         int               `json:"status"` // 0: unknown 1: offline, 2: online
	Bind           bool              `json:"bind"`   // due to node name is empty
	Config         *DeviceConfigView `json:"config,omitempty"`
}

func (n *NodeDevice) ToNodeBindDeviceView() (*NodeBindDeviceView, error) {
	var view NodeBindDeviceView
	if err := copier.Copy(&view, n); err != nil {
		return nil, errors.Trace(err)
	}
	return &view, nil
}

type NodeDeviceList struct {
	Items        []NodeDevice `json:"items"`
	*ListOptions `json:",inline"`
	Total        int `json:"total"`
}

func EqualNodeDevice(old, new *NodeDevice) bool {
	if old.Name != new.Name || old.AccessTemplate != new.AccessTemplate || old.DeviceModel != new.DeviceModel {
		return false
	}
	if !reflect.DeepEqual(old.Config, new.Config) {
		return false
	}
	if old.NodeName != new.NodeName || old.DriverName != new.DriverName {
		return false
	}
	return true
}
