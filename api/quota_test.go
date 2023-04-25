package api

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/baetyl/baetyl-go/v2/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/common"
	ms "github.com/baetyl/baetyl-cloud/v2/mock/service"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
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
		quota.GET("/:namespace", common.WrapperMis(api.GetQuotaForMis))
		quota.PUT("", common.WrapperMis(api.UpdateQuota))
	}

	return api, router, mockCtl
}

func TestAPI_GetQuota(t *testing.T) {
	api, router, mockCtl := initQuotaAPI(t)
	defer mockCtl.Finish()

	mQuota := ms.NewMockQuotaService(mockCtl)
	api.Quota = mQuota

	quotas := map[string]int{
		"maxNodeCount": 10,
	}

	mQuota.EXPECT().GetQuota(namespace).Return(quotas, nil)
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

	mQuota := ms.NewMockQuotaService(mockCtl)
	api.Quota = mQuota

	number := 10
	body := `{"quotaName":"maxNodeCount","quota": 10,"namespace":"default"}`
	mQuota.EXPECT().UpdateQuota(namespace, plugin.QuotaNode, number).Return(nil)
	// 200
	req, _ := http.NewRequest(http.MethodPut, "/v1/quotas", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

}

func TestAPI_CreateQuota(t *testing.T) {
	api, router, mockCtl := initQuotaAPI(t)
	defer mockCtl.Finish()

	mQuota := ms.NewMockQuotaService(mockCtl)
	api.Quota = mQuota

	quotas := map[string]int{
		"maxNodeCount": 10,
	}
	body := `{"quotaName":"maxNodeCount","quota": 10,"namespace":"default"}`

	mQuota.EXPECT().CreateQuota(namespace, quotas).Return(nil)
	// 200
	req, _ := http.NewRequest(http.MethodPost, "/v1/quotas", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAPI_DeleteQuota(t *testing.T) {
	api, router, mockCtl := initQuotaAPI(t)
	defer mockCtl.Finish()

	mQuota := ms.NewMockQuotaService(mockCtl)
	api.Quota = mQuota

	body := `{"quotaName":"maxNodeCount","quota": 10,"namespace":"default"}`

	mQuota.EXPECT().DeleteQuota(namespace, plugin.QuotaNode).Return(nil)
	// 200
	req, _ := http.NewRequest(http.MethodDelete, "/v1/quotas", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAPI_GetQuotaForMis(t *testing.T) {
	api, router, mockCtl := initQuotaAPI(t)
	defer mockCtl.Finish()

	mQuota := ms.NewMockQuotaService(mockCtl)
	api.Quota = mQuota

	quotas := map[string]int{
		"maxNodeCount": 10,
	}

	body := `{"quotaName":"maxNodeCount","quota": 10,"namespace":"default"}`

	mQuota.EXPECT().GetQuota(namespace).Return(quotas, nil)
	// 200
	req, _ := http.NewRequest(http.MethodGet, "/v1/quotas/default", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAPI_InitQuotas(t *testing.T) {
	api, _, mockCtl := initQuotaAPI(t)
	defer mockCtl.Finish()

	mQuota := ms.NewMockQuotaService(mockCtl)
	api.Quota = mQuota

	ns := "testError"
	testErr := fmt.Errorf("testError")

	mQuota.EXPECT().GetDefaultQuotas(ns).Return(nil, testErr)
	err := api.InitQuotas(ns)
	assert.Error(t, testErr, err)

	ns = "testSucc"
	quotas := map[string]int{plugin.QuotaNode: 10}

	mQuota.EXPECT().GetDefaultQuotas(ns).Return(quotas, nil)
	mQuota.EXPECT().CreateQuota(ns, quotas).Return(testErr)
	err = api.InitQuotas(ns)
	assert.Error(t, testErr, err)

	mQuota.EXPECT().GetDefaultQuotas(ns).Return(quotas, nil)
	mQuota.EXPECT().CreateQuota(ns, quotas).Return(nil)
	err = api.InitQuotas(ns)
	assert.NoError(t, err)
}

func TestAPI_DeleteAllQuotas(t *testing.T) {
	api, _, mockCtl := initQuotaAPI(t)
	defer mockCtl.Finish()

	mQuota := ms.NewMockQuotaService(mockCtl)
	api.Quota = mQuota

	ns := "testError"
	testErr := fmt.Errorf("testError")

	mQuota.EXPECT().DeleteQuotaByNamespace(ns).Return(testErr)
	err := api.DeleteQuotaByNamespace(ns)
	assert.Error(t, testErr, err)
	ns = "testSucc"

	mQuota.EXPECT().DeleteQuotaByNamespace(ns).Return(nil)
	err = api.DeleteQuotaByNamespace(ns)
	assert.NoError(t, err)
}

func TestAPI_RealseNodeQuota(t *testing.T) {
	api, _, mockCtl := initQuotaAPI(t)
	defer mockCtl.Finish()

	mQuota := ms.NewMockQuotaService(mockCtl)
	api.Quota = mQuota

	ns := "testError"
	number := 1
	testErr := fmt.Errorf("testError")

	mQuota.EXPECT().ReleaseQuota(ns, plugin.QuotaNode, number).Return(testErr)
	err := api.ReleaseQuota(ns, plugin.QuotaNode, number)
	assert.Error(t, testErr, err)
	ns = "testSucc"

	mQuota.EXPECT().ReleaseQuota(ns, plugin.QuotaNode, number).Return(nil)
	err = api.ReleaseQuota(ns, plugin.QuotaNode, number)
	assert.NoError(t, err)
}
