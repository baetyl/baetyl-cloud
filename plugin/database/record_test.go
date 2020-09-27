package database

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

var (
	recordTables = []string{
		`
CREATE TABLE baetyl_batch_record
(
	id                integer       PRIMARY KEY AUTOINCREMENT,
    name              varchar(128)  NOT NULL DEFAULT '',
    batch_name        varchar(128)  NOT NULL DEFAULT '',
    namespace         varchar(64)   NOT NULL DEFAULT '',
    fingerprint_value varchar(512)  NOT NULL DEFAULT '',
    active            int(1)        NOT NULL DEFAULT '0',
    node_name         varchar(64)   NOT NULL DEFAULT '',
    active_ip         varchar(64)   NOT NULL DEFAULT '0.0.0.0',
    active_time       timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    create_time       timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time       timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`,
	}
)

func (d *dbStorage) MockCreateRecordTable() {
	for _, sql := range recordTables {
		_, err := d.exec(nil, sql)
		if err != nil {
			panic(fmt.Sprintf("create table exception: %s", err.Error()))
		}
	}
}

func TestRecord(t *testing.T) {
	record := &models.Record{
		Name:             "4b9424842b9c11ea954b186590da6863",
		Namespace:        "default",
		BatchName:        "d6cb4c5e2b9611eaa104186590da6863",
		FingerprintValue: "1",
		Active:           0,
		NodeName:         "node test name",
		ActiveIP:         "127.0.0.1",
		ActiveTime:       time.Now(),
		CreateTime:       time.Now(),
		UpdateTime:       time.Now(),
	}

	db, err := MockNewDB()
	if err != nil {
		fmt.Printf("get mock sqlite3 error = %s", err.Error())
		t.Fail()
		return
	}
	db.MockCreateRecordTable()
	res, err := db.CreateRecord([]models.Record{*record})
	assert.NoError(t, err)
	num, err := res.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), num)
	count, err := db.CountRecord(record.BatchName, "%", record.Namespace)
	assert.Equal(t, 1, count)

	resRecord, err := db.GetRecord(record.BatchName, record.Name, record.Namespace)
	assert.NoError(t, err)
	checkRecord(t, record, resRecord)
	resRecord, err = db.GetRecordByFingerprint(record.BatchName, record.Namespace, record.FingerprintValue)
	assert.NoError(t, err)
	checkRecord(t, record, resRecord)

	resRecordList, err := db.ListRecord(record.BatchName, record.FingerprintValue, record.Namespace, 1, 10)
	assert.NoError(t, err)
	checkRecord(t, record, &resRecordList[0])

	records, err := db.ListRecordByBatchTx(nil, record.BatchName, record.Namespace)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(records))

	record.Active = 1
	res, err = db.UpdateRecord(record)
	assert.NoError(t, err)
	num, err = res.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), num)
	resRecord, err = db.GetRecord(record.BatchName, record.Name, record.Namespace)
	assert.NoError(t, err)
	checkRecord(t, record, resRecord)

	res, err = db.DeleteRecord(record.BatchName, record.Name, record.Namespace)
	assert.NoError(t, err)
	num, err = res.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), num)
}

func checkRecord(t *testing.T, expect, actual *models.Record) {
	assert.Equal(t, expect.Name, actual.Name)
	assert.Equal(t, expect.Namespace, actual.Namespace)
	assert.Equal(t, expect.BatchName, actual.BatchName)
	assert.Equal(t, expect.FingerprintValue, actual.FingerprintValue)
	assert.Equal(t, expect.Active, actual.Active)
	assert.Equal(t, expect.NodeName, actual.NodeName)
	assert.Equal(t, expect.ActiveIP, actual.ActiveIP)
}
