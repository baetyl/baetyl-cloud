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

// GetDMPLinkInfo 这里做空实现即可
func (d *BaetylCloudDB) GetDMPLinkInfo(_, _, _ string) (*models.DMPLinkInfo, error) {
	return &models.DMPLinkInfo{}, nil
}

// ListAllInstance 这里做空实现即可
func (d *BaetylCloudDB) ListAllInstance(_, _ string) (*models.ListView, error) {
	return &models.ListView{Total: 0, Items: []models.DMPInstance{}}, nil
}

func (d *BaetylCloudDB) GetDeviceModel(ns, _, name string) (*models.DeviceModel, error) {
	return d.GetDeviceModelTx(nil, ns, name, common.TypeDeviceModel)
}

func (d *BaetylCloudDB) GetDeviceModelDraft(ns, _, name string) (*models.DeviceModel, error) {
	return d.GetDeviceModelTx(nil, ns, name, common.TypeDeviceModelDraft)
}

func (d *BaetylCloudDB) ListDeviceModel(ns, _ string, listOptions *models.ListOptions) (*models.DeviceModelList, error) {
	deviceModels, count, err := d.ListDeviceModelTx(nil, ns, listOptions)
	if err != nil {
		return nil, err
	}
	res := &models.DeviceModelList{
		Items:       deviceModels,
		Total:       count,
		ListOptions: listOptions,
	}
	return res, nil
}

func (d *BaetylCloudDB) ListDeviceModelByNames(ns, _ string, names []string) ([]models.DeviceModel, error) {
	return d.ListDeviceModelByNameTx(nil, ns, names)
}

func (d *BaetylCloudDB) CreateDeviceModel(_ string, deviceModel *models.DeviceModel) (*models.DeviceModel, error) {
	var res *models.DeviceModel
	err := d.Transact(func(tx *sqlx.Tx) error {
		_, err := d.CreateDeviceModelTx(tx, deviceModel)
		if err != nil {
			return err
		}
		res, err = d.GetDeviceModelTx(tx, deviceModel.Namespace, deviceModel.Name, deviceModel.Type)
		return err
	})
	return res, err
}

func (d *BaetylCloudDB) UpdateDeviceModel(_ string, deviceModel *models.DeviceModel) (*models.DeviceModel, error) {
	var res *models.DeviceModel
	err := d.Transact(func(tx *sqlx.Tx) error {
		old, err := d.GetDeviceModelTx(tx, deviceModel.Namespace, deviceModel.Name, deviceModel.Type)
		if err != nil {
			return err
		}
		if models.EqualDeviceModel(old, deviceModel) {
			res = old
			return nil
		}
		_, err = d.UpdateDeviceModelTx(tx, deviceModel)
		if err != nil {
			return err
		}
		res, err = d.GetDeviceModelTx(tx, deviceModel.Namespace, deviceModel.Name, deviceModel.Type)
		return err
	})
	return res, err
}

func (d *BaetylCloudDB) DeleteDeviceModel(ns, _, name string) error {
	_, err := d.DeleteDeviceModelTx(nil, ns, name, common.TypeDeviceModel)
	return err
}

func (d *BaetylCloudDB) DeleteDeviceModelDraft(ns, _, name string) error {
	_, err := d.DeleteDeviceModelTx(nil, ns, name, common.TypeDeviceModelDraft)
	return err
}

func (d *BaetylCloudDB) GetDeviceModelTx(tx *sqlx.Tx, ns, name string, typ byte) (*models.DeviceModel, error) {
	selectSQL := `
SELECT  
name, namespace, version, description, protocol, labels, type,
attributes, properties, create_time, update_time
FROM baetyl_device_model WHERE namespace=? AND name=? AND type=? LIMIT 0,1
`
	var deviceModels []entities.DeviceModel
	if err := d.Query(tx, selectSQL, &deviceModels, ns, name, typ); err != nil {
		return nil, err
	}
	if len(deviceModels) > 0 {
		return entities.ToModelDeviceModel(&deviceModels[0])
	}
	return nil, common.Error(
		common.ErrResourceNotFound,
		common.Field("type", "device model"),
		common.Field("name", name))
}

func (d *BaetylCloudDB) ListDeviceModelTx(tx *sqlx.Tx, ns string, listOptions *models.ListOptions) ([]models.DeviceModel, int, error) {
	selectSQL := `
SELECT  
name, namespace, version, description, protocol, labels, type,
attributes, properties, create_time, update_time
FROM baetyl_device_model WHERE namespace=? AND type=? AND name LIKE ? ORDER BY create_time DESC 
`
	var deviceModels []entities.DeviceModel
	args := []interface{}{ns, common.TypeDeviceModel, listOptions.GetFuzzyName()}
	if err := d.Query(tx, selectSQL, &deviceModels, args...); err != nil {
		return nil, 0, errors.Trace(err)
	}
	var res []models.DeviceModel
	for _, dm := range deviceModels {
		labels := map[string]string{}
		if err := json.Unmarshal([]byte(dm.Labels), &labels); err != nil {
			return nil, 0, errors.Trace(err)
		}
		if ok, err := utils.IsLabelMatch(listOptions.LabelSelector, labels); err != nil || !ok {
			continue
		}
		devModel, err := entities.ToModelDeviceModel(&dm)
		if err != nil {
			return nil, 0, errors.Trace(err)
		}
		res = append(res, *devModel)
	}
	start, end := models.GetPagingParam(listOptions, len(res))
	return res[start:end], len(res), nil
}

func (d *BaetylCloudDB) ListDeviceModelByNameTx(tx *sqlx.Tx, ns string, names []string) ([]models.DeviceModel, error) {
	selectSQL := `
SELECT  
name, namespace, version, description, protocol, labels, type,
attributes, properties, create_time, update_time
FROM baetyl_device_model WHERE namespace=? AND type=? AND name IN (?)`
	qry, args, err := sqlx.In(selectSQL, ns, common.TypeDeviceModel, names)
	if err != nil {
		return nil, err
	}
	var deviceModels []entities.DeviceModel
	if err := d.Query(tx, qry, &deviceModels, args...); err != nil {
		return nil, err
	}
	var res []models.DeviceModel
	for _, dm := range deviceModels {
		devModel, err := entities.ToModelDeviceModel(&dm)
		if err != nil {
			return nil, err
		}
		res = append(res, *devModel)
	}
	return res, nil
}

func (d *BaetylCloudDB) CreateDeviceModelTx(tx *sqlx.Tx, deviceModel *models.DeviceModel) (sql.Result, error) {
	insertSQL := `
INSERT INTO baetyl_device_model
(name, namespace, version, description, protocol, labels, type,
attributes, properties)
VALUES 
(?,?,?,?,?,?,?,?,?)
`
	devModel, err := entities.FromModelDeviceModel(deviceModel)
	if err != nil {
		return nil, err
	}
	return d.Exec(tx, insertSQL, devModel.Name, devModel.Namespace, devModel.Version, devModel.Description,
		devModel.Protocol, devModel.Labels, devModel.Type, devModel.Attributes, devModel.Properties)
}

func (d *BaetylCloudDB) UpdateDeviceModelTx(tx *sqlx.Tx, deviceModel *models.DeviceModel) (sql.Result, error) {
	updateSQL := `
UPDATE baetyl_device_model SET version=?, description=?, labels=?,
attributes=?, properties=? WHERE namespace=? AND name=? AND type=?
`
	devModel, err := entities.FromModelDeviceModel(deviceModel)
	if err != nil {
		return nil, err
	}
	return d.Exec(tx, updateSQL, devModel.Version, devModel.Description, devModel.Labels,
		devModel.Attributes, devModel.Properties, devModel.Namespace, devModel.Name, devModel.Type)
}

func (d *BaetylCloudDB) DeleteDeviceModelTx(tx *sqlx.Tx, ns, name string, typ int) (sql.Result, error) {
	deleteSQL := `
DELETE FROM baetyl_device_model where namespace=? AND name=? AND type=?
`
	return d.Exec(tx, deleteSQL, ns, name, typ)
}

func (d *BaetylCloudDB) CountDeviceModelTx(tx *sqlx.Tx, ns, name string) (int, error) {
	selectSQL := `
SELECT count(name) AS count
FROM baetyl_device_model WHERE namespace=? AND type=? AND name LIKE ?
`
	var res []struct {
		Count int `db:"count"`
	}
	if err := d.Query(tx, selectSQL, &res, ns, common.TypeDeviceModel, name); err != nil {
		return 0, err
	}
	return res[0].Count, nil
}
