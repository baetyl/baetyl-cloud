package api

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
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

func initInitAPI(t *testing.T) (*InitAPI, *gin.Engine, *gomock.Controller) {
	api := &InitAPI{}
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
	auth := ms.NewMockAuthService(mockCtl)
	api.Auth = auth
	// 构造token
	info := map[string]interface{}{
		service.InfoName:      "n0",
		service.InfoNamespace: "default",
		service.InfoExpiry:    time.Now().Unix() + 60*60*24*3650,
	}
	data, err := json.Marshal(info)
	assert.NoError(t, err)
	encode := hex.EncodeToString(data)
	sign := "0123456789"
	token := sign + encode
	// 构造form
	sendUrl, _ := url.Parse("/v1/init/kube-init-setup.sh?")
	val := sendUrl.Query()
	val.Set("token", token)
	sendUrl.RawQuery = val.Encode()

	// ResourceSetup
	mInit.EXPECT().GetResource("default", "n0", "kube-init-setup.sh", gomock.Any()).Return([]byte("setup"), nil)
	auth.EXPECT().GenToken(gomock.Any()).Return(token, nil).Times(3)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, sendUrl.String(), nil)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// ResourceMetrics
	sendUrl, _ = url.Parse("/v1/init/" + "kube-api-metrics.yml?")
	val = sendUrl.Query()
	val.Set("token", token)
	sendUrl.RawQuery = val.Encode()

	mInit.EXPECT().GetResource("default", "n0", "kube-api-metrics.yml", gomock.Any()).Return([]byte("metrics"), nil)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, sendUrl.String(), nil)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// ResourceLocalPathStorage
	sendUrl, _ = url.Parse("/v1/init/" + "kube-local-path-storage.yml?")
	val = sendUrl.Query()
	val.Set("token", token)
	sendUrl.RawQuery = val.Encode()

	mInit.EXPECT().GetResource("default", "n0", "kube-local-path-storage.yml", gomock.Any()).Return([]byte("metrics"), nil)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, sendUrl.String(), nil)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestInitAPIImpl_CheckAndParseToken(t *testing.T) {
	as := InitAPI{}
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	auth := ms.NewMockAuthService(mockCtl)
	as.Auth = auth
	info := map[string]interface{}{
		service.InfoName:      "n0",
		service.InfoNamespace: "default",
		service.InfoExpiry:    time.Now().Unix() + 60*60*24*3650,
	}
	data, err := json.Marshal(info)
	assert.NoError(t, err)
	encode := hex.EncodeToString(data)
	sign := "0123456789"
	token := sign + encode

	auth.EXPECT().GenToken(gomock.Any()).Return(token, nil).Times(1)

	res, err := CheckAndParseToken(token, as.Auth.GenToken)
	assert.NoError(t, err)
	assert.Equal(t, info[service.InfoName], res[service.InfoName].(string))
	assert.Equal(t, info[service.InfoNamespace], res[service.InfoNamespace].(string))
}
