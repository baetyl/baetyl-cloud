package api

import (
	"bytes"
	"encoding/json"
	"github.com/baetyl/baetyl-cloud/v2/common"
	ms "github.com/baetyl/baetyl-cloud/v2/mock/service"
	"github.com/baetyl/baetyl-cloud/v2/models"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	defaultAppActive = &specV1.Application{
		Name:      "template-config-app-active",
		Namespace: "default",
		Version:   "12345",
		Selector:  "",
		Services: []specV1.Service{
			{
				Name:    "agent",
				Image:   "hub.baidubce.com/baetyl/baetyl-agent:1.0.0",
				Replica: 1,
				VolumeMounts: []specV1.VolumeMount{
					{
						Name:      "agent-config",
						MountPath: "etc/baetyl",
					},
					{
						Name:      "agent-log",
						MountPath: "var/log/baetyl",
					},
				},
			},
		},
		Volumes: []specV1.Volume{
			{
				Name: "agent-config",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name: "default-agent-active",
					},
				},
			},
			{
				Name: "agent-log",
				VolumeSource: specV1.VolumeSource{
					HostPath: &specV1.HostPathVolumeSource{
						Path: "var/db/baetyl/agent-log",
					},
				},
			},
		},
	}

	defaultConfigurationActive = &specV1.Configuration{
		Name:      "default-agent-active",
		Namespace: "default",
		Data: map[string]string{
			common.FilenameYamlService: `
remote:
	link:
	  address: 0.0.0.0:8273

interval: 45s
fingerprints:
    - proof: sn
      value: {{.BAETYL_FINGERPRINT_SN}}
attributes:
    - name: batch
      label: batch
      value: {{.BAETYL_BATCH_NAME}}
    - name: namespace
      label: namespace
      value: {{.BAETYL_BATCH_NAMESPACE}}
    - name: {{.BAETYL_BATCH_SECURITY_TYPE}}
      label: SecurityType
      value: {{.BAETYL_BATCH_SECURITY_KEY}}
logger:
    path: var/log/baetyl/service.log
    level: debug
`,
		},
		Version: "23456",
	}
)

// GetRegister get a register
func initRegisterAPI(t *testing.T) (*API, *gin.Engine, *gomock.Controller) {
	api := &API{}
	router := gin.Default()
	mockCtl := gomock.NewController(t)
	mockIM := func(c *gin.Context) { common.NewContext(c).SetNamespace("default") }
	v1 := router.Group("v1")
	{
		register := v1.Group("/register")
		register.POST("", mockIM, common.Wrapper(api.CreateBatch))
		register.PUT("/:batchName", mockIM, common.Wrapper(api.UpdateBatch))
		register.DELETE("/:batchName", mockIM, common.Wrapper(api.DeleteBatch))
		register.GET("/:batchName", mockIM, common.Wrapper(api.GetBatch))
		register.GET("/:batchName/init", mockIM, common.Wrapper(api.GenInitCmdFromBatch))

		register.POST("/:batchName/record", mockIM, common.Wrapper(api.CreateRecord))
		register.PUT("/:batchName/record/:recordName", mockIM, common.Wrapper(api.UpdateRecord))
		register.GET("/:batchName/record/:recordName", mockIM, common.Wrapper(api.GetRecord))
		register.DELETE("/:batchName/record/:recordName", mockIM, common.Wrapper(api.DeleteRecord))
		register.POST("/:batchName/generate", mockIM, common.Wrapper(api.GenRecordRandom))
		register.GET("", mockIM, common.Wrapper(api.ListBatch))
		register.GET("/:batchName/record", mockIM, common.Wrapper(api.ListRecord))
		register.GET("/:batchName/download", mockIM, common.Wrapper(api.DownloadRecords))

	}
	return api, router, mockCtl
}

func TestAPI_CreateBatch(t *testing.T) {
	api, router, mockCtl := initRegisterAPI(t)
	defer mockCtl.Finish()
	rs := ms.NewMockRegisterService(mockCtl)
	cs := ms.NewMockCallbackService(mockCtl)
	api.registerService = rs
	api.callbackService = cs

	mBatch := &models.Batch{
		Namespace:       "default",
		Description:     "test desc",
		QuotaNum:        20,
		EnableWhitelist: 0,
		SecurityType:    common.Token,
		SecurityKey:     "123123",
		CallbackName:    "123",
		Labels:          nil,
		Fingerprint:     models.Fingerprint{},
	}

	rs.EXPECT().CreateBatch(gomock.Any()).Return(mBatch, nil).Times(1)
	cs.EXPECT().Get(mBatch.CallbackName, mBatch.Namespace).Return(nil, nil).Times(1)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(mBatch)

	req, _ := http.NewRequest(http.MethodPost, "/v1/register", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// test for Fingerprint.Type
	mBatch.Fingerprint.Type = -1
	w = httptest.NewRecorder()
	body, _ = json.Marshal(mBatch)
	req, _ = http.NewRequest(http.MethodPost, "/v1/register", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// test for EnableWhitelist
	mBatch.Fingerprint.Type = common.FingerprintSN
	mBatch.EnableWhitelist = -1
	w = httptest.NewRecorder()
	body, _ = json.Marshal(mBatch)
	req, _ = http.NewRequest(http.MethodPost, "/v1/register", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// bad case
	// len(name)=64
	mBatch = generateDefaultBatch("test")
	mBatch.Name = "1234567812345678123456781234567812345678123456781234567812345678"
	body, _ = json.Marshal(mBatch)
	req, _ = http.NewRequest(http.MethodPost, "/v1/register", bytes.NewReader(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// bad case
	// len(label)=64
	mBatch = generateDefaultBatch("test")
	mBatch.Labels["a"] = "1234567812345678123456781234567812345678123456781234567812345678"
	body, _ = json.Marshal(mBatch)
	req, _ = http.NewRequest(http.MethodPost, "/v1/register", bytes.NewReader(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAPI_UpdateBatch(t *testing.T) {
	api, router, mockCtl := initRegisterAPI(t)
	defer mockCtl.Finish()
	rs := ms.NewMockRegisterService(mockCtl)
	api.registerService = rs

	mBatch := &models.Batch{
		Name:            "test",
		Namespace:       "default",
		Description:     "test desc",
		QuotaNum:        20,
		EnableWhitelist: 0,
		SecurityType:    common.None,
		Labels:          map[string]string{"batch": "test"},
		Fingerprint:     models.Fingerprint{},
	}

	rs.EXPECT().GetBatch("test", gomock.Any()).Return(mBatch, nil).AnyTimes()
	rs.EXPECT().UpdateBatch(gomock.Any()).Return(mBatch, nil).AnyTimes()

	updateBatch := &models.Batch{}
	w := httptest.NewRecorder()
	body, _ := json.Marshal(updateBatch)
	req, _ := http.NewRequest(http.MethodPut, "/v1/register/test", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	updateBatch = &models.Batch{
		Labels: map[string]string{"a": "b"},
	}
	body, _ = json.Marshal(updateBatch)
	req, _ = http.NewRequest(http.MethodPut, "/v1/register/test", bytes.NewReader(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAPI_UpdateBatch_Err(t *testing.T) {
	api, router, mockCtl := initRegisterAPI(t)
	defer mockCtl.Finish()
	rs := ms.NewMockRegisterService(mockCtl)
	api.registerService = rs

	rs.EXPECT().GetBatch("test", gomock.Any()).Return(nil,
		common.Error(common.ErrResourceNotFound)).Times(1)

	updateBatch := &models.Batch{}
	w := httptest.NewRecorder()
	body, _ := json.Marshal(updateBatch)
	req, _ := http.NewRequest(http.MethodPut, "/v1/register/test", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAPI_GetBatch(t *testing.T) {
	api, router, mockCtl := initRegisterAPI(t)
	defer mockCtl.Finish()
	rs := ms.NewMockRegisterService(mockCtl)
	api.registerService = rs
	rs.EXPECT().GetBatch("notExist", gomock.Any()).Return(nil, common.Error(common.ErrResourceNotFound, common.Field("type", "batch"), common.Field("name", "notExist"))).AnyTimes()

	req, _ := http.NewRequest(http.MethodGet, "/v1/register/notExist", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAPI_CreateRecord(t *testing.T) {
	api, router, mockCtl := initRegisterAPI(t)
	defer mockCtl.Finish()
	rs := ms.NewMockRegisterService(mockCtl)
	api.registerService = rs

	mRecord := models.Record{
		FingerprintValue: "123",
	}
	rs.EXPECT().CreateRecord(gomock.Any()).Return(&mRecord, nil).Times(1)

	body, _ := json.Marshal(mRecord)
	req, _ := http.NewRequest(http.MethodPost, "/v1/register/test/record", bytes.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAPI_GetRecord(t *testing.T) {
	api, router, mockCtl := initRegisterAPI(t)
	defer mockCtl.Finish()
	rs := ms.NewMockRegisterService(mockCtl)
	api.registerService = rs

	rs.EXPECT().GetRecord("test", "notExist", gomock.Any()).Return(nil, common.Error(common.ErrResourceNotFound, common.Field("type", "record"), common.Field("name", "notExist"))).AnyTimes()

	req, _ := http.NewRequest(http.MethodGet, "/v1/register/test/record/notExist", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAPI_DeleteBatch(t *testing.T) {
	api, router, mockCtl := initRegisterAPI(t)
	defer mockCtl.Finish()
	rs := ms.NewMockRegisterService(mockCtl)
	api.registerService = rs

	batch := &models.Batch{
		Name:      "test",
		Namespace: "default",
	}
	rs.EXPECT().DeleteBatch(batch.Name, batch.Namespace).Return(nil).Times(1)

	req, _ := http.NewRequest(http.MethodDelete, "/v1/register/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAPI_UpdateRecord(t *testing.T) {
	api, router, mockCtl := initRegisterAPI(t)
	defer mockCtl.Finish()
	rs := ms.NewMockRegisterService(mockCtl)
	api.registerService = rs

	batch := &models.Batch{
		Name:         "test",
		Namespace:    "default",
		Description:  "test desc",
		SecurityType: common.Token,
		SecurityKey:  "token",
	}
	record := &models.Record{
		Name:             "123",
		Namespace:        "default",
		BatchName:        "test",
		FingerprintValue: "123",
		NodeName:         "123",
	}
	rs.EXPECT().GetBatch(batch.Name, batch.Namespace).Return(nil, nil).Times(1)
	rs.EXPECT().GetRecord(batch.Name, record.Name, record.Namespace).Return(record, nil).Times(1)
	rs.EXPECT().UpdateRecord(record).Return(record, nil).Times(1)

	body, _ := json.Marshal(record)
	req, _ := http.NewRequest(http.MethodPut, "/v1/register/test/record/123", bytes.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAPI_GenInitCmdFromBatch(t *testing.T) {
	api, router, mockCtl := initRegisterAPI(t)
	defer mockCtl.Finish()
	rs := ms.NewMockRegisterService(mockCtl)
	auth := ms.NewMockAuthService(mockCtl)
	ss := ms.NewMockSysConfigService(mockCtl)
	api.authService = auth
	api.sysConfigService = ss
	api.registerService = rs

	sc := &models.SysConfig{Key: common.AddressActive, Type: "address", Value: "baetyl.com"}

	// good case
	rs.EXPECT().GetBatch("baetyl-default-batch", "default").Return(nil, nil).Times(1)
	auth.EXPECT().GenToken(gomock.Any()).Return("token", nil).Times(1)
	ss.EXPECT().GetSysConfig(sc.Type, sc.Key).Return(sc, nil).Times(1)

	req, _ := http.NewRequest(http.MethodGet, "/v1/register/baetyl-default-batch/init", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// bad case 0
	rs.EXPECT().GetBatch("baetyl-default-batch", "default").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	req, _ = http.NewRequest(http.MethodGet, "/v1/register/baetyl-default-batch/init", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	// bad case 1
	err := common.Error(
		common.ErrResourceNotFound,
		common.Field("type", "test"),
		common.Field("name", "test"))
	rs.EXPECT().GetBatch("baetyl-default-batch", "default").Return(nil, nil).Times(1)
	auth.EXPECT().GenToken(gomock.Any()).Return("token", nil).Times(1)
	ss.EXPECT().GetSysConfig(sc.Type, sc.Key).Return(nil, err).Times(1)
	req, _ = http.NewRequest(http.MethodGet, "/v1/register/baetyl-default-batch/init", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAPI_DeleteRecord(t *testing.T) {
	api, router, mockCtl := initRegisterAPI(t)
	defer mockCtl.Finish()
	rs := ms.NewMockRegisterService(mockCtl)
	api.registerService = rs

	rs.EXPECT().DeleteRecord("test", "123", "default").Return(nil).Times(1)

	req, _ := http.NewRequest(http.MethodDelete, "/v1/register/test/record/123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAPI_GenRecordRandom(t *testing.T) {
	api, router, mockCtl := initRegisterAPI(t)
	defer mockCtl.Finish()
	rs := ms.NewMockRegisterService(mockCtl)
	api.registerService = rs

	rs.EXPECT().GenRecordRandom("default", "test", 2).Return([]string{"123", "456"}, nil).Times(1)

	body, _ := json.Marshal(struct {
		Num int
	}{Num: 2})
	req, _ := http.NewRequest(http.MethodPost, "/v1/register/test/generate", bytes.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAPI_ListBatch(t *testing.T) {
	api, router, mockCtl := initRegisterAPI(t)
	defer mockCtl.Finish()
	rs := ms.NewMockRegisterService(mockCtl)
	api.registerService = rs

	page := &models.Filter{
		PageNo:   1,
		PageSize: 2,
		Name:     "%2%",
	}

	batch := generateDefaultBatch("default")
	rs.EXPECT().ListBatch("default", page).Return(&models.ListView{
		Total:    1,
		PageNo:   1,
		PageSize: 2,
		Items:    []models.Batch{*batch},
	}, nil).Times(1)

	req, _ := http.NewRequest(http.MethodGet, "/v1/register?name=2&pageNo=1&pageSize=2", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAPI_ListRecord(t *testing.T) {
	api, router, mockCtl := initRegisterAPI(t)
	defer mockCtl.Finish()
	rs := ms.NewMockRegisterService(mockCtl)
	api.registerService = rs

	page := &models.Filter{
		PageNo:   1,
		PageSize: 20,
		Name:     "%",
	}

	record := models.Record{
		Name:             "r1",
		Namespace:        "default",
		BatchName:        "test",
		FingerprintValue: "qweqwe",
	}
	rs.EXPECT().ListRecord(record.BatchName, "default", page).Return(&models.ListView{
		Total:    1,
		PageNo:   1,
		PageSize: 20,
		Items:    []models.Record{record},
	}, nil).Times(1)

	req, _ := http.NewRequest(http.MethodGet, "/v1/register/test/record", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAPI_DownloadRecords(t *testing.T) {
	api, router, mockCtl := initRegisterAPI(t)
	defer mockCtl.Finish()
	rs := ms.NewMockRegisterService(mockCtl)
	api.registerService = rs

	rs.EXPECT().DownloadRecords("test", "default").Return(nil, nil).Times(1)

	req, _ := http.NewRequest(http.MethodGet, "/v1/register/test/download", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
