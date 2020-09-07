package api

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	ms "github.com/baetyl/baetyl-cloud/v2/mock/service"
	"github.com/baetyl/baetyl-cloud/v2/service"
)

func initInitAPI(t *testing.T) (*InitAPIImpl, *gin.Engine, *gomock.Controller) {
	api := &InitAPIImpl{}
	router := gin.Default()
	mockCtl := gomock.NewController(t)
	mockIM := func(c *gin.Context) { common.NewContext(c).SetNamespace("default") }
	v1 := router.Group("v1")
	{
		init := v1.Group("/init")
		init.GET("/:resource", mockIM, common.WrapperRaw(api.GetResource))
	}
	return api, router, mockCtl
}

func TestNewInitAPI(t *testing.T) {
	// bad case
	_, err := NewInitAPI(&config.CloudConfig{})
	assert.Error(t, err)
}

func TestInitAPIImpl_GetResource(t *testing.T) {
	api, router, mockCtl := initInitAPI(t)
	defer mockCtl.Finish()
	mInit := ms.NewMockInitService(mockCtl)
	api.Init = mInit

	// ResourceSetup
	mInit.EXPECT().GetResource("kube-init-setup.sh", "", "", nil).Return([]byte("setup"), nil)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/v1/init/kube-init-setup.sh", nil)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// ResourceMetrics
	mInit.EXPECT().GetResource("kube-api-metrics.yml", "", "", nil).Return([]byte("metrics"), nil)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/v1/init/"+"kube-api-metrics.yml", nil)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// ResourceLocalPathStorage
	mInit.EXPECT().GetResource("kube-local-path-storage.yml", "", "", nil).Return([]byte("metrics"), nil)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/v1/init/"+"kube-local-path-storage.yml", nil)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestInitAPIImpl_CheckAndParseToken(t *testing.T) {
	as := InitAPIImpl{}
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	auth := ms.NewMockAuthService(mockCtl)
	as.Auth = auth
	info := map[string]interface{}{
		service.InfoKind:      "node",
		service.InfoName:      "n0",
		service.InfoNamespace: "default",
		service.InfoTimestamp: time.Now().Unix(),
		service.InfoExpiry:    60 * 60 * 24 * 3650,
	}
	data, err := json.Marshal(info)
	assert.NoError(t, err)
	encode := hex.EncodeToString(data)
	sign := "0123456789"
	token := sign + encode

	auth.EXPECT().GenToken(gomock.Any()).Return(token, nil).Times(1)

	res, err := as.CheckAndParseToken(token, "baetyl-init-deployment.yml")

	assert.Equal(t, res[service.InfoKind], info[service.InfoKind])
	assert.Equal(t, res[service.InfoName], info[service.InfoName])
	assert.Equal(t, res[service.InfoNamespace], info[service.InfoNamespace])
	assert.NoError(t, err)
}
