package service

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	mockPlugin "github.com/baetyl/baetyl-cloud/v2/mock/plugin"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

type MockServices struct {
	conf           *config.CloudConfig
	ctl            *gomock.Controller
	node           *mockPlugin.MockResource
	namespace      *mockPlugin.MockResource
	configuration  *mockPlugin.MockResource
	secret         *mockPlugin.MockResource
	app            *mockPlugin.MockResource
	index          *mockPlugin.MockIndex
	appHis         *mockPlugin.MockAppHistory
	objectStorage  *mockPlugin.MockObject
	functionPlugin *mockPlugin.MockFunction
	pki            *mockPlugin.MockPKI
	auth           *mockPlugin.MockAuth
	shadow         *mockPlugin.MockShadow
	license        *mockPlugin.MockLicense
	property       *mockPlugin.MockProperty
	module         *mockPlugin.MockModule
	task           *mockPlugin.MockTask
}

func (m *MockServices) Close() {
	if m.ctl != nil {
		m.ctl.Finish()
	}
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

func mockProperty(mock plugin.Property) plugin.Factory {
	factory := func() (plugin.Plugin, error) {
		return mock, nil
	}
	return factory
}

func mockModule(mock plugin.Module) plugin.Factory {
	factory := func() (plugin.Plugin, error) {
		return mock, nil
	}
	return factory
}

func mockResource(mock plugin.Resource) plugin.Factory {
	factory := func() (plugin.Plugin, error) {
		return mock, nil
	}
	return factory
}

func mockIndex(mock plugin.Index) plugin.Factory {
	factory := func() (plugin.Plugin, error) {
		return mock, nil
	}
	return factory
}

func mockAppHis(mock plugin.AppHistory) plugin.Factory {
	factory := func() (plugin.Plugin, error) {
		return mock, nil
	}
	return factory
}

func mockTask(task plugin.Task) plugin.Factory {
	factory := func() (plugin.Plugin, error) {
		return task, nil
	}
	return factory
}

func mockTestConfig() *config.CloudConfig {
	conf := &config.CloudConfig{}
	conf.Plugin.Resource = common.RandString(9)
	conf.Plugin.Objects = []string{common.RandString(9)}
	conf.Plugin.PKI = common.RandString(9)
	conf.Plugin.Auth = common.RandString(9)
	conf.Plugin.Functions = []string{common.RandString(9)}
	conf.Plugin.Shadow = common.RandString(9)
	conf.Plugin.Index = common.RandString(9)
	conf.Plugin.AppHistory = common.RandString(9)
	conf.Plugin.License = common.RandString(9)
	conf.Plugin.Property = common.RandString(9)
	conf.Plugin.Task = common.RandString(9)
	conf.Template.Path = "../scripts/native/templates"
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

	mProperty := mockPlugin.NewMockProperty(mockCtl)
	plugin.RegisterFactory(conf.Plugin.Property, mockProperty(mProperty))

	mModule := mockPlugin.NewMockModule(mockCtl)
	plugin.RegisterFactory(conf.Plugin.Module, mockModule(mModule))

	mResource := mockPlugin.NewMockResource(mockCtl)
	plugin.RegisterFactory(conf.Plugin.Resource, mockResource(mResource))

	mShadow := mockPlugin.NewMockShadow(mockCtl)
	plugin.RegisterFactory(conf.Plugin.Shadow, mockShadowStorage(mShadow))

	mIndex := mockPlugin.NewMockIndex(mockCtl)
	plugin.RegisterFactory(conf.Plugin.Index, mockIndex(mIndex))

	mAppHis := mockPlugin.NewMockAppHistory(mockCtl)
	plugin.RegisterFactory(conf.Plugin.AppHistory, mockAppHis(mAppHis))

	mTask := mockPlugin.NewMockTask(mockCtl)
	plugin.RegisterFactory(conf.Plugin.Task, mockTask(mTask))

	_, err := NewSyncService(conf)
	assert.Nil(t, err)

	return &MockServices{
		conf:           conf,
		ctl:            mockCtl,
		node:           mResource,
		shadow:         mShadow,
		namespace:      mResource,
		configuration:  mResource,
		secret:         mResource,
		app:            mResource,
		index:          mIndex,
		appHis:         mAppHis,
		objectStorage:  mockObjectStorage,
		functionPlugin: mockFunctionPlugin,
		pki:            mPKI,
		auth:           mAuth,
		license:        mLicense,
		property:       mProperty,
		module:         mModule,
		task:           mTask,
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
	mProperty := mockPlugin.NewMockProperty(mockCtl)
	plugin.RegisterFactory(conf.Plugin.Property, mockProperty(mProperty))
	return &MockServices{
		conf:           conf,
		ctl:            mockCtl,
		objectStorage:  mockObjectStorage,
		functionPlugin: mockFunctionPlugin,
		property:       mProperty,
	}
}
