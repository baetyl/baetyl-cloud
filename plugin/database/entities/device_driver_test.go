// Package entities 数据库存储基本结构与方法
package entities

import (
	"testing"

	dm "github.com/baetyl/baetyl-go/v2/dmcontext"
	v1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

func TestTransDeviceDriver(t *testing.T) {
	device := &models.DeviceDriver{
		NodeName:   "node-1",
		DriverName: "driver-1",
		Namespace:  "default",
		Protocol:   "pro-1",
		Application: &v1.ObjectReference{
			Name:    "app1",
			Version: "v1",
		},
		Configuration: &v1.ObjectReference{
			Name:    "cfg1",
			Version: "c1",
		},
		DriverConfig: &models.DriverConfig{
			Channels: []models.ChannelConfig{
				{
					ChannelID: "1",
					Modbus: &models.ModbusChannel{
						TCP: &dm.TCPConfig{
							Address: "local",
							Port:    50200,
						},
					},
				},
			},
		},
	}
	dbDevice, err := FromModelDeviceDriver(device)
	assert.NoError(t, err)

	_, err = ToModelDeviceDriver(dbDevice)
	assert.NoError(t, err)

	dbDevice.DriverConfig = "test"
	_, err = ToModelDeviceDriver(dbDevice)
	assert.NotNil(t, err)

	dbDevice.Configuration = "test"
	_, err = ToModelDeviceDriver(dbDevice)
	assert.NotNil(t, err)

	dbDevice.Application = "test"
	_, err = ToModelDeviceDriver(dbDevice)
	assert.NotNil(t, err)
}
