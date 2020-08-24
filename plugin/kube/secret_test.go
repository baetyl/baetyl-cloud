package kube

import (
	"github.com/baetyl/baetyl-go/v2/log"
	"testing"

	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin/kube/apis/cloud/v1alpha1"
	"github.com/baetyl/baetyl-cloud/v2/plugin/kube/client/clientset/versioned/fake"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func genSecretRuntime() []runtime.Object {
	rs := []runtime.Object{
		&v1alpha1.Secret{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Secret",
				APIVersion: v1alpha1.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-get",
				Namespace: "default",
				Annotations: map[string]string{
					"prefix/certId": "certId",
				},
			},
		},
		&v1alpha1.Secret{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Secret",
				APIVersion: v1alpha1.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-update",
				Namespace: "default",
			},
		},
		&v1alpha1.Secret{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Secret",
				APIVersion: v1alpha1.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-create",
				Namespace: "default",
			},
		},
		&v1alpha1.Secret{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Secret",
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

func initSecretMapClient() *client {
	fc := fake.NewSimpleClientset(genSecretRuntime()...)
	return &client{
		customClient: fc,
		aesKey:       []byte("0123456789abcdef"),
		log:          log.With(log.Any("plugin", "kube")),
	}
}

func TestGetSecret(t *testing.T) {
	c := initSecretMapClient()
	_, err := c.GetSecret("default", "test", "")
	assert.NotNil(t, err)
	secret, err := c.GetSecret("default", "test-get", "")
	assert.Equal(t, secret.Name, "test-get")
	assert.Equal(t, secret.Annotations["prefix/certId"], "certId")
}

func TestCreateSecret(t *testing.T) {
	c := initSecretMapClient()
	secret := &specV1.Secret{
		Name:      "test-create1",
		Namespace: "default",
		Data: map[string][]byte{
			"key": []byte("value"),
		},
		Labels: map[string]string{
			"Name":      "test-update",
			"Namespace": "default",
		},
		Annotations: map[string]string{
			"Annotations": "test-Annotations",
		},
		System:      true,
		Description: "test",
		Version:     "2",
	}
	secret2, err := c.CreateSecret(secret.Namespace, secret)
	assert.NoError(t, err)
	assert.Equal(t, secret2.Name, secret.Name)
	assert.Equal(t, secret2.Namespace, secret.Namespace)
	assert.Equal(t, secret2.Data, secret.Data)
	assert.Equal(t, secret2.Labels, secret.Labels)
	assert.Equal(t, secret2.System, secret2.System)
	assert.Equal(t, secret2.Annotations, secret.Annotations)
	assert.Equal(t, secret2.CreationTimestamp, secret.CreationTimestamp)
}

func TestUpdateSecret(t *testing.T) {
	c := initSecretMapClient()
	cfg := &specV1.Secret{
		Name:      "test-update",
		Namespace: "default",
		Data: map[string][]byte{
			"service.yml": []byte("test"),
		},
	}
	cfg2, err := c.UpdateSecret(cfg.Namespace, cfg)
	assert.NoError(t, err)
	assert.Equal(t, cfg.Name, cfg2.Name)
	v, _ := cfg2.Data["service.yml"]
	assert.Equal(t, v, []byte("test"))

	cfg.Name = cfg.Name + "NULL"
	_, err = c.UpdateSecret(cfg.Namespace, cfg)
	assert.NotNil(t, err)

}

func TestDeleteSecret(t *testing.T) {
	c := initSecretMapClient()
	err := c.DeleteSecret("default", "test-delete")
	assert.NoError(t, err)
	err = c.DeleteSecret("default", "test-delete")
	assert.NoError(t, err)
}

func TestListSecret(t *testing.T) {
	c := initSecretMapClient()
	l, err := c.ListSecret("default", &models.ListOptions{})
	assert.NoError(t, err)
	assert.Equal(t, l.Total, 4)
}
