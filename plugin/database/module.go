package database

import (
	"github.com/jmoiron/sqlx"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin/database/entities"
)

func (d *DB) GetModule(name string) (*models.Module, error) {
	return d.GetModuleTx(nil, name)
}

func (d *DB) GetModuleByVersion(name, version string) (*models.Module, error) {
	return d.GetModuleByVersionTx(nil, name, version)
}

func (d *DB) CreateModule(module *models.Module) error {
	return d.CreateModuleTx(nil, module)
}

func (d *DB) UpdateModule(module *models.Module) error {
	return d.UpdateModuleTx(nil, module)
}

func (d *DB) DeleteModule(name string) error {
	return d.DeleteModuleTx(nil, name)
}

func (d *DB) ListModule(filter *models.Filter) ([]models.Module, error) {
	return d.ListModuleTx(nil, filter)
}

func (d *DB) ListModuleWithOptions(tp string, hidden bool, filter *models.Filter) ([]models.Module, error) {
	return d.ListModuleWithOptionsTx(nil, tp, hidden, filter)
}

func (d *DB) GetModuleTx(tx *sqlx.Tx, name string) (*models.Module, error) {
	selectSQL := `
SELECT  
id, name, image, programs, version, type, is_hidden, description, create_time, update_time
FROM baetyl_module 
WHERE name=?
`
	return d.getModuleTx(tx, selectSQL, name)
}

func (d *DB) GetModuleByVersionTx(tx *sqlx.Tx, name, version string) (*models.Module, error) {
	selectSQL := `
SELECT  
id, name, image, programs, version, type, is_hidden, description, create_time, update_time
FROM baetyl_module 
WHERE name=? AND version=?
`
	return d.getModuleTx(tx, selectSQL, name, version)
}

func (d *DB) CreateModuleTx(tx *sqlx.Tx, module *models.Module) error {
	insertSQL := `
INSERT INTO baetyl_module (name, image, programs, version, type, is_hidden, description)
VALUES (?, ?, ?, ?, ?, ?, ?)
`
	res, err := entities.FromModuleModel(module)
	if err != nil {
		return err
	}
	_, err = d.Exec(tx, insertSQL, res.Name, res.Image, res.Programs, res.Version, res.Type, res.IsHidden, res.Description)
	return err
}

func (d *DB) UpdateModuleTx(tx *sqlx.Tx, module *models.Module) error {
	updateSQL := `
UPDATE baetyl_module
SET image=?, programs=?, version=?, type=?, is_hidden=?, description=? 
WHERE name=?
`
	res, err := entities.FromModuleModel(module)
	if err != nil {
		return err
	}
	_, err = d.Exec(tx, updateSQL, res.Image, res.Programs, res.Version, res.Type, res.IsHidden, res.Description, res.Name)
	return err
}

func (d *DB) DeleteModuleTx(tx *sqlx.Tx, name string) error {
	deleteSQL := `
DELETE FROM baetyl_module 
WHERE name=?
	`
	_, err := d.Exec(tx, deleteSQL, name)
	return err
}

func (d *DB) ListModuleTx(tx *sqlx.Tx, filter *models.Filter) ([]models.Module, error) {
	selectSQL := `
SELECT 
id, name, image, programs, version, type, is_hidden, description, create_time, update_time
FROM baetyl_module WHERE name LIKE ? ORDER BY create_time DESC
`
	args := []interface{}{filter.GetFuzzyName()}
	if filter.GetLimitNumber() > 0 {
		selectSQL = selectSQL + "LIMIT ?,?"
		args = append(args, filter.GetLimitOffset(), filter.GetLimitNumber())
	}

	return d.listModuleWithOptionsTx(tx, selectSQL, args...)
}

func (d *DB) ListModuleWithOptionsTx(tx *sqlx.Tx, tp string, hidden bool, filter *models.Filter) ([]models.Module, error) {
	selectSQL := `
SELECT 
id, name, image, programs, version, type, is_hidden, description, create_time, update_time
FROM baetyl_module WHERE name LIKE ? AND type=? AND is_hidden=? ORDER BY create_time DESC
`

	args := []interface{}{filter.GetFuzzyName(), tp, hidden}
	if filter.GetLimitNumber() > 0 {
		selectSQL = selectSQL + "LIMIT ?,?"
		args = append(args, filter.GetLimitOffset(), filter.GetLimitNumber())
	}

	return d.listModuleWithOptionsTx(tx, selectSQL, args...)
}

func (d *DB) getModuleTx(tx *sqlx.Tx, sql string, args ...interface{}) (*models.Module, error) {
	var modules []entities.Module
	if err := d.Query(tx, sql, &modules, args...); err != nil {
		return nil, err
	}
	if len(modules) > 0 {
		return entities.ToModuleModel(&modules[0])
	}
	return nil, common.Error(
		common.ErrResourceNotFound,
		common.Field("type", "module"),
		common.Field("args", args))
}

func (d *DB) listModuleWithOptionsTx(tx *sqlx.Tx, sql string, args ...interface{}) ([]models.Module, error) {
	ms := make([]entities.Module, 0)
	if err := d.Query(tx, sql, &ms, args...); err != nil {
		return nil, err
	}

	result := make([]models.Module, 0)
	for _, module := range ms {
		m, err := entities.ToModuleModel(&module)
		if err != nil {
			return nil, err
		}
		result = append(result, *m)
	}
	return result, nil
}
