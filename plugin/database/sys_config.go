package database

import (
	"database/sql"

	"github.com/baetyl/baetyl-cloud/common"
	"github.com/baetyl/baetyl-cloud/models"
	"github.com/jmoiron/sqlx"
)

func (d *dbStorage) GetSysConfig(tp, key string) (*models.SysConfig, error) {
	return d.GetSysConfigTx(nil, tp, key)
}

func (d *dbStorage) ListSysConfig(tp string, page, size int) ([]models.SysConfig, error) {
	return d.ListSysConfigTx(nil, tp, page, size)
}

func (d *dbStorage) ListSysConfigAll(tp string) ([]models.SysConfig, error) {
	return d.ListSysConfigAllTx(nil, tp)
}

func (d *dbStorage) CreateSysConfig(sysConfig *models.SysConfig) (sql.Result, error) {
	return d.CreateSysConfigTx(nil, sysConfig)
}

func (d *dbStorage) UpdateSysConfig(sysConfig *models.SysConfig) (sql.Result, error) {
	return d.UpdateSysConfigTx(nil, sysConfig)
}

func (d *dbStorage) DeleteSysConfig(tp, key string) (sql.Result, error) {
	return d.DeleteSysConfigTx(nil, tp, key)
}

func (d *dbStorage) GetSysConfigTx(tx *sqlx.Tx, tp, key string) (*models.SysConfig, error) {
	selectSQL := `
SELECT  
type, name, value, create_time, update_time 
FROM baetyl_system_config WHERE type=? AND name=? LIMIT 0,1
`
	var cs []models.SysConfig
	if err := d.query(tx, selectSQL, &cs, tp, key); err != nil {
		return nil, err
	}
	if len(cs) > 0 {
		return &cs[0], nil
	}
	return nil, common.Error(
		common.ErrResourceNotFound,
		common.Field("type", tp),
		common.Field("name", key))
}

func (d *dbStorage) ListSysConfigTx(tx *sqlx.Tx, tp string, pageNo, pageSize int) ([]models.SysConfig, error) {
	selectSQL := `
SELECT  
type, name, value, create_time, update_time 
FROM baetyl_system_config WHERE type=? LIMIT ?,?
`
	cs := []models.SysConfig{}
	if err := d.query(tx, selectSQL, &cs, tp, (pageNo-1)*pageSize, pageSize); err != nil {
		return nil, err
	}
	return cs, nil
}

func (d *dbStorage) ListSysConfigAllTx(tx *sqlx.Tx, tp string) ([]models.SysConfig, error) {
	selectSQL := `
SELECT  
type, name, value, create_time, update_time 
FROM baetyl_system_config WHERE type=?
`
	cs := []models.SysConfig{}
	if err := d.query(tx, selectSQL, &cs, tp); err != nil {
		return nil, err
	}
	return cs, nil
}

func (d *dbStorage) CreateSysConfigTx(tx *sqlx.Tx, sysConfig *models.SysConfig) (sql.Result, error) {
	insertSQL := `
INSERT INTO baetyl_system_config 
(type, name, value) 
VALUES 
(?,?,?)
`
	return d.exec(tx, insertSQL, sysConfig.Type, sysConfig.Key, sysConfig.Value)
}

func (d *dbStorage) UpdateSysConfigTx(tx *sqlx.Tx, sysConfig *models.SysConfig) (sql.Result, error) {
	updateSQL := `
UPDATE baetyl_system_config SET value=?
WHERE type=? AND name=?
`
	return d.exec(tx, updateSQL, sysConfig.Value, sysConfig.Type, sysConfig.Key)
}

func (d *dbStorage) DeleteSysConfigTx(tx *sqlx.Tx, tp, key string) (sql.Result, error) {
	deleteSQL := `
DELETE FROM baetyl_system_config where type=? AND name=?
`
	return d.exec(tx, deleteSQL, tp, key)
}
