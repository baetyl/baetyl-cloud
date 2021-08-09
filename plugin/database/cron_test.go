package database

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
)

var (
	cronAppTables = []string{
		`
CREATE TABLE baetyl_cron_app(
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        VARCHAR(128) NOT NULL DEFAULT '',
    namespace   VARCHAR(64) NOT NULL DEFAULT '',
	selector    VARCHAR(2048) NOT NULL DEFAULT '',
	cron_time   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`,
	}
)

func (d *DB) MockCreateCronAppTable() {
	for _, sql := range cronAppTables {
		_, err := d.Exec(nil, sql)
		if err != nil {
			panic(fmt.Sprintf("create table exception: %s", err.Error()))
		}
	}
}

func TestCronApp(t *testing.T) {
	db, err := MockNewDB()
	if err != nil {
		fmt.Printf("get mock sqlite3 error = %s", err.Error())
		t.Fail()
		return
	}
	db.MockCreateCronAppTable()

	name, ns, selector := "baetyl", "cloud", "baetyl-node-name=node1"
	cronApp := &models.Cron{
		Name: name,
		Namespace: ns,
		Selector: selector,
		CronTime: time.Now(),
	}

	err = db.CreateCron(cronApp)
	assert.NoError(t, err)

	_, err = db.GetCron(name, ns)
	assert.NoError(t, err)

	cronApp.Selector = "baetyl-node-name=node2"
	err = db.UpdateCron(cronApp)
	assert.NoError(t, err)

	_, err = db.ListExpiredApps()
	assert.NotEqual(t, err, nil)

	err = db.DeleteExpiredApps([]uint64{1})
	assert.NoError(t, err)

	err = db.DeleteCron(name, ns)
	assert.NoError(t, err)

	_, err = db.GetCron(name, ns)
	assert.Error(t, err, common.ErrResourceNotFound)
}