package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/baetyl/baetyl-cloud/v2/api"

	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	mockPlugin "github.com/baetyl/baetyl-cloud/v2/mock/plugin"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

func initAdminServerMock(t *testing.T) (*AdminServer, *mockPlugin.MockAuth, *mockPlugin.MockLicense, *gomock.Controller) {
	c := &config.CloudConfig{}
	c.Plugin.Auth = common.RandString(9)
	c.Plugin.ModelStorage = common.RandString(9)
	c.Plugin.DatabaseStorage = common.RandString(9)
	c.Plugin.Shadow = c.Plugin.DatabaseStorage
	c.Plugin.Objects = []string{common.RandString(9)}
	c.Plugin.PKI = common.RandString(9)
	c.Plugin.Functions = []string{common.RandString(9)}
	c.Plugin.License = common.RandString(9)
	c.Plugin.Property = common.RandString(9)
	mockCtl := gomock.NewController(t)

	mockModelStorage := mockPlugin.NewMockModelStorage(mockCtl)
	plugin.RegisterFactory(c.Plugin.ModelStorage, func() (plugin.Plugin, error) {
		return mockModelStorage, nil
	})
	res := &models.ConfigurationList{}
	mockModelStorage.EXPECT().ListConfig(gomock.Any(), gomock.Any()).Return(res, nil).AnyTimes()
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

	mockAPI, err := api.NewAPI(c)
	assert.NoError(t, err)

	s, err := NewAdminServer(c)
	assert.NoError(t, err)
	s.SetAPI(mockAPI)

	return s, mockAuth, mLicense, mockCtl
}

func TestAdminServer_Handler(t *testing.T) {
	s, mkAuth, mkLicense, mockCtl := initAdminServerMock(t)
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
	mkLicense.EXPECT().CheckQuota(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	req, _ = http.NewRequest(http.MethodPost, "/v1/nodes", nil)
	w4 = httptest.NewRecorder()
	s.GetRoute().ServeHTTP(w4, req)
	assert.Equal(t, http.StatusInternalServerError, w4.Code)

	go s.Run()
	defer s.Close()
}
