package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/api"
	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	mockPlugin "github.com/baetyl/baetyl-cloud/v2/mock/plugin"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

func initMisServerMock(t *testing.T) (*MisServer, *gomock.Controller) {
	c := &config.CloudConfig{}
	c.Plugin.Auth = common.RandString(9)
	c.Plugin.Resource = common.RandString(9)
	c.Plugin.Shadow = common.RandString(9)
	c.Plugin.Index = common.RandString(9)
	c.Plugin.AppHistory = common.RandString(9)
	c.Plugin.Objects = []string{common.RandString(9)}
	c.Plugin.PKI = common.RandString(9)
	c.Plugin.Functions = []string{common.RandString(9)}
	c.Plugin.License = common.RandString(9)
	c.Plugin.Property = common.RandString(9)
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
	mockResource := mockPlugin.NewMockResource(mockCtl)
	plugin.RegisterFactory(c.Plugin.Resource, func() (plugin.Plugin, error) {
		return mockResource, nil
	})
	mockShadow := mockPlugin.NewMockShadow(mockCtl)
	plugin.RegisterFactory(c.Plugin.Shadow, func() (plugin.Plugin, error) {
		return mockShadow, nil
	})
	plugin.RegisterFactory(c.Plugin.Resource, func() (plugin.Plugin, error) {
		return mockResource, nil
	})
	plugin.RegisterFactory(c.Plugin.Resource, func() (plugin.Plugin, error) {
		return mockResource, nil
	})
	res := &models.ConfigurationList{}
	mockResource.EXPECT().ListConfig(gomock.Any(), gomock.Any()).Return(res, nil).AnyTimes()

	plugin.RegisterFactory(c.Plugin.Resource, func() (plugin.Plugin, error) {
		return mockResource, nil
	})

	plugin.RegisterFactory(c.Plugin.Resource, func() (plugin.Plugin, error) {
		return mockResource, nil
	})

	mockIndex := mockPlugin.NewMockIndex(mockCtl)
	plugin.RegisterFactory(c.Plugin.Index, func() (plugin.Plugin, error) {
		return mockIndex, nil
	})
	mockAppHis := mockPlugin.NewMockAppHistory(mockCtl)
	plugin.RegisterFactory(c.Plugin.AppHistory, func() (plugin.Plugin, error) {
		return mockAppHis, nil
	})

	mockAPI, err := api.NewAPI(c)
	assert.NoError(t, err)

	s, err := NewMisServer(c)
	assert.NoError(t, err)
	s.SetAPI(mockAPI)

	return s, mockCtl
}

func TestMisServer_Handler(t *testing.T) {
	s, mockCtl := initMisServerMock(t)
	defer mockCtl.Finish()

	s.InitRoute()
	s.router.Use(s.authHandler)
	r := s.GetRoute()
	assert.NotNil(t, r)
	// https 200
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set("baetyl-cloud-token", "baetyl-cloud-token")
	req.Header.Set("baetyl-cloud-user", "1")
	w := httptest.NewRecorder()
	s.GetRoute().ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	go s.Run()
	defer s.Close()
	// http 200
	req, _ = http.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set("baetyl-cloud-token", "baetyl-cloud-token")
	req.Header.Set("baetyl-cloud-user", "")
	w = httptest.NewRecorder()
	s.GetRoute().ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	go s.Run()
	defer s.Close()
	// http 200
	req, _ = http.NewRequest(http.MethodGet, "/health", nil)
	w = httptest.NewRecorder()
	s.GetRoute().ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	go s.Run()
	defer s.Close()
}
