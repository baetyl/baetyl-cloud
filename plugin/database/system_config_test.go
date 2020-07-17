package database

import (
	"fmt"
	"github.com/baetyl/baetyl-cloud/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	systemTables = []string{
		"CREATE TABLE baetyl_cloud_system_config(" +
			"    `id`               integer       PRIMARY KEY AUTOINCREMENT," +
			"    `key`              varchar(128)  NOT NULL DEFAULT ''," +
			"    `value`            varchar(2048) NOT NULL DEFAULT '',  " +
			"    `create_time`      timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP," +
			"    `update_time`      timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP);",
	}
)

func genSystemConfig() *models.SystemConfig {
	return &models.SystemConfig{
		Key:   "baetyl_0.1.0",
		Value: "http://test.baetyl/0.1.0",
	}
}
func (d *dbStorage) MockCreateSystemTable() {
	for _, sql := range systemTables {
		_, err := d.exec(nil, sql)
		if err != nil {
			panic(fmt.Sprintf("create table exception: %s", err.Error()))
		}
	}
}
func TestSystemConfig(t *testing.T) {
	systemConfig := genSystemConfig()

	db, err := MockNewDB()
	if err != nil {
		fmt.Printf("get mock sqlite3 error = %s", err.Error())
		t.Fail()
		return
	}
	db.MockCreateSystemTable()
	res, err := db.CreateSystemConfig(systemConfig)
	assert.NoError(t, err)
	num, err := res.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), num)

	resSystemConfig, err := db.GetSystemConfig(systemConfig.Key)
	assert.NoError(t, err)
	checkSystemConfig(t, systemConfig, resSystemConfig)

	page := &models.Filter{
		PageNo:   1,
		PageSize: 2,
		Name:     "%",
	}
	resSystemConfigList, err := db.ListSystemConfig(page.Name, page.PageNo, page.PageSize)
	assert.NoError(t, err)
	checkSystemConfig(t, systemConfig, &resSystemConfigList[0])

	systemConfig.Value = "updated_value"
	res, err = db.UpdateSystemConfig(systemConfig)
	assert.NoError(t, err)
	num, err = res.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), num)
	resSystemConfig, err = db.GetSystemConfig(systemConfig.Key)
	assert.NoError(t, err)
	checkSystemConfig(t, systemConfig, resSystemConfig)

	res, err = db.DeleteSystemConfig(systemConfig.Key)
	assert.NoError(t, err)
	num, err = res.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), num)
}

func checkSystemConfig(t *testing.T, expect, actual *models.SystemConfig) {
	assert.Equal(t, expect.Key, actual.Key)
	assert.Equal(t, expect.Value, actual.Value)
}
