package database

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

func (d *dbStorage) CreateTask(task *models.Task) (sql.Result, error) {
	return d.CreateTaskTx(nil, task)
}

func (d *dbStorage) GetTask(traceId string) (*models.Task, error) {
	return d.GetTaskTx(nil, traceId)
}

func (d *dbStorage) UpdateTask(task *models.Task) (sql.Result, error) {
	return d.UpdateTaskTx(nil, task)
}

func (d *dbStorage) DeleteTask(traceId string) (sql.Result, error) {
	return d.DeleteTaskTx(nil, traceId)
}

func (d *dbStorage) CountTask(task *models.Task) (int, error) {
	selectSQL := `SELECT count(*) AS count FROM baetyl_task where node=? and namespace=? and old_version=? and new_version=?`

	var res []struct {
		Count int `db:"count"`
	}
	if err := d.query(nil, selectSQL, &res, task.Node, task.Namespace, task.OldVersion, task.NewVersion); err != nil {
		return 0, err
	}
	return res[0].Count, nil
}

func (d *dbStorage) GetTaskTx(tx *sqlx.Tx, traceId string) (*models.Task, error) {
	selectSQL := `SELECT * FROM baetyl_task where trace_id=?`
	var tasks []models.Task
	if err := d.query(tx, selectSQL, &tasks, traceId); err != nil {
		return nil, err
	}
	if len(tasks) > 0 {
		return &tasks[0], nil
	}
	return nil, nil
}

func (d *dbStorage) CreateTaskTx(tx *sqlx.Tx, task *models.Task) (sql.Result, error) {
	insertSQL := `INSERT INTO baetyl_task
(trace_id, namespace, node, type, state, step, old_version, new_version, create_time, update_time)
VALUES (?,?,?,?,?,?,?,?,?,?)`
	return d.exec(tx, insertSQL, task.TraceId, task.Namespace, task.Node,
		task.Type, task.State, task.Step, task.OldVersion, task.NewVersion, time.Now(), time.Now())
}

func (d *dbStorage) UpdateTaskTx(tx *sqlx.Tx, task *models.Task) (sql.Result, error) {
	updateSQL := `UPDATE baetyl_task SET state=?,step=?,update_time=? WHERE trace_id=?`
	return d.exec(tx, updateSQL, task.State, task.Step, time.Now(), task.TraceId)
}

func (d *dbStorage) DeleteTaskTx(tx *sqlx.Tx, traceId string) (sql.Result, error) {
	deleteSQL := `DELETE FROM baetyl_task WHERE trace_id=?`
	return d.exec(tx, deleteSQL, traceId)
}
