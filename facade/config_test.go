package facade

import (
	"testing"

	"github.com/baetyl/baetyl-go/v2/errors"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
)

func TestCreateConfig(t *testing.T) {
	mFacade, mCtl := InitMockEnvironment(t)
	defer mCtl.Finish()
	cfgFacade := &facade{
		node:      mFacade.sNode,
		app:       mFacade.sApp,
		config:    mFacade.sConfig,
		index:     mFacade.sIndex,
		txFactory: mFacade.txFactory,
	}
	ns := "test"
	mFacade.txFactory.EXPECT().BeginTx().Return(nil, nil).AnyTimes()
	mFacade.txFactory.EXPECT().Rollback(nil).Return().Times(1)
	mFacade.sConfig.EXPECT().Create(nil, ns, gomock.Any()).Return(nil, unknownErr).Times(1)
	_, err := cfgFacade.CreateConfig(ns, nil)
	assert.Error(t, err, unknownErr)

	mFacade.txFactory.EXPECT().Commit(nil).Return().Times(1)
	mFacade.sConfig.EXPECT().Create(nil, ns, gomock.Any()).Return(nil, nil).Times(1)
	_, err = cfgFacade.CreateConfig(ns, nil)
	assert.NoError(t, err)
}

func TestUpdateConfig(t *testing.T) {
	mFacade, mCtl := InitMockEnvironment(t)
	defer mCtl.Finish()
	cfgFacade := &facade{
		node:      mFacade.sNode,
		app:       mFacade.sApp,
		config:    mFacade.sConfig,
		index:     mFacade.sIndex,
		txFactory: mFacade.txFactory,
	}
	ns, name := "default", "abc"
	res := &specV1.Configuration{
		Name:      name,
		Namespace: ns,
		Version:   "10",
		Labels: map[string]string{
			"test": "test",
		},
		Data: map[string]string{
			common.ConfigObjectPrefix + "function": `{"metadata":{"bucket":"baetyl","function":"process","handler":"index.handler","object":"a.zip","runtime":"python36","type":"function","userID":"default","version":"1"}}`,
		},
	}
	mConf2 := &models.ConfigurationView{
		Namespace: ns,
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
	res3 := &specV1.Configuration{
		Name:      name,
		Namespace: ns,
		Labels: map[string]string{
			"test": "test",
		},
		Data: map[string]string{
			common.ConfigObjectPrefix + "function": `{"metadata":{"bucket":"baetyl","function":"process","handler":"index.handler","object":"a.zip","runtime":"python36","type":"function","version":"1"}}`,
		},
		Description: "diff",
	}

	mFacade.sConfig.EXPECT().Update(nil, ns, gomock.Any()).Return(res, unknownErr).Times(1)
	_, err := cfgFacade.UpdateConfig(ns, res3)
	assert.Error(t, err, unknownErr)

	mFacade.sConfig.EXPECT().Update(nil, ns, gomock.Any()).Return(res, nil).AnyTimes()
	mFacade.sIndex.EXPECT().ListAppIndexByConfig(mConf2.Namespace, "abc").Return(nil, unknownErr).Times(1)
	_, err = cfgFacade.UpdateConfig(ns, res3)
	assert.Error(t, err, unknownErr)

	appNames := make([]string, 0)
	mFacade.sIndex.EXPECT().ListAppIndexByConfig(ns, name).Return(appNames, nil).Times(1)
	_, err = cfgFacade.UpdateConfig(ns, res3)
	assert.NoError(t, err)

	appNames = []string{"app01", "app02"}
	mFacade.sIndex.EXPECT().ListAppIndexByConfig(ns, name).Return(appNames, nil).Times(1)
	mFacade.sApp.EXPECT().Get(ns, "app01", "").Return(nil, errors.New("err")).Times(1)
	_, err = cfgFacade.UpdateConfig(ns, res3)
	assert.Error(t, err, unknownErr)

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
	}
	mFacade.sIndex.EXPECT().ListAppIndexByConfig(ns, name).Return(appNames, nil).Times(1)
	mFacade.sApp.EXPECT().Get(ns, appNames[0], "").Return(apps[0], nil).Times(1)
	mFacade.sApp.EXPECT().Get(ns, appNames[1], "").Return(apps[1], nil).Times(1)
	_, err = cfgFacade.UpdateConfig(ns, res3)
	assert.NoError(t, err)

	appNames = []string{"app01"}
	apps = []*specV1.Application{
		{
			Namespace: "default",
			Name:      appNames[0],
			Volumes: []specV1.Volume{
				{
					Name:         "vol0",
					VolumeSource: specV1.VolumeSource{Config: &specV1.ObjectReference{Name: "abc", Version: "1"}},
				},
			},
		},
	}
	mFacade.sIndex.EXPECT().ListAppIndexByConfig(ns, name).Return(appNames, nil).AnyTimes()
	mFacade.sApp.EXPECT().Get(ns, appNames[0], "").Return(apps[0], nil).Times(1)
	mFacade.sApp.EXPECT().Update(nil, ns, gomock.Any()).Return(nil, unknownErr).Times(1)
	_, err = cfgFacade.UpdateConfig(ns, res3)
	assert.Error(t, err, unknownErr)

	apps[0].Volumes[0].Config.Version = "1"
	mFacade.sApp.EXPECT().Get(ns, appNames[0], "").Return(apps[0], nil).Times(1)
	mFacade.sApp.EXPECT().Update(nil, ns, gomock.Any()).Return(nil, nil).Times(1)
	mFacade.sNode.EXPECT().UpdateNodeAppVersion(nil, ns, gomock.Any()).Return(nil, unknownErr).Times(1)
	_, err = cfgFacade.UpdateConfig(ns, res3)
	assert.Error(t, err, unknownErr)
}

func TestDeleteConfig(t *testing.T) {
	mFacade, mCtl := InitMockEnvironment(t)
	defer mCtl.Finish()
	cfgFacade := &facade{
		node:      mFacade.sNode,
		app:       mFacade.sApp,
		config:    mFacade.sConfig,
		index:     mFacade.sIndex,
		txFactory: mFacade.txFactory,
	}
	ns, n := "test", "test"
	mFacade.sConfig.EXPECT().Delete(nil, ns, n).Return(nil)
	err := cfgFacade.DeleteConfig(ns, n)
	assert.NoError(t, err)
}
