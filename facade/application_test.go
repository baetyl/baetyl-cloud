package facade

import (
	"testing"

	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/common"
)

func TestCreateApplication(t *testing.T) {
	mAppFacade, mCtl := InitMockEnvironment(t)
	defer mCtl.Finish()
	appFacade := &facade{
		node:      mAppFacade.sNode,
		app:       mAppFacade.sApp,
		config:    mAppFacade.sConfig,
		index:     mAppFacade.sIndex,
		txFactory: mAppFacade.txFactory,
	}
	mAppFacade.txFactory.EXPECT().BeginTx().Return(nil, nil).AnyTimes()
	mAppFacade.txFactory.EXPECT().Rollback(nil).Return().AnyTimes()

	config := &specV1.Configuration{}
	app := &specV1.Application{
		Name: "abc",
	}
	configs := []specV1.Configuration{*config}
	ns := "baetyl-cloud"

	mAppFacade.sConfig.EXPECT().Upsert(nil, ns, gomock.Any()).Return(nil, unknownErr)
	_, err := appFacade.CreateApp(ns, app, app, configs)
	assert.Error(t, err, unknownErr)

	mAppFacade.sConfig.EXPECT().Upsert(nil, ns, gomock.Any()).Return(nil, nil).AnyTimes()
	mAppFacade.sApp.EXPECT().CreateWithBase(nil, ns, gomock.Any(), gomock.Any()).Return(nil, unknownErr).Times(1)
	_, err = appFacade.CreateApp(ns, app, app, configs)
	assert.Error(t, err, unknownErr)

	mAppFacade.sApp.EXPECT().CreateWithBase(nil, ns, gomock.Any(), gomock.Any()).Return(app, nil).AnyTimes()
	mAppFacade.sNode.EXPECT().UpdateNodeAppVersion(nil, ns, gomock.Any()).Return(nil, unknownErr).Times(1)
	_, err = appFacade.CreateApp(ns, app, app, configs)
	assert.Error(t, err, unknownErr)

	mAppFacade.txFactory.EXPECT().Commit(nil).Return().AnyTimes()
	mAppFacade.sNode.EXPECT().UpdateNodeAppVersion(nil, ns, gomock.Any()).Return(nil, nil)
	mAppFacade.sIndex.EXPECT().RefreshNodesIndexByApp(nil, ns, gomock.Any(), gomock.Any()).Return(nil)
	_, err = appFacade.CreateApp(ns, app, app, configs)
	assert.NoError(t, err)
}

func TestDeleteApplication(t *testing.T) {
	mAppFacade, mCtl := InitMockEnvironment(t)
	defer mCtl.Finish()
	appFacade := &facade{
		node:      mAppFacade.sNode,
		app:       mAppFacade.sApp,
		config:    mAppFacade.sConfig,
		index:     mAppFacade.sIndex,
		txFactory: mAppFacade.txFactory,
	}
	ns := "baetyl-cloud"

	// Function
	app := &specV1.Application{
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
				Name: "baetyl-function-code-agent",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name: "func1",
					},
				},
			},
			{
				Name: "baetyl-function-conf-agent",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name: "baetyl-function-config-app-service-xxxxxxxxx",
					},
				},
			},
		},
	}
	mAppFacade.sApp.EXPECT().Delete(ns, app.Name, "").Return(unknownErr).Times(1)
	err := appFacade.DeleteApp(ns, app.Name, app)
	assert.Error(t, err, unknownErr)

	mAppFacade.sApp.EXPECT().Delete(ns, app.Name, "").Return(nil).AnyTimes()
	mAppFacade.sNode.EXPECT().DeleteNodeAppVersion(nil, ns, app).Return(nil, unknownErr).Times(1)
	err = appFacade.DeleteApp(ns, app.Name, app)
	assert.Error(t, err, unknownErr)

	mAppFacade.sNode.EXPECT().DeleteNodeAppVersion(nil, ns, app).Return(nil, nil).Times(1)
	mAppFacade.sIndex.EXPECT().RefreshNodesIndexByApp(nil, ns, app.Name, gomock.Any()).Return(nil).AnyTimes()
	mAppFacade.sConfig.EXPECT().Delete(nil, ns, gomock.Any()).Return(unknownErr)
	err = appFacade.DeleteApp(ns, app.Name, app)
	assert.NoError(t, err)
}

func TestUpdateApplication(t *testing.T) {
	mAppFacade, mCtl := InitMockEnvironment(t)
	defer mCtl.Finish()
	appFacade := &facade{
		node:      mAppFacade.sNode,
		app:       mAppFacade.sApp,
		config:    mAppFacade.sConfig,
		index:     mAppFacade.sIndex,
		txFactory: mAppFacade.txFactory,
	}

	config := &specV1.Configuration{}
	app := &specV1.Application{
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
				Name: "baetyl-function-code-agent",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name: "func1",
					},
				},
			},
		},
	}
	configs := []specV1.Configuration{*config}
	ns := "baetyl-cloud"

	mAppFacade.sConfig.EXPECT().Upsert(nil, ns, gomock.Any()).Return(nil, unknownErr).Times(1)
	_, err := appFacade.UpdateApp(ns, app, app, configs)
	assert.Error(t, err, unknownErr)

	mAppFacade.sConfig.EXPECT().Upsert(nil, ns, gomock.Any()).Return(nil, nil).AnyTimes()
	mAppFacade.sApp.EXPECT().Update(ns, app).Return(nil, unknownErr).Times(1)
	_, err = appFacade.UpdateApp(ns, app, app, configs)
	assert.Error(t, err, unknownErr)

	mAppFacade.sApp.EXPECT().Update(ns, app).Return(app, nil).AnyTimes()
	oldApp := &specV1.Application{
		Selector: "test",
	}
	mAppFacade.sNode.EXPECT().DeleteNodeAppVersion(nil, ns, oldApp).Return(nil, unknownErr).Times(1)
	_, err = appFacade.UpdateApp(ns, oldApp, app, configs)
	assert.Error(t, err, unknownErr)

	mAppFacade.sNode.EXPECT().UpdateNodeAppVersion(nil, ns, gomock.Any()).Return(nil, unknownErr).Times(1)
	_, err = appFacade.UpdateApp(ns, app, app, configs)
	assert.Error(t, err, unknownErr)

	mAppFacade.sNode.EXPECT().UpdateNodeAppVersion(nil, ns, gomock.Any()).Return(nil, nil).AnyTimes()
	mAppFacade.sIndex.EXPECT().RefreshNodesIndexByApp(nil, ns, app.Name, gomock.Any()).Return(nil).AnyTimes()
	mAppFacade.sConfig.EXPECT().Delete(nil, ns, gomock.Any()).Return(nil).AnyTimes()
	_, err = appFacade.UpdateApp(ns, app, app, configs)
	assert.NoError(t, err)
}
