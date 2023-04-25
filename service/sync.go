package service

import (
	"fmt"
	"strings"

	"github.com/baetyl/baetyl-go/v2/json"
	"github.com/baetyl/baetyl-go/v2/log"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/models"
)

//go:generate mockgen -destination=../mock/service/sync.go -package=service github.com/baetyl/baetyl-cloud/v2/service SyncService

// SyncService sync service
type SyncService interface {
	Report(namespace, name string, report specV1.Report) (specV1.Delta, error)
	Desire(namespace string, infos []specV1.ResourceInfo, metadata map[string]string) ([]specV1.ResourceValue, error)
}

type HandlerPopulateConfig func(cfg *specV1.Configuration, metadata map[string]string) error

const (
	HookNamePopulateConfig = "populateConfig"
)

type SyncServiceImpl struct {
	ConfigService ConfigService
	NodeService   NodeService
	AppService    ApplicationService
	SecretService SecretService
	ObjectService ObjectService
	Hooks         map[string]interface{}
}

// NewSyncService new SyncService
func NewSyncService(config *config.CloudConfig) (SyncService, error) {
	es := &SyncServiceImpl{
		Hooks: map[string]interface{}{},
	}
	var err error
	es.ConfigService, err = NewConfigService(config)
	if err != nil {
		return nil, err
	}
	es.NodeService, err = NewNodeService(config)
	if err != nil {
		return nil, err
	}
	es.AppService, err = NewApplicationService(config)
	if err != nil {
		return nil, err
	}
	es.SecretService, err = NewSecretService(config)
	if err != nil {
		return nil, err
	}
	es.ObjectService, err = NewObjectService(config)
	if err != nil {
		return nil, err
	}
	es.Hooks[HookNamePopulateConfig] = HandlerPopulateConfig(es.PopulateConfig)
	return es, nil
}

func (t *SyncServiceImpl) Report(namespace, name string, report specV1.Report) (specV1.Delta, error) {
	shadow, err := t.NodeService.UpdateReport(namespace, name, report)
	if err != nil {
		log.L().Error("failed to update node reported status",
			log.Any(common.KeyContextNamespace, namespace),
			log.Any("name", name),
			log.Error(err))
		return nil, err
	}

	err = checkSysapp(name, &shadow.Desire)

	if err != nil {
		log.L().Error("system app wasnot ready",
			log.Any(common.KeyContextNamespace, namespace),
			log.Any("name", name),
			log.Error(err))
		return nil, err
	}

	node, err := t.NodeService.Get(nil, namespace, name)
	if err != nil {
		log.L().Error("failed to get node",
			log.Any(common.KeyContextNamespace, namespace),
			log.Any("name", name),
			log.Error(err))
		return nil, err
	}

	syncMode := specV1.CloudMode
	if node.Attributes != nil {
		syncMode, _ = node.Attributes[specV1.KeySyncMode].(specV1.SyncMode)
	}

	var delta specV1.Delta
	if syncMode != specV1.LocalMode {
		delta, err = shadow.Desire.DiffWithNil(extractComparingReport(shadow.Report))
		if err != nil {
			log.L().Error("failed to calculate node delta",
				log.Any(common.KeyContextNamespace, namespace),
				log.Any("name", name),
				log.Error(err))
			return nil, err
		}
	}
	// TODO remove in the future
	if delta != nil && shadow.Desire[common.NodeProps] != nil {
		delta[common.NodeProps] = shadow.Desire[common.NodeProps]
	}

	return delta, nil
}

func extractComparingReport(report specV1.Report) specV1.Report {
	res := map[string]interface{}{}
	if apps, ok := report["apps"]; ok {
		res["apps"] = apps
	}
	if sysapps, ok := report["sysapps"]; ok {
		res["sysapps"] = sysapps
	}
	if nodeProps, ok := report[common.NodeProps]; ok {
		res[common.NodeProps] = nodeProps
	}
	return res
}

func (t *SyncServiceImpl) Desire(namespace string, crdInfos []specV1.ResourceInfo, metadata map[string]string) ([]specV1.ResourceValue, error) {
	var crdDatas []specV1.ResourceValue
	for _, info := range crdInfos {
		crdData := specV1.ResourceValue{
			ResourceInfo: info,
		}
		log.L().Info("sync get crd", log.Any("kind", info.Kind), log.Any("name", info.Name))
		switch info.Kind {
		case specV1.KindApplication, specV1.KindApp:
			app, err := t.AppService.Get(namespace, info.Name, info.Version)
			if err != nil {
				log.L().Error("failed to get application", log.Any(common.KeyContextNamespace, namespace), log.Any("name", info.Name))
				return nil, err
			}
			crdData.Value.Value = app
		case specV1.KindConfiguration, specV1.KindConfig:
			cfg, err := t.ConfigService.Get(namespace, info.Name, info.Version)
			if err != nil {
				log.L().Error("failed to get config", log.Any(common.KeyContextNamespace, namespace), log.Any("name", info.Name))
				return nil, err
			}
			if err = t.Hooks[HookNamePopulateConfig].(HandlerPopulateConfig)(cfg, metadata); err != nil {
				log.L().Error("failed to populate config", log.Any(common.KeyContextNamespace, namespace), log.Any("name", info.Name))
				return nil, err
			}
			crdData.Value.Value = cfg
		case specV1.KindSecret:
			secret, err := t.SecretService.Get(namespace, info.Name, info.Version)
			if err != nil {
				log.L().Error("failed to get secret", log.Any(common.KeyContextNamespace, namespace), log.Any("name", info.Name))
				return nil, err
			}
			crdData.Value.Value = secret
		default:
			return nil, fmt.Errorf("unsupported request type")
		}
		crdDatas = append(crdDatas, crdData)
	}
	return crdDatas, nil
}

func (t *SyncServiceImpl) PopulateConfig(cfg *specV1.Configuration, metadata map[string]string) error {
	for k, v := range cfg.Data {
		if strings.HasPrefix(k, common.ConfigObjectPrefix) {
			err := t.PopulateConfigObject(k, v, cfg)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *SyncServiceImpl) PopulateConfigObject(k, v string, cfg *specV1.Configuration) error {
	obj := new(specV1.ConfigurationObject)
	err := json.Unmarshal([]byte(v), obj)
	if err != nil {
		return err
	}
	if obj.URL != "" {
		obj.Metadata = nil
		return nil
	}

	bytes, err := json.Marshal(obj.Metadata)
	if err != nil {
		return err
	}
	var item models.ConfigObjectItem
	err = json.Unmarshal(bytes, &item)
	if err != nil {
		return err
	}

	var res *models.ObjectURL
	if item.Endpoint == "" {
		res, err = t.ObjectService.GenInternalObjectURL(obj.Metadata["userID"], item.Bucket, item.Object, item.Source)
	} else {
		res, err = t.ObjectService.GenExternalObjectURL(models.ExternalObjectInfo{
			Endpoint:      item.Endpoint,
			Ak:            item.Ak,
			Sk:            item.Sk,
			AddressFormat: item.AddressFormat,
		}, item.Bucket, item.Object, item.Source)
	}
	if err != nil {
		return err
	}
	obj.URL = res.URL
	obj.Token = res.Token
	if res.MD5 != "" {
		obj.MD5 = res.MD5
	}
	obj.Unpack = item.Unpack
	obj.Metadata = nil

	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	cfg.Data[k] = string(data)
	return nil
}

func checkSysapp(name string, desire *specV1.Desire) error {
	if desire == nil {
		return common.Error(common.ErrNodeNotReady, common.Field("name", name))
	}

	if len(desire.AppInfos(true)) == 0 {
		return common.Error(common.ErrNodeNotReady, common.Field("name", name))
	}

	return nil
}
