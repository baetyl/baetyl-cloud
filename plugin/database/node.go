// Package database 数据库存储实现
package database

import (
	"database/sql"
	"strings"

	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/baetyl/baetyl-go/v2/json"
	"github.com/baetyl/baetyl-go/v2/log"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/baetyl/baetyl-go/v2/utils"
	"github.com/jmoiron/sqlx"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin/database/entities"
)

func (d *BaetylCloudDB) GetNode(tx interface{}, namespace, name string) (*specV1.Node, error) {
	defer utils.Trace(d.Log.Debug, "GetNode")()
	transaction, err := d.InterfaceToTx(tx)
	if err != nil {
		return nil, err
	}
	return d.GetNodeTx(transaction, namespace, name)
}

func (d *BaetylCloudDB) CreateNode(tx interface{}, namespace string, node *specV1.Node) (*specV1.Node, error) {
	var nd *specV1.Node
	transaction, err := d.InterfaceToTx(tx)
	if err != nil {
		return nil, err
	}
	defer utils.Trace(d.Log.Debug, "CreateNode")()
	if transaction == nil {
		err = d.Transact(func(sqlTx *sqlx.Tx) error {
			return d.createAndGetNode(sqlTx, namespace, node, &nd)
		})
	} else {
		err = d.createAndGetNode(transaction, namespace, node, &nd)
	}
	return nd, err
}

func (d *BaetylCloudDB) createAndGetNode(tx *sqlx.Tx, ns string, in *specV1.Node, out **specV1.Node) error {
	_, err := d.CreateNodeTx(tx, ns, in)
	if err != nil {
		return err
	}
	*out, err = d.GetNodeTx(tx, ns, in.Name)
	return err
}

// TODO: The node update returns the source node. If the node contains a field that changes every time it is updated,
// TODO: the update node needs to return the updated node
func (d *BaetylCloudDB) UpdateNode(tx interface{}, namespace string, nodes []*specV1.Node) ([]*specV1.Node, error) {
	defer utils.Trace(d.Log.Debug, "UpdateNode")()
	transaction, err := d.InterfaceToTx(tx)
	if err != nil {
		return nil, err
	}
	if transaction == nil {
		err = d.Transact(func(sqlTx *sqlx.Tx) error {
			return d.getAndUpdateNodeTx(sqlTx, namespace, nodes)
		})
	} else {
		err = d.getAndUpdateNodeTx(transaction, namespace, nodes)
	}
	return nodes, err
}

func (d *BaetylCloudDB) getAndUpdateNodeTx(tx *sqlx.Tx, namespace string, nodes []*specV1.Node) error {
	for _, node := range nodes {
		oldNd, err := d.GetNodeTx(tx, namespace, node.Name)
		if err != nil {
			return err
		}
		if entities.EqualNode(node, oldNd) {
			return nil
		}
	}
	_, err := d.UpdateNodeTx(tx, namespace, nodes)
	if err != nil {
		return err
	}
	return err
}

func (d *BaetylCloudDB) DeleteNode(tx interface{}, namespace, name string) error {
	defer utils.Trace(d.Log.Debug, "DeleteNode")()
	transaction, err := d.InterfaceToTx(tx)
	if err != nil {
		return err
	}
	_, err = d.DeleteNodeTx(transaction, namespace, name)
	return err
}

func (d *BaetylCloudDB) ListNode(tx interface{}, namespace string, listOptions *models.ListOptions) (*models.NodeList, error) {
	defer utils.Trace(d.Log.Debug, "ListNode")()
	transaction, err := d.InterfaceToTx(tx)
	if err != nil {
		return nil, err
	}
	nodes, resLen, err := d.ListNodeTx(transaction, namespace, listOptions)
	if err != nil {
		return nil, err
	}
	result := &models.NodeList{
		Total:       resLen,
		ListOptions: listOptions,
		Items:       nodes,
	}
	return result, nil
}

func (d *BaetylCloudDB) CountAllNode(tx interface{}) (int, error) {
	defer utils.Trace(d.Log.Debug, "CountAllNode")()
	transaction, err := d.InterfaceToTx(tx)
	if err != nil {
		return 0, err
	}
	return d.CountAllNodeTx(transaction)
}

func (d *BaetylCloudDB) CountAllNodeTx(tx *sqlx.Tx) (int, error) {
	selectSQL := `
SELECT count(name) AS count FROM baetyl_node
`
	var res []struct {
		Count int `db:"count"`
	}
	if err := d.Query(tx, selectSQL, &res); err != nil {
		return 0, err
	}
	return res[0].Count, nil
}

func (d *BaetylCloudDB) GetNodeTx(tx *sqlx.Tx, namespace, name string) (*specV1.Node, error) {
	selectSQL := `
SELECT 
id, namespace, name, version, core_version, node_mode, description, create_time, labels, annotations, attributes
FROM baetyl_node WHERE namespace=? AND name=?
`
	var nodes []entities.Node
	if err := d.Query(tx, selectSQL, &nodes, namespace, name); err != nil {
		return nil, err
	}
	if len(nodes) > 0 {
		return entities.ToNodeModel(&nodes[0])
	}
	return nil, common.Error(
		common.ErrResourceNotFound,
		common.Field("type", "node"),
		common.Field("name", name))
}

func (d *BaetylCloudDB) CreateNodeTx(tx *sqlx.Tx, namespace string, node *specV1.Node) (sql.Result, error) {
	insertSQL := `
INSERT INTO baetyl_node (namespace, name, version, core_version, node_mode, description, labels, annotations, attributes)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
`
	nd, err := entities.FromNodeModel(namespace, node)
	if err != nil {
		return nil, err
	}
	return d.Exec(tx, insertSQL, nd.Namespace, nd.Name, nd.Version, nd.CoreVersion, nd.NodeMode, nd.Description, nd.Labels, nd.Annotations, nd.Attributes)
}

// UpdateNodeTx 更新语句示例:（有两个节点的情况）
//
// UPDATE baetyl_node
// SET version  = CASE name WHEN ? THEN ? WHEN ? THEN ? END,
// core_version = CASE name WHEN ? THEN ? WHEN ? THEN ? END,
// description  = CASE name WHEN ? THEN ? WHEN ? THEN ? END,
// labels       = CASE name WHEN ? THEN ? WHEN ? THEN ? END,
// annotations  = CASE name WHEN ? THEN ? WHEN ? THEN ? END,
// attributes   = CASE name WHEN ? THEN ? WHEN ? THEN ? END
// WHERE name IN (?, ?)
// AND namespace = ?
func (d *BaetylCloudDB) UpdateNodeTx(tx *sqlx.Tx, namespace string, nodes []*specV1.Node) (sql.Result, error) {
	if len(nodes) < 1 {
		return nil, nil
	}
	fields := []string{
		"version",
		"core_version",
		"node_mode",
		"description",
		"labels",
		"annotations",
		"attributes",
	}
	params := []interface{}{}

	for i := range fields {
		fields[i] += " = CASE name"
	}

	for i, f := range fields {
		field := strings.SplitN(f, " ", 2)[0]
		for _, node := range nodes {
			nd, err := entities.FromNodeModel(namespace, node)
			if err != nil {
				return nil, err
			}
			node.Version = nd.Version
			switch field {
			case "version":
				params = append(params, nd.Name, nd.Version)
			case "core_version":
				params = append(params, nd.Name, nd.CoreVersion)
			case "node_mode":
				params = append(params, nd.Name, nd.NodeMode)
			case "description":
				params = append(params, nd.Name, nd.Description)
			case "labels":
				params = append(params, nd.Name, nd.Labels)
			case "annotations":
				params = append(params, nd.Name, nd.Annotations)
			case "attributes":
				params = append(params, nd.Name, nd.Attributes)
			}
			fields[i] += " WHEN ? THEN ?"
		}
	}

	for _, node := range nodes {
		params = append(params, node.Name)
	}
	params = append(params, nodes[0].Namespace)

	updateSQL := `UPDATE baetyl_node SET ` + strings.Join(fields, ` END, `) + ` END WHERE name IN (` + strings.Repeat(`?, `, len(nodes)-1) + `?) AND namespace=?`
	d.log.Debug("UpdateNodeTx", log.Any("updateSQL", updateSQL))
	return d.Exec(tx, updateSQL, params...)
}

func (d *BaetylCloudDB) ListNodeTx(_ *sqlx.Tx, namespace string, listOptions *models.ListOptions) ([]specV1.Node, int, error) {
	selectSQL := `
SELECT 
id, namespace, name, version, core_version, node_mode, description, create_time, labels, annotations, attributes
FROM baetyl_node WHERE namespace=? AND name LIKE ? ORDER BY create_time DESC
`
	var nodes []entities.Node
	if err := d.Query(nil, selectSQL, &nodes, namespace, listOptions.GetFuzzyName()); err != nil {
		return nil, 0, err
	}
	var result []specV1.Node
	for _, node := range nodes {
		labels := map[string]string{}
		if err := json.Unmarshal([]byte(node.Labels), &labels); err != nil {
			return nil, 0, errors.Trace(err)
		}
		if ok, err := utils.IsLabelMatch(listOptions.LabelSelector, labels); err != nil || !ok {
			continue
		}
		nd, err := entities.ToNodeModel(&node)
		if err != nil {
			return nil, 0, errors.Trace(err)
		}
		result = append(result, *nd)
	}
	return result, len(result), nil
}

func (d *BaetylCloudDB) DeleteNodeTx(tx *sqlx.Tx, namespace, name string) (sql.Result, error) {
	deleteSQL := `
DELETE FROM baetyl_node WHERE namespace=? AND name=?
`
	return d.Exec(tx, deleteSQL, namespace, name)
}

func (d *BaetylCloudDB) GetNodeByNames(tx interface{}, namespace string, names []string) ([]specV1.Node, error) {
	defer utils.Trace(d.Log.Debug, "GetNodeByNames")()
	transaction, err := d.InterfaceToTx(tx)
	if err != nil {
		return nil, err
	}
	return d.GetNodeByNamesTx(transaction, namespace, names)

}
func (d *BaetylCloudDB) GetNodeByNamesTx(tx *sqlx.Tx, namespace string, names []string) ([]specV1.Node, error) {
	selectSQL := `SELECT id, namespace, name, version, core_version, node_mode, description, create_time, labels, annotations, attributes
FROM baetyl_node WHERE name in (?)
`
	qry, args, err := sqlx.In(selectSQL, names)
	if err != nil {
		return nil, err
	}
	qry += " AND namespace=?  ORDER BY create_time DESC"
	args = append(args, namespace)
	var nodes []entities.Node
	err = d.Query(tx, qry, &nodes, args...)
	if err != nil {
		return nil, err
	}
	var result []specV1.Node
	for _, node := range nodes {
		nd, err := entities.ToNodeModel(&node)
		if err != nil {
			return nil, errors.Trace(err)
		}
		result = append(result, *nd)
	}
	return result, nil
}
