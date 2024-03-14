// Package database 数据库存储实现
package database

import (
	ejson "encoding/json"
	"fmt"
	"strconv"
	"testing"

	dm "github.com/baetyl/baetyl-go/v2/dmcontext"
	"github.com/baetyl/baetyl-go/v2/log"
	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

var (
	deviceModelTables = []string{
		`
CREATE TABLE baetyl_device_model
(
	id                integer       PRIMARY KEY AUTOINCREMENT,
    name              varchar(128)  NOT NULL DEFAULT '',
    namespace         varchar(128)  NOT NULL DEFAULT '',
	version           varchar(36)   NOT NULL DEFAULT '',
	description       varchar(255)  NOT NULL DEFAULT '',
    protocol          varchar(64)   NOT NULL DEFAULT '',
    type              tinyint(4)    NOT NULL DEFAULT 0,       
    labels            varchar(2048) NOT NULL DEFAULT '',
    attributes        varchar(1024) NOT NULL DEFAULT '',
	properties        varchar(1024) NOT NULL DEFAULT '',
    create_time       timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time       timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`,
	}
)

func (d *BaetylCloudDB) MockCreateDeviceModelTable() {
	for _, sql := range deviceModelTables {
		_, err := d.Exec(nil, sql)
		if err != nil {
			panic(fmt.Sprintf("create device model exception: %s", err.Error()))
		}
	}
}

func TestDeviceModel(t *testing.T) {
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
	listOptions := &models.ListOptions{}
	log.L().Info("Test device model", log.Any("device model", dmm))

	db, err := MockNewDB()
	if err != nil {
		fmt.Printf("get mock sqlite3 error = %s", err.Error())
		t.Fail()
		return
	}
	db.MockCreateDeviceModelTable()

	res, err := db.CreateDeviceModel(userID, dmm)
	assert.NoError(t, err)
	checkDeviceModel(t, dmm, res)

	res, err = db.GetDeviceModel(dmm.Namespace, userID, dmm.Name)
	assert.NoError(t, err)
	checkDeviceModel(t, dmm, res)

	// update device model with equal value
	res, err = db.UpdateDeviceModel(userID, dmm)
	assert.NoError(t, err)
	checkDeviceModel(t, dmm, res)

	dmm.Labels = map[string]string{"b": "b"}
	res, err = db.UpdateDeviceModel(userID, dmm)
	assert.NoError(t, err)
	checkDeviceModel(t, dmm, res)

	res, err = db.GetDeviceModel(dmm.Namespace, userID, dmm.Name)
	assert.NoError(t, err)
	checkDeviceModel(t, dmm, res)

	resList, err := db.ListDeviceModel(dmm.Namespace, userID, listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 1)

	err = db.DeleteDeviceModel(dmm.Namespace, userID, dmm.Name)
	assert.NoError(t, err)

	res, err = db.GetDeviceModel(dmm.Namespace, userID, dmm.Name)
	assert.Nil(t, res)

	_, err = db.ListAllInstance("", userID)
	assert.NoError(t, err)
}

func TestListDeviceModel(t *testing.T) {
	dm1 := &models.DeviceModel{
		Name:        "dm-123-1",
		Namespace:   "default",
		Description: "desc",
		Labels:      map[string]string{"label": "aaa"},
	}
	dm2 := &models.DeviceModel{
		Name:        "dm-abc-1",
		Namespace:   "default",
		Description: "desc",
		Labels:      map[string]string{"label": "aaa"},
	}
	dm3 := &models.DeviceModel{
		Name:        "dm-123-2",
		Namespace:   "default",
		Description: "desc",
		Labels:      map[string]string{"label": "bbb"},
	}
	dm4 := &models.DeviceModel{
		Name:        "dm-abc-2",
		Namespace:   "default",
		Description: "desc",
		Labels:      map[string]string{"label": "bbb"},
	}
	listOptions := &models.ListOptions{}

	db, err := MockNewDB()
	if err != nil {
		fmt.Printf("get mock sqlite3 error = %s", err.Error())
		t.Fail()
		return
	}
	db.MockCreateDeviceModelTable()

	res, err := db.CreateDeviceModel(userID, dm1)
	assert.NoError(t, err)
	checkDeviceModel(t, dm1, res)

	res, err = db.CreateDeviceModel(userID, dm2)
	assert.NoError(t, err)
	checkDeviceModel(t, dm2, res)

	res, err = db.CreateDeviceModel(userID, dm3)
	assert.NoError(t, err)
	checkDeviceModel(t, dm3, res)

	res, err = db.CreateDeviceModel(userID, dm4)
	assert.NoError(t, err)
	checkDeviceModel(t, dm4, res)

	// list option nil, return all cfgs
	resList, err := db.ListDeviceModel("default", userID, listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 4)
	assert.Equal(t, dm1.Name, resList.Items[0].Name)
	assert.Equal(t, dm2.Name, resList.Items[1].Name)
	assert.Equal(t, dm3.Name, resList.Items[2].Name)
	assert.Equal(t, dm4.Name, resList.Items[3].Name)
	// page 1 num 2
	listOptions.PageNo = 1
	listOptions.PageSize = 2
	resList, err = db.ListDeviceModel("default", userID, listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 4)
	assert.Equal(t, dm1.Name, resList.Items[0].Name)
	assert.Equal(t, dm2.Name, resList.Items[1].Name)
	// page 2 num 2
	listOptions.PageNo = 2
	listOptions.PageSize = 2
	resList, err = db.ListDeviceModel("default", userID, listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 4)
	assert.Equal(t, dm3.Name, resList.Items[0].Name)
	assert.Equal(t, dm4.Name, resList.Items[1].Name)
	// page 3 num 0
	listOptions.PageNo = 3
	listOptions.PageSize = 2
	resList, err = db.ListDeviceModel("default", userID, listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 4)
	// page 1 num 2 name like dm
	listOptions.PageNo = 1
	listOptions.PageSize = 2
	listOptions.Name = "dm"
	resList, err = db.ListDeviceModel("default", userID, listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 4)
	assert.Equal(t, dm1.Name, resList.Items[0].Name)
	assert.Equal(t, dm2.Name, resList.Items[1].Name)
	// page 1 num 2 name like abc
	listOptions.PageNo = 1
	listOptions.PageSize = 2
	listOptions.Name = "abc"
	resList, err = db.ListDeviceModel("default", userID, listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 2)
	assert.Equal(t, dm2.Name, resList.Items[0].Name)
	assert.Equal(t, dm4.Name, resList.Items[1].Name)
	// page 1 num2 label : aaa
	listOptions.PageNo = 1
	listOptions.PageSize = 4
	listOptions.Name = "123"
	resList, err = db.ListDeviceModel("default", userID, listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 2)
	assert.Equal(t, dm1.Name, resList.Items[0].Name)
	assert.Equal(t, dm3.Name, resList.Items[1].Name)

	// page 1 num2 label : aaa
	listOptions.PageNo = 1
	listOptions.PageSize = 4
	listOptions.Name = ""
	listOptions.LabelSelector = "label=aaa"
	resList, err = db.ListDeviceModel("default", userID, listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 2)
	assert.Equal(t, dm1.Name, resList.Items[0].Name)
	assert.Equal(t, dm2.Name, resList.Items[1].Name)

	dModels, err := db.ListDeviceModelByNames("default", userID, []string{dm1.Name, dm2.Name})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(dModels))
	assert.Equal(t, dm1.Name, dModels[0].Name)
	assert.Equal(t, dm2.Name, dModels[1].Name)

	count, err := db.CountDeviceModelTx(nil, "default", "%")
	assert.NoError(t, err)
	assert.Equal(t, 4, count)

	err = db.DeleteDeviceModel("default", userID, dm1.Name)
	assert.NoError(t, err)
	err = db.DeleteDeviceModel("default", userID, dm2.Name)
	assert.NoError(t, err)
	err = db.DeleteDeviceModel("default", userID, dm3.Name)
	assert.NoError(t, err)
	err = db.DeleteDeviceModel("default", userID, dm4.Name)
	assert.NoError(t, err)

	res, err = db.GetDeviceModel("default", userID, dm1.Name)
	assert.Nil(t, res)
	res, err = db.GetDeviceModel("default", userID, dm2.Name)
	assert.Nil(t, res)
	res, err = db.GetDeviceModel("default", userID, dm3.Name)
	assert.Nil(t, res)
	res, err = db.GetDeviceModel("default", userID, dm4.Name)
	assert.Nil(t, res)
}

func checkDeviceModel(t *testing.T, expect, actual *models.DeviceModel) {
	assert.Equal(t, expect.Name, actual.Name)
	assert.Equal(t, expect.Namespace, actual.Namespace)
	assert.Equal(t, expect.Description, actual.Description)
	assert.Equal(t, expect.Protocol, actual.Protocol)
	assert.EqualValues(t, expect.Labels, actual.Labels)
	assert.EqualValues(t, expect.Attributes, actual.Attributes)
	assert.EqualValues(t, expect.Properties, actual.Properties)
}
