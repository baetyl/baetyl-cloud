package service

import (
	"errors"
	"testing"

	"github.com/baetyl/baetyl-cloud/models"
	"github.com/stretchr/testify/assert"
)

func TestDefaultFunctionService_List(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	functions := []models.Function{
		{
			Name:    "test1",
			Handler: "index.handler1",
			Version: "latest",
			Runtime: "python3",
		},
		{
			Name:    "test2",
			Handler: "index.handler2",
			Version: "latest",
			Runtime: "nodejs10",
		},
		{
			Name:    "test3",
			Handler: "index.handler3",
			Version: "latest",
			Runtime: "other",
		},
	}

	namespace := "default"

	mockObject.functionPlugin.EXPECT().List(namespace).Return(functions, nil)
	cs, err := NewFunctionService(mockObject.conf)
	assert.NoError(t, err)
	res, err := cs.List(namespace, mockObject.conf.Plugin.Functions[0])
	assert.NoError(t, err)
	assert.Equal(t, len(res), 3)
	assert.Equal(t, functions[0].Name, res[0].Name)
	assert.Equal(t, functions[0].Handler, res[0].Handler)
	assert.Equal(t, functions[0].Version, res[0].Version)
	assert.Equal(t, functions[0].Runtime, res[0].Runtime)
	assert.Equal(t, functions[1].Name, res[1].Name)
	assert.Equal(t, functions[0].Handler, res[0].Handler)
	assert.Equal(t, functions[0].Version, res[0].Version)
	assert.Equal(t, functions[0].Runtime, res[0].Runtime)
}

func TestDefaultFunctionService_ListFunctionVersions(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	functions := []models.Function{
		{
			Name:    "test1",
			Version: "v1",
		},
		{
			Name:    "test1",
			Version: "v2",
		},
	}

	namespace := "default"
	name := "test1"
	mockObject.functionPlugin.EXPECT().ListFunctionVersions(namespace, name).Return(functions, nil).Times(1)
	cs, err := NewFunctionService(mockObject.conf)
	assert.NoError(t, err)
	res, err := cs.ListFunctionVersions(namespace, name, mockObject.conf.Plugin.Functions[0])
	assert.NoError(t, err)
	assert.Equal(t, len(functions), len(res))
	assert.Equal(t, functions[0].Name, res[0].Name)
	assert.Equal(t, functions[0].Version, res[0].Version)
	assert.Equal(t, functions[1].Name, res[1].Name)
	assert.Equal(t, functions[0].Version, res[0].Version)

	name2 := "test2"
	mockObject.functionPlugin.EXPECT().ListFunctionVersions(namespace, name2).Return(nil, errors.New("err")).Times(1)
	_, err2 := cs.ListFunctionVersions(namespace, name2, mockObject.conf.Plugin.Functions[0])
	assert.Error(t, err2)
	assert.Equal(t, err2.Error(), "err")
}

func TestDefaultFunctionService_ListSources(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	cs, err := NewFunctionService(mockObject.conf)
	assert.NoError(t, err)
	res := cs.ListSources()
	assert.NotNil(t, res)
	assert.Equal(t, len(res), 1)
}

func TestDefaultFunctionService_ListSourcesWithEmptySource(t *testing.T) {
	mockObject := InitEmptyMockEnvironment(t)
	defer mockObject.Close()

	cs, err := NewFunctionService(mockObject.conf)
	assert.NoError(t, err)
	res := cs.ListSources()
	assert.NotNil(t, res)
	assert.Equal(t, len(res), 0)
}

func TestDefaultFunctionService_GetFunction(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	function := models.Function{
		Name:    "test0",
		Handler: "index.handler",
		Version: "v1",
		Runtime: "python36",
		Code: models.FunctionCode{
			Size:     120,
			Sha256:   "sha",
			Location: "func.zip",
		},
	}

	namespace, name, version := "default", "test1", "v1"
	mockObject.functionPlugin.EXPECT().Get(namespace, name, version).Return(&function, nil).Times(1)

	cs, err := NewFunctionService(mockObject.conf)
	assert.NoError(t, err)
	res, err := cs.GetFunction(namespace, name, version, mockObject.conf.Plugin.Functions[0])
	assert.NoError(t, err)
	assert.Equal(t, *res, function)

	name2, version2 := "test2", "v2"
	mockObject.functionPlugin.EXPECT().Get(namespace, name2, version2).Return(nil, errors.New("err")).Times(1)
	_, err2 := cs.GetFunction(namespace, name2, version2, mockObject.conf.Plugin.Functions[0])
	assert.Error(t, err2)
	assert.Equal(t, err2.Error(), "err")
}
