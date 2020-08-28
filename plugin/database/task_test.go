package database

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

var (
	taskTables = []string{
		`
CREATE TABLE baetyl_task
(
    trace_id          varchar(36)  NOT NULL DEFAULT '' PRIMARY KEY,
    namespace         varchar(64)  NOT NULL DEFAULT '',
    node              varchar(128)   NOT NULL DEFAULT '',
    type              varchar(32)  NOT NULL DEFAULT '',
    state             varchar(16)       NOT NULL DEFAULT '0',
    step              text   NOT NULL,
    old_version       varchar(36)   NOT NULL DEFAULT '',
    new_version       varchar(36)     NOT NULL DEFAULT '',
    create_time       timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time       timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`,
	}
)

func (d *dbStorage) MockCreateTaskTable() {
	for _, sql := range taskTables {
		_, err := d.exec(nil, sql)
		if err != nil {
			panic(fmt.Sprintf("create table exception: %s", err.Error()))
		}
	}
}

func TestTask(t *testing.T) {
	task := &models.Task{
		TraceId:    "d6cb4c5e2b9611eaa104186590da6863",
		Namespace:  "default",
		Node:       "test node",
		Type:       "APP",
		State:      "1",
		Step:       "2",
		OldVersion: "123",
		NewVersion: "345",
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	db, err := MockNewDB()
	if err != nil {
		fmt.Printf("get mock sqlite3 error = %s", err.Error())
		t.Fail()
		return
	}
	db.MockCreateTaskTable()
	res, err := db.CreateTask(task)
	assert.NoError(t, err)
	num, err := res.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), num)

	resTask, err := db.GetTask(task.TraceId)
	assert.NoError(t, err)
	checkTask(t, task, resTask)

	task.Step = "3"
	res, err = db.UpdateTask(task)
	assert.NoError(t, err)
	num, err = res.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), num)
	resTask, err = db.GetTask(task.TraceId)
	assert.NoError(t, err)
	checkTask(t, task, resTask)

	taskNum, err := db.CountTask(task)
	assert.NoError(t, err)
	assert.Equal(t, 1, taskNum)

	res, err = db.DeleteTask(task.TraceId)
	assert.NoError(t, err)
	num, err = res.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), num)
}

func checkTask(t *testing.T, expect, actual *models.Task) {
	assert.Equal(t, expect.TraceId, actual.TraceId)
	assert.Equal(t, expect.Namespace, actual.Namespace)
	assert.Equal(t, expect.Node, actual.Node)
	assert.Equal(t, expect.Type, actual.Type)
	assert.Equal(t, expect.State, actual.State)
	assert.Equal(t, expect.Step, actual.Step)
	assert.Equal(t, expect.OldVersion, actual.OldVersion)
	assert.Equal(t, expect.NewVersion, actual.NewVersion)
}
