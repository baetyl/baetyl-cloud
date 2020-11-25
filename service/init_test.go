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
	as := InitServiceImpl{}
	tp := service.NewMockTemplateService(mockCtl)
	ns := service.NewMockNodeService(mockCtl)
	sc := service.NewMockSecretService(mockCtl)
	sApp := service.NewMockApplicationService(mockCtl)
	sConfig := service.NewMockConfigService(mockCtl)
	sSecret := service.NewMockSecretService(mockCtl)
	as.AppCombinedService = &AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}
	as.ResourceMapFunc = map[string]GetInitResource{}
	as.TemplateService = tp
	as.NodeService = ns
	as.Secret = sc
	as.ResourceMapFunc[templateInitDeploymentYaml] = as.getInitDeploymentYaml
	desire := &v1.Desire{
		"sysapps": []specV1.AppInfo{{
			Name:    "baetyl-init-node01",
			Version: "123",
		}},
	}
	app := &specV1.Application{
		Namespace: "default",
		Name:      "baetyl-init-node01",
		Selector:  "test",
		Version:   "1",
		Services:  []specV1.Service{},
		Volumes: []specV1.Volume{
			{
				Name: "node-cert",
				VolumeSource: specV1.VolumeSource{
					Secret: &specV1.ObjectReference{
						Name: "agent-conf",
					},
				},
			},
		},
	}
	sec := &specV1.Secret{
		Namespace: "default",
		Name:      "abc",
	}
	//good case : setup
	tp.EXPECT().ParseTemplate(templateInitDeploymentYaml, gomock.Any()).Return([]byte("init"), nil).Times(1)
	ns.EXPECT().GetDesire("default", "node1").Return(desire, nil)
	sApp.EXPECT().Get("default", "baetyl-init-node01", "").Return(app, nil)
	sc.EXPECT().Get("default", "agent-conf", "").Return(sec, nil)

	res, _ := as.GetResource("default", "node1", templateInitDeploymentYaml, nil)
	assert.Equal(t, res, []byte("init"))

	// bad case : not found
	_, err := as.GetResource("", "", "dummy", nil)
	assert.EqualError(t, err, "The (resource) resource (dummy) is not found.")
}

func TestInitService_getInitYaml(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	sNode := service.NewMockNodeService(mockCtl)
	as := InitServiceImpl{}
	as.NodeService = sNode

	sNode.EXPECT().GetDesire("default", "n0").Return(nil, common.Error(common.ErrResourceNotFound))
	res, err := as.getInitDeploymentYaml("default", "n0", nil)
	assert.EqualError(t, err, "The resource is not found.")
	assert.Nil(t, res)
}

func TestInitService_getSyncCert(t *testing.T) {
	as := InitServiceImpl{}
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	sApp := service.NewMockApplicationService(mockCtl)
	sConfig := service.NewMockConfigService(mockCtl)
	sSecret := service.NewMockSecretService(mockCtl)
	as.AppCombinedService = &AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}

	app1 := &specV1.Application{
		Name:      "baetyl-core",
		Namespace: "default",
	}

	sSecret.EXPECT().Get("default", "", "").Return(nil, nil).Times(1)
	res, err := as.GetNodeCert(app1)
	assert.EqualError(t, err, "The (secret) resource is not found in namespace(default).")
	assert.Nil(t, res)
}

func TestInitService_GenCmd(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	sAuth := service.NewMockAuthService(mockCtl)
	sTemplate := service.NewMockTemplateService(mockCtl)
	sProp := service.NewMockPropertyService(mockCtl)
	as := InitServiceImpl{}
	as.AuthService = sAuth
	as.TemplateService = sTemplate
	as.Property = sProp
	info := map[string]interface{}{
		InfoName:      "name",
		InfoNamespace: "ns",
		InfoExpiry:    time.Now().Unix() + CmdExpirationInSeconds,
	}
	expect := "curl -skfL 'https://1.2.3.4:9003/v1/active/setup.sh?token=tokenexpect' -osetup.sh && sh setup.sh"
	params := map[string]interface{}{
		"InitApplyYaml": "baetyl-init-deployment.yml",
		"mode":          "",
	}
	sAuth.EXPECT().GenToken(info).Return("tokenexpect", nil).Times(1)
	sProp.EXPECT().GetPropertyValue(TemplateKubeInitCommand).Return(TemplateBaetylInitCommand, nil)
	sTemplate.EXPECT().Execute("setup-command", TemplateBaetylInitCommand, gomock.Any()).Return([]byte(expect), nil).Times(1)

	res, err := as.GetInitCommand("ns", "name", params)
	assert.NoError(t, err)
	assert.Equal(t, string(res), expect)
}

func TestInitService_getDesireAppInfo(t *testing.T) {
	as := InitServiceImpl{}
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	node := service.NewMockNodeService(mockCtl)
	app := service.NewMockApplicationService(mockCtl)
	as.NodeService = node
	as.AppCombinedService = &AppCombinedService{
		App: app,
	}

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

	res, err := as.GetAppFromDesire("default", "node01", "baetyl-core", true)
	assert.NoError(t, err)
	assert.Equal(t, res, app1)
}

func TestInitService_GenApps(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	sApp := service.NewMockApplicationService(mock)
	sConfig := service.NewMockConfigService(mock)
	sSecret := service.NewMockSecretService(mock)
	sTemplate := service.NewMockTemplateService(mock)
	sNode := service.NewMockNodeService(mock)
	sAuth := service.NewMockAuthService(mock)
	sPKI := service.NewMockPKIService(mock)

	is := InitServiceImpl{}
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

	sTemplate.EXPECT().UnmarshalTemplate("baetyl-core-conf.yml", gomock.Any(), gomock.Any()).Return(nil)
	sTemplate.EXPECT().UnmarshalTemplate("baetyl-core-app.yml", gomock.Any(), gomock.Any()).Return(nil)
	sTemplate.EXPECT().UnmarshalTemplate("baetyl-function-conf.yml", gomock.Any(), gomock.Any()).Return(nil)
	sTemplate.EXPECT().UnmarshalTemplate("baetyl-function-app.yml", gomock.Any(), gomock.Any()).Return(nil)
	sTemplate.EXPECT().UnmarshalTemplate("baetyl-broker-conf.yml", gomock.Any(), gomock.Any()).Return(nil)
	sTemplate.EXPECT().UnmarshalTemplate("baetyl-broker-app.yml", gomock.Any(), gomock.Any()).Return(nil)
	sTemplate.EXPECT().UnmarshalTemplate("baetyl-init-app.yml", gomock.Any(), gomock.Any()).Return(nil)
	sTemplate.EXPECT().UnmarshalTemplate("baetyl-init-conf.yml", gomock.Any(), gomock.Any()).Return(nil)
	sPKI.EXPECT().SignClientCertificate("ns.abc", gomock.Any()).Return(cert, nil)
	sPKI.EXPECT().GetCA().Return([]byte("RootCA"), nil)
	sConfig.EXPECT().Create("ns", gomock.Any()).Return(config, nil).Times(4)
	sSecret.EXPECT().Create("ns", gomock.Any()).Return(secret, nil).Times(1)
	sApp.EXPECT().Create("ns", gomock.Any()).Return(app, nil).Times(4)

	out, err := is.GenApps("ns", "abc")
	assert.NoError(t, err)
	assert.Equal(t, 4, len(out))
}
