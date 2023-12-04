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
	deviceTables = []string{
		`
CREATE TABLE baetyl_device
(
	id                integer       PRIMARY KEY AUTOINCREMENT,
    name              varchar(128)  NOT NULL DEFAULT '',
    namespace         varchar(128)  NOT NULL DEFAULT '',
	version           varchar(64)   NOT NULL DEFAULT '',
	description       varchar(255)  NOT NULL DEFAULT '',
	ready             tinyint(1)    NOT NULL DEFAULT 0,
    active            tinyint(1)    NOT NULL DEFAULT 0,
    protocol          varchar(64)   NOT NULL DEFAULT '',
    labels            varchar(2048) NOT NULL DEFAULT '',
	alias             varchar(128)  NOT NULL DEFAULT '',
	device_model      varchar(128)  NOT NULL DEFAULT '',
    node_name         varchar(128)  NOT NULL DEFAULT '',
	driver_name       varchar(128)  NOT NULL DEFAULT '',
    attributes        varchar(2048) NOT NULL DEFAULT '',
	properties        varchar(2048) NOT NULL DEFAULT '',
	shadow            varchar(64)   NOT NULL DEFAULT '',
	config            varchar(4096) NOT NULL DEFAULT '',
   	create_time       timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time       timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP
);

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

	userID = "default"
)

func (d *BaetylCloudDB) MockCreateDeviceTable() {
	for _, sql := range deviceTables {
		_, err := d.Exec(nil, sql)
		if err != nil {
			panic(fmt.Sprintf("create device exception: %s", err.Error()))
		}
	}
}

func TestBatchCreatDevice(t *testing.T) {
	d := &models.Device{
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
	log.L().Info("Test device model", log.Any("device model", d))

	db, err := MockNewDB()
	if err != nil {
		fmt.Printf("get mock sqlite3 error = %s", err.Error())
		t.Fail()
		return
	}
	db.MockCreateDeviceTable()

	res, err := db.BatchCreateDevice(userID, []*models.Device{d})
	assert.NoError(t, err)
	checkDevice(t, d, &res[0])
}

func TestDevice(t *testing.T) {
	d := &models.Device{
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
	listOptions := &models.ListOptions{}
	log.L().Info("Test device model", log.Any("device model", d))

	db, err := MockNewDB()
	if err != nil {
		fmt.Printf("get mock sqlite3 error = %s", err.Error())
		t.Fail()
		return
	}
	db.MockCreateDeviceTable()

	res, err := db.CreateDevice(userID, d)
	assert.NoError(t, err)
	checkDevice(t, d, res)

	res, err = db.GetDevice(d.Namespace, userID, d.Name)
	assert.NoError(t, err)
	checkDevice(t, d, res)

	// update device with equal value
	res, err = db.UpdateDevice(userID, d)
	assert.NoError(t, err)
	checkDevice(t, d, res)

	d.Labels = map[string]string{"b": "b"}
	res, err = db.UpdateDevice(userID, d)
	assert.NoError(t, err)
	checkDevice(t, d, res)

	res, err = db.GetDevice(d.Namespace, userID, d.Name)
	assert.NoError(t, err)
	checkDevice(t, d, res)

	resList, err := db.ListDevice(d.Namespace, userID, listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 1)

	err = db.DeleteDevice(d.Namespace, userID, d.Name)
	assert.NoError(t, err)

	res, err = db.GetDevice(d.Namespace, userID, d.Name)
	assert.Nil(t, res)

	names := []string{"dm-1", "dm-2"}
	err = db.UpdateDeviceStateByName(false, userID, names)
	assert.NoError(t, err)

	res, err = db.CreateDevice(userID, d)
	assert.NoError(t, err)
	checkDevice(t, d, res)

	nodeDevice := &models.NodeDevice{
		Name:           "dm-1",
		Namespace:      "default",
		DeviceModel:    "model-1",
		NodeName:       "node-1",
		DriverName:     "driver-1",
		DriverInstName: "driver-1-123456789",
		Config:         &models.DeviceConfig{},
	}

	_, err = db.CreateNodeDevice(nodeDevice)
	assert.NoError(t, err)
	d.Labels = map[string]string{"1": "2"}
	devList := []models.Device{*d}
	err = db.BatchUpdateDeviceAndDeleteNodeDevice(userID, devList)
	assert.NoError(t, err)
}

func TestListDevice(t *testing.T) {
	d1 := &models.Device{
		Name:        "device-123-1",
		Namespace:   "default",
		Description: "desc",
		Alias:       "d1",
		Protocol:    "protocol-1",
		DriverName:  "driver-1",
		DeviceModel: "model-1",
		NodeName:    "node-1",
		Labels:      map[string]string{"label": "aaa"},
	}
	d2 := &models.Device{
		Name:        "device-abc-1",
		Namespace:   "default",
		Description: "desc",
		Alias:       "d2",
		Protocol:    "protocol-2",
		DriverName:  "driver-2",
		DeviceModel: "model-3",
		NodeName:    "node-1",
		Labels:      map[string]string{"label": "aaa"},
	}
	d3 := &models.Device{
		Name:        "device-123-2",
		Namespace:   "default",
		Description: "desc",
		Protocol:    "protocol-1",
		DriverName:  "driver-1",
		NodeName:    "node-2",
		DeviceModel: "model-1",
		Labels:      map[string]string{"label": "bbb"},
	}
	d4 := &models.Device{
		Name:        "device-abc-2",
		Namespace:   "default",
		Description: "desc",
		NodeName:    "node-2",
		Protocol:    "protocol-2",
		DriverName:  "driver-2",
		DeviceModel: "model-4",
		Labels:      map[string]string{"label": "bbb"},
	}
	listOptions := &models.ListOptions{}

	db, err := MockNewDB()
	if err != nil {
		fmt.Printf("get mock sqlite3 error = %s", err.Error())
		t.Fail()
		return
	}
	db.MockCreateDeviceTable()

	res, err := db.CreateDevice(userID, d1)
	assert.NoError(t, err)
	checkDevice(t, d1, res)

	res, err = db.CreateDevice(userID, d2)
	assert.NoError(t, err)
	checkDevice(t, d2, res)

	res, err = db.CreateDevice(userID, d3)
	assert.NoError(t, err)
	checkDevice(t, d3, res)

	res, err = db.CreateDevice(userID, d4)
	assert.NoError(t, err)
	checkDevice(t, d4, res)

	// list option nil, return all cfgs
	resList, err := db.ListDevice("default", userID, listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 4)
	assert.Equal(t, d1.Name, resList.Items[0].Name)
	assert.Equal(t, d2.Name, resList.Items[1].Name)
	assert.Equal(t, d3.Name, resList.Items[2].Name)
	assert.Equal(t, d4.Name, resList.Items[3].Name)
	// page 1 num 2
	listOptions.PageNo = 1
	listOptions.PageSize = 2
	resList, err = db.ListDevice("default", userID, listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 4)
	assert.Equal(t, d1.Name, resList.Items[0].Name)
	assert.Equal(t, d2.Name, resList.Items[1].Name)
	// page 2 num 2
	listOptions.PageNo = 2
	listOptions.PageSize = 2
	resList, err = db.ListDevice("default", userID, listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 4)
	assert.Equal(t, d3.Name, resList.Items[0].Name)
	assert.Equal(t, d4.Name, resList.Items[1].Name)
	// page 3 num 0
	listOptions.PageNo = 3
	listOptions.PageSize = 2
	resList, err = db.ListDevice("default", userID, listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 0)
	// page 1 num 2 name like device
	listOptions.PageNo = 1
	listOptions.PageSize = 2
	listOptions.Keyword = "device"
	resList, err = db.ListDevice("default", userID, listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 4)
	assert.Equal(t, d1.Name, resList.Items[0].Name)
	assert.Equal(t, d2.Name, resList.Items[1].Name)
	// page 1 num 2 name like abc
	listOptions.PageNo = 1
	listOptions.PageSize = 2
	listOptions.Name = "abc"
	resList, err = db.ListDevice("default", userID, listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 2)
	assert.Equal(t, d2.Name, resList.Items[0].Name)
	assert.Equal(t, d4.Name, resList.Items[1].Name)
	// page 1 num2 label : aaa
	listOptions.PageNo = 1
	listOptions.PageSize = 4
	listOptions.Name = "123"
	resList, err = db.ListDevice("default", userID, listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 2)
	assert.Equal(t, d1.Name, resList.Items[0].Name)
	assert.Equal(t, d3.Name, resList.Items[1].Name)

	// list by alias
	listOptions.PageNo = 1
	listOptions.PageSize = 2
	listOptions.Name = ""
	listOptions.Alias = "d"
	resList, err = db.ListDevice("default", userID, listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 2)
	assert.Equal(t, d1.Name, resList.Items[0].Name)
	assert.Equal(t, d2.Name, resList.Items[1].Name)

	// list device by protocol
	listOptions.PageNo = 1
	listOptions.PageSize = 4
	resList, err = db.ListDeviceByProtocol("default", userID, "protocol-1", listOptions)
	assert.NoError(t, err)
	assert.Equal(t, d1.Name, resList.Items[0].Name)
	assert.Equal(t, d3.Name, resList.Items[1].Name)

	// list device by device model
	listOptions.LabelSelector = "label=aaa"
	resDevices, err := db.ListDeviceByDeviceModel("default", userID, "model-1", listOptions)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(resDevices))
	assert.Equal(t, d1.Name, resDevices[0].Name)

	// list device by device model and protocol
	listOptions.LabelSelector = ""
	listOptions.PageNo = 1
	listOptions.PageSize = 4
	listOptions.Alias = ""
	resList, err = db.ListDeviceByDriverAndNode("default", userID, "driver-2", "node-2", listOptions)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(resList.Items))
	assert.Equal(t, d4.Name, resList.Items[0].Name)

	// batch update device and bind to node
	deviceNames := []string{d1.Name, d3.Name}
	nodeName := "node-4"
	driverName := "driver-1"
	devList, err := db.BatchUpdateDeviceNodeAndDriver("default", userID, nodeName, driverName, deviceNames)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(devList))
	assert.Equal(t, d1.Name, devList[0].Name)
	assert.Equal(t, d3.Name, devList[1].Name)

	resList, err = db.ListDeviceByDriverAndNode("default", userID, driverName, nodeName, listOptions)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(resList.Items))
	assert.Equal(t, d1.Name, resList.Items[0].Name)
	assert.Equal(t, d3.Name, resList.Items[1].Name)

	/// batch update device properties
	deviceNames = []string{d2.Name, d4.Name}
	properties := []dm.DeviceProperty{{ID: "prop-1"}}
	err = db.BatchUpdateDeviceProps("default", userID, deviceNames, properties)
	assert.NoError(t, err)
	res, err = db.GetDevice(d2.Namespace, userID, d2.Name)
	assert.NoError(t, err)
	assert.Equal(t, properties, res.Properties)
	res, err = db.GetDevice(d4.Namespace, userID, d4.Name)
	assert.NoError(t, err)
	assert.Equal(t, properties, res.Properties)

	// batch update device attributes
	deviceNames = []string{d2.Name, d4.Name}
	attributes := []models.DeviceAttribute{{ID: "attr-1"}}
	err = db.BatchUpdateDeviceAttrs("default", userID, deviceNames, attributes)
	assert.NoError(t, err)
	res, err = db.GetDevice(d2.Namespace, userID, d2.Name)
	assert.NoError(t, err)
	assert.Equal(t, attributes, res.Attributes)
	res, err = db.GetDevice(d4.Namespace, userID, d4.Name)
	assert.NoError(t, err)
	assert.Equal(t, attributes, res.Attributes)

	names := map[string]string{d2.Name: d2.DeviceModel}
	deMap, err := db.GetDevicesByDeviceModelAndNameList("default", userID, names)
	assert.NoError(t, err)
	assert.Equal(t, d2.Name, deMap[d2.Name].Name)

	err = db.DeleteDevice(d1.Name, userID, "default")
	assert.NoError(t, err)
	err = db.DeleteDevice(d2.Name, userID, "default")
	assert.NoError(t, err)
	err = db.DeleteDevice(d3.Name, userID, "default")
	assert.NoError(t, err)
	err = db.DeleteDevice(d4.Name, userID, "default")
	assert.NoError(t, err)

	res, err = db.GetDevice(d1.Name, userID, "default")
	assert.Nil(t, res)
	res, err = db.GetDevice(d2.Name, userID, "default")
	assert.Nil(t, res)
	res, err = db.GetDevice(d3.Name, userID, "default")
	assert.Nil(t, res)
	res, err = db.GetDevice(d4.Name, userID, "default")
	assert.Nil(t, res)

	d3.Active = true
	list := []models.Device{*d3}
	err = db.BatchUpdateDevice(userID, list)
	assert.NoError(t, err)

	_, err = db.GetDeviceNumByDeviceModel(d4.Namespace, d4.DeviceModel)
	assert.NoError(t, err)

	devNames := []string{d4.Name}
	de, err := db.BatchGetDeviceByNames(d4.Namespace, d4.Namespace, devNames)
	assert.NoError(t, err)
	assert.Equal(t, d4.DeviceModel, de[d4.Name])

	devNames = []string{}
	de, err = db.BatchGetDeviceByNames(d4.Namespace, d4.Namespace, devNames)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(de))
}

func checkDevice(t *testing.T, expect, actual *models.Device) {
	assert.Equal(t, expect.Name, actual.Name)
	assert.Equal(t, expect.Alias, actual.Alias)
	assert.Equal(t, expect.NodeName, actual.NodeName)
	assert.Equal(t, expect.DriverName, actual.DriverName)
	assert.Equal(t, expect.DeviceModel, actual.DeviceModel)
	assert.Equal(t, expect.Shadow, actual.Shadow)
	assert.Equal(t, expect.Namespace, actual.Namespace)
	assert.Equal(t, expect.Description, actual.Description)
	assert.Equal(t, expect.Protocol, actual.Protocol)
	assert.EqualValues(t, expect.Labels, actual.Labels)
	assert.EqualValues(t, expect.Attributes, actual.Attributes)
	assert.EqualValues(t, expect.Properties, actual.Properties)
	assert.EqualValues(t, expect.Config, actual.Config)
}
