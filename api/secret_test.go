package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/common"
	ms "github.com/baetyl/baetyl-cloud/v2/mock/service"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/service"
)

// TODO: optimize this layer, general abstraction

func initSecretAPI(t *testing.T) (*API, *gin.Engine, *gomock.Controller) {
	api := &API{}
	router := gin.Default()
	mockCtl := gomock.NewController(t)
	mockIM := func(c *gin.Context) { common.NewContext(c).SetNamespace("default") }
	v1 := router.Group("v1")
	{
		configs := v1.Group("/secrets")
		configs.GET("/:name", mockIM, common.Wrapper(api.GetSecret))
		configs.GET("/:name/apps", mockIM, common.Wrapper(api.GetAppBySecret))
		configs.PUT("/:name", mockIM, common.Wrapper(api.UpdateSecret))
		configs.DELETE("/:name", mockIM, common.Wrapper(api.DeleteSecret))
		configs.POST("", mockIM, common.Wrapper(api.CreateSecret))
		configs.GET("", mockIM, common.Wrapper(api.ListSecret))
	}

	return api, router, mockCtl
}

func TestGetSecret(t *testing.T) {
	api, router, mockCtl := initSecretAPI(t)
	defer mockCtl.Finish()

	sSecret := ms.NewMockSecretService(mockCtl)
	api.AppCombinedService = &service.AppCombinedService{
		Secret: sSecret,
	}

	mConf := &specV1.Secret{
		Namespace: "default",
		Name:      "abc",
	}

	mConf2 := &specV1.Secret{
		Namespace: "default",
		Name:      "abc",
		Labels: map[string]string{
			specV1.SecretLabel: specV1.SecretConfig,
		},
	}

	sSecret.EXPECT().Get(mConf.Namespace, mConf.Name, "").Return(mConf2, nil)
	sSecret.EXPECT().Get(mConf.Namespace, "cba", "").Return(nil, fmt.Errorf("error"))

	// 200
	req, _ := http.NewRequest(http.MethodGet, "/v1/secrets/abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 404
	req, _ = http.NewRequest(http.MethodGet, "/v1/secrets/cba", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req)
	assert.Equal(t, http.StatusNotFound, w2.Code)
}

func TestListSecret(t *testing.T) {
	api, router, mockCtl := initSecretAPI(t)
	defer mockCtl.Finish()

	sSecret := ms.NewMockSecretService(mockCtl)
	api.AppCombinedService = &service.AppCombinedService{
		Secret: sSecret,
	}

	mClist := &models.SecretList{
		Total: 0,
	}

	sSecret.EXPECT().List("default", gomock.Any()).Return(mClist, nil)

	// 200
	req, _ := http.NewRequest(http.MethodGet, "/v1/secrets", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateSecret(t *testing.T) {
	api, router, mockCtl := initSecretAPI(t)
	defer mockCtl.Finish()

	sSecret := ms.NewMockSecretService(mockCtl)
	api.AppCombinedService = &service.AppCombinedService{
		Secret: sSecret,
	}

	mConf := &models.SecretView{
		Namespace: "default",
		Name:      "abc",
		Data: map[string]string{
			"a": "b",
		},
	}
	mConf2 := &specV1.Secret{
		Namespace: "default",
		Name:      "abc",
		Labels: map[string]string{
			specV1.SecretLabel: specV1.SecretConfig,
		},
	}
	sSecret.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
	sSecret.EXPECT().Create(mConf.Namespace, gomock.Any()).Return(mConf2, nil)
	w := httptest.NewRecorder()
	body, _ := json.Marshal(mConf)
	req, _ := http.NewRequest(http.MethodPost, "/v1/secrets", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	sSecret.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(mConf2, nil)
	w2 := httptest.NewRecorder()
	body2, _ := json.Marshal(mConf)
	req2, _ := http.NewRequest(http.MethodPost, "/v1/secrets", bytes.NewReader(body2))
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusBadRequest, w2.Code)
}

func TestUpdateSecret(t *testing.T) {
	api, router, mockCtl := initSecretAPI(t)
	defer mockCtl.Finish()

	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sSecret := ms.NewMockSecretService(mockCtl)
	api.AppCombinedService = &service.AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}

	sNode, sIndex := ms.NewMockNodeService(mockCtl), ms.NewMockIndexService(mockCtl)
	api.Node, api.Index = sNode, sIndex

	mConf := &models.SecretView{
		Namespace:   "default",
		Name:        "abc",
		Description: "haha",
		Data: map[string]string{
			"a": "b",
		},
	}
	mConfSecret := &specV1.Secret{
		Namespace:   "default",
		Name:        "abc",
		Description: "haha",
		Labels: map[string]string{
			specV1.SecretLabel: specV1.SecretConfig,
		},
		Data: map[string][]byte{
			"a": []byte("b"),
		},
	}
	sSecret.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(mConfSecret, nil)
	w := httptest.NewRecorder()
	body, _ := json.Marshal(mConf)
	req, _ := http.NewRequest(http.MethodPut, "/v1/secrets/abc", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	sSecret.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
	w2 := httptest.NewRecorder()
	body2, _ := json.Marshal(mConf)
	req2, _ := http.NewRequest(http.MethodPut, "/v1/secrets/abc", bytes.NewReader(body2))
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusInternalServerError, w2.Code)

	mConfSecret2 := &specV1.Secret{
		Namespace:   "default",
		Name:        "cba",
		Description: "haha modify",
		Labels: map[string]string{
			specV1.SecretLabel: specV1.SecretConfig,
		},
	}
	sSecret.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(mConfSecret2, nil)
	sSecret.EXPECT().Update(mConfSecret2.Namespace, gomock.Any()).Return(nil, common.Error(common.ErrRequestParamInvalid))
	w3 := httptest.NewRecorder()
	body3, _ := json.Marshal(mConf)
	req3, _ := http.NewRequest(http.MethodPut, "/v1/secrets/cba", bytes.NewReader(body3))
	router.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusBadRequest, w3.Code)

	appNames := []string{"app1", "app2", "app3"}
	apps := []*specV1.Application{
		{
			Namespace: "default",
			Name:      appNames[0],
			Selector:  "tag=abc",
			Volumes: []specV1.Volume{
				{
					Name:         "vol0",
					VolumeSource: specV1.VolumeSource{Secret: &specV1.ObjectReference{Name: "abc", Version: "1"}},
				},
				{
					Name:         "vol1",
					VolumeSource: specV1.VolumeSource{Secret: &specV1.ObjectReference{Name: "cba", Version: "2"}},
				},
				{
					Name:         "vol2",
					VolumeSource: specV1.VolumeSource{Secret: &specV1.ObjectReference{Name: "cba", Version: "3"}},
				},
			},
		},
		{
			Namespace: "default",
			Name:      appNames[1],
			Selector:  "tag=abc",
			Volumes: []specV1.Volume{
				{
					Name:         "vol1",
					VolumeSource: specV1.VolumeSource{Secret: &specV1.ObjectReference{Name: "cba", Version: "3"}},
				},
			},
		},
		{
			Namespace: "default",
			Name:      appNames[2],
			Selector:  "tag=abc",
			Volumes: []specV1.Volume{
				{
					Name:         "vol2",
					VolumeSource: specV1.VolumeSource{Secret: &specV1.ObjectReference{Name: "abc", Version: "4"}},
				},
			},
		},
	}

	mConfSecret3 := &specV1.Secret{
		Namespace:   "default",
		Name:        "abc",
		Description: "haha",
		Version:     "5",
		Labels: map[string]string{
			specV1.SecretLabel: specV1.SecretConfig,
		},
	}

	mConf2 := &models.Registry{
		Namespace: "default",
		Name:      "abc",
		Password:  "haha",
	}
	sSecret.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(mConfSecret3, nil)
	sSecret.EXPECT().Update(mConf2.Namespace, gomock.Any()).Return(mConfSecret3, nil)
	sIndex.EXPECT().ListAppIndexBySecret(mConf2.Namespace, mConf2.Name).Return(appNames, nil).Times(1)
	sApp.EXPECT().Get(mConf2.Namespace, appNames[0], "").Return(apps[0], nil).AnyTimes()
	sApp.EXPECT().Get(mConf2.Namespace, appNames[1], "").Return(apps[1], nil).AnyTimes()
	sApp.EXPECT().Get(mConf2.Namespace, appNames[2], "").Return(apps[2], nil).AnyTimes()
	sApp.EXPECT().Update(mConf2.Namespace, gomock.Any()).Return(apps[0], nil).AnyTimes()
	sNode.EXPECT().UpdateNodeAppVersion(mConf2.Namespace, gomock.Any()).Return(nil, nil).AnyTimes()

	w4 := httptest.NewRecorder()
	body4, _ := json.Marshal(mConf)
	req4, _ := http.NewRequest(http.MethodPut, "/v1/secrets/abc", bytes.NewReader(body4))
	router.ServeHTTP(w4, req4)
	assert.Equal(t, http.StatusOK, w4.Code)
}

func TestDeleteSecret(t *testing.T) {
	api, router, mockCtl := initSecretAPI(t)
	defer mockCtl.Finish()

	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sSecret := ms.NewMockSecretService(mockCtl)
	api.AppCombinedService = &service.AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}

	sNode, sIndex := ms.NewMockNodeService(mockCtl), ms.NewMockIndexService(mockCtl)
	api.Node, api.Index = sNode, sIndex

	mConfSecret := &specV1.Secret{
		Namespace:   "default",
		Name:        "abc",
		Description: "haha",
		Labels: map[string]string{
			specV1.SecretLabel: specV1.SecretConfig,
		},
	}
	sSecret.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(mConfSecret, nil).AnyTimes()
	sSecret.EXPECT().Delete(mConfSecret.Namespace, mConfSecret.Name).Return(nil).AnyTimes()
	sIndex.EXPECT().ListAppIndexBySecret(gomock.Any(), gomock.Any()).Return(nil, nil)
	// 200
	req, _ := http.NewRequest(http.MethodDelete, "/v1/secrets/abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	sIndex.EXPECT().ListAppIndexBySecret(gomock.Any(), gomock.Any()).Return([]string{"app1", "app2"}, nil)
	// 400
	req2, _ := http.NewRequest(http.MethodDelete, "/v1/secrets/abc", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusForbidden, w2.Code)
}

func TestGetAppBySecret(t *testing.T) {
	api, router, mockCtl := initSecretAPI(t)
	defer mockCtl.Finish()

	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sSecret := ms.NewMockSecretService(mockCtl)
	api.AppCombinedService = &service.AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}

	sNode, sIndex := ms.NewMockNodeService(mockCtl), ms.NewMockIndexService(mockCtl)
	api.Node, api.Index = sNode, sIndex

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
			specV1.SecretLabel: specV1.SecretConfig,
		},
	}

	sSecret.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(mConfSecret3, nil)

	sIndex.EXPECT().ListAppIndexBySecret(mConfSecret3.Namespace, mConfSecret3.Name).Return(appNames, nil).Times(1)
	sApp.EXPECT().Get(mConfSecret3.Namespace, appNames[0], "").Return(apps[0], nil).AnyTimes()
	sApp.EXPECT().Get(mConfSecret3.Namespace, appNames[1], "").Return(apps[1], nil).AnyTimes()
	sApp.EXPECT().Get(mConfSecret3.Namespace, appNames[2], "").Return(apps[2], nil).AnyTimes()

	w4 := httptest.NewRecorder()
	req4, _ := http.NewRequest(http.MethodGet, "/v1/secrets/abc/apps", nil)
	router.ServeHTTP(w4, req4)
	assert.Equal(t, http.StatusOK, w4.Code)
}
