package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/baetyl/baetyl-cloud/common"
	ms "github.com/baetyl/baetyl-cloud/mock/service"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func initSyncAPI(t *testing.T) (*API, *gin.Engine, *gomock.Controller) {
	mockCtl := gomock.NewController(t)

	api := &API{}
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		cc := common.NewContext(c)
		cc.SetNamespace("default")
		cc.SetName("test")
	})
	v1 := router.Group("v1")
	{
		sync := v1.Group("/sync")
		sync.POST("/report", common.Wrapper(api.Report))
		sync.POST("/desire", common.Wrapper(api.Desire))
	}
	return api, router, mockCtl
}

func TestReport(t *testing.T) {
	api, router, mockCtl := initSyncAPI(t)
	defer mockCtl.Finish()

	mSync := ms.NewMockSyncService(mockCtl)
	api.syncService = mSync

	// TODO: use real tls cert
	// info := &specV1.Report{}
	// data, err := json.Marshal(info)
	// assert.NoError(t, err)
	// r := bytes.NewReader(data)
	// req, err := http.NewRequest(http.MethodPost, "/v1/sync/report", r)
	// w := httptest.NewRecorder()
	// router.ServeHTTP(w, req)
	// assert.Equal(t, http.StatusOK, w.Code)
	// assert.NoError(t, err)

	response := specV1.Desire{}
	mSync.EXPECT().Report(gomock.Any(), gomock.Any(), gomock.Any()).Return(response, nil)
	info := &specV1.Report{
		"apps": []specV1.AppInfo{
			{
				Name:    "app01",
				Version: "v1",
			},
		},
	}
	data, err := json.Marshal(info)
	r := bytes.NewReader(data)
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/v1/sync/report", r)
	router.ServeHTTP(w, req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, w.Code)

	report := &specV1.Report{}
	data, _ = json.Marshal(report)

	mSync.EXPECT().Report(gomock.Any(), gomock.Any(), gomock.Any()).Return(response, nil)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/v1/sync/report", bytes.NewReader(data))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDesire(t *testing.T) {
	api, router, mockCtl := initSyncAPI(t)
	defer mockCtl.Finish()
	mSync := ms.NewMockSyncService(mockCtl)
	api.syncService = mSync
	var response []specV1.ResourceValue
	var request []specV1.ResourceInfo
	mSync.EXPECT().Desire(gomock.Any(), gomock.Any()).Return(response, nil)
	data, err := json.Marshal(request)
	assert.NoError(t, err)
	r := bytes.NewReader(data)

	req, err := http.NewRequest(http.MethodPost, "/v1/sync/desire", r)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
