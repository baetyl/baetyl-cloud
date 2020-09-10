package service

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/baetyl/baetyl-go/v2/errors"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/models"
)

//go:generate mockgen -destination=../mock/service/init.go -package=service github.com/baetyl/baetyl-cloud/v2/service InitService

const (
	InfoName      = "n"
	InfoNamespace = "ns"
	InfoTimestamp = "ts"
	InfoExpiry    = "e"
)

const (
	templateCoreConfYaml             = "baetyl-core-conf.yml"
	templateCoreAppYaml              = "baetyl-core-app.yml"
	templateFuncConfYaml             = "baetyl-function-conf.yml"
	templateFuncAppYaml              = "baetyl-function-app.yml"
	TemplateInitDeploymentYaml       = "baetyl-init-deployment.yml"
	templateKubeAPIMetricsYaml       = "kube-api-metrics.yml"
	templateKubeLocalPathStorageYaml = "kube-local-path-storage.yml"
	TemplateInitSetupShell           = "kube-init-setup.sh"
	templateKubeInitCommand          = `curl -skfL '{{GetProperty "init-server-address"}}/v1/init/kube-init-setup.sh?token={{.Token}}' -osetup.sh && sh setup.sh`
)

var (
	CmdExpirationInSeconds = int64(60 * 60)
	HookNamePopulateParams = "populateParams"
)

type HandlerPopulateParams func(ns string, params map[string]interface{}) error
type GetInitResource func(ns, nodeName, resourceName string, params map[string]interface{}) ([]byte, error)

// InitService
type InitService interface {
	GetResource(ns, nodeName, resourceName string, params map[string]interface{}) (interface{}, error)
	GenApps(ns, nodeName string) ([]*specV1.Application, error)
	GenCmd(ns, nameNode string) (string, error)
}

type InitServiceImpl struct {
	cfg             *config.CloudConfig
	AuthService     AuthService
	NodeService     NodeService
	SecretService   SecretService
	TemplateService TemplateService
	*AppCombinedService
	PKI             PKIService
	Hooks           map[string]interface{}
	ResourceMapFunc map[string]GetInitResource
}

// NewSyncService new SyncService
func NewInitService(config *config.CloudConfig) (InitService, error) {
	authService, err := NewAuthService(config)
	if err != nil {
		return nil, errors.Trace(err)
	}
	nodeService, err := NewNodeService(config)
	if err != nil {
		return nil, errors.Trace(err)
	}
	secretService, err := NewSecretService(config)
	if err != nil {
		return nil, errors.Trace(err)
	}
	propertyService, err := NewPropertyService(config)
	if err != nil {
		return nil, errors.Trace(err)
	}
	templateService, err := NewTemplateService(config, map[string]interface{}{
		"GetProperty": propertyService.GetPropertyValue,
		"RandString":  common.RandString,
	})
	if err != nil {
		return nil, errors.Trace(err)
	}
	acs, err := NewAppCombinedService(config)
	if err != nil {
		return nil, errors.Trace(err)
	}
	pki, err := NewPKIService(config)
	if err != nil {
		return nil, errors.Trace(err)
	}
	initService := &InitServiceImpl{
		cfg:                config,
		AuthService:        authService,
		NodeService:        nodeService,
		TemplateService:    templateService,
		SecretService:      secretService,
		AppCombinedService: acs,
		PKI:                pki,
		Hooks:              map[string]interface{}{},
		ResourceMapFunc:    map[string]GetInitResource{},
	}
	initService.ResourceMapFunc[templateKubeAPIMetricsYaml] = initService.getCommonResource
	initService.ResourceMapFunc[templateKubeLocalPathStorageYaml] = initService.getCommonResource
	initService.ResourceMapFunc[TemplateInitSetupShell] = initService.getInitSetupShell
	initService.ResourceMapFunc[TemplateInitDeploymentYaml] = initService.getInitDeploymentYaml

	return initService, nil
}

func (s *InitServiceImpl) GetResource(ns, nodeName, resourceName string, params map[string]interface{}) (interface{}, error) {
	if handler, ok := s.ResourceMapFunc[resourceName]; ok {
		if params == nil {
			params = map[string]interface{}{}
		}
		return handler(ns, nodeName, resourceName, params)
	}
	return nil, common.Error(
		common.ErrResourceNotFound,
		common.Field("type", "resource"),
		common.Field("name", resourceName))
}

func (s *InitServiceImpl) getCommonResource(_, _, resourceName string, _ map[string]interface{}) ([]byte, error) {
	res, err := s.TemplateService.GetTemplate(resourceName)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return []byte(res), nil
}

func (s *InitServiceImpl) getInitSetupShell(_, _, resourceName string, params map[string]interface{}) ([]byte, error) {
	params["DeploymentYml"] = TemplateInitDeploymentYaml
	data, err := s.TemplateService.ParseTemplate(resourceName, params)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return data, nil
}

func (s *InitServiceImpl) getInitDeploymentYaml(ns, nodeName, _ string, params map[string]interface{}) ([]byte, error) {
	app, err := s.GetCoreAppFromDesire(ns, nodeName)
	if err != nil {
		return nil, errors.Trace(err)
	}
	cert, err := s.GetNodeCert(app)
	if err != nil {
		return nil, errors.Trace(err)
	}

	params["Namespace"] = ns
	params["NodeName"] = nodeName
	params["NodeCertName"] = cert.Name
	params["NodeCertVersion"] = cert.Version
	params["NodeCertPem"] = base64.StdEncoding.EncodeToString(cert.Data["client.pem"])
	params["NodeCertKey"] = base64.StdEncoding.EncodeToString(cert.Data["client.key"])
	params["NodeCertCa"] = base64.StdEncoding.EncodeToString(cert.Data["ca.pem"])
	params["EdgeNamespace"] = common.DefaultBaetylEdgeNamespace
	params["EdgeSystemNamespace"] = common.DefaultBaetylEdgeSystemNamespace
	return s.TemplateService.ParseTemplate(TemplateInitDeploymentYaml, params)
}

func (s *InitServiceImpl) GetNodeCert(app *specV1.Application) (*specV1.Secret, error) {
	certName := ""
	for _, vol := range app.Volumes {
		if vol.Name == "node-cert" || vol.Name == "sync-cert" {
			certName = vol.Secret.Name
			break
		}
	}
	cert, _ := s.SecretService.Get(app.Namespace, certName, "")
	if cert == nil {
		return nil, common.Error(
			common.ErrResourceNotFound,
			common.Field("type", "secret"),
			common.Field("name", certName),
			common.Field("namespace", app.Namespace))
	}
	return cert, nil
}

func (s *InitServiceImpl) GenCmd(ns, nameNode string) (string, error) {
	info := map[string]interface{}{
		InfoNamespace: ns,
		InfoName:      nameNode,
		InfoExpiry:    CmdExpirationInSeconds,
		InfoTimestamp: time.Now().Unix(),
	}
	token, err := s.AuthService.GenToken(info)
	if err != nil {
		return "", err
	}
	params := map[string]interface{}{
		"Token": token,
	}
	data, err := s.TemplateService.Execute("setup-command", templateKubeInitCommand, params)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *InitServiceImpl) GetCoreAppFromDesire(ns, nodeName string) (*specV1.Application, error) {
	shadowDesire, err := s.NodeService.GetDesire(ns, nodeName)
	if err != nil {
		return nil, errors.Trace(err)
	}
	apps := shadowDesire.AppInfos(true)
	for _, appInfo := range apps {
		if strings.Contains(appInfo.Name, "baetyl-core") {
			app, _ := s.App.Get(ns, appInfo.Name, "")
			if app == nil {
				return nil, common.Error(
					common.ErrResourceNotFound,
					common.Field("type", "application"),
					common.Field("name", appInfo.Name),
					common.Field("namespace", ns))
			}
			return app, nil
		}
	}
	return nil, common.Error(
		common.ErrResourceNotFound,
		common.Field("type", "sysapp"),
		common.Field("name", nodeName),
		common.Field("namespace", ns))
}

func (s *InitServiceImpl) GenApps(ns, nodeName string) ([]*specV1.Application, error) {
	params := map[string]interface{}{
		"Namespace": ns,
		"NodeName":  nodeName,
	}
	if handler, ok := s.Hooks[HookNamePopulateParams]; ok {
		err := handler.(HandlerPopulateParams)(ns, params)
		if err != nil {
			return nil, errors.Trace(err)
		}
	}

	var apps []*specV1.Application
	ca, err := s.genCoreApp(ns, nodeName, params)
	if err != nil {
		return nil, errors.Trace(err)
	}
	fa, err := s.genFunctionApp(ns, nodeName, params)
	if err != nil {
		return nil, errors.Trace(err)
	}
	apps = append(apps, ca, fa)
	return apps, nil
}

func (s *InitServiceImpl) genCoreApp(ns, nodeName string, params map[string]interface{}) (*specV1.Application, error) {
	appName := fmt.Sprintf("baetyl-core-%s", common.RandString(9))
	confName := fmt.Sprintf("baetyl-core-conf-%s", common.RandString(9))
	params["CoreAppName"] = appName
	params["CoreConfName"] = confName

	// create config
	conf, err := s.genConfig(ns, templateCoreConfYaml, params)
	if err != nil {
		return nil, errors.Trace(err)
	}

	// create secret
	cert, err := s.genNodeCerts(ns, nodeName, appName)
	if err != nil {
		return nil, errors.Trace(err)
	}

	params["CoreConfVersion"] = conf.Version
	params["NodeCertName"] = cert.Name
	params["NodeCertVersion"] = cert.Version

	// create application
	return s.genApp(ns, templateCoreAppYaml, params)
}

func (s *InitServiceImpl) genFunctionApp(ns, nodeName string, params map[string]interface{}) (*specV1.Application, error) {
	appName := fmt.Sprintf("baetyl-function-%s", common.RandString(9))
	confName := fmt.Sprintf("baetyl-function-conf-%s", common.RandString(9))
	// create config
	confMap := map[string]interface{}{
		"Namespace":        ns,
		"NodeName":         nodeName,
		"FunctionAppName":  appName,
		"FunctionConfName": confName,
	}
	for k, v := range params {
		confMap[k] = v
	}
	conf, err := s.genConfig(ns, templateFuncConfYaml, confMap)
	if err != nil {
		return nil, errors.Trace(err)
	}

	// create application
	appMap := map[string]interface{}{
		"Namespace":           ns,
		"NodeName":            nodeName,
		"FunctionAppName":     appName,
		"FunctionConfName":    conf.Name,
		"FunctionConfVersion": conf.Version,
	}
	for k, v := range params {
		appMap[k] = v
	}
	return s.genApp(ns, templateFuncAppYaml, appMap)
}

func (s *InitServiceImpl) genNodeCerts(ns, nodeName, appName string) (*specV1.Secret, error) {
	confName := fmt.Sprintf("crt-%s-%s", nodeName, common.RandString(9))
	certName := fmt.Sprintf(`%s.%s`, ns, nodeName)
	certPEM, err := s.PKI.SignClientCertificate(certName, models.AltNames{})
	if err != nil {
		return nil, errors.Trace(err)
	}

	ca, err := s.PKI.GetCA()
	if err != nil {
		return nil, errors.Trace(err)
	}
	srt := &specV1.Secret{
		Name:      confName,
		Namespace: ns,
		Labels: map[string]string{
			common.LabelAppName:  appName,
			common.LabelNodeName: nodeName,
			specV1.SecretLabel:   specV1.SecretCertificate,
			common.LabelSystem:   "true",
		},
		Data: map[string][]byte{
			"client.pem": certPEM.CertPEM,
			"client.key": certPEM.KeyPEM,
			"ca.pem":     ca,
		},
		Annotations: map[string]string{
			common.AnnotationPkiCertID: certPEM.CertId,
		},
		System: true,
	}
	return s.Secret.Create(ns, srt)
}

func (s *InitServiceImpl) genConfig(ns, template string, params map[string]interface{}) (*specV1.Configuration, error) {
	config := &specV1.Configuration{}
	err := s.TemplateService.UnmarshalTemplate(template, params, config)
	if err != nil {
		return nil, errors.Trace(err)
	}
	conf, err := s.Config.Create(ns, config)
	if err != nil {
		res, err := s.Config.Get(ns, config.Name, "")
		if err != nil {
			return nil, errors.Trace(err)
		}
		conf = res
	}
	return conf, nil
}

func (s *InitServiceImpl) genApp(ns, template string, params map[string]interface{}) (*specV1.Application, error) {
	application := &specV1.Application{}
	err := s.TemplateService.UnmarshalTemplate(template, params, application)
	if err != nil {
		return nil, errors.Trace(err)
	}
	app, err := s.App.Create(ns, application)
	if err != nil {
		res, err := s.App.Get(ns, application.Name, "")
		if err != nil {
			return nil, errors.Trace(err)
		}
		app = res
	}
	return app, nil
}
