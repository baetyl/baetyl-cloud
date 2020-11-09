package database

import (
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin/database/entities"
)

// TODO: 抽象 batch 操作的 interface
func (d *DB) GetBatch(name, ns string) (*models.Batch, error) {
	return d.GetBatchTx(nil, name, ns)
}

func (d *DB) ListBatch(ns string, filter *models.Filter) ([]models.Batch, error) {
	return d.ListBatchTx(nil, ns, filter)
}

func (d *DB) CreateBatch(batch *models.Batch) (sql.Result, error) {
	return d.CreateBatchTx(nil, batch)
}

func (d *DB) UpdateBatch(batch *models.Batch) (sql.Result, error) {
	return d.UpdateBatchTx(nil, batch)
}

func (d *DB) DeleteBatch(name, ns string) (sql.Result, error) {
	return d.DeleteBatchTx(nil, name, ns)
}

func (d *DB) CountBatch(ns, name string) (int, error) {
	return d.CountBatchTx(nil, ns, name)
}

func (d *DB) CountBatchByCallback(callbackName, ns string) (int, error) {
	return d.CountBatchByCallbackTx(nil, callbackName, ns)
}

func (d *DB) GetBatchTx(tx *sqlx.Tx, name, ns string) (*models.Batch, error) {
	selectSQL := `
SELECT  
name, namespace, description, quota_num, enable_whitelist,
security_type, security_key, callback_name,
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

func (d *DB) ListBatchTx(tx *sqlx.Tx, ns string, filter *models.Filter) ([]models.Batch, error) {
	selectSQL := `
SELECT  
name, namespace, description, quota_num, enable_whitelist,
security_type, security_key, callback_name,
labels, fingerprint, create_time, update_time 
FROM baetyl_batch WHERE namespace=? AND name LIKE ? ORDER BY create_time DESC 
`
	batchs := []entities.Batch{}
	args := []interface{}{ns, filter.GetFuzzyName()}
	if filter.GetLimitNumber() > 0 {
		selectSQL = selectSQL + "LIMIT ?,?"
		args = append(args, filter.GetLimitOffset(), filter.GetLimitNumber())
	}
	if err := d.Query(tx, selectSQL, &batchs, args...); err != nil {
		return nil, err
	}
	var res []models.Batch
	for _, b := range batchs {
		res = append(res, *entities.ToBatchModel(&b))
	}
	return res, nil
}

func (d *DB) CreateBatchTx(tx *sqlx.Tx, batch *models.Batch) (sql.Result, error) {
	insertSQL := `
INSERT INTO baetyl_batch 
(name, namespace, description, quota_num, 
enable_whitelist, security_type, security_key, 
callback_name, labels, fingerprint) 
VALUES 
(?,?,?,?,?,?,?,?,?,?)
`
	batchDB := entities.FromBatchModel(batch)
	return d.Exec(tx, insertSQL, batchDB.Name, batchDB.Namespace, batchDB.Description,
		batchDB.QuotaNum, batchDB.EnableWhitelist, batchDB.SecurityType, batchDB.SecurityKey,
		batchDB.CallbackName, batchDB.Labels, batchDB.Fingerprint)
}

func (d *DB) UpdateBatchTx(tx *sqlx.Tx, batch *models.Batch) (sql.Result, error) {
	updateSQL := `
UPDATE baetyl_batch SET description=?,quota_num=?,
callback_name=?,labels=?,fingerprint=?
WHERE namespace=? AND name=?
`
	batchDB := entities.FromBatchModel(batch)
	return d.Exec(tx, updateSQL, batchDB.Description, batchDB.QuotaNum,
		batchDB.CallbackName, batchDB.Labels, batchDB.Fingerprint,
		batchDB.Namespace, batchDB.Name)
}

func (d *DB) DeleteBatchTx(tx *sqlx.Tx, name, ns string) (sql.Result, error) {
	deleteSQL := `
DELETE FROM baetyl_batch where namespace=? AND name=?
`
	return d.Exec(tx, deleteSQL, ns, name)
}

func (d *DB) CountBatchTx(tx *sqlx.Tx, ns, name string) (int, error) {
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

func (d *DB) CountBatchByCallbackTx(tx *sqlx.Tx, callbackName, ns string) (int, error) {
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
