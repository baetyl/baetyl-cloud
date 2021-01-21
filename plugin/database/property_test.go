package database

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

var (
	propertyTables = []string{
		"CREATE TABLE baetyl_property(" +
			"    `id`               integer       PRIMARY KEY AUTOINCREMENT," +
			"    `name`             varchar(128)  UNIQUE NOT NULL DEFAULT ''," +
			"    `value`            varchar(2048) NOT NULL DEFAULT ''," +
			"    `create_time`      timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP," +
			"    `update_time`      timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP);",
	}
)

func genProperty() *models.Property {
	return &models.Property{
		Name:  "baetyl_0.1.0",
		Value: "http://test.baetyl/0.1.0",
	}
}
func (d *DB) MockCreatePropertyTable() {
	for _, sql := range propertyTables {
		_, err := d.Exec(nil, sql)
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

	err = db.CreateProperty(property)
	assert.NoError(t, err)
	// name existed
	err = db.CreateProperty(property)
	assert.Error(t, err)

	property.Value = "updated_" + property.Value
	err = db.UpdateProperty(property)
	assert.NoError(t, err)

	value, err := db.GetPropertyValue(property.Name)
	assert.NoError(t, err)
	assert.Equal(t, property.Value, value)
	_, err = db.GetPropertyValue("bad name")
	assert.Error(t, err)

	page := &models.Filter{
		PageNo:   1,
		PageSize: 2,
	}
	properties, err := db.ListProperty(page)
	assert.NoError(t, err)
	checkProperty(t, property, &properties[0])
	count, err := db.CountProperty(page.GetFuzzyName())
	assert.Equal(t, count, 1)

	err = db.CreateProperty(&models.Property{
		Name:  "baetyl-function-runtime-python3",
		Value: "python3",
	})
	assert.NoError(t, err)
	err = db.CreateProperty(&models.Property{
		Name:  "baetyl-function-runtime-nodejs10",
		Value: "nodejs10",
	})
	assert.NoError(t, err)
	err = db.CreateProperty(&models.Property{
		Name:  "baetyl-function-runtime-sql",
		Value: "sql",
	})
	assert.NoError(t, err)

	page = &models.Filter{
		Name: "baetyl-function-runtime-",
	}
	properties, err = db.ListProperty(page)
	assert.NoError(t, err)
	assert.Len(t, properties, 3)
	count, err = db.CountProperty(page.GetFuzzyName())
	assert.Equal(t, count, 3)

	err = db.DeleteProperty(property.Name)
	assert.NoError(t, err)
}

func checkProperty(t *testing.T, expect, actual *models.Property) {
	assert.Equal(t, expect.Name, actual.Name)
	assert.Equal(t, expect.Value, actual.Value)
}
