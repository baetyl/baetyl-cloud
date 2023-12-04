// Package models 模型定义
package models

import (
	"reflect"
	"time"

	dm "github.com/baetyl/baetyl-go/v2/dmcontext"
	v1 "github.com/baetyl/baetyl-go/v2/spec/v1"
)

type DeviceDriverList struct {
	Items        []DeviceDriver `json:"items"`
	*ListOptions `json:",inline"`
	Total        int `json:"total"`
}

type DriverRef struct {
	Name        string `json:"name,omitempty"`
	Type        byte   `json:"type,omitempty"`
	Protocol    string `json:"protocol,omitempty"`
	Description string `json:"description,omitempty"`
	Bind        bool   `json:"bind"`
}

func EqualDeviceDriver(old, new *DeviceDriver) bool {
	if old.NodeName != new.NodeName || old.DriverName != new.DriverName || old.DriverInstName != new.DriverInstName || old.Protocol != new.Protocol {
		return false
	}
	if !reflect.DeepEqual(old.Application, new.Application) || !reflect.DeepEqual(old.Configuration, new.Configuration) ||
		!reflect.DeepEqual(old.DriverConfig, new.DriverConfig) {
		return false
	}
	return true
}

type DeviceDriver struct {
	NodeName       string              `yaml:"nodeName,omitempty" json:"nodeName,omitempty"`
	DriverName     string              `yaml:"driverName,omitempty" json:"driverName,omitempty"`
	DriverInstName string              `yaml:"driverInstName,omitempty" json:"driverInstName,omitempty"`
	Namespace      string              `yaml:"namespace,omitempty" json:"namespace,omitempty"`
	Version        string              `yaml:"version,omitempty" json:"version,omitempty"`
	Protocol       string              `yaml:"protocol,omitempty" json:"protocol,omitempty"`
	Application    *v1.ObjectReference `yaml:"application,omitempty" json:"application,omitempty"`
	Configuration  *v1.ObjectReference `yaml:"configuration,omitempty" json:"configuration,omitempty"`
	DriverConfig   *DriverConfig       `yaml:"driverConfig,omitempty" json:"driverConfig,omitempty"`
	CreateTime     time.Time           `yaml:"createTime,omitempty" json:"createTime,omitempty"`
	UpdateTime     time.Time           `yaml:"updateTime,omitempty" json:"updateTime,omitempty"`
}

type DeviceDriverView struct {
	NodeName       string        `yaml:"nodeName,omitempty" json:"nodeName,omitempty"`
	DriverName     string        `yaml:"driverName,omitempty" json:"driverName,omitempty"`
	DriverInstName string        `yaml:"driverInstName,omitempty" json:"driverInstName,omitempty"`
	Version        string        `yaml:"version,omitempty" json:"version,omitempty"`
	Protocol       string        `yaml:"protocol,omitempty" json:"protocol,omitempty"`
	DriverConfig   *DriverConfig `yaml:"driverConfig,omitempty" json:"driverConfig,omitempty"`
	CreateTime     time.Time     `yaml:"createTime,omitempty" json:"createTime,omitempty"`
	UpdateTime     time.Time     `yaml:"updateTime,omitempty" json:"updateTime,omitempty"`
}

type DriverConfig struct {
	Channels    []ChannelConfig       `yaml:"channels,omitempty" json:"channels,omitempty"`
	IpcServices []dm.IpcServiceConfig `yaml:"ipcServices,omitempty" json:"ipcServices,omitempty"`
	Custom      *CustomDriverConfig   `yaml:"custom,omitempty" json:"custom,omitempty"`
}

type CustomDriverConfig string

type DeviceDriverCsv struct {
	DeviceName     string
	DeviceModel    string
	AccessTemplate string
	Custom         *CustomDeviceConfig
	Modbus         *ModbusConfig
	OpcUA          *OpcuaConfig
	Ipc            *dm.IpcDeviceConfig
	IEC104         *IEC104Config
	Opcda          *OpcdaConfig
	Bacnet         *BacnetConfig
}
