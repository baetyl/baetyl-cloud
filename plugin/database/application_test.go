package database

import (
	"fmt"
	"testing"

	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

var (
	applicationTebles = []string{`
CREATE TABLE baetyl_application_history
(
    id          integer             PRIMARY KEY AUTOINCREMENT,
    namespace   varchar(64)         NOT NULL DEFAULT '' ,
    name        varchar(128)        NOT NULL DEFAULT '' ,
    version     varchar(36)         NOT NULL DEFAULT '' ,
    is_deleted  smallint            NOT NULL DEFAULT 0  ,
    create_time timestamp           NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time timestamp           NOT NULL DEFAULT CURRENT_TIMESTAMP ,
    content     BLOB                NOT NULL DEFAULT '' 

);
`,
	}
)

func (d *DB) MockCreateApplicationTable() {
	for _, sql := range applicationTebles {
		_, err := d.Exec(nil, sql)
		if err != nil {
			panic(fmt.Sprintf("create table exception: %s", err.Error()))
		}
	}
}

func mockDb(t *testing.T) *DB {
	db, err := MockNewDB()
	if err != nil {
		fmt.Printf("get mock sqlite3 error = %s", err.Error())
		t.Fail()
		return nil
	}
	db.MockCreateApplicationTable()
	return db
}

func TestDB_CreateApplication(t *testing.T) {
	app := &specV1.Application{
		Name:        "29987d6a2b8f11eabc62186590da6863",
		Namespace:   "default",
		Description: "desc",
		Version:     "1",
	}
	db := mockDb(t)
	res, err := db.CreateApplicationHis(app)
	assert.NoError(t, err)
	num, err := res.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), num)

}

func TestDB_UpdateApplication(t *testing.T) {
	db := mockDb(t)
	app := &specV1.Application{
		Name:        "test",
		Namespace:   "default",
		Description: "desc",
		Version:     "2",
	}

	app2 := &specV1.Application{
		Name:        "test",
		Namespace:   "default",
		Description: "desc",
		Version:     "1",
	}

	res, err := db.CreateApplicationHis(app2)
	assert.NoError(t, err)

	res, err = db.UpdateApplicationHis(app, "1")
	assert.NoError(t, err)
	num, err := res.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), num)

	newApp, err := db.GetApplicationHis(app.Name, app.Namespace, app.Version)
	assert.NoError(t, err)
	checkApplication(t, app, newApp)

	res, err = db.DeleteApplicationHis(newApp.Name, newApp.Namespace, newApp.Version)
	assert.NoError(t, err)
	num, err = res.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), num)
}

func TestDB_GetApplication(t *testing.T) {
	db := mockDb(t)
	app := &specV1.Application{
		Name:        "29987d6a2b8f11eabc62186590da6863",
		Namespace:   "default",
		Description: "desc",
		Version:     "1",
	}

	_, err := db.CreateApplicationHis(app)
	assert.NoError(t, err)

	res, err := db.GetApplicationHis(app.Name, app.Namespace, app.Version)
	assert.NoError(t, err)
	checkApplication(t, app, res)

	num, err := db.CountApplicationHis(nil, app.Name, app.Namespace)
	assert.NoError(t, err)
	assert.Equal(t, 1, num)

	filter := &models.Filter{
		PageNo:   1,
		PageSize: 10,
		Name:     app.Name,
	}

	apps, err := db.ListApplicationHis(app.Namespace, filter)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(apps))
	assert.Equal(t, app.Name, apps[0].Name)
	assert.Equal(t, app.Namespace, apps[0].Namespace)
	assert.Equal(t, app.Description, apps[0].Description)

}

func checkApplication(t *testing.T, expect, actual *specV1.Application) {
	assert.Equal(t, expect.Name, actual.Name)
	assert.Equal(t, expect.Namespace, actual.Namespace)
	assert.Equal(t, expect.Description, actual.Description)
}
