package service

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/baetyl/baetyl-go/v2/context"
	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/baetyl/baetyl-go/v2/log"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/models"
)

//go:generate mockgen -destination=../mock/service/init.go -package=service github.com/baetyl/baetyl-cloud/v2/service InitService

const (
	InfoName      = "n"
	InfoNamespace = "ns"
	InfoExpiry    = "e"
)

const (
	templateInitConfYaml       = "baetyl-init-conf.yml"
	templateInitAppYaml        = "baetyl-init-app.yml"
	templateCoreAppYaml        = "baetyl-core-app.yml"
	templateFuncConfYaml       = "baetyl-function-conf.yml"
	templateFuncAppYaml        = "baetyl-function-app.yml"
	templateBrokerConfYaml     = "baetyl-broker-conf.yml"
	templateBrokerAppYaml      = "baetyl-broker-app.yml"
	templateRuleConfYaml       = "baetyl-rule-conf.yml"
	templateRuleAppYaml        = "baetyl-rule-app.yml"
	templateInitDeploymentYaml = "baetyl-init-deployment.yml"
	templateBaetylInstallShell = "baetyl-install.sh"

	TemplateCoreConfYaml      = "baetyl-core-conf.yml"
	TemplateBaetylInitCommand = "baetyl-init-command"
)

var (
	CmdExpirationInSeconds        = int64(60 * 60)
	HookNamePopulateParams        = "populateParams"
	HookNamePopulateOptAppsParams = "populateOptAppsParams"
	HookNameGenAppsByOption       = "genAppsByOption"
)

type HandlerPopulateParams func(ns string, params map[string]interface{}) error
type HandlerPopulateOptAppsParams func(ns string, params map[string]interface{}, appAlias []string) error
type GetInitResource func(ns, nodeName string, params map[string]interface{}) ([]byte, error)
type GenAppsByOption func(ns string, node *specV1.Node, params map[string]interface{}) ([]*specV1.Application, error)
type GenAppFunc func(ns, nodeName string, params map[string]interface{}) (*specV1.Application, error)

// InitService
type InitService interface {
	GetResource(ns, nodeName, resourceName string, params map[string]interface{}) (interface{}, error)
	GetOptionalApps() []string

	GenApps(ns string, nodeName *specV1.Node) ([]*specV1.Application, error)
	GenOptionalApps(ns string, nodeName string, apps []string) ([]*specV1.Application, error)
}

type InitServiceImpl struct {
	cfg             *config.CloudConfig
	AuthService     AuthService
	NodeService     NodeService
	Property        PropertyService
	TemplateService TemplateService
	*AppCombinedService
	PKI              PKIService
	Hooks            map[string]interface{}
	ResourceMapFunc  map[string]GetInitResource
	OptionalAppFuncs map[string]GenAppFunc
	log              *log.Logger
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
	propertyService, err := NewPropertyService(config)
	if err != nil {
		return nil, errors.Trace(err)
	}
	moduleService, err := NewModuleService(config)
	if err != nil {
		return nil, err
	}
	templateService, err := NewTemplateService(config, map[string]interface{}{
		"GetProperty":      propertyService.GetPropertyValue,
		"RandString":       common.RandString,
		"GetModuleImage":   moduleService.GetLatestModuleImage,
		"GetModuleProgram": moduleService.GetLatestModuleProgram,
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
		Property:           propertyService,
		TemplateService:    templateService,
		AppCombinedService: acs,
		PKI:                pki,
		Hooks:              map[string]interface{}{},
		ResourceMapFunc:    map[string]GetInitResource{},
		log:                log.L().With(log.Any("service", "init")),
	}
	initService.ResourceMapFunc[templateInitDeploymentYaml] = initService.getInitDeploymentYaml
	initService.ResourceMapFunc[TemplateBaetylInitCommand] = initService.GetInitCommand
	initService.ResourceMapFunc[TemplateCoreConfYaml] = initService.getCoreConfig
	initService.ResourceMapFunc[templateBaetylInstallShell] = initService.getInstallShell

	initService.OptionalAppFuncs = map[string]GenAppFunc{
		specV1.BaetylFunction: initService.genFunctionApp,
		specV1.BaetylRule:     initService.genRuleApp,
	}

	return initService, nil
}

func (s *InitServiceImpl) GetResource(ns, nodeName, resourceName string, params map[string]interface{}) (interface{}, error) {
	if handler, ok := s.ResourceMapFunc[resourceName]; ok {
		if params == nil {
			params = map[string]interface{}{}
		}
		return handler(ns, nodeName, params)
	}
	return nil, common.Error(
		common.ErrResourceNotFound,
		common.Field("type", "resource"),
		common.Field("name", resourceName))
}

func (s *InitServiceImpl) getInitDeploymentYaml(ns, nodeName string, params map[string]interface{}) ([]byte, error) {
	init, err := s.GetAppFromDesire(ns, nodeName, specV1.BaetylInit, true)
	if err != nil {
		s.log.Warn("failed to get init app from desire, the init module will use the default name and version",
			log.Any("InitAppName", "baetyl-init"), log.Any("InitVersion", "1"), log.Error(err))

		core, err := s.GetAppFromDesire(ns, nodeName, specV1.BaetylCore, true)
		if err != nil {
			return nil, errors.Trace(err)
		}

		// for node cert
		init = &specV1.Application{
			Name:      specV1.BaetylInit,
			Namespace: core.Namespace,
			Version:   "1",
			Volumes:   core.Volumes,
		}
	}
	cert, err := s.GetNodeCert(init)
	if err != nil {
		return nil, errors.Trace(err)
	}

	node, err := s.NodeService.Get(ns, nodeName)
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
	params["EdgeNamespace"] = context.EdgeNamespace()
	params["EdgeSystemNamespace"] = context.EdgeSystemNamespace()
	params["InitAppName"] = init.Name
	params["InitVersion"] = init.Version
	params["GPUStats"] = node.Accelerator == specV1.NVAccelerator

	return s.TemplateService.ParseTemplate(templateInitDeploymentYaml, params)
}

func (s *InitServiceImpl) GetNodeCert(app *specV1.Application) (*specV1.Secret, error) {
	certName := ""
	for _, vol := range app.Volumes {
		if vol.Name == "node-cert" || vol.Name == "cert-sync" {
			certName = vol.Secret.Name
			break
		}
	}
	cert, _ := s.Secret.Get(app.Namespace, certName, "")
	if cert == nil {
		return nil, common.Error(
			common.ErrResourceNotFound,
			common.Field("type", "secret"),
			common.Field("name", certName),
			common.Field("namespace", app.Namespace))
	}
	return cert, nil
}

func (s *InitServiceImpl) getInstallShell(ns, nodeName string, params map[string]interface{}) ([]byte, error) {
	params["DBPath"] = "/var/lib/baetyl"
	data, err := s.TemplateService.ParseTemplate(templateBaetylInstallShell, params)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return data, nil
}

func (s *InitServiceImpl) GetInitCommand(ns, nodeName string, params map[string]interface{}) ([]byte, error) {
	info := map[string]interface{}{
		InfoNamespace: ns,
		InfoName:      nodeName,
		InfoExpiry:    time.Now().Unix() + CmdExpirationInSeconds,
	}
	initCommand, err := s.Property.GetPropertyValue(TemplateBaetylInitCommand)
	if err != nil {
		return nil, err
	}
	token, err := s.AuthService.GenToken(info)
	if err != nil {
		return nil, err
	}
	params["Token"] = token
	data, err := s.TemplateService.Execute("setup-command", initCommand, params)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *InitServiceImpl) GetAppFromDesire(ns, nodeName, moduleName string, isSys bool) (*specV1.Application, error) {
	shadowDesire, err := s.NodeService.GetDesire(ns, nodeName)
	if err != nil {
		return nil, errors.Trace(err)
	}
	apps := shadowDesire.AppInfos(isSys)
	for _, appInfo := range apps {
		if strings.Contains(appInfo.Name, moduleName) {
			return s.App.Get(ns, appInfo.Name, "")
		}
	}
	return nil, common.Error(
		common.ErrResourceNotFound,
		common.Field("type", "app"),
		common.Field("name", nodeName),
		common.Field("namespace", ns))
}

func (s *InitServiceImpl) GenApps(ns string, node *specV1.Node) ([]*specV1.Application, error) {
	params := map[string]interface{}{
		"Namespace":                  ns,
		"NodeName":                   node.Name,
		context.KeyBaetylHostPathLib: "{{." + context.KeyBaetylHostPathLib + "}}",
		"GPUStats":                   node.Accelerator == specV1.NVAccelerator,
	}
	if handler, ok := s.Hooks[HookNamePopulateParams]; ok {
		err := handler.(HandlerPopulateParams)(ns, params)
		if err != nil {
			return nil, errors.Trace(err)
		}
	}

	var apps []*specV1.Application
	ca, err := s.genCoreApp(ns, node.Name, params)
	if err != nil {
		return nil, errors.Trace(err)
	}
	ia, err := s.genInitApp(ns, node.Name, params)
	if err != nil {
		return nil, errors.Trace(err)
	}
	ba, err := s.genBrokerApp(ns, node.Name, params)
	if err != nil {
		return nil, errors.Trace(err)
	}
	apps = append(apps, ca, ia, ba)

	optionalSysApps, err := s.GenOptionalApps(ns, node.Name, node.SysApps)
	if err != nil {
		return nil, err
	}
	apps = append(apps, optionalSysApps...)

	if gen, ok := s.Hooks[HookNameGenAppsByOption]; ok {
		extApps, err := gen.(GenAppsByOption)(ns, node, params)
		if err != nil {
			return nil, errors.Trace(err)
		}
		apps = append(apps, extApps...)
	}
	return apps, nil
}

func (s *InitServiceImpl) GenOptionalApps(ns string, node string, appAlias []string) ([]*specV1.Application, error) {
	params := map[string]interface{}{
		"Namespace":                  ns,
		"NodeName":                   node,
		context.KeyBaetylHostPathLib: "{{." + context.KeyBaetylHostPathLib + "}}",
	}
	if handler, ok := s.Hooks[HookNamePopulateOptAppsParams]; ok {
		err := handler.(HandlerPopulateOptAppsParams)(ns, params, appAlias)
		if err != nil {
			return nil, errors.Trace(err)
		}
	}

	var apps []*specV1.Application
	for _, v := range appAlias {
		if f, ok := s.OptionalAppFuncs[v]; ok {
			app, err := f(ns, node, params)
			if err != nil {
				return nil, err
			}
			apps = append(apps, app)
		}
	}
	return apps, nil
}

func (s *InitServiceImpl) GetOptionalApps() []string {
	var res []string
	for k := range s.OptionalAppFuncs {
		res = append(res, k)
	}
	return res
}

func (s *InitServiceImpl) genCoreApp(ns, nodeName string, params map[string]interface{}) (*specV1.Application, error) {
	appName := fmt.Sprintf("baetyl-core-%s", common.RandString(9))
	confName := fmt.Sprintf("baetyl-core-conf-%s", common.RandString(9))
	params["CoreAppName"] = appName
	params["CoreConfName"] = confName
	params["CoreFrequency"] = fmt.Sprintf("%ss", common.DefaultCoreFrequency)
	params["CoreAPIPort"] = common.DefaultCoreAPIPort

	// create config
	conf, err := s.GenConfig(ns, TemplateCoreConfYaml, params)
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
	return s.GenApp(ns, templateCoreAppYaml, params)
}

func (s *InitServiceImpl) genInitApp(ns, nodeName string, params map[string]interface{}) (*specV1.Application, error) {
	appName := fmt.Sprintf("baetyl-init-%s", common.RandString(9))
	confName := fmt.Sprintf("baetyl-init-conf-%s", common.RandString(9))
	params["InitAppName"] = appName
	params["InitConfName"] = confName

	// create config
	conf, err := s.GenConfig(ns, templateInitConfYaml, params)
	if err != nil {
		return nil, errors.Trace(err)
	}

	params["InitConfVersion"] = conf.Version
	// create application
	return s.GenApp(ns, templateInitAppYaml, params)
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
	conf, err := s.GenConfig(ns, templateFuncConfYaml, confMap)
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
	return s.GenApp(ns, templateFuncAppYaml, appMap)
}

func (s *InitServiceImpl) genBrokerApp(ns, nodeName string, params map[string]interface{}) (*specV1.Application, error) {
	appName := fmt.Sprintf("baetyl-broker-%s", common.RandString(9))
	confName := fmt.Sprintf("baetyl-broker-conf-%s", common.RandString(9))
	// create config
	confMap := map[string]interface{}{
		"Namespace":      ns,
		"NodeName":       nodeName,
		"BrokerAppName":  appName,
		"BrokerConfName": confName,
	}
	for k, v := range params {
		confMap[k] = v
	}
	conf, err := s.GenConfig(ns, templateBrokerConfYaml, confMap)
	if err != nil {
		return nil, errors.Trace(err)
	}

	// create application
	appMap := map[string]interface{}{
		"Namespace":         ns,
		"NodeName":          nodeName,
		"BrokerAppName":     appName,
		"BrokerConfName":    conf.Name,
		"BrokerConfVersion": conf.Version,
	}
	for k, v := range params {
		appMap[k] = v
	}
	return s.GenApp(ns, templateBrokerAppYaml, appMap)
}

func (s *InitServiceImpl) genRuleApp(ns, nodeName string, params map[string]interface{}) (*specV1.Application, error) {
	appName := fmt.Sprintf("baetyl-rule-%s", common.RandString(9))
	confName := fmt.Sprintf("baetyl-rule-conf-%s", common.RandString(9))
	// create config
	confMap := map[string]interface{}{
		"Namespace":    ns,
		"NodeName":     nodeName,
		"RuleAppName":  appName,
		"RuleConfName": confName,
	}
	for k, v := range params {
		confMap[k] = v
	}
	conf, err := s.GenConfig(ns, templateRuleConfYaml, confMap)
	if err != nil {
		return nil, errors.Trace(err)
	}

	// create application
	appMap := map[string]interface{}{
		"Namespace":       ns,
		"NodeName":        nodeName,
		"RuleAppName":     appName,
		"RuleConfName":    conf.Name,
		"RuleConfVersion": conf.Version,
	}
	for k, v := range params {
		appMap[k] = v
	}
	return s.GenApp(ns, templateRuleAppYaml, appMap)
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
			common.LabelAppName:      appName,
			common.LabelNodeName:     nodeName,
			specV1.SecretLabel:       specV1.SecretConfig,
			common.LabelSystem:       "true",
			common.ResourceInvisible: "true",
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

func (s *InitServiceImpl) getCoreConfig(ns, nodeName string, params map[string]interface{}) ([]byte, error) {
	params["Namespace"] = ns
	params["NodeName"] = nodeName
	return s.TemplateService.ParseTemplate(TemplateCoreConfYaml, params)
}

func (s *InitServiceImpl) GenConfig(ns, template string, params map[string]interface{}) (*specV1.Configuration, error) {
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

func (s *InitServiceImpl) GenApp(ns, template string, params map[string]interface{}) (*specV1.Application, error) {
	application := &specV1.Application{}
	err := s.TemplateService.UnmarshalTemplate(template, params, application)
	if err != nil {
		return nil, errors.Trace(err)
	}
	app, err := s.App.Create(ns, application)
	if err != nil {
		var res *specV1.Application
		res, err = s.App.Get(ns, application.Name, "")
		if err != nil {
			return nil, errors.Trace(err)
		}
		app = res
	}
	return app, nil
}
