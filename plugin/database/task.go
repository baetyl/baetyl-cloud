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
	return d.GetTaskTx(name)
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

func (d *DB) GetNeedProcessTask(number int, seconds float32) ([]*models.Task, error) {
	selectSQL := `SELECT * FROM baetyl_task 
WHERE updated_time < DATE_ADD(NOW(), INTERVAL ? SECOND) AND status < 3 
limit ?`
	var tArr []*entities.Task
	var tasks []*models.Task
	if err := d.Query(nil, selectSQL, &tArr, -1*seconds, number); err != nil {
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
(task_name, namespace, resource_name, resource_type, expire_time, content)
VALUES (?,?,?,?,?,?)`

	return d.Exec(tx, insertSQL, task.TaskName, task.Namespace, task.ResourceName,
		task.ResourceType, task.ExpireTime, task.Content)
}

func (d *DB) AcquireTaskLockTx(tx *sqlx.Tx, task *entities.Task) (sql.Result, error) {
	updateSQL := `UPDATE baetyl_task SET status=?, version=version + 1, expire_time=? WHERE task_name=? and version=?`
	return d.Exec(tx, updateSQL, task.Status, task.ExpireTime, task.TaskName, task.Version)
}

func (d *DB) UpdateTaskTx(tx *sqlx.Tx, task *entities.Task) (sql.Result, error) {
	updateSQL := `UPDATE baetyl_task SET status=?,content=?,version=version + 1 WHERE task_name=? and version=?`
	return d.Exec(tx, updateSQL, task.Status, task.Content, task.TaskName, task.Version)
}

func (d *DB) DeleteTaskTx(tx *sqlx.Tx, taskName string) (sql.Result, error) {
	deleteSQL := `DELETE FROM baetyl_task WHERE task_name=?`
	return d.Exec(tx, deleteSQL, taskName)
}

func (d *DB) GetTaskTx(name string) (*models.Task, error) {
	selectSQL := `
SELECT  
id, task_name, namespace, resource_name, resource_type, version, expire_time, status, content, created_time, updated_time
FROM baetyl_task 
WHERE name=? 
`
	var task *entities.Task
	if err := d.Query(nil, selectSQL, &task, name); err != nil {
		return nil, err
	}

	return entities.ToTaskModel(task)
}

func isOperatedSuccess(result sql.Result) (bool, error) {
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}
