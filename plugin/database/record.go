// Package database 数据库存储实现
package database

import (
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

func (d *BaetylCloudDB) CountRecord(batchName, fingerprintValue, ns string) (int, error) {
	return d.CountRecordTx(nil, batchName, fingerprintValue, ns)
}

func (d *BaetylCloudDB) GetRecord(batchName, recordName, ns string) (*models.Record, error) {
	return d.GetRecordTx(nil, batchName, recordName, ns)
}

func (d *BaetylCloudDB) GetRecordByFingerprint(batchName, ns, value string) (*models.Record, error) {
	return d.GetRecordByFingerprintTx(nil, batchName, ns, value)
}

func (d *BaetylCloudDB) ListRecord(batchName, ns string, filter *models.Filter) ([]models.Record, error) {
	return d.ListRecordTx(nil, batchName, ns, filter)
}

func (d *BaetylCloudDB) CreateRecord(records []models.Record) (sql.Result, error) {
	return d.CreateRecordTx(nil, records)
}

func (d *BaetylCloudDB) UpdateRecord(record *models.Record) (sql.Result, error) {
	return d.UpdateRecordTx(nil, record)
}

func (d *BaetylCloudDB) DeleteRecord(batchName, recordName, ns string) (sql.Result, error) {
	return d.DeleteRecordTx(nil, batchName, recordName, ns)
}

func (d *BaetylCloudDB) CountRecordTx(tx *sqlx.Tx, batchName, fingerprintValue, ns string) (int, error) {
	selectSQL := `
SELECT count(name) AS count
FROM baetyl_batch_record WHERE namespace=? AND batch_name=? AND fingerprint_value LIKE ?
`
	var res []struct {
		Count int `db:"count"`
	}
	if err := d.Query(tx, selectSQL, &res, ns, batchName, fingerprintValue); err != nil {
		return 0, err
	}
	return res[0].Count, nil
}

func (d *BaetylCloudDB) GetRecordTx(tx *sqlx.Tx, batchName, recordName, ns string) (*models.Record, error) {
	selectSQL := `
SELECT 
name, batch_name, namespace, fingerprint_value, 
active, node_name, active_ip, active_time, create_time, 
update_time 
FROM baetyl_batch_record WHERE namespace=? AND batch_name=? AND name=?
`
	var record []models.Record
	if err := d.Query(tx, selectSQL, &record, ns, batchName, recordName); err != nil {
		return nil, err
	}
	if len(record) > 0 {
		return &record[0], nil
	}
	return nil, nil
}

func (d *BaetylCloudDB) GetRecordByFingerprintTx(tx *sqlx.Tx, batchName, ns, value string) (*models.Record, error) {
	selectSQL := `
SELECT 
name, batch_name, namespace, fingerprint_value, 
active, node_name, active_ip, active_time, create_time, 
update_time 
FROM baetyl_batch_record WHERE namespace=? AND batch_name=? AND fingerprint_value=?
`
	var record []models.Record
	if err := d.Query(tx, selectSQL, &record, ns, batchName, value); err != nil {
		return nil, err
	}
	if len(record) > 0 {
		return &record[0], nil
	}
	return nil, nil
}

func (d *BaetylCloudDB) ListRecordByBatchTx(tx *sqlx.Tx, batchName, namespace string) ([]models.Record, error) {
	selectSQL := `
SELECT 
name, batch_name, namespace, fingerprint_value, 
active, node_name, active_ip, active_time, create_time, 
update_time 
FROM baetyl_batch_record WHERE namespace=? and batch_name=? ORDER BY create_time DESC
`
	records := []models.Record{}
	if err := d.Query(tx, selectSQL, &records, namespace, batchName); err != nil {
		return nil, err
	}
	return records, nil
}

func (d *BaetylCloudDB) ListRecordTx(tx *sqlx.Tx, batchName, ns string, filter *models.Filter) ([]models.Record, error) {
	selectSQL := `
SELECT 
name, batch_name, namespace, fingerprint_value, 
active, node_name, active_ip, active_time, create_time, 
update_time 
FROM baetyl_batch_record WHERE namespace=? AND batch_name=? AND fingerprint_value LIKE ? ORDER BY create_time DESC 
`
	records := []models.Record{}
	args := []interface{}{ns, batchName, filter.GetFuzzyName()}
	if filter.GetLimitNumber() > 0 {
		selectSQL = selectSQL + "LIMIT ?,?"
		args = append(args, filter.GetLimitOffset(), filter.GetLimitNumber())
	}
	if err := d.Query(tx, selectSQL, &records, args...); err != nil {
		return nil, err
	}
	return records, nil
}

func (d *BaetylCloudDB) CreateRecordTx(tx *sqlx.Tx, records []models.Record) (sql.Result, error) {
	selectSQL := `
INSERT INTO baetyl_batch_record (
name, batch_name, namespace, fingerprint_value,
active, node_name, active_ip, active_time)
VALUES 
`
	vals := []interface{}{}
	for _, record := range records {
		selectSQL += "(?,?,?,?,?,?,?,?),"
		vals = append(vals, record.Name, record.BatchName, record.Namespace, record.FingerprintValue,
			record.Active, record.NodeName, record.ActiveIP, record.ActiveTime)
	}
	return d.Exec(tx, selectSQL[0:len(selectSQL)-1], vals...)
}

func (d *BaetylCloudDB) UpdateRecordTx(tx *sqlx.Tx, record *models.Record) (sql.Result, error) {
	selectSQL := `
UPDATE baetyl_batch_record
SET active=?,
    node_name=?,
    active_ip=?,
    active_time=?
WHERE namespace=? AND batch_name=? AND name = ?;
`
	return d.Exec(tx, selectSQL, record.Active, record.NodeName,
		record.ActiveIP, record.ActiveTime, record.Namespace, record.BatchName, record.Name)
}

func (d *BaetylCloudDB) DeleteRecordTx(tx *sqlx.Tx, batchName, recordName, ns string) (sql.Result, error) {
	selectSQL := `
DELETE FROM baetyl_batch_record WHERE namespace=? AND batch_name=? AND name=?
`
	return d.Exec(tx, selectSQL, ns, batchName, recordName)
}
