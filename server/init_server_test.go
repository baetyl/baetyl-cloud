package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/baetyl/baetyl-cloud/v2/api"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	mockPlugin "github.com/baetyl/baetyl-cloud/v2/mock/plugin"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

func initInitServerMock(t *testing.T) (*InitServer, *gomock.Controller, *config.CloudConfig) {
	c := &config.CloudConfig{}
	c.Plugin.Auth = common.RandString(9)
	c.Plugin.ModelStorage = common.RandString(9)
	c.Plugin.DatabaseStorage = common.RandString(9)
	c.Plugin.Shadow = common.RandString(9)
	c.Plugin.Node = common.RandString(9)
	c.Plugin.Namespace = common.RandString(9)
	c.Plugin.Configuration = common.RandString(9)
	c.Plugin.Secret = common.RandString(9)
	c.Plugin.Application = common.RandString(9)
	c.Plugin.Index = common.RandString(9)
	c.Plugin.Objects = []string{common.RandString(9)}
	c.Plugin.PKI = common.RandString(9)
	c.Plugin.Functions = []string{common.RandString(9)}
	c.Plugin.License = common.RandString(9)
	c.Plugin.Property = common.RandString(9)
	c.InitServer.Certificate.CA = "../scripts/demo/native/certs/client_ca.crt"
	c.InitServer.Certificate.Cert = "../scripts/demo/native/certs/server.crt"
	c.InitServer.Certificate.Key = "../scripts/demo/native/certs/server.key"
	mockCtl := gomock.NewController(t)

	mockObjectStorage := mockPlugin.NewMockObject(mockCtl)
	for _, v := range c.Plugin.Objects {
		plugin.RegisterFactory(v, func() (plugin.Plugin, error) {
			return mockObjectStorage, nil
		})
	}
	mockFunction := mockPlugin.NewMockFunction(mockCtl)
	for _, v := range c.Plugin.Functions {
		plugin.RegisterFactory(v, func() (plugin.Plugin, error) {
			return mockFunction, nil
		})
	}

	mockAuth := mockPlugin.NewMockAuth(mockCtl)
	plugin.RegisterFactory(c.Plugin.Auth, func() (plugin.Plugin, error) {
		return mockAuth, nil
	})

	mPKI := mockPlugin.NewMockPKI(mockCtl)
	plugin.RegisterFactory(c.Plugin.PKI, func() (plugin.Plugin, error) {
		return mPKI, nil
	})

	mLicense := mockPlugin.NewMockLicense(mockCtl)
	plugin.RegisterFactory(c.Plugin.License, func() (plugin.Plugin, error) {
		return mLicense, nil
	})
	mockProperty := mockPlugin.NewMockProperty(mockCtl)
	plugin.RegisterFactory(c.Plugin.Property, func() (plugin.Plugin, error) {
		return mockProperty, nil
	})
	mockNamespace := mockPlugin.NewMockNamespace(mockCtl)
	plugin.RegisterFactory(c.Plugin.Namespace, func() (plugin.Plugin, error) {
		return mockNamespace, nil
	})
	mockConfig := mockPlugin.NewMockConfiguration(mockCtl)
	plugin.RegisterFactory(c.Plugin.Configuration, func() (plugin.Plugin, error) {
		return mockConfig, nil
	})
	res := &models.ConfigurationList{}
	mockConfig.EXPECT().ListConfig(gomock.Any(), gomock.Any()).Return(res, nil).AnyTimes()

	mockSecret := mockPlugin.NewMockSecret(mockCtl)
	plugin.RegisterFactory(c.Plugin.Secret, func() (plugin.Plugin, error) {
		return mockSecret, nil
	})

	mockApplication := mockPlugin.NewMockApplication(mockCtl)
	plugin.RegisterFactory(c.Plugin.Application, func() (plugin.Plugin, error) {
		return mockApplication, nil
	})

	mockIndex := mockPlugin.NewMockIndex(mockCtl)
	plugin.RegisterFactory(c.Plugin.Index, func() (plugin.Plugin, error) {
		return mockIndex, nil
	})

	mockInitAPI, err := api.NewInitAPI(c)
	assert.NoError(t, err)

	s, err := NewInitServer(c)
	assert.NoError(t, err)
	s.SetAPI(mockInitAPI)

	return s, mockCtl, c
}

func TestInitServer_Handler(t *testing.T) {
	t.Skip()
	s, mockCtl, c := initInitServerMock(t)
	defer mockCtl.Finish()

	// https 200
	s.InitRoute()
	r := s.GetRoute()
	assert.NotNil(t, r)

	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	s.GetRoute().ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	go s.Run()
	defer s.Close()
	// http 200
	c.InitServer.Certificate.Key = ""
	c.InitServer.Certificate.Cert = ""
	aHttp, err := NewInitServer(c)
	assert.NoError(t, err)
	aHttp.InitRoute()
	req, err = http.NewRequest(http.MethodGet, "/health", nil)
	assert.NoError(t, err)
	w = httptest.NewRecorder()
	aHttp.GetRoute().ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	go aHttp.Run()
	defer aHttp.Close()
}
