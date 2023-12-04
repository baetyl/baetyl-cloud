// Package entities 数据库存储基本结构与方法
package entities

import (
	"testing"

	"github.com/baetyl/baetyl-go/v2/dmcontext"
	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

func TestTransAccessTemplate(t *testing.T) {
	at := &models.AccessTemplate{
		Name:        "dm-1",
		Namespace:   "default",
		Description: "desc",
		Protocol:    "pro-1",
		DeviceModel: "dm-1",
		Labels:      map[string]string{"a": "b"},
		Mappings: []dmcontext.ModelMapping{{
			Attribute:  "attr-1",
			Expression: "equal(x1)",
		}},
		Properties: []dmcontext.DeviceProperty{{
			Name: "prop-1",
			ID:   "1",
			Type: "float32",
			Mode: "rw",
			Visitor: dmcontext.PropertyVisitor{
				Modbus: &dmcontext.ModbusVisitor{
					Function: 3,
					Address:  "0x3",
				},
			},
		}},
	}
	dbAT, err := FromAccessTemplate(at)
	assert.NoError(t, err)

	_, err = ToAccessTemplate(dbAT)
	assert.NoError(t, err)

	dbAT.Properties = "test"
	_, err = ToAccessTemplate(dbAT)
	assert.NotNil(t, err)

	dbAT.Mappings = "test"
	_, err = ToAccessTemplate(dbAT)
	assert.NotNil(t, err)

	dbAT.Labels = "test"
	_, err = ToAccessTemplate(dbAT)
	assert.NotNil(t, err)
}
