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
	secrettabales = []string{
		`
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

func (d *BaetylCloudDB) MockCreateSecretTable() {
	for _, sql := range secrettabales {
		_, err := d.Exec(nil, sql)
		if err != nil {
			panic(fmt.Sprintf("create secret exception: %s", err.Error()))
		}
	}
}

func TestSecret(t *testing.T) {
	secret := &specV1.Secret{
		Name:              "secret123",
		Namespace:         "default",
		Description:       "desc",
		CreationTimestamp: time.Unix(1000, 1000),
		Labels:            map[string]string{"a": "b"},
		Data: map[string][]byte{
			"cfg1": []byte("123"),
			"cfg2": []byte("abc"),
		},
	}
	listOptions := &models.ListOptions{}
	log.L().Info("Test secret", log.Any("secret", secret))

	db, err := MockNewDB()
	if err != nil {
		fmt.Printf("get mock sqlite3 error = %s", err.Error())
		t.Fail()
		return
	}
	db.MockCreateSecretTable()

	res, err := db.CreateSecret(nil, "default", secret)
	assert.NoError(t, err)
	checkSecret(t, secret, res)

	secret2 := &specV1.Secret{Name: "tx", Namespace: "tx"}
	tx, err := db.BeginTx()
	assert.NoError(t, err)
	res, err = db.CreateSecret(tx, "tx", secret2)
	assert.NoError(t, err)
	assert.Equal(t, res.Namespace, "tx")
	assert.NoError(t, tx.Commit())

	res, err = db.GetSecret(nil, secret.Namespace, secret.Name, secret.Version)
	assert.NoError(t, err)
	checkSecret(t, secret, res)

	secret.Labels = map[string]string{"b": "b"}
	res, err = db.UpdateSecret("default", secret)
	assert.NoError(t, err)
	checkSecret(t, secret, res)

	res, err = db.GetSecret(nil, secret.Namespace, secret.Name, secret.Version)
	assert.NoError(t, err)
	checkSecret(t, secret, res)

	resList, err := db.ListSecret(secret.Namespace, listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 1)

	err = db.DeleteSecret(nil, secret.Namespace, secret.Name)
	assert.NoError(t, err)

	res, err = db.GetSecret(nil, secret.Namespace, secret.Name, secret.Version)
	assert.Nil(t, res)
}

func TestListSecret(t *testing.T) {
	secret1 := &specV1.Secret{
		Name:              "secret_123",
		Namespace:         "default",
		Description:       "desc",
		CreationTimestamp: time.Unix(1000, 1000),
		Labels:            map[string]string{"label": "aaa"},
		Data: map[string][]byte{
			"cfg1": []byte("123"),
			"cfg2": []byte("abc"),
		},
	}
	secret2 := &specV1.Secret{
		Name:              "secret_abc",
		Namespace:         "default",
		Description:       "desc",
		CreationTimestamp: time.Unix(1000, 1000),
		Labels:            map[string]string{"label": "aaa"},
		Data: map[string][]byte{
			"cfg1": []byte("123"),
			"cfg2": []byte("abc"),
		},
	}
	secret3 := &specV1.Secret{
		Name:              "secret_test",
		Namespace:         "default",
		Description:       "desc",
		CreationTimestamp: time.Unix(1000, 1000),
		Labels:            map[string]string{"label": "bbb"},
		Data: map[string][]byte{
			"cfg1": []byte("123"),
			"cfg2": []byte("abc"),
		},
	}
	secret4 := &specV1.Secret{
		Name:              "secret_testabc",
		Namespace:         "default",
		Description:       "desc",
		CreationTimestamp: time.Unix(1000, 1000),
		Labels:            map[string]string{"label": "bbb"},
		Data: map[string][]byte{
			"cfg1": []byte("123"),
			"cfg2": []byte("abc"),
		},
	}
	listOptions := &models.ListOptions{}

	db, err := MockNewDB()
	if err != nil {
		fmt.Printf("get mock sqlite3 error = %s", err.Error())
		t.Fail()
		return
	}
	db.MockCreateSecretTable()

	res, err := db.CreateSecret(nil, "default", secret1)
	assert.NoError(t, err)
	checkSecret(t, secret1, res)

	res, err = db.CreateSecret(nil, "default", secret2)
	assert.NoError(t, err)
	checkSecret(t, secret2, res)

	res, err = db.CreateSecret(nil, "default", secret3)
	assert.NoError(t, err)
	checkSecret(t, secret3, res)

	res, err = db.CreateSecret(nil, "default", secret4)
	assert.NoError(t, err)
	checkSecret(t, secret4, res)

	// list option nil, return all cfgs
	resList, err := db.ListSecret("default", listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 4)
	assert.Equal(t, secret1.Name, resList.Items[0].Name)
	assert.Equal(t, secret2.Name, resList.Items[1].Name)
	assert.Equal(t, secret3.Name, resList.Items[2].Name)
	assert.Equal(t, secret4.Name, resList.Items[3].Name)
	// page 1 num 2
	listOptions.PageNo = 1
	listOptions.PageSize = 2
	resList, err = db.ListSecret("default", listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 4)
	assert.Equal(t, secret1.Name, resList.Items[0].Name)
	assert.Equal(t, secret2.Name, resList.Items[1].Name)
	// page 2 num 2
	listOptions.PageNo = 2
	listOptions.PageSize = 2
	resList, err = db.ListSecret("default", listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 4)
	assert.Equal(t, secret3.Name, resList.Items[0].Name)
	assert.Equal(t, secret4.Name, resList.Items[1].Name)
	// page 3 num 0
	listOptions.PageNo = 3
	listOptions.PageSize = 2
	resList, err = db.ListSecret("default", listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 4)
	// page 1 num 2 name like cfg
	listOptions.PageNo = 1
	listOptions.PageSize = 2
	listOptions.Name = "secret"
	resList, err = db.ListSecret("default", listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 4)
	assert.Equal(t, secret1.Name, resList.Items[0].Name)
	assert.Equal(t, secret2.Name, resList.Items[1].Name)
	// page 1 num 2 name like abc
	listOptions.PageNo = 1
	listOptions.PageSize = 2
	listOptions.Name = "abc"
	resList, err = db.ListSecret("default", listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 2)
	assert.Equal(t, secret2.Name, resList.Items[0].Name)
	assert.Equal(t, secret4.Name, resList.Items[1].Name)
	// page 1 num2 label : aaa
	listOptions.PageNo = 1
	listOptions.PageSize = 4
	listOptions.Name = ""
	listOptions.LabelSelector = "label=aaa"
	resList, err = db.ListSecret("default", listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 2)
	assert.Equal(t, secret1.Name, resList.Items[0].Name)
	assert.Equal(t, secret2.Name, resList.Items[1].Name)

	listOptions.PageNo = 1
	listOptions.PageSize = 4
	listOptions.Name = "abc"
	listOptions.LabelSelector = "label=aaa"
	resList, err = db.ListSecret("default", listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 1)
	assert.Equal(t, secret2.Name, resList.Items[0].Name)

	err = db.DeleteSecret(nil, "default", secret1.Name)
	assert.NoError(t, err)
	err = db.DeleteSecret(nil, "default", secret2.Name)
	assert.NoError(t, err)
	err = db.DeleteSecret(nil, "default", secret3.Name)
	assert.NoError(t, err)
	err = db.DeleteSecret(nil, "default", secret4.Name)
	assert.NoError(t, err)

	res, err = db.GetSecret(nil, "default", secret1.Name, "")
	assert.Nil(t, res)
	res, err = db.GetSecret(nil, "default", secret2.Name, "")
	assert.Nil(t, res)
	res, err = db.GetSecret(nil, "default", secret3.Name, "")
	assert.Nil(t, res)
	res, err = db.GetSecret(nil, "default", secret4.Name, "")
	assert.Nil(t, res)
}

func checkSecret(t *testing.T, expect, actual *specV1.Secret) {
	assert.Equal(t, expect.Name, actual.Name)
	assert.Equal(t, expect.Namespace, actual.Namespace)
	assert.Equal(t, expect.System, actual.System)
	assert.Equal(t, expect.Description, actual.Description)
	assert.Equal(t, expect.Data, actual.Data)
	assert.EqualValues(t, expect.Labels, actual.Labels)
}
