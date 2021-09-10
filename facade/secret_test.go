package facade

import (
	"testing"

	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateSecret(t *testing.T) {
	mFacade, mCtl := InitMockEnvironment(t)
	defer mCtl.Finish()
	sFacade := &facade{
		node:      mFacade.sNode,
		app:       mFacade.sApp,
		config:    mFacade.sConfig,
		secret:    mFacade.sSecret,
		index:     mFacade.sIndex,
		txFactory: mFacade.txFactory,
	}
	ns := "test"
	mFacade.txFactory.EXPECT().BeginTx().Return(nil, nil).AnyTimes()
	mFacade.txFactory.EXPECT().Rollback(nil).Return().Times(1)
	mFacade.sSecret.EXPECT().Create(nil, ns, gomock.Any()).Return(nil, unknownErr).Times(1)
	_, err := sFacade.CreateSecret(ns, nil)
	assert.Error(t, err, unknownErr)

	mFacade.txFactory.EXPECT().Commit(nil).Return().Times(1)
	mFacade.sSecret.EXPECT().Create(nil, ns, gomock.Any()).Return(nil, nil).Times(1)
	_, err = sFacade.CreateSecret(ns, nil)
	assert.NoError(t, err)
}

func TestUpdateSecret(t *testing.T) {
	mFacade, mCtl := InitMockEnvironment(t)
	defer mCtl.Finish()
	sFacade := &facade{
		node:      mFacade.sNode,
		app:       mFacade.sApp,
		config:    mFacade.sConfig,
		secret:    mFacade.sSecret,
		index:     mFacade.sIndex,
		txFactory: mFacade.txFactory,
	}
	ns, name := "default", "abc"
	mConf := &specV1.Secret{
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
	appNames := []string{"app1", "app2"}
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

	mFacade.sSecret.EXPECT().Update(ns, gomock.Any()).Return(nil, unknownErr)
	_, err := sFacade.UpdateSecret(ns, mConf)
	assert.Error(t, err, unknownErr)

	mFacade.sSecret.EXPECT().Update(ns, gomock.Any()).Return(mConfSecret3, nil).AnyTimes()
	mFacade.sIndex.EXPECT().ListAppIndexBySecret(ns, name).Return(appNames, nil).Times(1)
	mFacade.sApp.EXPECT().Get(ns, appNames[0], "").Return(apps[0], nil).Times(1)
	mFacade.sApp.EXPECT().Get(ns, appNames[1], "").Return(apps[1], nil).Times(1)
	mFacade.sApp.EXPECT().Update(nil, ns, gomock.Any()).Return(apps[0], nil).Times(1)
	mFacade.sNode.EXPECT().UpdateNodeAppVersion(nil, ns, gomock.Any()).Return(nil, nil).Times(1)
	_, err = sFacade.UpdateSecret(ns, mConf)
	assert.NoError(t, err)

	appNames = []string{"app01"}
	apps = []*specV1.Application{
		{
			Namespace: "default",
			Name:      appNames[0],
			Volumes: []specV1.Volume{
				{
					Name:         "vol0",
					VolumeSource: specV1.VolumeSource{Secret: &specV1.ObjectReference{Name: "abc", Version: "1"}},
				},
			},
		},
	}
	mFacade.sIndex.EXPECT().ListAppIndexBySecret(ns, name).Return(appNames, nil).AnyTimes()
	mFacade.sApp.EXPECT().Get(ns, appNames[0], "").Return(apps[0], nil).Times(1)
	mFacade.sApp.EXPECT().Update(nil, ns, gomock.Any()).Return(nil, unknownErr).Times(1)
	_, err = sFacade.UpdateSecret(ns, mConf)
	assert.Error(t, err, unknownErr)

	apps[0].Volumes[0].Secret.Version = "1"
	mFacade.sApp.EXPECT().Get(ns, appNames[0], "").Return(apps[0], nil).Times(1)
	mFacade.sApp.EXPECT().Update(nil, ns, gomock.Any()).Return(nil, nil).Times(1)
	mFacade.sNode.EXPECT().UpdateNodeAppVersion(nil, ns, gomock.Any()).Return(nil, unknownErr).Times(1)
	_, err = sFacade.UpdateSecret(ns, mConf)
	assert.Error(t, err, unknownErr)
}

func TestDeleteSecret(t *testing.T) {
	mFacade, mCtl := InitMockEnvironment(t)
	defer mCtl.Finish()
	sFacade := &facade{
		node:      mFacade.sNode,
		app:       mFacade.sApp,
		config:    mFacade.sConfig,
		secret:    mFacade.sSecret,
		index:     mFacade.sIndex,
		txFactory: mFacade.txFactory,
	}
	ns, n := "test", "test"

	mFacade.sSecret.EXPECT().Delete(ns, n).Return(nil).Times(1)
	err := sFacade.DeleteSecret(ns, n)
	assert.NoError(t, err)
}
