package database

import (
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin/database/entities"
)

func (d *BaetylCloudDB) ListDeviceUplink(ns, nodeName string, listOptions *models.ListOptions) ([]entities.DeviceUplink, error) {
	return d.ListDeviceUplinkTx(nil, ns, nodeName, listOptions)
}

func (d *BaetylCloudDB) GetDeviceUplink(ns, nodeName, destinationName string) (*entities.DeviceUplink, error) {
	return d.GetDeviceUplinkTx(nil, ns, nodeName, destinationName)
}

func (d *BaetylCloudDB) CreateDeviceUplink(deviceUplink *entities.DeviceUplink) (*entities.DeviceUplink, error) {
	_, err := d.CreateDeviceUplinkTx(nil, deviceUplink)
	if err != nil {
		return nil, err
	}
	return d.GetDeviceUplinkTx(nil, deviceUplink.Namespace, deviceUplink.NodeName, deviceUplink.DestinationName)
}

func (d *BaetylCloudDB) UpdateDeviceUplink(deviceUplink *entities.DeviceUplink) (*entities.DeviceUplink, error) {
	_, err := d.UpdateDeviceUplinkTx(nil, deviceUplink)
	if err != nil {
		return nil, err
	}
	return d.GetDeviceUplinkTx(nil, deviceUplink.Namespace, deviceUplink.NodeName, deviceUplink.DestinationName)

}

func (d *BaetylCloudDB) DeleteDeviceUplink(ns, nodeName, destinationName string) error {
	_, err := d.DeleteDeviceUplinkTx(nil, ns, nodeName, destinationName)
	return err
}

func (d *BaetylCloudDB) ListDeviceUplinkTx(tx *sqlx.Tx, ns, nodeName string, listOptions *models.ListOptions) ([]entities.DeviceUplink, error) {
	selectSQL := `
SELECT  
node_name, namespace, protocol, destination, destination_name, address,
mqtt_user, mqtt_password, http_method, http_path,create_time,update_time,private_key,ca,cert,passphrase
FROM baetyl_device_uplink WHERE namespace=? AND node_name=?  ORDER BY create_time DESC
`
	var DeviceUplink []entities.DeviceUplink
	args := []any{ns, nodeName}
	if listOptions.GetLimitNumber() > 0 {
		selectSQL = selectSQL + " LIMIT ?,?"
		args = append(args, listOptions.GetLimitOffset(), listOptions.GetLimitNumber())
	}
	if err := d.Query(tx, selectSQL, &DeviceUplink, args...); err != nil {
		return nil, err
	}
	return DeviceUplink, nil

}

func (d *BaetylCloudDB) GetDeviceUplinkCount(ns, nodeName string, listOptions *models.ListOptions) (int, error) {
	return d.GetDeviceUplinkCountTx(nil, ns, nodeName, listOptions)
}

func (d *BaetylCloudDB) GetDeviceUplinkCountTx(tx *sqlx.Tx, ns, nodeName string, _ *models.ListOptions) (int, error) {
	selectSQL := `SELECT  COUNT(id) FROM baetyl_device_uplink  WHERE namespace=? AND node_name=?`
	var count []int
	args := []any{ns, nodeName}
	if err := d.Query(tx, selectSQL, &count, args...); err != nil {
		return 0, err
	}
	return count[0], nil
}

func (d *BaetylCloudDB) GetDeviceUplinkTx(tx *sqlx.Tx, ns, nodeName, destinationName string) (*entities.DeviceUplink, error) {
	selectSQL := `
SELECT  
node_name, namespace, protocol, destination, destination_name, address,
mqtt_user, mqtt_password, http_method, http_path,create_time,update_time,ca,private_key,cert,passphrase
FROM baetyl_device_uplink WHERE namespace=? AND node_name=? AND destination_name =? ORDER BY create_time DESC
`
	var DeviceUplink []entities.DeviceUplink
	args := []any{ns, nodeName, destinationName}
	if err := d.Query(tx, selectSQL, &DeviceUplink, args...); err != nil {
		return nil, err
	}
	if len(DeviceUplink) > 0 {
		return &DeviceUplink[0], nil
	}
	return nil, common.Error(
		common.ErrResourceNotFound,
		common.Field("type", "device link"),
		common.Field("name", nodeName))
}

func (d *BaetylCloudDB) CreateDeviceUplinkTx(tx *sqlx.Tx, deviceUplink *entities.DeviceUplink) (sql.Result, error) {
	insertSQL := `
INSERT INTO baetyl_device_uplink
(node_name, namespace, protocol, destination, destination_name, address,
mqtt_user, mqtt_password, http_method, http_path,ca,private_key,cert,passphrase)
VALUES 
(?,?,?,?,?,?,?,?,?,?,?,?,?,?)
`
	return d.Exec(tx, insertSQL, deviceUplink.NodeName, deviceUplink.Namespace, deviceUplink.Protocol, deviceUplink.Destination,
		deviceUplink.DestinationName, deviceUplink.Address, deviceUplink.MQTTUser,
		deviceUplink.MQTTPassword, deviceUplink.HTTPMethod, deviceUplink.HTTPPath, deviceUplink.CA, deviceUplink.PrivateKey, deviceUplink.Cert,
		deviceUplink.Passphrase)
}

func (d *BaetylCloudDB) UpdateDeviceUplinkTx(tx *sqlx.Tx, deviceUplink *entities.DeviceUplink) (sql.Result, error) {
	updateSQL := `
UPDATE baetyl_device_uplink
SET     address=?,mqtt_user=?, mqtt_password=?, http_method=?, http_path=? ,ca = ? , private_key=? ,
        cert=? , passphrase =? WHERE  namespace=? AND node_name=? AND destination_name=? 
`
	return d.Exec(tx, updateSQL, deviceUplink.Address, deviceUplink.MQTTUser,
		deviceUplink.MQTTPassword, deviceUplink.HTTPMethod, deviceUplink.HTTPPath, deviceUplink.CA, deviceUplink.PrivateKey, deviceUplink.Cert,
		deviceUplink.Passphrase, deviceUplink.Namespace, deviceUplink.NodeName, deviceUplink.DestinationName)
}

func (d *BaetylCloudDB) DeleteDeviceUplinkTx(tx *sqlx.Tx, ns, nodeName, destinationName string) (sql.Result, error) {
	deleteSQL := ` DELETE FROM baetyl_device_uplink WHERE namespace=? AND node_name=? AND destination_name=?`
	return d.Exec(tx, deleteSQL, ns, nodeName, destinationName)
}
