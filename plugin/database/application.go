package database

import (
	"database/sql"

	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/jmoiron/sqlx"

	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin/database/entities"
)

func (d *DB) CreateApplicationHis(app *specV1.Application) (sql.Result, error) {
	return d.CreateApplicationHisWithTx(nil, app)
}

func (d *DB) UpdateApplicationHis(app *specV1.Application, oldVersion string) (sql.Result, error) {
	return d.UpdateApplicationHisWithTx(nil, app, oldVersion)
}

func (d *DB) DeleteApplicationHis(namespace, name, version string) (sql.Result, error) {
	return d.DeleteApplicationHisWithTx(nil, namespace, name, version)
}

func (d *DB) GetApplicationHis(namespace, name, version string) (*specV1.Application, error) {
	selectSQL := `
SELECT  
id, namespace, name, version, is_deleted, create_time, update_time, content
FROM baetyl_application_history 
WHERE namespace = ? AND name=? AND version = ? AND is_deleted = 0
`
	var apps []entities.Application
	if err := d.Query(nil, selectSQL, &apps, namespace, name, version); err != nil {
		return nil, err
	}
	if len(apps) > 0 {
		return entities.ToApplicationModel(&apps[0])
	}
	return nil, nil
}

func (d *DB) ListApplicationHis(namespace string, filter *models.Filter) ([]specV1.Application, error) {
	selectSQL := `
SELECT  
id, namespace, name, version, is_deleted, create_time, update_time, content
FROM baetyl_application_history WHERE namespace = ? AND name LIKE ? AND is_deleted = 0 
`
	var apps []entities.Application
	args := []interface{}{namespace, filter.GetFuzzyName()}
	if filter.GetLimitNumber() > 0 {
		selectSQL = selectSQL + "LIMIT ?,?"
		args = append(args, filter.GetLimitOffset(), filter.GetLimitNumber())
	}
	if err := d.Query(nil, selectSQL, &apps, args...); err != nil {
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

func (d *DB) DeleteAllAppsHis(namespace, name string) (sql.Result, error) {
	return d.DeleteAllAppsHisTx(nil, namespace, name)
}

func (d *DB) CreateApplicationHisWithTx(tx *sqlx.Tx, app *specV1.Application) (sql.Result, error) {
	insertSQL := `
INSERT INTO baetyl_application_history 
(namespace, name, version, content) 
VALUES (?, ?, ?, ?)
`
	application, err := entities.FromApplicationModel(app)
	if err != nil {
		return nil, err
	}
	return d.Exec(tx, insertSQL, application.Namespace, application.Name, application.Version, application.Content)
}

func (d *DB) UpdateApplicationHisWithTx(tx *sqlx.Tx, app *specV1.Application, oldVersion string) (sql.Result, error) {
	updateSQL := `
UPDATE baetyl_application_history SET namespace = ?, name = ?, version = ?, content = ?
WHERE namespace = ? AND name = ? AND version = ?  
`
	application, err := entities.FromApplicationModel(app)
	if err != nil {
		return nil, err
	}
	return d.Exec(tx, updateSQL, application.Namespace, application.Name, application.Version, application.Content,
		app.Namespace, app.Name, oldVersion)
}

func (d *DB) DeleteApplicationHisWithTx(tx *sqlx.Tx, namespace, name, version string) (sql.Result, error) {
	deleteSQL := `
UPDATE baetyl_application_history 
SET is_deleted = 1
where namespace=? AND name=? AND version=?
`
	return d.Exec(tx, deleteSQL, namespace, name, version)
}

func (d *DB) CountApplicationHis(tx *sqlx.Tx, namespace, name string) (int, error) {
	selectSQL := `
SELECT count(name) AS count
FROM baetyl_application_history WHERE namespace=? AND name=?
`
	var res []struct {
		Count int `db:"count"`
	}
	if err := d.Query(tx, selectSQL, &res, namespace, name); err != nil {
		return 0, err
	}
	return res[0].Count, nil
}

func (d *DB) DeleteAllAppsHisTx(tx *sqlx.Tx, namespace, name string) (sql.Result, error) {
	deleteSQL := `
DELETE FROM baetyl_application_history
WHERE namespace=? and name =?
`
	return d.Exec(tx, deleteSQL, namespace, name)
}
