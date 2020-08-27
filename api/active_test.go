package api

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/mock/service"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func initActiveAPI(t *testing.T) (*API, *gin.Engine, *gomock.Controller) {
	api := &API{}
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

func TestAPI_GetResource(t *testing.T) {
	api, router, ctl := initActiveAPI(t)
	init := service.NewMockInitializeService(ctl)
	ss := service.NewMockSysConfigService(ctl)
	tp := service.NewMockTemplateService(ctl)
	api.sysConfigService = ss
	api.initService = init
	api.templateService = tp

	// good case : metrics
	init.EXPECT().GetResource(common.ResourceMetrics).Return("metrics", nil).Times(1)
	tp.EXPECT().GenSetupShell(gomock.Any()).Return([]byte("shell"), nil).Times(1)

	req, _ := http.NewRequest(http.MethodGet, "/v1/active/"+common.ResourceMetrics, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// good case : local_path_storage
	init.EXPECT().GetResource(common.ResourceLocalPathStorage).Return("local-path-storage", nil).Times(1)

	req, _ = http.NewRequest(http.MethodGet, "/v1/active/"+common.ResourceLocalPathStorage, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// good case : setup
	sc := &models.SysConfig{
		Type:  "address",
		Key:   common.AddressActive,
		Value: "baetyl.com",
	}
	ss.EXPECT().GetSysConfig(sc.Type, sc.Key).Return(sc, nil).Times(1)
	init.EXPECT().GetResource(common.ResourceSetup).Return("{}", nil).Times(1)

	req, _ = http.NewRequest(http.MethodGet, "/v1/active/"+common.ResourceSetup, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// bad case : init
	req, _ = http.NewRequest(http.MethodGet, "/v1/active/baetyl-init.yml", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// bad case : not found
	req, _ = http.NewRequest(http.MethodGet, "/v1/active/notfound", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestApi_genInitYaml(t *testing.T) {
	// expiry token
	token := "ac40cc632e217d7675abfdfbf64e285f7b22657870697279223a333630302c226b696e64223a226e6f6465222c226e616d65223a22303431353031222c226e616d657370616365223a2264656661756c74222c2274696d657374616d70223a313538363935363931367d"
	kube := "k3s"
	api, _, ctl := initActiveAPI(t)
	auth := service.NewMockAuthService(ctl)
	api.authService = auth
	auth.EXPECT().GenToken(gomock.Any()).Return(token, nil).Times(1)
	res, err := api.getInitYaml(token, kube)
	assert.Error(t, err, ErrInvalidToken)
	assert.Nil(t, res)
}

func TestApi_genCmd(t *testing.T) {
	token := "ac40cc632e217d7675abfdfbf64e285f7b22657870697279223a333630302c226b696e64223a226e6f6465222c226e616d65223a22303431353031222c226e616d657370616365223a2264656661756c74222c2274696d657374616d70223a313538363935363931367d"
	api, _, ctl := initActiveAPI(t)
	ss := service.NewMockSysConfigService(ctl)
	auth := service.NewMockAuthService(ctl)
	api.authService = auth
	api.sysConfigService = ss

	// bad case 0: gen Token error
	auth.EXPECT().GenToken(gomock.Any()).Return("", fmt.Errorf("gen token err")).Times(1)

	res, err := api.genCmd("batch", "default", "test")
	assert.Error(t, err)
	assert.Equal(t, "", res)

	// bad case 1: get sys config error
	auth.EXPECT().GenToken(gomock.Any()).Return(token, nil).Times(1)
	ss.EXPECT().GetSysConfig("address", common.AddressActive).Return(nil, fmt.Errorf("not found")).Times(1)

	res, err = api.genCmd("batch", "default", "test")
	assert.Error(t, err, common.Error(common.ErrResourceNotFound,
		common.Field("type", "address"),
		common.Field("name", common.AddressActive)))
	assert.Equal(t, "", res)
}

func TestAPI_getInitYaml(t *testing.T) {
	api, _, ctl := initActiveAPI(t)
	as := service.NewMockAuthService(ctl)
	init := service.NewMockInitializeService(ctl)
	api.authService = as
	api.initService = init

	info := map[string]interface{}{
		InfoKind:      "node",
		InfoName:      "n0",
		InfoNamespace: "default",
		InfoTimestamp: time.Now().Unix(),
		InfoExpiry:    60 * 60 * 24 * 3650,
	}
	data, err := json.Marshal(info)
	assert.NoError(t, err)
	encode := hex.EncodeToString(data)
	sign := "0123456789"
	token := sign + encode
	kube := "k3s"

	// good case 0
	as.EXPECT().GenToken(gomock.Any()).Return(token, nil).Times(1)
	init.EXPECT().InitWithNode("default", "n0", kube).Return(nil, nil).Times(1)
	_, err = api.getInitYaml(token, kube)
	assert.NoError(t, err)

	// bad case 0
	info[InfoKind] = "error"
	data, err = json.Marshal(info)
	assert.NoError(t, err)
	encode = hex.EncodeToString(data)
	token = sign + encode
	as.EXPECT().GenToken(gomock.Any()).Return(token, nil).Times(1)
	_, err = api.getInitYaml(token, kube)
	assert.Error(t, err)
}
