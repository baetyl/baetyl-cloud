package service

import (
	"encoding/base64"
	"strings"
	"time"

	"github.com/baetyl/baetyl-go/v2/context"
	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/baetyl/baetyl-go/v2/log"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
)

//go:generate mockgen -destination=../mock/service/init.go -package=service github.com/baetyl/baetyl-cloud/v2/service InitService

const (
	InfoName      = "n"
	InfoNamespace = "ns"
	InfoExpiry    = "e"
)

const (
	templateInitDeploymentYaml = "baetyl-init-deployment.yml"
	templateBaetylInstallShell = "baetyl-install.sh"

	TemplateCoreConfYaml      = "baetyl-core-conf.yml"
	TemplateInitConfYaml      = "baetyl-init-conf.yml"
	TemplateBaetylInitCommand = "baetyl-init-command"
	TemplateInitCommandWget   = "baetyl-init-command-wget"
)

var (
	CmdExpirationInSeconds = int64(60 * 60)
)

type GetInitResource func(ns, nodeName string, params map[string]interface{}) ([]byte, error)

// InitService
type InitService interface {
	GetResource(ns, nodeName, resourceName string, params map[string]interface{}) (interface{}, error)
}

type InitServiceImpl struct {
	cfg             *config.CloudConfig
	SignService     SignService
	NodeService     NodeService
	Property        PropertyService
	TemplateService TemplateService
	*AppCombinedService
	PKI             PKIService
	ResourceMapFunc map[string]GetInitResource
	log             *log.Logger
}

// NewSyncService new SyncService
func NewInitService(config *config.CloudConfig) (InitService, error) {
	signService, err := NewSignService(config)
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
		SignService:        signService,
		NodeService:        nodeService,
		Property:           propertyService,
		TemplateService:    templateService,
		AppCombinedService: acs,
		PKI:                pki,
		ResourceMapFunc:    map[string]GetInitResource{},
		log:                log.L().With(log.Any("service", "init")),
	}
	initService.ResourceMapFunc[templateInitDeploymentYaml] = initService.getInitDeploymentYaml
	initService.ResourceMapFunc[TemplateBaetylInitCommand] = initService.GetInitCommand
	initService.ResourceMapFunc[TemplateCoreConfYaml] = initService.getCoreConfig
	initService.ResourceMapFunc[templateBaetylInstallShell] = initService.getInstallShell
	initService.ResourceMapFunc[TemplateInitConfYaml] = initService.getInitConfig

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

	node, err := s.NodeService.Get(nil, ns, nodeName)
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
	initCommand, err := s.Property.GetPropertyValue(params["template"].(string))
	if err != nil {
		return nil, err
	}
	token, err := s.SignService.GenToken(info)
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

func (s *InitServiceImpl) getCoreConfig(ns, nodeName string, params map[string]interface{}) ([]byte, error) {
	params["Namespace"] = ns
	params["NodeName"] = nodeName
	return s.TemplateService.ParseTemplate(TemplateCoreConfYaml, params)
}

func (s *InitServiceImpl) getInitConfig(ns, nodeName string, params map[string]interface{}) ([]byte, error) {
	params["Namespace"] = ns
	params["NodeName"] = nodeName
	return s.TemplateService.ParseTemplate(TemplateInitConfYaml, params)
}
