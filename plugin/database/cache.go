package database

import (
	"database/sql"

	"github.com/baetyl/baetyl-cloud/common"
	"github.com/baetyl/baetyl-cloud/models"
	"github.com/jmoiron/sqlx"
)

func (d *dbStorage) AddCache(key ,value string) (sql.Result, error) {
	return d.AddCacheTx(nil, key,value)
}
func (d *dbStorage) DeleteCache(key string) (sql.Result, error) {
	return d.DeleteCacheTx(nil, key)
}
func (d *dbStorage) GetCache(key string) (*models.Cache, error) {
	return d.GetCacheTx(nil, key)
}
func (d *dbStorage) ListCache(key string, page, size int) ([]models.Cache, error) {
	return d.ListCacheTx(nil, key, page, size)
}
func (d *dbStorage) CountCache(key string)(int, error){
	return d.CountCacheTx(nil, key)
}
func (d *dbStorage) ReplaceCache(key ,value string) (sql.Result, error) {
	return d.ReplaceCacheTx(nil, key ,value)
}

func (d *dbStorage) AddCacheTx(tx *sqlx.Tx, key,value string) (sql.Result, error) {
	insertSQL := "INSERT INTO baetyl_cloud_system_config(`key`, value) VALUES(?,?)"
	return d.exec(tx, insertSQL, key,value)
}

func (d *dbStorage) DeleteCacheTx(tx *sqlx.Tx, key string) (sql.Result, error) {
	deleteSQL := "DELETE FROM baetyl_cloud_system_config " +
		"where `key`=?"
	return d.exec(tx, deleteSQL, key)
}

func (d *dbStorage) GetCacheTx(tx *sqlx.Tx, key string) (*models.Cache, error) {
	selectSQL := "SELECT `key`,  value, create_time, update_time " +
		"FROM baetyl_cloud_system_config " +
		"WHERE `key`=?"

	var cs []models.Cache
	if err := d.query(tx, selectSQL, &cs, key); err != nil {
		return nil, err
	}
	if len(cs) > 0 {
		return &cs[0], nil
	}
	return nil, common.Error(
		common.ErrResourceNotFound,
		common.Field("key", key))
}

func (d *dbStorage) ListCacheTx(tx *sqlx.Tx, key string, pageNo, pageSize int) ([]models.Cache, error) {
	selectSQL := "SELECT `key`, value, create_time, update_time " +
		"FROM baetyl_cloud_system_config " +
		"WHERE `key` LIKE ?" +
		" LIMIT ?,?"
	cs := []models.Cache{}
	if err := d.query(tx, selectSQL, &cs, key, (pageNo-1)*pageSize, pageSize); err != nil {
		return nil, err
	}
	return cs, nil
}

func (d *dbStorage) CountCacheTx(tx *sqlx.Tx, key string)(int, error){
	selectSQL := "SELECT count(`key`) AS count " +
		"FROM baetyl_cloud_system_config " +
		"WHERE `key` LIKE ?"
	var res []struct {
		Count int `db:"count"`
	}
	if err := d.query(tx, selectSQL, &res, key); err != nil {
		return 0, err
	}
	return res[0].Count, nil
}

func (d *dbStorage) ReplaceCacheTx(tx *sqlx.Tx, key ,value string) (sql.Result, error) {
	updateSQL := "UPDATE baetyl_cloud_system_config SET value=? " +
		"WHERE `key`=?"
	return d.exec(tx, updateSQL, key ,value)
}
