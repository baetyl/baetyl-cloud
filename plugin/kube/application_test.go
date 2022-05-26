package kube

import (
	"testing"

	"github.com/baetyl/baetyl-go/v2/log"

	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/stretchr/testify/assert"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin/kube/apis/cloud/v1alpha1"
	"github.com/baetyl/baetyl-cloud/v2/plugin/kube/client/clientset/versioned/fake"
)

func TestConvert(t *testing.T) {
	app := &v1alpha1.Application{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:       "default",
			Name:            "test_name",
			ResourceVersion: "1",
		},
		Spec: v1alpha1.ApplicationSpec{
			Selector: "a=a",
			Services: []v1alpha1.Service{{
				Container: v1.Container{
					Name:         "test",
					Image:        "image",
					Args:         []string{"arg"},
					Ports:        []v1.ContainerPort{{HostPort: 1000, ContainerPort: 1000}},
					Env:          []v1.EnvVar{{Name: "key", Value: "value"}},
					VolumeMounts: []v1.VolumeMount{{Name: "mount", MountPath: "path", ReadOnly: false}},
				},
				Devices: []v1alpha1.Device{{DevicePath: "dev"}},
				Labels: map[string]string{
					"tag": "function",
				},
				Hostname:  "hostname",
				Replica:   1,
				Resources: &v1alpha1.Resources{},
				Runtime:   "runtime",
			}},
			Volumes: []v1alpha1.Volume{{
				Name:     "test",
				HostPath: &v1alpha1.HostPathVolumeSource{Path: "hostPath"},
				Config:   &v1alpha1.ObjectReference{Name: "config"},
			}},
		},
	}

	expected := &specV1.Application{
		Name:      "test_name",
		Namespace: "default",
		Version:   "1",
		Selector:  "a=a",
		Services: []specV1.Service{{
			Name: "test",
			Labels: map[string]string{
				"tag": "function",
			},
			Hostname: "hostname",
			Image:    "image",
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
	res := toAppModel(app)
	assert.Equal(t, expected, res)
}

func TestGenerate(t *testing.T) {
	app := &specV1.Application{
		Name: "test_name",
		//Version:      "1",
		Selector: "a=a",
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
	expected := &v1alpha1.Application{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Application",
			APIVersion: v1alpha1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace:   "default",
			Name:        "test_name",
			Annotations: map[string]string{},
			//ResourceVersion: "1",
		},
		Spec: v1alpha1.ApplicationSpec{
			Selector: "a=a",
			Services: []v1alpha1.Service{{
				Container: v1.Container{
					Name:         "test",
					Image:        "test_image",
					Args:         []string{"arg"},
					Ports:        []v1.ContainerPort{{HostPort: 1000, ContainerPort: 1000}},
					Env:          []v1.EnvVar{{Name: "key", Value: "value"}},
					VolumeMounts: []v1.VolumeMount{{Name: "mount", MountPath: "path", ReadOnly: false}},
				},
				Devices:   []v1alpha1.Device{{DevicePath: "dev"}},
				Hostname:  "hostname",
				Replica:   1,
				Resources: &v1alpha1.Resources{},
				Runtime:   "runtime",
			}},
			Volumes: []v1alpha1.Volume{{
				Name:     "test",
				HostPath: &v1alpha1.HostPathVolumeSource{Path: "hostPath"},
				Config:   &v1alpha1.ObjectReference{Name: "config"},
			}},
		},
	}
	res := fromAppModel("default", app)
	assert.Equal(t, expected, res)
}

func genApplicationRuntime() []runtime.Object {
	rs := []runtime.Object{
		&v1alpha1.Application{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Application",
				APIVersion: v1alpha1.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "default",
				Name:      "test_name",
				//ResourceVersion: "1",
			},
			Spec: v1alpha1.ApplicationSpec{
				Selector: "a=a",
				Services: []v1alpha1.Service{{
					Container: v1.Container{
						Name:          "test",
						Image:         "test_image",
						Args:          []string{"arg"},
						Ports:         []v1.ContainerPort{{HostPort: 1000, ContainerPort: 1000}},
						Env:           []v1.EnvVar{{Name: "key", Value: "value"}},
						VolumeMounts:  []v1.VolumeMount{{Name: "mount", MountPath: "path", ReadOnly: false}},
						VolumeDevices: []v1.VolumeDevice{{DevicePath: "dev"}},
					},
					Hostname:  "hostname",
					Replica:   1,
					Resources: &v1alpha1.Resources{},
					Runtime:   "runtime",
				},
				},
				Volumes: []v1alpha1.Volume{{
					Name:     "test",
					HostPath: &v1alpha1.HostPathVolumeSource{Path: "hostPath"},
					Config: &v1alpha1.ObjectReference{
						Name: "config",
					},
				}},
			},
		},
	}
	return rs
}

func initApplicationClient() *client {
	fc := fake.NewSimpleClientset(genApplicationRuntime()...)
	return &client{
		customClient: fc,
		log:          log.With(log.Any("plugin", "kube")),
	}
}

func TestGetApplication(t *testing.T) {
	c := initApplicationClient()
	_, err := c.GetApplication("default", "test", "")
	assert.NotNil(t, err)
	cfg, err := c.GetApplication("default", "test_name", "")
	assert.Equal(t, cfg.Name, "test_name")
}

func TestCreateApplication(t *testing.T) {
	c := initApplicationClient()
	cfg := &specV1.Application{
		Name:      "test_name_2",
		Namespace: "default",
	}
	cfg2, err := c.CreateApplication(nil, cfg.Namespace, cfg)
	assert.NoError(t, err)
	assert.Equal(t, cfg.Name, cfg2.Name)

	cfg.Name = "test_name"
	_, err = c.CreateApplication(nil, cfg.Namespace, cfg)
	assert.NotNil(t, err)
}

func TestUpdateApplication(t *testing.T) {
	c := initApplicationClient()
	cfg := &specV1.Application{
		Name:      "test_name",
		Namespace: "default",
	}
	cfg2, err := c.UpdateApplication(nil, cfg.Namespace, cfg)
	assert.NoError(t, err)
	assert.Equal(t, cfg.Name, cfg2.Name)
	assert.Equal(t, 0, len(cfg2.Services))

	cfg.Name = cfg.Name + "NULL"
	_, err = c.UpdateApplication(nil, cfg.Namespace, cfg)
	assert.NotNil(t, err)
}

func TestDeleteApplication(t *testing.T) {
	c := initApplicationClient()
	err := c.DeleteApplication(nil, "default", "test_name")
	assert.NoError(t, err)
}

func TestListListApplication(t *testing.T) {
	c := initApplicationClient()
	l, err := c.ListApplication(nil, "default", &models.ListOptions{})
	assert.NoError(t, err)
	assert.Equal(t, l.Total, 1)
}

func TestListListApplicationByName(t *testing.T) {
	c := initApplicationClient()
	l, num, err := c.ListApplicationsByNames(nil, "default", []string{"test_name"})
	assert.NoError(t, err)
	assert.Equal(t, 1, num)
	assert.Equal(t, "test_name", l[0].Name)

	// case 2
	l, num, err = c.ListApplicationsByNames(nil, "default", []string{"no"})
	assert.NoError(t, err)
	assert.Equal(t, 0, num)
}
