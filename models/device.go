// Package models 模型定义
package models

import (
	"bytes"
	"encoding/json"
	"reflect"
	"time"

	dm "github.com/baetyl/baetyl-go/v2/dmcontext"
	v1 "github.com/baetyl/baetyl-go/v2/spec/v1"
)

const (
	DMPTypeDefault = "defaultdmp"

	UnitStorageCapacityByte     = "字节/B"
	UnitStorageCapacityMegaByte = "兆字节/MB"
	UnitStorageCapacityGigaByte = "吉字节/GB"
)

type MultiDeviceEvent struct {
	Devices []string `json:"devices,omitempty"`
	Event   dm.Event `json:"event,omitempty"`
}

type NodeDevices struct {
	Devices []string `json:"devices,omitempty"`
}

type NodeDevicesWithModel struct {
	Devices []DeviceAndModel `json:"devices,omitempty"`
}

type DeviceAndModel struct {
	Name        string `json:"name"`
	DeviceModel string `json:"deviceModel"`
}

type NodeDeviceView struct {
	Name        string `json:"name"`
	DeviceModel string `json:"deviceModel"`
	NodeName    string `json:"nodeName"`
	Ready       bool   `json:"ready"`
	Bind        bool   `json:"bind"`
	Status      int    `json:"status"`
}

type DeviceShadowList struct {
	Items        []DeviceMetaShadow `json:"items"`
	*ListOptions `json:",inline"`
	Total        int `json:"total"`
}

type DeviceList struct {
	Items        []Device `json:"items"`
	*ListOptions `json:",inline"`
	Total        int `json:"total"`
}

type DeviceWithShadow struct {
	Device `json:",inline"`
	Shadow DeviceMetaShadow `json:"shadow"`
}

type DeviceViewWithShadow struct {
	DeviceView `json:",inline"`
	Shadow     DeviceMetaShadow `json:"shadow"`
}

type DeviceMetaShadow struct {
	Attributes map[string]interface{} `json:"attrs,omitempty"`
	State      ShadowState            `json:"state,omitempty"`
	Metadata   ShadowMetadata         `json:"metadata,omitempty"`
}

type ShadowState struct {
	Report v1.Report `json:"report,omitempty"`
	Desire v1.Desire `json:"desire,omitempty"`
}

type ShadowMetadata struct {
	Report v1.Report `json:"report,omitempty"`
	Desire v1.Desire `json:"desire,omitempty"`
}

type DeviceView struct {
	Name        string              `json:"name,omitempty"`
	Protocol    string              `json:"protocol,omitempty"`
	Description string              `json:"description,omitempty"`
	Labels      map[string]string   `json:"labels,omitempty"`
	Ready       bool                `json:"ready"`
	DeviceModel string              `json:"deviceModel,omitempty"`
	Alias       string              `json:"alias,omitempty"`
	Status      int                 `json:"status"`
	Bind        bool                `json:"bind"`
	NodeName    string              `json:"nodeName,omitempty"`
	Attributes  []DeviceAttribute   `json:"attributes,omitempty"`
	Properties  []dm.DeviceProperty `json:"properties,omitempty"`
	Config      *DeviceConfigView   `json:"config,omitempty"`
	CreateTime  time.Time           `json:"createTime,omitempty"`
	UpdateTime  time.Time           `json:"updateTime,omitempty"`
}

type Device struct {
	Name        string              `json:"name,omitempty" binding:"omitempty,res_name"`
	Namespace   string              `json:"namespace,omitempty"`
	Version     string              `json:"version,omitempty"`
	Description string              `json:"description,omitempty"`
	Protocol    string              `json:"protocol,omitempty"`
	Alias       string              `json:"alias,omitempty"`
	Ready       bool                `json:"ready"`
	Active      bool                `json:"active,omitempty"`
	Labels      map[string]string   `json:"labels,omitempty"`
	DeviceModel string              `json:"deviceModel,omitempty"`
	Attributes  []DeviceAttribute   `json:"attributes,omitempty"`
	Properties  []dm.DeviceProperty `json:"properties,omitempty"`
	// Deprecated: 移动到 baetyl_node_device 表实现
	NodeName string `json:"nodeName,omitempty"`
	// Deprecated: 移动到 baetyl_node_device 表实现
	DriverName string `json:"driverName,omitempty"`
	// Deprecated: 移动到 baetyl_node_device 表实现
	Config     *DeviceConfig `json:"config,omitempty"`
	Shadow     string        `json:"shadow,omitempty"`
	CreateTime time.Time     `json:"createTime,omitempty"`
	UpdateTime time.Time     `json:"updateTime,omitempty"`
}

func EqualDevice(old, new *Device) bool {
	if old.Name != new.Name || old.Protocol != new.Protocol || old.DeviceModel != new.DeviceModel {
		return false
	}
	if old.Alias != new.Alias || old.Description != new.Description {
		return false
	}
	if len(old.Labels) != len(new.Labels) || !reflect.DeepEqual(old.Labels, new.Labels) {
		return false
	}
	if len(old.Attributes) != len(new.Attributes) || !reflect.DeepEqual(old.Attributes, new.Attributes) {
		return false
	}
	if len(old.Properties) != len(new.Properties) || !reflect.DeepEqual(old.Properties, new.Properties) {
		return false
	}
	if !reflect.DeepEqual(old.Config, new.Config) {
		return false
	}
	if old.NodeName != new.NodeName || old.DriverName != new.DriverName || old.Shadow != new.Shadow || old.Ready != new.Ready || old.Active != new.Active {
		return false
	}
	return true
}

type DeviceAttribute struct {
	Name     string      `json:"name,omitempty"`
	ID       string      `json:"id,omitempty"`
	Type     string      `json:"type,omitempty" binding:"data_type"`
	Unit     string      `json:"unit,omitempty"`
	Required bool        `json:"required"`
	Value    interface{} `json:"value"`
}

type deviceAttribute struct {
	Name     string      `json:"name,omitempty"`
	ID       string      `json:"id,omitempty"`
	Type     string      `json:"type,omitempty" binding:"data_type"`
	Unit     string      `json:"unit,omitempty"`
	Required bool        `json:"required"`
	Value    interface{} `json:"value"`
}

func (da *DeviceAttribute) UnmarshalJSON(b []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(b))
	decoder.UseNumber()
	var attr deviceAttribute
	if err := decoder.Decode(&attr); err != nil {
		return err
	}
	da.Name = attr.Name
	da.ID = attr.ID
	da.Unit = attr.Unit
	da.Type = attr.Type
	da.Required = attr.Required
	da.Value = attr.Value
	return nil
}

type DeviceConfig struct {
	Infos *dm.DeviceInfo `yaml:"infos,omitempty" json:"infos,omitempty"`
	// props are not inherited from device model and have not preserved yet,
	// might be updated and saved when device model update in future
	Props  []dm.DeviceProperty `yaml:"props,omitempty" json:"props,omitempty"`
	Modbus *ModbusConfig       `yaml:"modbus,omitempty" json:"modbus,omitempty"`
	OpcUA  *OpcuaConfig        `yaml:"opcua,omitempty" json:"opcua,omitempty"`
	IEC104 *IEC104Config       `yaml:"iec104,omitempty" json:"iec104,omitempty"`
	OpcDA  *OpcdaConfig        `yaml:"opcda,omitempty" json:"opcda,omitempty"`
	Bacnet *BacnetConfig       `yaml:"bacnet,omitempty" json:"bacnet,omitempty"`
	Ipc    *dm.IpcDeviceConfig `yaml:"ipc,omitempty" json:"ipc,omitempty"`
	Custom *CustomDeviceConfig `yaml:"custom,omitempty" json:"custom,omitempty"`
}

type FullDeviceConfig struct {
	Infos *dm.DeviceInfo `yaml:"infos,omitempty" json:"infos,omitempty"`
	Props interface{}    `yaml:"props,omitempty" json:"props,omitempty"`
}

type DeviceConfigView struct {
	Modbus *ModbusConfig       `yaml:"modbus,omitempty" json:"modbus,omitempty"`
	OpcUA  *OpcuaConfig        `yaml:"opcua,omitempty" json:"opcua,omitempty"`
	Ipc    *dm.IpcDeviceConfig `yaml:"ipc,omitempty" json:"ipc,omitempty"`
	IEC104 *IEC104Config       `yaml:"iec104,omitempty" json:"iec104,omitempty"`
	Bacnet *BacnetConfig       `yaml:"bacnet,omitempty" json:"bacnet,omitempty"`
	OpcDA  *OpcdaConfig        `yaml:"opcda,omitempty" json:"opcda,omitempty"`
	Custom *CustomDeviceConfig `yaml:"custom,omitempty" json:"custom,omitempty"`
}

type ChannelConfig struct {
	ChannelID string         `yaml:"channelId,omitempty" json:"channelId,omitempty"`
	Modbus    *ModbusChannel `yaml:"modbus,omitempty" json:"modbus,omitempty"`
	OpcUA     *OpcuaChannel  `yaml:"opcua,omitempty" json:"opcua,omitempty"`
	IEC104    *IEC104Channel `yaml:"iec104,omitempty" json:"iec104,omitempty"`
	Opcda     *OpcdaChannel  `yaml:"opcda,omitempty" json:"opcda,omitempty"`
	Bacnet    *BacnetChannel `yaml:"bacnet,omitempty" json:"bacnet,omitempty"`
}

type ModbusConfig struct {
	ChannelID string `yaml:"channelId,omitempty" json:"channelId,omitempty" binding:"required"`
	SlaveID   byte   `yaml:"slaveId,omitempty" json:"slaveId,omitempty" binding:"required"`
	Interval  int    `yaml:"interval,omitempty" json:"interval,omitempty"` // unit is second
}

type FullDriverConfig struct {
	Devices []dm.DeviceInfo `yaml:"devices" json:"devices"`
	Driver  string          `yaml:"driver" json:"driver"`
}

type OpcuaConfig struct {
	ChannelID string `yaml:"channelId,omitempty" json:"channelId,omitempty" binding:"required"`
	Subscribe bool   `yaml:"subscribe,omitempty" json:"subscribe,omitempty"`
	Interval  int    `yaml:"interval" json:"interval,omitempty"`
	NsOffset  int    `yaml:"nsOffset" json:"nsOffset,omitempty"`
	IDOffset  int    `yaml:"idOffset" json:"idOffset,omitempty"`
}

type IEC104Config struct {
	ChannelID string `yaml:"channelId,omitempty" json:"channelId,omitempty" binding:"required"`
	Interval  int    `yaml:"interval" json:"interval,omitempty"`
	AIOffset  int    `yaml:"aiOffset" json:"aiOffset"`
	DIOffset  int    `yaml:"diOffset" json:"diOffset"`
	AOOffset  int    `yaml:"aoOffset" json:"aoOffset"`
	DOOffset  int    `yaml:"doOffset" json:"doOffset"`
	PIOffset  int    `yaml:"piOffset" json:"piOffset"`
}

type OpcdaConfig struct {
	ChannelID string `yaml:"channelId,omitempty" json:"channelId,omitempty" binding:"required"`
	Interval  int    `yaml:"interval" json:"interval,omitempty"`
}

type BacnetConfig struct {
	ChannelID     string `yaml:"channelId,omitempty" json:"channelId,omitempty" binding:"required"`
	Interval      int    `yaml:"interval,omitempty" json:"interval,omitempty"`
	DeviceID      uint32 `yaml:"deviceId,omitempty" json:"deviceId,omitempty"`
	AddressOffset uint   `yaml:"addressOffset,omitempty" json:"addressOffset,omitempty"`
}

type ModbusChannel struct {
	TCP *dm.TCPConfig `yaml:"tcp,omitempty" json:"tcp,omitempty"`
	RTU *dm.RTUConfig `yaml:"rtu,omitempty" json:"rtu,omitempty"`
}

type OpcuaChannel struct {
	ID          byte                 `yaml:"id,omitempty" json:"id,omitempty"`
	Endpoint    string               `yaml:"endpoint,omitempty" json:"endpoint,omitempty"`
	Timeout     int                  `yaml:"timeout,omitempty" json:"timeout,omitempty"`
	Security    dm.OpcuaSecurity     `yaml:"security,omitempty" json:"security,omitempty"`
	Auth        *dm.OpcuaAuth        `yaml:"auth,omitempty" json:"auth,omitempty"`
	Certificate *dm.OpcuaCertificate `yaml:"certificate,omitempty" json:"certificate,omitempty"`
}

type IEC104Channel struct {
	Protocol string `yaml:"protocol,omitempty" json:"protocol,omitempty"`
	Address  string `yaml:"address,omitempty" json:"address,omitempty" binding:"required"`
	Port     uint16 `yaml:"port,omitempty" json:"port,omitempty" binding:"required"`
}

type OpcdaChannel struct {
	Host      string `yaml:"host,omitempty" json:"host,omitempty" binding:"required"`
	ProgramID string `yaml:"progid,omitempty" json:"progid,omitempty"`
	Clsid     string `yaml:"clsid,omitempty" json:"clsid,omitempty"`
	UserName  string `yaml:"username,omitempty" json:"username,omitempty" binding:"required"`
	Password  string `yaml:"password,omitempty" json:"password,omitempty" binding:"required"`
}

type BacnetChannel struct {
	Address string `yaml:"address,omitempty" json:"address,omitempty"`
	Port    int    `yaml:"port,omitempty" json:"port,omitempty"`
}

type CustomChannel string

type CustomDeviceConfig string
