// Package database 数据库存储实现
package database

import (
	"fmt"
	"testing"

	dm "github.com/baetyl/baetyl-go/v2/dmcontext"
	"github.com/baetyl/baetyl-go/v2/log"
	v1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

var (
	deviceDriverTables = []string{
		`
CREATE TABLE baetyl_device_driver
(
	id                integer       PRIMARY KEY AUTOINCREMENT,
    node_name         varchar(128)  NOT NULL DEFAULT '',
    driver_name       varchar(128)  NOT NULL DEFAULT '',
	driver_inst_name  varchar(128)  NOT NULL DEFAULT '',
    namespace         varchar(128)  NOT NULL DEFAULT '',
	version           varchar(36)   NOT NULL DEFAULT '',
    protocol          varchar(64)   NOT NULL DEFAULT '',
    application       varchar(128)  NOT NULL DEFAULT '',
    configuration     varchar(128)  NOT NULL DEFAULT '',
    driver_config     varchar(1024) NOT NULL DEFAULT '',
   	create_time       timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP,
   	update_time       timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`,
	}
)

func (d *BaetylCloudDB) MockCreateDeviceDriverTable() {
	for _, sql := range deviceDriverTables {
		_, err := d.Exec(nil, sql)
		if err != nil {
			panic(fmt.Sprintf("create device driver exception: %s", err.Error()))
		}
	}
}

func TestDeviceDriver(t *testing.T) {
	d := &models.DeviceDriver{
		NodeName:       "node-1",
		DriverName:     "driver-1",
		DriverInstName: "driver-1-123456789",
		Namespace:      "default",
		Protocol:       "pro-1",
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
	listOptions := &models.ListOptions{}
	log.L().Info("Test device driver", log.Any("device driver", d))

	db, err := MockNewDB()
	if err != nil {
		fmt.Printf("get mock sqlite3 error = %s", err.Error())
		t.Fail()
		return
	}
	db.MockCreateDeviceDriverTable()

	res, err := db.CreateDeviceDriver(nil, d)
	assert.NoError(t, err)
	checkDeviceDriver(t, d, res)

	res, err = db.GetDeviceDriver(nil, d.Namespace, d.NodeName, d.DriverInstName)
	assert.NoError(t, err)
	checkDeviceDriver(t, d, res)

	d.Application = &v1.ObjectReference{
		Name:    "app1",
		Version: "v2",
	}
	res, err = db.UpdateDeviceDriver(nil, d)
	assert.NoError(t, err)
	checkDeviceDriver(t, d, res)

	res, err = db.GetDeviceDriver(nil, d.Namespace, d.NodeName, d.DriverInstName)
	assert.NoError(t, err)
	checkDeviceDriver(t, d, res)

	resList, err := db.ListDeviceDriver(nil, d.Namespace, d.NodeName, listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 1)

	err = db.DeleteDeviceDriver(nil, d.Namespace, d.NodeName, d.DriverInstName)
	assert.NoError(t, err)

	res, err = db.GetDeviceDriver(nil, d.Namespace, d.NodeName, d.DriverInstName)
	assert.Nil(t, res)
}

func TestListDeviceDriver(t *testing.T) {
	d1 := &models.DeviceDriver{
		NodeName:       "node-1",
		DriverName:     "driver-abc-1",
		DriverInstName: "driver-abc-1-123456789",
		Namespace:      "default",
	}
	d2 := &models.DeviceDriver{
		NodeName:       "node-1",
		DriverName:     "driver-123-1",
		DriverInstName: "driver-123-1-123456789",
		Namespace:      "default",
	}
	d3 := &models.DeviceDriver{
		NodeName:       "node-1",
		DriverName:     "driver-abc-2",
		DriverInstName: "driver-abc-2-123456789",
		Namespace:      "default",
	}
	d4 := &models.DeviceDriver{
		NodeName:       "node-1",
		Namespace:      "default",
		DriverName:     "driver-123-2",
		DriverInstName: "driver-123-2-123456789",
	}
	listOptions := &models.ListOptions{}

	db, err := MockNewDB()
	if err != nil {
		fmt.Printf("get mock sqlite3 error = %s", err.Error())
		t.Fail()
		return
	}
	db.MockCreateDeviceDriverTable()

	res, err := db.CreateDeviceDriver(nil, d1)
	assert.NoError(t, err)
	checkDeviceDriver(t, d1, res)

	res, err = db.CreateDeviceDriver(nil, d2)
	assert.NoError(t, err)
	checkDeviceDriver(t, d2, res)

	res, err = db.CreateDeviceDriver(nil, d3)
	assert.NoError(t, err)
	checkDeviceDriver(t, d3, res)

	res, err = db.CreateDeviceDriver(nil, d4)
	assert.NoError(t, err)
	checkDeviceDriver(t, d4, res)

	// list option nil, return all cfgs
	resList, err := db.ListDeviceDriver(nil, "default", "node-1", listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 4)
	assert.Equal(t, d1.DriverInstName, resList.Items[0].DriverInstName)
	assert.Equal(t, d2.DriverInstName, resList.Items[1].DriverInstName)
	assert.Equal(t, d3.DriverInstName, resList.Items[2].DriverInstName)
	assert.Equal(t, d4.DriverInstName, resList.Items[3].DriverInstName)
	// page 1 num 2
	listOptions.PageNo = 1
	listOptions.PageSize = 2
	resList, err = db.ListDeviceDriver(nil, "default", "node-1", listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 4)
	assert.Equal(t, d1.DriverInstName, resList.Items[0].DriverInstName)
	assert.Equal(t, d2.DriverInstName, resList.Items[1].DriverInstName)
	// page 2 num 2
	listOptions.PageNo = 2
	listOptions.PageSize = 2
	resList, err = db.ListDeviceDriver(nil, "default", "node-1", listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 4)
	assert.Equal(t, d3.DriverInstName, resList.Items[0].DriverInstName)
	assert.Equal(t, d4.DriverInstName, resList.Items[1].DriverInstName)
	// page 3 num 0
	listOptions.PageNo = 3
	listOptions.PageSize = 2
	resList, err = db.ListDeviceDriver(nil, "default", "node-1", listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 4)

	resDrivers, err := db.ListDeviceDriverByName(nil, "default", "driver-abc-1")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(resDrivers))
	assert.Equal(t, d1.DriverInstName, resDrivers[0].DriverInstName)

	// page 1 num 2 name like driver
	//listOptions.PageNo = 1
	//listOptions.PageSize = 2
	//listOptions.Name = "driver"
	//resList, err = db.ListDeviceDriver("default", "node-1", &listOptions.Filter)
	//assert.NoError(t, err)
	//assert.Equal(t, resList.Total, 4)
	//assert.Equal(t, d1.DriverName, resList.Items[0].DriverName)
	//assert.Equal(t, d2.DriverName, resList.Items[1].DriverName)
	// page 1 num 2 name like abc
	//listOptions.PageNo = 1
	//listOptions.PageSize = 2
	//listOptions.Name = "abc"
	//resList, err = db.ListDeviceDriver("default", "node-1", &listOptions.Filter)
	//assert.NoError(t, err)
	//assert.Equal(t, resList.Total, 2)
	//assert.Equal(t, d2.DriverName, resList.Items[0].DriverName)
	//assert.Equal(t, d4.DriverName, resList.Items[1].DriverName)
	// page 1 num2 label : aaa
	//listOptions.PageNo = 1
	//listOptions.PageSize = 4
	//listOptions.Name = "123"
	//resList, err = db.ListDeviceDriver("default", "node-1", &listOptions.Filter)
	//assert.NoError(t, err)
	//assert.Equal(t, resList.Total, 2)
	//assert.Equal(t, d1.DriverName, resList.Items[0].DriverName)
	//assert.Equal(t, d3.DriverName, resList.Items[1].DriverName)

	err = db.DeleteDeviceDriver(nil, "default", d1.NodeName, d1.DriverInstName)
	assert.NoError(t, err)
	err = db.DeleteDeviceDriver(nil, "default", d2.NodeName, d2.DriverInstName)
	assert.NoError(t, err)
	err = db.DeleteDeviceDriver(nil, "default", d3.NodeName, d3.DriverInstName)
	assert.NoError(t, err)
	err = db.DeleteDeviceDriver(nil, "default", d4.NodeName, d4.DriverInstName)
	assert.NoError(t, err)

	res, err = db.GetDeviceDriver(nil, "default", d1.NodeName, d1.DriverInstName)
	assert.Nil(t, res)
	res, err = db.GetDeviceDriver(nil, "default", d2.NodeName, d2.DriverInstName)
	assert.Nil(t, res)
	res, err = db.GetDeviceDriver(nil, "default", d3.NodeName, d3.DriverInstName)
	assert.Nil(t, res)
	res, err = db.GetDeviceDriver(nil, "default", d4.NodeName, d4.DriverInstName)
	assert.Nil(t, res)
}

func checkDeviceDriver(t *testing.T, expect, actual *models.DeviceDriver) {
	assert.Equal(t, expect.NodeName, actual.NodeName)
	assert.Equal(t, expect.DriverName, actual.DriverName)
	assert.Equal(t, expect.Namespace, actual.Namespace)
	assert.Equal(t, expect.Protocol, actual.Protocol)
	assert.EqualValues(t, expect.DriverConfig, actual.DriverConfig)
}
