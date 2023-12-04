// Package database 数据库存储实现
package database

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
)

var (
	lockTable = []string{
		`
CREATE TABLE baetyl_lock
(
	id                integer       PRIMARY KEY AUTOINCREMENT,
    name              varchar(128)  NOT NULL DEFAULT '',
	version           varchar(128)  NOT NULL DEFAULT '',
    expire            Integer       NOT NULL DEFAULT 0,
    create_time       timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time       timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`,
	}
)

func (d *BaetylCloudDB) MockCreateLockTable() {
	for _, sql := range lockTable {
		_, err := d.Exec(nil, sql)
		if err != nil {
			panic(fmt.Sprintf("create lock exception: %s", err.Error()))
		}
	}
}

func TestLock(t *testing.T) {
	locker1 := &models.Lock{Name: "1", Version: "1", TTL: 0}

	db, err := MockNewDB()
	if err != nil {
		fmt.Printf("get mock sqlite3 error = %s", err.Error())
		t.Fail()
		return
	}
	db.MockCreateLockTable()

	err = db.InsertLock(locker1)
	assert.NoError(t, err)

	err = db.DeleteLock(locker1)
	assert.NoError(t, err)

	err = db.DeleteLock(locker1)
	assert.Error(t, err, common.ErrResourceNotFound)
}
