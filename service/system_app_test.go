package service

import (
	"testing"

	v1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/mock/service"
	"github.com/baetyl/baetyl-cloud/v2/models"
)

func TestInitService_GenApps(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	sApp := service.NewMockApplicationService(mock)
	sConfig := service.NewMockConfigService(mock)
	sSecret := service.NewMockSecretService(mock)
	sTemplate := service.NewMockTemplateService(mock)
	sPKI := service.NewMockPKIService(mock)

	is := SystemAppServiceImpl{}
	is.TemplateService = sTemplate
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
	sTemplate.EXPECT().UnmarshalTemplate("baetyl-broker-conf.yml", gomock.Any(), gomock.Any()).Return(nil)
	sTemplate.EXPECT().UnmarshalTemplate("baetyl-broker-app.yml", gomock.Any(), gomock.Any()).Return(nil)
	sTemplate.EXPECT().UnmarshalTemplate("baetyl-init-app.yml", gomock.Any(), gomock.Any()).Return(nil)
	sTemplate.EXPECT().UnmarshalTemplate("baetyl-init-conf.yml", gomock.Any(), gomock.Any()).Return(nil)
	sPKI.EXPECT().SignClientCertificate("ns.abc", gomock.Any()).Return(cert, nil)
	sPKI.EXPECT().GetCA().Return([]byte("RootCA"), nil)
	sConfig.EXPECT().Create(gomock.Any(), "ns", gomock.Any()).Return(config, nil).Times(3)
	sSecret.EXPECT().Create(gomock.Any(), "ns", gomock.Any()).Return(secret, nil).Times(1)
	sApp.EXPECT().Create(gomock.Any(), "ns", gomock.Any()).Return(app, nil).Times(3)

	node := &v1.Node{
		Namespace: "ns",
		Name:      "abc",
	}
	out, err := is.GenApps(gomock.Any(), "ns", node)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(out))
}

func TestInitService_GenOptionalApps(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	sApp := service.NewMockApplicationService(mock)
	sConfig := service.NewMockConfigService(mock)
	sSecret := service.NewMockSecretService(mock)
	sTemplate := service.NewMockTemplateService(mock)
	sPKI := service.NewMockPKIService(mock)

	is := SystemAppServiceImpl{}
	is.TemplateService = sTemplate
	is.PKI = sPKI
	is.AppCombinedService = &AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}

	config := &v1.Configuration{
		Namespace: "ns",
		Name:      "config",
	}

	app := &v1.Application{
		Namespace: "ns",
		Name:      "app",
	}

	sTemplate.EXPECT().UnmarshalTemplate("baetyl-function-conf.yml", gomock.Any(), gomock.Any()).Return(nil)
	sTemplate.EXPECT().UnmarshalTemplate("baetyl-function-app.yml", gomock.Any(), gomock.Any()).Return(nil)
	sTemplate.EXPECT().UnmarshalTemplate("baetyl-rule-conf.yml", gomock.Any(), gomock.Any()).Return(nil)
	sTemplate.EXPECT().UnmarshalTemplate("baetyl-rule-app.yml", gomock.Any(), gomock.Any()).Return(nil)
	sConfig.EXPECT().Create(gomock.Any(), "ns", gomock.Any()).Return(config, nil).Times(2)
	sApp.EXPECT().Create(gomock.Any(), "ns", gomock.Any()).Return(app, nil).Times(2)

	node := &v1.Node{
		Namespace: "ns",
		Name:      "abc",
		SysApps: []string{
			"baetyl-function",
			"baetyl-rule",
		},
	}

	is.OptionalAppFuncs = map[string]GenAppFunc{
		"baetyl-function": is.genFunctionApp,
		"baetyl-rule":     is.genRuleApp,
	}

	out, err := is.GenOptionalApps(gomock.Any(), "ns", node.Name, node.SysApps)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(out))
}

func TestInitService_GetOptionalApps(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	sApp := service.NewMockApplicationService(mock)
	sConfig := service.NewMockConfigService(mock)
	sSecret := service.NewMockSecretService(mock)
	sTemplate := service.NewMockTemplateService(mock)
	sPKI := service.NewMockPKIService(mock)
	sProperty := service.NewMockPropertyService(mock)

	is := SystemAppServiceImpl{}
	is.TemplateService = sTemplate
	is.PKI = sPKI
	is.AppCombinedService = &AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}
	is.OptionalAppFuncs = map[string]GenAppFunc{
		"a": nil,
	}
	is.Property = sProperty

	apps := is.GetOptionalApps()
	assert.Len(t, apps, 1)
	assert.Equal(t, apps[0], "a")
}
