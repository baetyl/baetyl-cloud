package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/baetyl/baetyl-cloud/common"
	ms "github.com/baetyl/baetyl-cloud/mock/service"
	"github.com/baetyl/baetyl-cloud/models"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func getMockContainerApp() *specV1.Application {
	mApp := &specV1.Application{
		Namespace: "baetyl-cloud",
		Name:      "abc",
		Type:      common.ContainerApp,
		Services: []specV1.Service{
			{
				Name:     "Agent",
				Hostname: "test-agent",
				Image:    "hub.baidubce.com/baetyl/baetyl-agent:1.0.0",
				Replica:  1,
				VolumeMounts: []specV1.VolumeMount{
					{
						Name:      "name",
						MountPath: "mountPath",
					},
				},
				Ports: []specV1.ContainerPort{
					{
						HostPort:      8080,
						ContainerPort: 8080,
					},
				},
				Devices: []specV1.Device{
					{
						DevicePath: "DevicePath",
					},
				},
				Args: []string{"test"},
				Env: []specV1.Environment{
					{
						Name:  "name",
						Value: "value",
					},
				},
			},
		},
		Volumes: []specV1.Volume{
			{
				Name: "name",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name: "agent-conf",
					},
				},
			},
			{
				Name: "secret",
				VolumeSource: specV1.VolumeSource{
					Secret: &specV1.ObjectReference{
						Name: "secret01",
					},
				},
			},
		},
	}
	return mApp
}

func getMockFunctionApp() *specV1.Application {
	mApp := &specV1.Application{
		Namespace: "baetyl-cloud",
		Name:      "abc",
		Type:      common.FunctionApp,
		Services: []specV1.Service{
			{
				Name:     "Agent",
				Hostname: "test-agent",
				Replica:  1,
				VolumeMounts: []specV1.VolumeMount{
					{
						Name:      "baetyl-function-code-Agent",
						MountPath: "mountPath",
					},
					{
						Name:      "baetyl-function-config-Agent",
						MountPath: "mountPath",
					},
				},
				Devices: []specV1.Device{
					{
						DevicePath: "DevicePath",
					},
				},
				FunctionConfig: &specV1.ServiceFunctionConfig{
					Name:    "func1",
					Runtime: "python36",
				},
				Functions: []specV1.ServiceFunction{
					{
						Name:    "process",
						Handler: "index.handler",
						CodeDir: "path",
					},
				},
			},
		},
		Volumes: []specV1.Volume{
			{
				Name: "baetyl-function-code-Agent",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name: "func1",
					},
				},
			},
			{
				Name: "baetyl-function-config-Agent",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name: "baetyl-function-app-service-xxxxxxxxx",
					},
				},
			},
		},
	}
	return mApp
}

func initApplicationAPI(t *testing.T) (*API, *gin.Engine, *gomock.Controller) {
	api := &API{}
	router := gin.Default()
	mockCtl := gomock.NewController(t)
	mockIM := func(c *gin.Context) { c.Set(common.KeyContextNamespace, "baetyl-cloud") }
	v1 := router.Group("v1")
	{
		configs := v1.Group("/apps")
		configs.GET("/:name", mockIM, common.Wrapper(api.GetApplication))
		configs.PUT("/:name", mockIM, common.Wrapper(api.UpdateApplication))
		configs.DELETE("/:name", mockIM, common.Wrapper(api.DeleteApplication))
		configs.POST("", mockIM, common.Wrapper(api.CreateApplication))
		configs.GET("", mockIM, common.Wrapper(api.ListApplication))
	}
	return api, router, mockCtl
}

func TestGetContainerApplication(t *testing.T) {
	api, router, mockCtl := initApplicationAPI(t)
	defer mockCtl.Finish()
	mkApplicationService := ms.NewMockApplicationService(mockCtl)
	mSecretService := ms.NewMockSecretService(mockCtl)
	api.applicationService = mkApplicationService
	api.secretService = mSecretService

	mApp := getMockContainerApp()
	mkApplicationService.EXPECT().Get(mApp.Namespace, "cba", "").Return(nil, errors.New("err")).Times(1)
	req, _ := http.NewRequest(http.MethodGet, "/v1/apps/cba", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	secret := &specV1.Secret{
		Name:    "secret01",
		Version: "1",
		Labels: map[string]string{
			specV1.SecretLabel: specV1.SecretRegistry,
		},
	}
	mkApplicationService.EXPECT().Get(mApp.Namespace, mApp.Name, "").Return(mApp, nil).Times(1)
	mSecretService.EXPECT().Get(mApp.Namespace, secret.Name, "").Return(secret, nil).Times(1)

	// 200
	req, _ = http.NewRequest(http.MethodGet, "/v1/apps/abc", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var view models.ApplicationView
	err := json.Unmarshal(w.Body.Bytes(), &view)
	assert.NoError(t, err)
	assert.Equal(t, view.Registries[0].Name, "secret01")
}

func TestGetFunctionApplication(t *testing.T) {
	api, router, mockCtl := initApplicationAPI(t)
	defer mockCtl.Finish()
	mkApplicationService := ms.NewMockApplicationService(mockCtl)
	mkConfigService := ms.NewMockConfigService(mockCtl)
	api.applicationService = mkApplicationService
	api.configService = mkConfigService

	mApp := getMockFunctionApp()

	mkApplicationService.EXPECT().Get(mApp.Namespace, "cba", "").Return(nil, errors.New("err")).Times(1)
	req, _ := http.NewRequest(http.MethodGet, "/v1/apps/cba", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	functions := models.ServiceFunction{
		Functions: []specV1.ServiceFunction{
			{
				Name:    "process",
				Handler: "index.handler",
				CodeDir: "path",
			},
		},
	}
	data, err := json.Marshal(&functions)
	assert.NoError(t, err)
	config := &specV1.Configuration{
		Data: map[string]string{
			"service.yml": string(data),
		},
	}
	mkApplicationService.EXPECT().Get(mApp.Namespace, mApp.Name, "").Return(mApp, nil).Times(1)
	mkConfigService.EXPECT().Get(mApp.Namespace, "baetyl-function-app-service-xxxxxxxxx", "").Return(config, nil).Times(1)

	// 200
	req, _ = http.NewRequest(http.MethodGet, "/v1/apps/abc", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var view models.ApplicationView
	err = json.Unmarshal(w.Body.Bytes(), &view)
	assert.NoError(t, err)
	assert.Equal(t, view.Services[0].FunctionConfig.Name, "func1")
	assert.Equal(t, view.Services[0].FunctionConfig.Runtime, "python36")
	assert.Equal(t, view.Services[0].Functions[0], functions.Functions[0])
}

func TestListApplication(t *testing.T) {
	api, router, mockCtl := initApplicationAPI(t)
	defer mockCtl.Finish()
	mkApplicationService := ms.NewMockApplicationService(mockCtl)
	api.applicationService = mkApplicationService

	mClist := &models.ApplicationList{}

	mkApplicationService.EXPECT().List("baetyl-cloud", &models.ListOptions{
		LabelSelector: "!" + common.LabelSystem,
	}).Return(mClist, nil)

	// 200
	req, _ := http.NewRequest(http.MethodGet, "/v1/apps", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	mkApplicationService.EXPECT().List("baetyl-cloud", &models.ListOptions{
		LabelSelector: "!" + common.LabelSystem,
	}).Return(nil, fmt.Errorf("error"))

	// 400
	req, _ = http.NewRequest(http.MethodGet, "/v1/apps", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mkApplicationService.EXPECT().List("baetyl-cloud", &models.ListOptions{
		LabelSelector: "!" + common.LabelSystem,
	}).Return(nil, fmt.Errorf("error"))

	// 400
	req, _ = http.NewRequest(http.MethodGet, "/v1/apps", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateContainerApplication(t *testing.T) {
	api, router, mockCtl := initApplicationAPI(t)
	defer mockCtl.Finish()

	mkApplicationService := ms.NewMockApplicationService(mockCtl)
	mkIndexService := ms.NewMockIndexService(mockCtl)
	mSecretService := ms.NewMockSecretService(mockCtl)
	mkConfigService := ms.NewMockConfigService(mockCtl)
	mkNodeService := ms.NewMockNodeService(mockCtl)

	api.applicationService = mkApplicationService
	api.indexService = mkIndexService
	api.secretService = mSecretService
	api.configService = mkConfigService
	api.nodeService = mkNodeService

	appView := &models.ApplicationView{
		Application: specV1.Application{
			Namespace: "baetyl-cloud",
			Name:      "abc",
			Type:      common.ContainerApp,
			Services: []specV1.Service{
				{
					Name:     "Agent",
					Hostname: "test-agent",
					Image:    "hub.baidubce.com/baetyl/baetyl-agent:1.0.0",
					Replica:  1,
					VolumeMounts: []specV1.VolumeMount{
						{
							Name:      "name",
							MountPath: "mountPath",
						},
					},
					Ports: []specV1.ContainerPort{
						{
							HostPort:      8080,
							ContainerPort: 8080,
						},
					},
					Devices: []specV1.Device{
						{
							DevicePath: "DevicePath",
						},
					},
					Args: []string{"test"},
					Env: []specV1.Environment{
						{
							Name:  "name",
							Value: "value",
						},
					},
				},
			},
			Volumes: []specV1.Volume{
				{
					Name: "name",
					VolumeSource: specV1.VolumeSource{
						Config: &specV1.ObjectReference{
							Name: "agent-conf",
						},
					},
				},
				{
					Name: "secret",
					VolumeSource: specV1.VolumeSource{
						Secret: &specV1.ObjectReference{
							Name: "secret01",
						},
					},
				},
			},
		},
		Registries: []models.RegistryView{
			{
				Name: "registry01",
			},
		},
	}

	config := &specV1.Configuration{}
	secret := &specV1.Secret{}
	app := &specV1.Application{}
	mkConfigService.EXPECT().Get(appView.Namespace, "agent-conf", "").Return(config, nil).Times(1)
	mSecretService.EXPECT().Get(appView.Namespace, "secret01", "").Return(secret, nil).Times(1)
	mSecretService.EXPECT().Get(appView.Namespace, "registry01", "").Return(secret, nil).Times(1)
	mkApplicationService.EXPECT().Get(appView.Namespace, "abc", "").Return(app, nil).Times(1)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(appView)
	req, _ := http.NewRequest(http.MethodPost, "/v1/apps?base=eden", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	mkConfigService.EXPECT().Get(appView.Namespace, "agent-conf", "").Return(nil, fmt.Errorf("config not found")).Times(1)

	w = httptest.NewRecorder()
	body, _ = json.Marshal(appView)
	req, _ = http.NewRequest(http.MethodPost, "/v1/apps?base=eden", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mkConfigService.EXPECT().Get(appView.Namespace, "agent-conf", "").Return(config, nil).Times(1)
	mSecretService.EXPECT().Get(appView.Namespace, "secret01", "").Return(nil, fmt.Errorf("secret not found")).Times(1)

	w = httptest.NewRecorder()
	body, _ = json.Marshal(appView)
	req, _ = http.NewRequest(http.MethodPost, "/v1/apps?base=eden", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mkConfigService.EXPECT().Get(appView.Namespace, "agent-conf", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	w = httptest.NewRecorder()
	body, _ = json.Marshal(appView)
	req, _ = http.NewRequest(http.MethodPost, "/v1/apps?base=eden", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	mkConfigService.EXPECT().Get(appView.Namespace, "agent-conf", "").Return(config, nil).Times(1)
	mSecretService.EXPECT().Get(appView.Namespace, "secret01", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)

	w = httptest.NewRecorder()
	body, _ = json.Marshal(appView)
	req, _ = http.NewRequest(http.MethodPost, "/v1/apps?base=eden", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	mkConfigService.EXPECT().Get(appView.Namespace, "agent-conf", "").Return(config, nil).Times(1)
	mSecretService.EXPECT().Get(appView.Namespace, "secret01", "").Return(secret, nil).Times(1)
	mSecretService.EXPECT().Get(appView.Namespace, "registry01", "").Return(secret, nil).Times(1)
	mkApplicationService.EXPECT().Get(appView.Namespace, "abc", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	mkApplicationService.EXPECT().Get(appView.Namespace, "eden", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	w = httptest.NewRecorder()
	body, _ = json.Marshal(appView)
	req, _ = http.NewRequest(http.MethodPost, "/v1/apps?base=eden", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	eden2 := &specV1.Application{
		Namespace: appView.Namespace,
		Name:      "abc",
		Type:      common.FunctionApp,
		Services:  []specV1.Service{},
		Volumes: []specV1.Volume{
			{
				Name: "testSecret01",
				VolumeSource: specV1.VolumeSource{
					Secret: &specV1.ObjectReference{
						Name:    "secret01",
						Version: "1",
					},
				},
			},
		},
	}

	mkConfigService.EXPECT().Get(appView.Namespace, "agent-conf", "").Return(config, nil).Times(1)
	mSecretService.EXPECT().Get(appView.Namespace, "secret01", "").Return(secret, nil).Times(1)
	mSecretService.EXPECT().Get(appView.Namespace, "registry01", "").Return(secret, nil).Times(1)
	mkApplicationService.EXPECT().Get(appView.Namespace, "abc", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	mkApplicationService.EXPECT().Get(appView.Namespace, "eden2", "").Return(eden2, nil).Times(1)
	w = httptest.NewRecorder()
	body, _ = json.Marshal(appView)
	req, _ = http.NewRequest(http.MethodPost, "/v1/apps?base=eden2", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	eden2 = &specV1.Application{
		Namespace: appView.Namespace,
		Name:      "abc",
		Type:      common.ContainerApp,
		Services:  []specV1.Service{},
		Volumes: []specV1.Volume{
			{
				Name: "testSecret01",
				VolumeSource: specV1.VolumeSource{
					Secret: &specV1.ObjectReference{
						Name:    "secret01",
						Version: "1",
					},
				},
			},
		},
	}

	mkConfigService.EXPECT().Get(appView.Namespace, "agent-conf", "").Return(config, nil).Times(1)
	mSecretService.EXPECT().Get(appView.Namespace, "secret01", "").Return(secret, nil).Times(1)
	mSecretService.EXPECT().Get(appView.Namespace, "registry01", "").Return(secret, nil).Times(1)
	mkApplicationService.EXPECT().Get(appView.Namespace, "abc", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	mkApplicationService.EXPECT().Get(appView.Namespace, "eden2", "").Return(eden2, nil).Times(1)
	mkApplicationService.EXPECT().CreateWithBase(appView.Namespace, gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error")).Times(1)
	w = httptest.NewRecorder()
	body, _ = json.Marshal(appView)
	req, _ = http.NewRequest(http.MethodPost, "/v1/apps?base=eden2", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mApp := getMockContainerApp()
	mkConfigService.EXPECT().Get(appView.Namespace, "agent-conf", "").Return(config, nil).Times(1)
	mSecretService.EXPECT().Get(appView.Namespace, "secret01", "").Return(secret, nil).Times(1)
	mSecretService.EXPECT().Get(appView.Namespace, "registry01", "").Return(secret, nil).Times(1)
	mkApplicationService.EXPECT().Get(appView.Namespace, "abc", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	mkApplicationService.EXPECT().Get(appView.Namespace, "eden2", "").Return(eden2, nil).Return(eden2, nil).Times(1)
	mkApplicationService.EXPECT().CreateWithBase(appView.Namespace, gomock.Any(), eden2).Return(mApp, nil).Times(1)
	mkNodeService.EXPECT().UpdateNodeAppVersion(appView.Namespace, gomock.Any()).Return(nil, fmt.Errorf("error")).Times(1)
	w = httptest.NewRecorder()
	body, _ = json.Marshal(appView)
	req, _ = http.NewRequest(http.MethodPost, "/v1/apps?base=eden2", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mkConfigService.EXPECT().Get(appView.Namespace, "agent-conf", "").Return(config, nil).Times(1)
	mSecretService.EXPECT().Get(appView.Namespace, "secret01", "").Return(secret, nil).Times(1)
	mSecretService.EXPECT().Get(appView.Namespace, "registry01", "").Return(secret, nil).Times(1)
	mkApplicationService.EXPECT().Get(appView.Namespace, "abc", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	mkApplicationService.EXPECT().Get(appView.Namespace, "eden2", "").Return(eden2, nil).Return(eden2, nil).Times(1)
	mkApplicationService.EXPECT().CreateWithBase(appView.Namespace, gomock.Any(), eden2).Return(eden2, nil)
	mkNodeService.EXPECT().UpdateNodeAppVersion(appView.Namespace, gomock.Any()).Return([]string{}, nil).Times(1)
	mkIndexService.EXPECT().RefreshNodesIndexByApp(appView.Namespace, gomock.Any(), gomock.Any()).Return(nil).Times(1)
	mSecretService.EXPECT().Get(appView.Namespace, "secret01", "").Return(secret, nil).Times(1)
	w = httptest.NewRecorder()
	body, _ = json.Marshal(appView)
	req, _ = http.NewRequest(http.MethodPost, "/v1/apps?base=eden2", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateContainerApplication(t *testing.T) {
	api, router, mockCtl := initApplicationAPI(t)
	defer mockCtl.Finish()
	mkApplicationService := ms.NewMockApplicationService(mockCtl)
	mkIndexService := ms.NewMockIndexService(mockCtl)
	mSecretService := ms.NewMockSecretService(mockCtl)
	mkConfigService := ms.NewMockConfigService(mockCtl)
	mkNodeService := ms.NewMockNodeService(mockCtl)
	api.applicationService = mkApplicationService
	api.indexService = mkIndexService
	api.secretService = mSecretService
	api.configService = mkConfigService
	api.nodeService = mkNodeService

	mApp := getMockContainerApp()
	mApp.Selector = "label = test"

	config := &specV1.Configuration{Name: "agent-conf", Version: "123"}
	secret1 := &specV1.Secret{Name: "registry01", Version: "123", Labels: map[string]string{specV1.SecretLabel: specV1.SecretRegistry}}
	secret2 := &specV1.Secret{Name: "secret01", Version: "123"}
	registry := &models.Registry{Name: "registry01", Version: "1"}
	mkConfigService.EXPECT().Get(gomock.Any(), gomock.Any(), "").Return(config, nil).AnyTimes()
	mSecretService.EXPECT().Get(gomock.Any(), secret2.Name, gomock.Any()).Return(secret2, nil).AnyTimes()
	mSecretService.EXPECT().Get(gomock.Any(), registry.Name, gomock.Any()).Return(secret1, nil).AnyTimes()

	mkApplicationService.EXPECT().Get(mApp.Namespace, "abc", gomock.Any()).Return(nil, fmt.Errorf("error")).Times(1)
	w := httptest.NewRecorder()
	body, _ := json.Marshal(mApp)
	req, _ := http.NewRequest(http.MethodPut, "/v1/apps/abc", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mApp2 := getMockContainerApp()
	mApp2.Selector = "name = test"

	mkApplicationService.EXPECT().Get(gomock.Any(), "abc", gomock.Any()).Return(mApp, nil).AnyTimes()
	mkApplicationService.EXPECT().Update(mApp.Namespace, gomock.Any()).Return(mApp2, nil)
	mkIndexService.EXPECT().RefreshNodesIndexByApp(mApp.Namespace, mApp.Name, gomock.Any()).Return(nil).Times(2)
	mkNodeService.EXPECT().DeleteNodeAppVersion(gomock.Any(), gomock.Any()).Return(nil, nil)
	mkNodeService.EXPECT().UpdateNodeAppVersion(gomock.Any(), gomock.Any()).Return([]string{}, nil)
	w = httptest.NewRecorder()
	body, _ = json.Marshal(mApp)
	req, _ = http.NewRequest(http.MethodPut, "/v1/apps/abc", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	mkApplicationService.EXPECT().Update(mApp.Namespace, gomock.Any()).Return(nil, fmt.Errorf("error"))
	w = httptest.NewRecorder()
	body, _ = json.Marshal(mApp)
	req, _ = http.NewRequest(http.MethodPut, "/v1/apps/abc", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mkApplicationService.EXPECT().Update(mApp.Namespace, gomock.Any()).Return(mApp2, nil)
	mkNodeService.EXPECT().DeleteNodeAppVersion(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
	w = httptest.NewRecorder()
	body, _ = json.Marshal(mApp)
	req, _ = http.NewRequest(http.MethodPut, "/v1/apps/abc", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// 500
	mkApplicationService.EXPECT().Get(gomock.Any(), "abc", gomock.Any()).Return(nil, nil).AnyTimes()
	mkApplicationService.EXPECT().Update(mApp.Namespace, gomock.Any()).Return(mApp, nil)
	mkNodeService.EXPECT().UpdateNodeAppVersion(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
	w = httptest.NewRecorder()
	body, _ = json.Marshal(mApp)
	req, _ = http.NewRequest(http.MethodPut, "/v1/apps/abc", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCreateFunctionApplication(t *testing.T) {
	api, router, mockCtl := initApplicationAPI(t)
	defer mockCtl.Finish()

	mkApplicationService := ms.NewMockApplicationService(mockCtl)
	mkIndexService := ms.NewMockIndexService(mockCtl)
	mSecretService := ms.NewMockSecretService(mockCtl)
	mkConfigService := ms.NewMockConfigService(mockCtl)
	mkNodeService := ms.NewMockNodeService(mockCtl)
	mkSysConfigService := ms.NewMockSysConfigService(mockCtl)

	api.applicationService = mkApplicationService
	api.indexService = mkIndexService
	api.secretService = mSecretService
	api.configService = mkConfigService
	api.nodeService = mkNodeService
	api.sysConfigService = mkSysConfigService

	appView := &models.ApplicationView{
		Application: specV1.Application{
			Namespace: "baetyl-cloud",
			Name:      "abc",
			Type:      common.FunctionApp,
			Services: []specV1.Service{
				{
					Name:     "Agent",
					Hostname: "test-agent",
					Replica:  1,
					VolumeMounts: []specV1.VolumeMount{
						{
							Name:      "baetyl-function-code-Agent",
							MountPath: "/var/lib/baetyl/code",
						},
					},
					Devices: []specV1.Device{
						{
							DevicePath: "DevicePath",
						},
					},
					FunctionConfig: &specV1.ServiceFunctionConfig{
						Name:    "func1",
						Runtime: "python36",
					},
					Functions: []specV1.ServiceFunction{
						{
							Name:    "process",
							Handler: "index.handler",
							CodeDir: "path",
						},
					},
				},
			},
			Volumes: []specV1.Volume{
				{
					Name: "baetyl-function-code-Agent",
					VolumeSource: specV1.VolumeSource{
						Config: &specV1.ObjectReference{
							Name: "func1",
						},
					},
				},
			},
		},
	}

	config := &specV1.Configuration{}
	app := &specV1.Application{}
	mkConfigService.EXPECT().Get(appView.Namespace, "func1", "").Return(config, nil).Times(1)
	mkApplicationService.EXPECT().Get(appView.Namespace, "abc", "").Return(app, nil).Times(1)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(appView)
	req, _ := http.NewRequest(http.MethodPost, "/v1/apps?base=eden", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	mkConfigService.EXPECT().Get(appView.Namespace, "func1", "").Return(nil, fmt.Errorf("config not found")).Times(1)

	w = httptest.NewRecorder()
	body, _ = json.Marshal(appView)
	req, _ = http.NewRequest(http.MethodPost, "/v1/apps?base=eden", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mkConfigService.EXPECT().Get(appView.Namespace, "func1", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	w = httptest.NewRecorder()
	body, _ = json.Marshal(appView)
	req, _ = http.NewRequest(http.MethodPost, "/v1/apps?base=eden", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	mkConfigService.EXPECT().Get(appView.Namespace, "func1", "").Return(config, nil).Times(1)
	mkApplicationService.EXPECT().Get(appView.Namespace, "abc", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	mkApplicationService.EXPECT().Get(appView.Namespace, "eden", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	w = httptest.NewRecorder()
	body, _ = json.Marshal(appView)
	req, _ = http.NewRequest(http.MethodPost, "/v1/apps?base=eden", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	eden2 := &specV1.Application{
		Namespace: appView.Namespace,
		Name:      "abc",
		Type:      common.ContainerApp,
		Services:  []specV1.Service{},
		Volumes: []specV1.Volume{
			{
				Name: "volume2",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name:    "config2",
						Version: "2",
					},
				},
			},
		},
	}

	mkConfigService.EXPECT().Get(appView.Namespace, "func1", "").Return(config, nil).Times(1)
	mkApplicationService.EXPECT().Get(appView.Namespace, "abc", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	mkApplicationService.EXPECT().Get(appView.Namespace, "eden2", "").Return(eden2, nil).Times(1)

	w = httptest.NewRecorder()
	body, _ = json.Marshal(appView)
	req, _ = http.NewRequest(http.MethodPost, "/v1/apps?base=eden2", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	eden2 = &specV1.Application{
		Namespace: appView.Namespace,
		Name:      "abc",
		Type:      common.FunctionApp,
		Services:  []specV1.Service{},
		Volumes: []specV1.Volume{
			{
				Name: "volume2",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name:    "config2",
						Version: "2",
					},
				},
			},
		},
	}

	mkConfigService.EXPECT().Get(appView.Namespace, "func1", "").Return(config, nil).Times(1)
	mkApplicationService.EXPECT().Get(appView.Namespace, "abc", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	mkApplicationService.EXPECT().Get(appView.Namespace, "eden2", "").Return(eden2, nil).Times(1)
	mkSysConfigService.EXPECT().GetSysConfig(gomock.Any(), gomock.Any()).Return(nil, errors.New("err")).Times(1)

	w = httptest.NewRecorder()
	body, _ = json.Marshal(appView)
	req, _ = http.NewRequest(http.MethodPost, "/v1/apps?base=eden2", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	functions := models.ServiceFunction{
		Functions: []specV1.ServiceFunction{
			{
				Name:    "process",
				Handler: "index.handler",
				CodeDir: "path",
			},
		},
	}
	data, err := json.Marshal(&functions)
	assert.NoError(t, err)
	config = &specV1.Configuration{
		Name: fmt.Sprintf("baetyl-function-%s-%s-%s", app.Name, "Agent", common.RandString(9)),
		Labels: map[string]string{
			common.LabelSystem: "true",
		},
		Data: map[string]string{
			"service.yml": string(data),
		},
	}
	mkConfigService.EXPECT().Get(appView.Namespace, "func1", "").Return(config, nil).Times(1)
	mkApplicationService.EXPECT().Get(appView.Namespace, "abc", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	mkApplicationService.EXPECT().Get(appView.Namespace, "eden2", "").Return(eden2, nil).Times(1)
	sysconfig := &models.SysConfig{
		Type:  common.BaetylFunctionRuntime,
		Key:   "python36",
		Value: "image",
	}
	mkSysConfigService.EXPECT().GetSysConfig(gomock.Any(), gomock.Any()).Return(sysconfig, nil).Times(1)
	mkConfigService.EXPECT().Upsert(appView.Namespace, gomock.Any()).Return(nil, errors.New("err")).Times(1)

	w = httptest.NewRecorder()
	body, _ = json.Marshal(appView)
	req, _ = http.NewRequest(http.MethodPost, "/v1/apps?base=eden2", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mkConfigService.EXPECT().Get(appView.Namespace, "func1", "").Return(config, nil).Times(1)
	mkApplicationService.EXPECT().Get(appView.Namespace, "abc", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	mkApplicationService.EXPECT().Get(appView.Namespace, "eden2", "").Return(eden2, nil).Times(1)
	mkSysConfigService.EXPECT().GetSysConfig(gomock.Any(), gomock.Any()).Return(sysconfig, nil).Times(1)
	mkConfigService.EXPECT().Upsert(appView.Namespace, gomock.Any()).Return(config, nil).Times(1)
	mkApplicationService.EXPECT().CreateWithBase(appView.Namespace, gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error")).Times(1)

	w = httptest.NewRecorder()
	body, _ = json.Marshal(appView)
	req, _ = http.NewRequest(http.MethodPost, "/v1/apps?base=eden2", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mkConfigService.EXPECT().Get(appView.Namespace, "func1", "").Return(config, nil).Times(1)
	mkApplicationService.EXPECT().Get(appView.Namespace, "abc", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	mkApplicationService.EXPECT().Get(appView.Namespace, "eden2", "").Return(eden2, nil).Times(1)
	mkSysConfigService.EXPECT().GetSysConfig(gomock.Any(), gomock.Any()).Return(sysconfig, nil).Times(1)
	mkConfigService.EXPECT().Upsert(appView.Namespace, gomock.Any()).Return(config, nil).Times(1)

	mApp := &specV1.Application{
		Namespace: "baetyl-cloud",
		Name:      "abc",
		Type:      common.FunctionApp,
		Services: []specV1.Service{
			{
				Name:     "Agent",
				Image:    "image",
				Hostname: "test-agent",
				Replica:  1,
				VolumeMounts: []specV1.VolumeMount{
					{
						Name:      "baetyl-function-code-Agent",
						MountPath: "/var/lib/baetyl/code",
					},
					{
						Name:      "baetyl-function-config-Agent",
						MountPath: "/etc/baetyl",
					},
				},
				Devices: []specV1.Device{
					{
						DevicePath: "DevicePath",
					},
				},
				FunctionConfig: &specV1.ServiceFunctionConfig{
					Name:    "func1",
					Runtime: "python36",
				},
				Functions: []specV1.ServiceFunction{
					{
						Name:    "process",
						Handler: "index.handler",
						CodeDir: "path",
					},
				},
			},
		},
		Volumes: []specV1.Volume{
			{
				Name: "baetyl-function-code-Agent",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name: "func1",
					},
				},
			},
			{
				Name: "baetyl-function-config-Agent",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name: "baetyl-function-app-service-xxx",
					},
				},
			},
		},
	}
	mkApplicationService.EXPECT().CreateWithBase(appView.Namespace, gomock.Any(), gomock.Any()).Return(mApp, nil).Times(1)
	mkNodeService.EXPECT().UpdateNodeAppVersion(appView.Namespace, gomock.Any()).Return(nil, fmt.Errorf("error")).Times(1)

	w = httptest.NewRecorder()
	body, _ = json.Marshal(appView)
	req, _ = http.NewRequest(http.MethodPost, "/v1/apps?base=eden2", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mkConfigService.EXPECT().Get(appView.Namespace, "func1", "").Return(config, nil).Times(1)
	mkApplicationService.EXPECT().Get(appView.Namespace, "abc", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	mkApplicationService.EXPECT().Get(appView.Namespace, "eden2", "").Return(eden2, nil).Times(1)
	mkSysConfigService.EXPECT().GetSysConfig(gomock.Any(), gomock.Any()).Return(sysconfig, nil).Times(1)
	mkConfigService.EXPECT().Upsert(appView.Namespace, gomock.Any()).Return(config, nil).Times(1)
	mkApplicationService.EXPECT().CreateWithBase(appView.Namespace, gomock.Any(), gomock.Any()).Return(mApp, nil).Times(1)
	mkNodeService.EXPECT().UpdateNodeAppVersion(appView.Namespace, gomock.Any()).Return([]string{}, nil).Times(1)
	mkIndexService.EXPECT().RefreshNodesIndexByApp(appView.Namespace, gomock.Any(), gomock.Any()).Return(nil).Times(1)
	mkConfigService.EXPECT().Get(appView.Namespace, gomock.Any(), "").Return(config, nil).AnyTimes()

	w = httptest.NewRecorder()
	body, _ = json.Marshal(appView)
	req, _ = http.NewRequest(http.MethodPost, "/v1/apps?base=eden2", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateFunctionApplication(t *testing.T) {
	api, router, mockCtl := initApplicationAPI(t)
	defer mockCtl.Finish()
	mkApplicationService := ms.NewMockApplicationService(mockCtl)
	mkIndexService := ms.NewMockIndexService(mockCtl)
	mSecretService := ms.NewMockSecretService(mockCtl)
	mkConfigService := ms.NewMockConfigService(mockCtl)
	mkNodeService := ms.NewMockNodeService(mockCtl)
	mkSysConfigService := ms.NewMockSysConfigService(mockCtl)

	api.applicationService = mkApplicationService
	api.indexService = mkIndexService
	api.secretService = mSecretService
	api.configService = mkConfigService
	api.nodeService = mkNodeService
	api.sysConfigService = mkSysConfigService

	namespace := "baetyl-cloud"
	oldApp := &specV1.Application{
		Namespace: namespace,
		Name:      "abc",
		Type:      common.FunctionApp,
		Services: []specV1.Service{
			{
				Name:     "Agent",
				Hostname: "test-agent",
				Replica:  1,
				VolumeMounts: []specV1.VolumeMount{
					{
						Name:      "baetyl-function-code-Agent",
						MountPath: "/var/lib/baetyl/code",
					},
					{
						Name:      "baetyl-function-config-Agent",
						MountPath: "/etc/baetyl",
					},
				},
				Devices: []specV1.Device{
					{
						DevicePath: "DevicePath",
					},
				},
				FunctionConfig: &specV1.ServiceFunctionConfig{
					Name:    "func1",
					Runtime: "python36",
				},
			},
			{
				Name:     "Agent2",
				Hostname: "test-agent2",
				Replica:  1,
				VolumeMounts: []specV1.VolumeMount{
					{
						Name:      "baetyl-function-code-Agent2",
						MountPath: "/var/lib/baetyl/code",
					},
					{
						Name:      "baetyl-function-config-Agent2",
						MountPath: "/etc/baetyl",
					},
				},
				Devices: []specV1.Device{
					{
						DevicePath: "DevicePath",
					},
				},
				Args: []string{"test"},
				FunctionConfig: &specV1.ServiceFunctionConfig{
					Name:    "func2",
					Runtime: "python36",
				},
			},
		},
		Volumes: []specV1.Volume{
			{
				Name: "baetyl-funciton-code-Agent",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name: "func1",
					},
				},
			},
			{
				Name: "baetyl-function-code-Agent2",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name: "func2",
					},
				},
			},
			{
				Name: "baetyl-function-config-Agent",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name: "baetyl-function-config-app-service-1",
					},
				},
			},
			{
				Name: "baetyl-function-config-Agent2",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name: "baetyl-function-config-app-service-2",
					},
				},
			},
		},
	}

	functions2 := models.ServiceFunction{
		Functions: []specV1.ServiceFunction{
			{
				Name:    "process2",
				Handler: "index.handler",
				CodeDir: "path",
			},
		},
	}
	data2, err := json.Marshal(&functions2)
	assert.NoError(t, err)
	config2 := &specV1.Configuration{
		Name: fmt.Sprintf("baetyl-function-config-%s-%s-%s", oldApp.Name, "Agent", common.RandString(9)),
		Labels: map[string]string{
			common.LabelSystem: "true",
		},
		Data: map[string]string{
			"service.yml": string(data2),
		},
	}

	functions2extra := models.ServiceFunction{
		Functions: []specV1.ServiceFunction{
			{
				Name:    "process2",
				Handler: "index.handler",
				CodeDir: "path",
			},
			{
				Name:    "process2extra",
				Handler: "index.handler",
				CodeDir: "path",
			},
		},
	}
	data2extra, err := json.Marshal(&functions2extra)
	assert.NoError(t, err)
	config2extra := &specV1.Configuration{
		Name: fmt.Sprintf("baetyl-function-config-%s-%s-%s", oldApp.Name, "Agent", common.RandString(9)),
		Labels: map[string]string{
			common.LabelSystem: "true",
		},
		Data: map[string]string{
			"service.yml": string(data2extra),
		},
	}

	functions3 := models.ServiceFunction{
		Functions: []specV1.ServiceFunction{
			{
				Name:    "process3",
				Handler: "index.handler",
				CodeDir: "path",
			},
		},
	}
	data3, err := json.Marshal(&functions3)
	assert.NoError(t, err)
	config3 := &specV1.Configuration{
		Name: fmt.Sprintf("baetyl-function-config-%s-%s-%s", oldApp.Name, "Agent", common.RandString(9)),
		Labels: map[string]string{
			common.LabelSystem: "true",
		},
		Data: map[string]string{
			"service.yml": string(data3),
		},
	}

	sysconfig := &models.SysConfig{
		Type:  common.BaetylFunctionRuntime,
		Key:   "python36",
		Value: "image",
	}

	newApp := &specV1.Application{
		Namespace: namespace,
		Name:      "abc",
		Type:      common.FunctionApp,
		Services: []specV1.Service{
			{
				Name:     "Agent2",
				Hostname: "test-agent2",
				Replica:  1,
				VolumeMounts: []specV1.VolumeMount{
					{
						Name:      "baetyl-function-code-Agent2",
						MountPath: "/var/lib/baetyl/code",
					},
					{
						Name:      "baetyl-function-config-Agent2",
						MountPath: "/etc/baetyl",
					},
				},
				Devices: []specV1.Device{
					{
						DevicePath: "DevicePath",
					},
				},
				FunctionConfig: &specV1.ServiceFunctionConfig{
					Name:    "func2",
					Runtime: "python36",
				},
				Functions: []specV1.ServiceFunction{
					{
						Name:    "process2",
						Handler: "index.handler",
						CodeDir: "path",
					},
					{
						Name:    "process2extra",
						Handler: "index.handler",
						CodeDir: "path",
					},
				},
			},
			{
				Name:     "Agent3",
				Hostname: "test-agent3",
				Replica:  1,
				VolumeMounts: []specV1.VolumeMount{
					{
						Name:      "baetyl-function-code-Agent3",
						MountPath: "/var/lib/baetyl/code",
					},
					{
						Name:      "baetyl-function-config-Agent3",
						MountPath: "/etc/baetyl",
					},
				},
				Devices: []specV1.Device{
					{
						DevicePath: "DevicePath",
					},
				},
				FunctionConfig: &specV1.ServiceFunctionConfig{
					Name:    "func3",
					Runtime: "python36",
				},
				Functions: []specV1.ServiceFunction{
					{
						Name:    "process3",
						Handler: "index.handler",
						CodeDir: "path",
					},
				},
			},
		},
		Volumes: []specV1.Volume{
			{
				Name: "baetyl-function-code-Agent2",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name: "func2",
					},
				},
			},
			{
				Name: "baetyl-function-code-Agent3",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name: "func3",
					},
				},
			},
			{
				Name: "baetyl-function-config-Agent2",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name: "baetyl-function-config-app-service-2",
					},
				},
			},
			{
				Name: "baetyl-function-config-Agent3",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name: "baetyl-function-config-app-service-3",
					},
				},
			},
		},
	}
	newAppView := &models.ApplicationView{
		Application: specV1.Application{
			Namespace: namespace,
			Type:      common.FunctionApp,
			Services: []specV1.Service{
				{
					Name:     "Agent2",
					Hostname: "test-agent2",
					Replica:  1,
					VolumeMounts: []specV1.VolumeMount{
						{
							Name:      "baetyl-funciton-code-Agent2",
							MountPath: "/var/lib/baetyl/code",
						},
						{
							Name:      "baetyl-function-config-Agent2",
							MountPath: "/etc/baetyl",
						},
					},
					Devices: []specV1.Device{
						{
							DevicePath: "DevicePath",
						},
					},
					Args: []string{"test"},
					FunctionConfig: &specV1.ServiceFunctionConfig{
						Name:    "func2",
						Runtime: "python36",
					},
					Functions: []specV1.ServiceFunction{
						{
							Name:    "process2",
							Handler: "index.handler",
							CodeDir: "path",
						},
						{
							Name:    "process2extra",
							Handler: "index.handler",
							CodeDir: "path",
						},
					},
				},
				{
					Name:     "Agent3",
					Hostname: "test-agent3",
					Replica:  1,
					VolumeMounts: []specV1.VolumeMount{
						{
							Name:      "baetyl-function-code-Agent2",
							MountPath: "/var/lib/baetyl/code",
						},
					},
					Devices: []specV1.Device{
						{
							DevicePath: "DevicePath",
						},
					},
					FunctionConfig: &specV1.ServiceFunctionConfig{
						Name:    "func3",
						Runtime: "python36",
					},
					Functions: []specV1.ServiceFunction{
						{
							Name:    "process3",
							Handler: "index.handler",
							CodeDir: "path",
						},
					},
				},
			},
			Volumes: []specV1.Volume{
				{
					Name: "baetyl-function-code-Agent2",
					VolumeSource: specV1.VolumeSource{
						Config: &specV1.ObjectReference{
							Name: "func2",
						},
					},
				},
				{
					Name: "baetyl-function-code-Agent3",
					VolumeSource: specV1.VolumeSource{
						Config: &specV1.ObjectReference{
							Name: "func3",
						},
					},
				},
				{
					Name: "baetyl-function-config-Agent2",
					VolumeSource: specV1.VolumeSource{
						Config: &specV1.ObjectReference{
							Name: "baetyl-function-config-app-service-2",
						},
					},
				},
			},
		},
	}

	configCode := &specV1.Configuration{}
	mkApplicationService.EXPECT().Get(namespace, "abc", "").Return(oldApp, nil).Times(1)
	mkConfigService.EXPECT().Get(namespace, "func2", "").Return(configCode, nil).Times(1)
	mkConfigService.EXPECT().Get(namespace, "func3", "").Return(configCode, nil).Times(1)
	mkConfigService.EXPECT().Get(namespace, "baetyl-function-config-app-service-2", "").Return(config2, nil).Times(1)
	mkSysConfigService.EXPECT().GetSysConfig(gomock.Any(), gomock.Any()).Return(sysconfig, nil).Times(2)
	mkConfigService.EXPECT().Upsert(namespace, gomock.Any()).Return(config2extra, nil).Times(1)
	mkConfigService.EXPECT().Upsert(namespace, gomock.Any()).Return(config3, nil).Times(1)
	mkApplicationService.EXPECT().Update(namespace, gomock.Any()).Return(newApp, nil).Times(1)
	mkNodeService.EXPECT().UpdateNodeAppVersion(namespace, gomock.Any()).Return([]string{}, nil).Times(1)
	mkIndexService.EXPECT().RefreshNodesIndexByApp(namespace, gomock.Any(), gomock.Any()).Return(nil).Times(1)
	mkConfigService.EXPECT().Delete(namespace, gomock.Any()).Return(nil).Times(1)
	mkConfigService.EXPECT().Get(namespace, "baetyl-function-config-app-service-2", "").Return(config2, nil).Times(1)
	mkConfigService.EXPECT().Get(namespace, "baetyl-function-config-app-service-3", "").Return(config2, nil).Times(1)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(newAppView)
	req, _ := http.NewRequest(http.MethodPut, "/v1/apps/abc", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteApplication(t *testing.T) {
	api, router, mockCtl := initApplicationAPI(t)
	defer mockCtl.Finish()
	mkApplicationService := ms.NewMockApplicationService(mockCtl)
	mkIndexService := ms.NewMockIndexService(mockCtl)
	mkNodeService := ms.NewMockNodeService(mockCtl)
	mkConfigService := ms.NewMockConfigService(mockCtl)

	api.applicationService = mkApplicationService
	api.indexService = mkIndexService
	api.nodeService = mkNodeService
	api.configService = mkConfigService

	app := &specV1.Application{
		Namespace: "baetyl-cloud",
		Name:      "abc",
	}

	// 404
	mkApplicationService.EXPECT().Get(gomock.Any(), "abc", gomock.Any()).Return(nil, common.Error(common.ErrResourceNotFound))
	req, _ := http.NewRequest(http.MethodDelete, "/v1/apps/abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	// 500
	mkApplicationService.EXPECT().Get(gomock.Any(), "abc", gomock.Any()).Return(app, nil).Times(1)
	mkApplicationService.EXPECT().Delete(app.Namespace, app.Name, "").Return(fmt.Errorf("error")).Times(1)
	req, _ = http.NewRequest(http.MethodDelete, "/v1/apps/abc", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// 500
	mkApplicationService.EXPECT().Get(gomock.Any(), "abc", gomock.Any()).Return(app, nil).Times(1)
	mkApplicationService.EXPECT().Delete(app.Namespace, app.Name, "").Return(nil).Times(1)
	mkNodeService.EXPECT().DeleteNodeAppVersion(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
	req, _ = http.NewRequest(http.MethodDelete, "/v1/apps/abc", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// 200
	mkApplicationService.EXPECT().Get(gomock.Any(), "abc", gomock.Any()).Return(app, nil).Times(1)
	mkApplicationService.EXPECT().Delete(app.Namespace, app.Name, "").Return(nil).Times(1)
	mkIndexService.EXPECT().RefreshNodesIndexByApp(app.Namespace, app.Name, gomock.Any()).Return(nil).Times(1)
	mkNodeService.EXPECT().DeleteNodeAppVersion(gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)
	req, _ = http.NewRequest(http.MethodDelete, "/v1/apps/abc", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	app.Name = "baetyl-core-test"
	mkApplicationService.EXPECT().Get(gomock.Any(), app.Name, gomock.Any()).Return(app, nil).AnyTimes()
	mkIndexService.EXPECT().ListNodesByApp(gomock.Any(), app.Name).Return(nil, fmt.Errorf("error"))
	req, _ = http.NewRequest(http.MethodDelete, "/v1/apps/baetyl-core-test", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mkIndexService.EXPECT().ListNodesByApp(gomock.Any(), app.Name).Return([]string{"test-node"}, nil)
	req, _ = http.NewRequest(http.MethodDelete, "/v1/apps/baetyl-core-test", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Function
	app = &specV1.Application{
		Namespace: "baetyl-cloud",
		Name:      "abc",
		Type:      common.FunctionApp,
		Services: []specV1.Service{
			{
				Name:     "Agent",
				Hostname: "test-agent",
				Replica:  1,
				VolumeMounts: []specV1.VolumeMount{
					{
						Name:      "baetyl-function-conf-abc",
						MountPath: "mountPath",
					},
					{
						Name:      "baetyl-function-conf-func1",
						MountPath: "mountPath",
					},
				},
				Devices: []specV1.Device{
					{
						DevicePath: "DevicePath",
					},
				},
				FunctionConfig: &specV1.ServiceFunctionConfig{
					Name:    "func1",
					Runtime: "python36",
				},
				Functions: []specV1.ServiceFunction{
					{
						Name:    "process",
						Handler: "index.handler",
						CodeDir: "path",
					},
				},
			},
		},
		Volumes: []specV1.Volume{
			{
				Name: "baetyl-function-code-Agent",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name: "func1",
					},
				},
			},
			{
				Name: "baetyl-function-conf-Agent",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name: "baetyl-function-config-app-service-xxxxxxxxx",
					},
				},
			},
		},
	}

	// 200
	mkApplicationService.EXPECT().Get(gomock.Any(), "abc", gomock.Any()).Return(app, nil).Times(1)
	mkIndexService.EXPECT().RefreshNodesIndexByApp(app.Namespace, app.Name, gomock.Any()).Return(nil)
	mkNodeService.EXPECT().DeleteNodeAppVersion(gomock.Any(), gomock.Any()).Return(nil, nil)
	mkApplicationService.EXPECT().Delete(app.Namespace, app.Name, "").Return(nil)
	mkConfigService.EXPECT().Delete(app.Namespace, "baetyl-function-config-app-service-xxxxxxxxx").Return(nil).Times(1)

	req, _ = http.NewRequest(http.MethodDelete, "/v1/apps/abc", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
