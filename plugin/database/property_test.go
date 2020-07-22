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

func genProperty() *models.Property {
	return &models.Property{
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
func TestProperty(t *testing.T) {
	property := genProperty()

	db, err := MockNewDB()
	assert.NoError(t, err)
	db.MockCreatePropertyTable()

	_, err = db.CreateProperty(property)
	assert.NoError(t, err)

	newValue := "updated_" + property.Value
	property.Value = newValue
	updated_property, err := db.UpdateProperty(property)
	assert.NoError(t, err)
	assert.Equal(t, newValue, updated_property.Value)

	get_property, err := db.GetProperty(property.Key)
	assert.NoError(t, err)
	checkProperty(t, get_property, property)
	_, err = db.GetProperty("bad key")
	assert.Error(t, err)

	page := &models.Filter{
		PageNo:   1,
		PageSize: 2,
		Name:     "%",
	}
	resPropertyListView, err := db.ListProperty(page)
	assert.NoError(t, err)
	checkProperty(t, property, &resPropertyListView.Data.Rows.([]models.Property)[0])

	err = db.DeleteProperty(property.Key)
	assert.NoError(t, err)

}

func checkProperty(t *testing.T, expect, actual *models.Property) {
	assert.Equal(t, expect.Key, actual.Key)
	assert.Equal(t, expect.Value, actual.Value)
}
