// Package entities 数据库存储基本结构与方法
package entities

import (
	"time"

	"github.com/baetyl/baetyl-go/v2/json"
	v1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/jinzhu/copier"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

type DeviceDriver struct {
	ID             int64     `db:"id"`
	NodeName       string    `db:"node_name"`
	DriverName     string    `db:"driver_name"`
	DriverInstName string    `db:"driver_inst_name"`
	Namespace      string    `db:"namespace"`
	Version        string    `db:"version"`
	Protocol       string    `db:"protocol"`
	Application    string    `db:"application"`
	Configuration  string    `db:"configuration"`
	DriverConfig   string    `db:"driver_config"`
	CreateTime     time.Time `db:"create_time"`
	UpdateTime     time.Time `db:"update_time"`
}

func FromModelDeviceDriver(dDriver *models.DeviceDriver) (*DeviceDriver, error) {
	if dDriver.Application == nil {
		dDriver.Application = new(v1.ObjectReference)
	}
	app, err := json.Marshal(dDriver.Application)
	if err != nil {
		return nil, err
	}
	if dDriver.Configuration == nil {
		dDriver.Configuration = new(v1.ObjectReference)
	}
	cfg, err := json.Marshal(dDriver.Configuration)
	if err != nil {
		return nil, err
	}
	if dDriver.DriverConfig == nil {
		dDriver.DriverConfig = new(models.DriverConfig)
	}
	driverConfig, err := json.Marshal(dDriver.DriverConfig)
	if err != nil {
		return nil, err
	}
	dDriver.Version = GenResourceVersion()
	deviceDriver := &DeviceDriver{
		Application:   string(app),
		Configuration: string(cfg),
		DriverConfig:  string(driverConfig),
	}
	if err := copier.Copy(deviceDriver, dDriver); err != nil {
		return nil, err
	}
	return deviceDriver, nil
}

func ToModelDeviceDriver(dDriver *DeviceDriver) (*models.DeviceDriver, error) {
	var app v1.ObjectReference
	if err := json.Unmarshal([]byte(dDriver.Application), &app); err != nil {
		return nil, err
	}
	var cfg v1.ObjectReference
	if err := json.Unmarshal([]byte(dDriver.Configuration), &cfg); err != nil {
		return nil, err
	}
	var driverConfig models.DriverConfig
	if err := json.Unmarshal([]byte(dDriver.DriverConfig), &driverConfig); err != nil {
		return nil, err
	}
	deviceDriver := &models.DeviceDriver{
		Application:   &app,
		Configuration: &cfg,
		DriverConfig:  &driverConfig,
	}
	if err := copier.Copy(deviceDriver, dDriver); err != nil {
		return nil, err
	}
	return deviceDriver, nil
}
