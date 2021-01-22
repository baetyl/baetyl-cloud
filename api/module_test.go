package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/common"
	ms "github.com/baetyl/baetyl-cloud/v2/mock/service"
	"github.com/baetyl/baetyl-cloud/v2/models"
)

func initModuleAPI(t *testing.T) (*API, *gin.Engine, *gomock.Controller) {
	api := &API{}
	router := gin.Default()
	mockCtl := gomock.NewController(t)
	mockIM := func(c *gin.Context) { common.NewContext(c).SetNamespace("default") }
	v1 := router.Group("v1")
	{
		module := v1.Group("modules")
		module.GET("", mockIM, common.Wrapper(api.ListModules))
		module.GET("/:name", mockIM, common.Wrapper(api.GetModules))
		module.GET("/:name/version/:version", mockIM, common.Wrapper(api.GetModuleByVersion))
		module.GET("/:name/latest", mockIM, common.Wrapper(api.GetLatestModule))
		module.POST("", mockIM, common.Wrapper(api.CreateModule))
		module.PUT("/:name/version/:version", mockIM, common.Wrapper(api.UpdateModule))
		module.DELETE("/:name", mockIM, common.Wrapper(api.DeleteModules))
		module.DELETE("/:name/version/:version", mockIM, common.Wrapper(api.DeleteModules))
	}
	return api, router, mockCtl
}

func TestGetModules(t *testing.T) {
	api, router, mockCtl := initModuleAPI(t)
	defer mockCtl.Finish()
	sModule := ms.NewMockModuleService(mockCtl)
	api.Module = sModule

	m1 := []models.Module{
		{
			Name:    "m1",
			Version: "m1v",
		},
	}
	sModule.EXPECT().GetModules("abc").Return(m1, nil).Times(1)
	req, _ := http.NewRequest(http.MethodGet, "/v1/modules/abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	sModule.EXPECT().GetModules("abc").Return(nil, errors.New("err")).Times(1)
	req, _ = http.NewRequest(http.MethodGet, "/v1/modules/abc", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetModuleByVersion(t *testing.T) {
	api, router, mockCtl := initModuleAPI(t)
	defer mockCtl.Finish()
	sModule := ms.NewMockModuleService(mockCtl)
	api.Module = sModule

	m1 := &models.Module{
		Name:    "m1",
		Version: "m1v",
	}
	sModule.EXPECT().GetModuleByVersion("abc", "m1v").Return(m1, nil).Times(1)
	req, _ := http.NewRequest(http.MethodGet, "/v1/modules/abc/version/m1v", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	sModule.EXPECT().GetModuleByVersion("abc", "m1v").Return(nil, errors.New("err")).Times(1)
	req, _ = http.NewRequest(http.MethodGet, "/v1/modules/abc/version/m1v", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetLatestModule(t *testing.T) {
	api, router, mockCtl := initModuleAPI(t)
	defer mockCtl.Finish()
	sModule := ms.NewMockModuleService(mockCtl)
	api.Module = sModule

	m1 := &models.Module{
		Name:    "m1",
		Version: "m1v",
	}
	sModule.EXPECT().GetLatestModule("abc").Return(m1, nil).Times(1)
	req, _ := http.NewRequest(http.MethodGet, "/v1/modules/abc/latest", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	sModule.EXPECT().GetLatestModule("abc").Return(nil, errors.New("err")).Times(1)
	req, _ = http.NewRequest(http.MethodGet, "/v1/modules/abc/latest", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCreateModule(t *testing.T) {
	api, router, mockCtl := initModuleAPI(t)
	defer mockCtl.Finish()
	sModule := ms.NewMockModuleService(mockCtl)
	sInit := ms.NewMockInitService(mockCtl)
	api.Module = sModule
	api.Init = sInit

	s := []string{
		"baetyl-function",
	}
	sInit.EXPECT().GetOptionalApps().Return(s).Times(1)
	m1 := &models.Module{
		Name:    "baetyl-function",
		Version: "m1v",
		Type:    string(common.TypeSystemOptional),
	}
	sModule.EXPECT().CreateModule(m1).Return(m1, nil).Times(1)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(m1)
	req, _ := http.NewRequest(http.MethodPost, "/v1/modules", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	sInit.EXPECT().GetOptionalApps().Return(s).Times(1)
	sModule.EXPECT().CreateModule(m1).Return(nil, errors.New("err")).Times(1)

	w = httptest.NewRecorder()
	body, _ = json.Marshal(m1)
	req, _ = http.NewRequest(http.MethodPost, "/v1/modules", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpdateModule(t *testing.T) {
	api, router, mockCtl := initModuleAPI(t)
	defer mockCtl.Finish()
	sModule := ms.NewMockModuleService(mockCtl)
	sInit := ms.NewMockInitService(mockCtl)
	api.Module = sModule
	api.Init = sInit

	s := []string{
		"baetyl-function",
	}
	sInit.EXPECT().GetOptionalApps().Return(s).Times(1)
	m1 := &models.Module{
		Name:    "baetyl-function",
		Version: "m1v",
		Type:    string(common.TypeSystemOptional),
	}
	res := &models.Module{
		Name:    "baetyl-function",
		Version: "m2v",
		Type:    string(common.TypeSystemOptional),
	}
	sModule.EXPECT().GetModuleByVersion(m1.Name, m1.Version).Return(res, nil).Times(1)
	sModule.EXPECT().UpdateModuleByVersion(m1).Return(m1, nil).Times(1)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(m1)
	req, _ := http.NewRequest(http.MethodPut, "/v1/modules/baetyl-function/version/m1v", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	sModule.EXPECT().GetModuleByVersion("baetyl", "m1v").Return(nil, errors.New("err")).Times(1)
	w = httptest.NewRecorder()
	body, _ = json.Marshal(m1)
	req, _ = http.NewRequest(http.MethodPut, "/v1/modules/baetyl/version/m1v", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeleteModules(t *testing.T) {
	api, router, mockCtl := initModuleAPI(t)
	defer mockCtl.Finish()
	sModule := ms.NewMockModuleService(mockCtl)
	sInit := ms.NewMockInitService(mockCtl)
	api.Module = sModule
	api.Init = sInit

	sModule.EXPECT().DeleteModules("baetyl").Return(nil).Times(1)

	req, _ := http.NewRequest(http.MethodDelete, "/v1/modules/baetyl", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	sModule.EXPECT().DeleteModuleByVersion("baetyl", "v1").Return(nil).Times(1)
	req, _ = http.NewRequest(http.MethodDelete, "/v1/modules/baetyl/version/v1", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestListModules(t *testing.T) {
	api, router, mockCtl := initModuleAPI(t)
	defer mockCtl.Finish()
	sModule := ms.NewMockModuleService(mockCtl)
	sInit := ms.NewMockInitService(mockCtl)
	api.Module = sModule
	api.Init = sInit

	res := []models.Module{
		{
			Name: "baetyl",
		},
	}
	sModule.EXPECT().ListModules(gomock.Any()).Return(res, nil).Times(1)

	req, _ := http.NewRequest(http.MethodGet, "/v1/modules", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	sModule.EXPECT().ListRuntimeModules(gomock.Any()).Return(res, nil).Times(1)

	req, _ = http.NewRequest(http.MethodGet, "/v1/modules?type=runtime_user", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	sModule.EXPECT().ListOptionalSysModules(gomock.Any()).Return(res, nil).Times(1)

	req, _ = http.NewRequest(http.MethodGet, "/v1/modules?type=opt_system", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	sModule.EXPECT().ListOptionalSysModules(gomock.Any()).Return(nil, errors.New("err")).Times(1)

	req, _ = http.NewRequest(http.MethodGet, "/v1/modules?type=opt_system", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
