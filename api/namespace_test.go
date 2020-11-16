package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"

	ms "github.com/baetyl/baetyl-cloud/v2/mock/service"
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
	api.NS = mkNamespaceService

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
	api.NS = mkNamespaceService

	mLicense := ms.NewMockLicenseService(mockCtl)
	api.License = mLicense

	nsa := getMockNS("testA")

	quotas := map[string]int{"maxNodeCount": 10}

	mLicense.EXPECT().GetDefaultQuotas(nsa.Name).Return(quotas, nil)
	mLicense.EXPECT().CreateQuota(nsa.Name, quotas).Return(nil)
	mkNamespaceService.EXPECT().Create(nsa).Return(nsa, nil)

	// 200
	req, _ := http.NewRequest(http.MethodPost, "/testA/namespace", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	err := fmt.Errorf("error")
	mkNamespaceService.EXPECT().Create(nsa).Return(nil, err)

	// 500
	req, _ = http.NewRequest(http.MethodPost, "/testA/namespace", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mkNamespaceService.EXPECT().Create(nsa).Return(nsa, nil)
	mLicense.EXPECT().GetDefaultQuotas(nsa.Name).Return(quotas, nil)
	mLicense.EXPECT().CreateQuota(nsa.Name, quotas).Return(err)
	// 200
	req, _ = http.NewRequest(http.MethodPost, "/testA/namespace", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAPI_DeleteNamespace(t *testing.T) {
	api, router, mockCtl := initNamespaceAPI(t)
	defer mockCtl.Finish()
	mkNamespaceService := ms.NewMockNamespaceService(mockCtl)
	api.NS = mkNamespaceService
	mLicense := ms.NewMockLicenseService(mockCtl)
	api.License = mLicense

	nsa := getMockNS("testA")

	mkNamespaceService.EXPECT().Delete(nsa).Return(nil)
	mLicense.EXPECT().DeleteQuotaByNamespace(nsa.Name).Return(nil)
	// 200
	req, _ := http.NewRequest(http.MethodDelete, "/testA/namespace", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	err := fmt.Errorf("error")

	mkNamespaceService.EXPECT().Delete(nsa).Return(err)

	// 500
	req, _ = http.NewRequest(http.MethodDelete, "/testA/namespace", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mkNamespaceService.EXPECT().Delete(nsa).Return(nil)
	mLicense.EXPECT().DeleteQuotaByNamespace(nsa.Name).Return(err)
	// 200
	req, _ = http.NewRequest(http.MethodDelete, "/testA/namespace", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
