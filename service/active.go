package service

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"time"

	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/models"
)

//go:generate mockgen -destination=../mock/service/active.go -package=service github.com/baetyl/baetyl-cloud/v2/service ActiveService

const (
	InfoKind      = "k"
	InfoName      = "n"
	InfoNamespace = "ns"
	InfoTimestamp = "ts"
	InfoExpiry    = "e"
)

var (
	ErrInvalidToken = fmt.Errorf("invalid token")
)
// ActiveService
type ActiveService interface {
	GetResource(resourceName, node, token string) (interface{}, error)
	GenCmd(kind, ns, name string) (string, error)
}

type ActiveServiceImpl struct {
	cfg           *config.CloudConfig
	AuthService   AuthService
	NodeService   NodeService
	SecretService SecretService
	CacheService  CacheService
	TemplateService TemplateService
}

// NewSyncService new SyncService
func NewActiveService(config *config.CloudConfig) (ActiveService, error) {
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
	templateService, err := NewTemplateService(config)
	if err != nil {
		return nil, err
	}
	return &ActiveServiceImpl{
		cfg:			config,
		AuthService:	authService,
		NodeService:    nodeService,
		TemplateService:templateService,
		SecretService:  secretService,
		CacheService:   cacheService,
	}, nil
}

func (a *ActiveServiceImpl) GetResource(resourceName, node, token string) (interface{}, error) {
	switch resourceName {
	case common.ResourceMetrics:
		res, err := a.TemplateService.GetTemplate(templateResourceMetrics)
		if err != nil {
			return nil, err
		}
		return []byte(res), nil
	case common.ResourceLocalPathStorage:
		res, err := a.TemplateService.GetTemplate(templateResourceLocalPathStorage)
		if err != nil {
			return nil, err
		}
		return []byte(res), nil
	case common.ResourceSetup:
		return a.TemplateService.GenSetupShell(token)

	case common.ResourceInitYaml:
		return a.getInitYaml(token, node)
	default:
		return nil, common.Error(
			common.ErrResourceNotFound,
			common.Field("type", "resource"),
			common.Field("name", resourceName))
	}
}

func (a *ActiveServiceImpl) getInitYaml(token, edgeKubeNodeName string) ([]byte, error) {
	info, err := a.CheckAndParseToken(token)
	if err != nil {
		return nil, common.Error(
			common.ErrRequestParamInvalid,
			common.Field("error", err))
	}
	switch common.Resource(info[InfoKind].(string)) {
	case common.Node:
		return a.initWithNode(info[InfoNamespace].(string), info[InfoName].(string), edgeKubeNodeName)
	default:
		return nil, common.Error(
			common.ErrRequestParamInvalid,
			common.Field("error", err))
	}
}

func (a *ActiveServiceImpl) CheckAndParseToken(token string) (map[string]interface{}, error) {
	// check len
	if len(token) < 10 {
		return nil, ErrInvalidToken
	}

	// check sign
	data, err := hex.DecodeString(token[10:])
	if err != nil {
		return nil, err
	}
	info := map[string]interface{}{}
	err = json.Unmarshal(data, &info)
	if err != nil {
		return nil, err
	}
	realToken, err := a.AuthService.GenToken(info)
	if err != nil {
		return nil, err
	}
	if realToken != token {
		return nil, ErrInvalidToken
	}

	expiry, ok := info[InfoExpiry].(float64)
	if !ok {
		return nil, ErrInvalidToken
	}

	ts, ok := info[InfoTimestamp].(float64)
	if !ok {
		return nil, ErrInvalidToken
	}
	// check expiration
	timestamp := time.Unix(int64(ts), 0)
	if timestamp.Add(time.Duration(int64(expiry))*time.Second).Unix() < time.Now().Unix() {
		return nil, ErrInvalidToken
	}
	return info, nil
}

func (a *ActiveServiceImpl) initWithNode(ns, nodeName, edgeKubeNodeName string) ([]byte, error) {
	params, err := a.getSysParams(ns, edgeKubeNodeName)
	if err != nil {
		return nil, err
	}
	params["NodeName"] = nodeName

	app, err := a.NodeService.GetDesireAppInfo(ns, nodeName)
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

func (a *ActiveServiceImpl) InitWithBatch(batch *models.Batch, edgeKubeNodeName string) ([]byte, error) {
	ns := batch.Namespace
	params, err := a.getSysParams(ns, edgeKubeNodeName)
	if err != nil {
		return nil, err
	}

	ca, err := ioutil.ReadFile(a.cfg.ActiveServer.Certificate.CA)
	if err != nil {
		return nil, err
	}
	params["CertActive"] = "baetyl-cloud.ca"
	params["CertActiveCa"] = base64.StdEncoding.EncodeToString(ca)

	params["BatchName"] = batch.Name
	params["SecurityType"] = string(batch.SecurityType)
	params["SecurityKey"] = batch.SecurityKey
	params["ProofType"] = common.FingerprintMap[batch.Fingerprint.Type]
	if (batch.Fingerprint.Type & common.FingerprintSN) == common.FingerprintSN {
		params["SnHostPath"] = path.Dir(batch.Fingerprint.SnPath)
		params["ProofValue"] = path.Base(batch.Fingerprint.SnPath)
	}
	if (batch.Fingerprint.Type & common.FingerprintInput) == common.FingerprintInput {
		params["ProofValue"] = path.Base(batch.Fingerprint.InputField)
		params["ContainerPort"] = common.DefaultActiveWebPort
		params["HostPort"] = common.DefaultActiveWebPort
	}

	return a.TemplateService.ParseTemplate(common.ResourceInitYaml, params)
}

func (a *ActiveServiceImpl) GetSyncCert(ns, nodeName string) (*specV1.Secret, error) {
	app, err := a.NodeService.GetDesireAppInfo(ns, nodeName)
	if err != nil {
		return nil, err
	}
	return a.getSyncCert(app)
}

func (a *ActiveServiceImpl) getSyncCert(app *specV1.Application) (*specV1.Secret, error) {
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

func (a *ActiveServiceImpl) GenCmd(kind, ns, name string) (string, error) {
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
	activeAddr, err := a.CacheService.GetProperty(propertyActiveServerAddress)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`curl -skfL '%s/v1/active/setup.sh?token=%s' -osetup.sh && sh setup.sh`, activeAddr, token), nil
}

func (a *ActiveServiceImpl) getSysParams(ns, nodeName string) (map[string]interface{}, error) {
	imageConf, err := a.CacheService.GetProperty(string(common.BaetylInit))
	if err != nil {
		return nil, common.Error(
			common.ErrResourceNotFound,
			common.Field("type", "system config baetyl module"),
			common.Field("name", common.BaetylInit),
			common.Field("namespace", ns))
	}
	nodeAddress, err := a.CacheService.GetProperty(propertySyncServerAddress)
	if err != nil {
		return nil, common.Error(
			common.ErrResourceNotFound,
			common.Field("type", "system config address"),
			common.Field("name", propertySyncServerAddress),
			common.Field("namespace", ns))
	}
	activeAddress, err := a.CacheService.GetProperty(propertyActiveServerAddress)
	if err != nil {
		return nil, common.Error(
			common.ErrResourceNotFound,
			common.Field("type", "system config address"),
			common.Field("name", propertyActiveServerAddress),
			common.Field("namespace", ns))
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
