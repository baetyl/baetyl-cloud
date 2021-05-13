package database

import (
	"database/sql"

	"github.com/baetyl/baetyl-cloud/v2/plugin/database/entities"
	"github.com/jmoiron/sqlx"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

func (d *DB) CreateTask(task *models.Task) (bool, error) {
	t, err := entities.FromTaskModel(task)
	if err != nil {
		return false, err
	}

	result, err := d.CreateTaskTx(nil, t)
	if err != nil {
		return false, err
	}

	return isOperatedSuccess(result)
}

func (d *DB) GetTask(name string) (*models.Task, error) {
	return d.GetTaskTx(nil, name)
}

func (d *DB) AcquireTaskLock(task *models.Task) (bool, error) {
	t, err := entities.FromTaskModel(task)
	if err != nil {
		return false, err
	}

	result, err := d.AcquireTaskLockTx(nil, t)
	if err != nil {
		return false, err
	}

	return isOperatedSuccess(result)
}

// GetNeedProcessTask only support for mysql
func (d *DB) GetNeedProcessTask(batchNum, expiredSeconds int32) ([]*models.Task, error) {
	selectSQL := `SELECT id, name, registration_name, namespace, resource_name, resource_type, version, expire_time, 
status, content, create_time, update_time
FROM baetyl_task 
WHERE update_time < DATE_ADD(NOW(), INTERVAL ? SECOND) AND status < ?
limit ?`
	var tArr []*entities.Task
	var tasks []*models.Task
	if err := d.Query(nil, selectSQL, &tArr, -1*expiredSeconds, models.TaskFinished, batchNum); err != nil {
		return nil, err
	}

	for _, t := range tArr {
		task, err := entities.ToTaskModel(t)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (d *DB) UpdateTask(task *models.Task) (bool, error) {
	t, err := entities.FromTaskModel(task)
	if err != nil {
		return false, err
	}

	result, err := d.UpdateTaskTx(nil, t)
	if err != nil {
		return false, err
	}

	return isOperatedSuccess(result)
}

func (d *DB) DeleteTask(taskName string) (bool, error) {
	result, err := d.DeleteTaskTx(nil, taskName)
	if err != nil {
		return false, err
	}

	return isOperatedSuccess(result)
}

func (d *DB) CreateTaskTx(tx *sqlx.Tx, task *entities.Task) (sql.Result, error) {
	insertSQL := `INSERT INTO baetyl_task
(name, registration_name, namespace, resource_name, resource_type, expire_time, content)
VALUES (?,?,?,?,?,?,?)`

	return d.Exec(tx, insertSQL, task.Name, task.RegistrationName, task.Namespace, task.ResourceName,
		task.ResourceType, task.ExpireTime, task.Content)
}

func (d *DB) AcquireTaskLockTx(tx *sqlx.Tx, task *entities.Task) (sql.Result, error) {
	updateSQL := `UPDATE baetyl_task SET version=version + 1, expire_time=? WHERE id=? and version=?`
	return d.Exec(tx, updateSQL, task.ExpireTime, task.Id, task.Version)
}

func (d *DB) UpdateTaskTx(tx *sqlx.Tx, task *entities.Task) (sql.Result, error) {
	updateSQL := `UPDATE baetyl_task SET status=?, content=?, version=version + 1 WHERE name=? and version=?`
	return d.Exec(tx, updateSQL, task.Status, task.Content, task.Name, task.Version)
}

func (d *DB) DeleteTaskTx(tx *sqlx.Tx, name string) (sql.Result, error) {
	deleteSQL := `DELETE FROM baetyl_task WHERE name=?`
	return d.Exec(tx, deleteSQL, name)
}

func (d *DB) GetTaskTx(tx *sqlx.Tx, name string) (*models.Task, error) {
	selectSQL := `
SELECT  
id, name, namespace, registration_name, resource_name, resource_type, version, expire_time, status, content, 
create_time, update_time
FROM baetyl_task 
WHERE name=?
`
	var task []*entities.Task
	if err := d.Query(tx, selectSQL, &task, name); err != nil {
		return nil, err
	}

	if len(task) > 0 {
		return entities.ToTaskModel(task[0])
	}

	return nil, nil
}

func isOperatedSuccess(result sql.Result) (bool, error) {
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}
