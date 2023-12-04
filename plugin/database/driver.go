// Package database 数据库存储实现
package database

import (
	"database/sql"

	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/baetyl/baetyl-go/v2/json"
	"github.com/baetyl/baetyl-go/v2/utils"
	"github.com/jmoiron/sqlx"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin/database/entities"
)

func (d *BaetylCloudDB) GetDriver(ns, name string) (*models.Driver, error) {
	return d.GetDriverTx(nil, ns, name)
}

func (d *BaetylCloudDB) ListDriver(ns string, filter *models.ListOptions) (*models.DriverList, error) {
	drivers, resLen, err := d.ListDriverTx(nil, ns, filter)
	if err != nil {
		return nil, err
	}

	return &models.DriverList{
		Items:       drivers,
		Total:       resLen,
		ListOptions: filter,
	}, nil
}

func (d *BaetylCloudDB) CreateDriver(driver *models.Driver) (*models.Driver, error) {
	var res *models.Driver
	err := d.Transact(func(tx *sqlx.Tx) error {
		_, err := d.CreateDriverTx(tx, driver)
		if err != nil {
			return err
		}
		res, err = d.GetDriverTx(tx, driver.Namespace, driver.Name)
		return err
	})
	return res, err
}

func (d *BaetylCloudDB) UpdateDriver(driver *models.Driver) (*models.Driver, error) {
	var res *models.Driver
	err := d.Transact(func(tx *sqlx.Tx) error {
		old, err := d.GetDriverTx(tx, driver.Namespace, driver.Name)
		if err != nil {
			return err
		}
		if models.EqualDriver(old, driver) {
			res = old
			return nil
		}
		_, err = d.UpdateDriverTx(tx, driver)
		if err != nil {
			return err
		}
		res, err = d.GetDriverTx(tx, driver.Namespace, driver.Name)
		return err
	})
	return res, err
}

func (d *BaetylCloudDB) DeleteDriver(ns, name string) error {
	_, err := d.DeleteDriverTx(nil, ns, name)
	return err
}

func (d *BaetylCloudDB) GetDriverTx(tx *sqlx.Tx, ns, name string) (*models.Driver, error) {
	selectSQL := `
SELECT  
name, namespace, version, type, mode, labels, protocol, arch, description, 
default_config, service, volumes, registries, program_config, create_time, update_time
FROM baetyl_driver WHERE namespace=? AND name=? LIMIT 0,1
`
	var drivers []entities.Driver
	if err := d.Query(tx, selectSQL, &drivers, ns, name); err != nil {
		return nil, err
	}
	if len(drivers) > 0 {
		return entities.ToModelDriver(&drivers[0])
	}
	return nil, common.Error(
		common.ErrResourceNotFound,
		common.Field("type", "driver"),
		common.Field("name", name))
}

func (d *BaetylCloudDB) ListDriverTx(tx *sqlx.Tx, ns string, filter *models.ListOptions) ([]models.Driver, int, error) {
	selectSQL := `
SELECT  
name, namespace, version, type, mode, labels, protocol, arch, description, 
default_config, service, volumes, registries, program_config, create_time, update_time
FROM baetyl_driver WHERE namespace=? AND name LIKE ? ORDER BY create_time DESC 
`
	var drivers []entities.Driver
	if err := d.Query(tx, selectSQL, &drivers, ns, filter.GetFuzzyName()); err != nil {
		return nil, 0, err
	}
	var res []models.Driver
	for _, d := range drivers {
		labels := map[string]string{}
		if err := json.Unmarshal([]byte(d.Labels), &labels); err != nil {
			return nil, 0, errors.Trace(err)
		}
		if ok, err := utils.IsLabelMatch(filter.LabelSelector, labels); err != nil || !ok {
			continue
		}
		driver, err := entities.ToModelDriver(&d)
		if err != nil {
			return nil, 0, err
		}
		res = append(res, *driver)
	}
	start, end := models.GetPagingParam(filter, len(res))
	return res[start:end], len(res), nil
}

func (d *BaetylCloudDB) CreateDriverTx(tx *sqlx.Tx, driver *models.Driver) (sql.Result, error) {
	insertSQL := `
INSERT INTO baetyl_driver
(name, namespace, version, type, mode, labels, protocol, arch, description, 
default_config, service, volumes, registries, program_config) 
VALUES 
(?,?,?,?,?,?,?,?,?,?,?,?,?,?)
`
	dv, err := entities.FromModelDriver(driver)
	if err != nil {
		return nil, err
	}
	return d.Exec(tx, insertSQL, dv.Name, dv.Namespace, dv.Version, dv.Type, dv.Mode, dv.Labels,
		dv.Protocol, dv.Architecture, dv.Description, dv.DefaultConfig, dv.Service,
		dv.Volumes, dv.Registries, dv.ProgramConfig)
}

func (d *BaetylCloudDB) UpdateDriverTx(tx *sqlx.Tx, driver *models.Driver) (sql.Result, error) {
	updateSQL := `
UPDATE baetyl_driver SET version=?, type=?, arch=?, labels=?,
description=?, default_config=?, service=?, volumes=?, registries=?, program_config=?
WHERE namespace=? AND name=?
`
	dv, err := entities.FromModelDriver(driver)
	if err != nil {
		return nil, err
	}
	return d.Exec(tx, updateSQL, dv.Version, dv.Type, dv.Architecture, dv.Labels, dv.Description,
		dv.DefaultConfig, dv.Service, dv.Volumes, dv.Registries, dv.ProgramConfig,
		dv.Namespace, dv.Name)
}

func (d *BaetylCloudDB) DeleteDriverTx(tx *sqlx.Tx, ns, name string) (sql.Result, error) {
	deleteSQL := `
DELETE FROM baetyl_driver where namespace=? AND name=?
`
	return d.Exec(tx, deleteSQL, ns, name)
}

func (d *BaetylCloudDB) CountDriverTx(tx *sqlx.Tx, ns, name string) (int, error) {
	selectSQL := `
SELECT count(name) AS count
FROM baetyl_driver WHERE namespace=? AND name LIKE ?
`
	var res []struct {
		Count int `db:"count"`
	}
	if err := d.Query(tx, selectSQL, &res, ns, name); err != nil {
		return 0, err
	}
	return res[0].Count, nil
}
