package task

import (
	"testing"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	mockPlugin "github.com/baetyl/baetyl-cloud/v2/mock/plugin"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
	"github.com/golang/mock/gomock"
)

type MockServices struct {
	conf          *config.CloudConfig
	ctl           *gomock.Controller
	node          *mockPlugin.MockResource
	namespace     *mockPlugin.MockResource
	configuration *mockPlugin.MockResource
	secret        *mockPlugin.MockResource
	app           *mockPlugin.MockResource
	index         *mockPlugin.MockIndex
	license       *mockPlugin.MockLicense
	task          *mockPlugin.MockTask
	Lock          *mockPlugin.MockLocker
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

func mockTask(task plugin.Task) plugin.Factory {
	factory := func() (plugin.Plugin, error) {
		return task, nil
	}
	return factory
}

func mockTestConfig() *config.CloudConfig {
	conf := &config.CloudConfig{}
	conf.Plugin.Resource = common.RandString(9)
	conf.Plugin.Index = common.RandString(9)
	conf.Plugin.License = common.RandString(9)
	conf.Plugin.Task = common.RandString(9)
	conf.Plugin.Locker = common.RandString(9)

	return conf
}

func InitMockEnvironment(t *testing.T) *MockServices {
	conf := mockTestConfig()
	mockCtl := gomock.NewController(t)

	mLicense := mockPlugin.NewMockLicense(mockCtl)
	plugin.RegisterFactory(conf.Plugin.License, mockLicense(mLicense))

	mResource := mockPlugin.NewMockResource(mockCtl)
	plugin.RegisterFactory(conf.Plugin.Resource, mockResource(mResource))

	mIndex := mockPlugin.NewMockIndex(mockCtl)
	plugin.RegisterFactory(conf.Plugin.Index, mockIndex(mIndex))

	mTask := mockPlugin.NewMockTask(mockCtl)
	plugin.RegisterFactory(conf.Plugin.Task, mockTask(mTask))

	mLock := mockPlugin.NewMockLocker(mockCtl)

	return &MockServices{
		conf:          conf,
		ctl:           mockCtl,
		node:          mResource,
		namespace:     mResource,
		configuration: mResource,
		secret:        mResource,
		app:           mResource,
		index:         mIndex,
		license:       mLicense,
		task:          mTask,
		Lock:          mLock,
	}
}
