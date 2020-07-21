package database

import (
	"database/sql"
	"io/ioutil"
	"strings"

	"github.com/baetyl/baetyl-cloud/common"
	"github.com/baetyl/baetyl-cloud/models"
	"github.com/jmoiron/sqlx"
)

func (d *dbStorage) GetCache(key string) (string, error){
	cache, err := d.GetCacheTx(nil, key)
	if err != nil {
		return "", err
	}
	if strings.HasPrefix(cache.Value,"file"){
		buf, err := ioutil.ReadFile(cache.Value) //mock?
		if err != nil {
			return cache.Value, nil
		}
		cache.Value = string(buf)
	}
	return cache.Value, nil
}
func (d *dbStorage) SetCache(key, value string) error{
	_, err := d.GetCache(key)
	if err != nil {
		_, err := d.AddCacheTx(nil, key, value)
		if err != nil {
			return common.Error(common.ErrDatabase, common.Field("error", err))
		}
	}else{
		_, err := d.ReplaceCacheTx(nil, key, value)
		if err != nil {
			return common.Error(common.ErrDatabase, common.Field("error", err))
		}
	}
	return nil
}
func (d *dbStorage) DeleteCache(key string) error{
	_, err := d.DeleteCacheTx(nil, key)
	return err
}
func (d *dbStorage) ListCache(page *models.Filter) (*models.ListView, error) {
	caches, err := d.ListCacheTx(nil, page.Name, page.PageNo, page.PageSize)
	if err != nil {
		return nil, common.Error(common.ErrDatabase, common.Field("error", err))
	}
	count, err := d.CountCacheTx(nil, page.Name)
	if err != nil {
		return nil, common.Error(common.ErrDatabase, common.Field("error", err))
	}
	return &models.ListView{
		Total:    count,
		PageNo:   page.PageNo,
		PageSize: page.PageSize,
		Items:    caches,
	}, nil
}

func (d *dbStorage) AddCacheTx(tx *sqlx.Tx, key,value string) (sql.Result, error) {
	insertSQL := "INSERT INTO baetyl_property(`key`, value) VALUES(?,?)"
	return d.exec(tx, insertSQL, key,value)
}
func (d *dbStorage) GetCacheTx(tx *sqlx.Tx, key string) (*models.Cache, error) {
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
func (d *dbStorage) DeleteCacheTx(tx *sqlx.Tx, key string) (sql.Result, error) {
	deleteSQL := "DELETE FROM baetyl_property " +
		"where `key`=?"
	return d.exec(tx, deleteSQL, key)
}

func (d *dbStorage) ListCacheTx(tx *sqlx.Tx, key string, pageNo, pageSize int) ([]models.Cache, error) {
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

func (d *dbStorage) CountCacheTx(tx *sqlx.Tx, key string)(int, error){
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

func (d *dbStorage) ReplaceCacheTx(tx *sqlx.Tx, key ,value string) (sql.Result, error) {
	updateSQL := "UPDATE baetyl_property SET value=?  WHERE `key`=?"
	res, err :=  d.exec(tx, updateSQL, value, key)
	return res, err
}
