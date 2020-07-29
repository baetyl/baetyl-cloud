package service

import (
	"testing"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	mockPlugin "github.com/baetyl/baetyl-cloud/v2/mock/plugin"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type MockServices struct {
	conf           *config.CloudConfig
	ctl            *gomock.Controller
	modelStorage   *mockPlugin.MockModelStorage
	dbStorage      *mockPlugin.MockDBStorage
	objectStorage  *mockPlugin.MockObject
	functionPlugin *mockPlugin.MockFunction
	pki            *mockPlugin.MockPKI
	auth           *mockPlugin.MockAuth
	shadowStorage  *mockPlugin.MockShadow
	license        *mockPlugin.MockLicense
	property       *mockPlugin.MockProperty
}

func (m *MockServices) Close() {
	if m.ctl != nil {
		m.ctl.Finish()
	}
}

func mockStorageModel(mock plugin.ModelStorage) plugin.Factory {
	factory := func() (plugin.Plugin, error) {
		return mock, nil
	}
	return factory
}

func mockStorageDB(mock plugin.DBStorage) plugin.Factory {
	factory := func() (plugin.Plugin, error) {
		return mock, nil
	}
	return factory
}

func mockStorageObject(mock plugin.Object) plugin.Factory {
	factory := func() (plugin.Plugin, error) {
		return mock, nil
	}
	return factory
}

func mockPKI(mock plugin.PKI) plugin.Factory {
	factory := func() (plugin.Plugin, error) {
		return mock, nil
	}
	return factory
}

func mockAuth(mock plugin.Auth) plugin.Factory {
	factory := func() (plugin.Plugin, error) {
		return mock, nil
	}
	return factory
}

func mockFunction(mock plugin.Function) plugin.Factory {
	factory := func() (plugin.Plugin, error) {
		return mock, nil
	}
	return factory
}

func mockShadowStorage(mock plugin.Shadow) plugin.Factory {
	ss := func() (plugin.Plugin, error) {
		return mock, nil
	}
	return ss
}

func mockLicense(mock plugin.License) plugin.Factory {
	qc := func() (plugin.Plugin, error) {
		return mock, nil
	}
	return qc
}

func mockTestConfig() *config.CloudConfig {
	conf := &config.CloudConfig{}
	conf.Plugin.ModelStorage = common.RandString(9)
	conf.Plugin.DatabaseStorage = common.RandString(9)
	conf.Plugin.Objects = []string{common.RandString(9)}
	conf.Plugin.PKI = common.RandString(9)
	conf.Plugin.Auth = common.RandString(9)
	conf.Plugin.Functions = []string{common.RandString(9)}
	conf.Plugin.Shadow = conf.Plugin.DatabaseStorage
	conf.Plugin.License = common.RandString(9)
	conf.Plugin.Property = common.RandString(9)
	return conf
}

func mockEmptyTestConfig() *config.CloudConfig {
	conf := &config.CloudConfig{}
	conf.Plugin.Objects = []string{}
	conf.Plugin.Functions = []string{}
	return conf
}

func InitMockEnvironment(t *testing.T) *MockServices {
	conf := mockTestConfig()
	mockCtl := gomock.NewController(t)
	mockModelStorage := mockPlugin.NewMockModelStorage(mockCtl)
	plugin.RegisterFactory(conf.Plugin.ModelStorage, mockStorageModel(mockModelStorage))
	mockDBStorage := mockPlugin.NewMockDBStorage(mockCtl)
	plugin.RegisterFactory(conf.Plugin.DatabaseStorage, mockStorageDB(mockDBStorage))
	mPKI := mockPlugin.NewMockPKI(mockCtl)
	plugin.RegisterFactory(conf.Plugin.PKI, mockPKI(mPKI))
	mAuth := mockPlugin.NewMockAuth(mockCtl)
	plugin.RegisterFactory(conf.Plugin.Auth, mockAuth(mAuth))
	mockObjectStorage := mockPlugin.NewMockObject(mockCtl)
	for _, v := range conf.Plugin.Objects {
		plugin.RegisterFactory(v, mockStorageObject(mockObjectStorage))
	}
	mockFunctionPlugin := mockPlugin.NewMockFunction(mockCtl)
	for _, v := range conf.Plugin.Functions {
		plugin.RegisterFactory(v, mockFunction(mockFunctionPlugin))
	}

	mLicense := mockPlugin.NewMockLicense(mockCtl)
	plugin.RegisterFactory(conf.Plugin.License, mockLicense(mLicense))
	_, err := NewSyncService(conf)
	assert.Nil(t, err)

	mProperty := mockPlugin.NewMockProperty(mockCtl)
	plugin.RegisterFactory(conf.Plugin.Property, mockProperty(mProperty))
	return &MockServices{
		conf:           conf,
		ctl:            mockCtl,
		modelStorage:   mockModelStorage,
		dbStorage:      mockDBStorage,
		objectStorage:  mockObjectStorage,
		functionPlugin: mockFunctionPlugin,
		pki:            mPKI,
		auth:           mAuth,
		license:        mLicense,
		property:       mProperty,
	}
}

func InitEmptyMockEnvironment(t *testing.T) *MockServices {
	conf := mockEmptyTestConfig()
	mockCtl := gomock.NewController(t)
	mockObjectStorage := mockPlugin.NewMockObject(mockCtl)
	for _, v := range conf.Plugin.Objects {
		plugin.RegisterFactory(v, mockStorageObject(mockObjectStorage))
	}
	mockFunctionPlugin := mockPlugin.NewMockFunction(mockCtl)
	for _, v := range conf.Plugin.Functions {
		plugin.RegisterFactory(v, mockFunction(mockFunctionPlugin))
	}

	return &MockServices{
		conf:           conf,
		ctl:            mockCtl,
		objectStorage:  mockObjectStorage,
		functionPlugin: mockFunctionPlugin,
	}
}

func mockProperty(mock plugin.Property) plugin.Factory {
	factory := func() (plugin.Plugin, error) {
		return mock, nil
	}
	return factory
}
