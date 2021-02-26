package database

import (
	"fmt"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin/database/entities"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	taskTables = []string{
		`
CREATE TABLE baetyl_task
(
	id                INTEGER PRIMARY KEY AUTOINCREMENT,
    name              VARCHAR(128) NOT NULL DEFAULT '',
    namespace         VARCHAR(64) NOT NULL DEFAULT '',
    registration_name VARCHAR(32) NOT NULL DEFAULT '',
    resource_type     VARCHAR(32) NOT NULL DEFAULT '',
    resource_name     VARCHAR(128) NOT NULL DEFAULT '',
    version           INTEGER NOT NULL DEFAULT 0,
    expire_time       INTEGER  NOT NULL DEFAULT 0,
    status            INTEGER  NOT NULL DEFAULT 0,
    content           VARCHAR(1024)   NOT NULL DEFAULT '',
    create_time       TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time       TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`,
	}
)

func (d *DB) MockCreateTaskTable() {
	for _, sql := range taskTables {
		_, err := d.Exec(nil, sql)
		if err != nil {
			panic(fmt.Sprintf("create table exception: %s", err.Error()))
		}
	}
}

func TestTask(t *testing.T) {
	db, err := MockNewDB()
	if err != nil {
		fmt.Printf("get mock sqlite3 error = %s", err.Error())
		t.Fail()
		return
	}
	db.MockCreateTaskTable()

	task := &entities.Task{
		Name:             "task01",
		Namespace:        "default",
		RegistrationName: "delete_test",
		ResourceType:     "namespace",
		ResourceName:     "node01",
	}

	mTask := &models.Task{
		Name:             "task02",
		Namespace:        "default",
		RegistrationName: "delete_test",
		ResourceType:     "namespace",
		ResourceName:     "node02",
	}

	err = db.Transact(func(tx *sqlx.Tx) error {
		res, ierr := db.CreateTaskTx(tx, task)
		assert.NoError(t, ierr)
		num, ierr := res.RowsAffected()
		assert.NoError(t, ierr)
		assert.Equal(t, int64(1), num)

		tk, ierr := db.GetTaskTx(tx, task.Name)
		assert.NoError(t, ierr)
		assert.Equal(t, task.Name, tk.Name)
		assert.Equal(t, task.Namespace, tk.Namespace)
		assert.Equal(t, task.RegistrationName, tk.RegistrationName)
		assert.Equal(t, 0, tk.Status)

		task.Version = tk.Version
		task.Id = tk.Id
		res, ierr = db.AcquireTaskLockTx(tx, task)
		assert.NoError(t, ierr)
		num, ierr = res.RowsAffected()
		assert.NoError(t, ierr)
		assert.Equal(t, int64(1), num)

		tk, ierr = db.GetTaskTx(tx, task.Name)
		assert.NoError(t, ierr)
		task.Status = 3
		task.Version = tk.Version
		res, ierr = db.UpdateTaskTx(tx, task)
		assert.NoError(t, ierr)
		num, ierr = res.RowsAffected()
		assert.NoError(t, ierr)
		assert.Equal(t, int64(1), num)

		res, ierr = db.DeleteTaskTx(tx, task.Name)
		assert.NoError(t, ierr)
		num, ierr = res.RowsAffected()
		assert.NoError(t, ierr)
		assert.Equal(t, int64(1), num)

		return nil
	})
	assert.NoError(t, err)

	res, err := db.CreateTask(mTask)
	assert.NoError(t, err)
	assert.True(t, res)

	tk, err := db.GetTask(mTask.Name)
	assert.NoError(t, err)
	assert.Equal(t, mTask.Name, tk.Name)
	assert.Equal(t, mTask.Namespace, tk.Namespace)
	assert.Equal(t, mTask.RegistrationName, tk.RegistrationName)
	assert.Equal(t, 0, tk.Status)

	mTask.Id = tk.Id
	mTask.Version = tk.Version
	res, err = db.UpdateTask(mTask)
	assert.NoError(t, err)
	assert.True(t, res)

	tk, err = db.GetTask(mTask.Name)
	assert.NoError(t, err)
	mTask.Version = tk.Version

	//tasks, err := db.GetNeedProcessTask(10, 10)
	//assert.NoError(t, err)
	//assert.Equal(t, 1, len(tasks))

	res, err = db.AcquireTaskLock(mTask)
	assert.NoError(t, err)
	assert.True(t, res)

	res, err = db.DeleteTask(mTask.Name)
	assert.NoError(t, err)
	assert.True(t, res)
}
