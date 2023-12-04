// Package database 数据库存储实现
package database

import (
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
)

func (d *BaetylCloudDB) InsertLock(lock *models.Lock) error {
	_, err := d.InsertLockTx(nil, lock)
	return err
}

func (d *BaetylCloudDB) DeleteLock(lock *models.Lock) error {
	res, err := d.DeleteLockTx(nil, lock)
	if res == nil || err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err == nil && rows == 1 {
		return nil
	}
	return common.Error(common.ErrResourceNotFound, common.Field("type", "lock"), common.Field("name", lock.Name))
}

func (d *BaetylCloudDB) DeleteExpiredLock(lock *models.Lock) error {
	res, err := d.DeleteExpiredLockTx(nil, lock)
	if res == nil || err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err == nil && rows == 1 {
		return nil
	}
	return common.Error(common.ErrResourceNotFound, common.Field("type", "lock"), common.Field("name", lock.Name))
}

func (d *BaetylCloudDB) InsertLockTx(tx *sqlx.Tx, lock *models.Lock) (sql.Result, error) {
	insertSQL := `INSERT INTO baetyl_lock (name, version, expire) VALUES (?, ?, ?)`
	return d.Exec(tx, insertSQL, lock.Name, lock.Version, lock.TTL)
}

func (d *BaetylCloudDB) DeleteLockTx(tx *sqlx.Tx, lock *models.Lock) (sql.Result, error) {
	deleteSQL := `DELETE FROM baetyl_lock WHERE name=? AND version=?`
	return d.Exec(tx, deleteSQL, lock.Name, lock.Version)
}

func (d *BaetylCloudDB) DeleteExpiredLockTx(tx *sqlx.Tx, lock *models.Lock) (sql.Result, error) {
	deleteSQL := `DELETE FROM baetyl_lock WHERE name=? AND NOW() > DATE_ADD(create_time, INTERVAL expire SECOND)`
	return d.Exec(tx, deleteSQL, lock.Name)
}
