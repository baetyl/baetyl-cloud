// Package database 数据库存储实现
package database

import (
	"database/sql"

	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/baetyl/baetyl-go/v2/utils"
	"github.com/jmoiron/sqlx"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin/database/entities"
)

func (d *BaetylCloudDB) GetNamespace(namespace string) (*models.Namespace, error) {
	defer utils.Trace(d.Log.Debug, "GetNamespace")()
	return d.GetNamespaceTx(nil, namespace)
}

func (d *BaetylCloudDB) CreateNamespace(namespace *models.Namespace) (*models.Namespace, error) {
	var ns *models.Namespace
	defer utils.Trace(d.Log.Debug, "CreateNamespace")()

	err := d.Transact(func(tx *sqlx.Tx) error {
		_, err := d.CreateNamespaceTx(tx, namespace)
		if err != nil {
			return err
		}
		ns, err = d.GetNamespaceTx(tx, namespace.Name)
		return err
	})
	return ns, err
}

func (d *BaetylCloudDB) ListNamespace(listOptions *models.ListOptions) (*models.NamespaceList, error) {
	defer utils.Trace(d.Log.Debug, "ListNamespace")()
	nss, resLen, err := d.ListNamespaceTx(listOptions)
	if err != nil {
		return nil, err
	}
	result := &models.NamespaceList{
		Total:       resLen,
		ListOptions: listOptions,
		Items:       nss,
	}
	return result, nil
}

// TODO：增加资源清理
func (d *BaetylCloudDB) DeleteNamespace(namespace *models.Namespace) error {
	defer utils.Trace(d.Log.Debug, "DeleteNamespace")()
	_, err := d.DeleteNamespaceTx(nil, namespace.Name)
	if err != nil {
		return err
	}
	_, err = d.DeleteNodesTx(nil, namespace.Name)
	if err != nil {
		return err
	}
	_, err = d.DeleteAppsTx(nil, namespace.Name)
	if err != nil {
		return err
	}
	_, err = d.DeleteConfigsTx(nil, namespace.Name)
	if err != nil {
		return err
	}
	_, err = d.DeleteSecretsTx(nil, namespace.Name)
	if err != nil {
		return err
	}

	return err
}

func (d *BaetylCloudDB) GetNamespaceTx(tx *sqlx.Tx, name string) (*models.Namespace, error) {
	selectSQL := `
SELECT 
id, name
FROM baetyl_namespace WHERE name=?
`
	var ns []entities.Namespace
	if err := d.Query(tx, selectSQL, &ns, name); err != nil {
		return nil, err
	}
	if len(ns) > 0 {
		return entities.ToNamespaceModel(&ns[0])
	}
	return nil, common.Error(
		common.ErrResourceNotFound,
		common.Field("type", "namespace"),
		common.Field("name", name))
}

func (d *BaetylCloudDB) CreateNamespaceTx(tx *sqlx.Tx, namespace *models.Namespace) (sql.Result, error) {
	insertSQL := `
INSERT INTO baetyl_namespace (name)
VALUES (?)
`
	ns, err := entities.FromNamespaceModel(namespace)
	if err != nil {
		return nil, err
	}

	return d.Exec(tx, insertSQL, ns.Name)
}

func (d *BaetylCloudDB) ListNamespaceTx(listOptions *models.ListOptions) ([]models.Namespace, int, error) {
	selectSQL := `
SELECT 
id, name
FROM baetyl_namespace WHERE name LIKE ? ORDER BY create_time DESC
`
	var nss []entities.Namespace
	if err := d.Query(nil, selectSQL, &nss, listOptions.GetFuzzyName()); err != nil {
		return nil, 0, err
	}
	var result []models.Namespace
	for _, ns := range nss {
		nd, err := entities.ToNamespaceModel(&ns)
		if err != nil {
			return nil, 0, errors.Trace(err)
		}
		result = append(result, *nd)
	}
	start, end := models.GetPagingParam(listOptions, len(result))
	return result[start:end], len(result), nil
}

func (d *BaetylCloudDB) DeleteNamespaceTx(tx *sqlx.Tx, name string) (sql.Result, error) {
	deleSQL := `
DELETE FROM baetyl_namespace WHERE name=?
`
	return d.Exec(tx, deleSQL, name)
}

func (d *BaetylCloudDB) DeleteAppsTx(tx *sqlx.Tx, ns string) (sql.Result, error) {
	deleSQL := `DELETE FROM baetyl_application WHERE namespace=?`
	return d.Exec(tx, deleSQL, ns)
}

func (d *BaetylCloudDB) DeleteNodesTx(tx *sqlx.Tx, ns string) (sql.Result, error) {
	deleSQL := `DELETE FROM baetyl_node WHERE namespace=?`
	return d.Exec(tx, deleSQL, ns)
}

func (d *BaetylCloudDB) DeleteSecretsTx(tx *sqlx.Tx, ns string) (sql.Result, error) {
	deleSQL := `DELETE FROM baetyl_secret WHERE namespace=?`
	return d.Exec(tx, deleSQL, ns)
}

func (d *BaetylCloudDB) DeleteConfigsTx(tx *sqlx.Tx, ns string) (sql.Result, error) {
	deleSQL := `DELETE FROM baetyl_configuration WHERE namespace=?`
	return d.Exec(tx, deleSQL, ns)
}
