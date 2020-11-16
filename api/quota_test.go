package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/baetyl/baetyl-cloud/v2/common"
	ms "github.com/baetyl/baetyl-cloud/v2/mock/service"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var namespace = "default"

func initQuotaAPI(t *testing.T) (*API, *gin.Engine, *gomock.Controller) {
	api := &API{}
	router := gin.Default()
	mockCtl := gomock.NewController(t)
	mockIM := func(c *gin.Context) { common.NewContext(c).SetNamespace(namespace) }
	v1 := router.Group("v1")
	{
		quota := v1.Group("/quotas")
		quota.GET("", mockIM, common.Wrapper(api.GetQuota))

		quota.POST("", common.WrapperMis(api.CreateQuota))
		quota.DELETE("", common.WrapperMis(api.DeleteQuota))
		quota.GET("/mis", common.WrapperMis(api.GetQuotaForMis))
		quota.PUT("", common.WrapperMis(api.UpdateQuota))
	}

	return api, router, mockCtl
}

func TestAPI_GetQuota(t *testing.T) {
	api, router, mockCtl := initQuotaAPI(t)
	defer mockCtl.Finish()

	mLicense := ms.NewMockLicenseService(mockCtl)
	api.License = mLicense

	quotas := map[string]int{
		"maxNodeCount": 10,
	}

	mLicense.EXPECT().GetQuota(namespace).Return(quotas, nil)
	// 200
	req, _ := http.NewRequest(http.MethodGet, "/v1/quotas", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	result, err := ioutil.ReadAll(w.Body)
	assert.NoError(t, err)
	actual := map[string]int{}
	err = json.Unmarshal(result, &actual)
	assert.NoError(t, err)
	assert.Equal(t, quotas, actual)
}

func TestAPI_UpdateQuota(t *testing.T) {
	api, router, mockCtl := initQuotaAPI(t)
	defer mockCtl.Finish()

	mLicense := ms.NewMockLicenseService(mockCtl)
	api.License = mLicense

	number := 10
	body := `{"quotaName":"maxNodeCount","quota": 10,"namespace":"default"}`
	mLicense.EXPECT().UpdateQuota(namespace, plugin.QuotaNode, number).Return(nil)
	// 200
	req, _ := http.NewRequest(http.MethodPut, "/v1/quotas", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

}

func TestAPI_CreateQuota(t *testing.T) {
	api, router, mockCtl := initQuotaAPI(t)
	defer mockCtl.Finish()

	mLicense := ms.NewMockLicenseService(mockCtl)
	api.License = mLicense

	quotas := map[string]int{
		"maxNodeCount": 10,
	}
	body := `{"quotaName":"maxNodeCount","quota": 10,"namespace":"default"}`

	mLicense.EXPECT().CreateQuota(namespace, quotas).Return(nil)
	// 200
	req, _ := http.NewRequest(http.MethodPost, "/v1/quotas", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAPI_DeleteQuota(t *testing.T) {
	api, router, mockCtl := initQuotaAPI(t)
	defer mockCtl.Finish()

	mLicense := ms.NewMockLicenseService(mockCtl)
	api.License = mLicense

	body := `{"quotaName":"maxNodeCount","quota": 10,"namespace":"default"}`

	mLicense.EXPECT().DeleteQuota(namespace, plugin.QuotaNode).Return(nil)
	// 200
	req, _ := http.NewRequest(http.MethodDelete, "/v1/quotas", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAPI_GetQuotaForMis(t *testing.T) {
	api, router, mockCtl := initQuotaAPI(t)
	defer mockCtl.Finish()

	mLicense := ms.NewMockLicenseService(mockCtl)
	api.License = mLicense

	quotas := map[string]int{
		"maxNodeCount": 10,
	}

	body := `{"quotaName":"maxNodeCount","quota": 10,"namespace":"default"}`

	mLicense.EXPECT().GetQuota(namespace).Return(quotas, nil)
	// 200
	req, _ := http.NewRequest(http.MethodGet, "/v1/quotas/mis", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAPI_InitQuotas(t *testing.T) {
	api, _, mockCtl := initQuotaAPI(t)
	defer mockCtl.Finish()

	mLicense := ms.NewMockLicenseService(mockCtl)
	api.License = mLicense

	namespace := "testError"
	testErr := fmt.Errorf("testError")

	mLicense.EXPECT().GetDefaultQuotas(namespace).Return(nil, testErr)
	err := api.InitQuotas(namespace)
	assert.Error(t, testErr, err)

	namespace = "testSucc"
	quotas := map[string]int{plugin.QuotaNode: 10}

	mLicense.EXPECT().GetDefaultQuotas(namespace).Return(quotas, nil)
	mLicense.EXPECT().CreateQuota(namespace, quotas).Return(testErr)
	err = api.InitQuotas(namespace)
	assert.Error(t, testErr, err)

	mLicense.EXPECT().GetDefaultQuotas(namespace).Return(quotas, nil)
	mLicense.EXPECT().CreateQuota(namespace, quotas).Return(nil)
	err = api.InitQuotas(namespace)
	assert.NoError(t, err)
}

func TestAPI_DeleteAllQuotas(t *testing.T) {
	api, _, mockCtl := initQuotaAPI(t)
	defer mockCtl.Finish()

	mLicense := ms.NewMockLicenseService(mockCtl)
	api.License = mLicense

	namespace := "testError"
	testErr := fmt.Errorf("testError")

	mLicense.EXPECT().DeleteQuotaByNamespace(namespace).Return(testErr)
	err := api.DeleteQuotaByNamespace(namespace)
	assert.Error(t, testErr, err)
	namespace = "testSucc"

	mLicense.EXPECT().DeleteQuotaByNamespace(namespace).Return(nil)
	err = api.DeleteQuotaByNamespace(namespace)
	assert.NoError(t, err)
}

func TestAPI_RealseNodeQuota(t *testing.T) {
	api, _, mockCtl := initQuotaAPI(t)
	defer mockCtl.Finish()

	mLicense := ms.NewMockLicenseService(mockCtl)
	api.License = mLicense

	namespace := "testError"
	number := 1
	testErr := fmt.Errorf("testError")

	mLicense.EXPECT().ReleaseQuota(namespace, plugin.QuotaNode, number).Return(testErr)
	err := api.ReleaseQuota(namespace, plugin.QuotaNode, number)
	assert.Error(t, testErr, err)
	namespace = "testSucc"

	mLicense.EXPECT().ReleaseQuota(namespace, plugin.QuotaNode, number).Return(nil)
	err = api.ReleaseQuota(namespace, plugin.QuotaNode, number)
	assert.NoError(t, err)
}
