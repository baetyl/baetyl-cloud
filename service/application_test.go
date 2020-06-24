package service

import (
	"fmt"
	"testing"

	"github.com/baetyl/baetyl-cloud/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"encoding/json"
	ms "github.com/baetyl/baetyl-cloud/mock/service"
	specV1 "github.com/baetyl/baetyl-go/spec/v1"
	"time"
)

func genAppTestCase() (*specV1.Application, *specV1.Application) {
	newApp := &specV1.Application{
		Namespace: "default",
		Name:      "abc",
		Selector:  "test",
		Version:   "2",
		Services: []specV1.Service{
			{
				Name:     "Agent",
				Hostname: "test-agent",
				Image:    "hub.baidubce.com/baetyl/baetyl-agent:1.0.0",
				Replica:  1,
				VolumeMounts: []specV1.VolumeMount{
					{
						Name:      "test",
						MountPath: "mountPath",
					},
				},
			},
		},
		Volumes: []specV1.Volume{
			{
				Name: "test",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name: "agent-conf",
					},
				},
			},
			{
				Name: "test-2",
				VolumeSource: specV1.VolumeSource{
					Secret: &specV1.ObjectReference{
						Name: "test-secret-02",
					},
				},
			},
		},
	}

	baseApp := &specV1.Application{
		Namespace: "default",
		Name:      "abc",
		Selector:  "test",
		Version:   "1",
		Services: []specV1.Service{
			{
				Name:     "Agent-02",
				Hostname: "test-agent",
				Image:    "hub.baidubce.com/baetyl/baetyl-agent:1.0.0",
				Replica:  1,
				VolumeMounts: []specV1.VolumeMount{
					{
						Name:      "test-02",
						MountPath: "mountPath",
					},
				},
			},
		},
		Volumes: []specV1.Volume{
			{
				Name: "test-02",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name:    "agent-conf",
						Version: "version01",
					},
				},
			},
		},
	}
	return newApp, baseApp
}

func TestDefaultApplicationService_Get(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	namespace := "default"
	name := "Deployment-get"

	mockObject.modelStorage.EXPECT().GetApplication(namespace, name, "").Return(nil, nil).AnyTimes()
	cs, err := NewApplicationService(mockObject.conf)
	assert.NoError(t, err)
	_, err = cs.Get(namespace, name, "")
	assert.NoError(t, err)
}

func TestDefaultApplicationService_List(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	namespace := "default"
	selector := &models.ListOptions{
		LabelSelector: "a=a",
	}

	mockObject.modelStorage.EXPECT().ListApplication(namespace, selector).Return(nil, nil).AnyTimes()
	cs, err := NewApplicationService(mockObject.conf)
	assert.NoError(t, err)
	_, err = cs.List(namespace, selector)
	assert.NoError(t, err)
}

func TestDefaultApplicationService_Delete(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	mockIndexService := ms.NewMockIndexService(mockObject.ctl)
	as := applicationService{
		storage:      mockObject.modelStorage,
		indexService: mockIndexService,
		dbStorage:    mockObject.dbStorage,
	}
	newApp, _ := genAppTestCase()

	mockObject.modelStorage.EXPECT().DeleteApplication(gomock.Any(), gomock.Any()).Return(fmt.Errorf("error")).Times(1)
	err := as.Delete(newApp.Namespace, newApp.Name, "")
	assert.NotNil(t, err)

	mockObject.modelStorage.EXPECT().DeleteApplication(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockIndexService.EXPECT().RefreshConfigIndexByApp(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("error"))
	mockIndexService.EXPECT().RefreshSecretIndexByApp(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("error"))
	mockObject.dbStorage.EXPECT().DeleteApplication(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
	err = as.Delete(newApp.Namespace, newApp.Name, "")
	assert.NoError(t, err)

	mockIndexService.EXPECT().RefreshConfigIndexByApp(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	mockIndexService.EXPECT().RefreshSecretIndexByApp(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	mockObject.dbStorage.EXPECT().DeleteApplication(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)

	err = as.Delete(newApp.Namespace, newApp.Name, "")
	assert.NoError(t, err)
}

func TestDefaultApplicationService_CreateWithBase(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	mockIndexService := ms.NewMockIndexService(mockObject.ctl)

	as := applicationService{
		storage:      mockObject.modelStorage,
		indexService: mockIndexService,
		dbStorage:    mockObject.dbStorage,
	}
	config := &specV1.Configuration{Name: "agent-conf", Version: "123"}
	secret2 := &specV1.Secret{Name: "test-secret-02", Version: "123"}

	newApp, baseApp := genAppTestCase()
	mockIndexService.EXPECT().RefreshConfigIndexByApp(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
	mockIndexService.EXPECT().RefreshSecretIndexByApp(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
	mockObject.modelStorage.EXPECT().CreateApplication(gomock.Any(), gomock.Any()).Return(newApp, nil).Times(1)
	mockObject.modelStorage.EXPECT().GetConfig(gomock.Any(), gomock.Any(), gomock.Any()).Return(config, nil).Times(2)
	mockObject.modelStorage.EXPECT().GetSecret(gomock.Any(), gomock.Any(), gomock.Any()).Return(secret2, nil)
	mockObject.dbStorage.EXPECT().CreateApplication(gomock.Any()).Return(nil, fmt.Errorf("error"))
	_, err := as.CreateWithBase(newApp.Namespace, newApp, baseApp)
	assert.NoError(t, err)

	mockObject.modelStorage.EXPECT().CreateConfig(gomock.Any(), gomock.Any()).Return(config, nil).AnyTimes()
	mockObject.dbStorage.EXPECT().CreateApplication(gomock.Any()).Return(nil, nil).AnyTimes()
	mockObject.modelStorage.EXPECT().GetConfig(gomock.Any(), gomock.Any(), gomock.Any()).Return(config, fmt.Errorf("error")).Times(1)
	baseApp.Namespace = "test01"
	_, err = as.CreateWithBase(newApp.Namespace, newApp, baseApp)
	assert.NotNil(t, err)

	newApp, baseApp = genAppTestCase()
	mockIndexService.EXPECT().RefreshConfigIndexByApp(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockIndexService.EXPECT().RefreshSecretIndexByApp(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockObject.modelStorage.EXPECT().CreateApplication(gomock.Any(), gomock.Any()).Return(newApp, nil).AnyTimes()
	mockObject.modelStorage.EXPECT().GetConfig(gomock.Any(), gomock.Any(), gomock.Any()).Return(config, nil).AnyTimes()
	mockObject.modelStorage.EXPECT().GetSecret(gomock.Any(), secret2.Name, gomock.Any()).Return(secret2, nil)
	baseApp.Namespace = "test02"
	_, err = as.CreateWithBase(newApp.Namespace, newApp, baseApp)
	assert.NoError(t, err)

	newApp, baseApp = genAppTestCase()
	mockObject.modelStorage.EXPECT().GetSecret(gomock.Any(), secret2.Name, gomock.Any()).Return(secret2, nil)
	baseApp.Namespace = "test02"
	_, err = as.CreateWithBase(newApp.Namespace, newApp, baseApp)
	assert.NoError(t, err)

	newApp, baseApp = genAppTestCase()
	newApp.Volumes = append(newApp.Volumes, specV1.Volume{Name: "test"})
	baseApp.Namespace = "test02"
	_, err = as.CreateWithBase(newApp.Namespace, newApp, baseApp)
	assert.Error(t, err)

	newApp, baseApp = genAppTestCase()
	newApp.Services = append(newApp.Services, specV1.Service{Name: "Agent"})
	baseApp.Namespace = "test02"
	_, err = as.CreateWithBase(newApp.Namespace, newApp, baseApp)
	assert.Error(t, err)

	newApp, baseApp = genAppTestCase()
	newApp.Services[0].VolumeMounts = append(newApp.Services[0].VolumeMounts,
		specV1.VolumeMount{
			Name:      "test_01",
			MountPath: "mountPath",
		})
	baseApp.Namespace = "test02"
	_, err = as.CreateWithBase(newApp.Namespace, newApp, baseApp)
	assert.Error(t, err)
}

func TestDefaultApplicationService_Update(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	mockIndexService := ms.NewMockIndexService(mockObject.ctl)
	as := applicationService{
		storage:      mockObject.modelStorage,
		indexService: mockIndexService,
		dbStorage:    mockObject.dbStorage,
	}

	newApp, oldApp := genAppTestCase()
	mockObject.modelStorage.EXPECT().GetConfig(gomock.Any(), gomock.Any(), "").Return(nil, fmt.Errorf("error")).Times(1)
	_, err := as.Update(newApp.Namespace, newApp)
	assert.NotNil(t, err)

	secret1 := &specV1.Secret{Name: "test-secret-01", Version: "123"}
	secret2 := &specV1.Secret{Name: "test-secret-02", Version: "123"}

	mockIndexService.EXPECT().RefreshConfigIndexByApp(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
	mockIndexService.EXPECT().RefreshSecretIndexByApp(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
	mockObject.modelStorage.EXPECT().UpdateApplication(newApp.Namespace, newApp).Return(oldApp, nil)
	mockObject.modelStorage.EXPECT().GetConfig(gomock.Any(), gomock.Any(), "").Return(&specV1.Configuration{Version: "1"}, nil).AnyTimes()
	mockObject.modelStorage.EXPECT().GetSecret(gomock.Any(), secret1.Name, gomock.Any()).Return(secret1, nil).AnyTimes()
	mockObject.modelStorage.EXPECT().GetSecret(gomock.Any(), secret2.Name, gomock.Any()).Return(secret2, nil).AnyTimes()
	mockObject.dbStorage.EXPECT().CreateApplication(gomock.Any()).Return(nil, fmt.Errorf("error"))
	_, err = as.Update(newApp.Namespace, newApp)
	assert.NoError(t, err)

	newApp, _ = genAppTestCase()
	mockObject.modelStorage.EXPECT().UpdateApplication(newApp.Namespace, newApp).Return(nil, fmt.Errorf("error"))
	_, err = as.Update(newApp.Namespace, newApp)
	assert.NotNil(t, err)

	_, oldApp = genAppTestCase()
	mockIndexService.EXPECT().RefreshConfigIndexByApp(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("error")).Times(1)
	mockObject.modelStorage.EXPECT().UpdateApplication(gomock.Any(), gomock.Any()).Return(oldApp, nil)
	_, err = as.Update(newApp.Namespace, newApp)
	assert.NotNil(t, err)

}

func TestDefaultApplicationService_constuctConfig(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	cs := applicationService{
		storage: mockObject.modelStorage,
	}

	_, baseApp := genAppTestCase()
	mockObject.modelStorage.EXPECT().GetConfig(baseApp.Namespace, baseApp.Volumes[0].Config.Name, "").Return(nil, fmt.Errorf("error")).Times(1)
	err := cs.constuctConfig("default", baseApp)
	assert.NotNil(t, err)

	config := &specV1.Configuration{Name: "agent-conf"}
	_, baseApp = genAppTestCase()
	mockObject.modelStorage.EXPECT().
		GetConfig(baseApp.Namespace, baseApp.Volumes[0].Config.Name, "").
		Return(config, nil)
	mockObject.modelStorage.EXPECT().CreateConfig(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
	mockObject.modelStorage.EXPECT().CreateConfig(gomock.Any(), gomock.Any()).Return(config, nil)
	err = cs.constuctConfig("default", baseApp)
	assert.NoError(t, err)
}

type Test1 struct {
	a  time.Time  `json:"a,omitempty"`
	b  *time.Time `json:"b,omitempty"`
	t1 string
}

func TestNewCallbackService(t *testing.T) {

	tt := Test1{
		a:  time.Now(),
		t1: "test",
	}
	b, _ := json.Marshal(&tt)

	fmt.Println(string(b))
}
