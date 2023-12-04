// Package entities 数据库存储基本结构与方法
package entities

import (
	ejson "encoding/json"
	"strconv"
	"testing"

	dm "github.com/baetyl/baetyl-go/v2/dmcontext"
	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

func TestTransDevice(t *testing.T) {
	device := &models.Device{
		Name:        "dm-1",
		Alias:       "dm-1",
		Namespace:   "default",
		Description: "desc",
		Protocol:    "pro-1",
		Labels:      map[string]string{"a": "b"},
		Shadow:      "shad",
		DeviceModel: "model-1",
		NodeName:    "node-1",
		DriverName:  "driver-1",
		Attributes: []models.DeviceAttribute{{
			Name:     "attr-1",
			ID:       "attr-1",
			Type:     "float64",
			Required: true,
			Value:    ejson.Number(strconv.FormatFloat(12.23, 'f', -1, 64)),
		}},
		Properties: []dm.DeviceProperty{{
			Name:    "prop-1",
			ID:      "prop-1",
			Type:    "float32",
			Mode:    "rw",
			Expect:  ejson.Number(strconv.FormatFloat(1.23, 'f', -1, 64)),
			Current: ejson.Number(strconv.FormatFloat(2.45, 'f', -1, 64)),
		}},
		Config: &models.DeviceConfig{},
	}
	dbDevcie, err := FromModelDevice(device)
	assert.NoError(t, err)

	_, err = ToModelDevice(dbDevcie)
	assert.NoError(t, err)

	dbDevcie.Config = "test"
	_, err = ToModelDevice(dbDevcie)
	assert.NotNil(t, err)

	dbDevcie.Properties = "test"
	_, err = ToModelDevice(dbDevcie)
	assert.NotNil(t, err)

	dbDevcie.Attributes = "test"
	_, err = ToModelDevice(dbDevcie)
	assert.NotNil(t, err)

	dbDevcie.Labels = "test"
	_, err = ToModelDevice(dbDevcie)
	assert.NotNil(t, err)
}
