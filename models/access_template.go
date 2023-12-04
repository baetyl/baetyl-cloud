// Package models 模型定义
package models

import (
	"reflect"
	"time"

	dm "github.com/baetyl/baetyl-go/v2/dmcontext"
	"github.com/jinzhu/copier"
)

type AccessTemplateList struct {
	Items        []AccessTemplate `json:"items"`
	Total        int              `json:"total"`
	*ListOptions `json:",inline"`
}

type AccessTemplate struct {
	Name        string              `json:"name"`
	Namespace   string              `json:"namespace"`
	Version     string              `json:"version,omitempty"`
	Description string              `json:"description,omitempty"`
	Labels      map[string]string   `json:"labels,omitempty"`
	Protocol    string              `json:"protocol" validate:"oneof=modbus opcua opcda iec104 ipc custom bacnet ice101"`
	DeviceModel string              `json:"deviceModel,omitempty"`
	Properties  []dm.DeviceProperty `json:"properties,omitempty"`
	Mappings    []dm.ModelMapping   `json:"mappings,omitempty"`
	CreateTime  time.Time           `json:"createTime,omitempty"`
	UpdateTime  time.Time           `json:"updateTime,omitempty"`
}

type AccessTemplateView struct {
	Name        string              `json:"name"`
	Description string              `json:"description,omitempty"`
	Labels      map[string]string   `json:"labels,omitempty"`
	Protocol    string              `json:"protocol,omitempty"`
	DeviceModel string              `json:"deviceModel,omitempty"`
	Properties  []dm.DeviceProperty `json:"properties,omitempty"`
	Mappings    []dm.ModelMapping   `json:"mappings,omitempty"`
	CreateTime  time.Time           `json:"createTime,omitempty"`
	UpdateTime  time.Time           `json:"updateTime,omitempty"`
}

func EqualAccessTemplate(old, new *AccessTemplate) bool {
	if old.Name != new.Name || old.Protocol != new.Protocol || old.Description != new.Description {
		return false
	}
	if len(old.Labels) != len(new.Labels) || !reflect.DeepEqual(old.Labels, new.Labels) {
		return false
	}
	if len(old.Mappings) != len(new.Mappings) || !reflect.DeepEqual(old.Mappings, new.Mappings) {
		return false
	}
	if len(old.Properties) != len(new.Properties) || !reflect.DeepEqual(old.Properties, new.Properties) {
		return false
	}
	return true
}

func (at *AccessTemplate) ToAccessTemplateView() (*AccessTemplateView, error) {
	res := new(AccessTemplateView)
	err := copier.Copy(res, at)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func ToAccessTemplateListView(templates *AccessTemplateList) *ListView {
	var templateList []AccessTemplateView
	for _, m := range templates.Items {
		template := AccessTemplateView{
			Name:        m.Name,
			Description: m.Description,
			Protocol:    m.Protocol,
			Labels:      m.Labels,
			DeviceModel: m.DeviceModel,
			CreateTime:  m.CreateTime,
			UpdateTime:  m.UpdateTime,
		}
		templateList = append(templateList, template)
	}
	return &ListView{
		Total:       templates.Total,
		ListOptions: templates.ListOptions,
		Items:       templateList,
	}
}
