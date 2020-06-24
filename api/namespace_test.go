package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/baetyl/baetyl-cloud/common"
	"github.com/baetyl/baetyl-cloud/models"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	ms "github.com/baetyl/baetyl-cloud/mock/service"
)

func getMockNS(name string) *models.Namespace {
	return &models.Namespace{
		Name: name,
	}
}

func initNamespaceAPI(t *testing.T) (*API, *gin.Engine, *gomock.Controller) {
	api := &API{}
	router := gin.Default()
	mockCtl := gomock.NewController(t)
	mockIMtestA := func(c *gin.Context) { c.Set(common.KeyContextNamespace, "testA") }
	mockIMtestB := func(c *gin.Context) { c.Set(common.KeyContextNamespace, "testB") }
	v1 := router.Group("testA")
	{
		testA := v1.Group("/namespace")
		testA.POST("", mockIMtestA, common.Wrapper(api.CreateNamespace))
		testA.GET("", mockIMtestA, common.Wrapper(api.GetNamespace))
		testA.DELETE("", mockIMtestA, common.Wrapper(api.DeleteNamespace))
	}
	v2 := router.Group("testB")
	{
		testB := v2.Group("/namespace")
		testB.POST("", mockIMtestB, common.Wrapper(api.CreateNamespace))
		testB.GET("", mockIMtestB, common.Wrapper(api.GetNamespace))

	}
	return api, router, mockCtl
}

func TestGetNamespace(t *testing.T) {
	api, router, mockCtl := initNamespaceAPI(t)
	defer mockCtl.Finish()
	mkNamespaceService := ms.NewMockNamespaceService(mockCtl)
	api.namespaceService = mkNamespaceService

	ns := getMockNS("testA")

	mkNamespaceService.EXPECT().Get("testA").Return(ns, nil)
	mkNamespaceService.EXPECT().Get("testB").Return(nil, nil)

	// 200
	req, _ := http.NewRequest(http.MethodGet, "/testA/namespace", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 404
	req, _ = http.NewRequest(http.MethodGet, "/testB/namespace", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req)
	assert.Equal(t, http.StatusNotFound, w2.Code)
}

func TestCreateNamespace(t *testing.T) {
	api, router, mockCtl := initNamespaceAPI(t)
	defer mockCtl.Finish()
	mkNamespaceService := ms.NewMockNamespaceService(mockCtl)
	api.namespaceService = mkNamespaceService

	nsa := getMockNS("testA")

	mkNamespaceService.EXPECT().Create(nsa).Return(nsa, nil)

	// 200
	req, _ := http.NewRequest(http.MethodPost, "/testA/namespace", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

}

func TestAPI_DeleteNamespace(t *testing.T) {
	api, router, mockCtl := initNamespaceAPI(t)
	defer mockCtl.Finish()
	mkNamespaceService := ms.NewMockNamespaceService(mockCtl)
	api.namespaceService = mkNamespaceService

	nsa := getMockNS("testA")

	mkNamespaceService.EXPECT().Delete(nsa).Return(nil)

	// 200
	req, _ := http.NewRequest(http.MethodDelete, "/testA/namespace", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

}
