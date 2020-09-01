package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	ms "github.com/baetyl/baetyl-cloud/v2/mock/service"
)

func initActiveAPI(t *testing.T) (*InitAPIImpl, *gin.Engine, *gomock.Controller) {
	api := &InitAPIImpl{}
	router := gin.Default()
	mockCtl := gomock.NewController(t)
	mockIM := func(c *gin.Context) { common.NewContext(c).SetNamespace("default") }
	v1 := router.Group("v1")
	{
		active := v1.Group("/active")
		active.GET("/:resource", mockIM, common.WrapperRaw(api.GetResource))
	}
	return api, router, mockCtl
}

func TestNewActiveAPI(t *testing.T) {
	// bad case
	_, err := NewInitAPI(&config.CloudConfig{})
	assert.Error(t, err)
}

func TestActiveAPIImpl_GetResource(t *testing.T) {
	api, router, mockCtl := initActiveAPI(t)
	defer mockCtl.Finish()
	mActive := ms.NewMockInitService(mockCtl)
	api.initService = mActive

	mActive.EXPECT().GetResource(common.ResourceSetup, "", "").Return([]byte("setup"), nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/v1/active/" + common.ResourceSetup, nil)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}