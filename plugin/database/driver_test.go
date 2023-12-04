// Package database 数据库存储实现
package database

import (
	"fmt"
	"testing"

	"github.com/baetyl/baetyl-go/v2/context"
	"github.com/baetyl/baetyl-go/v2/log"
	v1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
)

var (
	driverTables = []string{
		`
CREATE TABLE baetyl_driver
(
	id                integer       PRIMARY KEY AUTOINCREMENT,
    name              varchar(128)  NOT NULL DEFAULT '',
    namespace         varchar(128)  NOT NULL DEFAULT '',
	version           varchar(36)   NOT NULL DEFAULT '',
	type              integer       NOT NULL DEFAULT 0,
    mode              varchar(32)   NOT NULL DEFAULT '',
    labels            varchar(2048) NOT NULL DEFAULT '',
    protocol          varchar(64)   NOT NULL DEFAULT '',
    arch              varchar(32)   NOT NULL DEFAULT '',
	description       varchar(255)  NOT NULL DEFAULT '',
    default_config    varchar(1024) NOT NULL DEFAULT '',
    service           text,
    volumes           text,
    registries        varchar(1024) NOT NULL DEFAULT '',
    program_config    varchar(128)  NOT NULL DEFAULT '',
    create_time       timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time       timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`,
	}
)

func (d *BaetylCloudDB) MockCreateDriverTable() {
	for _, sql := range driverTables {
		_, err := d.Exec(nil, sql)
		if err != nil {
			panic(fmt.Sprintf("create device model exception: %s", err.Error()))
		}
	}
}

func TestDriver(t *testing.T) {
	driver := &models.Driver{
		Name:      "dm-1",
		Namespace: "default",
		Type:      1,
		Mode:      context.RunModeKube,
		Labels: map[string]string{
			common.LabelAppMode: context.RunModeKube,
		},
		Description:   "desc",
		Protocol:      "pro-1",
		Architecture:  "amd64",
		DefaultConfig: "default config",
		Service: &models.Service{
			Image: "image",
			Resources: &v1.Resources{
				Limits:   map[string]string{"cpu": "2"},
				Requests: map[string]string{"cpu": "1"},
			},
			Ports:           []v1.ContainerPort{{HostPort: 80, ContainerPort: 80}},
			Env:             []v1.Environment{{Name: "a", Value: "b"}},
			SecurityContext: &v1.SecurityContext{Privileged: true},
			HostNetwork:     true,
			Args:            []string{"abc", "def"},
			VolumeMounts:    []v1.VolumeMount{{Name: "data", MountPath: "/data"}},
		},
		Volumes: []v1.Volume{{
			Name:         "data",
			VolumeSource: v1.VolumeSource{HostPath: &v1.HostPathVolumeSource{Path: "/data"}},
		}},
		Registries: []models.RegistryView{{
			Name:     "test",
			Address:  "test.com",
			Username: "test",
		}},
	}
	listOptions := &models.ListOptions{}
	log.L().Info("Test driver", log.Any("driver", driver))

	db, err := MockNewDB()
	if err != nil {
		fmt.Printf("get mock sqlite3 error = %s", err.Error())
		t.Fail()
		return
	}
	db.MockCreateDriverTable()

	res, err := db.CreateDriver(driver)
	assert.NoError(t, err)
	checkDriver(t, driver, res)

	res, err = db.GetDriver(driver.Namespace, driver.Name)
	assert.NoError(t, err)
	checkDriver(t, driver, res)

	driver.Architecture = "arm64"
	driver.Labels["aa"] = "bb"
	res, err = db.UpdateDriver(driver)
	assert.NoError(t, err)
	res.Labels["aa"] = "bb"
	checkDriver(t, driver, res)

	res, err = db.GetDriver(driver.Namespace, driver.Name)
	assert.NoError(t, err)
	checkDriver(t, driver, res)

	resList, err := db.ListDriver(driver.Namespace, listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 1)

	err = db.DeleteDriver(driver.Namespace, driver.Name)
	assert.NoError(t, err)

	res, err = db.GetDriver(driver.Namespace, driver.Name)
	assert.Nil(t, res)
}

func TestListDriver(t *testing.T) {
	d1 := &models.Driver{
		Name:        "driver-123-1",
		Namespace:   "default",
		Description: "desc",
	}
	d2 := &models.Driver{
		Name:        "driver-abc-1",
		Namespace:   "default",
		Description: "desc",
	}
	d3 := &models.Driver{
		Name:        "driver-123-2",
		Namespace:   "default",
		Description: "desc",
	}
	d4 := &models.Driver{
		Name:        "driver-abc-2",
		Namespace:   "default",
		Description: "desc",
	}
	listOptions := &models.ListOptions{}

	db, err := MockNewDB()
	if err != nil {
		fmt.Printf("get mock sqlite3 error = %s", err.Error())
		t.Fail()
		return
	}
	db.MockCreateDriverTable()

	res, err := db.CreateDriver(d1)
	assert.NoError(t, err)
	checkDriver(t, d1, res)

	res, err = db.CreateDriver(d2)
	assert.NoError(t, err)
	checkDriver(t, d2, res)

	res, err = db.CreateDriver(d3)
	assert.NoError(t, err)
	checkDriver(t, d3, res)

	res, err = db.CreateDriver(d4)
	assert.NoError(t, err)
	checkDriver(t, d4, res)

	// list option nil, return all drivers
	resList, err := db.ListDriver("default", listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 4)
	assert.Equal(t, d1.Name, resList.Items[0].Name)
	assert.Equal(t, d2.Name, resList.Items[1].Name)
	assert.Equal(t, d3.Name, resList.Items[2].Name)
	assert.Equal(t, d4.Name, resList.Items[3].Name)
	// page 1 num 2
	listOptions.PageNo = 1
	listOptions.PageSize = 2
	resList, err = db.ListDriver("default", listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 4)
	assert.Equal(t, d1.Name, resList.Items[0].Name)
	assert.Equal(t, d2.Name, resList.Items[1].Name)
	// page 2 num 2
	listOptions.PageNo = 2
	listOptions.PageSize = 2
	resList, err = db.ListDriver("default", listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 4)
	assert.Equal(t, d3.Name, resList.Items[0].Name)
	assert.Equal(t, d4.Name, resList.Items[1].Name)
	// page 3 num 0
	listOptions.PageNo = 3
	listOptions.PageSize = 2
	resList, err = db.ListDriver("default", listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 4)
	// page 1 num 2 name like driver
	listOptions.PageNo = 1
	listOptions.PageSize = 2
	listOptions.Name = "driver"
	resList, err = db.ListDriver("default", listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 4)
	assert.Equal(t, d1.Name, resList.Items[0].Name)
	assert.Equal(t, d2.Name, resList.Items[1].Name)
	// page 1 num 2 name like abc
	listOptions.PageNo = 1
	listOptions.PageSize = 2
	listOptions.Name = "abc"
	resList, err = db.ListDriver("default", listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 2)
	assert.Equal(t, d2.Name, resList.Items[0].Name)
	assert.Equal(t, d4.Name, resList.Items[1].Name)
	// page 1 num2 name like 123
	listOptions.PageNo = 1
	listOptions.PageSize = 4
	listOptions.Name = "123"
	resList, err = db.ListDriver("default", listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 2)
	assert.Equal(t, d1.Name, resList.Items[0].Name)
	assert.Equal(t, d3.Name, resList.Items[1].Name)

	err = db.DeleteDriver("default", d1.Name)
	assert.NoError(t, err)
	err = db.DeleteDriver("default", d2.Name)
	assert.NoError(t, err)
	err = db.DeleteDriver("default", d3.Name)
	assert.NoError(t, err)
	err = db.DeleteDriver("default", d4.Name)
	assert.NoError(t, err)

	res, err = db.GetDriver("default", d1.Name)
	assert.Nil(t, res)
	res, err = db.GetDriver("default", d2.Name)
	assert.Nil(t, res)
	res, err = db.GetDriver("default", d3.Name)
	assert.Nil(t, res)
	res, err = db.GetDriver("default", d4.Name)
	assert.Nil(t, res)
}

func checkDriver(t *testing.T, expect, actual *models.Driver) {
	assert.Equal(t, expect.Name, actual.Name)
	assert.Equal(t, expect.Namespace, actual.Namespace)
	assert.Equal(t, expect.Description, actual.Description)
	assert.Equal(t, expect.Protocol, actual.Protocol)
	assert.Equal(t, expect.Architecture, actual.Architecture)
	assert.Equal(t, expect.Type, actual.Type)
	assert.Equal(t, expect.Mode, actual.Mode)
	assert.Equal(t, expect.Labels, actual.Labels)
	assert.Equal(t, expect.DefaultConfig, actual.DefaultConfig)
	assert.EqualValues(t, expect.Service, actual.Service)
	assert.EqualValues(t, expect.Volumes, actual.Volumes)
	assert.EqualValues(t, expect.Registries, actual.Registries)
}
