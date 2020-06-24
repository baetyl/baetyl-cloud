package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/baetyl/baetyl-cloud/common"
	ms "github.com/baetyl/baetyl-cloud/mock/service"
	"github.com/baetyl/baetyl-cloud/models"
	specV1 "github.com/baetyl/baetyl-go/spec/v1"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// TODO: optimize this layer, general abstraction

func initRegistryAPI(t *testing.T) (*API, *gin.Engine, *gomock.Controller) {
	api := &API{}
	router := gin.Default()
	mockCtl := gomock.NewController(t)
	mockIM := func(c *gin.Context) { common.NewContext(c).SetNamespace("default") }
	v1 := router.Group("v1")
	{
		configs := v1.Group("/registries")
		configs.GET("/:name", mockIM, common.Wrapper(api.GetRegistry))
		configs.GET("/:name/apps", mockIM, common.Wrapper(api.GetAppByRegistry))
		configs.PUT("/:name", mockIM, common.Wrapper(api.UpdateRegistry))
		configs.POST(":name/refresh", mockIM, common.Wrapper(api.RefreshRegistryPassword))
		configs.DELETE("/:name", mockIM, common.Wrapper(api.DeleteRegistry))
		configs.POST("", mockIM, common.Wrapper(api.CreateRegistry))
		configs.GET("", mockIM, common.Wrapper(api.ListRegistry))
	}

	return api, router, mockCtl
}

func TestGetRegistry(t *testing.T) {
	api, router, mockCtl := initRegistryAPI(t)
	defer mockCtl.Finish()
	mkSecretService := ms.NewMockSecretService(mockCtl)
	api.secretService = mkSecretService

	mConf := &models.Registry{
		Namespace: "default",
		Name:      "abc",
	}

	mConf2 := &specV1.Secret{
		Namespace: "default",
		Name:      "abc",
		Labels: map[string]string{
			specV1.SecretLabel: specV1.SecretRegistry,
		},
	}

	mkSecretService.EXPECT().Get(mConf.Namespace, mConf.Name, "").Return(mConf2, nil)
	mkSecretService.EXPECT().Get(mConf.Namespace, "cba", "").Return(nil, fmt.Errorf("error"))

	// 200
	req, _ := http.NewRequest(http.MethodGet, "/v1/registries/abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 404
	req, _ = http.NewRequest(http.MethodGet, "/v1/registries/cba", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req)
	assert.Equal(t, http.StatusInternalServerError, w2.Code)
}

func TestListRegistry(t *testing.T) {
	api, router, mockCtl := initRegistryAPI(t)
	defer mockCtl.Finish()
	mkSecretService := ms.NewMockSecretService(mockCtl)
	api.secretService = mkSecretService

	mClist := &models.SecretList{
		Total: 0,
	}

	mkSecretService.EXPECT().List("default", gomock.Any()).Return(mClist, nil)

	// 200
	req, _ := http.NewRequest(http.MethodGet, "/v1/registries", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateRegistry(t *testing.T) {
	api, router, mockCtl := initRegistryAPI(t)
	defer mockCtl.Finish()
	mkSecretService := ms.NewMockSecretService(mockCtl)
	api.secretService = mkSecretService

	mConf := &models.Registry{
		Namespace: "default",
		Name:      "abc",
		Username:  "username",
		Password:  "password",
		Address:   "address",
	}
	mConf2 := &specV1.Secret{
		Namespace: "default",
		Name:      "abc",
		Labels: map[string]string{
			specV1.SecretLabel: specV1.SecretRegistry,
		},
		Data: map[string][]byte{
			"address":  []byte("address"),
			"password": []byte("password"),
			"username": []byte("username"),
		},
	}
	mkSecretService.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
	mkSecretService.EXPECT().Create(mConf.Namespace, gomock.Any()).Return(mConf2, nil)
	w := httptest.NewRecorder()
	body, _ := json.Marshal(mConf)
	req, _ := http.NewRequest(http.MethodPost, "/v1/registries", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	mkSecretService.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(mConf2, nil)
	w2 := httptest.NewRecorder()
	body2, _ := json.Marshal(mConf)
	req2, _ := http.NewRequest(http.MethodPost, "/v1/registries", bytes.NewReader(body2))
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusBadRequest, w2.Code)
}

func TestUpdateRegistry(t *testing.T) {
	api, router, mockCtl := initRegistryAPI(t)
	defer mockCtl.Finish()
	mkSecretService := ms.NewMockSecretService(mockCtl)
	api.secretService = mkSecretService

	mConf := &models.Registry{
		Namespace:   "default",
		Name:        "abc",
		Description: "haha",
	}
	mConfSecret := &specV1.Secret{
		Namespace:   "default",
		Name:        "abc",
		Description: "haha",
		Labels: map[string]string{
			specV1.SecretLabel: specV1.SecretRegistry,
		},
		Data: map[string][]byte{
			"address":  []byte("address"),
			"password": []byte("password"),
			"username": []byte("username"),
		},
	}
	mkSecretService.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(mConfSecret, nil)
	w := httptest.NewRecorder()
	body, _ := json.Marshal(mConf)
	req, _ := http.NewRequest(http.MethodPut, "/v1/registries/abc", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	mkSecretService.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
	w2 := httptest.NewRecorder()
	body2, _ := json.Marshal(mConf)
	req2, _ := http.NewRequest(http.MethodPut, "/v1/registries/abc", bytes.NewReader(body2))
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusInternalServerError, w2.Code)

	mConfSecret2 := &specV1.Secret{
		Namespace:   "default",
		Name:        "cba",
		Description: "haha modify",
	}
	mkSecretService.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(mConfSecret, nil)
	mkSecretService.EXPECT().Update(mConfSecret2.Namespace, gomock.Any()).Return(nil, common.Error(common.ErrRequestParamInvalid))
	w3 := httptest.NewRecorder()
	body3, _ := json.Marshal(mConfSecret2)
	req3, _ := http.NewRequest(http.MethodPut, "/v1/registries/cba", bytes.NewReader(body3))
	router.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusBadRequest, w3.Code)
}

func TestRefreshRegistryPassword(t *testing.T) {
	api, router, mockCtl := initRegistryAPI(t)
	defer mockCtl.Finish()
	mkSecretService, mkAppService := ms.NewMockSecretService(mockCtl), ms.NewMockApplicationService(mockCtl)
	mkNodeService, mkIndexService := ms.NewMockNodeService(mockCtl), ms.NewMockIndexService(mockCtl)

	api.secretService, api.nodeService = mkSecretService, mkNodeService
	api.applicationService, api.indexService = mkAppService, mkIndexService

	mConf := &models.Registry{
		Namespace: "default",
		Name:      "abc",
		Password:  "haha",
	}

	mkSecretService.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
	w2 := httptest.NewRecorder()
	body2, _ := json.Marshal(mConf)
	req2, _ := http.NewRequest(http.MethodPut, "/v1/registries/abc", bytes.NewReader(body2))
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusInternalServerError, w2.Code)

	appNames := []string{"app1", "app2", "app3"}
	apps := []*specV1.Application{
		{
			Namespace: "default",
			Name:      appNames[0],
			Selector:  "tag=abc",
		},
		{
			Namespace: "default",
			Name:      appNames[1],
			Selector:  "tag=abc",
		},
		{
			Namespace: "default",
			Name:      appNames[2],
			Selector:  "tag=abc",
		},
	}

	mConfSecret := &specV1.Secret{
		Namespace:   "default",
		Name:        "abc",
		Description: "haha",
		Version:     "5",
		Labels: map[string]string{
			specV1.SecretLabel: specV1.SecretRegistry,
		},
		Data: map[string][]byte{
			"address":  []byte("address"),
			"password": []byte("password"),
			"username": []byte("username"),
		},
	}

	mConf2 := &models.Registry{
		Namespace: "default",
		Name:      "abc",
		Password:  "haha",
	}
	mkSecretService.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(mConfSecret, nil)
	mkSecretService.EXPECT().Update(mConf2.Namespace, gomock.Any()).Return(mConfSecret, nil)
	mkIndexService.EXPECT().ListAppIndexBySecret(mConf2.Namespace, mConf2.Name).Return(appNames, nil).Times(1)
	mkAppService.EXPECT().Get(mConf2.Namespace, appNames[0], "").Return(apps[0], nil).AnyTimes()
	mkAppService.EXPECT().Get(mConf2.Namespace, appNames[1], "").Return(apps[1], nil).AnyTimes()
	mkAppService.EXPECT().Get(mConf2.Namespace, appNames[2], "").Return(apps[2], nil).AnyTimes()
	mkAppService.EXPECT().Update(mConf2.Namespace, gomock.Any()).Return(apps[0], nil).AnyTimes()
	mkNodeService.EXPECT().UpdateNodeAppVersion(mConf2.Namespace, gomock.Any()).Return(nil, nil).AnyTimes()
	w3 := httptest.NewRecorder()
	body3, _ := json.Marshal(mConf2)
	req3, _ := http.NewRequest(http.MethodPost, "/v1/registries/cba/refresh", bytes.NewReader(body3))
	router.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusOK, w3.Code)
}

func TestDeleteRegistry(t *testing.T) {
	api, router, mockCtl := initRegistryAPI(t)
	defer mockCtl.Finish()
	mkSecretService := ms.NewMockSecretService(mockCtl)
	mkIndexService := ms.NewMockIndexService(mockCtl)
	api.secretService = mkSecretService
	api.indexService = mkIndexService
	mConfSecret := &specV1.Secret{
		Namespace:   "default",
		Name:        "abc",
		Description: "haha",
		Labels: map[string]string{
			specV1.SecretLabel: specV1.SecretRegistry,
		},
	}
	mkSecretService.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(mConfSecret, nil)
	mkSecretService.EXPECT().Delete(mConfSecret.Namespace, mConfSecret.Name).Return(nil)
	mkIndexService.EXPECT().ListAppIndexBySecret(gomock.Any(), gomock.Any()).Return(nil, nil)
	// 200
	req, _ := http.NewRequest(http.MethodDelete, "/v1/registries/abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetAppByRegistry(t *testing.T) {
	api, router, mockCtl := initRegistryAPI(t)
	defer mockCtl.Finish()
	mkSecretService, mkAppService := ms.NewMockSecretService(mockCtl), ms.NewMockApplicationService(mockCtl)
	mkNodeService, mkIndexService := ms.NewMockNodeService(mockCtl), ms.NewMockIndexService(mockCtl)

	api.secretService, api.nodeService = mkSecretService, mkNodeService
	api.applicationService, api.indexService = mkAppService, mkIndexService

	appNames := []string{"app1", "app2", "app3"}
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

	mConfSecret3 := &specV1.Secret{
		Namespace:   "default",
		Name:        "abc",
		Description: "haha",
		Version:     "5",
		Labels: map[string]string{
			specV1.SecretLabel: specV1.SecretRegistry,
		},
	}

	mkSecretService.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(mConfSecret3, nil)

	mkIndexService.EXPECT().ListAppIndexBySecret(mConfSecret3.Namespace, mConfSecret3.Name).Return(appNames, nil).Times(1)
	mkAppService.EXPECT().Get(mConfSecret3.Namespace, appNames[0], "").Return(apps[0], nil).AnyTimes()
	mkAppService.EXPECT().Get(mConfSecret3.Namespace, appNames[1], "").Return(apps[1], nil).AnyTimes()
	mkAppService.EXPECT().Get(mConfSecret3.Namespace, appNames[2], "").Return(apps[2], nil).AnyTimes()

	w4 := httptest.NewRecorder()
	req4, _ := http.NewRequest(http.MethodGet, "/v1/registries/abc/apps", nil)
	router.ServeHTTP(w4, req4)
	assert.Equal(t, http.StatusOK, w4.Code)
}
