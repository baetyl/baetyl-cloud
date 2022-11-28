package service

import (
	"encoding/json"
	"fmt"

	"github.com/baetyl/baetyl-go/v2/context"
	"github.com/baetyl/baetyl-go/v2/errors"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/models"
)

//go:generate mockgen -destination=../mock/service/system_app.go -package=service github.com/baetyl/baetyl-cloud/v2/service SystemAppService

const (
	templateInitConfYaml   = "baetyl-init-conf.yml"
	templateInitAppYaml    = "baetyl-init-app.yml"
	templateCoreAppYaml    = "baetyl-core-app.yml"
	templateFuncConfYaml   = "baetyl-function-conf.yml"
	templateFuncAppYaml    = "baetyl-function-app.yml"
	templateBrokerConfYaml = "baetyl-broker-conf.yml"
	templateBrokerAppYaml  = "baetyl-broker-app.yml"
	templateRuleConfYaml   = "baetyl-rule-conf.yml"
	templateRuleAppYaml    = "baetyl-rule-app.yml"
	templateEkuiperAppYaml = "baetyl-ekuiper-app.yml"
)

var (
	HookNamePopulateParams        = "populateParams"
	HookNamePopulateOptAppsParams = "populateOptAppsParams"
	HookNameGenAppsByOption       = "genAppsByOption"
	HookNameGenSyncExtResource    = "genSyncExtResource"
)

type SystemAppService interface {
	GenApps(tx interface{}, ns string, nodeName *specV1.Node) ([]*specV1.Application, error)
	GetOptionalApps() []string
	GenOptionalApps(tx interface{}, ns string, node *specV1.Node, apps []string) ([]*specV1.Application, error)
}

type GenAppsByOption func(tx interface{}, ns string, node *specV1.Node, params map[string]interface{}) ([]*specV1.Application, error)
type GenAppFunc func(tx interface{}, ns string, node *specV1.Node, params map[string]interface{}) (*specV1.Application, error)
type HandlerPopulateParams func(tx interface{}, ns string, params map[string]interface{}) error
type HandlerPopulateOptAppsParams func(tx interface{}, ns string, params map[string]interface{}, appAlias []string) error
type GenSyncExtResource func(tx interface{}, ns string, node *specV1.Node, params map[string]interface{}) error

type SystemAppServiceImpl struct {
	cfg             *config.CloudConfig
	Property        PropertyService
	TemplateService TemplateService
	PKI             PKIService
	*AppCombinedService
	Hooks            map[string]interface{}
	OptionalAppFuncs map[string]GenAppFunc
}

func NewSystemAppService(config *config.CloudConfig) (SystemAppService, error) {
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
	pki, err := NewPKIService(config)
	if err != nil {
		return nil, errors.Trace(err)
	}
	acs, err := NewAppCombinedService(config)
	if err != nil {
		return nil, errors.Trace(err)
	}

	systemService := &SystemAppServiceImpl{
		cfg:                config,
		Property:           propertyService,
		TemplateService:    templateService,
		PKI:                pki,
		AppCombinedService: acs,
		Hooks:              map[string]interface{}{},
	}

	systemService.OptionalAppFuncs = map[string]GenAppFunc{
		specV1.BaetylFunction: systemService.genFunctionApp,
		specV1.BaetylRule:     systemService.genRuleApp,
		specV1.BaetylEkuiper:  systemService.genEkuiperApp,
	}
	return systemService, nil
}

func (s *SystemAppServiceImpl) GenApps(tx interface{}, ns string, node *specV1.Node) ([]*specV1.Application, error) {
	params := map[string]interface{}{
		"Namespace":                  ns,
		"NodeName":                   node.Name,
		"NodeMode":                   node.NodeMode,
		"AppMode":                    node.NodeMode,
		context.KeyBaetylHostPathLib: "{{." + context.KeyBaetylHostPathLib + "}}",
		"GPUStats":                   node.Accelerator != "",
		"DiskNetStats":               node.NodeMode == context.RunModeKube,
		"QPSStats":                   node.NodeMode == context.RunModeKube,
	}
	if handler, ok := s.Hooks[HookNamePopulateParams]; ok {
		err := handler.(HandlerPopulateParams)(tx, ns, params)
		if err != nil {
			return nil, errors.Trace(err)
		}
	}

	coreAppName := fmt.Sprintf("%s-%s", specV1.BaetylCore, common.RandString(9))
	// create secret
	cert, err := s.genNodeCerts(tx, ns, node.Name, coreAppName)
	if err != nil {
		return nil, errors.Trace(err)
	}
	params["NodeCertName"] = cert.Name
	params["NodeCertVersion"] = cert.Version
	params["CoreAppName"] = coreAppName

	if gen, ok := s.Hooks[HookNameGenSyncExtResource]; ok {
		err := gen.(GenSyncExtResource)(tx, ns, node, params)
		if err != nil {
			return nil, errors.Trace(err)
		}
	}

	var apps []*specV1.Application
	ca, err := s.genCoreApp(tx, ns, node.Name, params)
	if err != nil {
		return nil, errors.Trace(err)
	}
	ia, err := s.genInitApp(tx, ns, node.Name, params)
	if err != nil {
		return nil, errors.Trace(err)
	}
	ba, err := s.genBrokerApp(tx, ns, node.Name, params)
	if err != nil {
		return nil, errors.Trace(err)
	}
	apps = append(apps, ca, ia, ba)

	optionalSysApps, err := s.GenOptionalApps(tx, ns, node, node.SysApps)
	if err != nil {
		return nil, err
	}
	apps = append(apps, optionalSysApps...)

	if gen, ok := s.Hooks[HookNameGenAppsByOption]; ok {
		extApps, err := gen.(GenAppsByOption)(tx, ns, node, params)
		if err != nil {
			return nil, errors.Trace(err)
		}
		apps = append(apps, extApps...)
	}
	return apps, nil
}

func (s *SystemAppServiceImpl) GetOptionalApps() []string {
	var res []string
	for k := range s.OptionalAppFuncs {
		res = append(res, k)
	}
	return res
}

func (s *SystemAppServiceImpl) GenOptionalApps(tx interface{}, ns string, node *specV1.Node, appAlias []string) ([]*specV1.Application, error) {
	params := map[string]interface{}{
		"Namespace":                  ns,
		"NodeName":                   node.Name,
		"NodeMode":                   node.NodeMode,
		"AppMode":                    node.NodeMode,
		context.KeyBaetylHostPathLib: "{{." + context.KeyBaetylHostPathLib + "}}",
	}
	if handler, ok := s.Hooks[HookNamePopulateOptAppsParams]; ok {
		err := handler.(HandlerPopulateOptAppsParams)(tx, ns, params, appAlias)
		if err != nil {
			return nil, errors.Trace(err)
		}
	}

	var apps []*specV1.Application
	for _, v := range appAlias {
		if f, ok := s.OptionalAppFuncs[v]; ok {
			app, err := f(tx, ns, node, params)
			if err != nil {
				return nil, err
			}
			apps = append(apps, app)
		}
	}
	return apps, nil
}

func (s *SystemAppServiceImpl) genCoreApp(tx interface{}, ns, nodeName string, params map[string]interface{}) (*specV1.Application, error) {
	confName := fmt.Sprintf("%s-conf-%s", specV1.BaetylCore, common.RandString(9))
	params["CoreConfName"] = confName
	params["CoreFrequency"] = fmt.Sprintf("%ss", common.DefaultCoreFrequency)
	params["CoreAPIPort"] = common.DefaultCoreAPIPort
	params["AgentPort"] = common.DefaultAgentPort

	// create config
	conf, err := s.GenConfig(tx, ns, TemplateCoreConfYaml, params)
	if err != nil {
		return nil, errors.Trace(err)
	}

	params["CoreConfVersion"] = conf.Version

	// create application
	return s.GenApp(tx, ns, templateCoreAppYaml, params)
}

func (s *SystemAppServiceImpl) genInitApp(tx interface{}, ns, nodeName string, params map[string]interface{}) (*specV1.Application, error) {
	appName := fmt.Sprintf("%s-%s", specV1.BaetylInit, common.RandString(9))
	confName := fmt.Sprintf("%s-conf-%s", specV1.BaetylInit, common.RandString(9))
	params["InitAppName"] = appName
	params["InitConfName"] = confName
	params["AgentPort"] = common.DefaultAgentPort

	// create config
	conf, err := s.GenConfig(tx, ns, templateInitConfYaml, params)
	if err != nil {
		return nil, errors.Trace(err)
	}

	params["InitConfVersion"] = conf.Version
	// create application
	return s.GenApp(tx, ns, templateInitAppYaml, params)
}

func (s *SystemAppServiceImpl) genFunctionApp(tx interface{}, ns string, node *specV1.Node, params map[string]interface{}) (*specV1.Application, error) {
	appName := fmt.Sprintf("%s-%s", specV1.BaetylFunction, common.RandString(9))
	confName := fmt.Sprintf("%s-conf-%s", specV1.BaetylFunction, common.RandString(9))
	// create config
	confMap := map[string]interface{}{
		"Namespace":        ns,
		"NodeName":         node.Name,
		"FunctionAppName":  appName,
		"FunctionConfName": confName,
	}
	for k, v := range params {
		confMap[k] = v
	}
	conf, err := s.GenConfig(tx, ns, templateFuncConfYaml, confMap)
	if err != nil {
		return nil, errors.Trace(err)
	}

	// create application
	appMap := map[string]interface{}{
		"Namespace":           ns,
		"NodeName":            node.Name,
		"FunctionAppName":     appName,
		"FunctionConfName":    conf.Name,
		"FunctionConfVersion": conf.Version,
	}
	for k, v := range params {
		appMap[k] = v
	}
	return s.GenApp(tx, ns, templateFuncAppYaml, appMap)
}

func (s *SystemAppServiceImpl) genBrokerApp(tx interface{}, ns, nodeName string, params map[string]interface{}) (*specV1.Application, error) {
	appName := fmt.Sprintf("%s-%s", specV1.BaetylBroker, common.RandString(9))
	confName := fmt.Sprintf("%s-conf-%s", specV1.BaetylBroker, common.RandString(9))
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
	conf, err := s.GenConfig(tx, ns, templateBrokerConfYaml, confMap)
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
	return s.GenApp(tx, ns, templateBrokerAppYaml, appMap)
}

func (s *SystemAppServiceImpl) genRuleApp(tx interface{}, ns string, node *specV1.Node, params map[string]interface{}) (*specV1.Application, error) {
	appName := fmt.Sprintf("%s-%s", specV1.BaetylRule, common.RandString(9))
	confName := fmt.Sprintf("%s-conf-%s", specV1.BaetylRule, common.RandString(9))
	// create config
	confMap := map[string]interface{}{
		"Namespace":    ns,
		"NodeName":     node.Name,
		"RuleAppName":  appName,
		"RuleConfName": confName,
	}
	for k, v := range params {
		confMap[k] = v
	}
	conf, err := s.GenConfig(tx, ns, templateRuleConfYaml, confMap)
	if err != nil {
		return nil, errors.Trace(err)
	}

	// create application
	appMap := map[string]interface{}{
		"Namespace":       ns,
		"NodeName":        node.Name,
		"RuleAppName":     appName,
		"RuleConfName":    conf.Name,
		"RuleConfVersion": conf.Version,
	}
	for k, v := range params {
		appMap[k] = v
	}
	return s.GenApp(tx, ns, templateRuleAppYaml, appMap)
}

func (s *SystemAppServiceImpl) genEkuiperApp(tx interface{}, ns string, node *specV1.Node, params map[string]interface{}) (*specV1.Application, error) {
	appName := fmt.Sprintf("%s-%s", specV1.BaetylEkuiper, common.RandString(9))

	// create application
	appMap := map[string]interface{}{
		"Namespace":      ns,
		"NodeName":       node.Name,
		"EkuiperAppName": appName,
	}
	for k, v := range params {
		appMap[k] = v
	}
	return s.GenApp(tx, ns, templateEkuiperAppYaml, appMap)
}

func (s *SystemAppServiceImpl) genNodeCerts(tx interface{}, ns, nodeName, appName string) (*specV1.Secret, error) {
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
	return s.Secret.Create(tx, ns, srt)
}

func (s *SystemAppServiceImpl) getCoreConfig(ns, nodeName string, params map[string]interface{}) ([]byte, error) {
	params["Namespace"] = ns
	params["NodeName"] = nodeName
	return s.TemplateService.ParseTemplate(TemplateCoreConfYaml, params)
}

func (s *SystemAppServiceImpl) GenConfig(tx interface{}, ns, template string, params map[string]interface{}) (*specV1.Configuration, error) {
	cfg := &specV1.Configuration{}
	err := s.TemplateService.UnmarshalTemplate(template, params, cfg)
	if err != nil {
		return nil, errors.Trace(err)
	}
	conf, err := s.Config.Create(tx, ns, cfg)
	if err != nil {
		res, err := s.Config.Get(ns, cfg.Name, "")
		if err != nil {
			return nil, errors.Trace(err)
		}
		conf = res
	}
	return conf, nil
}

func (s *SystemAppServiceImpl) GenApp(tx interface{}, ns, template string, params map[string]interface{}) (*specV1.Application, error) {
	// Create registry secret for system app
	params["RegistryAuth"] = ""
	if registryAuth, err := s.Property.GetPropertyValue(common.RegistryAuth); err == nil {
		secretVersion, createErr := s.GenSystemRegistry(tx, ns, registryAuth)
		if createErr != nil {
			return nil, errors.Trace(createErr)
		}
		params["RegistryAuth"] = common.RegistryAuth
		params["RegistryAuthVersion"] = secretVersion
	}
	application := &specV1.Application{}
	err := s.TemplateService.UnmarshalTemplate(template, params, application)
	if err != nil {
		return nil, errors.Trace(err)
	}
	app, err := s.App.Create(tx, ns, application)
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

func (s *SystemAppServiceImpl) GenSystemRegistry(tx interface{}, ns, registryAuth string) (string, error) {
	var registryModel models.Registry
	err := json.Unmarshal([]byte(registryAuth), &registryModel)
	if err != nil {
		return "", errors.Trace(err)
	}
	registrySecret := &specV1.Secret{
		Name:      common.RegistryAuth,
		Namespace: ns,
		Labels:    map[string]string{specV1.SecretLabel: specV1.SecretRegistry, common.ResourceInvisible: "true", common.LabelSystem: "true"},
		System:    true,
		Data: map[string][]byte{
			"address":  []byte(registryModel.Address),
			"username": []byte(registryModel.Username),
			"password": []byte(registryModel.Password),
		},
	}
	rs, err := s.Secret.GetTx(tx, ns, common.RegistryAuth, "")
	if err != nil {
		if e, ok := err.(errors.Coder); ok && e.Code() == common.ErrResourceNotFound {
			rs, err = s.Secret.Create(tx, ns, registrySecret)
			if err != nil {
				return "", errors.Trace(err)
			}
		} else {
			return "", errors.Trace(err)
		}
	}
	return rs.Version, nil
}
