package database

import (
	"fmt"
	"github.com/baetyl/baetyl-cloud/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	propertyTables = []string{
		"CREATE TABLE baetyl_property(" +
			"    `id`               integer       PRIMARY KEY AUTOINCREMENT," +
			"    `key`              varchar(128)  NOT NULL DEFAULT ''," +
			"    `value`            varchar(2048) NOT NULL DEFAULT '',  " +
			"    `create_time`      timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP," +
			"    `update_time`      timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP);",
	}
)

func genCache() *models.Cache {
	return &models.Cache{
		Key:   "baetyl_0.1.0",
		Value: "http://test.baetyl/0.1.0",
	}
}
func (d *dbStorage) MockCreatePropertyTable() {
	for _, sql := range propertyTables {
		_, err := d.exec(nil, sql)
		if err != nil {
			panic(fmt.Sprintf("create table exception: %s", err.Error()))
		}
	}
}
func TestCache(t *testing.T) {
	cache := genCache()

	db, err := MockNewDB()
	assert.NoError(t, err)
	db.MockCreatePropertyTable()

	err = db.SetCache(cache.Key, cache.Value)
	assert.NoError(t, err)

	cache.Value = "updated_" + cache.Value
	err = db.SetCache(cache.Key, cache.Value)
	assert.NoError(t, err)

	value, err := db.GetCache(cache.Key)
	assert.NoError(t, err)
	assert.Equal(t, value, cache.Value)
	value, err = db.GetCache("bad key")
	assert.Error(t, err)

	page := &models.Filter{
		PageNo:   1,
		PageSize: 2,
		Name:     "%",
	}
	resCacheListView, err := db.ListCache(page)
	assert.NoError(t, err)
	checkCache(t, cache, &resCacheListView.Items.([]models.Cache)[0])

	err = db.DeleteCache(cache.Key)
	assert.NoError(t, err)

}

func checkCache(t *testing.T, expect, actual *models.Cache) {
	assert.Equal(t, expect.Key, actual.Key)
	assert.Equal(t, expect.Value, actual.Value)
}
