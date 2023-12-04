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
	apptabales = []string{
		`
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
    init_services     text          NULL,
	services          text          NULL,
	volumes           text          NULL,
	description       varchar(1024) NOT NULL DEFAULT '',
	cron_time         timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    host_network      tinyint(1)    NOT NULL DEFAULT 0,
    dns_policy        varchar(64)   NOT NULL DEFAULT 'ClusterFirst',
    replica           int           NOT NULL DEFAULT 1,
    job_config        varchar(512)  NOT NULL DEFAULT '{}',
    workload          varchar(36)   NULL DEFAULT '',
	ota               varchar(2048) NOT NULL DEFAULT '{}',
	autoScaleCfg      varchar(4096) NOT NULL DEFAULT '{}',
    create_time       timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time       timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP,
	preserve_updates tinyint(1)     NOT NULL DEFAULT 0
);
`,
	}
)

func (d *BaetylCloudDB) MockCreateAppTable() {
	for _, sql := range apptabales {
		_, err := d.Exec(nil, sql)
		if err != nil {
			panic(fmt.Sprintf("create app exception: %s", err.Error()))
		}
	}
}

func TestApplication(t *testing.T) {
	app := &specV1.Application{
		Name:              "app123",
		Namespace:         "default",
		Version:           "",
		Type:              "system",
		Description:       "desc",
		CreationTimestamp: time.Unix(1000, 1000),
		Labels:            map[string]string{"a": "b"},
		Selector:          "a=b",
		NodeSelector:      "1=2",
		HostNetwork:       true,
		Replica:           2,
		Workload:          "deployment",
		JobConfig: &specV1.AppJobConfig{
			Completions:   1,
			Parallelism:   2,
			BackoffLimit:  3,
			RestartPolicy: "Never",
		},
		InitServices: []specV1.Service{
			{
				Name:  "init",
				Image: "init_image",
			},
		},
		Services: []specV1.Service{{
			Name:     "test",
			Hostname: "hostname",
			Image:    "test_image",
			Replica:  1,
			VolumeMounts: []specV1.VolumeMount{{
				Name:      "mount",
				MountPath: "path",
				ReadOnly:  false,
			}},
			Ports: []specV1.ContainerPort{{
				HostPort:      1000,
				ContainerPort: 1000,
			}},
			Devices: []specV1.Device{{DevicePath: "dev"}},
			Args:    []string{"arg"},
			Env: []specV1.Environment{{
				Name:  "key",
				Value: "value",
			}},
			Resources: &specV1.Resources{},
			Runtime:   "runtime",
		}},
		Volumes: []specV1.Volume{{
			Name: "test",
			VolumeSource: specV1.VolumeSource{
				HostPath: &specV1.HostPathVolumeSource{Path: "hostPath"},
				Config: &specV1.ObjectReference{
					Name: "config",
				},
			},
		}},
		Ota: specV1.OtaInfo{
			ApkInfo:     &specV1.ApkInfo{},
			DeviceGroup: &specV1.OtaDeviceGroup{},
			Task:        &specV1.OtaTask{},
		},
		AutoScaleCfg: &specV1.AutoScaleCfg{
			MinReplicas: 1,
			MaxReplicas: 10,
			Metrics: []specV1.MetricSpec{
				{
					Type: "Resource",
					Resource: &specV1.ResourceMetric{
						Name:               "cpu",
						TargetType:         "Utilization",
						AverageUtilization: 50,
					},
				},
			},
		},
	}
	log.L().Info("Test app", log.Any("application", app))

	db, err := MockNewDB()
	if err != nil {
		fmt.Printf("get mock sqlite3 error = %s", err.Error())
		t.Fail()
		return
	}
	db.MockCreateAppTable()

	res, err := db.CreateApplication(nil, "default", app)
	assert.NoError(t, err)
	checkApp(t, app, res)

	app2 := &specV1.Application{Name: "tx", Namespace: "tx"}
	tx, err := db.BeginTx()
	assert.NoError(t, err)
	res, err = db.CreateApplication(tx, "tx", app2)
	assert.NoError(t, err)
	assert.Equal(t, res.Namespace, "tx")
	assert.NoError(t, tx.Commit())

	res, err = db.GetApplication(nil, app.Namespace, app.Name, app.Version)
	assert.NoError(t, err)
	checkApp(t, app, res)

	app.Labels = map[string]string{"b": "b"}
	res, err = db.UpdateApplication(nil, "default", app)
	assert.NoError(t, err)
	checkApp(t, app, res)

	res, err = db.GetApplication(nil, app.Namespace, app.Name, app.Version)
	assert.NoError(t, err)
	checkApp(t, app, res)

	listOptions := &models.ListOptions{}
	resList, err := db.ListApplication(nil, app.Namespace, listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 1)

	listOptions.LabelSelector = "b=b"
	resList, err = db.ListApplication(nil, app.Namespace, listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 1)

	listOptions.LabelSelector = "b!=b"
	resList, err = db.ListApplication(nil, app.Namespace, listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 0)

	listOptions.LabelSelector = "!b"
	resList, err = db.ListApplication(nil, app.Namespace, listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 0)

	err = db.DeleteApplication(nil, app.Namespace, app.Name)
	assert.NoError(t, err)

	res, err = db.GetApplication(nil, app.Namespace, app.Name, app.Version)
	assert.Nil(t, res)
}

func TestListApp(t *testing.T) {
	app1 := &specV1.Application{
		Name:              "app_123",
		Namespace:         "default",
		Type:              "system",
		Description:       "desc",
		CreationTimestamp: time.Unix(1000, 1000),
		Labels:            map[string]string{"label": "aaa"},
		Selector:          "a=b",
		NodeSelector:      "1=2",
		HostNetwork:       true,
		Replica:           2,
		Workload:          "deployment",
		JobConfig: &specV1.AppJobConfig{
			Completions:   1,
			Parallelism:   2,
			BackoffLimit:  3,
			RestartPolicy: "Never",
		},
		InitServices: []specV1.Service{
			{
				Name:  "init",
				Image: "init_image",
			},
		},
		Services: []specV1.Service{{
			Name:     "test",
			Hostname: "hostname",
			Image:    "test_image",
			Replica:  1,
			VolumeMounts: []specV1.VolumeMount{{
				Name:      "mount",
				MountPath: "path",
				ReadOnly:  false,
			}},
			Ports: []specV1.ContainerPort{{
				HostPort:      1000,
				ContainerPort: 1000,
			}},
			Devices: []specV1.Device{{DevicePath: "dev"}},
			Args:    []string{"arg"},
			Env: []specV1.Environment{{
				Name:  "key",
				Value: "value",
			}},
			Resources: &specV1.Resources{},
			Runtime:   "runtime",
		}},
		Volumes: []specV1.Volume{{
			Name: "test",
			VolumeSource: specV1.VolumeSource{
				HostPath: &specV1.HostPathVolumeSource{Path: "hostPath"},
				Config: &specV1.ObjectReference{
					Name: "config",
				},
			},
		}},
	}
	app2 := &specV1.Application{
		Name:              "app_abc",
		Namespace:         "default",
		Type:              "system",
		Description:       "desc",
		CreationTimestamp: time.Unix(1000, 1000),
		Labels:            map[string]string{"label": "aaa"},
		Selector:          "a=b",
		NodeSelector:      "1=2",
		HostNetwork:       true,
		Replica:           2,
		Workload:          "deployment",
		JobConfig: &specV1.AppJobConfig{
			Completions:   1,
			Parallelism:   2,
			BackoffLimit:  3,
			RestartPolicy: "Never",
		},
		InitServices: []specV1.Service{
			{
				Name:  "init",
				Image: "init_image",
			},
		},
		Services: []specV1.Service{{
			Name:     "test",
			Hostname: "hostname",
			Image:    "test_image",
			Replica:  1,
			VolumeMounts: []specV1.VolumeMount{{
				Name:      "mount",
				MountPath: "path",
				ReadOnly:  false,
			}},
			Ports: []specV1.ContainerPort{{
				HostPort:      1000,
				ContainerPort: 1000,
			}},
			Devices: []specV1.Device{{DevicePath: "dev"}},
			Args:    []string{"arg"},
			Env: []specV1.Environment{{
				Name:  "key",
				Value: "value",
			}},
			Resources: &specV1.Resources{},
			Runtime:   "runtime",
		}},
		Volumes: []specV1.Volume{{
			Name: "test",
			VolumeSource: specV1.VolumeSource{
				HostPath: &specV1.HostPathVolumeSource{Path: "hostPath"},
				Config: &specV1.ObjectReference{
					Name: "config",
				},
			},
		}},
	}
	app3 := &specV1.Application{
		Name:              "app_test",
		Namespace:         "default",
		Type:              "system",
		Description:       "desc",
		CreationTimestamp: time.Unix(1000, 1000),
		Labels:            map[string]string{"label": "bbb"},
		Selector:          "a=b",
		NodeSelector:      "1=2",
		HostNetwork:       true,
		Replica:           2,
		Workload:          "deployment",
		JobConfig: &specV1.AppJobConfig{
			Completions:   1,
			Parallelism:   2,
			BackoffLimit:  3,
			RestartPolicy: "Never",
		},
		InitServices: []specV1.Service{
			{
				Name:  "init",
				Image: "init_image",
			},
		},
		Services: []specV1.Service{{
			Name:     "test",
			Hostname: "hostname",
			Image:    "test_image",
			Replica:  1,
			VolumeMounts: []specV1.VolumeMount{{
				Name:      "mount",
				MountPath: "path",
				ReadOnly:  false,
			}},
			Ports: []specV1.ContainerPort{{
				HostPort:      1000,
				ContainerPort: 1000,
			}},
			Devices: []specV1.Device{{DevicePath: "dev"}},
			Args:    []string{"arg"},
			Env: []specV1.Environment{{
				Name:  "key",
				Value: "value",
			}},
			Resources: &specV1.Resources{},
			Runtime:   "runtime",
		}},
		Volumes: []specV1.Volume{{
			Name: "test",
			VolumeSource: specV1.VolumeSource{
				HostPath: &specV1.HostPathVolumeSource{Path: "hostPath"},
				Config: &specV1.ObjectReference{
					Name: "config",
				},
			},
		}},
	}
	app4 := &specV1.Application{
		Name:              "app_testabc",
		Namespace:         "default",
		Type:              "system",
		Description:       "desc",
		CreationTimestamp: time.Unix(1000, 1000),
		Labels:            map[string]string{"label": "bbb"},
		Selector:          "a=b",
		NodeSelector:      "1=2",
		HostNetwork:       true,
		Replica:           2,
		Workload:          "deployment",
		JobConfig: &specV1.AppJobConfig{
			Completions:   1,
			Parallelism:   2,
			BackoffLimit:  3,
			RestartPolicy: "Never",
		},
		InitServices: []specV1.Service{
			{
				Name:  "init",
				Image: "init_image",
			},
		},
		Services: []specV1.Service{{
			Name:     "test",
			Hostname: "hostname",
			Image:    "test_image",
			Replica:  1,
			VolumeMounts: []specV1.VolumeMount{{
				Name:      "mount",
				MountPath: "path",
				ReadOnly:  false,
			}},
			Ports: []specV1.ContainerPort{{
				HostPort:      1000,
				ContainerPort: 1000,
			}},
			Devices: []specV1.Device{{DevicePath: "dev"}},
			Args:    []string{"arg"},
			Env: []specV1.Environment{{
				Name:  "key",
				Value: "value",
			}},
			Resources: &specV1.Resources{},
			Runtime:   "runtime",
		}},
		Volumes: []specV1.Volume{{
			Name: "test",
			VolumeSource: specV1.VolumeSource{
				HostPath: &specV1.HostPathVolumeSource{Path: "hostPath"},
				Config: &specV1.ObjectReference{
					Name: "config",
				},
			},
		}},
		Ota: specV1.OtaInfo{},
		AutoScaleCfg: &specV1.AutoScaleCfg{
			MinReplicas: 1,
			MaxReplicas: 10,
			Metrics: []specV1.MetricSpec{
				{
					Type: "Resource",
					Resource: &specV1.ResourceMetric{
						Name:               "cpu",
						TargetType:         "Utilization",
						AverageUtilization: 50,
					},
				},
			},
		},
	}

	listOptions := &models.ListOptions{}

	db, err := MockNewDB()
	if err != nil {
		fmt.Printf("get mock sqlite3 error = %s", err.Error())
		t.Fail()
		return
	}
	db.MockCreateAppTable()

	res, err := db.CreateApplication(nil, "default", app1)
	assert.NoError(t, err)
	checkApp(t, app1, res)

	res, err = db.CreateApplication(nil, "default", app2)
	assert.NoError(t, err)
	checkApp(t, app2, res)

	res, err = db.CreateApplication(nil, "default", app3)
	assert.NoError(t, err)
	checkApp(t, app3, res)

	res, err = db.CreateApplication(nil, "default", app4)
	assert.NoError(t, err)
	checkApp(t, app4, res)

	// list option nil, return all apps
	resList, err := db.ListApplication(nil, "default", listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 4)
	assert.Equal(t, app1.Name, resList.Items[0].Name)
	assert.Equal(t, app2.Name, resList.Items[1].Name)
	assert.Equal(t, app3.Name, resList.Items[2].Name)
	assert.Equal(t, app4.Name, resList.Items[3].Name)
	// page 1 num 2
	listOptions.PageNo = 1
	listOptions.PageSize = 2
	resList, err = db.ListApplication(nil, "default", listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 4)
	assert.Equal(t, app1.Name, resList.Items[0].Name)
	assert.Equal(t, app2.Name, resList.Items[1].Name)
	// page 2 num 2
	listOptions.PageNo = 2
	listOptions.PageSize = 2
	resList, err = db.ListApplication(nil, "default", listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 4)
	assert.Equal(t, app3.Name, resList.Items[0].Name)
	assert.Equal(t, app4.Name, resList.Items[1].Name)
	// page 3 num 0
	listOptions.PageNo = 3
	listOptions.PageSize = 2
	resList, err = db.ListApplication(nil, "default", listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 4)
	// page 1 num 2 name like app
	listOptions.PageNo = 1
	listOptions.PageSize = 2
	listOptions.Name = "app"
	resList, err = db.ListApplication(nil, "default", listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 4)
	assert.Equal(t, app1.Name, resList.Items[0].Name)
	assert.Equal(t, app2.Name, resList.Items[1].Name)
	// page 1 num 2 name like abc
	listOptions.PageNo = 1
	listOptions.PageSize = 2
	listOptions.Name = "abc"
	resList, err = db.ListApplication(nil, "default", listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 2)
	assert.Equal(t, app2.Name, resList.Items[0].Name)
	assert.Equal(t, app4.Name, resList.Items[1].Name)
	// page 1 num2 label : aaa
	listOptions.PageNo = 1
	listOptions.PageSize = 4
	listOptions.Name = ""
	listOptions.LabelSelector = "label=aaa"
	resList, err = db.ListApplication(nil, "default", listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 2)
	assert.Equal(t, app1.Name, resList.Items[0].Name)
	assert.Equal(t, app2.Name, resList.Items[1].Name)
	// list by name
	apps, num, err := db.ListApplicationsByNames(nil, "default", []string{"app_test"})
	assert.NoError(t, err)
	assert.Equal(t, 1, num)
	assert.Equal(t, "app_test", apps[0].Name)

	listOptions.PageNo = 1
	listOptions.PageSize = 4
	listOptions.Name = "abc"
	listOptions.LabelSelector = "label=aaa"
	resList, err = db.ListApplication(nil, "default", listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 1)
	assert.Equal(t, app2.Name, resList.Items[0].Name)

	err = db.DeleteApplication(nil, "default", app1.Name)
	assert.NoError(t, err)
	err = db.DeleteApplication(nil, "default", app2.Name)
	assert.NoError(t, err)
	err = db.DeleteApplication(nil, "default", app3.Name)
	assert.NoError(t, err)
	err = db.DeleteApplication(nil, "default", app4.Name)
	assert.NoError(t, err)

	res, err = db.GetApplication(nil, "default", app1.Name, "")
	assert.Nil(t, res)
	res, err = db.GetApplication(nil, "default", app2.Name, "")
	assert.Nil(t, res)
	res, err = db.GetApplication(nil, "default", app3.Name, "")
	assert.Nil(t, res)
	res, err = db.GetApplication(nil, "default", app4.Name, "")
	assert.Nil(t, res)
}

func checkApp(t *testing.T, expect, actual *specV1.Application) {
	assert.Equal(t, expect.Name, actual.Name)
	assert.Equal(t, expect.Namespace, actual.Namespace)
	assert.Equal(t, expect.System, actual.System)
	assert.Equal(t, expect.Description, actual.Description)
	assert.EqualValues(t, expect.Labels, actual.Labels)
	assert.EqualValues(t, expect.JobConfig, actual.JobConfig)
	assert.Equal(t, expect.Replica, actual.Replica)
	assert.Equal(t, expect.HostNetwork, actual.HostNetwork)
	assert.Equal(t, expect.Workload, actual.Workload)
	assert.Equal(t, expect.Selector, actual.Selector)
	assert.Equal(t, expect.NodeSelector, actual.NodeSelector)
	assert.Equal(t, expect.InitServices, actual.InitServices)
	assert.Equal(t, expect.Services, actual.Services)
	assert.Equal(t, expect.Volumes, actual.Volumes)
	assert.EqualValues(t, expect.Ota, actual.Ota)
}
