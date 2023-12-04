// Package database 数据库存储实现
package database

import (
	"fmt"
	"testing"
	"time"

	"github.com/baetyl/baetyl-go/v2/log"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

var (
	cfgtabales = []string{
		`
CREATE TABLE baetyl_configuration
(
	id                integer       PRIMARY KEY AUTOINCREMENT,
    name              varchar(128)  NOT NULL DEFAULT '',
    namespace         varchar(64)   NOT NULL DEFAULT '',
	version           varchar(36)   NOT NULL DEFAULT '',
	is_system         integer(1)    NOT NULL,
	labels            varchar(2048) NOT NULL DEFAULT '{}',
	data              varchar(2048) NOT NULL DEFAULT '{}',
	description       varchar(1024) NOT NULL DEFAULT '',
    create_time       timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time       timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`,
	}
)

func (d *BaetylCloudDB) MockCreateCfgTable() {
	for _, sql := range cfgtabales {
		_, err := d.Exec(nil, sql)
		if err != nil {
			panic(fmt.Sprintf("create cfg exception: %s", err.Error()))
		}
	}
}

func TestConfiguration(t *testing.T) {
	cfg := &specV1.Configuration{
		Name:              "cfg123",
		Namespace:         "default",
		Description:       "desc",
		CreationTimestamp: time.Unix(1000, 1000),
		Labels:            map[string]string{"a": "b"},
		Data:              map[string]string{"cfg1": "123", "cfg2": "abc"},
	}
	listOptions := &models.ListOptions{}
	log.L().Info("Test cfg", log.Any("configuration", cfg))

	db, err := MockNewDB()
	if err != nil {
		fmt.Printf("get mock sqlite3 error = %s", err.Error())
		t.Fail()
		return
	}
	db.MockCreateCfgTable()

	res, err := db.CreateConfig(nil, "default", cfg)
	assert.NoError(t, err)
	checkCfg(t, cfg, res)

	cfg2 := &specV1.Configuration{Name: "tx", Namespace: "tx"}
	tx, err := db.BeginTx()
	assert.NoError(t, err)
	res, err = db.CreateConfig(tx, "tx", cfg2)
	assert.NoError(t, err)
	assert.Equal(t, res.Namespace, "tx")
	assert.NoError(t, tx.Commit())

	res, err = db.GetConfig(nil, cfg.Namespace, cfg.Name, cfg.Version)
	assert.NoError(t, err)
	checkCfg(t, cfg, res)

	cfg.Labels = map[string]string{"b": "b"}
	res, err = db.UpdateConfig(nil, "default", cfg)
	assert.NoError(t, err)
	checkCfg(t, cfg, res)

	res, err = db.UpdateConfig(nil, "default", cfg)
	assert.NoError(t, err)
	checkCfg(t, cfg, res)

	res, err = db.GetConfig(nil, cfg.Namespace, cfg.Name, cfg.Version)
	assert.NoError(t, err)
	checkCfg(t, cfg, res)

	resList, err := db.ListConfig(cfg.Namespace, listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 1)

	err = db.DeleteConfig(nil, cfg.Namespace, cfg.Name)
	assert.NoError(t, err)

	res, err = db.GetConfig(nil, cfg.Namespace, cfg.Name, cfg.Version)
	assert.Nil(t, res)
}

func TestListCfg(t *testing.T) {
	cfg1 := &specV1.Configuration{
		Name:              "cfg_123",
		Namespace:         "default",
		Description:       "desc",
		CreationTimestamp: time.Unix(1000, 1000),
		Labels:            map[string]string{"label": "aaa"},
		Data:              map[string]string{"cfg1": "123", "cfg2": "abc"},
	}
	cfg2 := &specV1.Configuration{
		Name:              "cfg_abc",
		Namespace:         "default",
		Description:       "desc",
		CreationTimestamp: time.Unix(1000, 1000),
		Labels:            map[string]string{"label": "aaa"},
		Data:              map[string]string{"cfg1": "123", "cfg2": "abc"},
	}
	cfg3 := &specV1.Configuration{
		Name:              "cfg_test",
		Namespace:         "default",
		Description:       "desc",
		CreationTimestamp: time.Unix(1000, 1000),
		Labels:            map[string]string{"label": "bbb"},
		Data:              map[string]string{"cfg1": "123", "cfg2": "abc"},
	}
	cfg4 := &specV1.Configuration{
		Name:              "cfg_testabc",
		Namespace:         "default",
		Description:       "desc",
		CreationTimestamp: time.Unix(1000, 1000),
		Labels:            map[string]string{"label": "bbb"},
		Data:              map[string]string{"cfg1": "123", "cfg2": "abc"},
	}
	listOptions := &models.ListOptions{}

	db, err := MockNewDB()
	if err != nil {
		fmt.Printf("get mock sqlite3 error = %s", err.Error())
		t.Fail()
		return
	}
	db.MockCreateCfgTable()

	res, err := db.CreateConfig(nil, "default", cfg1)
	assert.NoError(t, err)
	checkCfg(t, cfg1, res)

	res, err = db.CreateConfig(nil, "default", cfg2)
	assert.NoError(t, err)
	checkCfg(t, cfg2, res)

	res, err = db.CreateConfig(nil, "default", cfg3)
	assert.NoError(t, err)
	checkCfg(t, cfg3, res)

	res, err = db.CreateConfig(nil, "default", cfg4)
	assert.NoError(t, err)
	checkCfg(t, cfg4, res)

	// list option nil, return all cfgs
	resList, err := db.ListConfig("default", listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 4)
	assert.Equal(t, cfg1.Name, resList.Items[0].Name)
	assert.Equal(t, cfg2.Name, resList.Items[1].Name)
	assert.Equal(t, cfg3.Name, resList.Items[2].Name)
	assert.Equal(t, cfg4.Name, resList.Items[3].Name)
	// page 1 num 2
	listOptions.PageNo = 1
	listOptions.PageSize = 2
	resList, err = db.ListConfig("default", listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 4)
	assert.Equal(t, cfg1.Name, resList.Items[0].Name)
	assert.Equal(t, cfg2.Name, resList.Items[1].Name)
	// page 2 num 2
	listOptions.PageNo = 2
	listOptions.PageSize = 2
	resList, err = db.ListConfig("default", listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 4)
	assert.Equal(t, cfg3.Name, resList.Items[0].Name)
	assert.Equal(t, cfg4.Name, resList.Items[1].Name)
	// page 3 num 0
	listOptions.PageNo = 3
	listOptions.PageSize = 2
	resList, err = db.ListConfig("default", listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 4)
	// page 1 num 2 name like cfg
	listOptions.PageNo = 1
	listOptions.PageSize = 2
	listOptions.Name = "cfg"
	resList, err = db.ListConfig("default", listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 4)
	assert.Equal(t, cfg1.Name, resList.Items[0].Name)
	assert.Equal(t, cfg2.Name, resList.Items[1].Name)
	// page 1 num 2 name like abc
	listOptions.PageNo = 1
	listOptions.PageSize = 2
	listOptions.Name = "abc"
	resList, err = db.ListConfig("default", listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 2)
	assert.Equal(t, cfg2.Name, resList.Items[0].Name)
	assert.Equal(t, cfg4.Name, resList.Items[1].Name)
	// page 1 num2 label : aaa
	listOptions.PageNo = 1
	listOptions.PageSize = 4
	listOptions.Name = ""
	listOptions.LabelSelector = "label=aaa"
	resList, err = db.ListConfig("default", listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 2)
	assert.Equal(t, cfg1.Name, resList.Items[0].Name)
	assert.Equal(t, cfg2.Name, resList.Items[1].Name)

	listOptions.PageNo = 1
	listOptions.PageSize = 4
	listOptions.Name = "abc"
	listOptions.LabelSelector = "label=aaa"
	resList, err = db.ListConfig("default", listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 1)
	assert.Equal(t, cfg2.Name, resList.Items[0].Name)

	err = db.DeleteConfig(nil, "default", cfg1.Name)
	assert.NoError(t, err)
	err = db.DeleteConfig(nil, "default", cfg2.Name)
	assert.NoError(t, err)
	err = db.DeleteConfig(nil, "default", cfg3.Name)
	assert.NoError(t, err)
	err = db.DeleteConfig(nil, "default", cfg4.Name)
	assert.NoError(t, err)

	res, err = db.GetConfig(nil, "default", cfg1.Name, "")
	assert.Nil(t, res)
	res, err = db.GetConfig(nil, "default", cfg2.Name, "")
	assert.Nil(t, res)
	res, err = db.GetConfig(nil, "default", cfg3.Name, "")
	assert.Nil(t, res)
	res, err = db.GetConfig(nil, "default", cfg4.Name, "")
	assert.Nil(t, res)
}

func checkCfg(t *testing.T, expect, actual *specV1.Configuration) {
	assert.Equal(t, expect.Name, actual.Name)
	assert.Equal(t, expect.Namespace, actual.Namespace)
	assert.Equal(t, expect.System, actual.System)
	assert.Equal(t, expect.Description, actual.Description)
	assert.Equal(t, expect.Data, actual.Data)
	assert.EqualValues(t, expect.Labels, actual.Labels)
}
