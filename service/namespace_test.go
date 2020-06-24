package service

import (
	"fmt"
	"testing"

	"github.com/baetyl/baetyl-cloud/models"
	"github.com/stretchr/testify/assert"
)

func TestNamespaceService_Create(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	name := "user-id-test"
	ns := &models.Namespace{Name: name}
	mockObject.modelStorage.EXPECT().CreateNamespace(ns).Return(ns, nil)
	cs, err := NewNamespaceService(mockObject.conf)
	assert.NoError(t, err)
	res, err := cs.Create(ns)
	assert.NoError(t, err)
	assert.Equal(t, name, res.Name)
}

func TestNamespaceService_Get(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	name := "user-id-test"
	ns := &models.Namespace{Name: name}

	mockObject.modelStorage.EXPECT().GetNamespace(ns.Name).Return(ns, nil)

	cs, err := NewNamespaceService(mockObject.conf)
	assert.NoError(t, err)
	res, err := cs.Get(ns.Name)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	mockObject.modelStorage.EXPECT().GetNamespace(ns.Name).Return(nil, fmt.Errorf("namespaces \"user-id-test\" not found"))
	res, err = cs.Get(ns.Name)
	assert.Error(t, err)
	assert.Equal(t, true, res == nil)
}

func TestNamespaceService_Delete(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	name := "user-id-test"
	ns := &models.Namespace{Name: name}
	mockObject.modelStorage.EXPECT().DeleteNamespace(ns).Return(nil)
	cs, err := NewNamespaceService(mockObject.conf)
	assert.NoError(t, err)
	err = cs.Delete(ns)
	assert.NoError(t, err)
}
