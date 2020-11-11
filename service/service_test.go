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
	dbStorage      *mockPlugin.MockDBStorage
	node           *mockPlugin.MockNode
	namespace      *mockPlugin.MockNamespace
	configuration  *mockPlugin.MockConfiguration
	secret         *mockPlugin.MockSecret
	app            *mockPlugin.MockApplication
	matcher        *mockPlugin.MockMatcher
	objectStorage  *mockPlugin.MockObject
	functionPlugin *mockPlugin.MockFunction
	pki            *mockPlugin.MockPKI
	auth           *mockPlugin.MockAuth
	shadow         *mockPlugin.MockShadow
	license        *mockPlugin.MockLicense
	property       *mockPlugin.MockProperty
}

func (m *MockServices) Close() {
	if m.ctl != nil {
		m.ctl.Finish()
	}
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

func mockProperty(mock plugin.Property) plugin.Factory {
	factory := func() (plugin.Plugin, error) {
		return mock, nil
	}
	return factory
}

func mockNode(mock plugin.Node) plugin.Factory {
	factory := func() (plugin.Plugin, error) {
		return mock, nil
	}
	return factory
}

func mockNamespace(mock plugin.Namespace) plugin.Factory {
	factory := func() (plugin.Plugin, error) {
		return mock, nil
	}
	return factory
}

func mockConfig(mock plugin.Configuration) plugin.Factory {
	factory := func() (plugin.Plugin, error) {
		return mock, nil
	}
	return factory
}

func mockSecret(mock plugin.Secret) plugin.Factory {
	factory := func() (plugin.Plugin, error) {
		return mock, nil
	}
	return factory
}

func mockApplication(mock plugin.Application) plugin.Factory {
	factory := func() (plugin.Plugin, error) {
		return mock, nil
	}
	return factory
}

func mockMatcher(mock plugin.Matcher) plugin.Factory {
	factory := func() (plugin.Plugin, error) {
		return mock, nil
	}
	return factory
}

func mockTestConfig() *config.CloudConfig {
	conf := &config.CloudConfig{}
	conf.Plugin.ModelStorage = common.RandString(9)
	conf.Plugin.DatabaseStorage = common.RandString(9)
	conf.Plugin.Objects = []string{common.RandString(9)}
	conf.Plugin.PKI = common.RandString(9)
	conf.Plugin.Auth = common.RandString(9)
	conf.Plugin.Functions = []string{common.RandString(9)}
	conf.Plugin.Shadow = common.RandString(9)
	conf.Plugin.Node = common.RandString(9)
	conf.Plugin.Namespace = common.RandString(9)
	conf.Plugin.Configuration = common.RandString(9)
	conf.Plugin.Application = common.RandString(9)
	conf.Plugin.Matcher = common.RandString(9)
	conf.Plugin.Secret = common.RandString(9)
	conf.Plugin.License = common.RandString(9)
	conf.Plugin.Property = common.RandString(9)
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

	mProperty := mockPlugin.NewMockProperty(mockCtl)
	plugin.RegisterFactory(conf.Plugin.Property, mockProperty(mProperty))

	mNode := mockPlugin.NewMockNode(mockCtl)
	plugin.RegisterFactory(conf.Plugin.Node, mockNode(mNode))

	mShadow := mockPlugin.NewMockShadow(mockCtl)
	plugin.RegisterFactory(conf.Plugin.Shadow, mockShadowStorage(mShadow))

	mNamespace := mockPlugin.NewMockNamespace(mockCtl)
	plugin.RegisterFactory(conf.Plugin.Namespace, mockNamespace(mNamespace))

	mConfig := mockPlugin.NewMockConfiguration(mockCtl)
	plugin.RegisterFactory(conf.Plugin.Configuration, mockConfig(mConfig))

	mSecret := mockPlugin.NewMockSecret(mockCtl)
	plugin.RegisterFactory(conf.Plugin.Secret, mockSecret(mSecret))

	mApp := mockPlugin.NewMockApplication(mockCtl)
	plugin.RegisterFactory(conf.Plugin.Application, mockApplication(mApp))

	mMatcher := mockPlugin.NewMockMatcher(mockCtl)
	plugin.RegisterFactory(conf.Plugin.Matcher, mockMatcher(mMatcher))

	_, err := NewSyncService(conf)
	assert.Nil(t, err)

	return &MockServices{
		conf:           conf,
		ctl:            mockCtl,
		dbStorage:      mockDBStorage,
		node:           mNode,
		shadow:         mShadow,
		namespace:      mNamespace,
		configuration:  mConfig,
		secret:         mSecret,
		app:            mApp,
		matcher:        mMatcher,
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
