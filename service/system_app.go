package service

import (
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
)

var (
	HookNamePopulateParams        = "populateParams"
	HookNamePopulateOptAppsParams = "populateOptAppsParams"
	HookNameGenAppsByOption       = "genAppsByOption"
)

type SystemAppService interface {
	GenApps(ns string, nodeName *specV1.Node) ([]*specV1.Application, error)
	GetOptionalApps() []string
	GenOptionalApps(ns string, nodeName string, apps []string) ([]*specV1.Application, error)
}

type GenAppsByOption func(ns string, node *specV1.Node, params map[string]interface{}) ([]*specV1.Application, error)
type GenAppFunc func(ns, nodeName string, params map[string]interface{}) (*specV1.Application, error)
type HandlerPopulateParams func(ns string, params map[string]interface{}) error
type HandlerPopulateOptAppsParams func(ns string, params map[string]interface{}, appAlias []string) error

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
	}
	return systemService, nil
}

func (s *SystemAppServiceImpl) GenApps(ns string, node *specV1.Node) ([]*specV1.Application, error) {
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

func (s *SystemAppServiceImpl) GetOptionalApps() []string {
	var res []string
	for k := range s.OptionalAppFuncs {
		res = append(res, k)
	}
	return res
}

func (s *SystemAppServiceImpl) GenOptionalApps(ns string, node string, appAlias []string) ([]*specV1.Application, error) {
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

func (s *SystemAppServiceImpl) genCoreApp(ns, nodeName string, params map[string]interface{}) (*specV1.Application, error) {
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

func (s *SystemAppServiceImpl) genInitApp(ns, nodeName string, params map[string]interface{}) (*specV1.Application, error) {
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

func (s *SystemAppServiceImpl) genFunctionApp(ns, nodeName string, params map[string]interface{}) (*specV1.Application, error) {
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

func (s *SystemAppServiceImpl) genBrokerApp(ns, nodeName string, params map[string]interface{}) (*specV1.Application, error) {
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

func (s *SystemAppServiceImpl) genRuleApp(ns, nodeName string, params map[string]interface{}) (*specV1.Application, error) {
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

func (s *SystemAppServiceImpl) genNodeCerts(ns, nodeName, appName string) (*specV1.Secret, error) {
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

func (s *SystemAppServiceImpl) getCoreConfig(ns, nodeName string, params map[string]interface{}) ([]byte, error) {
	params["Namespace"] = ns
	params["NodeName"] = nodeName
	return s.TemplateService.ParseTemplate(TemplateCoreConfYaml, params)
}

func (s *SystemAppServiceImpl) GenConfig(ns, template string, params map[string]interface{}) (*specV1.Configuration, error) {
	cfg := &specV1.Configuration{}
	err := s.TemplateService.UnmarshalTemplate(template, params, cfg)
	if err != nil {
		return nil, errors.Trace(err)
	}
	conf, err := s.Config.Create(ns, cfg)
	if err != nil {
		res, err := s.Config.Get(ns, cfg.Name, "")
		if err != nil {
			return nil, errors.Trace(err)
		}
		conf = res
	}
	return conf, nil
}

func (s *SystemAppServiceImpl) GenApp(ns, template string, params map[string]interface{}) (*specV1.Application, error) {
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
