// Package database 数据库存储实现
package database

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/plugin/database/entities"

	"github.com/jmoiron/sqlx"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

func (d *BaetylCloudDB) GetNodeDevice(ns, name, model string) (*models.NodeDevice, error) {
	return d.GetNodeDeviceTx(nil, ns, name, model)
}

func (d *BaetylCloudDB) CreateNodeDevice(device *models.NodeDevice) (*models.NodeDevice, error) {
	var res *models.NodeDevice
	err := d.Transact(func(tx *sqlx.Tx) error {
		_, err := d.CreateNodeDeviceTx(tx, device)
		if err != nil {
			return err
		}
		res, err = d.GetNodeDeviceTx(tx, device.Namespace, device.Name, device.DeviceModel)
		return err
	})
	return res, err
}

func (d *BaetylCloudDB) UpdateNodeDevice(device *models.NodeDevice) (*models.NodeDevice, error) {
	var res *models.NodeDevice
	err := d.Transact(func(tx *sqlx.Tx) error {
		old, err := d.GetNodeDeviceTx(tx, device.Namespace, device.Name, device.DeviceModel)
		if err != nil {
			return err
		}
		if models.EqualNodeDevice(old, device) {
			res = old
			return nil
		}
		_, err = d.UpdateNodeDeviceTx(tx, device)
		if err != nil {
			return err
		}
		res, err = d.GetNodeDeviceTx(tx, device.Namespace, device.Name, device.DeviceModel)
		return err
	})
	return res, err
}

func (d *BaetylCloudDB) DeleteNodeDevice(ns, name, model string) error {
	_, err := d.DeleteNodeDeviceTx(nil, ns, name, model)
	return err
}

func (d *BaetylCloudDB) ListNodeDeviceByDriverAndNode(ns, driverInstName, nodeName string, listOptions *models.ListOptions) (*models.NodeDeviceList, error) {
	devices, count, err := d.ListNodeDeviceByDriverAndNodeTx(nil, ns, driverInstName, nodeName, listOptions)
	if err != nil {
		return nil, err
	}
	return &models.NodeDeviceList{
		Items:       devices,
		ListOptions: listOptions,
		Total:       count,
	}, nil
}

func (d *BaetylCloudDB) GetNodeDeviceTx(tx *sqlx.Tx, ns, name, model string) (*models.NodeDevice, error) {
	selectSQL := `
SELECT  
name, namespace, version, device_model, node_name, access_template, driver_name, driver_inst_name, config, create_time, update_time
FROM baetyl_node_device WHERE namespace=? AND name=? AND device_model=? LIMIT 0,1
`
	var devices []entities.NodeDevice
	if err := d.Query(tx, selectSQL, &devices, ns, name, model); err != nil {
		return nil, err
	}
	if len(devices) > 0 {
		return entities.ToNodeDevice(&devices[0])
	}
	return nil, common.Error(
		common.ErrResourceNotFound,
		common.Field("type", "node device"),
		common.Field("name", name))
}

func (d *BaetylCloudDB) CreateNodeDeviceTx(tx *sqlx.Tx, device *models.NodeDevice) (sql.Result, error) {
	insertSQL := `
INSERT INTO baetyl_node_device 
(name, namespace, version, device_model, access_template, node_name, driver_name, driver_inst_name, config) 
VALUES 
(?,?,?,?,?,?,?,?,?)
`
	dev, err := entities.FromNodeDevice(device)
	if err != nil {
		return nil, err
	}
	return d.Exec(tx, insertSQL, dev.Name, dev.Namespace, dev.Version, dev.DeviceModel, dev.AccessTemplate, dev.NodeName,
		dev.DriverName, dev.DriverInstName, dev.Config)
}

func (d *BaetylCloudDB) UpdateNodeDeviceTx(tx *sqlx.Tx, device *models.NodeDevice) (sql.Result, error) {
	updateSQL := `
UPDATE baetyl_node_device SET version=?, access_template=?, node_name=?, driver_name=?, driver_inst_name=?, config=?
WHERE namespace=? AND name=? AND device_model=?
`
	dev, err := entities.FromNodeDevice(device)
	if err != nil {
		return nil, err
	}
	return d.Exec(tx, updateSQL, dev.Version, dev.AccessTemplate, dev.NodeName,
		dev.DriverName, dev.DriverInstName, dev.Config, dev.Namespace, dev.Name, dev.DeviceModel)
}

func (d *BaetylCloudDB) DeleteNodeDeviceTx(tx *sqlx.Tx, ns, name, model string) (sql.Result, error) {
	deleteSQL := `
DELETE FROM baetyl_node_device where namespace=? AND name=? AND device_model=?
`
	return d.Exec(tx, deleteSQL, ns, name, model)
}

func (d *BaetylCloudDB) ListNodeDeviceByDriverAndNodeTx(
	tx *sqlx.Tx, ns, driverInstName, nodeName string, listOptions *models.ListOptions) ([]models.NodeDevice, int, error) {
	selectSQL := `
SELECT  
name, namespace, version, device_model, access_template, node_name, driver_name, driver_inst_name, config, create_time, update_time
FROM baetyl_node_device WHERE namespace=? AND driver_inst_name=? And node_name=? AND name LIKE ? ORDER BY create_time DESC
`
	args := []any{ns, driverInstName, nodeName, listOptions.GetFuzzyName()}
	var devices []entities.NodeDevice
	if listOptions.GetLimitNumber() > 0 {
		selectSQL = selectSQL + "LIMIT ?,?"
		args = append(args, listOptions.GetLimitOffset(), listOptions.GetLimitNumber())
	}

	if err := d.Query(tx, selectSQL, &devices, args...); err != nil {
		return nil, 0, err
	}
	var total []int
	if len(devices) == 0 {
		return []models.NodeDevice{}, 0, nil
	}
	numSQL := `SELECT COUNT(id)  FROM  baetyl_node_device WHERE namespace=? AND driver_inst_name=? And node_name=? AND name LIKE ?`
	if err := d.Query(tx, numSQL, &total, ns, driverInstName, nodeName, listOptions.GetFuzzyName()); err != nil {
		return nil, 0, err
	}
	var res []models.NodeDevice
	for _, device := range devices {
		dev, err := entities.ToNodeDevice(&device)
		if err != nil {
			return nil, 0, err
		}
		res = append(res, *dev)
	}
	return res, total[0], nil
}

func (d *BaetylCloudDB) ListNodeDeviceByNode(ns, nodeName string) ([]string, error) {
	selectSQL := `SELECT name FROM baetyl_node_device WHERE namespace=?  And node_name=?`

	var devicesName []string
	if err := d.Query(nil, selectSQL, &devicesName, ns, nodeName); err != nil {
		return nil, err
	}
	return devicesName, nil
}

func (d *BaetylCloudDB) BatchCreateNodeDeviceTx(tx *sqlx.Tx, devices []models.NodeDevice) error {
	var valueStrings []string
	var valueArgs []any
	for _, dev := range devices {
		dev, err := entities.FromNodeDevice(&dev)
		if err != nil {
			return err
		}
		valueStrings = append(valueStrings, "(?,?,?,?,?,?,?,?,?)")
		valueArgs = append(valueArgs, dev.Name, dev.Namespace, dev.Version, dev.DeviceModel, dev.AccessTemplate, dev.NodeName,
			dev.DriverName, dev.DriverInstName, dev.Config)
	}
	batchInsertSQL := fmt.Sprintf("INSERT INTO baetyl_node_device (name, namespace, version, device_model,"+
		" access_template, node_name, driver_name, driver_inst_name, config)"+
		" VALUES %s", strings.Join(valueStrings, ","))

	res, err := d.Exec(tx, batchInsertSQL, valueArgs...)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if affected != int64(len(devices)) {
		return common.Error(
			common.ErrUnknown,
			common.Field("type", "node device"),
			common.Field("info", "batch creation of node device data failed"))
	}
	return err
}

func (d *BaetylCloudDB) BatchCreateNodeDevice(device []models.NodeDevice) error {
	return d.BatchCreateNodeDeviceTx(nil, device)
}

func (d *BaetylCloudDB) BatchDeleteNodeDevices(ns, model string, names []string) error {
	return d.BatchDeleteNodeDevicesTx(nil, ns, model, names)
}

func (d *BaetylCloudDB) BatchDeleteNodeDevicesTx(tx *sqlx.Tx, ns, model string, names []string) error {
	deleteSQL := `DELETE FROM baetyl_node_device where namespace=? AND device_model=? AND name in (?)`
	qry, args, err := sqlx.In(deleteSQL, ns, model, names)
	if err != nil {
		return err
	}
	res, err := d.Exec(tx, qry, args...)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected != int64(len(names)) {
		return common.Error(
			common.ErrUnknown,
			common.Field("type", "node device"),
			common.Field("info", "batch delete node device data failed"))
	}
	return nil
}

func (d *BaetylCloudDB) BatchGetNodeDevicesByNames(ns, model string, names []string) (map[string]string, error) {
	selectSQL := `SELECT name,node_name FROM baetyl_node_device where namespace=? AND device_model=? AND name in (?)`
	var devices []entities.NodeDevice
	qry, args, err := sqlx.In(selectSQL, ns, model, names)
	if err != nil {
		return nil, err
	}
	err = d.Query(nil, qry, &devices, args...)
	if err != nil {
		return nil, err
	}

	res := make(map[string]string, 100)
	for _, device := range devices {
		res[device.Name] = device.NodeName
	}
	return res, err
}
