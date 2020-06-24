package database

import (
	"fmt"
	"github.com/baetyl/baetyl-cloud/common"
	"github.com/baetyl/baetyl-cloud/models"
	"github.com/stretchr/testify/assert"
	"path"
	"testing"
)

var (
	batchTables = []string{
		`
CREATE TABLE baetyl_batch
(
    id               integer       PRIMARY KEY AUTOINCREMENT,
    name             varchar(128)  NOT NULL DEFAULT '',
    namespace        varchar(64)   NOT NULL DEFAULT '',
    description      varchar(1024) NOT NULL DEFAULT '',
    quota_num        int(11)       NOT NULL DEFAULT '200',
    enable_whitelist int(11)       NOT NULL DEFAULT '0',
    security_type    varchar(32)   NOT NULL DEFAULT 'None',
    security_key     varchar(64)   NOT NULL DEFAULT '',
    callback_name    varchar(64)   NOT NULL DEFAULT '',
    labels           varchar(2048) NOT NULL DEFAULT '{}',
    fingerprint      varchar(1024) NOT NULL DEFAULT '{}',
    create_time      timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time      timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`,
	}
)

func genBatch() *models.Batch {
	return &models.Batch{
		Name:            "zx",
		Namespace:       "default",
		Description:     "desc",
		QuotaNum:        20,
		EnableWhitelist: 0,
		SecurityType:    common.None,
		SecurityKey:     "",
		CallbackName:    "test",
		Labels:          map[string]string{"a": "a"},
		Fingerprint: models.Fingerprint{
			Type:   common.FingerprintSN,
			SnPath: path.Join(common.DefaultSNPath, common.DefaultSNFile),
		},
	}
}

func (d *dbStorage) MockCreateBatchTable() {
	for _, sql := range batchTables {
		_, err := d.exec(nil, sql)
		if err != nil {
			panic(fmt.Sprintf("create table exception: %s", err.Error()))
		}
	}
}

func TestBatch(t *testing.T) {
	_, err := New()
	assert.NotNil(t, err)
	batch := genBatch()

	db, err := MockNewDB()
	if err != nil {
		fmt.Printf("get mock sqlite3 error = %s", err.Error())
		t.Fail()
		return
	}
	db.MockCreateBatchTable()
	res, err := db.CreateBatch(batch)
	assert.NoError(t, err)
	num, err := res.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), num)

	resBatch, err := db.GetBatch(batch.Name, batch.Namespace)
	assert.NoError(t, err)
	checkBatch(t, batch, resBatch)

	resBatchList, err := db.ListBatch(batch.Namespace, "%", 1, 10)
	assert.NoError(t, err)
	checkBatch(t, batch, &resBatchList[0])

	batch.Labels = map[string]string{"b": "b"}
	res, err = db.UpdateBatch(batch)
	assert.NoError(t, err)
	num, err = res.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), num)
	resBatch, err = db.GetBatch(batch.Name, batch.Namespace)
	assert.NoError(t, err)
	checkBatch(t, batch, resBatch)

	c1, err := db.CountBatch(batch.Namespace, "%")
	assert.NoError(t, err)
	assert.Equal(t, 1, c1)

	c2, err := db.CountBatchByCallback(batch.CallbackName, batch.Namespace)
	assert.NoError(t, err)
	assert.Equal(t, 1, c2)

	res, err = db.DeleteBatch(batch.Name, batch.Namespace)
	assert.NoError(t, err)
	num, err = res.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), num)
}

func checkBatch(t *testing.T, expect, actual *models.Batch) {
	assert.Equal(t, expect.Name, actual.Name)
	assert.Equal(t, expect.Namespace, actual.Namespace)
	assert.Equal(t, expect.Description, actual.Description)
	assert.Equal(t, expect.QuotaNum, actual.QuotaNum)
	assert.Equal(t, expect.EnableWhitelist, actual.EnableWhitelist)
	assert.Equal(t, expect.SecurityType, actual.SecurityType)
	assert.Equal(t, expect.SecurityKey, actual.SecurityKey)
	assert.Equal(t, expect.CallbackName, actual.CallbackName)
	assert.Equal(t, expect.Labels, actual.Labels)
	assert.Equal(t, expect.Fingerprint, actual.Fingerprint)
}
