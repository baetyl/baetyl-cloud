// Package database 数据库存储实现
package database

import (
	"fmt"
	"testing"

	"github.com/baetyl/baetyl-go/v2/dmcontext"
	"github.com/baetyl/baetyl-go/v2/log"
	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

var (
	accessTemplateTables = []string{
		`
CREATE TABLE baetyl_access_template
(
	id                integer       PRIMARY KEY AUTOINCREMENT,
    name              varchar(128)  NOT NULL DEFAULT '',
    namespace         varchar(128)  NOT NULL DEFAULT '',
	version           varchar(36)   NOT NULL DEFAULT '',
	description       varchar(255)  NOT NULL DEFAULT '',
    protocol          varchar(64)   NOT NULL DEFAULT '',
    device_model      varchar(64)   NOT NULL DEFAULT 0,       
    labels            varchar(2048) NOT NULL DEFAULT '',
    mappings          varchar(1024) NOT NULL DEFAULT '',
	properties        varchar(1024) NOT NULL DEFAULT '',
    create_time       timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time       timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`,
	}
)

func (d *BaetylCloudDB) MockCreateAccessTemplateTable() {
	for _, sql := range accessTemplateTables {
		_, err := d.Exec(nil, sql)
		if err != nil {
			panic(fmt.Sprintf("create device model exception: %s", err.Error()))
		}
	}
}

func TestAccessTemplate(t *testing.T) {
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
	listOptions := &models.ListOptions{}
	log.L().Info("Test access template", log.Any("access template", at))

	db, err := MockNewDB()
	if err != nil {
		fmt.Printf("get mock sqlite3 error = %s", err.Error())
		t.Fail()
		return
	}
	db.MockCreateAccessTemplateTable()

	res, err := db.CreateAccessTemplate(at)
	assert.NoError(t, err)
	checkAccessTemplate(t, at, res)

	res, err = db.GetAccessTemplate(at.Namespace, at.Name)
	assert.NoError(t, err)
	checkAccessTemplate(t, at, res)

	// update device model with equal value
	res, err = db.UpdateAccessTemplate(at)
	assert.NoError(t, err)
	checkAccessTemplate(t, at, res)

	at.Labels = map[string]string{"b": "b"}
	res, err = db.UpdateAccessTemplate(at)
	assert.NoError(t, err)
	checkAccessTemplate(t, at, res)

	res, err = db.GetAccessTemplate(at.Namespace, at.Name)
	assert.NoError(t, err)
	checkAccessTemplate(t, at, res)

	resList, err := db.ListAccessTemplate(at.Namespace, listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 1)

	resList, err = db.ListAccessTemplateByModelAndProtocol(at.Namespace, at.DeviceModel, "protocol", listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 0)

	resList, err = db.ListAccessTemplateByModelAndProtocol(at.Namespace, at.DeviceModel, at.Protocol, listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 1)

	err = db.DeleteAccessTemplate(at.Namespace, at.Name)
	assert.NoError(t, err)

	res, err = db.GetAccessTemplate(at.Namespace, at.Name)
	assert.Nil(t, res)
}

func checkAccessTemplate(t *testing.T, expect, actual *models.AccessTemplate) {
	assert.Equal(t, expect.Name, actual.Name)
	assert.Equal(t, expect.Namespace, actual.Namespace)
	assert.Equal(t, expect.Description, actual.Description)
	assert.Equal(t, expect.Protocol, actual.Protocol)
	assert.Equal(t, expect.DeviceModel, actual.DeviceModel)
	assert.EqualValues(t, expect.Labels, actual.Labels)
	assert.EqualValues(t, expect.Mappings, actual.Mappings)
	assert.EqualValues(t, expect.Properties, actual.Properties)
}
