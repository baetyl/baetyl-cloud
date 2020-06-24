package kube

import (
	"testing"

	"github.com/baetyl/baetyl-cloud/models"
	"github.com/baetyl/baetyl-cloud/plugin/kube/apis/cloud/v1alpha1"
	"github.com/baetyl/baetyl-cloud/plugin/kube/client/clientset/versioned/fake"
	specV1 "github.com/baetyl/baetyl-go/spec/v1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func genConfigRuntime() []runtime.Object {
	rs := []runtime.Object{
		&v1alpha1.Configuration{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Configuration",
				APIVersion: v1alpha1.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-get",
				Namespace: "default",
			},
		},
		&v1alpha1.Configuration{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Configuration",
				APIVersion: v1alpha1.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-update",
				Namespace: "default",
			},
		},
		&v1alpha1.Configuration{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Configuration",
				APIVersion: v1alpha1.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-create",
				Namespace: "default",
			},
		},
		&v1alpha1.Configuration{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Configuration",
				APIVersion: v1alpha1.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-delete",
				Namespace: "default",
			},
		},
	}
	return rs
}

func initConfigurationClient() *client {
	fc := fake.NewSimpleClientset(genConfigRuntime()...)
	return &client{
		customClient: fc,
	}
}

func TestGetConfig(t *testing.T) {
	c := initConfigurationClient()
	_, err := c.GetConfig("default", "test", "")
	assert.NotNil(t, err)
	cfg, err := c.GetConfig("default", "test-get", "")
	assert.Equal(t, cfg.Name, "test-get")
}

func TestCreateConfig(t *testing.T) {
	c := initConfigurationClient()
	cfg := &specV1.Configuration{
		Name:      "test-add",
		Namespace: "default",
		Data: map[string]string{
			"key": "value",
		},
	}
	cfg2, err := c.CreateConfig(cfg.Namespace, cfg)
	assert.NoError(t, err)
	assert.Equal(t, cfg.Name, cfg2.Name)
	assert.Equal(t, "value", cfg2.Data["key"])
}

func TestUpdateConfig(t *testing.T) {
	c := initConfigurationClient()
	cfg := &specV1.Configuration{
		Name:      "test-update",
		Namespace: "default",
		Data: map[string]string{
			"service.yml": "test",
		},
	}
	cfg2, err := c.UpdateConfig(cfg.Namespace, cfg)
	assert.NoError(t, err)
	assert.Equal(t, cfg.Name, cfg2.Name)
	v, _ := cfg2.Data["service.yml"]
	assert.Equal(t, v, "test")

	cfg.Name = cfg.Name + "NULL"
	_, err = c.UpdateConfig(cfg.Namespace, cfg)
	assert.NotNil(t, err)
}

func TestDeleteConfig(t *testing.T) {
	c := initConfigurationClient()
	err := c.DeleteConfig("default", "test-delete")
	assert.NoError(t, err)
}

func TestListConfig(t *testing.T) {
	c := initConfigurationClient()
	l, err := c.ListConfig("default", &models.ListOptions{})
	assert.NoError(t, err)
	assert.Equal(t, l.Total, 4)
}
