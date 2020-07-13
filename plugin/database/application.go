package database

import (
	"database/sql"
	"github.com/baetyl/baetyl-cloud/plugin/database/entities"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/jmoiron/sqlx"
)

func (d *dbStorage) CreateApplication(app *specV1.Application) (sql.Result, error) {
	return d.CreateApplicationWithTx(nil, app)
}

func (d *dbStorage) UpdateApplication(app *specV1.Application, oldVersion string) (sql.Result, error) {
	return d.UpdateApplicationWithTx(nil, app, oldVersion)
}

func (d *dbStorage) DeleteApplication(name, namespace, version string) (sql.Result, error) {
	return d.DeleteApplicationWithTx(nil, name, namespace, version)
}

func (d *dbStorage) GetApplication(name, namespace, version string) (*specV1.Application, error) {
	selectSQL := `
SELECT  
id, namespace, name, version, is_deleted, create_time, update_time, content
FROM baetyl_application_history 
WHERE namespace = ? AND name=? AND version = ? AND is_deleted = 0
`
	var apps []entities.Application
	if err := d.query(nil, selectSQL, &apps, namespace, name, version); err != nil {
		return nil, err
	}
	if len(apps) > 0 {
		return entities.ToApplicationModel(&apps[0])
	}
	return nil, nil
}

func (d *dbStorage) ListApplication(name, namespace string, pageNo, pageSize int) ([]specV1.Application, error) {
	selectSQL := `
SELECT  
id, namespace, name, version, is_deleted, create_time, update_time, content
FROM baetyl_application_history WHERE namespace = ? AND name = ? AND is_deleted = 0 LIMIT ?,?
`
	var apps []entities.Application
	if err := d.query(nil, selectSQL, &apps, namespace, name, (pageNo-1)*pageSize, pageSize); err != nil {
		return nil, err
	}
	var result []specV1.Application
	for _, app := range apps {
		application, err := entities.ToApplicationModel(&app)
		if err != nil {
			return nil, err
		}
		result = append(result, *application)
	}
	return result, nil
}

func (d *dbStorage) CreateApplicationWithTx(tx *sqlx.Tx, app *specV1.Application) (sql.Result, error) {
	insertSQL := `
INSERT INTO baetyl_application_history 
(namespace, name, version, content) 
VALUES (?, ?, ?, ?)
`
	application, err := entities.FromApplicationModel(app)
	if err != nil {
		return nil, err
	}
	return d.exec(tx, insertSQL, application.Namespace, application.Name, application.Version, application.Content)
}

func (d *dbStorage) UpdateApplicationWithTx(tx *sqlx.Tx, app *specV1.Application, oldVersion string) (sql.Result, error) {
	updateSQL := `
UPDATE baetyl_application_history SET namespace = ?, name = ?, version = ?, content = ?
WHERE namespace = ? AND name = ? AND version = ?  
`
	application, err := entities.FromApplicationModel(app)
	if err != nil {
		return nil, err
	}
	return d.exec(tx, updateSQL, application.Namespace, application.Name, application.Version, application.Content,
		app.Namespace, app.Name, oldVersion)
}

func (d *dbStorage) DeleteApplicationWithTx(tx *sqlx.Tx, name, namespace, version string) (sql.Result, error) {
	deleteSQL := `
UPDATE baetyl_application_history 
SET is_deleted = 1
where namespace=? AND name=? AND version=?
`
	return d.exec(tx, deleteSQL, namespace, name, version)
}

func (d *dbStorage) CountApplication(tx *sqlx.Tx, name, namespace string) (int, error) {
	selectSQL := `
SELECT count(name) AS count
FROM baetyl_application_history WHERE namespace=? AND name=?
`
	var res []struct {
		Count int `db:"count"`
	}
	if err := d.query(tx, selectSQL, &res, namespace, name); err != nil {
		return 0, err
	}
	return res[0].Count, nil
}
