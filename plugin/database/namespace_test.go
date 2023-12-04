// Package database 数据库存储实现
package database

import (
	"fmt"
	"testing"

	"github.com/baetyl/baetyl-go/v2/log"
	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

var (
	nstabales = []string{
		`
CREATE TABLE baetyl_namespace
(
	id                integer       PRIMARY KEY AUTOINCREMENT,
    name              varchar(128)  NOT NULL DEFAULT '',
    create_time       timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time       timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE baetyl_node
(
	id                integer       PRIMARY KEY AUTOINCREMENT,
    name              varchar(128)  NOT NULL DEFAULT '',
    namespace         varchar(64)   NOT NULL DEFAULT '',
	version           varchar(36)   NOT NULL DEFAULT '',
    nodeMode          varchar(36)   NOT NULL DEFAULT '',
	core_version      varchar(36)   NOT NULL DEFAULT '',
	labels            varchar(2048) NOT NULL DEFAULT '{}',
	annotations       varchar(2048) NOT NULL DEFAULT '',
	attributes        varchar(2048) NOT NULL DEFAULT '{}',
	description       varchar(1024) NOT NULL DEFAULT '',
    create_time       timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time       timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE baetyl_application
(
	id                integer       PRIMARY KEY AUTOINCREMENT,
    name              varchar(128)  NOT NULL DEFAULT '',
    namespace         varchar(64)   NOT NULL DEFAULT '',
	version           varchar(36)   NOT NULL DEFAULT '',
	type              varchar(36)   NOT NULL DEFAULT '',
	mode              varchar(36)   NOT NULL DEFAULT '',
	is_system         integer(1)    NOT NULL DEFAULT 0,
	cron_status       integer(1)    NOT NULL DEFAULT 0,
	labels            varchar(2048) NOT NULL DEFAULT '{}',
    selector          varchar(255)  NOT NULL DEFAULT '{}',
    node_selector     varchar(255)  NOT NULL DEFAULT '',
	services          text          NULL,
	volumes           text          NULL,
	description       varchar(1024) NOT NULL DEFAULT '',
	cron_time         timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    create_time       timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time       timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP
);
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
CREATE TABLE baetyl_secret
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

func (d *BaetylCloudDB) MockCreateNsTable() {
	for _, sql := range nstabales {
		_, err := d.Exec(nil, sql)
		if err != nil {
			panic(fmt.Sprintf("create ns exception: %s", err.Error()))
		}
	}
}

func TestNamespace(t *testing.T) {
	ns1 := &models.Namespace{
		Name: "ns1",
	}
	ns2 := &models.Namespace{
		Name: "ns2",
	}
	log.L().Info("Test ns", log.Any("namespace", ns1))

	db, err := MockNewDB()
	if err != nil {
		fmt.Printf("get mock sqlite3 error = %s", err.Error())
		t.Fail()
		return
	}
	db.MockCreateNsTable()

	res, err := db.CreateNamespace(ns1)
	assert.NoError(t, err)
	checkNamespace(t, ns1, res)

	res, err = db.GetNamespace(ns1.Name)
	assert.NoError(t, err)
	checkNamespace(t, ns1, res)

	err = db.DeleteNamespace(ns1)
	assert.NoError(t, err)

	res, err = db.GetNamespace(ns1.Name)
	assert.Nil(t, res)

	res, err = db.CreateNamespace(ns1)
	assert.NoError(t, err)

	res, err = db.CreateNamespace(ns2)
	assert.NoError(t, err)

	params := &models.ListOptions{}
	resList, err := db.ListNamespace(params)
	assert.NoError(t, err)
	assert.Equal(t, 2, resList.Total)
	assert.Equal(t, ns1.Name, resList.Items[0].Name)
	assert.Equal(t, ns2.Name, resList.Items[1].Name)
}

func checkNamespace(t *testing.T, expect, actual *models.Namespace) {
	assert.Equal(t, expect.Name, actual.Name)
}
