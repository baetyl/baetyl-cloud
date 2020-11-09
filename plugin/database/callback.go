package database

import (
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin/database/entities"
)

func (d *DB) GetCallback(name, namespace string) (*models.Callback, error) {
	return d.GetCallbackTx(nil, name, namespace)
}

func (d *DB) CreateCallback(callback *models.Callback) (sql.Result, error) {
	return d.CreateCallbackTx(nil, callback)
}

func (d *DB) UpdateCallback(callback *models.Callback) (sql.Result, error) {
	return d.UpdateCallbackTx(nil, callback)
}

func (d *DB) DeleteCallback(name, ns string) (sql.Result, error) {
	return d.DeleteCallbackTx(nil, name, ns)
}

func (d *DB) GetCallbackTx(tx *sqlx.Tx, name, namespace string) (*models.Callback, error) {
	selectSQL := `
SELECT name, namespace, method, params, 
header, body, url, description, create_time, 
update_time 
FROM baetyl_callback 
WHERE namespace=? and name=? LIMIT 0,1
`
	var callback []entities.Callback
	if err := d.Query(tx, selectSQL, &callback, namespace, name); err != nil {
		return nil, err
	}
	if len(callback) > 0 {
		return entities.ToCallbackModel(&callback[0]), nil
	}
	return nil, nil
}

func (d *DB) CreateCallbackTx(tx *sqlx.Tx, callback *models.Callback) (sql.Result, error) {
	insertSQL := `
INSERT INTO baetyl_callback (
name, namespace, method, params, 
header, body, url, description) 
VALUES (?,?,?,?,?,?,?,?)
`
	callbackDB := entities.FromCallbackModel(callback)
	return d.Exec(tx, insertSQL, callbackDB.Name,
		callbackDB.Namespace, callbackDB.Method, callbackDB.Params,
		callbackDB.Header, callbackDB.Body, callbackDB.Url, callbackDB.Description)
}

func (d *DB) UpdateCallbackTx(tx *sqlx.Tx, callback *models.Callback) (sql.Result, error) {
	updateSQL := `
UPDATE baetyl_callback SET method=?,params=?,
header=?,body=?,url=?,description=? 
WHERE namespace=? AND name=?
`
	callbackDB := entities.FromCallbackModel(callback)
	return d.Exec(tx, updateSQL, callbackDB.Method, callbackDB.Params,
		callbackDB.Header, callbackDB.Body, callbackDB.Url, callbackDB.Description,
		callbackDB.Namespace, callbackDB.Name)
}

func (d *DB) DeleteCallbackTx(tx *sqlx.Tx, name, ns string) (sql.Result, error) {
	deleteSQL := `
DELETE FROM baetyl_callback where namespace=? AND name=?
`
	return d.Exec(tx, deleteSQL, ns, name)
}
