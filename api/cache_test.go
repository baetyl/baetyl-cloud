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

func initCacheAPI(t *testing.T) (*API, *gin.Engine, *gomock.Controller) {
	api := &API{}
	router := gin.Default()
	mockCtl := gomock.NewController(t)
	mockIM := func(c *gin.Context) { common.NewContext(c).SetNamespace("default") }

	v1 := router.Group("v1")
	{
		cache := v1.Group("/caches")

		cache.GET("/:key", mockIM, common.Wrapper(api.GetCache))
		cache.GET("", mockIM, common.Wrapper(api.ListCache))

		cache.POST("", mockIM, common.Wrapper(api.CreateCache))
		cache.DELETE("/:key", mockIM, common.Wrapper(api.DeleteCache))
		cache.PUT("/:key", mockIM, common.Wrapper(api.UpdateCache))
	}
	return api, router, mockCtl
}

func genCache() *models.Cache {
	return &models.Cache{
		Key:   "bae",
		Value: "http://test",
	}
}

func TestCreateCache(t *testing.T) {
	api, router, ctl := initCacheAPI(t)
	rs := plugin.NewMockCacheService(ctl)
	api.cacheService = rs

	cache := genCache()

	rs.EXPECT().Set(cache.Key, cache.Value).Return(nil).Times(1)

	body, _ := json.Marshal(cache)
	req, _ := http.NewRequest(http.MethodPost, "/v1/caches", bytes.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetCache(t *testing.T) {
	api, router, ctl := initCacheAPI(t)
	rs := plugin.NewMockCacheService(ctl)
	api.cacheService = rs

	cache := genCache()

	rs.EXPECT().Get(cache.Key).Return(cache.Value, nil)

	req, _ := http.NewRequest(http.MethodGet, "/v1/caches/"+cache.Key, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteCache(t *testing.T) {
	api, router, ctl := initCacheAPI(t)
	rs := plugin.NewMockCacheService(ctl)
	api.cacheService = rs

	cache := genCache()

	rs.EXPECT().Delete(gomock.Any()).Return(nil).Times(1)

	req, _ := http.NewRequest(http.MethodDelete, "/v1/caches/"+cache.Key, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateCache(t *testing.T) {
	api, router, ctl := initCacheAPI(t)
	rs := plugin.NewMockCacheService(ctl)
	api.cacheService = rs

	cache := genCache()

	rs.EXPECT().Get(cache.Key).Return(cache.Value, nil).Times(1)
	rs.EXPECT().Set(cache.Key, cache.Value).Return(nil).Times(1)

	body, _ := json.Marshal(cache)
	req, _ := http.NewRequest(http.MethodPut, "/v1/caches/"+cache.Key, bytes.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestListCache(t *testing.T) {
	api, router, ctl := initCacheAPI(t)
	rs := plugin.NewMockCacheService(ctl)
	api.cacheService = rs

	mConf := genCache()

	page := &models.Filter{
		PageNo:   1,
		PageSize: 2,
		Name:     "%",
	}
	rs.EXPECT().List(page).Return(&models.ListView{
		Total:    1,
		PageNo:   1,
		PageSize: 2,
		Items:    []models.Cache{*mConf},
	}, nil)

	req, _ := http.NewRequest(http.MethodGet, "/v1/caches?pageNo=1&pageSize=2", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
