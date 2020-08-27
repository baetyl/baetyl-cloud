package api

import (
	"fmt"
	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/service"
	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/baetyl/baetyl-go/v2/log"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
)

var (
	certInfo = models.CertStorage{
		CertName: "client.pem",
		KeyName:  "client.key",
	}
)

// todo more general, initialize according to list
func (api *API) GenSysApp(nodeName, ns string, appList []common.SystemApplication) ([]specV1.Application, error) {
	var apps []specV1.Application
	for _, appName := range appList {
		app := specV1.Application{}
		switch appName {
		case common.BaetylCore:
			res, err := api.GenCoreApp(ns, nodeName)
			if err != nil {
				return nil, err
			}
			app = *res
		case common.BaetylFunction:
			res, err := api.GenFunctionApp(ns, nodeName)
			if err != nil {
				return nil, err
			}
			app = *res
			//case common.BaetylBroker:
			//	res, err := api.GenBrokerApp(nodeName, ns, true)
			//	if err != nil {
			//		return nil, err
			//	}
			//	app = *res
			//case common.BaetylRule:
			//	res, err := api.GenRuleApp(nodeName, ns, true)
			//	if err != nil {
			//		return nil, err
			//	}
			//	app = *res
		}
		apps = append(apps, app)
	}
	return apps, nil
}

func (api *API) GenCoreApp(ns, nodeName string) (*specV1.Application, error) {
	syncAddress, err := api.sysConfigService.GetSysConfig("address", common.AddressNode)
	if err != nil {
		log.L().Error("GenCoreApp", log.Any("GetSysConfig", "imageConf get nil"))
		return nil, err
	}
	appName := fmt.Sprintf("%s-%s", common.BaetylCore, common.RandString(9))
	// create config
	confMap := map[string]string{
		"Namespace": ns,
		"NodeName":  nodeName,
		"AppName":   appName,
		"SyncAddr":  syncAddress.Value,
		"ConfName":  fmt.Sprintf("%s-%s-config-%s", common.BaetylCore, nodeName, common.RandString(9)),
	}
	conf, err := api.genConfig(ns, service.TemplateCoreConf, confMap)
	if err != nil {
		return nil, err
	}

	// create secret
	cert, err := api.genNodeCerts(appName, nodeName, ns, common.BaetylCore, true)
	if err != nil {
		return nil, err
	}

	// create application
	appMap := map[string]string{
		"Namespace":   ns,
		"AppName":     appName,
		"NodeName":    nodeName,
		"CertName":    cert.Name,
		"CertVersion": cert.Version,
		"ConfName":    conf.Name,
		"ConfVersion": conf.Version,
	}
	return api.genApp(ns, service.TemplateCoreApp, appMap)
}

func (api *API) GenFunctionApp(ns, nodeName string) (*specV1.Application, error) {
	appName := fmt.Sprintf("%s-%s", common.BaetylFunction, common.RandString(9))
	// create config
	confMap := map[string]string{
		"Namespace": ns,
		"AppName":   appName,
		"NodeName":  nodeName,
		"ConfName":  fmt.Sprintf("%s-%s-config-%s", common.BaetylFunction, nodeName, common.RandString(9)),
	}
	conf, err := api.genConfig(ns, service.TemplateFuncConf, confMap)
	if err != nil {
		return nil, err
	}

	// create application
	appMap := map[string]string{
		"Namespace":   ns,
		"AppName":     appName,
		"NodeName":    nodeName,
		"ConfName":    conf.Name,
		"ConfVersion": conf.Version,
	}
	return api.genApp(ns, service.TemplateFuncApp, appMap)
}

//
//func (api *API) GenBrokerApp(nodeName, ns string, isSys bool) (*specV1.Application, error) {
//	// get sys config
//	imageConf, err := api.sysConfigService.GetSysConfig(common.BaetylModule, string(common.BaetylBroker))
//	if err != nil {
//		return nil, err
//	}
//
//	appName := fmt.Sprintf("%s-%s", common.BaetylBroker, common.RandString(9))
//	// create config
//	confMap := map[string]string{
//		"AppName":    appName,
//		"NodeName":   nodeName,
//		"Namespace":  ns,
//		"ConfigName": fmt.Sprintf("%s-%s-config-%s", common.BaetylBroker, nodeName, common.RandString(9)),
//	}
//	conf, err := api.genConfig(ns, common.TemplateJsonConfigBroker, confMap, isSys)
//	if err != nil {
//		return nil, err
//	}
//
//	// create application
//	appMap := map[string]string{
//		"AppName":       appName,
//		"Image":         imageConf.Value,
//		"NodeName":      nodeName,
//		"Namespace":     ns,
//		"ConfigName":    conf.Name,
//		"ConfigVersion": conf.Version,
//		"AppType":       common.ContainerApp,
//	}
//	return api.genApp(ns, common.TemplateJsonAppBroker, appMap, isSys)
//}
//
//func (api *API) GenRuleApp(nodeName, ns string, isSys bool) (*specV1.Application, error) {
//	// get sys config
//	imageConf, err := api.sysConfigService.GetSysConfig(common.BaetylModule, string(common.BaetylRule))
//	if err != nil {
//		return nil, err
//	}
//
//	appName := fmt.Sprintf("%s-%s", common.BaetylRule, common.RandString(9))
//	// create config
//	confMap := map[string]string{
//		"AppName":    appName,
//		"NodeName":   nodeName,
//		"Namespace":  ns,
//		"ConfigName": fmt.Sprintf("%s-%s-config-%s", common.BaetylRule, nodeName, common.RandString(9)),
//	}
//	conf, err := api.genConfig(ns, common.TemplateJsonConfigRule, confMap, isSys)
//	if err != nil {
//		return nil, err
//	}
//
//	// create application
//	appMap := map[string]string{
//		"AppName":       appName,
//		"Image":         imageConf.Value,
//		"NodeName":      nodeName,
//		"Namespace":     ns,
//		"ConfigName":    conf.Name,
//		"ConfigVersion": conf.Version,
//		"AppType":       common.ContainerApp,
//	}
//	return api.genApp(ns, common.TemplateJsonAppRule, appMap, isSys)
//}

func (api *API) genNodeCerts(appName, nodeName, ns string, module common.SystemApplication, isSys bool) (*specV1.Secret, error) {
	name := "sync-" + nodeName + "-" + string(module) + "-" + common.RandString(9)
	cerName := fmt.Sprintf(`%s.%s`, ns, nodeName)
	certPEM, err := api.pkiService.SignClientCertificate(cerName, models.AltNames{})
	if err != nil {
		return nil, err
	}

	ca, err := api.pkiService.GetCA()
	if err != nil {
		return nil, err
	}
	s := &specV1.Secret{
		Name:      name,
		Namespace: ns,
		Labels: map[string]string{
			common.LabelAppName:  appName,
			common.LabelNodeName: nodeName,
			specV1.SecretLabel:   specV1.SecretCertificate,
			common.LabelSystem:   "true",
		},
		Data: map[string][]byte{
			certInfo.CertName: certPEM.CertPEM,
			certInfo.KeyName:  certPEM.KeyPEM,
			"ca.pem":          ca,
		},
		Annotations: map[string]string{
			common.AnnotationPkiCertID: certPEM.CertId,
		},
		System: isSys,
	}
	return api.secretService.Create(ns, s)
}

func (api *API) genConfig(ns, template string, params map[string]string) (*specV1.Configuration, error) {
	config := &specV1.Configuration{}
	err := api.templateService.ParseSystemTemplate(template, params, config)
	if err != nil {
		return nil, err
	}
	conf, err := api.configService.Create(ns, config)
	if err != nil {
		log.L().Error("API", log.Any("func", "genApp"), log.Any("err", err.Error()))
		res, err := api.configService.Get(ns, config.Name, "")
		if err != nil {
			return nil, err
		}
		conf = res
	}
	return conf, nil
}

func (api *API) genApp(ns, template string, params map[string]string) (*specV1.Application, error) {
	application := &specV1.Application{}
	err := api.templateService.ParseSystemTemplate(template, params, application)
	if err != nil {
		return nil, errors.Trace(err)
	}
	app, err := api.applicationService.Create(ns, application)
	if err != nil {
		log.L().Error("API", log.Any("func", "genApp"), log.Any("err", err.Error()))
		res, err := api.applicationService.Get(ns, application.Name, "")
		if err != nil {
			return nil, err
		}
		app = res
	}

	err = api.updateNodeAndAppIndex(ns, app)
	if err != nil {
		return nil, err
	}

	return app, nil
}
