package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin/database/entities"
)

func (d *DB) GetModules(name string) ([]models.Module, error) {
	return d.GetModuleTx(nil, name)
}

func (d *DB) GetModuleByVersion(name, version string) (*models.Module, error) {
	return d.GetModuleByVersionTx(nil, name, version)
}

func (d *DB) GetModuleByImage(name, image string) (*models.Module, error) {
	return d.GetModuleByImageTx(nil, name, image)
}

func (d *DB) GetLatestModule(name string) (*models.Module, error) {
	return d.GetLatestModuleTx(nil, name)
}

func (d *DB) CreateModule(module *models.Module) (*models.Module, error) {
	var res *models.Module
	err := d.Transact(func(tx *sqlx.Tx) error {
		if module.IsLatest {
			err := d.disableLatestModuleTx(tx, module.Name)
			if err != nil {
				return err
			}
		}
		err := d.CreateModuleTx(tx, module)
		if err != nil {
			return err
		}
		res, err = d.GetModuleByVersionTx(tx, module.Name, module.Version)
		return err
	})
	return res, err
}

func (d *DB) UpdateModuleByVersion(module *models.Module) (*models.Module, error) {
	var res *models.Module
	err := d.Transact(func(tx *sqlx.Tx) error {
		if module.IsLatest {
			err := d.disableLatestModuleTx(tx, module.Name)
			if err != nil {
				return err
			}
		}
		err := d.UpdateModuleByVersionTx(tx, module)
		if err != nil {
			return err
		}
		res, err = d.GetModuleByVersionTx(tx, module.Name, module.Version)
		return err
	})
	return res, err
}

func (d *DB) DeleteModules(name string) error {
	return d.DeleteModulesTx(nil, name)
}

func (d *DB) DeleteModuleByVersion(name, version string) error {
	return d.DeleteModuleByVersionTx(nil, name, version)
}

func (d *DB) ListModules(filter *models.Filter) ([]models.Module, error) {
	return d.ListModulesTx(nil, filter)
}

func (d *DB) ListOptionalSysModules(filter *models.Filter) ([]models.Module, error) {
	return d.ListOptionalSysModulesTx(nil, filter)
}

func (d *DB) ListRuntimeModules(filter *models.Filter) ([]models.Module, error) {
	return d.ListRuntimeModulesTx(nil, filter)
}

func (d *DB) GetLatestModuleImage(name string) (string, error) {
	return d.GetLatestModuleImageTx(nil, name)
}

func (d *DB) GetLatestModuleProgram(name, platform string) (string, error) {
	return d.GetLatestModuleProgramTx(nil, name, platform)
}

func (d *DB) GetModuleTx(tx *sqlx.Tx, name string) ([]models.Module, error) {
	selectSQL := `
SELECT  
id, name, image, programs, version, type, is_latest, description, create_time, update_time
FROM baetyl_module 
WHERE name=? ORDER BY create_time DESC
`
	modules, err := d.listModuleTx(tx, selectSQL, name)
	if err != nil {
		return nil, err
	}
	if len(modules) > 0 {
		return modules, nil
	}
	return nil, common.Error(
		common.ErrResourceNotFound,
		common.Field("type", "module"),
		common.Field("name", name))
}

func (d *DB) GetLatestModuleTx(tx *sqlx.Tx, name string) (*models.Module, error) {
	selectSQL := `
SELECT  
id, name, image, programs, version, type, is_latest, description, create_time, update_time
FROM baetyl_module 
WHERE name=? AND is_latest=?
`
	return d.getModuleTx(tx, selectSQL, name, true)
}

func (d *DB) GetModuleByVersionTx(tx *sqlx.Tx, name, version string) (*models.Module, error) {
	selectSQL := `
SELECT  
id, name, image, programs, version, type, is_latest, description, create_time, update_time
FROM baetyl_module 
WHERE name=? AND version=?
`
	return d.getModuleTx(tx, selectSQL, name, version)
}

func (d *DB) GetModuleByImageTx(tx *sqlx.Tx, name, image string) (*models.Module, error) {
	selectSQL := `
SELECT  
id, name, image, programs, version, type, is_latest, description, create_time, update_time
FROM baetyl_module 
WHERE name=? AND image=?
`
	return d.getModuleTx(tx, selectSQL, name, image)
}

func (d *DB) CreateModuleTx(tx *sqlx.Tx, module *models.Module) error {
	insertSQL := `
INSERT INTO baetyl_module (name, image, programs, version, type, is_latest, description)
VALUES (?, ?, ?, ?, ?, ?, ?)
`
	res, err := entities.FromModuleModel(module)
	if err != nil {
		return err
	}
	_, err = d.Exec(tx, insertSQL, res.Name, res.Image, res.Programs, res.Version, res.Type, res.IsLatest, res.Description)
	return err
}

func (d *DB) UpdateModuleByVersionTx(tx *sqlx.Tx, module *models.Module) error {
	updateSQL := `
UPDATE baetyl_module
SET image=?, programs=?, version=?, type=?, is_latest=?, description=? 
WHERE name=? AND version=?
`
	res, err := entities.FromModuleModel(module)
	if err != nil {
		return err
	}
	_, err = d.Exec(tx, updateSQL, res.Image, res.Programs, res.Version, res.Type, res.IsLatest, res.Description, res.Name, res.Version)
	return err
}

func (d *DB) DeleteModulesTx(tx *sqlx.Tx, name string) error {
	deleteSQL := `
DELETE FROM baetyl_module 
WHERE name=?
	`
	_, err := d.Exec(tx, deleteSQL, name)
	return err
}

func (d *DB) DeleteModuleByVersionTx(tx *sqlx.Tx, name, version string) error {
	deleteSQL := `
DELETE FROM baetyl_module 
WHERE name=? AND version=?
	`
	_, err := d.Exec(tx, deleteSQL, name, version)
	return err
}

func (d *DB) ListModulesTx(tx *sqlx.Tx, filter *models.Filter) ([]models.Module, error) {
	selectSQL := `
SELECT 
id, name, image, programs, version, type, is_latest, description, create_time, update_time
FROM baetyl_module WHERE name LIKE ? ORDER BY create_time DESC
`
	args := []interface{}{filter.GetFuzzyName()}
	if filter.GetLimitNumber() > 0 {
		selectSQL = selectSQL + "LIMIT ?,?"
		args = append(args, filter.GetLimitOffset(), filter.GetLimitNumber())
	}

	return d.listModuleTx(tx, selectSQL, args...)
}

func (d *DB) ListOptionalSysModulesTx(tx *sqlx.Tx, filter *models.Filter) ([]models.Module, error) {
	t := common.TypeSystemOptional
	return d.listModulesByTypeTx(tx, t, filter)
}

func (d *DB) ListRuntimeModulesTx(tx *sqlx.Tx, filter *models.Filter) ([]models.Module, error) {
	t := common.TypeUserRuntime
	return d.listModulesByTypeTx(tx, t, filter)
}

func (d *DB) GetLatestModuleImageTx(tx *sqlx.Tx, name string) (string, error) {
	module, err := d.GetLatestModuleTx(tx, name)
	if err != nil {
		return "", err
	}
	return module.Image, nil
}

func (d *DB) GetLatestModuleProgramTx(tx *sqlx.Tx, name, platform string) (string, error) {
	module, err := d.GetLatestModuleTx(tx, name)
	if err != nil {
		return "", err
	}
	for k, v := range module.Programs {
		if k == platform {
			return v, nil
		}
	}
	return "", common.Error(common.ErrResourceNotFound,
		common.Field("type", "program"),
		common.Field("name", fmt.Sprintf("%s-%s", name, platform)))
}

func (d *DB) listModulesByTypeTx(tx *sqlx.Tx, tp common.ModuleType, filter *models.Filter) ([]models.Module, error) {
	selectSQL := `
SELECT 
id, name, image, programs, version, type, is_latest, description, create_time, update_time
FROM baetyl_module WHERE name LIKE ? AND type=? AND is_latest=? ORDER BY create_time DESC
`

	args := []interface{}{filter.GetFuzzyName(), string(tp), true}
	if filter.GetLimitNumber() > 0 {
		selectSQL = selectSQL + "LIMIT ?,?"
		args = append(args, filter.GetLimitOffset(), filter.GetLimitNumber())
	}

	return d.listModuleTx(tx, selectSQL, args...)
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
		common.Field("name", args))
}

func (d *DB) listModuleTx(tx *sqlx.Tx, sql string, args ...interface{}) ([]models.Module, error) {
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

func (d *DB) disableLatestModuleTx(tx *sqlx.Tx, name string) error {
	updateSQL := `
UPDATE baetyl_module
SET is_latest=?
WHERE name=?
`
	_, err := d.Exec(tx, updateSQL, false, name)
	return err
}
