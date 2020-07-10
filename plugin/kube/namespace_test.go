package kube

import (
	"github.com/baetyl/baetyl-go/log"
	"testing"

	"github.com/baetyl/baetyl-cloud/models"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
)

func genNamespaceRuntime() []runtime.Object {
	rs := []runtime.Object{
		&v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: "default",
			},
		},
	}
	return rs
}

func initNamespaceClient() *client {
	fc := fake.NewSimpleClientset(genNamespaceRuntime()...)
	return &client{
		coreV1: fc.CoreV1(),
		log:    log.With(log.Any("plugin", "kube")),
	}
}

func TestGetNamespace(t *testing.T) {
	c := initNamespaceClient()
	namespace, err := c.GetNamespace("default")
	assert.NoError(t, err)
	assert.Equal(t, "default", namespace.Name)

	_, err = c.GetNamespace("test")
	assert.Error(t, err)
}

func TestCreateNamespace(t *testing.T) {
	c := initNamespaceClient()
	ns := &models.Namespace{
		Name: "default",
	}
	_, err := c.CreateNamespace(ns)
	assert.Error(t, err)

	ns.Name = "test"
	namespace, err := c.CreateNamespace(ns)
	assert.NoError(t, err)
	assert.Equal(t, "test", namespace.Name)
}

func TestDeleteNamespace(t *testing.T) {
	c := initNamespaceClient()
	ns := &models.Namespace{
		Name: "default",
	}
	err := c.DeleteNamespace(ns)
	assert.NoError(t, err)

	ns.Name = "test"
	err = c.DeleteNamespace(ns)
	assert.Error(t, err)
}
