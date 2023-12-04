// Package database 数据库存储实现
package database

import (
	"fmt"
	"testing"

	"github.com/baetyl/baetyl-go/v2/log"
	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

var (
	nodeDeviceTables = []string{
		`
CREATE TABLE baetyl_node_device
(
	id                integer       PRIMARY KEY AUTOINCREMENT,
    name              varchar(128)  NOT NULL DEFAULT '',
    namespace         varchar(128)  NOT NULL DEFAULT '',
	version           varchar(64)   NOT NULL DEFAULT '',
	access_template   varchar(128)  NOT NULL DEFAULT '',
	device_model      varchar(128)  NOT NULL DEFAULT '',
    node_name         varchar(128)  NOT NULL DEFAULT '',
	driver_name       varchar(128)  NOT NULL DEFAULT '',
	driver_inst_name  varchar(128)  NOT NULL DEFAULT '',
	config            varchar(4096) NOT NULL DEFAULT '',
   	create_time       timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time       timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`,
	}
)

func (d *BaetylCloudDB) MockCreateNodeDeviceTable() {
	for _, sql := range nodeDeviceTables {
		_, err := d.Exec(nil, sql)
		if err != nil {
			panic(fmt.Sprintf("create device exception: %s", err.Error()))
		}
	}
}

func TestNodeDevice(t *testing.T) {
	d := &models.NodeDevice{
		Name:           "dm-1",
		Namespace:      "default",
		DeviceModel:    "model-1",
		NodeName:       "node-1",
		DriverName:     "driver-1",
		DriverInstName: "driver-1-123456789",
		Config:         &models.DeviceConfig{},
	}
	listOptions := &models.ListOptions{}
	log.L().Info("Test device model", log.Any("device model", d))

	db, err := MockNewDB()
	if err != nil {
		fmt.Printf("get mock sqlite3 error = %s", err.Error())
		t.Fail()
		return
	}
	db.MockCreateNodeDeviceTable()

	res, err := db.CreateNodeDevice(d)
	assert.NoError(t, err)
	checkNodeDevice(t, d, res)

	res, err = db.GetNodeDevice(d.Namespace, d.Name, d.DeviceModel)
	assert.NoError(t, err)
	checkNodeDevice(t, d, res)

	// update device with equal value
	res, err = db.UpdateNodeDevice(d)
	assert.NoError(t, err)
	checkNodeDevice(t, d, res)

	res, err = db.GetNodeDevice(d.Namespace, d.Name, d.DeviceModel)
	assert.NoError(t, err)
	checkNodeDevice(t, d, res)

	resList, err := db.ListNodeDeviceByDriverAndNode(d.Namespace, d.DriverInstName, d.NodeName, listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 1)

	newListOptions := &models.ListOptions{Filter: models.Filter{PageNo: 1, PageSize: 10}}
	resList, err = db.ListNodeDeviceByDriverAndNode(d.Namespace, d.DriverInstName, d.NodeName, newListOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 1)

	resList, err = db.ListNodeDeviceByDriverAndNode(d.Namespace, "test", d.NodeName, newListOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 0)

	err = db.DeleteNodeDevice(d.Namespace, d.Name, d.DeviceModel)
	assert.NoError(t, err)

	res, err = db.GetNodeDevice(d.Namespace, d.Name, d.DeviceModel)
	assert.Nil(t, res)

	_, err = db.ListNodeDeviceByNode(d.Namespace, d.NodeName)
	assert.NoError(t, err)

	list := []models.NodeDevice{*d}
	err = db.BatchCreateNodeDevice(list)
	assert.NoError(t, err)

	deNames := []string{d.Name}
	byNames, err := db.BatchGetNodeDevicesByNames(d.Namespace, d.DeviceModel, deNames)
	assert.NoError(t, err)
	assert.Equal(t, byNames[d.Name], d.NodeName)

	err = db.BatchDeleteNodeDevices(userID, "model-1", deNames)
	assert.NoError(t, err)
}

func checkNodeDevice(t *testing.T, expect, actual *models.NodeDevice) {
	assert.Equal(t, expect.Name, actual.Name)
	assert.Equal(t, expect.NodeName, actual.NodeName)
	assert.Equal(t, expect.DriverName, actual.DriverName)
	assert.Equal(t, expect.DriverInstName, actual.DriverInstName)
	assert.Equal(t, expect.DeviceModel, actual.DeviceModel)
	assert.Equal(t, expect.AccessTemplate, actual.AccessTemplate)
	assert.Equal(t, expect.Namespace, actual.Namespace)
	assert.EqualValues(t, expect.Config, actual.Config)
}
