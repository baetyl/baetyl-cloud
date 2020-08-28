package database

import (
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/common"
)

var (
	tables = []string{
		`
CREATE TABLE baetyl_index_application_config
(
    id          integer             PRIMARY KEY AUTOINCREMENT,
    namespace   varchar(64)         NOT NULL DEFAULT '',
    application varchar(128)        NOT NULL DEFAULT '',
    config      varchar(128)        NOT NULL DEFAULT '',
    create_time timestamp           NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time timestamp           NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`,
		`CREATE TABLE baetyl_index_application_node
(
    id          integer             PRIMARY KEY AUTOINCREMENT,
    namespace   varchar(64)         NOT NULL DEFAULT '',
    application varchar(128)        NOT NULL DEFAULT '',
    node      varchar(128)        NOT NULL DEFAULT '',
    create_time timestamp           NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time timestamp           NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`,
	}
)

func (d *dbStorage) MockCreateIndexTable() {
	for _, sql := range tables {
		_, err := d.exec(nil, sql)
		if err != nil {
			panic(fmt.Sprintf("create table exception: %s", err.Error()))
		}
	}
}

func TestGetTable(t *testing.T) {
	delete(cache, "application_config")
	delete(cache, "config_application")
	res := getTable(common.Application, common.Config)
	assert.Equal(t, "baetyl_index_application_config", res)
	assert.Equal(t, "baetyl_index_application_config", cache["application_config"])
	assert.Equal(t, "baetyl_index_application_config", cache["config_application"])
	assert.Equal(t, "baetyl_index_application_config", getTable(common.Application, common.Config))
	assert.Equal(t, "baetyl_index_application_config", getTable(common.Config, common.Application))

	delete(cache, "application_config")
	delete(cache, "config_application")
	res = getTable(common.Config, common.Application)
	assert.Equal(t, "baetyl_index_application_config", res)
	assert.Equal(t, "baetyl_index_application_config", cache["application_config"])
	assert.Equal(t, "baetyl_index_application_config", cache["config_application"])
}

func TestIndex(t *testing.T) {
	db, err := MockNewDB()
	if err != nil {
		fmt.Printf("get mock sqlite3 error = %s", err.Error())
		t.Fail()
		return
	}

	namespace := "default"

	db.MockCreateIndexTable()
	err = db.Transact(func(tx *sqlx.Tx) error {
		res, ierr := db.CreateIndexTx(tx, namespace, common.Application, common.Config, "app0", "config0")
		assert.NoError(t, ierr)
		num, ierr := res.RowsAffected()
		assert.NoError(t, ierr)
		assert.Equal(t, int64(1), num)

		arr, ierr := db.ListIndexTx(tx, namespace, common.Application, common.Config, "config0")
		assert.NoError(t, ierr)
		assert.Equal(t, 1, len(arr))
		assert.Equal(t, "app0", arr[0])

		res, ierr = db.DeleteIndexTx(tx, namespace, common.Application, common.Config, "config0")
		assert.NoError(t, ierr)
		num, ierr = res.RowsAffected()
		assert.NoError(t, ierr)
		assert.Equal(t, int64(1), num)

		arr, ierr = db.ListIndexTx(tx, namespace, common.Application, common.Config, "config0")
		assert.NoError(t, ierr)
		assert.Equal(t, 0, len(arr))

		return nil
	})
	assert.NoError(t, err)

	res, err := db.CreateIndex(namespace, common.Application, common.Config, "app0", "config0")
	assert.NoError(t, err)
	num, err := res.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), num)

	arr, err := db.ListIndex(namespace, common.Application, common.Config, "config0")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(arr))
	assert.Equal(t, "app0", arr[0])

	res, err = db.DeleteIndex(namespace, common.Application, common.Config, "config0")
	assert.NoError(t, err)
	num, err = res.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), num)

	arr, err = db.ListIndex(namespace, common.Application, common.Config, "config0")
	assert.NoError(t, err)
	assert.Equal(t, 0, len(arr))
	err = db.Close()
	assert.NoError(t, err)

	err = db.Transact(func(tx *sqlx.Tx) error {
		return fmt.Errorf("rollback test")
	})
	assert.Error(t, err, "rollback test")
}

func TestDbStorage_RefreshIndex(t *testing.T) {
	db, err := MockNewDB()
	if err != nil {
		fmt.Printf("get mock sqlite3 error = %s", err.Error())
		t.Fail()
		return
	}
	db.MockCreateIndexTable()

	namespace := "default"

	valueA := "nodeName"
	valueB := "application"
	db.RefreshIndex(namespace, common.Node, common.Application, valueA, []string{valueB})
	db.RefreshIndex(namespace, common.Application, common.Node, valueB, []string{valueA})
	db.RefreshIndex(namespace, common.Node, common.Application, valueB, []string{valueA})
	db.RefreshIndex(namespace, common.Application, common.Node, valueA, []string{valueB})
}
