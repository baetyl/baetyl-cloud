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

var (
	CmdExpirationInSeconds = int64(60 * 60)
)

// InitService
type InitService interface {
	GetResource(resourceName, node, token string, info map[string]interface{}) (interface{}, error)
	GenCmd(kind, ns, name string) (string, error)
}

type InitServiceImpl struct {
	cfg             *config.CloudConfig
	AuthService     AuthService
	NodeService     NodeService
	SecretService   SecretService
	CacheService    CacheService
	AppService      ApplicationService
	TemplateService TemplateService
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
	appService, err := NewApplicationService(config)
	if err != nil {
		return nil, err
	}
	templateService, err := NewTemplateService(config)
	if err != nil {
		return nil, err
	}
	return &InitServiceImpl{
		cfg:             config,
		AuthService:     authService,
		NodeService:     nodeService,
		TemplateService: templateService,
		SecretService:   secretService,
		CacheService:    cacheService,
		AppService:      appService,
	}, nil
}

func (a *InitServiceImpl) GetResource(resourceName, node, token string, info map[string]interface{}) (interface{}, error) {
	switch resourceName {
	case common.ResourceMetrics:
		res, err := a.TemplateService.GetTemplate(ResourceMetrics)
		if err != nil {
			return nil, err
		}
		return []byte(res), nil
	case common.ResourceLocalPathStorage:
		res, err := a.TemplateService.GetTemplate(ResourceLocalPathStorage)
		if err != nil {
			return nil, err
		}
		return []byte(res), nil
	case common.ResourceSetup:
		return a.TemplateService.GenSetupShell(token)

	case common.ResourceInitYaml:
		return a.getInitYaml(info, node)
	default:
		return nil, common.Error(
			common.ErrResourceNotFound,
			common.Field("type", "resource"),
			common.Field("name", resourceName))
	}
}

func (a *InitServiceImpl) getInitYaml(info map[string]interface{}, edgeKubeNodeName string) ([]byte, error) {
	switch common.Resource(info[InfoKind].(string)) {
	case common.Node:
		return a.genInitYml(info[InfoNamespace].(string), info[InfoName].(string), edgeKubeNodeName)
	default:
		return nil, common.Error(
			common.ErrRequestParamInvalid,
			common.Field("error", "invalid info kind"))
	}
}

func (a *InitServiceImpl) genInitYml(ns, nodeName, edgeKubeNodeName string) ([]byte, error) {
	params, err := a.getSysParams(ns, edgeKubeNodeName)
	if err != nil {
		return nil, err
	}
	params["NodeName"] = nodeName

	app, err := a.getDesireAppInfo(ns, nodeName)
	if err != nil {
		return nil, err
	}
	params["CoreVersion"] = app.Version

	sync, err := a.getSyncCert(app)
	if err != nil {
		return nil, err
	}

	params["CertSync"] = sync.Name
	params["CertSyncPem"] = base64.StdEncoding.EncodeToString(sync.Data["client.pem"])
	params["CertSyncKey"] = base64.StdEncoding.EncodeToString(sync.Data["client.key"])
	params["CertSyncCa"] = base64.StdEncoding.EncodeToString(sync.Data["ca.pem"])

	return a.TemplateService.ParseTemplate(common.ResourceInitYaml, params)
}

func (a *InitServiceImpl) getSyncCert(app *specV1.Application) (*specV1.Secret, error) {
	certName := ""
	for _, vol := range app.Volumes {
		if vol.Name == "node-certs" {
			certName = vol.Secret.Name
			break
		}
	}
	cert, _ := a.SecretService.Get(app.Namespace, certName, "")
	if cert == nil {
		return nil, common.Error(
			common.ErrResourceNotFound,
			common.Field("type", "secret"),
			common.Field("name", certName),
			common.Field("namespace", app.Namespace))
	}
	return cert, nil
}

func (a *InitServiceImpl) GenCmd(kind, ns, name string) (string, error) {
	info := map[string]interface{}{
		InfoKind:      kind,
		InfoName:      name,
		InfoNamespace: ns,
		InfoExpiry:    CmdExpirationInSeconds,
		InfoTimestamp: time.Now().Unix(),
	}
	token, err := a.AuthService.GenToken(info)
	if err != nil {
		return "", err
	}
	activeAddr, err := a.CacheService.GetProperty(propertyInitServerAddress)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`curl -skfL '%s/v1/active/setup.sh?token=%s' -osetup.sh && sh setup.sh`, activeAddr, token), nil
}

func (a *InitServiceImpl) getSysParams(ns, nodeName string) (map[string]interface{}, error) {
	imageConf, err := a.CacheService.GetProperty("baetyl-image")
	if err != nil {
		return nil, errors.Trace(err)
	}
	nodeAddress, err := a.CacheService.GetProperty(propertySyncServerAddress)
	if err != nil {
		return nil, errors.Trace(err)
	}
	activeAddress, err := a.CacheService.GetProperty(propertyInitServerAddress)
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

func (a *InitServiceImpl) getDesireAppInfo(ns, nodeName string) (*specV1.Application, error) {
	shadowDesire, err := a.NodeService.GetDesire(ns, nodeName)
	if err != nil {
		return nil, err
	}
	apps := shadowDesire.AppInfos(true)
	for _, appInfo := range apps {
		if strings.Contains(appInfo.Name, string(common.BaetylCore)) {
			app, _ := a.AppService.Get(ns, appInfo.Name, "")
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
