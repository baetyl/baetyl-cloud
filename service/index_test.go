package service

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/common"
)

func TestDefaultIndexService_RefreshIndex(t *testing.T) {
	namespace := "default"
	mockObject := InitMockEnvironment(t)
	mockObject.index.EXPECT().RefreshIndex(nil, namespace, common.Config, common.Application, "123", []string{}).Return(nil).AnyTimes()
	is, err := NewIndexService(mockObject.conf)
	assert.NoError(t, err)
	err = is.RefreshIndex(nil, namespace, common.Config, common.Application, "123", []string{})
	assert.NoError(t, err)
	mockObject.Close()

	mockObject = InitMockEnvironment(t)
	mockObject.index.EXPECT().RefreshIndex(nil, namespace, common.Config, common.Node, "123", []string{}).Return(fmt.Errorf("delete : table not exist")).AnyTimes()
	is, err = NewIndexService(mockObject.conf)
	assert.NoError(t, err)
	err = is.RefreshIndex(nil, namespace, common.Config, common.Node, "123", []string{})
	assert.Error(t, err, "delete : table not exist")
	mockObject.Close()

	mockObject = InitMockEnvironment(t)
	mockObject.index.EXPECT().RefreshIndex(nil, namespace, common.Config, common.Node, "123", []string{}).Return(fmt.Errorf("create : table not exist")).AnyTimes()
	is, err = NewIndexService(mockObject.conf)
	assert.NoError(t, err)
	err = is.RefreshIndex(nil, namespace, common.Config, common.Node, "123", []string{})
	assert.Error(t, err, "create : table not exist")
	mockObject.Close()
}

func TestDefaultIndexService_ListIndex(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	data := "to list by"
	namespace := "default"

	mockObject.index.EXPECT().ListIndex(namespace, common.Application, common.Config, data).Return([]string{}, nil).AnyTimes()
	is, err := NewIndexService(mockObject.conf)
	assert.NoError(t, err)
	_, err = is.ListIndex(namespace, common.Application, common.Config, data)
	assert.NoError(t, err)
}

func TestResourceList(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	data := "to update"
	namespace := "default"

	mockObject.index.EXPECT().ListIndex(namespace, gomock.Any(), gomock.Any(), data).Return([]string{}, nil).AnyTimes()
	is, err := NewIndexService(mockObject.conf)
	assert.NoError(t, err)
	_, err = is.ListIndex(namespace, common.Application, common.Config, data)
	assert.NoError(t, err)
	_, err = is.ListAppIndexByConfig(namespace, data)
	assert.NoError(t, err)
	_, err = is.ListAppIndexBySecret(namespace, data)
	assert.NoError(t, err)
	_, err = is.ListConfigIndexByApp(namespace, data)
	assert.NoError(t, err)
	_, err = is.ListNodesByApp(namespace, data)
	assert.NoError(t, err)
	_, err = is.ListAppsByNode(namespace, data)
	assert.NoError(t, err)
	_, err = is.ListSecretIndexByApp(namespace, data)
	assert.NoError(t, err)
}

func TestResourceRefresh(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	data := "to update"
	arr := []string{"r0", "r1", "r2"}
	namespace := "default"

	mockObject.index.EXPECT().RefreshIndex(nil, namespace, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	is, err := NewIndexService(mockObject.conf)
	assert.NoError(t, err)
	err = is.RefreshAppIndexByConfig(nil, namespace, data, arr)
	assert.NoError(t, err)
	err = is.RefreshConfigIndexByApp(nil, namespace, data, arr)
	assert.NoError(t, err)

	err = is.RefreshNodesIndexByApp(nil, namespace, data, arr)

	err = is.RefreshSecretIndexByApp(nil, namespace, data, arr)
	assert.NoError(t, err)

	err = is.RefreshAppsIndexByNode(nil, namespace, data, arr)
	assert.NoError(t, err)
}
