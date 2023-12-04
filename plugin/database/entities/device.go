// Package entities 数据库存储基本结构与方法
package entities

import (
	"time"

	"github.com/baetyl/baetyl-go/v2/dmcontext"
	"github.com/baetyl/baetyl-go/v2/json"
	"github.com/jinzhu/copier"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

type Device struct {
	ID          int64     `db:"id"`
	Name        string    `db:"name"`
	Namespace   string    `db:"namespace"`
	Version     string    `db:"version"`
	Ready       bool      `db:"ready"`
	Active      bool      `db:"active"`
	Description string    `db:"description"`
	Protocol    string    `db:"protocol"`
	Labels      string    `db:"labels"`
	Alias       string    `db:"alias"`
	DeviceModel string    `db:"device_model"`
	NodeName    string    `db:"node_name"`
	DriverName  string    `db:"driver_name"`
	Attributes  string    `db:"attributes"`
	Properties  string    `db:"properties"`
	Shadow      string    `db:"shadow"`
	Config      string    `db:"config"`
	CreateTime  time.Time `db:"create_time"`
	UpdateTime  time.Time `db:"update_time"`
}

func FromModelDevice(dev *models.Device) (*Device, error) {
	if dev.Labels == nil {
		dev.Labels = make(map[string]string)
	}
	labels, err := json.Marshal(dev.Labels)
	if err != nil {
		return nil, err
	}
	if dev.Attributes == nil {
		dev.Attributes = make([]models.DeviceAttribute, 0)
	}
	attrs, err := json.Marshal(dev.Attributes)
	if err != nil {
		return nil, err
	}
	if dev.Properties == nil {
		dev.Properties = make([]dmcontext.DeviceProperty, 0)
	}
	props, err := json.Marshal(dev.Properties)
	if err != nil {
		return nil, err
	}
	if dev.Config == nil {
		dev.Config = new(models.DeviceConfig)
	}
	cfg, err := json.Marshal(dev.Config)
	if err != nil {
		return nil, err
	}
	dev.Version = GenResourceVersion()
	device := &Device{
		Labels:     string(labels),
		Attributes: string(attrs),
		Properties: string(props),
		Config:     string(cfg),
	}
	if err := copier.Copy(device, dev); err != nil {
		return nil, err
	}
	return device, nil
}

func ToModelDevice(dev *Device) (*models.Device, error) {
	var labels map[string]string
	if err := json.Unmarshal([]byte(dev.Labels), &labels); err != nil {
		return nil, err
	}
	var attrs []models.DeviceAttribute
	if err := json.Unmarshal([]byte(dev.Attributes), &attrs); err != nil {
		return nil, err
	}
	var props []dmcontext.DeviceProperty
	if err := json.Unmarshal([]byte(dev.Properties), &props); err != nil {
		return nil, err
	}
	var cfg models.DeviceConfig
	if err := json.Unmarshal([]byte(dev.Config), &cfg); err != nil {
		return nil, err
	}
	device := &models.Device{
		Labels:     labels,
		Attributes: attrs,
		Properties: props,
		Config:     &cfg,
	}
	if err := copier.Copy(device, dev); err != nil {
		return nil, err
	}
	return device, nil
}
