// Package entities 数据库存储基本结构与方法
package entities

import (
	"time"

	"github.com/baetyl/baetyl-go/v2/json"
	"github.com/jinzhu/copier"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

type NodeDevice struct {
	ID             int64     `db:"id"`
	Name           string    `db:"name"`
	Namespace      string    `db:"namespace"`
	Version        string    `db:"version"`
	NodeName       string    `db:"node_name"`
	DeviceModel    string    `db:"device_model"`
	AccessTemplate string    `db:"access_template"`
	DriverName     string    `db:"driver_name"`
	DriverInstName string    `db:"driver_inst_name"`
	Config         string    `db:"config"`
	CreateTime     time.Time `db:"create_time"`
	UpdateTime     time.Time `db:"update_time"`
}

func FromNodeDevice(dev *models.NodeDevice) (*NodeDevice, error) {
	if dev.Config == nil {
		dev.Config = new(models.DeviceConfig)
	}
	cfg, err := json.Marshal(dev.Config)
	if err != nil {
		return nil, err
	}
	dev.Version = GenResourceVersion()
	device := &NodeDevice{
		Config: string(cfg),
	}
	if err = copier.Copy(device, dev); err != nil {
		return nil, err
	}
	return device, nil
}

func ToNodeDevice(dev *NodeDevice) (*models.NodeDevice, error) {
	var cfg models.DeviceConfig
	if err := json.Unmarshal([]byte(dev.Config), &cfg); err != nil {
		return nil, err
	}
	device := &models.NodeDevice{
		Config: &cfg,
	}
	if err := copier.Copy(device, dev); err != nil {
		return nil, err
	}
	return device, nil
}
