package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/baetyl/baetyl-cloud/v2/common"
	ms "github.com/baetyl/baetyl-cloud/v2/mock/service"
	"github.com/baetyl/baetyl-cloud/v2/models"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// TODO: optimize this layer, general abstraction

// GetConfig get a config
func initConfigAPI(t *testing.T) (*API, *gin.Engine, *gomock.Controller) {
	api := &API{}
	router := gin.Default()
	mockCtl := gomock.NewController(t)
	mockIM := func(c *gin.Context) {
		common.NewContext(c).SetNamespace("default")
		common.NewContext(c).SetUser(common.User{ID: "default"})
	}
	v1 := router.Group("v1")
	{
		configs := v1.Group("/configs")
		configs.GET("/:name", mockIM, common.Wrapper(api.GetConfig))
		configs.GET("/:name/apps", mockIM, common.Wrapper(api.GetAppByConfig))
		configs.PUT("/:name", mockIM, common.Wrapper(api.UpdateConfig))
		configs.DELETE("/:name", mockIM, common.Wrapper(api.DeleteConfig))
		configs.POST("", mockIM, common.Wrapper(api.CreateConfig))
		configs.GET("", mockIM, common.Wrapper(api.ListConfig))
	}

	return api, router, mockCtl
}

func TestGetConfig(t *testing.T) {
	api, router, mockCtl := initConfigAPI(t)
	defer mockCtl.Finish()
	mkConfigService := ms.NewMockConfigService(mockCtl)
	api.configService = mkConfigService

	mConf := &specV1.Configuration{
		Namespace: "default",
		Name:      "abc",
	}

	mkConfigService.EXPECT().Get(mConf.Namespace, mConf.Name, "").Return(mConf, nil)
	mkConfigService.EXPECT().Get(mConf.Namespace, "cba", "").Return(nil, fmt.Errorf("error"))

	// 200
	req, _ := http.NewRequest(http.MethodGet, "/v1/configs/abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 404
	req, _ = http.NewRequest(http.MethodGet, "/v1/configs/cba", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req)
	assert.Equal(t, http.StatusInternalServerError, w2.Code)
}

func TestListConfig(t *testing.T) {
	api, router, mockCtl := initConfigAPI(t)
	defer mockCtl.Finish()
	mkConfigService := ms.NewMockConfigService(mockCtl)
	api.configService = mkConfigService

	mClist := &models.ConfigurationList{}

	mkConfigService.EXPECT().List("default", &models.ListOptions{
		LabelSelector: "!" + common.LabelSystem,
	}).Return(mClist, nil)

	// 200
	req, _ := http.NewRequest(http.MethodGet, "/v1/configs", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	mkConfigService.EXPECT().List("default", &models.ListOptions{
		LabelSelector: "!" + common.LabelSystem,
	}).Return(nil, fmt.Errorf("error"))

	// 500
	req, _ = http.NewRequest(http.MethodGet, "/v1/configs", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateConfig(t *testing.T) {
	api, router, mockCtl := initConfigAPI(t)
	defer mockCtl.Finish()
	mkConfigService := ms.NewMockConfigService(mockCtl)
	api.configService = mkConfigService

	mConf := &models.ConfigurationView{
		Name:      "abc",
		Namespace: "default",
		Labels: map[string]string{
			"test": "test",
		},
		Data: []models.ConfigDataItem{
			{
				Key: "function",
				Value: map[string]string{
					"type":    "function",
					"unknown": "unknown",
				},
			},
		},
	}

	w := httptest.NewRecorder()
	body, _ := json.Marshal(mConf)
	req, _ := http.NewRequest(http.MethodPost, "/v1/configs", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mConf = &models.ConfigurationView{
		Name:      "abc",
		Namespace: "default",
		Labels: map[string]string{
			"test": "test",
		},
		Data: []models.ConfigDataItem{
			{
				Key: "object",
				Value: map[string]string{
					"type":    "object",
					"unknown": "unknown",
				},
			},
		},
	}

	w = httptest.NewRecorder()
	body, _ = json.Marshal(mConf)
	req, _ = http.NewRequest(http.MethodPost, "/v1/configs", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mConf = &models.ConfigurationView{
		Name:      "abc",
		Namespace: "default",
		Labels: map[string]string{
			"test": "test",
		},
		Data: []models.ConfigDataItem{
			{
				Key: "function",
				Value: map[string]string{
					"type":     "function",
					"function": "process",
					"version":  "1",
					"runtime":  "python36",
					"handler":  "index.handler",
					"bucket":   "baetyl",
					"object":   "a.zip",
				},
			},
		},
	}

	res := &specV1.Configuration{}
	mkConfigService.EXPECT().Get(mConf.Namespace, mConf.Name, gomock.Any()).Return(res, nil).Times(1)

	// 400: configuration already exist
	w = httptest.NewRecorder()
	body, _ = json.Marshal(mConf)
	req, _ = http.NewRequest(http.MethodPost, "/v1/configs", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mConf = &models.ConfigurationView{
		Name:      "abc",
		Namespace: "default",
		Labels: map[string]string{
			"test": "test",
		},
		Data: []models.ConfigDataItem{
			{
				Key: "function",
				Value: map[string]string{
					"type":     "function",
					"function": "process",
					"version":  "1",
					"runtime":  "python36",
					"handler":  "index.handler",
					"bucket":   "baetyl",
					"object":   "a.zip",
				},
			},
		},
	}

	res = &specV1.Configuration{
		Name:      "abc",
		Namespace: "default",
		Labels: map[string]string{
			"test": "test",
		},
		Data: map[string]string{
			common.ConfigObjectPrefix + "function": `{"metadata":{"bucket":"baetyl","function":"process","handler":"index.handler","object":"a.zip","runtime":"python36","type":"function","version":"1"}}`,
		},
		CreationTimestamp: time.Now(),
		UpdateTimestamp:   time.Now(),
		Description:       "des",
		Version:           "12",
		System:            false,
	}
	mkConfigService.EXPECT().Get(mConf.Namespace, mConf.Name, gomock.Any()).Return(nil, nil).Times(1)
	mkConfigService.EXPECT().Create(mConf.Namespace, gomock.Any()).Return(res, nil).Times(1)

	w = httptest.NewRecorder()
	body, _ = json.Marshal(mConf)
	req, _ = http.NewRequest(http.MethodPost, "/v1/configs", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var view models.ConfigurationView
	err := json.Unmarshal(w.Body.Bytes(), &view)
	assert.NoError(t, err)
	assert.Equal(t, view.Name, res.Name)
	assert.Equal(t, view.Namespace, res.Namespace)
	assert.Equal(t, view.Labels, res.Labels)
	assert.Equal(t, view.Data[0].Key, "function")
	assert.Equal(t, view.Data[0].Value["type"], "function")
	assert.Equal(t, view.Data[0].Value["bucket"], "baetyl")
	assert.Equal(t, view.Data[0].Value["object"], "a.zip")
	assert.Equal(t, view.Data[0].Value["function"], "process")
	assert.Equal(t, view.Data[0].Value["handler"], "index.handler")
	assert.Equal(t, view.Data[0].Value["runtime"], "python36")
	assert.Equal(t, view.Data[0].Value["version"], "1")

	mConf = &models.ConfigurationView{
		Name:      "abc",
		Namespace: "default",
		Labels: map[string]string{
			"test": "test",
		},
		Data: []models.ConfigDataItem{
			{
				Key: "object",
				Value: map[string]string{
					"type":   "object",
					"source": "baidubos",
					"bucket": "baetyl",
					"object": "a.zip",
				},
			},
		},
	}

	res = &specV1.Configuration{
		Name:      "abc",
		Namespace: "default",
		Labels: map[string]string{
			"test": "test",
		},
		Data: map[string]string{
			common.ConfigObjectPrefix + "object": `{"metadata":{"bucket":"baetyl","object":"a.zip","source":"baidubos","type":"object"}}`,
		},
		CreationTimestamp: time.Now(),
		UpdateTimestamp:   time.Now(),
		Description:       "des",
		Version:           "12",
		System:            false,
	}
	mkConfigService.EXPECT().Get(mConf.Namespace, mConf.Name, gomock.Any()).Return(nil, nil).Times(1)
	mkConfigService.EXPECT().Create(mConf.Namespace, gomock.Any()).Return(res, nil).Times(1)

	w = httptest.NewRecorder()
	body, _ = json.Marshal(mConf)
	req, _ = http.NewRequest(http.MethodPost, "/v1/configs", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	err = json.Unmarshal(w.Body.Bytes(), &view)
	assert.NoError(t, err)
	assert.Equal(t, view.Name, res.Name)
	assert.Equal(t, view.Namespace, res.Namespace)
	assert.Equal(t, view.Labels, res.Labels)
	assert.Equal(t, view.Data[0].Key, "object")
	assert.Equal(t, view.Data[0].Value["type"], "object")
	assert.Equal(t, view.Data[0].Value["bucket"], "baetyl")
	assert.Equal(t, view.Data[0].Value["object"], "a.zip")
	assert.Equal(t, view.Data[0].Value["source"], "baidubos")

	mConf = &models.ConfigurationView{
		Name:      "abc",
		Namespace: "default",
		Labels: map[string]string{
			"test": "test",
		},
		Data: []models.ConfigDataItem{
			{
				Key: "object",
				Value: map[string]string{
					"type":   "object",
					"source": "baidubos",
					"bucket": "baetyl",
					"object": "a.zip",
				},
			},
		},
	}

	mkConfigService.EXPECT().Get(mConf.Namespace, mConf.Name, gomock.Any()).Return(nil, nil).Times(1)
	mkConfigService.EXPECT().Create(mConf.Namespace, gomock.Any()).Return(nil, errors.New("err")).Times(1)

	w = httptest.NewRecorder()
	body, _ = json.Marshal(mConf)
	req, _ = http.NewRequest(http.MethodPost, "/v1/configs", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpdateConfig(t *testing.T) {
	api, router, mockCtl := initConfigAPI(t)
	defer mockCtl.Finish()
	mkConfigService, mkAppService := ms.NewMockConfigService(mockCtl), ms.NewMockApplicationService(mockCtl)
	mkNodeService, mkIndexService := ms.NewMockNodeService(mockCtl), ms.NewMockIndexService(mockCtl)

	api.configService, api.nodeService = mkConfigService, mkNodeService
	api.applicationService, api.indexService = mkAppService, mkIndexService

	namespace, name := "default", "abc"
	mConf := &models.ConfigurationView{
		Name:      name,
		Namespace: namespace,
		Labels: map[string]string{
			"test": "test",
		},
		Data: []models.ConfigDataItem{
			{
				Key: "function",
				Value: map[string]string{
					"type":     "function",
					"function": "process",
					"version":  "1",
					"runtime":  "python36",
					"handler":  "index.handler",
					"bucket":   "baetyl",
					"object":   "a.zip",
				},
			},
			{
				Key: "function",
				Value: map[string]string{
					"type":     "function",
					"function": "process",
					"version":  "1",
					"runtime":  "python36",
					"handler":  "index.handler",
					"bucket":   "baetyl",
					"object":   "a.zip",
				},
			},
		},
	}

	mkConfigService.EXPECT().Get(namespace, name, gomock.Any()).Return(nil, errors.New("err")).Times(1)
	w := httptest.NewRecorder()
	body, _ := json.Marshal(mConf)
	req, _ := http.NewRequest(http.MethodPut, "/v1/configs/"+name, bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	res := &specV1.Configuration{
		Name:      name,
		Namespace: namespace,
		Labels: map[string]string{
			"test": "test",
		},
		Data: map[string]string{
			common.ConfigObjectPrefix + "function": `{"metadata":{"bucket":"baetyl","function":"process","handler":"index.handler","object":"a.zip","runtime":"python36","type":"function","userID":"default","version":"1"}}`,
		},
	}
	// 200: config is unchanged
	mkConfigService.EXPECT().Get(namespace, name, gomock.Any()).Return(res, nil).Times(1)
	w = httptest.NewRecorder()
	body, _ = json.Marshal(mConf)
	req, _ = http.NewRequest(http.MethodPut, "/v1/configs/"+name, bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	mConf2 := &models.ConfigurationView{
		Namespace: namespace,
		Data: []models.ConfigDataItem{
			{
				Key: "function",
				Value: map[string]string{
					"type":     "function",
					"function": "process",
					"version":  "1",
					"runtime":  "python36",
					"handler":  "index.handler",
					"bucket":   "baetyl",
					"object":   "a.zip",
				},
			},
		},
		Description: "update",
	}

	res2 := &specV1.Configuration{}
	mkConfigService.EXPECT().Get(namespace, name, gomock.Any()).Return(res2, nil).Times(1)
	mkConfigService.EXPECT().Update(namespace, gomock.Any()).Return(nil, errors.New("err")).Times(1)
	w = httptest.NewRecorder()
	body, _ = json.Marshal(mConf2)
	req, _ = http.NewRequest(http.MethodPut, "/v1/configs/abc", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	res3 := &specV1.Configuration{
		Name:      name,
		Namespace: namespace,
		Labels: map[string]string{
			"test": "test",
		},
		Data: map[string]string{
			common.ConfigObjectPrefix + "function": `{"metadata":{"bucket":"baetyl","function":"process","handler":"index.handler","object":"a.zip","runtime":"python36","type":"function","version":"1"}}`,
		},
		Description: "diff",
	}
	mkConfigService.EXPECT().Get(namespace, name, gomock.Any()).Return(res3, nil).Times(1)
	mkConfigService.EXPECT().Update(namespace, gomock.Any()).Return(res, nil).Times(1)
	mkIndexService.EXPECT().ListAppIndexByConfig(mConf2.Namespace, "abc").Return(nil, errors.New("err")).Times(1)
	w = httptest.NewRecorder()
	body, _ = json.Marshal(mConf2)
	req, _ = http.NewRequest(http.MethodPut, "/v1/configs/"+name, bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	appNames := make([]string, 0)
	mkConfigService.EXPECT().Get(namespace, name, gomock.Any()).Return(res3, nil).Times(1)
	mkConfigService.EXPECT().Update(namespace, gomock.Any()).Return(res, nil).Times(1)
	mkIndexService.EXPECT().ListAppIndexByConfig(namespace, name).Return(appNames, nil).Times(1)
	w = httptest.NewRecorder()
	body, _ = json.Marshal(mConf2)
	req, _ = http.NewRequest(http.MethodPut, "/v1/configs/"+name, bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	appNames = []string{"app01", "app02", "app03"}
	mkConfigService.EXPECT().Get(namespace, name, gomock.Any()).Return(res3, nil).Times(1)
	mkConfigService.EXPECT().Update(namespace, gomock.Any()).Return(res, nil).Times(1)
	mkIndexService.EXPECT().ListAppIndexByConfig(namespace, name).Return(appNames, nil).Times(1)
	mkAppService.EXPECT().Get(namespace, "app01", "").Return(nil, errors.New("err")).Times(1)
	w = httptest.NewRecorder()
	body, _ = json.Marshal(mConf2)
	req, _ = http.NewRequest(http.MethodPut, "/v1/configs/"+name, bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	apps := []*specV1.Application{
		{
			Namespace: "default",
			Name:      appNames[0],
			Volumes: []specV1.Volume{
				{
					Name:         "vol0",
					VolumeSource: specV1.VolumeSource{Config: &specV1.ObjectReference{Name: "cba", Version: "1"}},
				},
				{
					Name:         "vol1",
					VolumeSource: specV1.VolumeSource{Config: &specV1.ObjectReference{Name: "cba", Version: "2"}},
				},
				{
					Name:         "vol2",
					VolumeSource: specV1.VolumeSource{Config: &specV1.ObjectReference{Name: "cba", Version: "3"}},
				},
			},
		},
		{
			Namespace: "default",
			Name:      appNames[1],
			Volumes: []specV1.Volume{
				{
					Name:         "vol1",
					VolumeSource: specV1.VolumeSource{Config: &specV1.ObjectReference{Name: "cba", Version: "3"}},
				},
			},
		},
		{
			Namespace: "default",
			Name:      appNames[2],
			Volumes: []specV1.Volume{
				{
					Name:         "vol2",
					VolumeSource: specV1.VolumeSource{Config: &specV1.ObjectReference{Name: "cba", Version: "4"}},
				},
			},
		},
	}

	mkConfigService.EXPECT().Get(namespace, name, gomock.Any()).Return(res3, nil).Times(1)
	mkConfigService.EXPECT().Update(namespace, gomock.Any()).Return(res, nil).Times(1)
	mkIndexService.EXPECT().ListAppIndexByConfig(namespace, name).Return(appNames, nil).Times(1)
	mkAppService.EXPECT().Get(namespace, appNames[0], "").Return(apps[0], nil).Times(1)
	mkAppService.EXPECT().Get(namespace, appNames[1], "").Return(apps[1], nil).Times(1)
	mkAppService.EXPECT().Get(namespace, appNames[2], "").Return(apps[2], nil).Times(1)
	w = httptest.NewRecorder()
	body, _ = json.Marshal(mConf2)
	req, _ = http.NewRequest(http.MethodPut, "/v1/configs/"+name, bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteConfig(t *testing.T) {
	api, router, mockCtl := initConfigAPI(t)
	defer mockCtl.Finish()
	mkConfigService, mkIndexService := ms.NewMockConfigService(mockCtl), ms.NewMockIndexService(mockCtl)
	api.configService = mkConfigService
	api.indexService = mkIndexService

	mConf := &specV1.Configuration{
		Namespace: "default",
		Name:      "abc",
	}

	// 404
	mkConfigService.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
	req, _ := http.NewRequest(http.MethodDelete, "/v1/configs/abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// 500
	mkConfigService.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
	req, _ = http.NewRequest(http.MethodDelete, "/v1/configs/abc", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// 500
	mkConfigService.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(mConf, nil)
	mkIndexService.EXPECT().ListAppIndexByConfig(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
	req, _ = http.NewRequest(http.MethodDelete, "/v1/configs/abc", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	appNames := []string{"app01"}

	// 403
	mkConfigService.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(mConf, nil)
	mkIndexService.EXPECT().ListAppIndexByConfig(gomock.Any(), gomock.Any()).Return(appNames, nil)
	req, _ = http.NewRequest(http.MethodDelete, "/v1/configs/abc", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	// 200
	mkConfigService.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(mConf, nil)
	mkIndexService.EXPECT().ListAppIndexByConfig(gomock.Any(), gomock.Any()).Return(nil, nil)
	mkConfigService.EXPECT().Delete(mConf.Namespace, mConf.Name).Return(nil)
	req, _ = http.NewRequest(http.MethodDelete, "/v1/configs/abc", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetAppByConfig(t *testing.T) {
	api, router, mockCtl := initConfigAPI(t)
	defer mockCtl.Finish()
	mkConfigService, mkAppService := ms.NewMockConfigService(mockCtl), ms.NewMockApplicationService(mockCtl)
	mkNodeService, mkIndexService := ms.NewMockNodeService(mockCtl), ms.NewMockIndexService(mockCtl)

	api.configService, api.nodeService = mkConfigService, mkNodeService
	api.applicationService, api.indexService = mkAppService, mkIndexService

	mConf := &specV1.Configuration{
		Namespace: "default",
		Name:      "abc",
	}
	appNames := []string{"app01", "app02", "app03"}
	mkConfigService.EXPECT().Get(mConf.Namespace, mConf.Name, "").Return(mConf, nil)
	apps := []*specV1.Application{
		{
			Namespace: "default",
			Name:      appNames[0],
		},
		{
			Namespace: "default",
			Name:      appNames[1],
		},
		{
			Namespace: "default",
			Name:      appNames[2],
		},
	}

	mkIndexService.EXPECT().ListAppIndexByConfig(mConf.Namespace, mConf.Name).Return(appNames, nil).Times(1)
	mkAppService.EXPECT().Get(mConf.Namespace, appNames[0], "").Return(apps[0], nil)
	mkAppService.EXPECT().Get(mConf.Namespace, appNames[1], "").Return(apps[1], nil)
	mkAppService.EXPECT().Get(mConf.Namespace, appNames[2], "").Return(apps[2], nil)

	// 200
	req, _ := http.NewRequest(http.MethodGet, "/v1/configs/abc/apps", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 404
	mkConfigService.EXPECT().Get(mConf.Namespace, mConf.Name, "").Return(nil, nil)
	req, _ = http.NewRequest(http.MethodGet, "/v1/configs/abc/apps", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// 500
	mkConfigService.EXPECT().Get(mConf.Namespace, mConf.Name, "").Return(nil, fmt.Errorf("error"))
	req, _ = http.NewRequest(http.MethodGet, "/v1/configs/abc/apps", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
