// Package entities 数据库存储基本结构与方法
package entities

import (
	"testing"

	"github.com/baetyl/baetyl-go/v2/context"
	v1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
)

func TestTransDriver(t *testing.T) {
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
	dbDriver, err := FromModelDriver(driver)
	assert.NoError(t, err)

	_, err = ToModelDriver(dbDriver)
	assert.NoError(t, err)

	dbDriver.Labels = "test"
	_, err = ToModelDriver(dbDriver)
	assert.NotNil(t, err)

	dbDriver.Registries = "test"
	_, err = ToModelDriver(dbDriver)
	assert.NotNil(t, err)

	dbDriver.Volumes = "test"
	_, err = ToModelDriver(dbDriver)
	assert.NotNil(t, err)

	dbDriver.Service = "test"
	_, err = ToModelDriver(dbDriver)
	assert.NotNil(t, err)
}
