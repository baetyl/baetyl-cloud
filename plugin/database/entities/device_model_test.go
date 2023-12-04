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

func TestTransDeviceModel(t *testing.T) {
	dmm := &models.DeviceModel{
		Name:        "dm-1",
		Namespace:   "default",
		Description: "desc",
		Protocol:    "pro-1",
		Labels:      map[string]string{"a": "b"},
		Attributes: []models.DeviceModelAttribute{{
			Name:         "attr-1",
			ID:           "attr-1",
			Type:         "float64",
			DefaultValue: ejson.Number(strconv.FormatFloat(12.3, 'f', -1, 64)),
			Required:     true,
		}},
		Properties: []models.DeviceModelProperty{{
			Name: "prop-1",
			ID:   "prop-1",
			Type: "float32",
			Mode: "rw",
			Visitor: dm.PropertyVisitor{
				Modbus: &dm.ModbusVisitor{
					Function: 3,
					Address:  "0x3",
				},
			},
		}},
	}
	dbDM, err := FromModelDeviceModel(dmm)
	assert.NoError(t, err)

	_, err = ToModelDeviceModel(dbDM)
	assert.NoError(t, err)

	dbDM.Properties = "test"
	_, err = ToModelDeviceModel(dbDM)
	assert.NotNil(t, err)

	dbDM.Attributes = "test"
	_, err = ToModelDeviceModel(dbDM)
	assert.NotNil(t, err)

	dbDM.Labels = "test"
	_, err = ToModelDeviceModel(dbDM)
	assert.NotNil(t, err)
}
