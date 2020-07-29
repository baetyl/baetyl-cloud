package database

import (
	"fmt"
	"testing"

	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/stretchr/testify/assert"
)

var (
	sysTables = []string{
		`
CREATE TABLE baetyl_system_config
(
    id               integer       PRIMARY KEY AUTOINCREMENT,
    type             varchar(128)  NOT NULL DEFAULT '',
    name             varchar(128)  NOT NULL DEFAULT '',
    value            varchar(2048) NOT NULL DEFAULT '',
    create_time      timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time      timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`,
	}
)

func genSysConfig() *models.SysConfig {
	return &models.SysConfig{
		Type:  "baetyl",
		Key:   "0.1.0",
		Value: "http://test.baetyl/0.1.0",
	}
}

func (d *dbStorage) MockCreateSysTable() {
	for _, sql := range sysTables {
		_, err := d.exec(nil, sql)
		if err != nil {
			panic(fmt.Sprintf("create table exception: %s", err.Error()))
		}
	}
}

func TestSysConfig(t *testing.T) {
	sysConfig := genSysConfig()

	db, err := MockNewDB()
	if err != nil {
		fmt.Printf("get mock sqlite3 error = %s", err.Error())
		t.Fail()
		return
	}
	db.MockCreateSysTable()
	res, err := db.CreateSysConfig(sysConfig)
	assert.NoError(t, err)
	num, err := res.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), num)

	resSysConfig, err := db.GetSysConfig(sysConfig.Type, sysConfig.Key)
	assert.NoError(t, err)
	checkSysConfig(t, sysConfig, resSysConfig)

	resSysConfigList, err := db.ListSysConfig(sysConfig.Type, 1, 10)
	assert.NoError(t, err)
	checkSysConfig(t, sysConfig, &resSysConfigList[0])

	resSysConfigList, err = db.ListSysConfigAll(sysConfig.Type)
	assert.NoError(t, err)
	checkSysConfig(t, sysConfig, &resSysConfigList[0])

	sysConfig.Value = "test"
	res, err = db.UpdateSysConfig(sysConfig)
	assert.NoError(t, err)
	num, err = res.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), num)
	resSysConfig, err = db.GetSysConfig(sysConfig.Type, sysConfig.Key)
	assert.NoError(t, err)
	checkSysConfig(t, sysConfig, resSysConfig)

	res, err = db.DeleteSysConfig(sysConfig.Type, sysConfig.Key)
	assert.NoError(t, err)
	num, err = res.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), num)
}

func checkSysConfig(t *testing.T, expect, actual *models.SysConfig) {
	assert.Equal(t, expect.Type, actual.Type)
	assert.Equal(t, expect.Key, actual.Key)
	assert.Equal(t, expect.Value, actual.Value)
}
