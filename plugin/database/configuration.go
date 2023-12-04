// Package database 数据库存储实现
package database

import (
	"database/sql"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin/database/entities"

	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/baetyl/baetyl-go/v2/json"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/baetyl/baetyl-go/v2/utils"
	"github.com/jmoiron/sqlx"
)

func (d *BaetylCloudDB) GetConfig(tx interface{}, namespace, name, _ string) (*specV1.Configuration, error) {
	defer utils.Trace(d.Log.Debug, "GetConfiguration")()
	transaction, err := d.InterfaceToTx(tx)
	if err != nil {
		return nil, err
	}
	return d.GetConfigurationTx(transaction, namespace, name)
}

func (d *BaetylCloudDB) CreateConfig(tx interface{}, namespace string, configuration *specV1.Configuration) (*specV1.Configuration, error) {
	var config *specV1.Configuration
	transaction, err := d.InterfaceToTx(tx)
	if err != nil {
		return nil, err
	}
	defer utils.Trace(d.Log.Debug, "CreateConfiguration")()
	if transaction == nil {
		err = d.Transact(func(tx *sqlx.Tx) error {
			return d.createAndGetConfig(tx, namespace, configuration, &config)
		})
	} else {
		err = d.createAndGetConfig(transaction, namespace, configuration, &config)
	}
	return config, err
}

func (d *BaetylCloudDB) createAndGetConfig(tx *sqlx.Tx, ns string, in *specV1.Configuration, out **specV1.Configuration) error {
	_, err := d.CreateConfigurationTx(tx, ns, in)
	if err != nil {
		return err
	}
	*out, err = d.GetConfigurationTx(tx, ns, in.Name)
	return err
}

func (d *BaetylCloudDB) UpdateConfig(tx interface{}, namespace string, configuration *specV1.Configuration) (*specV1.Configuration, error) {
	defer utils.Trace(d.Log.Debug, "UpdateConfiguration")()
	transaction, err := d.InterfaceToTx(tx)
	if err != nil {
		return nil, err
	}
	oldConfiguration, err := d.GetConfigurationTx(transaction, namespace, configuration.Name)
	if err != nil {
		return nil, err
	}
	if entities.EqualConfig(oldConfiguration, configuration) {
		configuration.Version = oldConfiguration.Version
		return configuration, nil
	}
	_, err = d.UpdateConfigurationTx(transaction, namespace, configuration)
	if err != nil {
		return nil, err
	}
	return configuration, err
}

func (d *BaetylCloudDB) DeleteConfig(tx interface{}, namespace, name string) error {
	defer utils.Trace(d.Log.Debug, "DeleteConfiguration")()
	transaction, err := d.InterfaceToTx(tx)
	if err != nil {
		return err
	}
	_, err = d.DeleteConfigurationTx(transaction, namespace, name)
	return err
}

func (d *BaetylCloudDB) ListConfig(namespace string, listOptions *models.ListOptions) (*models.ConfigurationList, error) {
	defer utils.Trace(d.Log.Debug, "ListConfiguration")()
	configs, resLen, err := d.ListConfigurationTx(nil, namespace, listOptions)
	if err != nil {
		return nil, err
	}

	result := &models.ConfigurationList{
		Total:       resLen,
		ListOptions: listOptions,
		Items:       configs,
	}
	return result, nil
}

func (d *BaetylCloudDB) GetConfigurationTx(tx *sqlx.Tx, namespace, name string) (*specV1.Configuration, error) {
	selectSQL := `
SELECT 
id, namespace, name, labels, data, version, is_system, description, create_time, update_time
FROM baetyl_configuration WHERE namespace=? AND name=?
`
	var configs []entities.Configuration
	if err := d.Query(tx, selectSQL, &configs, namespace, name); err != nil {
		return nil, err
	}
	if len(configs) > 0 {
		return entities.ToConfigModel(&configs[0])
	}
	return nil, common.Error(
		common.ErrResourceNotFound,
		common.Field("type", "configuration"),
		common.Field("name", name))
}

func (d *BaetylCloudDB) CreateConfigurationTx(tx *sqlx.Tx, namespace string, configuration *specV1.Configuration) (sql.Result, error) {
	insertSQL := `
INSERT INTO baetyl_configuration (namespace, name, labels, data, version, is_system, description)
VALUES (?, ?, ?, ?, ?, ?, ?)
`
	config, err := entities.FromConfigModel(namespace, configuration)
	if err != nil {
		return nil, err
	}
	return d.Exec(tx, insertSQL, config.Namespace, config.Name, config.Labels, config.Data, config.Version, config.System, config.Description)
}

func (d *BaetylCloudDB) UpdateConfigurationTx(tx *sqlx.Tx, namespace string, configuration *specV1.Configuration) (sql.Result, error) {
	updateSQL := `
UPDATE baetyl_configuration
SET version = ?, description = ?, labels = ?, data = ?
WHERE namespace=? AND name=?
`
	config, err := entities.FromConfigModel(namespace, configuration)
	if err != nil {
		return nil, err
	}
	configuration.Version = config.Version
	return d.Exec(tx, updateSQL, config.Version, config.Description, config.Labels, config.Data, config.Namespace, config.Name)
}

func (d *BaetylCloudDB) ListConfigurationTx(_ *sqlx.Tx, namespace string, listOptions *models.ListOptions) ([]specV1.Configuration, int, error) {
	selectSQL := `
SELECT 
id, namespace, name, labels, data, version, is_system, description, create_time, update_time
FROM baetyl_configuration WHERE namespace=? AND name LIKE ? ORDER BY create_time DESC
`
	var configs []entities.Configuration
	if err := d.Query(nil, selectSQL, &configs, namespace, listOptions.GetFuzzyName()); err != nil {
		return nil, 0, err
	}
	result := make([]specV1.Configuration, 0)
	for _, config := range configs {
		labels := map[string]string{}
		if err := json.Unmarshal([]byte(config.Labels), &labels); err != nil {
			return nil, 0, errors.Trace(err)
		}
		if ok, err := utils.IsLabelMatch(listOptions.LabelSelector, labels); err != nil || !ok {
			continue
		}
		cfg, err := entities.ToConfigModel(&config)
		if err != nil {
			return nil, 0, errors.Trace(err)
		}
		result = append(result, *cfg)
	}
	start, end := models.GetPagingParam(listOptions, len(result))
	return result[start:end], len(result), nil
}

func (d *BaetylCloudDB) DeleteConfigurationTx(tx *sqlx.Tx, namespace, name string) (sql.Result, error) {
	deleteSQL := `
DELETE FROM baetyl_configuration WHERE namespace=? AND name=?
`
	return d.Exec(tx, deleteSQL, namespace, name)
}
