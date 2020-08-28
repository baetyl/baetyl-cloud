package service

import (
	"bytes"
	"encoding/base64"
	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"io/ioutil"
	"path"
	"strings"
	"text/template"
)

//go:generate mockgen -destination=../mock/service/initialize.go -package=service github.com/baetyl/baetyl-cloud/v2/service InitializeService

type InitializeService interface {
	InitWithNode(ns, nodeName, edgeKubeNodeName string) ([]byte, error)
	InitWithBatch(batch *models.Batch, edgeKubeNodeName string) ([]byte, error)
	GetResource(key string) (string, error)
	GetSyncCert(ns, nodeName string) (*specV1.Secret, error)
}

type initializeService struct {
	cfg          *config.CloudConfig
	modelStorage plugin.ModelStorage
	dbStorage    plugin.DBStorage
	shadow       plugin.Shadow
}

// NewInitializeService New Initialize Service
func NewInitializeService(config *config.CloudConfig) (InitializeService, error) {
	ms, err := plugin.GetPlugin(config.Plugin.ModelStorage)
	if err != nil {
		return nil, err
	}
	ds, err := plugin.GetPlugin(config.Plugin.DatabaseStorage)
	if err != nil {
		return nil, err
	}
	ss, err := plugin.GetPlugin(config.Plugin.Shadow)
	if err != nil {
		return nil, err
	}
	return &initializeService{
		cfg:          config,
		modelStorage: ms.(plugin.ModelStorage),
		dbStorage:    ds.(plugin.DBStorage),
		shadow:       ss.(plugin.Shadow),
	}, nil
}

func (init *initializeService) GetResource(key string) (string, error) {
	if common.Cache == nil {
		common.Cache = map[string]string{}
	}
	if tl, ok := common.Cache[key]; ok {
		return tl, nil
	}
	sysConf, _ := init.dbStorage.GetSysConfig("resource", key)
	if sysConf == nil {
		return "", common.Error(
			common.ErrResourceNotFound,
			common.Field("type", "resource"),
			common.Field("name", key))
	}
	data, err := base64.StdEncoding.DecodeString(sysConf.Value)
	if err != nil {
		return "", err
	}
	tl := string(data)
	common.Cache[key] = tl
	return tl, nil
}

func (init *initializeService) InitWithNode(ns, nodeName, edgeKubeNodeName string) ([]byte, error) {
	params, err := init.getSysParams(ns, edgeKubeNodeName)
	if err != nil {
		return nil, err
	}
	params["NodeName"] = nodeName

	app, err := init.getCoreApp(ns, nodeName)
	if err != nil {
		return nil, err
	}
	params["CoreVersion"] = app.Version

	sync, err := init.getSyncCert(app)
	if err != nil {
		return nil, err
	}

	params["CertSync"] = sync.Name
	params["CertSyncPem"] = base64.StdEncoding.EncodeToString(sync.Data["client.pem"])
	params["CertSyncKey"] = base64.StdEncoding.EncodeToString(sync.Data["client.key"])
	params["CertSyncCa"] = base64.StdEncoding.EncodeToString(sync.Data["ca.pem"])

	return init.parseInitApp(common.ResourceInitYaml, params)
}

func (init *initializeService) InitWithBatch(batch *models.Batch, edgeKubeNodeName string) ([]byte, error) {
	ns := batch.Namespace
	params, err := init.getSysParams(ns, edgeKubeNodeName)
	if err != nil {
		return nil, err
	}

	ca, err := ioutil.ReadFile(init.cfg.ActiveServer.Certificate.CA)
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

	return init.parseInitApp(common.ResourceInitYaml, params)
}

func (init *initializeService) GetSyncCert(ns, nodeName string) (*specV1.Secret, error) {
	app, err := init.getCoreApp(ns, nodeName)
	if err != nil {
		return nil, err
	}
	return init.getSyncCert(app)
}

func (init *initializeService) getCoreApp(ns, nodeName string) (*specV1.Application, error) {
	shadow, _ := init.shadow.Get(ns, nodeName)
	if shadow == nil {
		return nil, common.Error(
			common.ErrResourceNotFound,
			common.Field("type", "node"),
			common.Field("name", nodeName),
			common.Field("namespace", ns))
	}
	apps := shadow.Desire.AppInfos(true)
	for _, appInfo := range apps {
		if strings.Contains(appInfo.Name, string(common.BaetylCore)) {
			app, _ := init.modelStorage.GetApplication(ns, appInfo.Name, "")
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

func (init *initializeService) getSyncCert(app *specV1.Application) (*specV1.Secret, error) {
	certName := ""
	for _, vol := range app.Volumes {
		if vol.Name == "cert-sync" {
			certName = vol.Secret.Name
			break
		}
	}
	cert, _ := init.modelStorage.GetSecret(app.Namespace, certName, "")
	if cert == nil {
		return nil, common.Error(
			common.ErrResourceNotFound,
			common.Field("type", "secret"),
			common.Field("name", certName),
			common.Field("namespace", app.Namespace))
	}
	return cert, nil
}

func (init *initializeService) parseInitApp(templateName string, params map[string]interface{}) ([]byte, error) {
	yaml, err := init.GetResource(templateName)
	if err != nil {
		return nil, err
	}
	tl, err := template.New(templateName).Parse(yaml)
	if err != nil {
		return nil, common.Error(common.ErrTemplate, common.Field("error", err))
	}
	buf := &bytes.Buffer{}
	err = tl.Execute(buf, params)
	if err != nil {
		return nil, common.Error(common.ErrTemplate, common.Field("error", err))
	}
	return buf.Bytes(), nil
}

func (init *initializeService) getSysParams(ns, edgeKubeNodeName string) (map[string]interface{}, error) {
	imageConf, _ := init.dbStorage.GetSysConfig(common.BaetylModule, string(common.BaetylInit))
	if imageConf == nil {
		return nil, common.Error(
			common.ErrResourceNotFound,
			common.Field("type", "system config baetyl module"),
			common.Field("name", common.BaetylInit),
			common.Field("namespace", ns))
	}
	nodeAddress, _ := init.dbStorage.GetSysConfig("address", common.AddressNode)
	if nodeAddress == nil {
		return nil, common.Error(
			common.ErrResourceNotFound,
			common.Field("type", "system config address"),
			common.Field("name", common.AddressNode),
			common.Field("namespace", ns))
	}
	activeAddress, _ := init.dbStorage.GetSysConfig("address", common.AddressActive)
	if activeAddress == nil {
		return nil, common.Error(
			common.ErrResourceNotFound,
			common.Field("type", "system config address"),
			common.Field("name", common.AddressActive),
			common.Field("namespace", ns))
	}
	return map[string]interface{}{
		"Namespace":           ns,
		"EdgeNamespace":       common.DefaultBaetylEdgeNamespace,
		"EdgeSystemNamespace": common.DefaultBaetylEdgeSystemNamespace,
		"NodeAddress":         nodeAddress.Value,
		"ActiveAddress":       activeAddress.Value,
		"Image":               imageConf.Value,
		"KubeNodeName":        edgeKubeNodeName,
	}, nil
}
