package service

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/baetyl/baetyl-cloud/v2/common"
	v1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/stretchr/testify/assert"
)

func TestTemplateService_GenApps(t *testing.T) {
	mock := InitMockEnvironment(t)
	defer mock.Close()

	ts, err := NewTemplateService(mock.conf)
	assert.NoError(t, err)

	node := &v1.Node{
		Namespace: "ns",
		Name:      "abc",
		Labels: map[string]string{
			common.LabelNodeName: "abc",
		},
	}

	var cfgs []*v1.Configuration
	var secrets []*v1.Secret
	var apps []*v1.Application
	mock.property.EXPECT().GetPropertyValue("sync-server-address").Return("https://localhost:50001", nil)
	mock.pki.EXPECT().GetRootCertId().Return("RootCertID").Times(2)
	mock.pki.EXPECT().CreateClientCert(gomock.Any(), "RootCertID").Return("ClientCertID", nil)
	mock.pki.EXPECT().GetClientCert("ClientCertID").Return([]byte("ClientCertData"), nil)
	mock.pki.EXPECT().GetRootCert("RootCertID").Return([]byte("RootCertData"), nil)
	mock.modelStorage.EXPECT().CreateConfig(node.Namespace, gomock.Any()).DoAndReturn(
		func(ns string, c *v1.Configuration) (*v1.Configuration, error) {
			c.Version = fmt.Sprintf("v%d", len(cfgs))
			cfgs = append(cfgs, c)
			fmt.Println("CreateConfig", c)
			return c, nil
		}).Times(2)
	mock.modelStorage.EXPECT().CreateSecret(node.Namespace, gomock.Any()).DoAndReturn(
		func(ns string, s *v1.Secret) (*v1.Secret, error) {
			s.Version = fmt.Sprintf("v%d", len(secrets))
			secrets = append(secrets, s)
			fmt.Println("CreateSecret", s)
			return s, nil
		}).Times(1)
	mock.modelStorage.EXPECT().CreateApplication(node.Namespace, gomock.Any()).DoAndReturn(
		func(ns string, a *v1.Application) (*v1.Application, error) {
			a.Version = fmt.Sprintf("v%d", len(apps))
			apps = append(apps, a)
			fmt.Println("CreateApplication", a)
			return a, nil
		}).Times(2)
	mock.modelStorage.EXPECT().GetConfig(node.Namespace, gomock.Any(), "").Return(&v1.Configuration{}, nil).Times(2)
	mock.modelStorage.EXPECT().GetSecret(node.Namespace, gomock.Any(), "").Return(&v1.Secret{}, nil)
	mock.dbStorage.EXPECT().RefreshIndex(node.Namespace, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(4)
	mock.dbStorage.EXPECT().CreateApplication(gomock.Any()).Return(nil, nil).Times(2)

	out, err := ts.GenSystemApps(node.Namespace, node.Name, nil)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(out))
	assert.Equal(t, "baetyltechtest/baetyl:git-122bf10", out[0].Services[0].Image)
	assert.Equal(t, "baetyltechtest/function:git-226e40d", out[1].Services[0].Image)
}

func TestTemplateService_GetSetupShell(t *testing.T) {
	mock := InitMockEnvironment(t)
	defer mock.Close()

	ts, err := NewTemplateService(mock.conf)
	assert.NoError(t, err)

	mock.property.EXPECT().GetPropertyValue("active-server-address").Return("https://localhost:50002", nil)

	// TOKEN="{{.Token}}"
	// CLOUD_ADDR="{{.CloudAddr}}"
	actual, err := ts.GenSetupShell("xxx")
	expected, err := ioutil.ReadFile("../scripts/demo/native/templates/setup.sh")
	expectedstr := strings.Replace(string(expected), "{{.Token}}", "xxx", 1)
	expectedstr = strings.Replace(expectedstr, "{{.CloudAddr}}", "https://localhost:50002", 1)
	assert.NoError(t, err)
	assert.Equal(t, expectedstr, string(actual))
}

// TODO: add more bad cases
//func TestTemplateService_GenApps_ErrUpdateNode(t *testing.T) {
//	ctl := gomock.NewController(t)
//	defer ctl.Finish()
//
//	as := ms.NewMockApplicationService(ctl)
//	cs := ms.NewMockConfigService(ctl)
//	ss := ms.NewMockSecretService(ctl)
//	scs := ms.NewMockSysConfigService(ctl)
//	pki := ms.NewMockPKIService(ctl)
//	ns := ms.NewMockNodeService(ctl)
//	is := ms.NewMockIndexService(ctl)
//	init := ms.NewMockInitializeService(ctl)
//
//	api := &API{
//		applicationService: as,
//		configService:      cs,
//		secretService:      ss,
//		sysConfigService:   scs,
//		pkiService:         pki,
//		nodeService:        ns,
//		indexService:       is,
//		initService:        init,
//	}
//	node := getMockNode()
//	list := []common.SystemApplication{
//		common.BaetylCore,
//		common.BaetylFunction,
//	}
//	conf := &v1.Configuration{
//		Name:      "testConf",
//		Namespace: node.Namespace,
//		Data:      nil,
//	}
//	app := &v1.Application{
//		Name:      "testApp",
//		Namespace: node.Namespace,
//	}
//	nodeList := []string{"s0", "s1", "s2"}
//	sysConf := &models.SysConfig{
//		Type:  "baetyl-edge",
//		Key:   "test",
//		Value: "123",
//	}
//	certPEM := &models.PEMCredential{
//		CertPEM: []byte("test"),
//		KeyPEM:  []byte("test"),
//	}
//	certMap := map[string][]byte{
//		"client.pem": certPEM.CertPEM,
//		"client.key": certPEM.KeyPEM,
//		"ca.pem":     []byte("test"),
//	}
//	secret := &v1.Secret{
//		Name:      "sync-" + node.Name + "-core",
//		Namespace: node.Namespace,
//		Data:      certMap,
//		Version:   "123",
//	}
//	cs.EXPECT().Create(node.Namespace, gomock.Any()).Return(conf, nil).Times(1)
//	as.EXPECT().Create(node.Namespace, gomock.Any()).Return(app, nil).Times(1)
//	ss.EXPECT().Get(node.Namespace, gomock.Any(), "").Return(secret, nil).AnyTimes()
//	scs.EXPECT().GetSysConfig(gomock.Any(), gomock.Any()).Return(sysConf, nil).AnyTimes()
//	pki.EXPECT().SignClientCertificate(gomock.Any(), gomock.Any()).Return(certPEM, nil).AnyTimes()
//	pki.EXPECT().GetCA().Return([]byte("test"), nil).AnyTimes()
//	ss.EXPECT().Create(node.Namespace, gomock.Any()).Return(secret, nil).AnyTimes()
//	ns.EXPECT().UpdateNodeAppVersion(node.Namespace, gomock.Any()).Return(nodeList, nil).Times(1)
//	is.EXPECT().RefreshNodesIndexByApp(node.Namespace, gomock.Any(), nodeList).Return(fmt.Errorf("update err")).Times(1)
//	init.EXPECT().GetResource(gomock.Any()).Return("{}", nil).AnyTimes()
//
//	_, err := api.GenSysApp(node.Name, node.Namespace, list)
//	assert.NotNil(t, err)
//	assert.Equal(t, "update err", err.Error())
//}

//
//func TestAPI_GenSysApp_ErrSysConfig(t *testing.T) {
//	ctl := gomock.NewController(t)
//	defer ctl.Finish()
//
//	scs := ms.NewMockSysConfigService(ctl)
//
//	api := &API{
//		sysConfigService: scs,
//	}
//	node := getMockNode()
//	list := []common.SystemApplication{
//		common.BaetylCore,
//	}
//	err1 := common.Error(
//		common.ErrResourceNotFound,
//		common.Field("type", "system config baetyl module"),
//		common.Field("name", "core"))
//	scs.EXPECT().GetSysConfig(common.BaetylModule, string(common.BaetylCore)).Return(nil, err1).Times(1)
//
//	_, err2 := api.GenSysApp(node.Name, node.Namespace, list)
//	assert.NotNil(t, err1)
//	assert.Equal(t, err1, err2)
//
//	sysConf := &models.SysConfig{
//		Type:  "baetyl-edge",
//		Key:   "test",
//		Value: "123",
//	}
//	err3 := common.Error(
//		common.ErrResourceNotFound,
//		common.Field("type", "system config address"),
//		common.Field("name", common.AddressNode))
//	scs.EXPECT().GetSysConfig(common.BaetylModule, string(common.BaetylCore)).Return(sysConf, nil).Times(1)
//	scs.EXPECT().GetSysConfig("address", common.AddressNode).Return(nil, err3).Times(1)
//	_, err4 := api.GenSysApp(node.Name, node.Namespace, list)
//	assert.NotNil(t, err4)
//	assert.Equal(t, err3, err4)
//
//	list = []common.SystemApplication{
//		common.BaetylFunction,
//	}
//	err5 := common.Error(
//		common.ErrResourceNotFound,
//		common.Field("type", "system config baetyl module"),
//		common.Field("name", "function"))
//	scs.EXPECT().GetSysConfig(common.BaetylModule, string(common.BaetylFunction)).Return(nil, err5).Times(1)
//
//	_, err6 := api.GenSysApp(node.Name, node.Namespace, list)
//	assert.NotNil(t, err6)
//	assert.Equal(t, err5, err6)
//}
//
//func Test_genConfig_Err(t *testing.T) {
//	ctl := gomock.NewController(t)
//	defer ctl.Finish()
//
//	cs := ms.NewMockConfigService(ctl)
//	init := ms.NewMockInitializeService(ctl)
//
//	api := &API{
//		configService: cs,
//		initService:   init,
//	}
//
//	ns := "default"
//	confName := "123"
//	templateKey := "test"
//	template := "{\n  \"name\": \"{{.ConfigName}}\",\n  \"namespace\": \"{{.Namespace}}\",\n  \"system\": true,\n  \"data\": {\n    \"service.yml\": \"logger:\\n  filename: var/log/baetyl/service.log\\n  level: info\"\n  }\n}"
//	params := map[string]string{
//		"ConfigName": confName,
//		"Namespace":  ns,
//	}
//
//	// bad case 0
//	init.EXPECT().GetResource(templateKey).Return("", fmt.Errorf("get template err")).Times(1)
//	_, err := api.genConfig(ns, templateKey, params)
//	assert.Error(t, err)
//
//	// bad case 1
//	init.EXPECT().GetResource(templateKey).Return("error json", nil).Times(1)
//	_, err = api.genConfig(ns, templateKey, params)
//	assert.Error(t, err)
//
//	// bad case 2
//	init.EXPECT().GetResource(templateKey).Return(template, nil).Times(1)
//	cs.EXPECT().Create(ns, gomock.Any()).Return(nil, os.ErrInvalid).Times(1)
//	cs.EXPECT().Get(ns, confName, "").Return(nil, os.ErrInvalid).Times(1)
//	_, err = api.genConfig(ns, templateKey, params)
//	assert.Error(t, err, common.Error(
//		common.ErrResourceNotFound,
//		common.Field("type", "config"),
//		common.Field("name", confName),
//		common.Field("namespace", ns)))
//}
//
//func Test_genApp_Err(t *testing.T) {
//	ctl := gomock.NewController(t)
//	defer ctl.Finish()
//
//	as := ms.NewMockApplicationService(ctl)
//	init := ms.NewMockInitializeService(ctl)
//
//	api := &API{
//		applicationService: as,
//		initService:        init,
//	}
//
//	ns := "default"
//	appName := "123"
//	nodeName := "456"
//	templateKey := "test"
//	template := "{\"name\":\"{{.Name}}\",\"namespace\":\"{{.Namespace}}\",\"selector\":\"baetyl-node-name={{.NodeName}}\",\"labels\":{\"baetyl-cloud-system\":\"{{.Name}}\"},\"services\":[],\"volumes\":[]}"
//	params := map[string]string{
//		"Name":      appName,
//		"Namespace": ns,
//		"NodeName":  nodeName,
//	}
//
//	// bad case 0
//	init.EXPECT().GetResource(templateKey).Return("", fmt.Errorf("get template err")).Times(1)
//	_, err := api.genApp(ns, templateKey, params)
//	assert.Error(t, err)
//
//	// bad case 1
//	init.EXPECT().GetResource(templateKey).Return("error json", nil).Times(1)
//	_, err = api.genApp(ns, templateKey, params)
//	assert.Error(t, err)
//
//	// bad case 2
//	init.EXPECT().GetResource(templateKey).Return(template, nil).Times(1)
//	as.EXPECT().Create(ns, gomock.Any()).Return(nil, os.ErrInvalid).Times(1)
//	as.EXPECT().Get(ns, appName, "").Return(nil, os.ErrInvalid).Times(1)
//	_, err = api.genApp(ns, templateKey, params)
//	assert.Error(t, err, common.Error(
//		common.ErrResourceNotFound,
//		common.Field("type", "config"),
//		common.Field("name", appName),
//		common.Field("namespace", ns)))
//
