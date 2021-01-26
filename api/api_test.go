package api

import (
	"testing"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	mockPlugin "github.com/baetyl/baetyl-cloud/v2/mock/plugin"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewAdminAPI(t *testing.T) {
	c := &config.CloudConfig{}
	c.Plugin.Pubsub = common.RandString(9)
	c.Plugin.PKI = common.RandString(9)
	c.Plugin.Auth = common.RandString(9)
	c.Plugin.License = common.RandString(9)
	c.Plugin.Resource = common.RandString(9)
	c.Plugin.Shadow = common.RandString(9)
	c.Plugin.Index = common.RandString(9)
	c.Plugin.Batch = common.RandString(9)
	c.Plugin.Record = common.RandString(9)
	c.Plugin.Callback = common.RandString(9)
	c.Plugin.AppHistory = common.RandString(9)
	c.Plugin.Objects = []string{common.RandString(9), common.RandString(9)}
	c.Plugin.Functions = []string{common.RandString(9), common.RandString(9)}
	c.Plugin.Property = common.RandString(9)
	c.Plugin.Module = common.RandString(9)
	c.Plugin.SyncLinks = []string{common.RandString(9), common.RandString(9)}

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	mockPubsub := mockPlugin.NewMockPubsub(mockCtl)
	plugin.RegisterFactory(c.Plugin.Pubsub, func() (plugin.Plugin, error) {
		return mockPubsub, nil
	})
	mockPKI := mockPlugin.NewMockPKI(mockCtl)
	plugin.RegisterFactory(c.Plugin.PKI, func() (plugin.Plugin, error) {
		return mockPKI, nil
	})
	mockAuth := mockPlugin.NewMockAuth(mockCtl)
	plugin.RegisterFactory(c.Plugin.Auth, func() (plugin.Plugin, error) {
		return mockAuth, nil
	})
	mockLicense := mockPlugin.NewMockLicense(mockCtl)
	plugin.RegisterFactory(c.Plugin.License, func() (plugin.Plugin, error) {
		return mockLicense, nil
	})
	mockResource := mockPlugin.NewMockResource(mockCtl)
	plugin.RegisterFactory(c.Plugin.Resource, func() (plugin.Plugin, error) {
		return mockResource, nil
	})
	mockShadow := mockPlugin.NewMockShadow(mockCtl)
	plugin.RegisterFactory(c.Plugin.Shadow, func() (plugin.Plugin, error) {
		return mockShadow, nil
	})

	mockIndex := mockPlugin.NewMockIndex(mockCtl)
	plugin.RegisterFactory(c.Plugin.Index, func() (plugin.Plugin, error) {
		return mockIndex, nil
	})
	mockAppHistory := mockPlugin.NewMockAppHistory(mockCtl)
	plugin.RegisterFactory(c.Plugin.AppHistory, func() (plugin.Plugin, error) {
		return mockAppHistory, nil
	})
	mockProperty := mockPlugin.NewMockProperty(mockCtl)
	plugin.RegisterFactory(c.Plugin.Property, func() (plugin.Plugin, error) {
		return mockProperty, nil
	})
	mockModule := mockPlugin.NewMockModule(mockCtl)
	plugin.RegisterFactory(c.Plugin.Module, func() (plugin.Plugin, error) {
		return mockModule, nil
	})

	mockObjectStorage := mockPlugin.NewMockObject(mockCtl)
	for _, v := range c.Plugin.Objects {
		plugin.RegisterFactory(v, func() (plugin.Plugin, error) {
			return mockObjectStorage, nil
		})
	}
	mockFunctions := mockPlugin.NewMockFunction(mockCtl)
	for _, v := range c.Plugin.Functions {
		plugin.RegisterFactory(v, func() (plugin.Plugin, error) {
			return mockFunctions, nil
		})
	}

	api, err := NewAPI(c)
	assert.NoError(t, err)
	assert.NotNil(t, api)
}
