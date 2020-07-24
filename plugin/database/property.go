package database

import (
	"database/sql"
	"github.com/baetyl/baetyl-cloud/common"
	"github.com/baetyl/baetyl-cloud/models"
	"github.com/jmoiron/sqlx"
)

func (d *dbStorage) CreateProperty(property *models.Property) (*models.Property, error) {
	var pro *models.Property
	err := d.Transact(func(tx *sqlx.Tx) error {
		_, err := d.insertPropertyTx(tx, property)
		if err != nil {
			return err
		}
		pro, err = d.queryPropertyTx(tx, property.Key)
		return err
	})
	return pro, err
}

func (d *dbStorage) DeleteProperty(key string) error {
	_, err := d.deletePropertyTx(nil, key)
	return err
}

func (d *dbStorage) GetProperty(key string) (*models.Property, error) {
	return d.queryPropertyTx(nil, key)
}

func (d *dbStorage) ListProperty(page *models.Filter) ([]models.Property, int, error) {
	properties, err := d.listPropertyTx(nil, page.Name, page.PageNo, page.PageSize)
	if err != nil {
		return nil, 0, err
	}
	count, err := d.countPropertyTx(nil, page.Name)
	if err != nil {
		return nil, 0, err
	}
	return properties, count, nil
}

func (d *dbStorage) UpdateProperty(property *models.Property) (*models.Property, error) {
	var pro *models.Property
	err := d.Transact(func(tx *sqlx.Tx) error {
		_, err := d.updatePropertyTx(tx, property)
		if err != nil {
			return err
		}
		pro, err = d.queryPropertyTx(tx, property.Key)
		return err
	})
	return pro, err
}

func (d *dbStorage) insertPropertyTx(tx *sqlx.Tx, property *models.Property) (sql.Result, error) {
	insertSQL := "INSERT INTO baetyl_property(`key`, value) VALUES(?,?)"
	_, err := d.exec(tx, insertSQL, property.Key, property.Value)
	if err != nil {
		return nil, common.Error(common.ErrDatabase, common.Field("error", "key existed"))
	}
	return nil, err
}

func (d *dbStorage) deletePropertyTx(tx *sqlx.Tx, key string) (sql.Result, error) {
	deleteSQL := "DELETE FROM baetyl_property " +
		"where `key`=?"
	return d.exec(tx, deleteSQL, key)
}

func (d *dbStorage) queryPropertyTx(tx *sqlx.Tx, key string) (*models.Property, error) {
	selectSQL := "SELECT `key`, value, create_time, update_time " +
		"FROM baetyl_property " +
		"WHERE `key`=?"

	var cs []models.Property
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

func (d *dbStorage) listPropertyTx(tx *sqlx.Tx, key string, pageNo, pageSize int) ([]models.Property, error) {
	selectSQL := "SELECT `key`, value, create_time, update_time " +
		"FROM baetyl_property " +
		"WHERE `key` LIKE ?" +
		" LIMIT ?,?"
	cs := []models.Property{}
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

func (d *dbStorage) updatePropertyTx(tx *sqlx.Tx, property *models.Property) (sql.Result, error) {
	updateSQL := "UPDATE baetyl_property SET value=?  WHERE `key`=?"
	res, err := d.exec(tx, updateSQL, property.Value, property.Key)
	return res, err
}
