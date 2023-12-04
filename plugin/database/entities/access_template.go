// Package entities 数据库存储基本结构与方法
package entities

import (
	"time"

	"github.com/baetyl/baetyl-go/v2/dmcontext"
	"github.com/baetyl/baetyl-go/v2/json"
	"github.com/jinzhu/copier"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

type AccessTemplate struct {
	ID          int64     `db:"id"`
	Name        string    `db:"name"`
	Namespace   string    `db:"namespace"`
	Version     string    `db:"version"`
	Description string    `db:"description"`
	Protocol    string    `db:"protocol"`
	Labels      string    `db:"labels"`
	DeviceModel string    `db:"device_model"`
	Mappings    string    `db:"mappings"`
	Properties  string    `db:"properties"`
	CreateTime  time.Time `db:"create_time"`
	UpdateTime  time.Time `db:"update_time"`
}

func FromAccessTemplate(template *models.AccessTemplate) (*AccessTemplate, error) {
	labels, err := json.Marshal(template.Labels)
	if err != nil {
		return nil, err
	}
	mappings, err := json.Marshal(template.Mappings)
	if err != nil {
		return nil, err
	}
	props, err := json.Marshal(template.Properties)
	if err != nil {
		return nil, err
	}
	template.Version = GenResourceVersion()
	accessTemplate := &AccessTemplate{
		Labels:     string(labels),
		Mappings:   string(mappings),
		Properties: string(props),
	}
	if err = copier.Copy(accessTemplate, template); err != nil {
		return nil, err
	}
	return accessTemplate, nil
}

func ToAccessTemplate(template *AccessTemplate) (*models.AccessTemplate, error) {
	var labels map[string]string
	if err := json.Unmarshal([]byte(template.Labels), &labels); err != nil {
		return nil, err
	}
	var mappings []dmcontext.ModelMapping
	if err := json.Unmarshal([]byte(template.Mappings), &mappings); err != nil {
		return nil, err
	}
	var props []dmcontext.DeviceProperty
	if err := json.Unmarshal([]byte(template.Properties), &props); err != nil {
		return nil, err
	}
	accessTemplate := &models.AccessTemplate{
		Labels:     labels,
		Mappings:   mappings,
		Properties: props,
	}
	if err := copier.Copy(accessTemplate, template); err != nil {
		return nil, err
	}
	return accessTemplate, nil
}
