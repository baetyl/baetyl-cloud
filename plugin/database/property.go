package database

import (
	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
)

func (d *dbStorage) CreateProperty(property *models.Property) error {
	insertSQL := `
		INSERT INTO baetyl_property(name, value) 
		VALUES (?,?)
	`
	_, err := d.exec(nil, insertSQL, property.Name, property.Value)
	if err != nil {
		return common.Error(common.ErrDatabase, common.Field("error", err.Error()))
	}
	return nil
}

func (d *dbStorage) DeleteProperty(name string) error {
	deleteSQL := `
		DELETE FROM baetyl_property 
		WHERE name=?
	`
	_, err := d.exec(nil, deleteSQL, name)
	return err
}

func (d *dbStorage) GetPropertyValue(name string) (string, error) {
	selectSQL := `
		SELECT 
		name, value, create_time, update_time 
			FROM baetyl_property 
		WHERE name=?
	`
	var cs []models.Property
	if err := d.query(nil, selectSQL, &cs, name); err != nil {
		return "", err
	}
	if len(cs) == 1 {
		return (&cs[0]).Value, nil
	}
	return "", common.Error(
		common.ErrResourceNotFound,
		common.Field("name", name))
}

func (d *dbStorage) ListProperty(page *models.Filter) ([]models.Property, error) {
	selectSQL := `
		SELECT 
		name, value, create_time, update_time 
			FROM baetyl_property 
		WHERE name LIKE ? 
		LIMIT ?,?
	`
	var cs []models.Property
	if err := d.query(nil, selectSQL, &cs, page.Name, (page.PageNo-1)*page.PageSize, page.PageSize); err != nil {
		return nil, err
	}
	return cs, nil
}

func (d *dbStorage) CountProperty(name string) (int, error) {
	selectSQL := `
		SELECT 
		count(name) AS count 
			FROM baetyl_property 
		WHERE name LIKE ?
	`
	var res []struct {
		Count int `db:"count"`
	}
	if err := d.query(nil, selectSQL, &res, name); err != nil {
		return 0, err
	}
	return res[0].Count, nil
}

func (d *dbStorage) UpdateProperty(property *models.Property) error {
	updateSQL := `
		UPDATE baetyl_property 
			SET value=? 
		WHERE name=?
	`
	_, err := d.exec(nil, updateSQL, property.Value, property.Name)
	return err
}
