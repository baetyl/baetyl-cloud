package service

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
	"github.com/baetyl/baetyl-go/v2/log"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
)

//go:generate mockgen -destination=../mock/service/sync.go -package=service github.com/baetyl/baetyl-cloud/v2/service SyncService

// SyncService sync service
type SyncService interface {
	Report(namespace, name string, report specV1.Report) (specV1.Desire, error)
	Desire(namespace string, infos []specV1.ResourceInfo) ([]specV1.ResourceValue, error)
}

type syncService struct {
	plugin.ModelStorage
	plugin.DBStorage
	cs            ConfigService
	ns            NodeService
	as            ApplicationService
	secretService SecretService
	objectService ObjectService
}

// NewSyncService new SyncService
func NewSyncService(config *config.CloudConfig) (SyncService, error) {
	ms, err := plugin.GetPlugin(config.Plugin.ModelStorage)
	if err != nil {
		return nil, err
	}
	db, err := plugin.GetPlugin(config.Plugin.DatabaseStorage)
	if err != nil {
		return nil, err
	}
	es := &syncService{
		ModelStorage: ms.(plugin.ModelStorage),
		DBStorage:    db.(plugin.DBStorage),
	}
	es.cs, err = NewConfigService(config)
	if err != nil {
		return nil, err
	}
	es.ns, err = NewNodeService(config)
	if err != nil {
		return nil, err
	}
	es.as, err = NewApplicationService(config)
	if err != nil {
		return nil, err
	}
	es.secretService, err = NewSecretService(config)
	if err != nil {
		return nil, err
	}
	es.objectService, err = NewObjectService(config)
	if err != nil {
		return nil, err
	}
	return es, nil
}

func (t *syncService) Report(namespace, name string, report specV1.Report) (specV1.Desire, error) {
	shadow, err := t.ns.UpdateReport(namespace, name, report)
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

	delta, err := shadow.Desire.Diff(shadow.Report)
	if err != nil {
		log.L().Error("failed to calculate node delta",
			log.Any(common.KeyContextNamespace, namespace),
			log.Any("name", name),
			log.Error(err))
		return nil, err
	}

	return delta, nil
}

func (t *syncService) Desire(namespace string, crdInfos []specV1.ResourceInfo) ([]specV1.ResourceValue, error) {
	var crdDatas []specV1.ResourceValue
	for _, info := range crdInfos {
		crdData := specV1.ResourceValue{
			ResourceInfo: info,
		}
		log.L().Info("sync get crd", log.Any("kind", info.Kind), log.Any("name", info.Name))
		switch info.Kind {
		case specV1.KindApplication, specV1.KindApp:
			app, err := t.as.Get(namespace, info.Name, info.Version)
			if err != nil {
				log.L().Error("failed to get application", log.Any(common.KeyContextNamespace, namespace), log.Any("name", info.Name))
				return nil, err
			}
			crdData.Value.Value = app
		case specV1.KindConfiguration, specV1.KindConfig:
			config, err := t.cs.Get(namespace, info.Name, "")
			if err != nil {
				log.L().Error("failed to get config", log.Any(common.KeyContextNamespace, namespace), log.Any("name", info.Name))
				return nil, err
			}
			if err = t.populateConfig(config); err != nil {
				log.L().Error("failed to populate config", log.Any(common.KeyContextNamespace, namespace), log.Any("name", info.Name))
				return nil, err
			}
			crdData.Value.Value = config
		case specV1.KindSecret:
			secret, err := t.secretService.Get(namespace, info.Name, "")
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

func (t *syncService) populateConfig(cfg *specV1.Configuration) error {
	for k, v := range cfg.Data {
		if strings.HasPrefix(k, common.ConfigObjectPrefix) {
			obj := new(specV1.ConfigurationObject)
			err := json.Unmarshal([]byte(v), &obj)
			if err != nil {
				return err
			}
			if obj.URL != "" {
				continue
			}

			bytes, _ := json.Marshal(obj.Metadata)
			var configObject models.ConfigObjectItem
			err = json.Unmarshal(bytes, &configObject)
			res, err := t.objectService.GenObjectURL(obj.Metadata["userID"], configObject)
			if err != nil {
				return err
			}
			obj.URL = res.URL
			obj.Token = res.Token
			if res.MD5 != "" {
				obj.MD5 = res.MD5
			}
			obj.Unpack = configObject.Unpack
			obj.Metadata = nil

			data, err := json.Marshal(obj)
			if err != nil {
				return err
			}
			cfg.Data[k] = string(data)
		}
	}
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
