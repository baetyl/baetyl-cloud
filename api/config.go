package api

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/baetyl/baetyl-cloud/common"
	"github.com/baetyl/baetyl-cloud/models"
	"github.com/baetyl/baetyl-go/errors"
	"github.com/baetyl/baetyl-go/log"
	specV1 "github.com/baetyl/baetyl-go/spec/v1"
	"github.com/jinzhu/copier"
)

const (
	ConfigTypeKV       = "kv"
	ConfigTypeObject   = "object"
	ConfigTypeFunction = "function"
)

// TODO: optimize this layer, general abstraction

// GetConfig get a config
func (api *API) GetConfig(c *common.Context) (interface{}, error) {
	ns, n := c.GetNamespace(), c.GetNameFromParam()
	config, err := api.configService.Get(ns, n, "")
	if err != nil {
		return nil, err
	}
	return api.toConfigurationView(config)
}

// ListConfig list config
func (api *API) ListConfig(c *common.Context) (interface{}, error) {
	ns := c.GetNamespace()
	list, err := api.configService.List(ns, api.parseListOptionsAppendSystemLabel(c))
	if err != nil {
		log.L().Error("list config error", log.Error(err))
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", err.Error()))
	}
	for i := range list.Items {
		list.Items[i].Data = nil
	}
	return list, err
}

// CreateConfig create one config
func (api *API) CreateConfig(c *common.Context) (interface{}, error) {
	config, err := api.parseAndCheckConfigView(c)
	if err != nil {
		log.L().Error("parse and check config model failed", log.Error(err))
		return nil, err
	}

	ns, name := c.GetNamespace(), config.Name
	// TODO: remove get method, return error inside service instead
	oldConfig, err := api.configService.Get(ns, name, "")
	if err != nil {
		if e, ok := err.(errors.Coder); !ok || e.Code() != common.ErrResourceNotFound {
			return nil, err
		}
	}

	if oldConfig != nil {
		return nil, common.Error(common.ErrRequestParamInvalid,
			common.Field("error", "this name is already in use"))
	}

	config, err = api.configService.Create(ns, config)
	if err != nil {
		return nil, err
	}

	return api.toConfigurationView(config)
}

// UpdateConfig update the config
func (api *API) UpdateConfig(c *common.Context) (interface{}, error) {
	config, err := api.parseAndCheckConfigView(c)
	if err != nil {
		return nil, err
	}
	ns, n := c.GetNamespace(), c.GetNameFromParam()
	res, err := api.configService.Get(ns, n, "")

	if err != nil {
		return nil, err
	}

	if models.EqualConfig(res, config) {
		return api.toConfigurationView(res)
	}

	config.Version = res.Version
	config.UpdateTimestamp = time.Now()
	res, err = api.configService.Update(ns, config)
	if err != nil {
		log.L().Error("Update config failed", log.Error(err))
		return nil, err
	}

	appNames, err := api.indexService.ListAppIndexByConfig(ns, res.Name)
	if err != nil {
		log.L().Error("list app index by config failed", log.Error(err))
		return nil, err
	}

	if err := api.updateNodeAndApp(ns, res, appNames); err != nil {
		log.L().Error("update node and app failed", log.Error(err))
		return nil, err
	}

	return api.toConfigurationView(res)
}

// DeleteConfig delete the config
func (api *API) DeleteConfig(c *common.Context) (interface{}, error) {
	ns, n := c.GetNamespace(), c.GetNameFromParam()
	res, err := api.configService.Get(ns, n, "")
	if err != nil {
		log.L().Error("get config failed", log.Error(err))
		return nil, err
	}

	appNames, err := api.indexService.ListAppIndexByConfig(ns, res.Name)
	if err != nil {
		log.L().Error("list app index by config failed", log.Error(err))
		return nil, err
	}

	if len(appNames) > 0 {
		return nil, common.Error(common.ErrResourceHasBeenUsed,
			common.Field("type", "config"),
			common.Field("name", n))
	}

	//TODO: should remove file(bos/aws) of a function Config
	return nil, api.configService.Delete(c.GetNamespace(), c.GetNameFromParam())
}

func (api *API) GetAppByConfig(c *common.Context) (interface{}, error) {
	ns, n := c.GetNamespace(), c.GetNameFromParam()
	res, err := api.configService.Get(ns, n, "")
	if err != nil {
		log.L().Error("get config failed", log.Error(err))
		return nil, err
	}
	return api.listAppByConfig(ns, res.Name)
}

// parseAndCheckConfigModel parse and check the config model
func (api *API) parseAndCheckConfigView(c *common.Context) (*specV1.Configuration, error) {
	configView := new(models.ConfigurationView)
	configView.Name = c.GetNameFromParam()
	err := c.LoadBody(configView)
	if err != nil {
		log.L().Error("parse config failed", log.Error(err))
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", err.Error()))
	}
	if name := c.GetNameFromParam(); name != "" {
		configView.Name = name
	}
	if configView.Name == "" {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "name is required"))
	}

	for _, item := range configView.Data {
		if _type, ok := item.Value["type"]; ok {
			switch _type {
			case ConfigTypeObject:
				res := checkElementsExist(item.Value, "source", "bucket", "object")
				if !res {
					return nil, common.Error(common.ErrRequestParamInvalid,
						common.Field("error", "failed to validate object data of config"))
				}
			case ConfigTypeFunction:
				res := checkElementsExist(item.Value, "function", "version", "runtime",
					"handler", "bucket", "object")
				if !res {
					return nil, common.Error(common.ErrRequestParamInvalid,
						common.Field("error", "failed to validate function data of config"))
				}
			case ConfigTypeKV:
				if strings.HasPrefix(item.Key, common.ObjectSource) {
					return nil, common.Error(common.ErrRequestParamInvalid,
						common.Field("error", "key of kv data can't start with "+common.ConfigObjectPrefix))
				}
			}
		}
	}

	config, err := api.toConfiguration(c.GetUser().ID, configView)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (api *API) updateNodeAndApp(namespace string, config *specV1.Configuration, appNames []string) error {
	for _, appName := range appNames {
		app, err := api.applicationService.Get(namespace, appName, "")
		if err != nil {
			if e, ok := err.(errors.Coder); ok && e.Code() == common.ErrResourceNotFound {
				continue
			}
			return err
		}

		if !needUpdateApp(config, app) {
			continue
		}
		// Todo remove by list watch
		app, err = api.applicationService.Update(namespace, app)
		if err != nil {
			return err
		}
		_, err = api.nodeService.UpdateNodeAppVersion(namespace, app)
		if err != nil {
			return err
		}
	}
	return nil
}

func needUpdateApp(config *specV1.Configuration, app *specV1.Application) bool {
	appNeedUpdate := false
	for _, volume := range app.Volumes {
		if volume.Config != nil &&
			volume.Config.Name == config.Name &&
			// config's version must increment
			common.CompareNumericalString(config.Version, volume.Config.Version) > 0 {
			volume.Config.Version = config.Version
			appNeedUpdate = true
		}
	}
	return appNeedUpdate
}

func (api *API) toConfigurationView(config *specV1.Configuration) (*models.ConfigurationView, error) {
	configView := new(models.ConfigurationView)
	copier.Copy(configView, config)

	for k, v := range config.Data {
		obj := models.ConfigDataItem{
			Value: map[string]string{},
		}
		if strings.HasPrefix(k, common.ConfigObjectPrefix) {
			obj.Key = strings.TrimPrefix(k, common.ConfigObjectPrefix)
			var object specV1.ConfigurationObject
			err := json.Unmarshal([]byte(v), &object)
			if err != nil {
				return nil, err
			}
			obj.Value = object.Metadata
		} else {
			obj.Key = k
			obj.Value = map[string]string{
				"type":  ConfigTypeKV,
				"value": v,
			}
		}
		configView.Data = append(configView.Data, obj)
	}
	return configView, nil
}

func (api *API) toConfiguration(userID string, configView *models.ConfigurationView) (*specV1.Configuration, error) {
	config := new(specV1.Configuration)
	copier.Copy(config, configView)

	config.Data = map[string]string{}
	for _, v := range configView.Data {
		switch v.Value["type"] {
		case ConfigTypeKV:
			config.Data[v.Key] = v.Value["value"]
		case ConfigTypeFunction, ConfigTypeObject:
			object := &specV1.ConfigurationObject{
				MD5:      v.Value["md5"],
				Metadata: map[string]string{},
			}
			object.Metadata = v.Value
			object.Metadata["userID"] = userID
			bytes, err := json.Marshal(object)
			if err != nil {
				return nil, err
			}
			config.Data[common.ConfigObjectPrefix+v.Key] = string(bytes)
		}
	}
	return config, nil
}

func checkElementsExist(m map[string]string, elems ...string) bool {
	for _, v := range elems {
		if _, ok := m[v]; !ok {
			return false
		}
	}
	return true
}
