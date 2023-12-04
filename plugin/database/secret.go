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

func (d *BaetylCloudDB) GetSecret(tx interface{}, namespace, name, _ string) (*specV1.Secret, error) {
	defer utils.Trace(d.Log.Debug, "GetSecret")()
	transaction, err := d.InterfaceToTx(tx)
	if err != nil {
		return nil, err
	}
	return d.GetSecretTx(transaction, namespace, name)
}

func (d *BaetylCloudDB) CreateSecret(tx interface{}, namespace string, secret *specV1.Secret) (*specV1.Secret, error) {
	var se *specV1.Secret
	transaction, err := d.InterfaceToTx(tx)
	if err != nil {
		return nil, err
	}
	defer utils.Trace(d.Log.Debug, "CreateSecret")()
	if transaction == nil {
		err = d.Transact(func(tx *sqlx.Tx) error {
			return d.createAndGetSecret(tx, namespace, secret, &se)
		})
	} else {
		err = d.createAndGetSecret(transaction, namespace, secret, &se)
	}
	return se, err
}

func (d *BaetylCloudDB) createAndGetSecret(tx *sqlx.Tx, ns string, in *specV1.Secret, out **specV1.Secret) error {
	_, err := d.CreateSecretTx(tx, ns, in)
	if err != nil {
		return err
	}
	*out, err = d.GetSecretTx(tx, ns, in.Name)
	return err
}

func (d *BaetylCloudDB) UpdateSecret(namespace string, secret *specV1.Secret) (*specV1.Secret, error) {
	var se *specV1.Secret
	defer utils.Trace(d.Log.Debug, "UpdateSecret")()
	err := d.Transact(func(tx *sqlx.Tx) error {
		_, err := d.UpdateSecretTx(tx, namespace, secret)
		if err != nil {
			return err
		}
		se, err = d.GetSecretTx(tx, namespace, secret.Name)
		return err
	})
	return se, err
}

func (d *BaetylCloudDB) DeleteSecret(tx interface{}, namespace, name string) error {
	defer utils.Trace(d.Log.Debug, "DeleteSecret")()
	transaction, err := d.InterfaceToTx(tx)
	if err != nil {
		return err
	}
	_, err = d.DeleteSecretTx(transaction, namespace, name)
	return err
}

func (d *BaetylCloudDB) ListSecret(namespace string, listOptions *models.ListOptions) (*models.SecretList, error) {
	defer utils.Trace(d.Log.Debug, "ListSecret")()
	secrets, resLen, err := d.ListSecretTx(nil, namespace, listOptions)
	if err != nil {
		return nil, err
	}

	result := &models.SecretList{
		Total:       resLen,
		ListOptions: listOptions,
		Items:       secrets,
	}
	return result, nil
}

func (d *BaetylCloudDB) GetSecretTx(tx *sqlx.Tx, namespace, name string) (*specV1.Secret, error) {
	selectSQL := `
SELECT 
id, namespace, name, labels, data, version, is_system, description, create_time, update_time
FROM baetyl_secret WHERE namespace=? AND name=?
`
	var secrets []entities.Secret
	if err := d.Query(tx, selectSQL, &secrets, namespace, name); err != nil {
		return nil, err
	}
	if len(secrets) > 0 {
		return entities.ToSecretModel(&secrets[0])

	}
	return nil, common.Error(
		common.ErrResourceNotFound,
		common.Field("type", "secret"),
		common.Field("name", name))
}

func (d *BaetylCloudDB) CreateSecretTx(tx *sqlx.Tx, namespace string, secret *specV1.Secret) (sql.Result, error) {
	insertSQL := `
INSERT INTO baetyl_secret (namespace, name, labels, data, version, is_system, description)
VALUES (?, ?, ?, ?, ?, ?, ?)
`
	se, err := entities.FromSecretModel(namespace, secret)
	if err != nil {
		return nil, err
	}
	return d.Exec(tx, insertSQL, se.Namespace, se.Name, se.Labels, se.Data, se.Version, se.System, se.Description)
}

func (d *BaetylCloudDB) UpdateSecretTx(tx *sqlx.Tx, namespace string, secret *specV1.Secret) (sql.Result, error) {
	updateSQL := `
UPDATE baetyl_secret
SET version = ?, description = ?, labels = ?, data = ?
WHERE namespace=? AND name=?
`
	se, err := entities.FromSecretModel(namespace, secret)
	if err != nil {
		return nil, err
	}
	return d.Exec(tx, updateSQL, se.Version, se.Description, se.Labels, se.Data, se.Namespace, se.Name)
}

func (d *BaetylCloudDB) ListSecretTx(_ *sqlx.Tx, namespace string, listOptions *models.ListOptions) ([]specV1.Secret, int, error) {
	selectSQL := `
SELECT 
id, namespace, name, labels, data, version, is_system, description, create_time, update_time
FROM baetyl_secret WHERE namespace=? AND name LIKE ? ORDER BY create_time DESC
`
	var secrets []entities.Secret
	if err := d.Query(nil, selectSQL, &secrets, namespace, listOptions.GetFuzzyName()); err != nil {
		return nil, 0, err
	}
	var result []specV1.Secret
	for _, se := range secrets {
		labels := map[string]string{}
		if err := json.Unmarshal([]byte(se.Labels), &labels); err != nil {
			return nil, 0, errors.Trace(err)
		}
		if ok, err := utils.IsLabelMatch(listOptions.LabelSelector, labels); err != nil || !ok {
			continue
		}
		s, err := entities.ToSecretModel(&se)
		if err != nil {
			return nil, 0, errors.Trace(err)
		}
		result = append(result, *s)
	}
	start, end := models.GetPagingParam(listOptions, len(result))
	return result[start:end], len(result), nil
}

func (d *BaetylCloudDB) DeleteSecretTx(tx *sqlx.Tx, namespace, name string) (sql.Result, error) {
	deleteSQL := `
DELETE FROM baetyl_secret WHERE namespace=? AND name=?
`
	return d.Exec(tx, deleteSQL, namespace, name)
}
