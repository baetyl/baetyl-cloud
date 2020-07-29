package api

import (
	"encoding/json"
	"fmt"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-go/v2/log"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
)

var (
	certInfo = models.CertStorage{
		CertName: "client.pem",
		KeyName:  "client.key",
	}
)

func (api *API) GenSysApp(nodeName, ns string, appList []common.SystemApplication) ([]specV1.Application, error) {
	var apps []specV1.Application
	isSysApp := true
	for _, appName := range appList {
		app := specV1.Application{}
		switch appName {
		case common.BaetylCore:
			res, err := api.GenCoreApp(nodeName, ns, isSysApp)
			if err != nil {
				return nil, err
			}
			app = *res
		case common.BaetylFunction:
			res, err := api.GenFunctionApp(nodeName, ns, isSysApp)
			if err != nil {
				return nil, err
			}
			app = *res
		}
		apps = append(apps, app)
	}
	return apps, nil
}

func (api *API) GenCoreApp(nodeName, ns string, isSys bool) (*specV1.Application, error) {
	// get sys config
	imageConf, err := api.sysConfigService.GetSysConfig(common.BaetylModule, string(common.BaetylCore))
	if err != nil {
		log.L().Error("GenCoreApp", log.Any("GetSysConfig", "imageConf get nil"))
		return nil, err
	}
	nodeAddress, err := api.sysConfigService.GetSysConfig("address", common.AddressNode)
	if err != nil {
		log.L().Error("GenCoreApp", log.Any("GetSysConfig", "imageConf get nil"))
		return nil, err
	}
	appName := fmt.Sprintf("%s-%s", common.BaetylCore, common.RandString(9))
	// create config
	confMap := map[string]string{
		"NodeAddress": nodeAddress.Value,
		"NodeName":    nodeName,
		"AppName":     appName,
		"Namespace":   ns,
		"ConfigName":  fmt.Sprintf("%s-%s-config-%s", common.BaetylCore, nodeName, common.RandString(9)),
	}
	conf, err := api.genConfig(ns, common.TemplateJsonConfigCore, confMap, isSys)
	if err != nil {
		return nil, err
	}

	// create secret
	sync, err := api.genCertSync(appName, nodeName, ns, common.BaetylCore, isSys)
	if err != nil {
		return nil, err
	}

	// create application
	appMap := map[string]string{
		"AppName":         appName,
		"Image":           imageConf.Value,
		"CertSync":        sync.Name,
		"CertSyncVersion": sync.Version,
		"NodeName":        nodeName,
		"Namespace":       ns,
		"ConfigName":      conf.Name,
		"ConfigVersion":   conf.Version,
		"AppType":         common.ContainerApp,
	}
	return api.genApp(ns, common.TemplateJsonAppCore, appMap, isSys)
}

func (api *API) GenFunctionApp(nodeName, ns string, isSys bool) (*specV1.Application, error) {
	// get sys config
	imageConf, err := api.sysConfigService.GetSysConfig(common.BaetylModule, string(common.BaetylFunction))
	if err != nil {
		return nil, err
	}

	appName := fmt.Sprintf("%s-%s", common.BaetylFunction, common.RandString(9))
	// create config
	confMap := map[string]string{
		"AppName":    appName,
		"NodeName":   nodeName,
		"Namespace":  ns,
		"ConfigName": fmt.Sprintf("%s-%s-config-%s", common.BaetylFunction, nodeName, common.RandString(9)),
	}
	conf, err := api.genConfig(ns, common.TemplateJsonConfigFunction, confMap, isSys)
	if err != nil {
		return nil, err
	}

	// create application
	appMap := map[string]string{
		"AppName":       appName,
		"Image":         imageConf.Value,
		"NodeName":      nodeName,
		"Namespace":     ns,
		"ConfigName":    conf.Name,
		"ConfigVersion": conf.Version,
		"AppType":       common.ContainerApp,
	}
	return api.genApp(ns, common.TemplateJsonAppFunction, appMap, isSys)
}

func (api *API) genCertSync(appName, nodeName, ns string, module common.SystemApplication, isSys bool) (*specV1.Secret, error) {
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

func (api *API) genConfig(ns, template string, params map[string]string, isSys bool) (*specV1.Configuration, error) {
	confJson, err := api.ParseTemplate(template, params)
	if err != nil {
		return nil, err
	}
	config := &specV1.Configuration{}
	err = json.Unmarshal(confJson, config)
	if err != nil {
		return nil, err
	}
	config.System = isSys
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

func (api *API) genApp(ns, template string, params map[string]string, isSys bool) (*specV1.Application, error) {
	appJson, err := api.ParseTemplate(template, params)
	if err != nil {
		return nil, err
	}
	application := &specV1.Application{}
	err = json.Unmarshal(appJson, application)
	if err != nil {
		return nil, err
	}
	application.System = isSys
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
