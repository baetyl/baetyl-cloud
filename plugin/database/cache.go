package database

import (
	"database/sql"
	"github.com/baetyl/baetyl-cloud/common"
	"github.com/baetyl/baetyl-cloud/models"
	"github.com/jmoiron/sqlx"
)

func (d *dbStorage) GetCache(key string) (string, error) {
	cache, err := d.queryPropertyTx(nil, key)
	if err != nil {
		return "", err
	}
	return cache.Value, err
}
func (d *dbStorage) SetCache(key, value string) error {
	_, err := d.GetCache(key)
	if err != nil {
		_, err = d.insertPropertyTx(nil, key, value)
	} else {
		_, err = d.updatePropertyTx(nil, key, value)
	}
	return err
}
func (d *dbStorage) DeleteCache(key string) error {
	_, err := d.deletePropertyTx(nil, key)
	return err
}
func (d *dbStorage) ListCache(page *models.Filter) (*models.AmisListView, error) {
	caches, err := d.listPropertyTx(nil, page.Name, page.PageNo, page.PageSize)
	if err != nil {
		return nil, err
	}
	count, err := d.countPropertyTx(nil, page.Name)
	if err != nil {
		return nil, err
	}
	return &models.AmisListView{
		Status: "0",
		Msg:    "ok",
		Data: models.AmisData{
			Count: count,
			Rows:  caches,
		},
	}, nil
}

func (d *dbStorage) insertPropertyTx(tx *sqlx.Tx, key, value string) (sql.Result, error) {
	insertSQL := "INSERT INTO baetyl_property(`key`, value) VALUES(?,?)"
	return d.exec(tx, insertSQL, key, value)
}
func (d *dbStorage) queryPropertyTx(tx *sqlx.Tx, key string) (*models.Cache, error) {
	selectSQL := "SELECT `key`, value, create_time, update_time " +
		"FROM baetyl_property " +
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
func (d *dbStorage) deletePropertyTx(tx *sqlx.Tx, key string) (sql.Result, error) {
	deleteSQL := "DELETE FROM baetyl_property " +
		"where `key`=?"
	return d.exec(tx, deleteSQL, key)
}

func (d *dbStorage) listPropertyTx(tx *sqlx.Tx, key string, pageNo, pageSize int) ([]models.Cache, error) {
	selectSQL := "SELECT `key`, value, create_time, update_time " +
		"FROM baetyl_property " +
		"WHERE `key` LIKE ?" +
		" LIMIT ?,?"
	cs := []models.Cache{}
	if err := d.query(tx, selectSQL, &cs, key, (pageNo-1)*pageSize, pageSize); err != nil {
		return nil, err
	}
	return cs, nil
}

func (d *dbStorage) countPropertyTx(tx *sqlx.Tx, key string) (int, error) {
	selectSQL := "SELECT count(`key`) AS count " +
		"FROM baetyl_property " +
		"WHERE `key` LIKE ?"
	var res []struct {
		Count int `db:"count"`
	}
	if err := d.query(tx, selectSQL, &res, key); err != nil {
		return 0, err
	}
	return res[0].Count, nil
}

func (d *dbStorage) updatePropertyTx(tx *sqlx.Tx, key, value string) (sql.Result, error) {
	updateSQL := "UPDATE baetyl_property SET value=?  WHERE `key`=?"
	res, err := d.exec(tx, updateSQL, value, key)
	return res, err
}
