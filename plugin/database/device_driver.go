// Package database 数据库存储实现
package database

import (
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin/database/entities"
)

func (d *BaetylCloudDB) GetDeviceDriver(tx any, ns, nodeName, driverInstName string) (*models.DeviceDriver, error) {
	transaction, err := d.InterfaceToTx(tx)
	if err != nil {
		return nil, err
	}
	return d.GetDeviceDriverTx(transaction, ns, nodeName, driverInstName)
}

func (d *BaetylCloudDB) ListDeviceDriver(tx any, ns string, nodeName string, listOptions *models.ListOptions) (*models.DeviceDriverList, error) {
	transaction, err := d.InterfaceToTx(tx)
	if err != nil {
		return nil, err
	}

	deviceDrivers, err := d.ListDeviceDriverTx(transaction, ns, nodeName, listOptions)
	if err != nil {
		return nil, err
	}
	count, err := d.CountDeviceDriverTx(transaction, ns, nodeName, listOptions.GetFuzzyName())
	if err != nil {
		return nil, err
	}
	return &models.DeviceDriverList{
		Items:       deviceDrivers,
		Total:       count,
		ListOptions: listOptions,
	}, nil
}

func (d *BaetylCloudDB) CreateDeviceDriver(tx any, deviceDriver *models.DeviceDriver) (*models.DeviceDriver, error) {
	var res *models.DeviceDriver

	var err error
	if tx == nil {
		err = d.Transact(func(tx *sqlx.Tx) error {
			_, exErr := d.CreateDeviceDriverTx(tx, deviceDriver)
			if exErr != nil {
				return exErr
			}
			res, exErr = d.GetDeviceDriverTx(tx, deviceDriver.Namespace, deviceDriver.NodeName, deviceDriver.DriverInstName)
			return exErr
		})
	} else {
		transaction, txErr := d.InterfaceToTx(tx)
		if txErr != nil {
			return nil, txErr
		}
		_, err = d.CreateDeviceDriverTx(transaction, deviceDriver)
		if err != nil {
			return nil, err
		}
		res, err = d.GetDeviceDriverTx(transaction, deviceDriver.Namespace, deviceDriver.NodeName, deviceDriver.DriverInstName)
		if err != nil {
			return nil, err
		}
	}
	return res, err
}

func (d *BaetylCloudDB) UpdateDeviceDriver(tx any, deviceDriver *models.DeviceDriver) (*models.DeviceDriver, error) {
	var res *models.DeviceDriver

	var err error
	if tx == nil {
		err = d.Transact(func(tx *sqlx.Tx) error {
			old, exErr := d.GetDeviceDriverTx(tx, deviceDriver.Namespace, deviceDriver.NodeName, deviceDriver.DriverInstName)
			if exErr != nil {
				return exErr
			}
			if models.EqualDeviceDriver(old, deviceDriver) {
				res = old
				return nil
			}
			_, exErr = d.UpdateDeviceDriverTx(tx, deviceDriver)
			if exErr != nil {
				return exErr
			}
			res, exErr = d.GetDeviceDriverTx(tx, deviceDriver.Namespace, deviceDriver.NodeName, deviceDriver.DriverInstName)
			return exErr
		})
	} else {
		transaction, txErr := d.InterfaceToTx(tx)
		if txErr != nil {
			return nil, txErr
		}
		old, exErr := d.GetDeviceDriverTx(transaction, deviceDriver.Namespace, deviceDriver.NodeName, deviceDriver.DriverInstName)
		if exErr != nil {
			return nil, exErr
		}
		if models.EqualDeviceDriver(old, deviceDriver) {
			res = old
			return res, nil
		}
		_, exErr = d.UpdateDeviceDriverTx(transaction, deviceDriver)
		if exErr != nil {
			return nil, exErr
		}
		res, exErr = d.GetDeviceDriverTx(transaction, deviceDriver.Namespace, deviceDriver.NodeName, deviceDriver.DriverInstName)
		return nil, exErr
	}
	return res, err
}

func (d *BaetylCloudDB) DeleteDeviceDriver(tx any, ns, nodeName, driverInstName string) error {
	transaction, err := d.InterfaceToTx(tx)
	if err != nil {
		return err
	}
	_, err = d.DeleteDeviceDriverTx(transaction, ns, nodeName, driverInstName)
	return err
}

func (d *BaetylCloudDB) ListDeviceDriverByName(tx any, ns, driverName string) ([]models.DeviceDriver, error) {
	transaction, err := d.InterfaceToTx(tx)
	if err != nil {
		return nil, err
	}
	return d.ListDeviceDriverByNameTx(transaction, ns, driverName)
}

func (d *BaetylCloudDB) GetDeviceDriverTx(tx *sqlx.Tx, ns, nodeName, driverInstName string) (*models.DeviceDriver, error) {
	selectSQL := `
SELECT  
node_name, driver_name, driver_inst_name, namespace, version, protocol, application,
configuration, driver_config, create_time, update_time 
FROM baetyl_device_driver WHERE namespace=? AND node_name=? AND driver_inst_name=? LIMIT 0,1
`
	var drivers []entities.DeviceDriver
	if err := d.Query(tx, selectSQL, &drivers, ns, nodeName, driverInstName); err != nil {
		return nil, err
	}
	if len(drivers) > 0 {
		return entities.ToModelDeviceDriver(&drivers[0])
	}
	return nil, common.Error(
		common.ErrResourceNotFound,
		common.Field("type", "device driver, driver inst name ["+driverInstName+"]"),
		common.Field("name", nodeName))
}

func (d *BaetylCloudDB) ListDeviceDriverTx(tx *sqlx.Tx, ns, nodeName string, listOptions *models.ListOptions) ([]models.DeviceDriver, error) {
	selectSQL := `
SELECT  
node_name, driver_name, driver_inst_name, namespace, version, protocol, application,
configuration, driver_config, create_time, update_time
FROM baetyl_device_driver WHERE namespace=? AND node_name=? AND driver_inst_name LIKE ? ORDER BY create_time DESC
`
	var drivers []entities.DeviceDriver
	args := []interface{}{ns, nodeName, listOptions.GetFuzzyName()}
	if listOptions.GetLimitNumber() > 0 {
		selectSQL = selectSQL + "LIMIT ?,?"
		args = append(args, listOptions.GetLimitOffset(), listOptions.GetLimitNumber())
	}
	if err := d.Query(tx, selectSQL, &drivers, args...); err != nil {
		return nil, err
	}
	var res []models.DeviceDriver
	for _, d := range drivers {
		driver, err := entities.ToModelDeviceDriver(&d)
		if err != nil {
			return nil, err
		}
		res = append(res, *driver)
	}
	return res, nil
}

func (d *BaetylCloudDB) ListDeviceDriverByNameTx(tx *sqlx.Tx, ns, driverName string) ([]models.DeviceDriver, error) {
	selectSQL := `
SELECT  
node_name, driver_name, driver_inst_name, namespace, version, protocol, application,
configuration, driver_config, create_time, update_time
FROM baetyl_device_driver WHERE namespace=? AND driver_name=? ORDER BY create_time DESC
`
	var drivers []entities.DeviceDriver
	args := []interface{}{ns, driverName}
	if err := d.Query(tx, selectSQL, &drivers, args...); err != nil {
		return nil, err
	}
	var res []models.DeviceDriver
	for _, dr := range drivers {
		driver, err := entities.ToModelDeviceDriver(&dr)
		if err != nil {
			return nil, err
		}
		res = append(res, *driver)
	}
	return res, nil
}

func (d *BaetylCloudDB) CreateDeviceDriverTx(tx *sqlx.Tx, driver *models.DeviceDriver) (sql.Result, error) {
	insertSQL := `
INSERT INTO baetyl_device_driver
(node_name, driver_name, driver_inst_name, namespace, version, protocol, 
application, configuration, driver_config)
VALUES 
(?,?,?,?,?,?,?,?,?)
`
	dv, err := entities.FromModelDeviceDriver(driver)
	if err != nil {
		return nil, err
	}
	return d.Exec(tx, insertSQL, dv.NodeName, dv.DriverName, dv.DriverInstName, dv.Namespace, dv.Version,
		dv.Protocol, dv.Application, dv.Configuration, dv.DriverConfig)
}

func (d *BaetylCloudDB) UpdateDeviceDriverTx(tx *sqlx.Tx, deviceDriver *models.DeviceDriver) (sql.Result, error) {
	updateSQL := `
UPDATE baetyl_device_driver SET version=?, application=?, configuration=?,
driver_config=? WHERE namespace=? AND node_name=? AND driver_inst_name=?
`
	dDriver, err := entities.FromModelDeviceDriver(deviceDriver)
	if err != nil {
		return nil, err
	}
	return d.Exec(tx, updateSQL, dDriver.Version, dDriver.Application, dDriver.Configuration,
		dDriver.DriverConfig, dDriver.Namespace, dDriver.NodeName, dDriver.DriverInstName)
}

func (d *BaetylCloudDB) DeleteDeviceDriverTx(tx *sqlx.Tx, ns, nodeName, driverInstName string) (sql.Result, error) {
	deleteSQL := `
DELETE FROM baetyl_device_driver where namespace=? AND node_name=? AND driver_inst_name=?
`
	return d.Exec(tx, deleteSQL, ns, nodeName, driverInstName)
}

func (d *BaetylCloudDB) CountDeviceDriverTx(tx *sqlx.Tx, ns, nodeName, driverInstName string) (int, error) {
	selectSQL := `
SELECT count(driver_inst_name) AS count
FROM baetyl_device_driver WHERE namespace=? AND node_name=? AND driver_inst_name LIKE ?
`
	var res []struct {
		Count int `db:"count"`
	}
	if err := d.Query(tx, selectSQL, &res, ns, nodeName, driverInstName); err != nil {
		return 0, err
	}
	return res[0].Count, nil
}
