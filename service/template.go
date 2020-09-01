package service

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path"
	"text/template"

	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"gopkg.in/yaml.v2"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"

	"github.com/baetyl/baetyl-go/v2/errors"

	"github.com/baetyl/baetyl-cloud/v2/config"
)

const (
	templateInitAppName = "baetyl-init"
	templateCoreAppName = "baetyl-core"
	templateFuncAppName = "baetyl-function"

	templateCoreConfYaml = "baetyl-core-conf.yml"
	templateCoreAppYaml  = "baetyl-core-app.yml"
	templateFuncConfYaml = "baetyl-function-conf.yml"
	templateFuncAppYaml  = "baetyl-function-app.yml"
	templateSetupShell   = "setup.sh"

	propertySyncServerAddress   = "sync-server-address"
	propertyActiveServerAddress = "active-server-address"
)

//go:generate mockgen -destination=../mock/service/template.go -package=service github.com/baetyl/baetyl-cloud/v2/service TemplateService

type TemplateService interface {
	GetTemplate(filename string) (string, error)
	ParseTemplate(filename string, params map[string]interface{}) ([]byte, error)
	UnmarshalTemplate(filename string, params map[string]interface{}, out interface{}) error

	// the following functions are business logic
	GenSetupShell(token string) ([]byte, error)
	GenSystemApps(ns, nodeName string, params map[string]string) ([]*specV1.Application, error)
}

// TemplateServiceImpl is a service to read and parse template files.
type TemplateServiceImpl struct {
	path  string
	cache CacheService
	// TODO: move the following services out of template, template service only generates models without creating
	pki   PKIService
	*AppCombinedService
}

func NewTemplateService(cfg *config.CloudConfig) (TemplateService, error) {
	cacheService, err := NewCacheService(cfg)
	if err != nil {
		return nil, errors.Trace(err)
	}
	pkiService, err := NewPKIService(cfg)
	if err != nil {
		return nil, errors.Trace(err)
	}
	rs, err := NewAppCombinedService(cfg)
	if err != nil {
		return nil, err
	}
	return &TemplateServiceImpl{
		path:               cfg.Template.Path,
		cache:              cacheService,
		pki:                pkiService,
		AppCombinedService: rs,
	}, nil
}

func (s *TemplateServiceImpl) GetTemplate(filename string) (string, error) {
	return s.cache.Get(filename, func(key string) (string, error) {
		data, err := ioutil.ReadFile(path.Join(s.path, key))
		if err != nil {
			return "", common.Error(common.ErrTemplate, common.Field("error", err))
		}
		return string(data), nil
	})
}

func (s *TemplateServiceImpl) ParseTemplate(filename string, params map[string]interface{}) ([]byte, error) {
	tl, err := s.GetTemplate(filename)
	if err != nil {
		return nil, errors.Trace(err)
	}
	t, err := template.New(filename).Parse(tl)
	if err != nil {
		return nil, common.Error(common.ErrTemplate, common.Field("error", err))
	}
	buf := &bytes.Buffer{}
	err = t.Execute(buf, params)
	if err != nil {
		return nil, common.Error(common.ErrTemplate, common.Field("error", err))
	}
	return buf.Bytes(), nil
}

func (s *TemplateServiceImpl) UnmarshalTemplate(filename string, params map[string]interface{}, out interface{}) error {
	tp, err := s.ParseTemplate(filename, params)
	if err != nil {
		return errors.Trace(err)
	}
	err = yaml.Unmarshal(tp, out)
	if err != nil {
		return common.Error(common.ErrTemplate, common.Field("error", err))
	}
	return nil
}

// business logic

func (s *TemplateServiceImpl) GenSetupShell(token string) ([]byte, error) {
	activeAddr, err := s.cache.GetProperty(propertyActiveServerAddress)
	if err != nil {
		return nil, errors.Trace(err)
	}
	params := map[string]interface{}{
		"Token":     token,
		"CloudAddr": activeAddr,
	}
	data, err := s.ParseTemplate(templateSetupShell, params)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return data, nil
}

func (s *TemplateServiceImpl) GenSystemApps(ns, nodeName string, params map[string]string) ([]*specV1.Application, error) {
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

func (s *TemplateServiceImpl) genCoreApp(ns, nodeName string, params map[string]string) (*specV1.Application, error) {
	syncAddr, err := s.cache.GetProperty(propertySyncServerAddress)
	if err != nil {
		return nil, errors.Trace(err)
	}
	appName := fmt.Sprintf("%s-%s", templateCoreAppName, common.RandString(9))
	confName := fmt.Sprintf("%s-conf-%s", templateCoreAppName, common.RandString(9))
	// create config
	confMap := map[string]interface{}{
		"Namespace":    ns,
		"NodeName":     nodeName,
		"SyncAddr":     syncAddr,
		"CoreAppName":  appName,
		"CoreConfName": confName,
	}
	for k, v := range params {
		confMap[k] = v
	}
	conf, err := s.genConfig(ns, templateCoreConfYaml, confMap)
	if err != nil {
		return nil, errors.Trace(err)
	}

	// create secret
	cert, err := s.genNodeCerts(ns, nodeName, appName)
	if err != nil {
		return nil, errors.Trace(err)
	}

	// create application
	appMap := map[string]interface{}{
		"Namespace":       ns,
		"NodeName":        nodeName,
		"CoreAppName":     appName,
		"CoreCertName":    cert.Name,
		"CoreCertVersion": cert.Version,
		"CoreConfName":    conf.Name,
		"CoreConfVersion": conf.Version,
	}
	for k, v := range params {
		appMap[k] = v
	}
	return s.genApp(ns, templateCoreAppYaml, appMap)
}

func (s *TemplateServiceImpl) genFunctionApp(ns, nodeName string, params map[string]string) (*specV1.Application, error) {
	appName := fmt.Sprintf("%s-%s", templateFuncAppName, common.RandString(9))
	confName := fmt.Sprintf("%s-conf-%s", templateFuncAppName, common.RandString(9))
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

func (s *TemplateServiceImpl) genNodeCerts(ns, nodeName, appName string) (*specV1.Secret, error) {
	confName := fmt.Sprintf("crt-%s-%s", nodeName, common.RandString(9))
	certName := fmt.Sprintf(`%s.%s`, ns, nodeName)
	certPEM, err := s.pki.SignClientCertificate(certName, models.AltNames{})
	if err != nil {
		return nil, errors.Trace(err)
	}

	ca, err := s.pki.GetCA()
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

func (s *TemplateServiceImpl) genConfig(ns, template string, params map[string]interface{}) (*specV1.Configuration, error) {
	config := &specV1.Configuration{}
	err := s.UnmarshalTemplate(template, params, config)
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

func (s *TemplateServiceImpl) genApp(ns, template string, params map[string]interface{}) (*specV1.Application, error) {
	application := &specV1.Application{}
	err := s.UnmarshalTemplate(template, params, application)
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
