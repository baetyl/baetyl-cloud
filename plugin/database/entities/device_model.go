// Package entities 数据库存储基本结构与方法
package entities

import (
	"time"

	"github.com/baetyl/baetyl-go/v2/json"
	"github.com/jinzhu/copier"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

type DeviceModel struct {
	ID          int64     `db:"id"`
	Name        string    `db:"name"`
	Namespace   string    `db:"namespace"`
	Version     string    `db:"version"`
	Description string    `db:"description"`
	Protocol    string    `db:"protocol"`
	Labels      string    `db:"labels"`
	Attributes  string    `db:"attributes"`
	Properties  string    `db:"properties"`
	Type        byte      `db:"type"`
	CreateTime  time.Time `db:"create_time"`
	UpdateTime  time.Time `db:"update_time"`
}

func FromModelDeviceModel(devModel *models.DeviceModel) (*DeviceModel, error) {
	labels, err := json.Marshal(devModel.Labels)
	if err != nil {
		return nil, err
	}
	attrs, err := json.Marshal(devModel.Attributes)
	if err != nil {
		return nil, err
	}
	props, err := json.Marshal(devModel.Properties)
	if err != nil {
		return nil, err
	}
	devModel.Version = GenResourceVersion()
	deviceModel := &DeviceModel{
		Labels:     string(labels),
		Attributes: string(attrs),
		Properties: string(props),
	}
	if err = copier.Copy(deviceModel, devModel); err != nil {
		return nil, err
	}
	return deviceModel, nil
}

func ToModelDeviceModel(devModel *DeviceModel) (*models.DeviceModel, error) {
	var labels map[string]string
	if err := json.Unmarshal([]byte(devModel.Labels), &labels); err != nil {
		return nil, err
	}
	var attrs []models.DeviceModelAttribute
	if err := json.Unmarshal([]byte(devModel.Attributes), &attrs); err != nil {
		return nil, err
	}
	var props []models.DeviceModelProperty
	if err := json.Unmarshal([]byte(devModel.Properties), &props); err != nil {
		return nil, err
	}
	deviceModel := &models.DeviceModel{
		Labels:     labels,
		Attributes: attrs,
		Properties: props,
	}
	if err := copier.Copy(deviceModel, devModel); err != nil {
		return nil, err
	}
	return deviceModel, nil
}
