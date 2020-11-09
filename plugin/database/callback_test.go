package database

import (
	"fmt"
	"testing"
	"time"

	"github.com/baetyl/baetyl-go/v2/log"
	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

var (
	callbackTables = []string{
		`
CREATE TABLE baetyl_callback
(
    name        varchar(128)  NOT NULL DEFAULT '' PRIMARY KEY,
    namespace   varchar(64)   NOT NULL DEFAULT '',
    method      varchar(36)   NOT NULL DEFAULT 'GET',
    params      varchar(2048) NOT NULL DEFAULT '{}',
    header      varchar(1024) NOT NULL DEFAULT '{}',
    body        varchar(2048) NOT NULL DEFAULT '{}',
    url         varchar(1024) NOT NULL DEFAULT '',
    description varchar(1024) NOT NULL DEFAULT '',
    create_time timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`,
	}
)

func (d *DB) MockCreateCallbackTable() {
	for _, sql := range callbackTables {
		_, err := d.Exec(nil, sql)
		if err != nil {
			panic(fmt.Sprintf("create table exception: %s", err.Error()))
		}
	}
}

func TestCallback(t *testing.T) {
	call := &models.Callback{
		Name:        "29987d6a2b8f11eabc62186590da6863",
		Namespace:   "default",
		Method:      "Get",
		Url:         "http://www.baidu.com",
		Description: "desc",
		CreateTime:  time.Unix(1000, 1000),
		UpdateTime:  time.Unix(1000, 1000),
		Params:      map[string]string{"a": "b"},
		Header:      map[string]string{"e": "g"},
		Body:        map[string]string{"f": "v"},
	}
	log.L().Info("Test callback", log.Any("call", call))

	db, err := MockNewDB()
	if err != nil {
		fmt.Printf("get mock sqlite3 error = %s", err.Error())
		t.Fail()
		return
	}
	db.MockCreateCallbackTable()

	res, err := db.CreateCallback(call)
	assert.NoError(t, err)
	num, err := res.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), num)

	resCall, err := db.GetCallback(call.Name, call.Namespace)
	assert.NoError(t, err)
	checkCallback(t, call, resCall)

	call.Params = map[string]string{"b": "b"}
	res, err = db.UpdateCallback(call)
	assert.NoError(t, err)
	num, err = res.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), num)
	resCall, err = db.GetCallback(call.Name, call.Namespace)
	assert.NoError(t, err)
	checkCallback(t, call, resCall)

	res, err = db.DeleteCallback(call.Name, call.Namespace)
	assert.NoError(t, err)
	num, err = res.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), num)
}

func checkCallback(t *testing.T, expect, actual *models.Callback) {
	assert.Equal(t, expect.Name, actual.Name)
	assert.Equal(t, expect.Namespace, actual.Namespace)
	assert.Equal(t, expect.Description, actual.Description)
	assert.Equal(t, expect.Method, actual.Method)
	assert.EqualValues(t, expect.Params, actual.Params)
	assert.EqualValues(t, expect.Header, actual.Header)
	assert.EqualValues(t, expect.Body, actual.Body)
	assert.Equal(t, expect.Url, actual.Url)
}
