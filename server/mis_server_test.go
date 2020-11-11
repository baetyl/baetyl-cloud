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
	c.Plugin.ModelStorage = common.RandString(9)
	c.Plugin.DatabaseStorage = common.RandString(9)
	c.Plugin.Shadow = common.RandString(9)
	c.Plugin.Node = common.RandString(9)
	c.Plugin.Namespace = common.RandString(9)
	c.Plugin.Configuration = common.RandString(9)
	c.Plugin.Secret = common.RandString(9)
	c.Plugin.Application = common.RandString(9)
	c.Plugin.Matcher = common.RandString(9)
	c.Plugin.Objects = []string{common.RandString(9)}
	c.Plugin.PKI = common.RandString(9)
	c.Plugin.Functions = []string{common.RandString(9)}
	c.Plugin.License = common.RandString(9)
	c.Plugin.Property = common.RandString(9)
	mockCtl := gomock.NewController(t)

	mockDBStorage := mockPlugin.NewMockDBStorage(mockCtl)
	plugin.RegisterFactory(c.Plugin.DatabaseStorage, func() (plugin.Plugin, error) {
		return mockDBStorage, nil
	})
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
	mockNode := mockPlugin.NewMockNode(mockCtl)
	plugin.RegisterFactory(c.Plugin.Node, func() (plugin.Plugin, error) {
		return mockNode, nil
	})
	mockShadow := mockPlugin.NewMockShadow(mockCtl)
	plugin.RegisterFactory(c.Plugin.Shadow, func() (plugin.Plugin, error) {
		return mockShadow, nil
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

	mockMatcher := mockPlugin.NewMockMatcher(mockCtl)
	plugin.RegisterFactory(c.Plugin.Matcher, func() (plugin.Plugin, error) {
		return mockMatcher, nil
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
