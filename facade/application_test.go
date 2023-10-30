package facade

import (
	"testing"
	"time"

	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

func TestCreateApplication(t *testing.T) {
	mAppFacade, mCtl := InitMockEnvironment(t)
	defer mCtl.Finish()
	appFacade := &facade{
		node:      mAppFacade.sNode,
		app:       mAppFacade.sApp,
		config:    mAppFacade.sConfig,
		index:     mAppFacade.sIndex,
		cron:      mAppFacade.sCron,
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

	app.CronStatus = specV1.CronWait
	mAppFacade.sCron.EXPECT().CreateCron(gomock.Any()).Return(nil)
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
		cron:      mAppFacade.sCron,
		txFactory: mAppFacade.txFactory,
	}
	ns := "baetyl-cloud"

	// Function
	app := &specV1.Application{
		Namespace: "baetyl-cloud",
		Name:      "abc",
		Type:      specV1.AppTypeFunction,
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
	mAppFacade.txFactory.EXPECT().BeginTx().Return(nil, nil).AnyTimes()
	mAppFacade.txFactory.EXPECT().Rollback(nil).Return().AnyTimes()
	mAppFacade.txFactory.EXPECT().Commit(nil).Return().AnyTimes()

	mAppFacade.sApp.EXPECT().Delete(nil, ns, app.Name, "").Return(unknownErr).Times(1)
	err := appFacade.DeleteApp(ns, app.Name, app)
	assert.Error(t, err, unknownErr)

	mAppFacade.sApp.EXPECT().Delete(nil, ns, app.Name, "").Return(nil).AnyTimes()
	mAppFacade.sNode.EXPECT().DeleteNodeAppVersion(nil, ns, app).Return(nil, unknownErr).Times(1)
	err = appFacade.DeleteApp(ns, app.Name, app)
	assert.Error(t, err, unknownErr)

	app.CronStatus = specV1.CronWait
	mAppFacade.sCron.EXPECT().DeleteCron(app.Name, ns).Return(nil)
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
		cron:      mAppFacade.sCron,
		txFactory: mAppFacade.txFactory,
	}

	config := &specV1.Configuration{}
	app := &specV1.Application{
		Namespace: "baetyl-cloud",
		Name:      "abc",
		Type:      specV1.AppTypeFunction,
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

	mAppFacade.txFactory.EXPECT().BeginTx().Return(nil, nil).AnyTimes()
	mAppFacade.txFactory.EXPECT().Rollback(nil).Return().AnyTimes()
	mAppFacade.txFactory.EXPECT().Commit(nil).Return().AnyTimes()

	mAppFacade.sConfig.EXPECT().Upsert(nil, ns, gomock.Any()).Return(nil, unknownErr).Times(1)
	_, err := appFacade.UpdateApp(ns, app, app, configs)
	assert.Error(t, err, unknownErr)

	mAppFacade.sConfig.EXPECT().Upsert(nil, ns, gomock.Any()).Return(nil, nil).AnyTimes()
	mAppFacade.sApp.EXPECT().Update(nil, ns, app).Return(nil, unknownErr).Times(1)
	_, err = appFacade.UpdateApp(ns, app, app, configs)
	assert.Error(t, err, unknownErr)

	mAppFacade.sApp.EXPECT().Update(nil, ns, app).Return(app, nil).AnyTimes()
	oldApp := &specV1.Application{
		Selector: "test",
	}
	mAppFacade.sNode.EXPECT().DeleteNodeAppVersion(nil, ns, oldApp).Return(nil, unknownErr).Times(1)
	_, err = appFacade.UpdateApp(ns, oldApp, app, configs)
	assert.Error(t, err, unknownErr)

	mAppFacade.sNode.EXPECT().UpdateNodeAppVersion(nil, ns, gomock.Any()).Return(nil, unknownErr).Times(1)
	_, err = appFacade.UpdateApp(ns, app, app, configs)
	assert.Error(t, err, unknownErr)

	app.CronStatus = specV1.CronWait
	appNew := &specV1.Application{
		Namespace:  "baetyl-cloud",
		Name:       "abc",
		Type:       specV1.AppTypeFunction,
		CronStatus: specV1.CronNotSet,
	}
	mAppFacade.sCron.EXPECT().UpdateCron(gomock.Any()).Return(nil).AnyTimes()
	mAppFacade.sCron.EXPECT().DeleteCron(app.Name, ns).Return(unknownErr).Times(1)
	_, err = appFacade.UpdateApp(ns, app, appNew, configs)
	assert.Error(t, err, unknownErr)

	mAppFacade.sNode.EXPECT().UpdateNodeAppVersion(nil, ns, gomock.Any()).Return(nil, nil).AnyTimes()
	mAppFacade.sIndex.EXPECT().RefreshNodesIndexByApp(nil, ns, app.Name, gomock.Any()).Return(nil).AnyTimes()
	mAppFacade.sConfig.EXPECT().Delete(nil, ns, gomock.Any()).Return(nil).AnyTimes()
	_, err = appFacade.UpdateApp(ns, app, app, configs)
	assert.NoError(t, err)
}

func TestGetApplication(t *testing.T) {
	mAppFacade, mCtl := InitMockEnvironment(t)
	defer mCtl.Finish()
	appFacade := &facade{
		app:  mAppFacade.sApp,
		cron: mAppFacade.sCron,
	}
	name, ns := "baetyl", "cloud"
	mAppFacade.sApp.EXPECT().Get(ns, name, "").Return(nil, unknownErr).Times(1)
	_, err := appFacade.GetApp(ns, name, "")
	assert.Error(t, err, unknownErr)

	app := &specV1.Application{
		Namespace:  ns,
		Name:       name,
		CronStatus: specV1.CronWait,
		CronTime:   time.Now(),
	}
	cronApp := &models.Cron{
		Namespace: ns,
		Name:      name,
		Selector:  "",
		CronTime:  time.Now(),
	}
	mAppFacade.sApp.EXPECT().Get(ns, name, "").Return(app, nil).Times(1)
	mAppFacade.sCron.EXPECT().GetCron(name, ns).Return(cronApp, nil).Times(1)
	_, err = appFacade.GetApp(ns, name, "")
	assert.NoError(t, err)
}
