// Package database 数据库存储实现
package database

import (
	"database/sql"

	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/baetyl/baetyl-go/v2/json"
	"github.com/baetyl/baetyl-go/v2/utils"
	"github.com/jmoiron/sqlx"

	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin/database/entities"
)

// TODO: 抽象 batch 操作的 interface
func (d *BaetylCloudDB) GetBatch(name, ns string) (*models.Batch, error) {
	return d.GetBatchTx(nil, name, ns)
}

func (d *BaetylCloudDB) ListBatch(ns string, filter *models.ListOptions) ([]models.Batch, error) {
	return d.ListBatchTx(nil, ns, filter)
}

func (d *BaetylCloudDB) CreateBatch(batch *models.Batch) (sql.Result, error) {
	return d.CreateBatchTx(nil, batch)
}

func (d *BaetylCloudDB) UpdateBatch(batch *models.Batch) (sql.Result, error) {
	return d.UpdateBatchTx(nil, batch)
}

func (d *BaetylCloudDB) DeleteBatch(name, ns string) (sql.Result, error) {
	return d.DeleteBatchTx(nil, name, ns)
}

func (d *BaetylCloudDB) CountBatch(ns, name string) (int, error) {
	return d.CountBatchTx(nil, ns, name)
}

func (d *BaetylCloudDB) CountBatchByCallback(callbackName, ns string) (int, error) {
	return d.CountBatchByCallbackTx(nil, callbackName, ns)
}

func (d *BaetylCloudDB) GetBatchTx(tx *sqlx.Tx, name, ns string) (*models.Batch, error) {
	selectSQL := `
SELECT  
name, namespace, description, quota_num, accelerator, sys_apps, enable_whitelist,
security_type, security_key, callback_name, cluster, 
labels, fingerprint, create_time, update_time 
FROM baetyl_batch WHERE namespace=? AND name=? LIMIT 0,1
`
	batchs := []entities.Batch{}
	if err := d.Query(tx, selectSQL, &batchs, ns, name); err != nil {
		return nil, err
	}
	if len(batchs) > 0 {
		return entities.ToBatchModel(&batchs[0]), nil
	}
	return nil, nil
}

func (d *BaetylCloudDB) ListBatchTx(tx *sqlx.Tx, ns string, filter *models.ListOptions) ([]models.Batch, error) {
	selectSQL := `
SELECT  
name, namespace, description, quota_num, accelerator, sys_apps, enable_whitelist,
security_type, security_key, callback_name, cluster, 
labels, fingerprint, create_time, update_time 
FROM baetyl_batch WHERE namespace=? AND name LIKE ? ORDER BY create_time DESC 
`
	batchs := []entities.Batch{}
	args := []interface{}{ns, filter.GetFuzzyName()}
	// label match not support filter
	// if filter.GetLimitNumber() > 0 {
	//	 selectSQL = selectSQL + "LIMIT ?,?"
	//	 args = append(args, filter.GetLimitOffset(), filter.GetLimitNumber())
	// }
	if err := d.Query(tx, selectSQL, &batchs, args...); err != nil {
		return nil, err
	}
	var res []models.Batch
	for _, b := range batchs {
		labels := map[string]string{}
		if err := json.Unmarshal([]byte(b.Labels), &labels); err != nil {
			return nil, errors.Trace(err)
		}
		if ok, err := utils.IsLabelMatch(filter.LabelSelector, labels); err != nil || !ok {
			continue
		}
		res = append(res, *entities.ToBatchModel(&b))
	}
	return res, nil
}

func (d *BaetylCloudDB) CreateBatchTx(tx *sqlx.Tx, batch *models.Batch) (sql.Result, error) {
	insertSQL := `
INSERT INTO baetyl_batch 
(name, namespace, description, quota_num, accelerator,
sys_apps, enable_whitelist, security_type, security_key, 
callback_name, labels, fingerprint, cluster) 
VALUES 
(?,?,?,?,?,?,?,?,?,?,?,?,?)
`
	batchDB := entities.FromBatchModel(batch)
	return d.Exec(tx, insertSQL, batchDB.Name, batchDB.Namespace, batchDB.Description,
		batchDB.QuotaNum, batchDB.Accelerator, batchDB.SysApps, batchDB.EnableWhitelist, batchDB.SecurityType, batchDB.SecurityKey,
		batchDB.CallbackName, batchDB.Labels, batchDB.Fingerprint, batchDB.Cluster)
}

func (d *BaetylCloudDB) UpdateBatchTx(tx *sqlx.Tx, batch *models.Batch) (sql.Result, error) {
	updateSQL := `
UPDATE baetyl_batch SET description=?,quota_num=?,accelerator=?,
sys_apps=?,callback_name=?,labels=?,fingerprint=?
WHERE namespace=? AND name=?
`
	batchDB := entities.FromBatchModel(batch)
	return d.Exec(tx, updateSQL, batchDB.Description, batchDB.QuotaNum, batchDB.Accelerator,
		batchDB.SysApps, batchDB.CallbackName, batchDB.Labels, batchDB.Fingerprint,
		batchDB.Namespace, batchDB.Name)
}

func (d *BaetylCloudDB) DeleteBatchTx(tx *sqlx.Tx, name, ns string) (sql.Result, error) {
	deleteSQL := `
DELETE FROM baetyl_batch where namespace=? AND name=?
`
	return d.Exec(tx, deleteSQL, ns, name)
}

func (d *BaetylCloudDB) CountBatchTx(tx *sqlx.Tx, ns, name string) (int, error) {
	selectSQL := `
SELECT count(name) AS count
FROM baetyl_batch WHERE namespace=? AND name LIKE ?
`
	var res []struct {
		Count int `db:"count"`
	}
	if err := d.Query(tx, selectSQL, &res, ns, name); err != nil {
		return 0, err
	}
	return res[0].Count, nil
}

func (d *BaetylCloudDB) CountBatchByCallbackTx(tx *sqlx.Tx, callbackName, ns string) (int, error) {
	selectSQL := `
SELECT count(name) AS count
FROM baetyl_batch WHERE namespace=? AND callback_name=?
`
	var res []struct {
		Count int `db:"count"`
	}
	if err := d.Query(tx, selectSQL, &res, ns, callbackName); err != nil {
		return 0, err
	}
	return res[0].Count, nil
}
