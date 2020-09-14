package api

import (
	"encoding/json"
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
		configs := v1.Group("/quotas")
		configs.GET("", mockIM, common.Wrapper(api.GetQuota))
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
