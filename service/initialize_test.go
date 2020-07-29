package service

import (
	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	mockPlugin "github.com/baetyl/baetyl-cloud/v2/mock/plugin"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path"
	"testing"
)

func genInitStruct(t *testing.T) (*config.CloudConfig, *mockPlugin.MockModelStorage, *mockPlugin.MockDBStorage) {
	conf := mockTestConfig()
	mockCtl := gomock.NewController(t)
	mockModelStorage := mockPlugin.NewMockModelStorage(mockCtl)
	plugin.RegisterFactory(conf.Plugin.ModelStorage, mockStorageModel(mockModelStorage))
	mockDBStorage := mockPlugin.NewMockDBStorage(mockCtl)
	plugin.RegisterFactory(conf.Plugin.DatabaseStorage, mockStorageDB(mockDBStorage))

	return conf, mockModelStorage, mockDBStorage
}

func TestInitializeService_GetSync(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	is, err := NewInitializeService(mockObject.conf)
	assert.Nil(t, err)

	// bad case 0
	node := &specV1.Node{Name: "error", Namespace: "error"}
	mockObject.dbStorage.EXPECT().Get(node.Namespace, node.Name).Return(nil, nil)

	_, err = is.GetSyncCert(node.Namespace, node.Name)
	assert.Error(t, err)

	shadow := genShadowTestCase()
	// bad case 1
	app := &specV1.Application{
		Name:      "baetyl-core-" + shadow.Name,
		Namespace: shadow.Namespace,
		Volumes:   []specV1.Volume{},
	}

	mockObject.dbStorage.EXPECT().Get(gomock.Any(), gomock.Any()).Return(shadow, nil)
	mockObject.modelStorage.EXPECT().GetApplication(shadow.Namespace, app.Name, "").Return(app, nil).Times(1)
	mockObject.modelStorage.EXPECT().GetSecret(app.Namespace, "", "").Return(nil, nil).Times(1)
	_, err = is.GetSyncCert(shadow.Namespace, shadow.Name)
	assert.Error(t, err)
}

func TestInitializeService_GetResource(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	is, err := NewInitializeService(mockObject.conf)
	assert.Nil(t, err)

	// good case 0
	common.Cache = map[string]string{
		"exist": "exist",
	}
	res, err := is.GetResource("exist")
	assert.NoError(t, err)
	assert.Equal(t, common.Cache["exist"], res)

	// bad case 0
	mockObject.dbStorage.EXPECT().GetSysConfig("resource", "error").Return(nil, nil).Times(1)
	res, err = is.GetResource("error")
	assert.Error(t, err)

	// bad case 1
	sc := &models.SysConfig{Type: "resource", Key: "base64 error", Value: "err"}
	mockObject.dbStorage.EXPECT().GetSysConfig("resource", sc.Key).Return(sc, nil).Times(1)
	res, err = is.GetResource(sc.Key)
	assert.Error(t, err)

	// good case 1
	common.Cache = nil
	sc = &models.SysConfig{Type: "resource", Key: "test", Value: "MTIz"}
	mockObject.dbStorage.EXPECT().GetSysConfig("resource", sc.Key).Return(sc, nil).Times(1)
	res, err = is.GetResource(sc.Key)
	assert.Nil(t, err)
	assert.Equal(t, "123", res)
}

func TestGenInitCmd(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	tempDir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	err = ioutil.WriteFile(path.Join(tempDir, "ca.pem"), []byte("ca.pem"), 777)
	assert.NoError(t, err)
	mockObject.conf.ActiveServer.Certificate.CA = path.Join(tempDir, "ca.pem")

	common.Cache = map[string]string{
		common.ResourceInitYaml: "123",
	}

	kube := "k3s"
	node := genNodeTestCase()
	sysConf := &models.SysConfig{
		Type:  "baetyl-edge",
		Key:   "test",
		Value: "123",
	}
	certMap := map[string][]byte{
		"client.pem": []byte("test"),
		"client.key": []byte("test"),
		"ca.pem":     []byte("test"),
	}
	secret := &specV1.Secret{
		Name:      "sync-" + node.Name + "-core",
		Namespace: node.Namespace,
		Data:      certMap,
		Version:   "123",
	}
	app := &specV1.Application{
		Namespace: node.Namespace,
		Version:   "123",
	}
	batch := genBatchTestCase()
	shadow := &models.Shadow{
		Name:      node.Name,
		Namespace: node.Namespace,
		Desire: specV1.Desire{
			common.DesiredSysApplications: []specV1.AppInfo{
				{
					Name:    "baetyl-core-" + node.Name,
					Version: "1",
				},
			},
		},
	}

	mockObject.modelStorage.EXPECT().GetNode(secret.Namespace, gomock.Any()).Return(node, nil).AnyTimes()
	mockObject.dbStorage.EXPECT().Get(gomock.Any(), gomock.Any()).Return(shadow, nil).AnyTimes()
	mockObject.dbStorage.EXPECT().GetSysConfig(gomock.Any(), gomock.Any()).Return(sysConf, nil).AnyTimes()
	mockObject.modelStorage.EXPECT().GetSecret(secret.Namespace, gomock.Any(), "").Return(secret, nil).AnyTimes()
	mockObject.modelStorage.EXPECT().CreateSecret(secret.Namespace, gomock.Any()).Return(secret, nil).AnyTimes()
	mockObject.modelStorage.EXPECT().GetApplication(node.Namespace, "baetyl-core-"+node.Name, "").Return(app, nil).AnyTimes()

	is, err := NewInitializeService(mockObject.conf)
	assert.NoError(t, err)

	res, err := is.InitWithNode(node.Namespace, node.Name, kube)
	assert.NoError(t, err)
	assert.Equal(t, "123", string(res))

	batch.Fingerprint.Type = common.FingerprintSN
	batch.Fingerprint.SnPath = "sn"
	res, err = is.InitWithBitch(batch, kube)
	assert.NoError(t, err)
	assert.Equal(t, "123", string(res))

	batch.Fingerprint.Type = common.FingerprintInput
	batch.Fingerprint.InputField = "sn"
	res, err = is.InitWithBitch(batch, kube)
	assert.NoError(t, err)
	assert.Equal(t, "123", string(res))
}

func TestInitializeService_getSysParams(t *testing.T) {
	cfg, mockModel, mockDB := genInitStruct(t)
	init := &initializeService{
		cfg:          cfg,
		modelStorage: mockModel,
		dbStorage:    mockDB,
	}
	kube := "k3s"

	// bad case 0
	sc0 := &models.SysConfig{Type: common.BaetylModule, Key: string(common.BaetylInit), Value: ""}
	mockDB.EXPECT().GetSysConfig(sc0.Type, sc0.Key).Return(nil, nil).Times(1)
	_, err := init.getSysParams("test", kube)
	assert.Error(t, err)

	// bad case 1
	sc1 := &models.SysConfig{Type: "address", Key: common.AddressNode, Value: ""}
	mockDB.EXPECT().GetSysConfig(sc0.Type, sc0.Key).Return(sc0, nil).Times(1)
	mockDB.EXPECT().GetSysConfig(sc1.Type, sc1.Key).Return(nil, nil).Times(1)
	_, err = init.getSysParams("test", kube)
	assert.Error(t, err)

	// bad case 2
	sc2 := &models.SysConfig{Type: "address", Key: common.AddressActive, Value: ""}
	mockDB.EXPECT().GetSysConfig(sc0.Type, sc0.Key).Return(sc0, nil).Times(1)
	mockDB.EXPECT().GetSysConfig(sc1.Type, sc1.Key).Return(sc1, nil).Times(1)
	mockDB.EXPECT().GetSysConfig(sc2.Type, sc2.Key).Return(nil, nil).Times(1)
	_, err = init.getSysParams("test", kube)
	assert.Error(t, err)
}

func TestInitializeService_getCoreApp(t *testing.T) {
	cfg, mockModel, mockDB := genInitStruct(t)
	init := &initializeService{
		cfg:          cfg,
		modelStorage: mockModel,
		dbStorage:    mockDB,
		shadow:       mockDB,
	}
	node := genNodeTestCase()

	shadow := &models.Shadow{
		Name:      "",
		Namespace: "",
		Desire: specV1.Desire{
			common.DesiredSysApplications: []specV1.AppInfo{
				{
					Name:    "core-" + node.Name,
					Version: "1",
				},
			},
		},
	}
	mockDB.EXPECT().Get(gomock.Any(), gomock.Any()).Return(shadow, nil).AnyTimes()

	// bad case 0
	mockModel.EXPECT().GetNode(node.Namespace, node.Name).Return(nil, nil).Times(1)
	_, err := init.getCoreApp(node.Namespace, node.Name)
	assert.Error(t, err)

	// bad case 1
	mockModel.EXPECT().GetNode(node.Namespace, node.Name).Return(node, nil).Times(1)
	mockModel.EXPECT().GetApplication(node.Namespace, "baetyl-core-abc", "").Return(nil, nil).Times(1)
	_, err = init.getCoreApp(node.Namespace, node.Name)
	assert.Error(t, err)

	// bad case 2
	nErr := &specV1.Node{
		Namespace: node.Namespace,
		Name:      node.Name,
		Version:   "123",
		Desire: map[string]interface{}{
			"sysapps": []specV1.AppInfo{},
		},
	}
	mockModel.EXPECT().GetNode(node.Namespace, node.Name).Return(nErr, nil).Times(1)
	_, err = init.getCoreApp(node.Namespace, node.Name)
	assert.Error(t, err)
}

func TestInitializeService_getSyncCert(t *testing.T) {
	cfg, mockModel, mockDB := genInitStruct(t)
	init := &initializeService{
		cfg:          cfg,
		modelStorage: mockModel,
		dbStorage:    mockDB,
	}
	app := &specV1.Application{
		Namespace: "default",
		Volumes: []specV1.Volume{
			{
				Name:         "cert-sync",
				VolumeSource: specV1.VolumeSource{Secret: &specV1.ObjectReference{Name: "cert-sync"}},
			},
		},
	}

	// bad case 0
	mockModel.EXPECT().GetSecret(app.Namespace, "cert-sync", "").Return(nil, nil).Times(1)
	_, err := init.getSyncCert(app)
	assert.Error(t, err)
}
