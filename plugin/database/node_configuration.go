package database

import (
	"database/sql"

	"github.com/baetyl/baetyl-go/v2/utils"
	"github.com/jinzhu/copier"
	"github.com/jmoiron/sqlx"

	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin/database/entities"
)

func (d *BaetylCloudDB) GetNodeConfig(tx interface{}, ns, nodeName, configurationType string) (*models.NodeConfiguration, error) {
	defer utils.Trace(d.Log.Debug, "GetNodeConfig")()
	transaction, err := d.InterfaceToTx(tx)
	if err != nil {
		return nil, err
	}
	return d.GetNodeConfigTx(transaction, ns, nodeName, configurationType)
}

func (d *BaetylCloudDB) CreateNodeConfig(tx interface{}, nodeConfig *models.NodeConfiguration) (*models.NodeConfiguration, error) {
	defer utils.Trace(d.Log.Debug, "CreateNodeConfig")()
	transaction, err := d.InterfaceToTx(tx)
	if err != nil {
		return nil, err
	}
	_, err = d.CreateDeviceCacheTx(transaction, nodeConfig)
	if err != nil {
		return nil, err
	}
	return d.GetNodeConfigTx(transaction, nodeConfig.Namespace, nodeConfig.NodeName, nodeConfig.ConfigurationType)
}

func (d *BaetylCloudDB) UpdateNodeConfig(tx interface{}, nodeConfig *models.NodeConfiguration) (*models.NodeConfiguration, error) {
	defer utils.Trace(d.Log.Debug, "UpdateNodeConfig")()
	transaction, err := d.InterfaceToTx(tx)
	if err != nil {
		return nil, err
	}
	_, err = d.UpdateNodeConfigTx(transaction, nodeConfig)
	if err != nil {
		return nil, err
	}
	return d.GetNodeConfigTx(transaction, nodeConfig.Namespace, nodeConfig.NodeName, nodeConfig.ConfigurationType)
}

func (d *BaetylCloudDB) GetNodeConfigTx(tx *sqlx.Tx, ns, nodeName, configurationType string) (*models.NodeConfiguration, error) {
	selectSQL := `
SELECT  
node_name, namespace, create_time, update_time,data, type  FROM baetyl_node_configuration WHERE node_name=? AND namespace=? AND type=?`
	var nodeConfig []entities.NodeConfiguration
	if err := d.Query(tx, selectSQL, &nodeConfig, nodeName, ns, configurationType); err != nil {
		return nil, err
	}
	if len(nodeConfig) > 0 {
		nc := new(models.NodeConfiguration)
		err := copier.Copy(&nc, &nodeConfig[0])
		if err != nil {
			return nil, err
		}
		return nc, nil
	}
	return nil, nil
}

func (d *BaetylCloudDB) CreateDeviceCacheTx(tx *sqlx.Tx, nodeConfig *models.NodeConfiguration) (sql.Result, error) {
	insertSQL := `
INSERT INTO baetyl_node_configuration
(node_name, namespace, type, data) VALUES (?,?,?,?)`
	return d.Exec(tx, insertSQL, nodeConfig.NodeName, nodeConfig.Namespace, nodeConfig.ConfigurationType, nodeConfig.Data)
}

func (d *BaetylCloudDB) UpdateNodeConfigTx(tx *sqlx.Tx, nodeConfig *models.NodeConfiguration) (sql.Result, error) {
	updateSQL := `UPDATE baetyl_node_configuration SET data=? WHERE  namespace=? AND node_name=? AND type=?`
	return d.Exec(tx, updateSQL, nodeConfig.Data, nodeConfig.Namespace, nodeConfig.NodeName, nodeConfig.ConfigurationType)
}
