package service

import (
	"testing"
	"time"

	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	v1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/mock/service"
	"github.com/baetyl/baetyl-cloud/v2/models"
)

func TestInitService_GetResource(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	tp := service.NewMockTemplateService(mockCtl)
	ns := service.NewMockNodeService(mockCtl)
	aus := service.NewMockAuthService(mockCtl)
	as := InitServiceImpl{}
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

func TestInitService_getInitYaml(t *testing.T) {
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

func TestInitService_getInitYaml(t *testing.T) {
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

	app1 := &specV1.Application{
		Name:      "baetyl-core",
		Namespace: "default",
	}

	secret.EXPECT().Get("default", "", "").Return(nil, nil).Times(1)
	res, err := as.getSyncCert(app1)
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

func TestInitService_GenApps(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	sApp := service.NewMockApplicationService(mock)
	sConfig := service.NewMockConfigService(mock)
	sSecret := service.NewMockSecretService(mock)
	sCache := service.NewMockCacheService(mock)
	sTemplate := service.NewMockTemplateService(mock)
	sNode := service.NewMockNodeService(mock)
	sAuth := service.NewMockAuthService(mock)
	sPKI := service.NewMockPKIService(mock)
	is := InitServiceImpl{}
	is.CacheService = sCache
	is.TemplateService = sTemplate
	is.NodeService = sNode
	is.AuthService = sAuth
	is.PKI = sPKI
	is.AppCombinedService = &AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}

	cert := &models.PEMCredential{
		CertPEM: []byte("CertPEM"),
		KeyPEM:  []byte("KeyPEM"),
		CertId:  "CertId",
	}

	config := &v1.Configuration{
		Namespace: "ns",
		Name:      "config",
	}

	secret := &v1.Secret{
		Namespace: "ns",
		Name:      "secret",
	}

	app := &v1.Application{
		Namespace: "ns",
		Name:      "app",
	}

	sCache.EXPECT().GetProperty("sync-server-address").Return("https://localhost:50001", nil)
	sTemplate.EXPECT().UnmarshalTemplate("baetyl-core-conf.yml", gomock.Any(), gomock.Any()).Return(nil)
	sTemplate.EXPECT().UnmarshalTemplate("baetyl-core-app.yml", gomock.Any(), gomock.Any()).Return(nil)
	sTemplate.EXPECT().UnmarshalTemplate("baetyl-function-conf.yml", gomock.Any(), gomock.Any()).Return(nil)
	sTemplate.EXPECT().UnmarshalTemplate("baetyl-function-app.yml", gomock.Any(), gomock.Any()).Return(nil)
	sPKI.EXPECT().SignClientCertificate("ns.abc", gomock.Any()).Return(cert, nil)
	sPKI.EXPECT().GetCA().Return([]byte("RootCA"), nil)
	sConfig.EXPECT().Create("ns", gomock.Any()).Return(config, nil).Times(2)
	sSecret.EXPECT().Create("ns", gomock.Any()).Return(secret, nil).Times(1)
	sApp.EXPECT().Create("ns", gomock.Any()).Return(app, nil).Times(2)

	out, err := is.GenApps("ns", "abc", nil)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(out))
}

func TestInitService_GetSetupShell(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	sCache := service.NewMockCacheService(mock)
	sTemplate := service.NewMockTemplateService(mock)
	ns := service.NewMockNodeService(mock)
	aus := service.NewMockAuthService(mock)
	is := InitServiceImpl{}
	is.CacheService = sCache
	is.TemplateService = sTemplate
	is.NodeService = ns
	is.AuthService = aus

	sCache.EXPECT().GetProperty("init-server-address").Return("https://localhost:50001", nil)
	sTemplate.EXPECT().ParseTemplate("setup.sh", gomock.Any()).Return([]byte("setup"), nil)

	actual, err := is.genSetupShell("xxx")
	assert.NoError(t, err)
	assert.Equal(t, "setup", string(actual))
}
