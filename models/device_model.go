// Package models 模型定义
package models

import (
	"bytes"
	"encoding/json"
	"reflect"
	"time"

	dm "github.com/baetyl/baetyl-go/v2/dmcontext"
	"github.com/jinzhu/copier"
)

const (
	BieGatewayProduct = "BIE-Product"

	DeviceTypePublish   = 0x0
	DeviceTypeDraft     = 0x1
	DeviceTypeDirect    = 0x2
	DeviceTypeGateway   = 0x3
	DeviceTypeSub       = 0x4
	DeviceTypeComponent = 0x5
)

var (
	NodeAttributes = []DeviceModelAttribute{
		{
			Name: "地址",
			ID:   "address",
			Type: dm.TypeString,
		},
		{
			Name: "主机名称",
			ID:   "hostName",
			Type: dm.TypeString,
		},
		{
			Name: "架构",
			ID:   "arch",
			Type: dm.TypeString,
		},
		{
			Name: "系统",
			ID:   "os",
			Type: dm.TypeString,
		},
		{
			Name: "系统唯一标识",
			ID:   "systemUUID",
			Type: dm.TypeString,
		},
		{
			Name: "节点模式",
			ID:   "nodeMode",
			Type: dm.TypeString,
		},
		{
			Name: "网络接收字节",
			ID:   "netBytesRecv",
			Type: dm.TypeString,
			Unit: UnitStorageCapacityByte,
		},
		{
			Name: "网络发送字节",
			ID:   "netBytesSent",
			Type: dm.TypeString,
			Unit: UnitStorageCapacityByte,
		},
		{
			Name: "网络接收包",
			ID:   "netPackRecv",
			Type: dm.TypeString,
		},
		{
			Name: "网络发送包",
			ID:   "netPackSent",
			Type: dm.TypeString,
		},
		{
			Name: "磁盘使用量",
			ID:   "diskUse",
			Type: dm.TypeString,
			Unit: UnitStorageCapacityByte,
		},
		{
			Name: "磁盘容量",
			ID:   "distCap",
			Type: dm.TypeString,
			Unit: UnitStorageCapacityByte,
		},
		{
			Name: "内存使用量",
			ID:   "memoryUse",
			Type: dm.TypeString,
			Unit: UnitStorageCapacityByte,
		},
		{
			Name: "内存容量",
			ID:   "memoryCap",
			Type: dm.TypeString,
			Unit: UnitStorageCapacityByte,
		},
		{
			Name: "处理器使用量",
			ID:   "cpuUse",
			Type: dm.TypeString,
		},
		{
			Name: "处理器容量",
			ID:   "cpuCap",
			Type: dm.TypeString,
		},
		{
			Name: "状态",
			ID:   "state",
			Type: dm.TypeString,
		},
	}
)

type DeviceModelList struct {
	Items        []DeviceModel `json:"items"`
	Total        int           `json:"total"`
	*ListOptions `json:",inline"`
}

type DeviceModelView struct {
	Name        string                 `json:"name,omitempty" binding:"omitempty,res_name"`
	Description string                 `json:"description,omitempty"`
	Protocol    string                 `json:"protocol,omitempty"`
	Labels      map[string]string      `json:"labels,omitempty"`
	Attributes  []DeviceModelAttribute `json:"attributes,omitempty"`
	Properties  []DeviceModelProperty  `json:"properties,omitempty"`
	CreateTime  time.Time              `json:"createTime,omitempty"`
	UpdateTime  time.Time              `json:"updateTime,omitempty"`
}

type DeviceModel struct {
	Name          string                 `json:"name,omitempty" binding:"omitempty,dev_model"`
	Namespace     string                 `json:"namespace,omitempty"`
	ProductSecret string                 `json:"productSecret,omitempty"`
	Version       string                 `json:"version,omitempty"`
	Description   string                 `json:"description,omitempty"`
	Protocol      string                 `json:"protocol,omitempty"`
	Labels        map[string]string      `json:"labels,omitempty"`
	Attributes    []DeviceModelAttribute `json:"attributes,omitempty" binding:"dive"`
	Properties    []DeviceModelProperty  `json:"properties,omitempty" binding:"dive"`
	CreateTime    time.Time              `json:"createTime,omitempty"`
	UpdateTime    time.Time              `json:"updateTime,omitempty"`
	Type          byte                   `json:"type"`
}

func (dm *DeviceModel) ToDeviceModelView() (*DeviceModelView, error) {
	res := new(DeviceModelView)
	err := copier.Copy(res, dm)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func EqualDeviceModel(old, new *DeviceModel) bool {
	if old.Name != new.Name || old.Protocol != new.Protocol || old.Description != new.Description {
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
	return true
}

type DeviceModelAttribute struct {
	Name         string      `json:"name,omitempty"`
	ID           string      `json:"id,omitempty"`
	Type         string      `json:"type,omitempty" binding:"data_type"`
	Unit         string      `json:"unit,omitempty"`
	DefaultValue interface{} `json:"defaultValue,omitempty"`
	Required     bool        `json:"required"`
}

type DeviceModelProperty struct {
	Name           string                   `json:"name,omitempty"`
	ID             string                   `json:"id,omitempty" binding:"required"`
	Type           string                   `json:"type,omitempty" binding:"data_plus_type"`
	Mode           string                   `json:"mode,omitempty" binding:"oneof=ro rw"`
	Unit           string                   `json:"unit,omitempty"`
	Format         string                   `json:"format,omitempty"`                    // 当 Type 为 date/time 时使用
	EnumType       dm.EnumType              `json:"enumType,omitempty" binding:"dive"`   // 当 Type 为 enum 时使用
	ArrayType      dm.ArrayType             `json:"arrayType,omitempty" binding:"dive"`  // 当 Type 为 array 时使用
	ObjectType     map[string]dm.ObjectType `json:"objectType,omitempty" binding:"dive"` // 当 Type 为 object 时使用
	ObjectRequired []string                 `json:"objectRequired,omitempty"`            // 当 Type 为 object 时, 记录必填字段
	Visitor        dm.PropertyVisitor       `json:"visitor,omitempty" binding:"dive"`
}

type DeviceModelPropertyYaml struct {
	Name           string                   `json:"name,omitempty" yaml:"name,omitempty"`
	Type           string                   `json:"type,omitempty" yaml:"type,omitempty"`
	Mode           string                   `json:"mode,omitempty" yaml:"mode,omitempty"`
	Format         string                   `json:"format,omitempty" yaml:"format,omitempty"`                        // 当 Type 为 date/time 时使用
	EnumType       dm.EnumType              `json:"enumType,omitempty" yaml:"enumType,omitempty" binding:"dive"`     // 当 Type 为 enum 时使用
	ArrayType      dm.ArrayType             `json:"arrayType,omitempty" yaml:"arrayType,omitempty" binding:"dive"`   // 当 Type 为 array 时使用
	ObjectType     map[string]dm.ObjectType `json:"objectType,omitempty" yaml:"objectType,omitempty" binding:"dive"` // 当 Type 为 object 时使用
	ObjectRequired []string                 `json:"objectRequired,omitempty" yaml:"objectRequired,omitempty"`        // 当 Type 为 object 时, 记录必填字段
}

type deviceModelAttribute struct {
	Name         string      `json:"name,omitempty"`
	ID           string      `json:"id,omitempty"`
	Type         string      `json:"type,omitempty" binding:"data_type"`
	Unit         string      `json:"unit,omitempty"`
	DefaultValue interface{} `json:"defaultValue,omitempty"`
	Required     bool        `json:"required"`
}

func (da *DeviceModelAttribute) UnmarshalJSON(b []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(b))
	decoder.UseNumber()
	var attr deviceModelAttribute
	if err := decoder.Decode(&attr); err != nil {
		return err
	}
	da.Name = attr.Name
	da.ID = attr.ID
	da.Type = attr.Type
	da.Unit = attr.Unit
	da.DefaultValue = attr.DefaultValue
	da.Required = attr.Required
	return nil
}
