package service

import (
	"testing"

	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

func genSecretTestCase() *specV1.Secret {
	r := &specV1.Secret{
		Namespace: "default",
		Name:      "abc",
	}
	return r
}

func genSecretLitsTestCase() *models.SecretList {
	l := &models.SecretList{
		Total: 0,
		Items: make([]specV1.Secret, 0),
	}
	return l
}

func TestDefaultRegistryService_Get(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	registry := genSecretTestCase()
	mockObject.secret.EXPECT().GetSecret(gomock.Any(), gomock.Any(), gomock.Any(), "").Return(genSecretTestCase(), nil).AnyTimes()
	cs, err := NewSecretService(mockObject.conf)
	assert.NoError(t, err)
	_, err = cs.Get(registry.Namespace, registry.Name, "")
	assert.NoError(t, err)
	_, err = cs.Get(registry.Namespace, registry.Name, "")
	assert.NoError(t, err)
}

func TestDefaultRegistryService_List(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	ns, s := "default", &models.ListOptions{}
	sl := genSecretLitsTestCase()
	mockObject.secret.EXPECT().ListSecret(ns, s).Return(sl, nil).AnyTimes()
	cs, err := NewSecretService(mockObject.conf)
	assert.NoError(t, err)
	_, err = cs.List(ns, s)
	assert.NoError(t, err)
}

func TestDefaultRegistryService_Delete(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	cs, err := NewSecretService(mockObject.conf)
	assert.NoError(t, err)
	registry := genSecretTestCase()
	mockObject.secret.EXPECT().DeleteSecret(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	err = cs.Delete(registry.Namespace, registry.Name)
	assert.NoError(t, err)
}
func TestDefaultRegistryService_Create(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	cs, err := NewSecretService(mockObject.conf)
	assert.NoError(t, err)
	registry := genSecretTestCase()
	mockObject.secret.EXPECT().CreateSecret(gomock.Any(), gomock.Any(), gomock.Any()).Return(genSecretTestCase(), nil).AnyTimes()
	_, err = cs.Create(nil, registry.Namespace, registry)
	assert.NoError(t, err)
}
func TestDefaultRegistryService_Update(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	cs, err := NewSecretService(mockObject.conf)
	assert.NoError(t, err)
	registry := genSecretTestCase()
	mockObject.secret.EXPECT().UpdateSecret(gomock.Any(), gomock.Any()).Return(genSecretTestCase(), nil)
	_, err = cs.Update(registry.Namespace, registry)
	assert.NoError(t, err)
}
