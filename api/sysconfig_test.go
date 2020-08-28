package api

import (
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

// TODO: optimize this layer, general abstraction

func initSysConfigAPI(t *testing.T) (*API, *gin.Engine, *gomock.Controller) {
	api := &API{}
	router := gin.Default()
	mockCtl := gomock.NewController(t)
	mockIM := func(c *gin.Context) { common.NewContext(c).SetNamespace("default") }
	v1 := router.Group("v1")
	{
		configs := v1.Group("/sysconfig")
		configs.GET("/:type/:key", mockIM, common.Wrapper(api.GetSysConfig))
		configs.GET("/:type", mockIM, common.Wrapper(api.ListSysConfig))
	}

	return api, router, mockCtl
}

func TestGetSysConfig(t *testing.T) {
	api, router, mockCtl := initSysConfigAPI(t)
	defer mockCtl.Finish()
	mkSysConfigService := ms.NewMockSysConfigService(mockCtl)
	api.sysConfigService = mkSysConfigService

	sysConfig := &models.SysConfig{
		Type:  "baetyl_version",
		Key:   "1.0.1",
		Value: "http://test",
	}

	mkSysConfigService.EXPECT().GetSysConfig(sysConfig.Type, sysConfig.Key).Return(sysConfig, nil)

	// 200
	req, _ := http.NewRequest(http.MethodGet, "/v1/sysconfig/baetyl_version/1.0.1", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestListSysConfig(t *testing.T) {
	api, router, mockCtl := initSysConfigAPI(t)
	defer mockCtl.Finish()
	mkSysConfigService := ms.NewMockSysConfigService(mockCtl)
	api.sysConfigService = mkSysConfigService

	sysConfig := &models.SysConfig{
		Type:  "baetyl_version",
		Key:   "1.0.1",
		Value: "http://test",
	}

	mkSysConfigService.EXPECT().ListSysConfigAll(sysConfig.Type).Return([]models.SysConfig{*sysConfig}, nil)

	// 200
	req, _ := http.NewRequest(http.MethodGet, "/v1/sysconfig/baetyl_version", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
