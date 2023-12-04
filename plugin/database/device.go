// Package database 数据库存储实现
package database

import (
	"database/sql"
	"fmt"
	"strings"

	dm "github.com/baetyl/baetyl-go/v2/dmcontext"
	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/baetyl/baetyl-go/v2/json"
	"github.com/baetyl/baetyl-go/v2/utils"
	"github.com/jmoiron/sqlx"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin/database/entities"
)

const KeywordAlias = "alias"

func (d *BaetylCloudDB) GetDevice(ns, _, name string) (*models.Device, error) {
	return d.GetDeviceTx(nil, ns, name)
}

func (d *BaetylCloudDB) GetDeviceByDeviceModel(ns, _, name, _ string) (*models.Device, error) {
	return d.GetDeviceTx(nil, ns, name)
}

func (d *BaetylCloudDB) ListDevice(ns, _ string, listOptions *models.ListOptions) (*models.DeviceList, error) {
	name := listOptions.GetFuzzyName()
	alias := listOptions.GetFuzzyAlias()
	selectSQL := `
SELECT  
name, namespace, version, description, protocol,
labels, alias, device_model, node_name, driver_name, attributes, properties,
shadow, ready, active, config, create_time, update_time
FROM baetyl_device WHERE namespace=? AND name LIKE ? AND alias LIKE ? ORDER BY create_time DESC 
`

	devices, count, err := d.ListDeviceByKeywordTx(nil, ns, selectSQL, name, alias, listOptions)
	if err != nil {
		return nil, err
	}
	res := &models.DeviceList{
		Items:       devices,
		Total:       count,
		ListOptions: listOptions,
	}
	return res, nil
}

func (d *BaetylCloudDB) CreateDevice(_ string, device *models.Device) (*models.Device, error) {
	var res *models.Device
	err := d.Transact(func(tx *sqlx.Tx) error {
		_, err := d.CreateDeviceTx(tx, device)
		if err != nil {
			return err
		}
		res, err = d.GetDeviceTx(tx, device.Namespace, device.Name)
		return err
	})
	return res, err
}

func (d *BaetylCloudDB) BatchCreateDevice(_ string, devices []*models.Device) ([]models.Device, error) {
	var res []models.Device
	if len(devices) == 0 {
		return res, nil
	}
	err := d.Transact(func(tx *sqlx.Tx) error {
		_, err := d.BatchCreateDeviceTx(tx, devices)
		if err != nil {
			return err
		}
		var names []string
		for i := range devices {
			names = append(names, devices[i].Name)
		}
		res, err = d.ListDeviceByNameTx(tx, devices[0].Namespace, names)
		return err
	})
	return res, err
}

func (d *BaetylCloudDB) BatchCreateDeviceTx(tx *sqlx.Tx, devices []*models.Device) (sql.Result, error) {
	insertSQL := `
INSERT INTO baetyl_device 
(name, namespace, version, description, protocol,
labels, alias, device_model, node_name, driver_name, attributes, properties, shadow, ready, active, config) 
VALUES 
`
	var args []any
	for _, device := range devices {
		insertSQL += `(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?),`
		dev, err := entities.FromModelDevice(device)
		if err != nil {
			return nil, err
		}
		args = append(args, dev.Name, dev.Namespace, dev.Version, dev.Description, dev.Protocol,
			dev.Labels, dev.Alias, dev.DeviceModel, dev.NodeName, dev.DriverName, dev.Attributes, dev.Properties, dev.Shadow, dev.Ready, dev.Active, dev.Config)
	}
	insertSQL = strings.TrimRight(insertSQL, ",")
	return d.Exec(tx, insertSQL, args...)
}

// Deprecated: BIE绑定逻辑移动到绑定中间表操作
func (d *BaetylCloudDB) BindingDevice(_, _, _, _, _ string) error {
	return nil
}

// Deprecated: BIE绑定解绑逻辑移动到绑定中间表操作
func (d *BaetylCloudDB) UnbindingDevice(_, _, _, _, _ string) error {
	return nil
}

func (d *BaetylCloudDB) UpdateDevice(_ string, device *models.Device) (*models.Device, error) {
	var res *models.Device
	err := d.Transact(func(tx *sqlx.Tx) error {
		old, err := d.GetDeviceTx(tx, device.Namespace, device.Name)
		if err != nil {
			return err
		}
		if models.EqualDevice(old, device) {
			res = old
			return nil
		}
		_, err = d.UpdateDeviceTx(tx, device)
		if err != nil {
			return err
		}
		res, err = d.GetDeviceTx(tx, device.Namespace, device.Name)
		return err
	})
	return res, err
}

func (d *BaetylCloudDB) BatchUpdateDeviceNodeAndDriver(ns, _, nodeName, driverName string, devices []string) ([]models.Device, error) {
	var res []models.Device
	err := d.Transact(func(tx *sqlx.Tx) error {
		_, err := d.BatchUpdateDeviceNodeAndDriverTx(tx, ns, nodeName, driverName, devices)
		if err != nil {
			return err
		}
		res, err = d.ListDeviceByNameTx(tx, ns, devices)
		return err
	})
	return res, err
}

func (d *BaetylCloudDB) BatchUpdateDeviceProps(ns, _ string, devices []string, properties []dm.DeviceProperty) error {
	return d.BatchUpdateDevicePropsTx(nil, ns, devices, properties)
}

func (d *BaetylCloudDB) BatchUpdateDeviceAttrs(ns, _ string, devices []string, attributes []models.DeviceAttribute) error {
	return d.BatchUpdateDeviceAttrsTx(nil, ns, devices, attributes)
}

func (d *BaetylCloudDB) DeleteDevice(ns, _, name string) error {
	_, err := d.DeleteDeviceTx(nil, ns, name)
	return err
}

func (d *BaetylCloudDB) ListDeviceByDeviceModel(ns, _, deviceModel string, listOptions *models.ListOptions) ([]models.Device, error) {
	return d.ListDeviceByDeviceModelTx(nil, ns, deviceModel, listOptions)
}

func (d *BaetylCloudDB) ListDeviceByProtocol(ns, _, protocol string, listOptions *models.ListOptions) (*models.DeviceList, error) {
	devices, err := d.ListDeviceByProtocolTx(nil, ns, protocol, listOptions)
	if err != nil {
		return nil, err
	}
	count, err := d.CountDeviceByProtocolTx(nil, ns, protocol, listOptions.GetFuzzyName())
	if err != nil {
		return nil, err
	}
	return &models.DeviceList{
		Items:       devices,
		Total:       count,
		ListOptions: listOptions,
	}, nil
}

func (d *BaetylCloudDB) ListDeviceByDriverAndNode(ns, _, driverName, nodeName string, listOptions *models.ListOptions) (*models.DeviceList, error) {
	devices, err := d.ListDeviceByDriverAndNodeTx(nil, ns, driverName, nodeName, listOptions)
	if err != nil {
		return nil, err
	}
	count, err := d.CountDeviceByDriverAndNodeTx(nil, ns, driverName, nodeName, listOptions.GetFuzzyName())
	if err != nil {
		return nil, err
	}
	return &models.DeviceList{
		Items:       devices,
		ListOptions: listOptions,
		Total:       count,
	}, nil
}

func (d *BaetylCloudDB) ListDeviceByNode(nodeName, ns string) ([]models.Device, error) {
	return d.ListDeviceByNodeTx(nil, nodeName, ns)
}

func (d *BaetylCloudDB) ListDeviceByName(ns, _ string, names []string) ([]models.Device, error) {
	return d.ListDeviceByNameTx(nil, ns, names)
}

func (d *BaetylCloudDB) GetDeviceTx(tx *sqlx.Tx, ns, name string) (*models.Device, error) {
	selectSQL := `
SELECT  
name, namespace, version, description, protocol,
labels, alias, device_model, node_name, driver_name, attributes, properties,
shadow, ready, active, config, create_time, update_time
FROM baetyl_device WHERE namespace=? AND name=? LIMIT 0,1
`
	var devices []entities.Device
	if err := d.Query(tx, selectSQL, &devices, ns, name); err != nil {
		return nil, err
	}
	if len(devices) > 0 {
		return entities.ToModelDevice(&devices[0])
	}
	return nil, common.Error(
		common.ErrResourceNotFound,
		common.Field("type", "device"),
		common.Field("name", name))
}

func (d *BaetylCloudDB) ListDeviceByKeywordTx(tx *sqlx.Tx, ns, selectSQL, name string, alias string, listOptions *models.ListOptions) (
	[]models.Device, int, error) {
	var devices []entities.Device
	args := []interface{}{ns, name, alias}
	numArgs := []interface{}{ns, name, alias}

	if listOptions.GetLimitNumber() > 0 && listOptions.LabelSelector == "" {
		selectSQL = selectSQL + "LIMIT ?,?"
		args = append(args, listOptions.GetLimitOffset(), listOptions.GetLimitNumber())
	}

	if err := d.Query(tx, selectSQL, &devices, args...); err != nil {
		return nil, 0, err
	}

	var total []int
	if len(devices) == 0 {
		return []models.Device{}, 0, nil
	}

	numSQL := `SELECT COUNT(id)  FROM baetyl_device WHERE namespace= ? AND name LIKE ? AND alias LIKE ?  `
	if err := d.Query(tx, numSQL, &total, numArgs...); err != nil {
		return nil, 0, err
	}
	var res []models.Device
	for _, device := range devices {
		labels := map[string]string{}
		if err := json.Unmarshal([]byte(device.Labels), &labels); err != nil {
			return nil, 0, errors.Trace(err)
		}
		if ok, err := utils.IsLabelMatch(listOptions.LabelSelector, labels); err != nil || !ok {
			continue
		}
		dev, err := entities.ToModelDevice(&device)
		if err != nil {
			return nil, 0, err
		}
		res = append(res, *dev)
	}
	return res, total[0], nil
}

func (d *BaetylCloudDB) CreateDeviceTx(tx *sqlx.Tx, device *models.Device) (sql.Result, error) {
	insertSQL := `
INSERT INTO baetyl_device 
(name, namespace, version, description, protocol,
labels, alias, device_model, node_name, driver_name, attributes, properties, shadow, ready, active, config) 
VALUES 
(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
`
	dev, err := entities.FromModelDevice(device)
	if err != nil {
		return nil, err
	}
	return d.Exec(tx, insertSQL, dev.Name, dev.Namespace, dev.Version, dev.Description, dev.Protocol,
		dev.Labels, dev.Alias, dev.DeviceModel, dev.NodeName, dev.DriverName, dev.Attributes, dev.Properties, dev.Shadow, dev.Ready, dev.Active, dev.Config)
}

func (d *BaetylCloudDB) BatchUpdateDevicePropsTx(tx *sqlx.Tx, ns string, names []string, properties []dm.DeviceProperty) error {
	updateSQLTemp := `
UPDATE baetyl_device SET properties=?, version = CASE name %s END WHERE namespace=? AND name IN (?)
`
	var versions string
	version := `WHEN '%s' THEN '%s' `
	for _, name := range names {
		versions += fmt.Sprintf(version, name, entities.GenResourceVersion())
	}
	updateSQL := fmt.Sprintf(updateSQLTemp, versions)
	props, err := json.Marshal(properties)
	if err != nil {
		return err
	}
	qry, args, err := sqlx.In(updateSQL, string(props), ns, names)
	if err != nil {
		return err
	}
	if _, err := d.Exec(tx, qry, args...); err != nil {
		return err
	}
	return nil
}

func (d *BaetylCloudDB) BatchUpdateDeviceAttrsTx(tx *sqlx.Tx, ns string, names []string, attributes []models.DeviceAttribute) error {
	updateSQLTemp := `
UPDATE baetyl_device SET attributes=?, version = CASE name %s END WHERE namespace=? AND name IN (?)
`
	var versions string
	version := `WHEN '%s' THEN '%s' `
	for _, name := range names {
		versions += fmt.Sprintf(version, name, entities.GenResourceVersion())
	}
	updateSQL := fmt.Sprintf(updateSQLTemp, versions)
	attrs, err := json.Marshal(attributes)
	if err != nil {
		return err
	}
	qry, args, err := sqlx.In(updateSQL, string(attrs), ns, names)
	if err != nil {
		return err
	}
	if _, err := d.Exec(tx, qry, args...); err != nil {
		return err
	}
	return nil
}

func (d *BaetylCloudDB) BatchUpdateDeviceNodeAndDriverTx(tx *sqlx.Tx, ns, nodeName, driverName string, names []string) (sql.Result, error) {
	updateSQLTemp := `
UPDATE baetyl_device SET node_name=?, driver_name=?, ready=?, version = CASE name %s END WHERE namespace=? AND name IN (?)
`
	version := `WHEN '%s' THEN '%s' `
	var versions string
	for _, name := range names {
		versions += fmt.Sprintf(version, name, entities.GenResourceVersion())
	}
	updateSQL := fmt.Sprintf(updateSQLTemp, versions)
	qry, args, err := sqlx.In(updateSQL, nodeName, driverName, false, ns, names)
	if err != nil {
		return nil, err
	}
	return d.Exec(tx, qry, args...)
}

func (d *BaetylCloudDB) ListDeviceByNameTx(tx *sqlx.Tx, ns string, names []string) ([]models.Device, error) {
	selectSQL := `SELECT  
name, namespace, version, description, protocol,
labels, alias, device_model, node_name, driver_name, attributes, properties,
shadow, ready, active, config, create_time, update_time
FROM baetyl_device WHERE namespace=? AND name IN (?)`
	qry, args, err := sqlx.In(selectSQL, ns, names)
	if err != nil {
		return nil, err
	}
	var devices []entities.Device
	if err := d.Query(tx, qry, &devices, args...); err != nil {
		return nil, err
	}
	var res []models.Device
	for _, d := range devices {
		dev, err := entities.ToModelDevice(&d)
		if err != nil {
			return nil, err
		}
		res = append(res, *dev)
	}
	return res, nil
}

func (d *BaetylCloudDB) UpdateDeviceTx(tx *sqlx.Tx, device *models.Device) (sql.Result, error) {
	updateSQL := `
UPDATE baetyl_device SET version=?, description=?, labels=?,
alias=?, device_model=?, node_name=?, driver_name=?, attributes=?, properties=?, shadow=?, ready=?, active=?, config=?
WHERE namespace=? AND name=?
`
	dev, err := entities.FromModelDevice(device)
	if err != nil {
		return nil, err
	}
	return d.Exec(tx, updateSQL, dev.Version, dev.Description, dev.Labels, dev.Alias, dev.DeviceModel, dev.NodeName,
		dev.DriverName, dev.Attributes, dev.Properties, dev.Shadow, dev.Ready, dev.Active, dev.Config, dev.Namespace, dev.Name)
}

func (d *BaetylCloudDB) DeleteDeviceTx(tx *sqlx.Tx, ns, name string) (sql.Result, error) {
	deleteSQL := `
DELETE FROM baetyl_device where namespace=? AND name=?
`
	return d.Exec(tx, deleteSQL, ns, name)
}

func (d *BaetylCloudDB) ListDeviceByDeviceModelTx(tx *sqlx.Tx, ns, deviceModel string, listOptions *models.ListOptions) ([]models.Device, error) {
	selectSQL := `
SELECT  
name, namespace, version, description, protocol,
labels, alias, device_model, node_name, driver_name, attributes, properties,
shadow, ready, active, config, create_time, update_time
FROM baetyl_device WHERE namespace=? AND device_model=? ORDER BY create_time DESC
`
	var devices []entities.Device
	args := []any{ns, deviceModel}
	if listOptions.GetLimitNumber() > 0 && listOptions.LabelSelector == "" {
		selectSQL = selectSQL + "LIMIT ?,?"
		args = append(args, listOptions.GetLimitOffset(), listOptions.GetLimitNumber())
	}
	if err := d.Query(tx, selectSQL, &devices, args...); err != nil {
		return nil, err
	}
	var res []models.Device
	for _, d := range devices {
		labels := map[string]string{}
		if err := json.Unmarshal([]byte(d.Labels), &labels); err != nil {
			return nil, errors.Trace(err)
		}
		if ok, err := utils.IsLabelMatch(listOptions.LabelSelector, labels); err != nil || !ok {
			continue
		}
		dev, err := entities.ToModelDevice(&d)
		if err != nil {
			return nil, err
		}
		res = append(res, *dev)
	}
	return res, nil
}

func (d *BaetylCloudDB) ListDeviceByProtocolTx(tx *sqlx.Tx, ns, protocol string, listOptions *models.ListOptions) ([]models.Device, error) {
	selectSQL := `
SELECT  
name, namespace, version, description, protocol,
labels, alias, device_model, node_name, driver_name, attributes, properties,
shadow, ready, active, config, create_time, update_time
FROM baetyl_device WHERE namespace=? AND protocol=? AND name LIKE ? ORDER BY create_time DESC
`
	var devices []entities.Device
	args := []any{ns, protocol, listOptions.GetFuzzyName()}
	if listOptions.GetLimitNumber() > 0 {
		selectSQL = selectSQL + "LIMIT ?,?"
		args = append(args, listOptions.GetLimitOffset(), listOptions.GetLimitNumber())
	}
	if err := d.Query(tx, selectSQL, &devices, args...); err != nil {
		return nil, err
	}
	var res []models.Device
	for _, d := range devices {
		dev, err := entities.ToModelDevice(&d)
		if err != nil {
			return nil, err
		}
		res = append(res, *dev)
	}
	return res, nil
}

func (d *BaetylCloudDB) ListDeviceByNodeTx(tx *sqlx.Tx, nodeName, ns string) ([]models.Device, error) {
	selectSQL := `
SELECT  
name, namespace, version, description, protocol,
labels, alias, device_model, node_name, driver_name, attributes, properties,
shadow, ready, active, config, create_time, update_time
FROM baetyl_device WHERE namespace=? AND node_name=? ORDER BY create_time DESC
`
	var devices []entities.Device
	args := []any{ns, nodeName}
	if err := d.Query(tx, selectSQL, &devices, args...); err != nil {
		return nil, err
	}
	var res []models.Device
	for _, d := range devices {
		dev, err := entities.ToModelDevice(&d)
		if err != nil {
			return nil, err
		}
		res = append(res, *dev)
	}
	return res, nil
}

func (d *BaetylCloudDB) CountDeviceByProtocolTx(tx *sqlx.Tx, ns, protocol, name string) (int, error) {
	selectSQL := `
SELECT count(name) AS count
FROM baetyl_device WHERE namespace=? AND protocol=? AND name LIKE ?
`
	var res []struct {
		Count int `db:"count"`
	}
	if err := d.Query(tx, selectSQL, &res, ns, protocol, name); err != nil {
		return 0, err
	}
	return res[0].Count, nil
}

func (d *BaetylCloudDB) CountDeviceByDriverAndNodeTx(tx *sqlx.Tx, ns, driverName, nodeName, name string) (int, error) {
	selectSQL := `
SELECT count(name) AS count
FROM baetyl_device WHERE namespace=? AND driver_name=? AND node_name=? AND name LIKE ?
`
	var res []struct {
		Count int `db:"count"`
	}
	if err := d.Query(tx, selectSQL, &res, ns, driverName, nodeName, name); err != nil {
		return 0, err
	}
	return res[0].Count, nil
}

func (d *BaetylCloudDB) ListDeviceByDriverAndNodeTx(tx *sqlx.Tx, ns, driverName, nodeName string, listOptions *models.ListOptions) ([]models.Device, error) {
	selectSQL := `
SELECT  
name, namespace, version, description, protocol,
labels, alias, device_model, node_name, driver_name, attributes, properties,
shadow, ready, active, config, create_time, update_time
FROM baetyl_device WHERE namespace=? AND driver_name=? And node_name=? AND name LIKE ? ORDER BY create_time DESC
`
	var devices []entities.Device
	args := []any{ns, driverName, nodeName, listOptions.GetFuzzyName()}
	if listOptions.GetLimitNumber() > 0 {
		selectSQL = selectSQL + "LIMIT ?,?"
		args = append(args, listOptions.GetLimitOffset(), listOptions.GetLimitNumber())
	}
	if err := d.Query(tx, selectSQL, &devices, args...); err != nil {
		return nil, err
	}
	var res []models.Device
	for _, d := range devices {
		dev, err := entities.ToModelDevice(&d)
		if err != nil {
			return nil, err
		}
		res = append(res, *dev)
	}
	return res, nil
}

func (d *BaetylCloudDB) UpdateDeviceStateByName(ready bool, ns string, name []string) error {
	updateSQL := ` UPDATE baetyl_device SET ready=? WHERE namespace=? AND name in (?) `
	qry, args, err := sqlx.In(updateSQL, ready, ns, name)
	if err != nil {
		return err
	}
	_, err = d.Exec(nil, qry, args...)
	if err != nil {
		return err
	}
	return nil
}

func (d *BaetylCloudDB) BatchUpdateDeviceTx(tx *sqlx.Tx, devices []models.Device) error {
	if len(devices) == 0 {
		return nil
	}
	updateSQL := `UPDATE baetyl_device SET  node_name=?, driver_name=?, ready=?, labels= CASE  %s  END, config= CASE %s END WHERE namespace=? AND name IN (?)`
	labelsCondition := make([]string, len(devices))
	configCondition := make([]string, len(devices))
	var device *entities.Device
	deviceNames := make([]string, len(devices))
	for i, de := range devices {
		dev, err := entities.FromModelDevice(&de)
		if err != nil {
			return err
		}
		labelsCondition[i] = fmt.Sprintf("WHEN name = '%s' THEN '%s'", dev.Name, dev.Labels)
		configCondition[i] = fmt.Sprintf("WHEN name = '%s' THEN '%s'", dev.Name, dev.Config)
		deviceNames[i] = dev.Name
		device = dev
	}
	labelsCn := strings.Join(labelsCondition, " ")
	configsCn := strings.Join(configCondition, " ")
	updateSQL = fmt.Sprintf(updateSQL, labelsCn, configsCn)

	qry, args, err := sqlx.In(updateSQL, device.NodeName, device.DriverName, device.Ready, device.Namespace, deviceNames)
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
	if affected != int64(len(devices)) {
		return common.Error(
			common.ErrUnknown,
			common.Field("type", "node device"),
			common.Field("info", "batch creation of node device data failed"))
	}
	return nil
}

func (d *BaetylCloudDB) BatchUpdateDevice(_ string, device []models.Device) error {
	err := d.Transact(func(tx *sqlx.Tx) error {
		err := d.BatchUpdateDeviceTx(tx, device)
		return err
	})
	return err
}

func (d *BaetylCloudDB) BatchUpdateDeviceAndDeleteNodeDevice(_ string, device []models.Device) error {
	err := d.Transact(func(tx *sqlx.Tx) error {
		if len(device) == 0 {
			return nil
		}
		err := d.BatchUpdateDeviceTx(tx, device)
		if err != nil {
			return err
		}
		devNames := make([]string, len(device))
		for i, v := range device {
			devNames[i] = v.Name
		}
		err = d.BatchDeleteNodeDevicesTx(tx, device[0].Namespace, device[0].DeviceModel, devNames)
		return err
	})
	return err
}

func (d *BaetylCloudDB) GetDevicesByDeviceModelAndNameList(ns, _ string, names map[string]string) (map[string]models.Device, error) {
	des := make(map[string]models.Device, 10)
	if len(names) == 0 {
		return des, nil
	}
	deviceNames := make([]string, len(names))
	for name := range names {
		deviceNames = append(deviceNames, name)
	}
	selectSQL := `SELECT name, ready, active FROM baetyl_device WHERE namespace=? AND name IN (?)`
	var devices []entities.Device
	qry, args, err := sqlx.In(selectSQL, ns, deviceNames)
	if err != nil {
		return nil, err
	}
	err = d.Query(nil, qry, &devices, args...)
	if err != nil {
		return nil, err
	}
	for _, v := range devices {
		des[v.Name] = models.Device{
			Name:   v.Name,
			Ready:  v.Ready,
			Active: v.Active,
		}
	}
	return des, nil
}

func (d *BaetylCloudDB) GetDeviceNumByDeviceModel(ns, deviceModel string) (int, error) {
	selectSQL := `SELECT COUNT(ID) FROM baetyl_device WHERE namespace=? AND device_model= ?`
	var total []int
	if err := d.Query(nil, selectSQL, &total, ns, deviceModel); err != nil {
		return 0, err
	}
	if len(total) > 0 {
		return total[0], nil
	}
	return 0, nil
}

func (d *BaetylCloudDB) BatchGetDeviceByNames(ns, _ string, deviceName []string) (map[string]string, error) {
	deviceModelMap := map[string]string{}
	if len(deviceName) == 0 {
		return deviceModelMap, nil
	}
	selectSQL := `SELECT * FROM baetyl_device WHERE namespace=? AND name in (?)`
	qry, args, err := sqlx.In(selectSQL, ns, deviceName)
	if err != nil {
		return nil, errors.Trace(err)
	}
	var devices []entities.Device
	if err := d.Query(nil, qry, &devices, args...); err != nil {
		return nil, errors.Trace(err)
	}
	for _, device := range devices {
		deviceModelMap[device.Name] = device.DeviceModel
	}
	return deviceModelMap, nil
}
