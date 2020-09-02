package service

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
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
	InfoKind                 = "k"
	InfoName                 = "n"
	InfoNamespace            = "ns"
	InfoTimestamp            = "ts"
	InfoExpiry               = "e"
	ResourceMetrics          = "metrics.yml"
	ResourceLocalPathStorage = "local-path-storage.yml"
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

	propertySyncServerAddress = "sync-server-address"
	propertyInitServerAddress = "init-server-address"
)

var (
	CmdExpirationInSeconds = int64(60 * 60)
)

// InitService
type InitService interface {
	GetResource(resourceName, node, token string, info map[string]interface{}) (interface{}, error)
	GenApps(ns, nodeName string, params map[string]string) ([]*specV1.Application, error)
	GenCmd(kind, ns, name string) (string, error)
}

type InitServiceImpl struct {
	cfg             *config.CloudConfig
	AuthService     AuthService
	NodeService     NodeService
	SecretService   SecretService
	CacheService    CacheService
	TemplateService TemplateService
	*AppCombinedService
	PKI PKIService
}

// NewSyncService new SyncService
func NewInitService(config *config.CloudConfig) (InitService, error) {
	authService, err := NewAuthService(config)
	if err != nil {
		return nil, err
	}
	nodeService, err := NewNodeService(config)
	if err != nil {
		return nil, err
	}
	secretService, err := NewSecretService(config)
	if err != nil {
		return nil, err
	}
	cacheService, err := NewCacheService(config)
	if err != nil {
		return nil, err
	}
	templateService, err := NewTemplateService(config, nil)
	if err != nil {
		return nil, err
	}
	acs, err := NewAppCombinedService(config)
	if err != nil {
		return nil, err
	}
	pki, err := NewPKIService(config)
	if err != nil {
		return nil, err
	}
	return &InitServiceImpl{
		cfg:                config,
		AuthService:        authService,
		NodeService:        nodeService,
		TemplateService:    templateService,
		SecretService:      secretService,
		CacheService:       cacheService,
		AppCombinedService: acs,
		PKI:                pki,
	}, nil
}

func (s *InitServiceImpl) GetResource(resourceName, node, token string, info map[string]interface{}) (interface{}, error) {
	switch resourceName {
	case common.ResourceMetrics:
		res, err := s.TemplateService.GetTemplate(ResourceMetrics)
		if err != nil {
			return nil, err
		}
		return []byte(res), nil
	case common.ResourceLocalPathStorage:
		res, err := s.TemplateService.GetTemplate(ResourceLocalPathStorage)
		if err != nil {
			return nil, err
		}
		return []byte(res), nil
	case common.ResourceSetup:
		return s.genSetupShell(token)
	case common.ResourceInitYaml:
		return s.getInitYaml(info, node)
	default:
		return nil, common.Error(
			common.ErrResourceNotFound,
			common.Field("type", "resource"),
			common.Field("name", resourceName))
	}
}

func (s *InitServiceImpl) getInitYaml(info map[string]interface{}, edgeKubeNodeName string) ([]byte, error) {
	switch common.Resource(info[InfoKind].(string)) {
	case common.Node:
		return s.genInitYml(info[InfoNamespace].(string), info[InfoName].(string), edgeKubeNodeName)
	default:
		return nil, common.Error(
			common.ErrRequestParamInvalid,
			common.Field("error", "invalid info kind"))
	}
}

func (s *InitServiceImpl) genInitYml(ns, nodeName, edgeKubeNodeName string) ([]byte, error) {
	params, err := s.getSysParams(ns, edgeKubeNodeName)
	if err != nil {
		return nil, err
	}
	params["NodeName"] = nodeName

	app, err := s.getDesireAppInfo(ns, nodeName)
	if err != nil {
		return nil, err
	}
	params["CoreVersion"] = app.Version

	sync, err := s.getSyncCert(app)
	if err != nil {
		return nil, err
	}

	params["CertSync"] = sync.Name
	params["CertSyncPem"] = base64.StdEncoding.EncodeToString(sync.Data["client.pem"])
	params["CertSyncKey"] = base64.StdEncoding.EncodeToString(sync.Data["client.key"])
	params["CertSyncCa"] = base64.StdEncoding.EncodeToString(sync.Data["ca.pem"])

	return s.TemplateService.ParseTemplate(common.ResourceInitYaml, params)
}

func (s *InitServiceImpl) getSyncCert(app *specV1.Application) (*specV1.Secret, error) {
	certName := ""
	for _, vol := range app.Volumes {
		if vol.Name == "node-certs" {
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

func (s *InitServiceImpl) GenCmd(kind, ns, name string) (string, error) {
	info := map[string]interface{}{
		InfoKind:      kind,
		InfoName:      name,
		InfoNamespace: ns,
		InfoExpiry:    CmdExpirationInSeconds,
		InfoTimestamp: time.Now().Unix(),
	}
	token, err := s.AuthService.GenToken(info)
	if err != nil {
		return "", err
	}
	activeAddr, err := s.CacheService.GetProperty(propertyInitServerAddress)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`curl -skfL '%s/v1/active/setup.sh?token=%s' -osetup.sh && sh setup.sh`, activeAddr, token), nil
}

func (s *InitServiceImpl) getSysParams(ns, nodeName string) (map[string]interface{}, error) {
	imageConf, err := s.CacheService.GetProperty("baetyl-image")
	if err != nil {
		return nil, errors.Trace(err)
	}
	nodeAddress, err := s.CacheService.GetProperty(propertySyncServerAddress)
	if err != nil {
		return nil, errors.Trace(err)
	}
	activeAddress, err := s.CacheService.GetProperty(propertyInitServerAddress)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return map[string]interface{}{
		"Namespace":           ns,
		"EdgeNamespace":       common.DefaultBaetylEdgeNamespace,
		"EdgeSystemNamespace": common.DefaultBaetylEdgeSystemNamespace,
		"NodeAddress":         nodeAddress,
		"ActiveAddress":       activeAddress,
		"Image":               imageConf,
		"KubeNodeName":        nodeName,
	}, nil
}

func (s *InitServiceImpl) getDesireAppInfo(ns, nodeName string) (*specV1.Application, error) {
	shadowDesire, err := s.NodeService.GetDesire(ns, nodeName)
	if err != nil {
		return nil, err
	}
	apps := shadowDesire.AppInfos(true)
	for _, appInfo := range apps {
		if strings.Contains(appInfo.Name, string(common.BaetylCore)) {
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

func (s *InitServiceImpl) genSetupShell(token string) ([]byte, error) {
	activeAddr, err := s.CacheService.GetProperty(propertyInitServerAddress)
	if err != nil {
		return nil, errors.Trace(err)
	}
	params := map[string]interface{}{
		"Token":     token,
		"CloudAddr": activeAddr,
	}
	data, err := s.TemplateService.ParseTemplate(templateSetupShell, params)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return data, nil
}

func (s *InitServiceImpl) GenApps(ns, nodeName string, params map[string]string) ([]*specV1.Application, error) {
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

func (s *InitServiceImpl) genCoreApp(ns, nodeName string, params map[string]string) (*specV1.Application, error) {
	syncAddr, err := s.CacheService.GetProperty(propertySyncServerAddress)
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

func (s *InitServiceImpl) genFunctionApp(ns, nodeName string, params map[string]string) (*specV1.Application, error) {
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
