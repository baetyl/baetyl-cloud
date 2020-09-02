package service

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/mock/service"
)

func TestAPI_GetResource(t *testing.T) {
	as := InitServiceImpl{}
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	tp := service.NewMockTemplateService(mockCtl)
	ns := service.NewMockNodeService(mockCtl)
	aus := service.NewMockAuthService(mockCtl)
	as.TemplateService = tp
	as.NodeService = ns
	as.AuthService = aus

	// good case : metrics
	tp.EXPECT().GetTemplate(common.ResourceMetrics).Return("metrics", nil).Times(1)
	res, _ := as.GetResource(common.ResourceMetrics, "", "", nil)
	assert.Equal(t, res, []byte("metrics"))

	// good case : local_path_storage
	tp.EXPECT().GetTemplate(common.ResourceLocalPathStorage).Return("local-path-storage", nil).Times(1)
	res, _ = as.GetResource(common.ResourceLocalPathStorage, "", "", nil)
	assert.Equal(t, res, []byte("local-path-storage"))

	// good case : setup
	tp.EXPECT().GenSetupShell(gomock.Any()).Return([]byte("shell"), nil).Times(1)
	res, _ = as.GetResource(common.ResourceSetup, "", "", nil)
	assert.Equal(t, res, []byte("shell"))

	// bad case : not found
	_, err := as.GetResource("", "", "", nil)
	assert.Error(t, err)
}

func TestApi_getInitYaml(t *testing.T) {
	info := map[string]interface{}{
		InfoKind:      "123",
		InfoName:      "n0",
		InfoNamespace: "default",
		InfoTimestamp: time.Now().Unix(),
		InfoExpiry:    60 * 60 * 24 * 3650,
	}
	as := InitServiceImpl{}
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	auth := service.NewMockAuthService(mockCtl)
	as.AuthService = auth
	res, err := as.getInitYaml(info, "kube")
	assert.Error(t, err, common.ErrRequestParamInvalid)
	assert.Nil(t, res)
}

func TestApi_getSysParams(t *testing.T) {
	as := InitServiceImpl{}
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	cache := service.NewMockCacheService(mockCtl)
	as.CacheService = cache

	cache.EXPECT().GetProperty("baetyl-image").Return("baetyl-image", nil).Times(1)
	cache.EXPECT().GetProperty(propertySyncServerAddress).Return("1.2.3.4", nil).Times(1)
	cache.EXPECT().GetProperty(propertyInitServerAddress).Return("5.6.7.8", nil).Times(1)
	expect := map[string]interface{}{
		"Namespace":           "default",
		"EdgeNamespace":       common.DefaultBaetylEdgeNamespace,
		"EdgeSystemNamespace": common.DefaultBaetylEdgeSystemNamespace,
		"NodeAddress":         "1.2.3.4",
		"ActiveAddress":       "5.6.7.8",
		"Image":               "baetyl-image",
		"KubeNodeName":        "kube",
	}
	res, err := as.getSysParams("default", "kube")
	assert.NoError(t, err)
	assert.Equal(t, res, expect)
}

func TestApi_getSyncCert(t *testing.T) {
	as := InitServiceImpl{}
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	secret := service.NewMockSecretService(mockCtl)
	as.SecretService = secret

	secret.EXPECT().Get("", "", "").Return(nil, nil).Times(1)
	res, err := as.getSyncCert(nil)
	assert.Error(t, err, common.ErrResourceNotFound)
	assert.Nil(t, res)
}

func TestApi_GenCmd(t *testing.T) {
	as := InitServiceImpl{}
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	auth := service.NewMockAuthService(mockCtl)
	cache := service.NewMockCacheService(mockCtl)
	as.AuthService = auth
	as.CacheService = cache

	info := map[string]interface{}{
		InfoKind:      "kind",
		InfoName:      "name",
		InfoNamespace: "ns",
		InfoExpiry:    CmdExpirationInSeconds,
		InfoTimestamp: time.Now().Unix(),
	}
	auth.EXPECT().GenToken(info).Return("tokenexpect", nil).Times(1)
	cache.EXPECT().GetProperty(propertyInitServerAddress).Return("https://1.2.3.4:9003", nil).Times(1)
	expect := "curl -skfL 'https://1.2.3.4:9003/v1/active/setup.sh?token=tokenexpect' -osetup.sh && sh setup.sh"

	res, err := as.GenCmd("kind", "ns", "name")
	assert.NoError(t, err)
	assert.Equal(t, res, expect)
}

func TestApi_getDesireAppInfo(t *testing.T) {
	as := InitServiceImpl{}
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	node := service.NewMockNodeService(mockCtl)
	app := service.NewMockApplicationService(mockCtl)
	as.NodeService = node
	as.AppService = app

	Desire := &specV1.Desire{
		"sysapps": []specV1.AppInfo{{
			Name:    "baetyl-core-node01",
			Version: "123",
		}},
	}
	app1 := &specV1.Application{
		Name:      "baetyl-core",
		Namespace: "default",
	}
	node.EXPECT().GetDesire("default", "node01").Return(Desire, nil).Times(1)
	app.EXPECT().Get("default", "baetyl-core-node01", "").Return(app1, nil).Times(1)

	res, err := as.getDesireAppInfo("default", "node01")
	assert.NoError(t, err)
	assert.Equal(t, res, app1)
}
