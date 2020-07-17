package api

import (
	"bytes"
	"encoding/json"
	"github.com/baetyl/baetyl-cloud/common"
	plugin "github.com/baetyl/baetyl-cloud/mock/service"
	"github.com/baetyl/baetyl-cloud/models"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func initSystemConfigAPI(t *testing.T) (*API, *gin.Engine, *gomock.Controller) {
	api := &API{}
	router := gin.Default()
	mockCtl := gomock.NewController(t)
	mockIM := func(c *gin.Context) { common.NewContext(c).SetNamespace("default") }

	v1 := router.Group("v1")
	{
		systemconfig := v1.Group("/system/configs")

		systemconfig.GET("/:key", mockIM, common.Wrapper(api.GetSystemConfig))
		systemconfig.GET("", mockIM, common.Wrapper(api.ListSystemConfig))

		systemconfig.POST("", mockIM, common.Wrapper(api.CreateSystemConfig))
		systemconfig.DELETE("/:key", mockIM, common.Wrapper(api.DeleteSystemConfig))
		systemconfig.PUT("/:key", mockIM, common.Wrapper(api.UpdateSystemConfig))
	}
	return api, router, mockCtl
}

func TestAPI_CreateSystemConfig(t *testing.T) {
	api, router, ctl := initSystemConfigAPI(t)
	rs := plugin.NewMockCacheService(ctl)
	api.cacheService = rs

	systemConfig := genSystemConfig()

	rs.EXPECT().CreateSystemConfig(gomock.Any()).Return(systemConfig, nil)

	body, _ := json.Marshal(systemConfig)
	req, _ := http.NewRequest(http.MethodPost, "/v1/system/configs", bytes.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAPI_GetSystemConfig(t *testing.T) {
	api, router, ctl := initSystemConfigAPI(t)
	rs := plugin.NewMockCacheService(ctl)
	api.cacheService = rs

	systemConfig := genSystemConfig()

	rs.EXPECT().GetSystemConfig(systemConfig.Key).Return(systemConfig, nil)

	req, _ := http.NewRequest(http.MethodGet, "/v1/system/configs/"+systemConfig.Key, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAPI_DeleteSystemConfig(t *testing.T) {
	api, router, ctl := initSystemConfigAPI(t)
	rs := plugin.NewMockCacheService(ctl)
	api.cacheService = rs

	systemConfig := genSystemConfig()

	rs.EXPECT().DeleteSystemConfig(gomock.Any()).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/v1/system/configs/"+systemConfig.Key, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAPI_UpdateSystemConfig(t *testing.T) {
	api, router, ctl := initSystemConfigAPI(t)
	rs := plugin.NewMockCacheService(ctl)
	api.cacheService = rs

	systemConfig := genSystemConfig()

	rs.EXPECT().GetSystemConfig(gomock.Any()).Return(systemConfig, nil).AnyTimes()
	rs.EXPECT().UpdateSystemConfig(gomock.Any()).Return(systemConfig, nil)

	body, _ := json.Marshal(systemConfig)
	req, _ := http.NewRequest(http.MethodPut, "/v1/system/configs/"+systemConfig.Key, bytes.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func genSystemConfig() *models.SystemConfig {
	return &models.SystemConfig{
		Key:   "bae",
		Value: "http://test",
	}
}

func TestAPI_ListSystemConfig(t *testing.T) {
	api, router, ctl := initSystemConfigAPI(t)
	rs := plugin.NewMockCacheService(ctl)
	api.cacheService = rs

	mConf := genSystemConfig()

	page := &models.Filter{
		PageNo:   1,
		PageSize: 2,
		Name:     "%",
	}
	rs.EXPECT().ListSystemConfig(page).Return(&models.ListView{
		Total:    1,
		PageNo:   1,
		PageSize: 2,
		Items: []models.SystemConfig{*mConf},
	},nil)

	req, _ := http.NewRequest(http.MethodGet, "/v1/system/configs?pageNo=1&pageSize=2", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
