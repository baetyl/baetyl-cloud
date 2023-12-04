// Package entities 数据库存储基本结构与方法
package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

func TestTransNodeDevice(t *testing.T) {
	d := &models.NodeDevice{
		Name:        "dm-1",
		Namespace:   "default",
		DeviceModel: "model-1",
		NodeName:    "node-1",
		DriverName:  "driver-1",
		Config:      &models.DeviceConfig{},
	}
	dbNodeDevice, err := FromNodeDevice(d)
	assert.NoError(t, err)

	_, err = ToNodeDevice(dbNodeDevice)
	assert.NoError(t, err)

	dbNodeDevice.Config = "test"
	_, err = ToNodeDevice(dbNodeDevice)
	assert.NotNil(t, err)
}
