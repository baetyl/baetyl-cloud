package database

import (
	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
)

func (d *dbStorage) CreateProperty(property *models.Property) error {
	insertSQL := "INSERT INTO baetyl_property(`key`, value) VALUES(?,?)"
	_, err := d.exec(nil, insertSQL, property.Key, property.Value)
	if err != nil {
		return common.Error(common.ErrDatabase, common.Field("error", "key existed"))
	}
	return nil
}

func (d *dbStorage) DeleteProperty(key string) error {
	deleteSQL := "DELETE FROM baetyl_property WHERE `key`=?"
	_, err := d.exec(nil, deleteSQL, key)
	return err
}

func (d *dbStorage) GetPropertyValue(key string) (string, error) {
	selectSQL := "SELECT `key`, value, create_time, update_time FROM baetyl_property WHERE `key`=?"
	var cs []models.Property
	if err := d.query(nil, selectSQL, &cs, key); err != nil {
		return "", err
	}
	if len(cs) == 1 {
		return (&cs[0]).Value, nil
	}
	return "", common.Error(
		common.ErrResourceNotFound,
		common.Field("key", key))
}

func (d *dbStorage) ListProperty(page *models.Filter) ([]models.Property, error) {
	selectSQL := "SELECT `key`, value, create_time, update_time FROM baetyl_property WHERE `key` LIKE ? LIMIT ?,?"
	cs := []models.Property{}
	if err := d.query(nil, selectSQL, &cs, page.Name, (page.PageNo-1)*page.PageSize, page.PageSize); err != nil {
		return nil, err
	}
	return cs, nil
}

func (d *dbStorage) CountProperty(key string) (int, error) {
	selectSQL := "SELECT count(`key`) AS count FROM baetyl_property WHERE `key` LIKE ?"
	var res []struct {
		Count int `db:"count"`
	}
	if err := d.query(nil, selectSQL, &res, key); err != nil {
		return 0, err
	}
	return res[0].Count, nil
}

func (d *dbStorage) UpdateProperty(property *models.Property) error {
	updateSQL := "UPDATE baetyl_property SET value=? WHERE `key`=?"
	_, err := d.exec(nil, updateSQL, property.Value, property.Key)
	return err
}
