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
	c.Plugin.Resource = common.RandString(9)
	c.Plugin.Shadow = common.RandString(9)
	c.Plugin.Index = common.RandString(9)
	c.Plugin.AppHistory = common.RandString(9)
	c.Plugin.Objects = []string{common.RandString(9)}
	c.Plugin.PKI = common.RandString(9)
	c.Plugin.Functions = []string{common.RandString(9)}
	c.Plugin.License = common.RandString(9)
	c.Plugin.Quota = common.RandString(9)
	c.Plugin.Property = common.RandString(9)
	c.Plugin.Module = common.RandString(9)
	c.Plugin.Task = common.RandString(9)
	c.Plugin.Locker = common.RandString(9)
	c.Plugin.Tx = common.RandString(9)
	c.Plugin.Sign = common.RandString(9)
	c.Plugin.Cron = common.RandString(9)
	c.Plugin.Cache = common.RandString(9)
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

	mockSign := mockPlugin.NewMockSign(mockCtl)
	plugin.RegisterFactory(c.Plugin.Sign, func() (plugin.Plugin, error) {
		return mockSign, nil
	})

	mPKI := mockPlugin.NewMockPKI(mockCtl)
	plugin.RegisterFactory(c.Plugin.PKI, func() (plugin.Plugin, error) {
		return mPKI, nil
	})

	mLicense := mockPlugin.NewMockLicense(mockCtl)
	plugin.RegisterFactory(c.Plugin.License, func() (plugin.Plugin, error) {
		return mLicense, nil
	})
	mQouta := mockPlugin.NewMockQuota(mockCtl)
	plugin.RegisterFactory(c.Plugin.Quota, func() (plugin.Plugin, error) {
		return mQouta, nil
	})
	mockProperty := mockPlugin.NewMockProperty(mockCtl)
	plugin.RegisterFactory(c.Plugin.Property, func() (plugin.Plugin, error) {
		return mockProperty, nil
	})
	mockModule := mockPlugin.NewMockModule(mockCtl)
	plugin.RegisterFactory(c.Plugin.Module, func() (plugin.Plugin, error) {
		return mockModule, nil
	})
	mockResource := mockPlugin.NewMockResource(mockCtl)
	plugin.RegisterFactory(c.Plugin.Resource, func() (plugin.Plugin, error) {
		return mockResource, nil
	})
	res := &models.ConfigurationList{}
	mockResource.EXPECT().ListConfig(gomock.Any(), gomock.Any()).Return(res, nil).AnyTimes()

	mockShadow := mockPlugin.NewMockShadow(mockCtl)
	plugin.RegisterFactory(c.Plugin.Shadow, func() (plugin.Plugin, error) {
		return mockShadow, nil
	})

	mockIndex := mockPlugin.NewMockIndex(mockCtl)
	plugin.RegisterFactory(c.Plugin.Index, func() (plugin.Plugin, error) {
		return mockIndex, nil
	})

	mockTask := mockPlugin.NewMockTask(mockCtl)
	plugin.RegisterFactory(c.Plugin.Task, func() (plugin.Plugin, error) {
		return mockTask, nil
	})

	mockLocker := mockPlugin.NewMockLocker(mockCtl)
	plugin.RegisterFactory(c.Plugin.Locker, func() (plugin.Plugin, error) {
		return mockLocker, nil
	})

	mockTx := mockPlugin.NewMockTransactionFactory(mockCtl)
	plugin.RegisterFactory(c.Plugin.Tx, func() (plugin.Plugin, error) {
		return mockTx, nil
	})

	mockCronApp := mockPlugin.NewMockCron(mockCtl)
	plugin.RegisterFactory(c.Plugin.Cron, func() (plugin.Plugin, error) {
		return mockCronApp, nil
	})
	mockCache := mockPlugin.NewMockDataCache(mockCtl)
	plugin.RegisterFactory(c.Plugin.Cache, func() (plugin.Plugin, error) {
		return mockCache, nil
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

	mLock := service.NewMockLockerService(mockCtl)
	s.api.Locker = mLock

	s.InitRoute()
	r := s.GetRoute()
	assert.NotNil(t, r)

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
	mQuota := service.NewMockQuotaService(mockCtl)
	s.api.Quota = mQuota
	mQuota.EXPECT().CheckQuota(gomock.Any(), gomock.Any()).Return(fmt.Errorf("quota error"))
	mLock.EXPECT().Lock(gomock.Any(), gomock.Any(), gomock.Any()).Return("", nil)
	mLock.EXPECT().Unlock(gomock.Any(), gomock.Any(), gomock.Any()).Return()
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
