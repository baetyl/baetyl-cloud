package database

import (
	"database/sql"

	"github.com/baetyl/baetyl-cloud/common"
	"github.com/baetyl/baetyl-cloud/models"
	"github.com/jmoiron/sqlx"
)

func (d *dbStorage) CreateSystemConfig(sysConfig *models.SystemConfig) (sql.Result, error) {
	return d.CreateSystemConfigTx(nil, sysConfig)
}
func (d *dbStorage) DeleteSystemConfig(key string) (sql.Result, error) {
	return d.DeleteSystemConfigTx(nil, key)
}
func (d *dbStorage) GetSystemConfig(key string) (*models.SystemConfig, error) {
	return d.GetSystemConfigTx(nil, key)
}
func (d *dbStorage) ListSystemConfig(key string, page, size int) ([]models.SystemConfig, error) {
	return d.ListSystemConfigTx(nil, key, page, size)
}
func (d *dbStorage) CountSystemConfig(key string)(int, error){
	return d.CountSystemConfigTx(nil, key)
}
func (d *dbStorage) UpdateSystemConfig(sysConfig *models.SystemConfig) (sql.Result, error) {
	return d.UpdateSystemConfigTx(nil, sysConfig)
}

func (d *dbStorage) CreateSystemConfigTx(tx *sqlx.Tx, sysConfig *models.SystemConfig) (sql.Result, error) {
	insertSQL := "INSERT INTO baetyl_cloud_system_config(`key`, value) VALUES(?,?)"
	return d.exec(tx, insertSQL, sysConfig.Key, sysConfig.Value)
}

func (d *dbStorage) DeleteSystemConfigTx(tx *sqlx.Tx, key string) (sql.Result, error) {
	deleteSQL := "DELETE FROM baetyl_cloud_system_config " +
		"where `key`=?"
	return d.exec(tx, deleteSQL, key)
}

func (d *dbStorage) GetSystemConfigTx(tx *sqlx.Tx, key string) (*models.SystemConfig, error) {
	selectSQL := "SELECT `key`,  value, create_time, update_time " +
		"FROM baetyl_cloud_system_config " +
		"WHERE `key`=?"
	var cs []models.SystemConfig
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

func (d *dbStorage) ListSystemConfigTx(tx *sqlx.Tx, key string, pageNo, pageSize int) ([]models.SystemConfig, error) {
	selectSQL := "SELECT `key`, value, create_time, update_time " +
		"FROM baetyl_cloud_system_config " +
		"WHERE `key` LIKE ?" +
		" LIMIT ?,?"
	cs := []models.SystemConfig{}
	if err := d.query(tx, selectSQL, &cs, key, (pageNo-1)*pageSize, pageSize); err != nil {
		return nil, err
	}
	return cs, nil
}

func (d *dbStorage) CountSystemConfigTx(tx *sqlx.Tx, key string)(int, error){
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

func (d *dbStorage) UpdateSystemConfigTx(tx *sqlx.Tx, sysConfig *models.SystemConfig) (sql.Result, error) {
	updateSQL := "UPDATE baetyl_cloud_system_config SET value=? " +
		"WHERE `key`=?"
	return d.exec(tx, updateSQL, sysConfig.Value, sysConfig.Key)
}
