package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/baetyl/baetyl-cloud/v2/models"

	"github.com/baetyl/baetyl-cloud/v2/mock/service"

	"github.com/baetyl/baetyl-cloud/v2/api"

	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	mockPlugin "github.com/baetyl/baetyl-cloud/v2/mock/plugin"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

func initAdminServerMock(t *testing.T) (*AdminServer, *mockPlugin.MockAuth, *mockPlugin.MockLicense, *gomock.Controller) {
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

	s, err := NewAdminServer(c)
	assert.NoError(t, err)
	s.SetAPI(mockAPI)

	return s, mockAuth, mLicense, mockCtl
}

func TestAdminServer_Handler(t *testing.T) {
	s, mkAuth, _, mockCtl := initAdminServerMock(t)
	defer mockCtl.Finish()

	s.InitRoute()
	r := s.GetRoute()
	assert.NotNil(t, r)

	mkAuth.EXPECT().Authenticate(gomock.Any()).Return(nil)
	// 200
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	s.GetRoute().ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	mkAuth.EXPECT().Authenticate(gomock.Any()).Return(nil)
	// 200
	req, _ = http.NewRequest(http.MethodGet, "/v1/configs", nil)
	w2 := httptest.NewRecorder()
	s.GetRoute().ServeHTTP(w2, req)
	assert.Equal(t, http.StatusOK, w2.Code)

	//mkAuth.EXPECT().Authenticate(gomock.Any()).Return(nil)
	// 400
	mConf := &specV1.Configuration{
		Namespace: "default",
		Name:      "abc",
		Labels:    make(map[string]string),
		Data:      make(map[string]string),
	}
	body, _ := json.Marshal(mConf)
	assert.Equal(t, http.StatusOK, w.Code)
	req, _ = http.NewRequest(http.MethodHead, "/zzz", bytes.NewReader(body))
	w3 := httptest.NewRecorder()
	s.GetRoute().ServeHTTP(w3, req)
	assert.Equal(t, http.StatusNotFound, w3.Code)

	// 401
	mkAuth.EXPECT().Authenticate(gomock.Any()).Return(fmt.Errorf("err"))
	req, _ = http.NewRequest(http.MethodGet, "/v1/configs", nil)
	w4 := httptest.NewRecorder()
	s.GetRoute().ServeHTTP(w4, req)
	assert.Equal(t, http.StatusUnauthorized, w4.Code)

	// 401
	mkAuth.EXPECT().Authenticate(gomock.Any()).Return(nil)
	mLicense := service.NewMockLicenseService(mockCtl)
	s.api.License = mLicense
	mLicense.EXPECT().CheckQuota(gomock.Any(), gomock.Any()).Return(fmt.Errorf("quota error"))
	req, _ = http.NewRequest(http.MethodPost, "/v1/nodes", nil)
	w4 = httptest.NewRecorder()
	s.GetRoute().ServeHTTP(w4, req)
	assert.Equal(t, http.StatusInternalServerError, w4.Code)

	mkAuth.EXPECT().Authenticate(gomock.Any()).Return(nil)

	mconf := &models.ConfigurationView{
		Name:      "baetyl-abc",
		Namespace: "default",
	}

	w = httptest.NewRecorder()
	body, _ = json.Marshal(mconf)
	req, _ = http.NewRequest(http.MethodPost, "/v1/configs", bytes.NewReader(body))
	s.GetRoute().ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mkAuth.EXPECT().Authenticate(gomock.Any()).Return(nil)
	req, _ = http.NewRequest(http.MethodDelete, "/v1/configs/baetyl-abc", nil)
	w = httptest.NewRecorder()
	s.GetRoute().ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	msecret := &models.SecretView{
		Name:      "baetyl-abc",
		Namespace: "default",
	}

	mkAuth.EXPECT().Authenticate(gomock.Any()).Return(nil)
	w = httptest.NewRecorder()
	body, _ = json.Marshal(msecret)
	req, _ = http.NewRequest(http.MethodPost, "/v1/secrets", bytes.NewReader(body))
	s.GetRoute().ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mkAuth.EXPECT().Authenticate(gomock.Any()).Return(nil)
	req, _ = http.NewRequest(http.MethodDelete, "/v1/secrets/baetyl-abc", nil)
	w = httptest.NewRecorder()
	s.GetRoute().ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mregistry := &models.Registry{
		Name:      "baetyl-abc",
		Namespace: "default",
	}

	mkAuth.EXPECT().Authenticate(gomock.Any()).Return(nil)
	w = httptest.NewRecorder()
	body, _ = json.Marshal(mregistry)
	req, _ = http.NewRequest(http.MethodPost, "/v1/registries", bytes.NewReader(body))
	s.GetRoute().ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mkAuth.EXPECT().Authenticate(gomock.Any()).Return(nil)
	req, _ = http.NewRequest(http.MethodDelete, "/v1/registries/baetyl-abc", nil)
	w = httptest.NewRecorder()
	s.GetRoute().ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mcertificate := &models.Registry{
		Name:      "baetyl-abc",
		Namespace: "default",
	}

	mkAuth.EXPECT().Authenticate(gomock.Any()).Return(nil)
	w = httptest.NewRecorder()
	body, _ = json.Marshal(mcertificate)
	req, _ = http.NewRequest(http.MethodPost, "/v1/certificates", bytes.NewReader(body))
	s.GetRoute().ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mkAuth.EXPECT().Authenticate(gomock.Any()).Return(nil)
	req, _ = http.NewRequest(http.MethodDelete, "/v1/certificates/baetyl-abc", nil)
	w = httptest.NewRecorder()
	s.GetRoute().ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mapp := &models.ApplicationView{
		Name:      "baetyl-abc",
		Namespace: "default",
	}

	mkAuth.EXPECT().Authenticate(gomock.Any()).Return(nil)
	w = httptest.NewRecorder()
	body, _ = json.Marshal(mapp)
	req, _ = http.NewRequest(http.MethodPost, "/v1/apps", bytes.NewReader(body))
	s.GetRoute().ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mkAuth.EXPECT().Authenticate(gomock.Any()).Return(nil)
	req, _ = http.NewRequest(http.MethodDelete, "/v1/apps/baetyl-abc", nil)
	w = httptest.NewRecorder()
	s.GetRoute().ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	go s.Run()
	defer s.Close()
}
