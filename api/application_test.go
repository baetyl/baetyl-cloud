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

	"github.com/baetyl/baetyl-go/v2/context"
	"github.com/baetyl/baetyl-go/v2/log"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/common"
	mf "github.com/baetyl/baetyl-cloud/v2/mock/facade"
	ms "github.com/baetyl/baetyl-cloud/v2/mock/service"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/service"
)

func getMockContainerApp() *specV1.Application {
	mApp := &specV1.Application{
		Namespace: "baetyl-cloud",
		Name:      "abc",
		Type:      common.ContainerApp,
		Replica:   1,
		Services: []specV1.Service{
			{
				Name:     "agent",
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
				Name:     "agent",
				Hostname: "test-agent",
				Replica:  1,
				VolumeMounts: []specV1.VolumeMount{
					{
						Name:      "baetyl-function-code-agent",
						MountPath: "mountPath",
					},
					{
						Name:      "baetyl-function-config-agent",
						MountPath: "mountPath",
					},
					{
						Name:      "baetyl-function-program-config-agent",
						MountPath: "/var/lib/baetyl/bin",
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
				Name: "baetyl-function-code-agent",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name: "func1",
					},
				},
			},
			{
				Name: "baetyl-function-config-agent",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name: "baetyl-function-app-service-xxxxxxxxx",
					},
				},
			},
			{
				Name: "baetyl-function-program-config-agent",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name: "baetyl-function-program-config-x3-xs3-uwredcfxb",
					},
				},
			},
		},
	}
	return mApp
}

func initApplicationAPI(t *testing.T) (*API, *gin.Engine, *gomock.Controller) {
	api := &API{log: log.L().With(log.Any("test", "api"))}
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
		configs.GET("/:name/configs", mockIM, common.Wrapper(api.GetSysAppConfigs))
		configs.GET("/:name/secrets", mockIM, common.Wrapper(api.GetSysAppSecrets))
		configs.GET("/:name/certificates", mockIM, common.Wrapper(api.GetSysAppCertificates))
		configs.GET("/:name/registries", mockIM, common.Wrapper(api.GetSysAppRegistries))
	}
	return api, router, mockCtl
}

func TestGetInvisibleApplication(t *testing.T) {
	api, router, mockCtl := initApplicationAPI(t)
	defer mockCtl.Finish()

	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sSecret := ms.NewMockSecretService(mockCtl)
	fApp := mf.NewMockFacade(mockCtl)
	api.Facade = fApp
	api.AppCombinedService = &service.AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}

	mApp := &specV1.Application{
		Namespace: "baetyl-cloud",
		Name:      "abc",
		Type:      common.ContainerApp,
		Labels: map[string]string{
			common.ResourceInvisible: "true",
		},
		Services: []specV1.Service{
			{
				Name:     "agent",
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
	fApp.EXPECT().GetApp(mApp.Namespace, "cba", "").Return(mApp, nil).Times(1)
	req, _ := http.NewRequest(http.MethodGet, "/v1/apps/cba", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetContainerApplication(t *testing.T) {
	api, router, mockCtl := initApplicationAPI(t)
	defer mockCtl.Finish()

	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sSecret := ms.NewMockSecretService(mockCtl)
	api.AppCombinedService = &service.AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}
	fApp := mf.NewMockFacade(mockCtl)
	api.Facade = fApp

	mApp := getMockContainerApp()
	fApp.EXPECT().GetApp(mApp.Namespace, "cba", "").Return(nil, errors.New("err")).Times(1)
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
	fApp.EXPECT().GetApp(mApp.Namespace, mApp.Name, "").Return(mApp, nil).Times(1)
	sSecret.EXPECT().Get(mApp.Namespace, secret.Name, "").Return(secret, nil).Times(1)

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

	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sSecret := ms.NewMockSecretService(mockCtl)
	api.AppCombinedService = &service.AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}
	fApp := mf.NewMockFacade(mockCtl)
	api.Facade = fApp

	mApp := getMockFunctionApp()

	fApp.EXPECT().GetApp(mApp.Namespace, "cba", "").Return(nil, errors.New("err")).Times(1)
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
	fApp.EXPECT().GetApp(mApp.Namespace, mApp.Name, "").Return(mApp, nil).Times(1)
	sConfig.EXPECT().Get(mApp.Namespace, "baetyl-function-app-service-xxxxxxxxx", "").Return(config, nil).Times(1)
	sConfig.EXPECT().Get(mApp.Namespace, "baetyl-function-program-config-x3-xs3-uwredcfxb", "").Return(config, nil).Times(1)

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

	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sSecret := ms.NewMockSecretService(mockCtl)
	api.AppCombinedService = &service.AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}

	mClist := &models.ApplicationList{}

	sApp.EXPECT().List("baetyl-cloud", &models.ListOptions{
		LabelSelector: "!" + common.LabelSystem,
	}).Return(mClist, nil)

	// 200
	req, _ := http.NewRequest(http.MethodGet, "/v1/apps", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	sApp.EXPECT().List("baetyl-cloud", &models.ListOptions{
		LabelSelector: "!" + common.LabelSystem,
	}).Return(nil, fmt.Errorf("error"))

	// 400
	req, _ = http.NewRequest(http.MethodGet, "/v1/apps", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	sApp.EXPECT().List("baetyl-cloud", &models.ListOptions{
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

	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sSecret := ms.NewMockSecretService(mockCtl)
	api.AppCombinedService = &service.AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}

	sNode := ms.NewMockNodeService(mockCtl)
	sIndex := ms.NewMockIndexService(mockCtl)
	fApp := mf.NewMockFacade(mockCtl)
	api.Facade = fApp
	api.Index = sIndex
	api.Node = sNode

	appView := &models.ApplicationView{
		Namespace: "baetyl-cloud",
		Name:      "abc",
		Type:      common.ContainerApp,
		Services: []models.ServiceView{
			{
				Service: specV1.Service{
					Name:     "agent",
					Hostname: "test-agent",
					Image:    "hub.baidubce.com/baetyl/baetyl-agent:1.0.0",
					Replica:  1,
					VolumeMounts: []specV1.VolumeMount{
						{
							Name:      "name",
							MountPath: "mountPath",
						},
						{
							Name:      "emptydir",
							MountPath: "mp",
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
		},
		Volumes: []models.VolumeView{
			{
				Name: "name",
				Config: &specV1.ObjectReference{
					Name: "agent-conf",
				},
			},
			{
				Name: "secret",
				Secret: &specV1.ObjectReference{
					Name: "secret01",
				},
			},
			{
				Name: "emptydir",
				EmptyDir: &specV1.EmptyDirVolumeSource{
					Medium:    "Memory",
					SizeLimit: "1",
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
	sConfig.EXPECT().Get(appView.Namespace, "agent-conf", "").Return(config, nil).Times(1)
	sSecret.EXPECT().Get(appView.Namespace, "secret01", "").Return(secret, nil).Times(1)
	sSecret.EXPECT().Get(appView.Namespace, "registry01", "").Return(secret, nil).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "abc", "").Return(app, nil).Times(1)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(appView)
	req, _ := http.NewRequest(http.MethodPost, "/v1/apps?base=eden", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	sConfig.EXPECT().Get(appView.Namespace, "agent-conf", "").Return(nil, fmt.Errorf("config not found")).Times(1)

	w = httptest.NewRecorder()
	body, _ = json.Marshal(appView)
	req, _ = http.NewRequest(http.MethodPost, "/v1/apps?base=eden", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	sConfig.EXPECT().Get(appView.Namespace, "agent-conf", "").Return(config, nil).Times(1)
	sSecret.EXPECT().Get(appView.Namespace, "secret01", "").Return(nil, fmt.Errorf("secret not found")).Times(1)

	w = httptest.NewRecorder()
	body, _ = json.Marshal(appView)
	req, _ = http.NewRequest(http.MethodPost, "/v1/apps?base=eden", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	sConfig.EXPECT().Get(appView.Namespace, "agent-conf", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	w = httptest.NewRecorder()
	body, _ = json.Marshal(appView)
	req, _ = http.NewRequest(http.MethodPost, "/v1/apps?base=eden", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	sConfig.EXPECT().Get(appView.Namespace, "agent-conf", "").Return(config, nil).Times(1)
	sSecret.EXPECT().Get(appView.Namespace, "secret01", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)

	w = httptest.NewRecorder()
	body, _ = json.Marshal(appView)
	req, _ = http.NewRequest(http.MethodPost, "/v1/apps?base=eden", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	sConfig.EXPECT().Get(appView.Namespace, "agent-conf", "").Return(config, nil).Times(1)
	sSecret.EXPECT().Get(appView.Namespace, "secret01", "").Return(secret, nil).Times(1)
	sSecret.EXPECT().Get(appView.Namespace, "registry01", "").Return(secret, nil).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "abc", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "eden", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
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

	sConfig.EXPECT().Get(appView.Namespace, "agent-conf", "").Return(config, nil).Times(1)
	sSecret.EXPECT().Get(appView.Namespace, "secret01", "").Return(secret, nil).Times(1)
	sSecret.EXPECT().Get(appView.Namespace, "registry01", "").Return(secret, nil).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "abc", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "eden2", "").Return(eden2, nil).Times(1)
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

	sConfig.EXPECT().Get(appView.Namespace, "agent-conf", "").Return(config, nil).Times(1)
	sSecret.EXPECT().Get(appView.Namespace, "secret01", "").Return(secret, nil).Times(1)
	sSecret.EXPECT().Get(appView.Namespace, "registry01", "").Return(secret, nil).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "abc", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "eden2", "").Return(eden2, nil).Times(1)
	fApp.EXPECT().CreateApp(appView.Namespace, gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error")).Times(1)
	w = httptest.NewRecorder()
	body, _ = json.Marshal(appView)
	req, _ = http.NewRequest(http.MethodPost, "/v1/apps?base=eden2", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mApp := getMockContainerApp()
	sConfig.EXPECT().Get(appView.Namespace, "agent-conf", "").Return(config, nil).Times(1)
	sSecret.EXPECT().Get(appView.Namespace, "secret01", "").Return(secret, nil).Times(2)
	sSecret.EXPECT().Get(appView.Namespace, "registry01", "").Return(secret, nil).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "abc", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "eden2", "").Return(eden2, nil).Return(eden2, nil).Times(1)
	fApp.EXPECT().CreateApp(appView.Namespace, gomock.Any(), gomock.Any(), gomock.Any()).Return(mApp, nil).Times(1)
	w = httptest.NewRecorder()
	body, _ = json.Marshal(appView)
	req, _ = http.NewRequest(http.MethodPost, "/v1/apps?base=eden2", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestInitContainerApp(t *testing.T) {
	api, router, mockCtl := initApplicationAPI(t)
	defer mockCtl.Finish()

	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sSecret := ms.NewMockSecretService(mockCtl)
	api.AppCombinedService = &service.AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}

	sNode := ms.NewMockNodeService(mockCtl)
	sIndex := ms.NewMockIndexService(mockCtl)
	fApp := mf.NewMockFacade(mockCtl)
	api.Facade = fApp
	api.Index = sIndex
	api.Node = sNode

	appView := &models.ApplicationView{
		Namespace: "baetyl-cloud",
		Name:      "abc",
		Type:      common.ContainerApp,
		InitServices: []models.ServiceView{
			{
				Service: specV1.Service{
					Name:  "init",
					Image: "hub.baidubce.com/baetyl/baetyl-init:1.0.0",
				},
			},
		},
		Services: []models.ServiceView{
			{
				Service: specV1.Service{
					Name:  "agent",
					Image: "hub.baidubce.com/baetyl/baetyl-agent:1.0.0",
				},
			},
		},
	}
	app := &specV1.Application{}
	copier.Copy(app, appView)

	sApp.EXPECT().Get(appView.Namespace, "abc", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	fApp.EXPECT().CreateApp(appView.Namespace, gomock.Any(), gomock.Any(), gomock.Any()).Return(app, nil).Times(1)
	w := httptest.NewRecorder()
	body, _ := json.Marshal(appView)
	req, _ := http.NewRequest(http.MethodPost, "/v1/apps", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateApplicationHasCertificates(t *testing.T) {
	api, router, mockCtl := initApplicationAPI(t)
	defer mockCtl.Finish()

	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sSecret := ms.NewMockSecretService(mockCtl)
	api.AppCombinedService = &service.AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}

	sNode := ms.NewMockNodeService(mockCtl)
	sIndex := ms.NewMockIndexService(mockCtl)
	fApp := mf.NewMockFacade(mockCtl)
	api.Facade = fApp
	api.Index = sIndex
	api.Node = sNode

	appView := &models.ApplicationView{
		Namespace: "baetyl-cloud",
		Name:      "abc",
		Type:      common.ContainerApp,
		Replica:   1,
		Workload:  specV1.WorkloadDeployment,
		Services: []models.ServiceView{
			{
				Service: specV1.Service{
					Name:     "agent",
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
		},
		Volumes: []models.VolumeView{
			{
				Name: "name",
				Config: &specV1.ObjectReference{
					Name: "agent-conf",
				},
			},
			{
				Name: "secret",
				Secret: &specV1.ObjectReference{
					Name: "secret01",
				},
			},
			{
				Name: "certificate",
				Certificate: &specV1.ObjectReference{
					Name: "certificate01",
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
	secretRegistry := &specV1.Secret{
		Name: "registry01",
		Labels: map[string]string{
			specV1.SecretLabel: specV1.SecretRegistry,
		},
	}
	secretCertificate := &specV1.Secret{
		Name: "certificate01",
		Labels: map[string]string{
			specV1.SecretLabel: specV1.SecretCertificate,
		},
	}
	app := &specV1.Application{
		Name:              "abc",
		Type:              common.ContainerApp,
		Labels:            nil,
		Mode:              "kube",
		Namespace:         "baetyl-cloud",
		CreationTimestamp: time.Time{},
		Version:           "",
		Selector:          "",
		Replica:           1,
		Workload:          specV1.WorkloadDeployment,
		JobConfig:         &specV1.AppJobConfig{RestartPolicy: "Never"},
		Services: []specV1.Service{
			{
				Name:     "agent",
				Hostname: "test-agent",
				Image:    "hub.baidubce.com/baetyl/baetyl-agent:1.0.0",
				Replica:  1,
				Type:     specV1.WorkloadDeployment,
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
			{
				Name: "certificate",
				VolumeSource: specV1.VolumeSource{
					Secret: &specV1.ObjectReference{
						Name: "certificate01",
					},
				},
			},
			{
				Name: "registry01",
				VolumeSource: specV1.VolumeSource{
					Secret: &specV1.ObjectReference{
						Name: "registry01",
					},
				},
			},
		},
		Description: "",
		System:      false,
	}
	app1 := app
	app1.Labels = map[string]string{
		common.LabelAppMode: context.RunModeKube,
	}

	log.L().Info("", log.Any("app1", app))

	sConfig.EXPECT().Get(appView.Namespace, "agent-conf", "").Return(config, nil).Times(1)
	sSecret.EXPECT().Get(appView.Namespace, "secret01", "").Return(secret, nil).Times(1)
	sSecret.EXPECT().Get(appView.Namespace, "registry01", "").Return(secret, nil).Times(1)
	sSecret.EXPECT().Get(appView.Namespace, "certificate01", "").Return(secret, nil).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "abc", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	fApp.EXPECT().CreateApp(appView.Namespace, nil, app, gomock.Any()).Return(app1, nil)
	sSecret.EXPECT().Get(appView.Namespace, "secret01", "").Return(secret, nil).Times(1)
	sSecret.EXPECT().Get(appView.Namespace, "registry01", "").Return(secretRegistry, nil).Times(1)
	sSecret.EXPECT().Get(appView.Namespace, "certificate01", "").Return(secretCertificate, nil).Times(1)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(appView)
	req, _ := http.NewRequest(http.MethodPost, "/v1/apps", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	appViewRes := &models.ApplicationView{}
	err := json.Unmarshal(w.Body.Bytes(), appViewRes)
	assert.NoError(t, err)
	assert.Len(t, appViewRes.Volumes, 3)
	assert.Len(t, appViewRes.Registries, 1)
	assert.Equal(t, appViewRes.Registries, appView.Registries)
	assert.Equal(t, appViewRes.Volumes, appView.Volumes)
}

func TestUpdateContainerApplication(t *testing.T) {
	api, router, mockCtl := initApplicationAPI(t)
	defer mockCtl.Finish()

	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sSecret := ms.NewMockSecretService(mockCtl)
	api.AppCombinedService = &service.AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}

	sIndex := ms.NewMockIndexService(mockCtl)
	sNode := ms.NewMockNodeService(mockCtl)
	fApp := mf.NewMockFacade(mockCtl)
	api.Facade = fApp
	api.Index = sIndex
	api.Node = sNode

	mApp := getMockContainerApp()
	mApp.Selector = "label = test"

	config := &specV1.Configuration{Name: "agent-conf", Version: "123"}
	secret1 := &specV1.Secret{Name: "registry01", Version: "123", Labels: map[string]string{specV1.SecretLabel: specV1.SecretRegistry}}
	secret2 := &specV1.Secret{Name: "secret01", Version: "123"}
	registry := &models.Registry{Name: "registry01", Version: "1"}
	sConfig.EXPECT().Get(gomock.Any(), gomock.Any(), "").Return(config, nil).AnyTimes()
	sSecret.EXPECT().Get(gomock.Any(), secret2.Name, gomock.Any()).Return(secret2, nil).AnyTimes()
	sSecret.EXPECT().Get(gomock.Any(), registry.Name, gomock.Any()).Return(secret1, nil).AnyTimes()

	sApp.EXPECT().Get(mApp.Namespace, "abc", gomock.Any()).Return(nil, fmt.Errorf("error")).Times(1)
	w := httptest.NewRecorder()
	body, _ := json.Marshal(mApp)
	req, _ := http.NewRequest(http.MethodPut, "/v1/apps/abc", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mApp2 := getMockContainerApp()
	mApp3 := mApp2
	mApp2.JobConfig = &specV1.AppJobConfig{RestartPolicy: "Never"}
	mApp3.Selector = "name = test"
	mApp3.Mode = context.RunModeKube
	mApp3.Labels = map[string]string{
		common.LabelAppMode: context.RunModeKube,
	}
	mApp3.Services[0].Type = specV1.WorkloadDeployment
	mApp3.Workload = specV1.WorkloadDeployment

	sApp.EXPECT().Get(gomock.Any(), "abc", gomock.Any()).Return(mApp, nil).AnyTimes()
	fApp.EXPECT().UpdateApp(mApp.Namespace, gomock.Any(), gomock.Any(), gomock.Any()).Return(mApp3, nil)
	w = httptest.NewRecorder()
	body, _ = json.Marshal(mApp2)
	req, _ = http.NewRequest(http.MethodPut, "/v1/apps/abc", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	fApp.EXPECT().UpdateApp(mApp.Namespace, gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
	w = httptest.NewRecorder()
	body, _ = json.Marshal(mApp2)
	req, _ = http.NewRequest(http.MethodPut, "/v1/apps/abc", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpdateNativeNotFunctionApplication(t *testing.T) {
	api, router, mockCtl := initApplicationAPI(t)
	defer mockCtl.Finish()

	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sSecret := ms.NewMockSecretService(mockCtl)
	api.AppCombinedService = &service.AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}

	sNode := ms.NewMockNodeService(mockCtl)
	sIndex := ms.NewMockIndexService(mockCtl)
	fApp := mf.NewMockFacade(mockCtl)
	api.Facade = fApp
	api.Index = sIndex
	api.Node = sNode

	appView := &models.ApplicationView{
		Namespace: "baetyl-cloud",
		Name:      "abc",
		Mode:      context.RunModeNative,
		Type:      common.ContainerApp,
		Services: []models.ServiceView{
			{
				Service: specV1.Service{
					Name:     "agent",
					Hostname: "test-agent",
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
				ProgramConfig: "test-program",
			},
		},
		Volumes: []models.VolumeView{
			{
				Name: "name",
				Config: &specV1.ObjectReference{
					Name: "agent-conf",
				},
			},
		},
	}

	programConfig := &specV1.Configuration{
		Name:      "test-program",
		Namespace: "baetyl-cloud",
		Labels: map[string]string{
			"baetyl-config-type": "baetyl-program",
		},
		Data: map[string]string{
			"darwin-amd64": "http://xx/xx.zip",
			"linux-amd64":  "http://xx/xx.zip",
			"linux-arm64":  "http://xx/xx.zip",
			"linux-arm-v7": "http://xx/xx.zip",
		},
	}

	agentConfig := &specV1.Configuration{}
	app1 := &specV1.Application{
		Namespace: "baetyl-cloud",
		Name:      "abc",
		Mode:      context.RunModeNative,
		Type:      common.ContainerApp,
		Services: []specV1.Service{
			{
				Name:     "agent",
				Hostname: "test-agent",
				Replica:  1,
				VolumeMounts: []specV1.VolumeMount{
					{
						Name:      "name",
						MountPath: "mountPath",
					},
					{
						Name:      "baetyl-program-config-agent",
						MountPath: "/var/lib/baetyl/bin",
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
				Name: "baetyl-program-config-agent",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name: "test-program",
					},
				},
			},
		},
	}
	sConfig.EXPECT().Get(appView.Namespace, "test-program", "").Return(programConfig, nil).Times(2)
	sConfig.EXPECT().Get(appView.Namespace, "agent-conf", "").Return(agentConfig, nil).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "abc", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	fApp.EXPECT().CreateApp(appView.Namespace, gomock.Any(), gomock.Any(), gomock.Any()).Return(app1, nil)
	w := httptest.NewRecorder()
	body, _ := json.Marshal(appView)
	req, _ := http.NewRequest(http.MethodPost, "/v1/apps", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	appView12 := &models.ApplicationView{
		Namespace: "baetyl-cloud",
		Name:      "abc",
		Mode:      context.RunModeNative,
		Type:      common.ContainerApp,
		Services: []models.ServiceView{
			{
				Service: specV1.Service{
					Name:     "agent",
					Hostname: "test-agent",
					Replica:  1,
					VolumeMounts: []specV1.VolumeMount{
						{
							Name:      "name",
							MountPath: "mountPath",
						},
						{
							Name:      "baetyl-program-config-agent",
							MountPath: "/var/lib/baetyl/bin",
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
				ProgramConfig: "test-program2",
			},
		},
		Volumes: []models.VolumeView{
			{
				Name: "name",
				Config: &specV1.ObjectReference{
					Name: "agent-conf",
				},
			},
			{
				Name: "baetyl-program-config-agent",
				Config: &specV1.ObjectReference{
					Name: "test-program",
				},
			},
		},
	}

	programConfig2 := &specV1.Configuration{
		Name:      "test-program2",
		Namespace: "baetyl-cloud",
		Labels: map[string]string{
			"baetyl-config-type": "baetyl-program",
		},
		Data: map[string]string{
			"darwin-amd64": "http://xx/xx.zip",
			"linux-amd64":  "http://xx/xx.zip",
			"linux-arm64":  "http://xx/xx.zip",
			"linux-arm-v7": "http://xx/xx.zip",
		},
	}

	app12 := &specV1.Application{
		Namespace: "baetyl-cloud",
		Name:      "abc",
		Mode:      context.RunModeNative,
		Type:      common.ContainerApp,
		Services: []specV1.Service{
			{
				Name:     "agent",
				Hostname: "test-agent",
				Replica:  1,
				VolumeMounts: []specV1.VolumeMount{
					{
						Name:      "name",
						MountPath: "mountPath",
					},
					{
						Name:      "baetyl-program-config-agent",
						MountPath: "/var/lib/baetyl/bin",
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
				Name: "baetyl-program-config-agent",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name: "test-program2",
					},
				},
			},
		},
	}
	sConfig.EXPECT().Get(appView.Namespace, "test-program2", "").Return(programConfig2, nil).Times(2)
	sConfig.EXPECT().Get(appView.Namespace, "agent-conf", "").Return(agentConfig, nil).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "abc", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	fApp.EXPECT().CreateApp(appView.Namespace, gomock.Any(), gomock.Any(), gomock.Any()).Return(app12, nil)
	w = httptest.NewRecorder()
	body, _ = json.Marshal(appView12)
	req, _ = http.NewRequest(http.MethodPost, "/v1/apps", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	appView2 := &models.ApplicationView{
		Namespace: "baetyl-cloud",
		Name:      "abc",
		Mode:      context.RunModeNative,
		Type:      common.ContainerApp,
		Services: []models.ServiceView{
			{
				Service: specV1.Service{
					Name:     "agent2",
					Hostname: "test-agent2",
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
				ProgramConfig: "test-program2",
			},
		},
		Volumes: []models.VolumeView{
			{
				Name: "name",
				Config: &specV1.ObjectReference{
					Name: "agent-conf2",
				},
			},
		},
	}

	app2 := &specV1.Application{
		Namespace: "baetyl-cloud",
		Name:      "abc",
		Mode:      context.RunModeNative,
		Type:      common.ContainerApp,
		Services: []specV1.Service{
			{
				Name:     "agent2",
				Hostname: "test-agent2",
				Replica:  1,
				VolumeMounts: []specV1.VolumeMount{
					{
						Name:      "name",
						MountPath: "mountPath",
					},
					{
						Name:      "baetyl-program-config-agent2",
						MountPath: "/var/lib/baetyl/bin",
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
						Name: "agent-conf2",
					},
				},
			},
			{
				Name: "baetyl-program-config-agent2",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name: "test-program2",
					},
				},
			},
		},
	}

	agentConfig2 := &specV1.Configuration{}
	sConfig.EXPECT().Get(appView.Namespace, "agent-conf2", "").Return(agentConfig2, nil).Times(1)
	sConfig.EXPECT().Get(appView.Namespace, "test-program2", "").Return(programConfig2, nil).Times(2)

	sApp.EXPECT().Get(appView.Namespace, "abc", "").Return(app1, nil).Times(1)
	fApp.EXPECT().UpdateApp(appView.Namespace, gomock.Any(), gomock.Any(), gomock.Any()).Return(app2, nil)

	w = httptest.NewRecorder()
	body, _ = json.Marshal(appView2)
	req, _ = http.NewRequest(http.MethodPut, "/v1/apps/abc", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateInvisibleApplication(t *testing.T) {
	api, router, mockCtl := initApplicationAPI(t)
	defer mockCtl.Finish()

	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sSecret := ms.NewMockSecretService(mockCtl)
	api.AppCombinedService = &service.AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}

	sIndex := ms.NewMockIndexService(mockCtl)
	sNode := ms.NewMockNodeService(mockCtl)
	api.Index = sIndex
	api.Node = sNode

	mApp := &specV1.Application{
		Namespace: "baetyl-cloud",
		Name:      "abc",
		Type:      common.ContainerApp,
		Labels: map[string]string{
			common.ResourceInvisible: "true",
		},
		Services: []specV1.Service{},
		Volumes:  []specV1.Volume{},
	}

	sApp.EXPECT().Get(mApp.Namespace, "abc", gomock.Any()).Return(mApp, nil).Times(1)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(mApp)
	req, _ := http.NewRequest(http.MethodPut, "/v1/apps/abc", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateSysApplication(t *testing.T) {
	api, router, mockCtl := initApplicationAPI(t)
	defer mockCtl.Finish()

	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sSecret := ms.NewMockSecretService(mockCtl)
	api.AppCombinedService = &service.AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}

	sIndex := ms.NewMockIndexService(mockCtl)
	sNode := ms.NewMockNodeService(mockCtl)
	api.Index = sIndex
	api.Node = sNode

	mOldApp := &specV1.Application{
		Namespace: "baetyl-cloud",
		Name:      "abc",
		Type:      common.ContainerApp,
		Labels: map[string]string{
			common.LabelSystem: "true",
		},
		Services: []specV1.Service{},
		Volumes:  []specV1.Volume{},
	}

	mApp := &specV1.Application{
		Namespace: "baetyl-cloud",
		Name:      "abc",
		Type:      common.ContainerApp,
		Labels: map[string]string{
			common.LabelSystem: "true",
			"extra":            "true",
		},
		Services: []specV1.Service{},
		Volumes:  []specV1.Volume{},
	}
	sApp.EXPECT().Get(mApp.Namespace, "abc", "").Return(mOldApp, nil).Times(1)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(mApp)
	req, _ := http.NewRequest(http.MethodPut, "/v1/apps/abc", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mOldApp2 := &specV1.Application{
		Namespace: "baetyl-cloud",
		Name:      "abcd",
		Type:      common.ContainerApp,
		Labels: map[string]string{
			common.LabelSystem: "true",
		},
		Services: []specV1.Service{},
		Volumes:  []specV1.Volume{},
		Selector: "a=a",
	}

	mApp2 := &specV1.Application{
		Namespace: "baetyl-cloud",
		Name:      "abcd",
		Type:      common.ContainerApp,
		Labels: map[string]string{
			common.LabelSystem: "true",
		},
		Services: []specV1.Service{},
		Volumes:  []specV1.Volume{},
		Selector: "b=b",
	}
	sApp.EXPECT().Get(mApp.Namespace, "abcd", "").Return(mOldApp2, nil).Times(1)

	w = httptest.NewRecorder()
	body, _ = json.Marshal(mApp2)
	req, _ = http.NewRequest(http.MethodPut, "/v1/apps/abcd", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateFunctionApplication(t *testing.T) {
	api, router, mockCtl := initApplicationAPI(t)
	defer mockCtl.Finish()

	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sSecret := ms.NewMockSecretService(mockCtl)
	api.AppCombinedService = &service.AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}

	sIndex := ms.NewMockIndexService(mockCtl)
	sNode := ms.NewMockNodeService(mockCtl)
	sFunc := ms.NewMockFunctionService(mockCtl)
	sTempalte := ms.NewMockTemplateService(mockCtl)
	fApp := mf.NewMockFacade(mockCtl)
	api.Facade = fApp

	api.App = sApp
	api.Index = sIndex
	api.Secret = sSecret
	api.Config = sConfig
	api.Node = sNode
	api.Func = sFunc
	api.Template = sTempalte

	appView := &models.ApplicationView{
		Namespace: "baetyl-cloud",
		Name:      "abc",
		Type:      common.FunctionApp,
		Services: []models.ServiceView{
			{
				Service: specV1.Service{
					Name:     "agent",
					Hostname: "test-agent",
					Replica:  1,
					Image:    "image",
					VolumeMounts: []specV1.VolumeMount{
						{
							Name:      "baetyl-function-code-agent",
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
		},
		Volumes: []models.VolumeView{
			{
				Name: "baetyl-function-code-agent",
				Config: &specV1.ObjectReference{
					Name: "func1",
				},
			},
		},
	}

	config := &specV1.Configuration{}
	app := &specV1.Application{}
	sConfig.EXPECT().Get(appView.Namespace, "func1", "").Return(config, nil).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "abc", "").Return(app, nil).Times(1)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(appView)
	req, _ := http.NewRequest(http.MethodPost, "/v1/apps?base=eden", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	sConfig.EXPECT().Get(appView.Namespace, "func1", "").Return(nil, fmt.Errorf("config not found")).Times(1)

	w = httptest.NewRecorder()
	body, _ = json.Marshal(appView)
	req, _ = http.NewRequest(http.MethodPost, "/v1/apps?base=eden", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	sConfig.EXPECT().Get(appView.Namespace, "func1", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	w = httptest.NewRecorder()
	body, _ = json.Marshal(appView)
	req, _ = http.NewRequest(http.MethodPost, "/v1/apps?base=eden", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	sConfig.EXPECT().Get(appView.Namespace, "func1", "").Return(config, nil).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "abc", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "eden", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
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

	sConfig.EXPECT().Get(appView.Namespace, "func1", "").Return(config, nil).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "abc", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "eden2", "").Return(eden2, nil).Times(1)

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

	sConfig.EXPECT().Get(appView.Namespace, "func1", "").Return(config, nil).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "abc", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "eden2", "").Return(eden2, nil).Times(1)
	config2 := &specV1.Configuration{}
	sTempalte.EXPECT().UnmarshalTemplate("baetyl-python36-program.yml", gomock.Any(), config2).Return(errors.New("err")).Times(1)

	w = httptest.NewRecorder()
	body, _ = json.Marshal(appView)
	req, _ = http.NewRequest(http.MethodPost, "/v1/apps?base=eden2", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	sConfig.EXPECT().Get(appView.Namespace, "func1", "").Return(config, nil).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "abc", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "eden2", "").Return(eden2, nil).Times(1)
	sTempalte.EXPECT().UnmarshalTemplate("baetyl-python36-program.yml", gomock.Any(), config2).Return(nil).Times(1)
	sFunc.EXPECT().ListRuntimes().Return(nil, errors.New("err")).Times(1)

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
		Name: fmt.Sprintf("baetyl-function-%s-%s-%s", app.Name, "agent", common.RandString(9)),
		Labels: map[string]string{
			common.LabelSystem: "true",
		},
		Data: map[string]string{
			"service.yml": string(data),
		},
	}
	sConfig.EXPECT().Get(appView.Namespace, "func1", "").Return(config, nil).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "abc", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "eden2", "").Return(eden2, nil).Times(1)
	sTempalte.EXPECT().UnmarshalTemplate("baetyl-python36-program.yml", gomock.Any(), config2).Return(nil).Times(1)
	funcs := map[string]string{
		"python36": "image",
	}
	sFunc.EXPECT().ListRuntimes().Return(funcs, nil).Times(1)
	fApp.EXPECT().CreateApp(appView.Namespace, gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("err")).Times(1)

	w = httptest.NewRecorder()
	body, _ = json.Marshal(appView)
	req, _ = http.NewRequest(http.MethodPost, "/v1/apps?base=eden2", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	sConfig.EXPECT().Get(appView.Namespace, "func1", "").Return(config, nil).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "abc", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "eden2", "").Return(eden2, nil).Times(1)
	sTempalte.EXPECT().UnmarshalTemplate("baetyl-python36-program.yml", gomock.Any(), config2).Return(nil).Times(1)
	sFunc.EXPECT().ListRuntimes().Return(funcs, nil).Times(1)
	// one more for program config
	fApp.EXPECT().CreateApp(appView.Namespace, gomock.Any(), gomock.Any(), gomock.Any()).Return(eden2, nil).Times(1)
	sConfig.EXPECT().Get(appView.Namespace, gomock.Any(), "").Return(config, nil).AnyTimes()

	w = httptest.NewRecorder()
	body, _ = json.Marshal(appView)
	req, _ = http.NewRequest(http.MethodPost, "/v1/apps?base=eden2", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateFunctionApplication(t *testing.T) {
	api, router, mockCtl := initApplicationAPI(t)
	defer mockCtl.Finish()

	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sSecret := ms.NewMockSecretService(mockCtl)
	api.AppCombinedService = &service.AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}

	sTempalte := ms.NewMockTemplateService(mockCtl)
	sIndex := ms.NewMockIndexService(mockCtl)
	sNode := ms.NewMockNodeService(mockCtl)
	sFunc := ms.NewMockFunctionService(mockCtl)
	fApp := mf.NewMockFacade(mockCtl)
	api.Facade = fApp
	api.Index = sIndex
	api.Node = sNode
	api.Func = sFunc
	api.Template = sTempalte

	namespace := "baetyl-cloud"
	oldApp := &specV1.Application{
		Namespace: namespace,
		Name:      "abc",
		Type:      common.FunctionApp,
		Services: []specV1.Service{
			{
				Name:     "agent",
				Hostname: "test-agent",
				Replica:  1,
				Image:    "image",
				VolumeMounts: []specV1.VolumeMount{
					{
						Name:      "baetyl-function-code-agent",
						MountPath: "/var/lib/baetyl/code",
					},
					{
						Name:      "baetyl-function-config-agent",
						MountPath: "/etc/baetyl",
					},
					{
						Name:      "baetyl-function-program-config-agent",
						MountPath: "/var/lib/baetyl/bin",
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
				Name:     "agent2",
				Hostname: "test-agent2",
				Replica:  1,
				Image:    "image2",
				VolumeMounts: []specV1.VolumeMount{
					{
						Name:      "baetyl-function-code-agent2",
						MountPath: "/var/lib/baetyl/code",
					},
					{
						Name:      "baetyl-function-config-agent2",
						MountPath: "/etc/baetyl",
					},
					{
						Name:      "baetyl-function-program-config-agent2",
						MountPath: "/var/lib/baetyl/bin",
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
				Name: "baetyl-funciton-code-agent",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name: "func1",
					},
				},
			},
			{
				Name: "baetyl-function-code-agent2",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name: "func2",
					},
				},
			},
			{
				Name: "baetyl-function-config-agent",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name: "baetyl-function-config-app-service-1",
					},
				},
			},
			{
				Name: "baetyl-function-config-agent2",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name: "baetyl-function-config-app-service-2",
					},
				},
			},
			{
				Name: "baetyl-function-program-config-agent",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name: "baetyl-function-program-config-app-service-aaaa",
					},
				},
			},
			{
				Name: "baetyl-function-program-config-agent2",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name: "baetyl-function-program-config-app-service-bbbb",
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
		Name: fmt.Sprintf("baetyl-function-config-%s-%s-%s", oldApp.Name, "agent", common.RandString(9)),
		Labels: map[string]string{
			common.LabelSystem: "true",
		},
		Data: map[string]string{
			"service.yml": string(data2),
		},
	}
	assert.NoError(t, err)
	assert.NoError(t, err)

	newApp := &specV1.Application{
		Namespace: namespace,
		Name:      "abc",
		Type:      common.FunctionApp,
		Services: []specV1.Service{
			{
				Name:     "agent2",
				Hostname: "test-agent2",
				Replica:  1,
				Image:    "image2",
				VolumeMounts: []specV1.VolumeMount{
					{
						Name:      "baetyl-function-code-agent2",
						MountPath: "/var/lib/baetyl/code",
					},
					{
						Name:      "baetyl-function-config-agent2",
						MountPath: "/etc/baetyl",
					},
					{
						Name:      "baetyl-function-program-config-agent2",
						MountPath: "/var/lib/baetyl/bin",
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
				Name:     "agent3",
				Hostname: "test-agent3",
				Replica:  1,
				Image:    "image2",
				VolumeMounts: []specV1.VolumeMount{
					{
						Name:      "baetyl-function-code-agent3",
						MountPath: "/var/lib/baetyl/code",
					},
					{
						Name:      "baetyl-function-config-agent3",
						MountPath: "/etc/baetyl",
					},
					{
						Name:      "baetyl-function-program-config-agent3",
						MountPath: "/var/lib/baetyl/bin",
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
				Name: "baetyl-function-code-agent2",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name: "func2",
					},
				},
			},
			{
				Name: "baetyl-function-code-agent3",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name: "func3",
					},
				},
			},
			{
				Name: "baetyl-function-config-agent2",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name: "baetyl-function-config-app-service-2",
					},
				},
			},
			{
				Name: "baetyl-function-config-agent3",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name: "baetyl-function-config-app-service-3",
					},
				},
			},
			{
				Name: "baetyl-function-program-config-agent2",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name: "baetyl-function-program-config-app-service-bbbb",
					},
				},
			},
			{
				Name: "baetyl-function-program-config-agent3",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name: "baetyl-function-program-config-app-service-cccc",
					},
				},
			},
		},
	}
	newAppView := &models.ApplicationView{
		Namespace: namespace,
		Type:      common.FunctionApp,
		Services: []models.ServiceView{
			{
				Service: specV1.Service{
					Name:     "agent2",
					Hostname: "test-agent2",
					Replica:  1,
					Image:    "image2",
					VolumeMounts: []specV1.VolumeMount{
						{
							Name:      "baetyl-funciton-code-agent2",
							MountPath: "/var/lib/baetyl/code",
						},
						{
							Name:      "baetyl-function-config-agent2",
							MountPath: "/etc/baetyl",
						},
						{
							Name:      "baetyl-function-program-config-agent2",
							MountPath: "/var/lib/baetyl/bin",
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
			},
			{
				Service: specV1.Service{
					Name:     "agent3",
					Hostname: "test-agent3",
					Replica:  1,
					Image:    "image3",
					VolumeMounts: []specV1.VolumeMount{
						{
							Name:      "baetyl-function-code-agent3",
							MountPath: "/var/lib/baetyl/code",
						},
						{
							Name:      "baetyl-function-config-agent3",
							MountPath: "/etc/baetyl",
						},
						{
							Name:      "baetyl-function-program-config-agent3",
							MountPath: "/var/lib/baetyl/bin",
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
		},
		Volumes: []models.VolumeView{
			{
				Name: "baetyl-function-code-agent2",
				Config: &specV1.ObjectReference{
					Name: "func2",
				},
			},
			{
				Name: "baetyl-function-code-agent3",
				Config: &specV1.ObjectReference{
					Name: "func3",
				},
			},
			{
				Name: "baetyl-function-config-agent2",
				Config: &specV1.ObjectReference{
					Name: "baetyl-function-config-app-service-2",
				},
			},
			{
				Name: "baetyl-function-config-agent3",
				Config: &specV1.ObjectReference{
					Name: "baetyl-function-config-app-service-3",
				},
			},
			{
				Name: "baetyl-function-program-config-agent2",
				Config: &specV1.ObjectReference{
					Name: "baetyl-function-program-config-app-service-bbbb",
				},
			},
			{
				Name: "baetyl-function-program-config-agent3",
				Config: &specV1.ObjectReference{
					Name: "baetyl-function-program-config-app-service-cccc",
				},
			},
		},
	}

	configCode := &specV1.Configuration{}
	sApp.EXPECT().Get(namespace, "abc", "").Return(oldApp, nil).Times(1)
	sConfig.EXPECT().Get(namespace, "func2", "").Return(configCode, nil).Times(1)
	sConfig.EXPECT().Get(namespace, "func3", "").Return(configCode, nil).Times(1)
	sConfig.EXPECT().Get(namespace, "baetyl-function-config-app-service-2", "").Return(config2, nil).Times(1)
	sConfig.EXPECT().Get(namespace, "baetyl-function-config-app-service-3", "").Return(config2, nil).Times(1)
	sConfig.EXPECT().Get(namespace, "baetyl-function-program-config-app-service-bbbb", "").Return(config2, nil).Times(1)
	sConfig.EXPECT().Get(namespace, "baetyl-function-program-config-app-service-cccc", "").Return(config2, nil).Times(1)
	config4 := &specV1.Configuration{}
	sTempalte.EXPECT().UnmarshalTemplate("baetyl-python36-program.yml", gomock.Any(), config4).Return(nil).Times(2)
	funcs := map[string]string{
		"python36": "image",
	}
	sFunc.EXPECT().ListRuntimes().Return(funcs, nil).Times(2)
	fApp.EXPECT().UpdateApp(namespace, gomock.Any(), gomock.Any(), gomock.Any()).Return(newApp, nil)
	sConfig.EXPECT().Get(namespace, "baetyl-function-config-app-service-2", "").Return(config2, nil).Times(1)
	sConfig.EXPECT().Get(namespace, "baetyl-function-config-app-service-3", "").Return(config2, nil).Times(1)
	sConfig.EXPECT().Get(namespace, "baetyl-function-program-config-app-service-bbbb", "").Return(config2, nil).Times(1)
	sConfig.EXPECT().Get(namespace, "baetyl-function-program-config-app-service-cccc", "").Return(config2, nil).Times(1)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(newAppView)
	req, _ := http.NewRequest(http.MethodPut, "/v1/apps/abc", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteApplication(t *testing.T) {
	api, router, mockCtl := initApplicationAPI(t)
	defer mockCtl.Finish()

	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sSecret := ms.NewMockSecretService(mockCtl)
	api.AppCombinedService = &service.AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}

	sIndex := ms.NewMockIndexService(mockCtl)
	sNode := ms.NewMockNodeService(mockCtl)
	fApp := mf.NewMockFacade(mockCtl)
	api.Facade = fApp
	api.Index = sIndex
	api.Node = sNode

	app := &specV1.Application{
		Namespace: "baetyl-cloud",
		Name:      "abc",
	}

	// 500
	sApp.EXPECT().Get(gomock.Any(), "abc", gomock.Any()).Return(app, nil).Times(1)
	fApp.EXPECT().DeleteApp(app.Namespace, app.Name, gomock.Any()).Return(fmt.Errorf("error")).Times(1)
	req, _ := http.NewRequest(http.MethodDelete, "/v1/apps/abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// 200
	sApp.EXPECT().Get(gomock.Any(), "abc", gomock.Any()).Return(app, nil).Times(1)
	fApp.EXPECT().DeleteApp(app.Namespace, app.Name, gomock.Any()).Return(nil).Times(1)
	req, _ = http.NewRequest(http.MethodDelete, "/v1/apps/abc", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 200 delete non-existent app
	sApp.EXPECT().Get(gomock.Any(), "abc", gomock.Any()).Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	req, _ = http.NewRequest(http.MethodDelete, "/v1/apps/abc", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	app.Name = "baetyl-core-test"
	sApp.EXPECT().Get(gomock.Any(), app.Name, gomock.Any()).Return(app, nil).AnyTimes()
	sIndex.EXPECT().ListNodesByApp(gomock.Any(), app.Name).Return(nil, fmt.Errorf("error"))
	req, _ = http.NewRequest(http.MethodDelete, "/v1/apps/baetyl-core-test", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	sIndex.EXPECT().ListNodesByApp(gomock.Any(), app.Name).Return([]string{"test-node"}, nil)
	req, _ = http.NewRequest(http.MethodDelete, "/v1/apps/baetyl-core-test", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateKubeFunctionApplication(t *testing.T) {
	api, router, mockCtl := initApplicationAPI(t)
	defer mockCtl.Finish()

	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sSecret := ms.NewMockSecretService(mockCtl)
	api.AppCombinedService = &service.AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}

	sIndex := ms.NewMockIndexService(mockCtl)
	sNode := ms.NewMockNodeService(mockCtl)
	sFunc := ms.NewMockFunctionService(mockCtl)
	sTemplate := ms.NewMockTemplateService(mockCtl)
	fApp := mf.NewMockFacade(mockCtl)
	api.Facade = fApp

	api.App = sApp
	api.Index = sIndex
	api.Secret = sSecret
	api.Config = sConfig
	api.Node = sNode
	api.Func = sFunc
	api.Template = sTemplate

	appView := &models.ApplicationView{
		Namespace: "baetyl-cloud",
		Name:      "abc",
		Mode:      context.RunModeKube,
		Type:      common.FunctionApp,
		Services: []models.ServiceView{
			{
				Service: specV1.Service{
					Name:     "agent",
					Hostname: "test-agent",
					Replica:  1,
					Image:    "image",
					VolumeMounts: []specV1.VolumeMount{
						{
							Name:      "baetyl-function-code-agent",
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
		},
		Volumes: []models.VolumeView{
			{
				Name: "baetyl-function-code-agent",
				Config: &specV1.ObjectReference{
					Name: "func1",
				},
			},
		},
	}

	config := &specV1.Configuration{}
	app := &specV1.Application{}
	sConfig.EXPECT().Get(appView.Namespace, "func1", "").Return(config, nil).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "abc", "").Return(app, nil).Times(1)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(appView)
	req, _ := http.NewRequest(http.MethodPost, "/v1/apps?base=eden", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	sConfig.EXPECT().Get(appView.Namespace, "func1", "").Return(nil, fmt.Errorf("config not found")).Times(1)

	w = httptest.NewRecorder()
	body, _ = json.Marshal(appView)
	req, _ = http.NewRequest(http.MethodPost, "/v1/apps?base=eden", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	sConfig.EXPECT().Get(appView.Namespace, "func1", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	w = httptest.NewRecorder()
	body, _ = json.Marshal(appView)
	req, _ = http.NewRequest(http.MethodPost, "/v1/apps?base=eden", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	sConfig.EXPECT().Get(appView.Namespace, "func1", "").Return(config, nil).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "abc", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "eden", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	w = httptest.NewRecorder()
	body, _ = json.Marshal(appView)
	req, _ = http.NewRequest(http.MethodPost, "/v1/apps?base=eden", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	eden2 := &specV1.Application{
		Namespace: appView.Namespace,
		Name:      "abc",
		Mode:      context.RunModeKube,
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

	sConfig.EXPECT().Get(appView.Namespace, "func1", "").Return(config, nil).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "abc", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "eden2", "").Return(eden2, nil).Times(1)

	w = httptest.NewRecorder()
	body, _ = json.Marshal(appView)
	req, _ = http.NewRequest(http.MethodPost, "/v1/apps?base=eden2", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	eden2 = &specV1.Application{
		Namespace: appView.Namespace,
		Name:      "abc",
		Type:      common.FunctionApp,
		Mode:      context.RunModeKube,
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

	sConfig.EXPECT().Get(appView.Namespace, "func1", "").Return(config, nil).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "abc", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "eden2", "").Return(eden2, nil).Times(1)
	config2 := &specV1.Configuration{}
	sTemplate.EXPECT().UnmarshalTemplate("baetyl-python36-program.yml", gomock.Any(), config2).Return(errors.New("err")).Times(1)

	w = httptest.NewRecorder()
	body, _ = json.Marshal(appView)
	req, _ = http.NewRequest(http.MethodPost, "/v1/apps?base=eden2", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	sConfig.EXPECT().Get(appView.Namespace, "func1", "").Return(config, nil).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "abc", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "eden2", "").Return(eden2, nil).Times(1)
	sTemplate.EXPECT().UnmarshalTemplate("baetyl-python36-program.yml", gomock.Any(), config2).Return(nil).Times(1)
	sFunc.EXPECT().ListRuntimes().Return(nil, errors.New("err")).Times(1)

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
		Name: fmt.Sprintf("baetyl-function-%s-%s-%s", app.Name, "agent", common.RandString(9)),
		Labels: map[string]string{
			common.LabelSystem: "true",
		},
		Data: map[string]string{
			"service.yml": string(data),
		},
	}
	sConfig.EXPECT().Get(appView.Namespace, "func1", "").Return(config, nil).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "abc", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "eden2", "").Return(eden2, nil).Times(1)
	sTemplate.EXPECT().UnmarshalTemplate("baetyl-python36-program.yml", gomock.Any(), config2).Return(nil).Times(1)
	funcs := map[string]string{
		"python36": "image",
	}
	sFunc.EXPECT().ListRuntimes().Return(funcs, nil).Times(1)
	fApp.EXPECT().CreateApp(appView.Namespace, gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("err")).Times(1)

	w = httptest.NewRecorder()
	body, _ = json.Marshal(appView)
	req, _ = http.NewRequest(http.MethodPost, "/v1/apps?base=eden2", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	sConfig.EXPECT().Get(appView.Namespace, "func1", "").Return(config, nil).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "abc", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "eden2", "").Return(eden2, nil).Times(1)
	sTemplate.EXPECT().UnmarshalTemplate("baetyl-python36-program.yml", gomock.Any(), config2).Return(nil).Times(1)
	sFunc.EXPECT().ListRuntimes().Return(funcs, nil).Times(1)
	// one more for program config
	fApp.EXPECT().CreateApp(appView.Namespace, gomock.Any(), gomock.Any(), gomock.Any()).Return(eden2, nil).Times(1)
	sConfig.EXPECT().Get(appView.Namespace, gomock.Any(), "").Return(config, nil).AnyTimes()

	w = httptest.NewRecorder()
	body, _ = json.Marshal(appView)
	req, _ = http.NewRequest(http.MethodPost, "/v1/apps?base=eden2", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateKubeNotFunctionApplication(t *testing.T) {
	api, router, mockCtl := initApplicationAPI(t)
	defer mockCtl.Finish()

	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sSecret := ms.NewMockSecretService(mockCtl)
	api.AppCombinedService = &service.AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}

	sNode := ms.NewMockNodeService(mockCtl)
	sIndex := ms.NewMockIndexService(mockCtl)
	fApp := mf.NewMockFacade(mockCtl)
	api.Facade = fApp
	api.Index = sIndex
	api.Node = sNode

	appView := &models.ApplicationView{
		Namespace: "baetyl-cloud",
		Name:      "abc",
		Mode:      context.RunModeKube,
		Type:      common.ContainerApp,
		Services: []models.ServiceView{
			{
				Service: specV1.Service{
					Name:     "agent",
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
		},
		Volumes: []models.VolumeView{
			{
				Name: "name",
				Config: &specV1.ObjectReference{
					Name: "agent-conf",
				},
			},
			{
				Name: "secret",
				Secret: &specV1.ObjectReference{
					Name: "secret01",
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

	eden2 := &specV1.Application{
		Namespace: appView.Namespace,
		Name:      "abc",
		Mode:      context.RunModeKube,
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

	sConfig.EXPECT().Get(appView.Namespace, "agent-conf", "").Return(config, nil).Times(1)
	sSecret.EXPECT().Get(appView.Namespace, "secret01", "").Return(secret, nil).Times(1)
	sSecret.EXPECT().Get(appView.Namespace, "registry01", "").Return(secret, nil).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "abc", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "eden2", "").Return(eden2, nil).Return(eden2, nil).Times(1)
	fApp.EXPECT().CreateApp(appView.Namespace, gomock.Any(), gomock.Any(), gomock.Any()).Return(eden2, nil).Times(1)
	sSecret.EXPECT().Get(appView.Namespace, "secret01", "").Return(secret, nil).Times(1)
	w := httptest.NewRecorder()
	body, _ := json.Marshal(appView)
	req, _ := http.NewRequest(http.MethodPost, "/v1/apps?base=eden2", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateNativeFunctionApplication(t *testing.T) {
	api, router, mockCtl := initApplicationAPI(t)
	defer mockCtl.Finish()

	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sSecret := ms.NewMockSecretService(mockCtl)
	api.AppCombinedService = &service.AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}

	sIndex := ms.NewMockIndexService(mockCtl)
	sNode := ms.NewMockNodeService(mockCtl)
	sFunc := ms.NewMockFunctionService(mockCtl)
	sTempalte := ms.NewMockTemplateService(mockCtl)
	fApp := mf.NewMockFacade(mockCtl)
	api.Facade = fApp

	api.App = sApp
	api.Index = sIndex
	api.Secret = sSecret
	api.Config = sConfig
	api.Node = sNode
	api.Func = sFunc
	api.Template = sTempalte

	appView := &models.ApplicationView{
		Namespace: "baetyl-cloud",
		Name:      "abc",
		Type:      common.FunctionApp,
		Services: []models.ServiceView{
			{
				Service: specV1.Service{
					Name:     "agent",
					Hostname: "test-agent",
					Replica:  1,
					Image:    "image",
					VolumeMounts: []specV1.VolumeMount{
						{
							Name:      "baetyl-function-code-agent",
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
		},
		Volumes: []models.VolumeView{
			{
				Name: "baetyl-function-code-agent",
				Config: &specV1.ObjectReference{
					Name: "func1",
				},
			},
		},
	}

	config := &specV1.Configuration{}
	app := &specV1.Application{}
	sConfig.EXPECT().Get(appView.Namespace, "func1", "").Return(config, nil).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "abc", "").Return(app, nil).Times(1)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(appView)
	req, _ := http.NewRequest(http.MethodPost, "/v1/apps?base=eden", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	sConfig.EXPECT().Get(appView.Namespace, "func1", "").Return(nil, fmt.Errorf("config not found")).Times(1)

	w = httptest.NewRecorder()
	body, _ = json.Marshal(appView)
	req, _ = http.NewRequest(http.MethodPost, "/v1/apps?base=eden", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	sConfig.EXPECT().Get(appView.Namespace, "func1", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	w = httptest.NewRecorder()
	body, _ = json.Marshal(appView)
	req, _ = http.NewRequest(http.MethodPost, "/v1/apps?base=eden", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	sConfig.EXPECT().Get(appView.Namespace, "func1", "").Return(config, nil).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "abc", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "eden", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
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

	sConfig.EXPECT().Get(appView.Namespace, "func1", "").Return(config, nil).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "abc", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "eden2", "").Return(eden2, nil).Times(1)

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

	sConfig.EXPECT().Get(appView.Namespace, "func1", "").Return(config, nil).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "abc", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "eden2", "").Return(eden2, nil).Times(1)
	config2 := &specV1.Configuration{}
	sTempalte.EXPECT().UnmarshalTemplate("baetyl-python36-program.yml", gomock.Any(), config2).Return(errors.New("err")).Times(1)

	w = httptest.NewRecorder()
	body, _ = json.Marshal(appView)
	req, _ = http.NewRequest(http.MethodPost, "/v1/apps?base=eden2", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	sConfig.EXPECT().Get(appView.Namespace, "func1", "").Return(config, nil).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "abc", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "eden2", "").Return(eden2, nil).Times(1)
	sTempalte.EXPECT().UnmarshalTemplate("baetyl-python36-program.yml", gomock.Any(), config2).Return(nil).Times(1)
	sFunc.EXPECT().ListRuntimes().Return(nil, errors.New("err")).Times(1)

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
		Name: fmt.Sprintf("baetyl-function-%s-%s-%s", app.Name, "agent", common.RandString(9)),
		Labels: map[string]string{
			common.LabelSystem: "true",
		},
		Data: map[string]string{
			"service.yml": string(data),
		},
	}
	sConfig.EXPECT().Get(appView.Namespace, "func1", "").Return(config, nil).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "abc", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "eden2", "").Return(eden2, nil).Times(1)
	sTempalte.EXPECT().UnmarshalTemplate("baetyl-python36-program.yml", gomock.Any(), config2).Return(nil).Times(1)
	funcs := map[string]string{
		"python36": "image",
	}
	sFunc.EXPECT().ListRuntimes().Return(funcs, nil).Times(1)
	fApp.EXPECT().CreateApp(appView.Namespace, gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("err")).Times(1)

	w = httptest.NewRecorder()
	body, _ = json.Marshal(appView)
	req, _ = http.NewRequest(http.MethodPost, "/v1/apps?base=eden2", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	sConfig.EXPECT().Get(appView.Namespace, "func1", "").Return(config, nil).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "abc", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "eden2", "").Return(eden2, nil).Times(1)
	sTempalte.EXPECT().UnmarshalTemplate("baetyl-python36-program.yml", gomock.Any(), config2).Return(nil).Times(1)
	sFunc.EXPECT().ListRuntimes().Return(funcs, nil).Times(1)
	// one more for program config
	fApp.EXPECT().CreateApp(appView.Namespace, gomock.Any(), gomock.Any(), gomock.Any()).Return(eden2, nil).Times(1)
	sConfig.EXPECT().Get(appView.Namespace, gomock.Any(), "").Return(config, nil).AnyTimes()

	w = httptest.NewRecorder()
	body, _ = json.Marshal(appView)
	req, _ = http.NewRequest(http.MethodPost, "/v1/apps?base=eden2", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateNativeNotFunctionApplication(t *testing.T) {
	api, router, mockCtl := initApplicationAPI(t)
	defer mockCtl.Finish()

	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sSecret := ms.NewMockSecretService(mockCtl)
	api.AppCombinedService = &service.AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}

	sNode := ms.NewMockNodeService(mockCtl)
	sIndex := ms.NewMockIndexService(mockCtl)
	fApp := mf.NewMockFacade(mockCtl)
	api.Facade = fApp
	api.Index = sIndex
	api.Node = sNode

	appView := &models.ApplicationView{
		Namespace: "baetyl-cloud",
		Name:      "abc",
		Mode:      context.RunModeNative,
		Type:      common.ContainerApp,
		Services: []models.ServiceView{
			{
				Service: specV1.Service{
					Name:     "agent",
					Hostname: "test-agent",
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
				ProgramConfig: "test-program",
			},
		},
		Volumes: []models.VolumeView{
			{
				Name: "name",
				Config: &specV1.ObjectReference{
					Name: "agent-conf",
				},
			},
			{
				Name: "secret",
				Secret: &specV1.ObjectReference{
					Name: "secret01",
				},
			},
		},
		Registries: []models.RegistryView{
			{
				Name: "registry01",
			},
		},
	}

	programConfig := &specV1.Configuration{
		Name:      "test-program",
		Namespace: "baetyl-cloud",
		Labels: map[string]string{
			"baetyl-config-type": "baetyl-program",
		},
		Data: map[string]string{
			"darwin-amd64": "http://xx/xx.zip",
			"linux-amd64":  "http://xx/xx.zip",
			"linux-arm64":  "http://xx/xx.zip",
			"linux-arm-v7": "http://xx/xx.zip",
		},
	}

	secret := &specV1.Secret{}

	eden2 := &specV1.Application{
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

	agentConfig := &specV1.Configuration{}
	sConfig.EXPECT().Get(appView.Namespace, "agent-conf", "").Return(agentConfig, nil).Times(1)
	sConfig.EXPECT().Get(appView.Namespace, "test-program", "").Return(programConfig, nil).Times(1)
	sSecret.EXPECT().Get(appView.Namespace, "secret01", "").Return(secret, nil).Times(1)
	sSecret.EXPECT().Get(appView.Namespace, "registry01", "").Return(secret, nil).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "abc", "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	sApp.EXPECT().Get(appView.Namespace, "eden2", "").Return(eden2, nil).Return(eden2, nil).Times(1)
	fApp.EXPECT().CreateApp(appView.Namespace, gomock.Any(), gomock.Any(), gomock.Any()).Return(eden2, nil).Times(1)
	sSecret.EXPECT().Get(appView.Namespace, "secret01", "").Return(secret, nil).Times(1)
	w := httptest.NewRecorder()
	body, _ := json.Marshal(appView)
	req, _ := http.NewRequest(http.MethodPost, "/v1/apps?base=eden2", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAPI_GetSysAppConfigs(t *testing.T) {
	api, router, mockCtl := initApplicationAPI(t)
	defer mockCtl.Finish()

	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sSecret := ms.NewMockSecretService(mockCtl)
	api.AppCombinedService = &service.AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}

	sIndex := ms.NewMockIndexService(mockCtl)
	sNode := ms.NewMockNodeService(mockCtl)
	api.Index = sIndex
	api.Node = sNode

	mClist := &models.ConfigurationList{}
	mApp := &specV1.Application{}
	sApp.EXPECT().Get("baetyl-cloud", gomock.Any(), "").Return(mApp, nil)
	sConfig.EXPECT().List("baetyl-cloud", gomock.Any()).Return(mClist, nil).Times(1)

	// 200
	req, _ := http.NewRequest(http.MethodGet, "/v1/apps/test/configs", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAPI_GetNodeSysAppSecrets(t *testing.T) {
	api, router, mockCtl := initApplicationAPI(t)
	defer mockCtl.Finish()

	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sSecret := ms.NewMockSecretService(mockCtl)
	api.AppCombinedService = &service.AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}

	sIndex := ms.NewMockIndexService(mockCtl)
	sNode := ms.NewMockNodeService(mockCtl)
	api.Index = sIndex
	api.Node = sNode

	mClist := &models.SecretList{}
	mApp := &specV1.Application{}
	sApp.EXPECT().Get("baetyl-cloud", gomock.Any(), "").Return(mApp, nil)
	sSecret.EXPECT().List("baetyl-cloud", gomock.Any()).Return(mClist, nil).Times(1)

	// 200
	req, _ := http.NewRequest(http.MethodGet, "/v1/apps/test/secrets", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAPI_GetNodeSysAppCertificates(t *testing.T) {
	api, router, mockCtl := initApplicationAPI(t)
	defer mockCtl.Finish()

	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sSecret := ms.NewMockSecretService(mockCtl)
	api.AppCombinedService = &service.AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}

	sIndex := ms.NewMockIndexService(mockCtl)
	sNode := ms.NewMockNodeService(mockCtl)
	api.Index = sIndex
	api.Node = sNode

	mClist := &models.SecretList{}
	mApp := &specV1.Application{}
	sApp.EXPECT().Get("baetyl-cloud", gomock.Any(), "").Return(mApp, nil)
	sSecret.EXPECT().List("baetyl-cloud", gomock.Any()).Return(mClist, nil).Times(1)

	// 200
	req, _ := http.NewRequest(http.MethodGet, "/v1/apps/test/certificates", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAPI_GetNodeSysAppRegistries(t *testing.T) {
	api, router, mockCtl := initApplicationAPI(t)
	defer mockCtl.Finish()

	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sSecret := ms.NewMockSecretService(mockCtl)
	api.AppCombinedService = &service.AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}

	sIndex := ms.NewMockIndexService(mockCtl)
	sNode := ms.NewMockNodeService(mockCtl)
	api.Index = sIndex
	api.Node = sNode

	mClist := &models.SecretList{}
	mApp := &specV1.Application{}
	sApp.EXPECT().Get("baetyl-cloud", gomock.Any(), "").Return(mApp, nil)
	sSecret.EXPECT().List("baetyl-cloud", gomock.Any()).Return(mClist, nil).Times(1)

	// 200
	req, _ := http.NewRequest(http.MethodGet, "/v1/apps/test/registries", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func Test_compatibleAppDeprecatedFiled(t *testing.T) {
	api, _, mockCtl := initApplicationAPI(t)
	defer mockCtl.Finish()

	app0 := &models.ApplicationView{
		Name: "a0",
		Services: []models.ServiceView{
			{
				Service: specV1.Service{
					Name:        "s1",
					Labels:      map[string]string{"a": "b"},
					HostNetwork: true,
					Replica:     3,
					JobConfig: &specV1.ServiceJobConfig{
						Completions:   1,
						Parallelism:   2,
						BackoffLimit:  3,
						RestartPolicy: "Never",
					},
					Type: specV1.WorkloadJob,
				},
			},
		},
	}

	expectApp0 := &models.ApplicationView{
		Name:        "a0",
		Labels:      map[string]string{"a": "b"},
		HostNetwork: true,
		Replica:     3,
		JobConfig: &specV1.AppJobConfig{
			Completions:   1,
			Parallelism:   2,
			BackoffLimit:  3,
			RestartPolicy: "Never",
		},
		Workload: specV1.WorkloadJob,
		Services: []models.ServiceView{
			{
				Service: specV1.Service{
					Name:        "s1",
					Labels:      map[string]string{"a": "b"},
					HostNetwork: true,
					Replica:     3,
					JobConfig: &specV1.ServiceJobConfig{
						Completions:   1,
						Parallelism:   2,
						BackoffLimit:  3,
						RestartPolicy: "Never",
					},
					Type: specV1.WorkloadJob,
				},
			},
		},
	}

	api.compatibleAppDeprecatedFiled(app0)
	assert.EqualValues(t, expectApp0, app0)

	// case 1
	app1 := &models.ApplicationView{
		Name:     "a1",
		Services: []models.ServiceView{},
	}

	expectApp1 := &models.ApplicationView{
		Name:        "a1",
		Labels:      map[string]string{},
		HostNetwork: false,
		Replica:     1,
		JobConfig: &specV1.AppJobConfig{
			Completions:   0,
			Parallelism:   0,
			BackoffLimit:  0,
			RestartPolicy: "Never",
		},
		Workload: specV1.WorkloadDeployment,
		Services: []models.ServiceView{},
	}

	api.compatibleAppDeprecatedFiled(app1)
	assert.EqualValues(t, expectApp1, app1)

	// case 2
	app2 := &models.ApplicationView{
		Name:        "a0",
		Replica:     1,
		HostNetwork: true,
		JobConfig: &specV1.AppJobConfig{
			Completions:   1,
			Parallelism:   2,
			BackoffLimit:  3,
			RestartPolicy: "Never",
		},
		Workload: specV1.WorkloadJob,
		Services: []models.ServiceView{
			{
				Service: specV1.Service{
					Name:   "s1",
					Labels: map[string]string{"a": "b"},
				},
			},
		},
	}

	expectApp2 := &models.ApplicationView{
		Name:        "a0",
		Labels:      map[string]string{"a": "b"},
		HostNetwork: true,
		Replica:     1,
		JobConfig: &specV1.AppJobConfig{
			Completions:   1,
			Parallelism:   2,
			BackoffLimit:  3,
			RestartPolicy: "Never",
		},
		Workload: specV1.WorkloadJob,
		Services: []models.ServiceView{
			{
				Service: specV1.Service{
					Name:        "s1",
					Labels:      map[string]string{"a": "b"},
					HostNetwork: true,
					Replica:     1,
					JobConfig: &specV1.ServiceJobConfig{
						Completions:   1,
						Parallelism:   2,
						BackoffLimit:  3,
						RestartPolicy: "Never",
					},
					Type: specV1.WorkloadJob,
				},
			},
		},
	}

	api.compatibleAppDeprecatedFiled(app2)
	assert.EqualValues(t, expectApp2.Services[0].JobConfig, app2.Services[0].JobConfig)
	expectApp2.Services[0].JobConfig = nil
	app2.Services[0].JobConfig = nil
	assert.EqualValues(t, expectApp2, app2)
}
