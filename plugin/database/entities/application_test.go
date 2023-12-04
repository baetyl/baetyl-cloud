// Package entities 数据库存储基本结构与方法
package entities

import (
	"testing"
	"time"

	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/common"
)

func TestFromApplicationModel(t *testing.T) {
	app := &Application{
		Name:        "testApp",
		Namespace:   "namespace",
		Labels:      `{"baetyl-cloud-system":"true"}`,
		HostNetwork: true,
		Workload:    "deployment",
		Replica:     1,
		JobConfig:   `{"completions":1,"parallelism":2,"backoffLimit":3,"restartPolicy":"Never"}`,
	}

	mApp := &specV1.Application{
		Namespace: "namespace",
		Name:      "testApp",
		Labels: map[string]string{
			common.LabelSystem: "true",
		},
		HostNetwork: true,
		Workload:    "deployment",
		Replica:     1,
		JobConfig: &specV1.AppJobConfig{
			Completions:   1,
			Parallelism:   2,
			BackoffLimit:  3,
			RestartPolicy: "Never",
		},
	}
	modelApp, err := FromAppModel("namespace", mApp)

	assert.NoError(t, err)
	assert.Equal(t, app.Name, modelApp.Name)
	assert.Equal(t, app.Namespace, modelApp.Namespace)
	assert.Equal(t, app.Labels, modelApp.Labels)
}

func TestEqualAppFail(t *testing.T) {
	app1 := &specV1.Application{
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
		InitServices: []specV1.Service{{
			Name:  "init",
			Image: "init_image",
		}},
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
		InitServices: []specV1.Service{{
			Name:  "init",
			Image: "init_image",
		}},
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

	app1.Volumes[0].HostPath.Path = "hostpath1"
	app2.Volumes[0].HostPath.Path = "hostpath2"
	flag := EqualApp(app1, app2)
	assert.Equal(t, flag, false)

	app1.Volumes[0].HostPath.Path = "hostpath"
	app2.Volumes[0].HostPath.Path = "hostpath"

	app1.Services[0].Hostname = "hostname1"
	app2.Services[0].Hostname = "hostname2"
	flag1 := EqualApp(app1, app2)
	assert.Equal(t, flag1, false)

	app1.Services[0].Hostname = "hostname"
	app2.Services[0].Hostname = "hostname"

	app1.InitServices[0].Name = "init1"
	app2.InitServices[0].Name = "init2"
	flag2 := EqualApp(app1, app2)
	assert.Equal(t, flag2, false)
}
