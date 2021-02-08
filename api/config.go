package api

import (
	"encoding/json"
	"reflect"
	"strings"
	"time"

	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/baetyl/baetyl-go/v2/log"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/jinzhu/copier"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
)

const (
	ConfigTypeKV         = "kv"
	ConfigTypeObject     = "object"
	ConfigTypeFunction   = "function"
	ConfigObjectTypeHttp = "http"
)

// TODO: optimize this layer, general abstraction

// GetConfig get a config
func (api *API) GetConfig(c *common.Context) (interface{}, error) {
	ns, n := c.GetNamespace(), c.GetNameFromParam()
	config, err := api.Config.Get(ns, n, "")
	if err != nil {
		return nil, err
	}
	return api.ToConfigurationView(config)
}

// ListConfig list config
func (api *API) ListConfig(c *common.Context) (interface{}, error) {
	ns := c.GetNamespace()
	params, err := api.parseListOptionsAppendSystemLabel(c)
	if err != nil {
		return nil, err
	}
	list, err := api.Config.List(ns, params)
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
	config.Namespace = ns

	// TODO: remove get method, return error inside service instead
	oldConfig, err := api.Config.Get(ns, name, "")
	if err != nil {
		if e, ok := err.(errors.Coder); !ok || e.Code() != common.ErrResourceNotFound {
			return nil, err
		}
	}

	if oldConfig != nil {
		return nil, common.Error(common.ErrRequestParamInvalid,
			common.Field("error", "this name is already in use"))
	}

	config, err = api.Config.Create(ns, config)
	if err != nil {
		return nil, err
	}

	return api.ToConfigurationView(config)
}

// UpdateConfig update the config
func (api *API) UpdateConfig(c *common.Context) (interface{}, error) {
	config, err := api.parseAndCheckConfigView(c)
	if err != nil {
		return nil, err
	}
	ns, n := c.GetNamespace(), c.GetNameFromParam()

	res, err := api.Config.Get(ns, n, "")
	if err != nil {
		return nil, err
	}

	// labels can't be modified of sys apps
	if checkIsSysResources(res.Labels) && !reflect.DeepEqual(res.Labels, config.Labels) {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "labels can't be modified of sys apps"))
	}

	if models.EqualConfig(res, config) {
		return api.ToConfigurationView(res)
	}

	config.Version = res.Version
	config.UpdateTimestamp = time.Now()
	res, err = api.Config.Update(ns, config)
	if err != nil {
		log.L().Error("Update config failed", log.Error(err))
		return nil, err
	}

	appNames, err := api.Index.ListAppIndexByConfig(ns, res.Name)
	if err != nil {
		log.L().Error("list app index by config failed", log.Error(err))
		return nil, err
	}

	if err := api.updateNodeAndApp(ns, res, appNames); err != nil {
		log.L().Error("update node and app failed", log.Error(err))
		return nil, err
	}

	return api.ToConfigurationView(res)
}

// DeleteConfig delete the config
func (api *API) DeleteConfig(c *common.Context) (interface{}, error) {
	ns, n := c.GetNamespace(), c.GetNameFromParam()
	res, err := api.Config.Get(ns, n, "")
	if err != nil {
		if e, ok := err.(errors.Coder); ok && e.Code() == common.ErrResourceNotFound {
			return nil, nil
		}
		log.L().Error("get config failed", log.Error(err), log.Any("name", n), log.Any("namespace", ns))
		return nil, err
	}

	appNames, err := api.Index.ListAppIndexByConfig(ns, res.Name)
	if err != nil {
		log.L().Error("list app index by config failed", log.Error(err), log.Any("name", n), log.Any("namespace", ns))
		return nil, err
	}

	if len(appNames) > 0 {
		return nil, common.Error(common.ErrResourceHasBeenUsed,
			common.Field("type", "config"),
			common.Field("name", n))
	}

	//TODO: should remove file(bos/aws) of a function Config
	return nil, api.Config.Delete(c.GetNamespace(), c.GetNameFromParam())
}

func (api *API) GetAppByConfig(c *common.Context) (interface{}, error) {
	ns, n := c.GetNamespace(), c.GetNameFromParam()
	res, err := api.Config.Get(ns, n, "")
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
	configView.Namespace = c.GetNamespace()
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
				ok = checkElementsExist(item.Value, "source")
				if !ok {
					return nil, common.Error(common.ErrRequestParamInvalid,
						common.Field("error", "failed to validate object data of config"))
				}
				if item.Value["source"] == ConfigObjectTypeHttp {
					ok = checkElementsExist(item.Value, "url")
				}
				if !ok {
					return nil, common.Error(common.ErrRequestParamInvalid,
						common.Field("error", "failed to validate object data of config"))
				}
			case ConfigTypeFunction:
				ok = checkElementsExist(item.Value, "function", "version", "runtime",
					"handler", "bucket", "object")
				if !ok {
					return nil, common.Error(common.ErrRequestParamInvalid,
						common.Field("error", "failed to validate function data of config"))
				}
			case ConfigTypeKV:
				if strings.HasPrefix(item.Key, common.ConfigObjectPrefix) {
					return nil, common.Error(common.ErrRequestParamInvalid,
						common.Field("error", "key of kv data can't start with "+common.ConfigObjectPrefix))
				}
			}
		}
	}

	config, err := api.ToConfiguration(c.GetUser().ID, configView)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (api *API) updateNodeAndApp(namespace string, config *specV1.Configuration, appNames []string) error {
	for _, appName := range appNames {
		app, err := api.App.Get(namespace, appName, "")
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
		app, err = api.App.Update(namespace, app)
		if err != nil {
			return err
		}
		_, err = api.Node.UpdateNodeAppVersion(namespace, app)
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
			strings.Compare(config.Version, volume.Config.Version) > 0 {
			volume.Config.Version = config.Version
			appNeedUpdate = true
		}
	}
	return appNeedUpdate
}

func (api *API) ToConfigurationView(config *specV1.Configuration) (*models.ConfigurationView, error) {
	configView := new(models.ConfigurationView)
	err := copier.Copy(configView, config)
	if err != nil {
		return nil, err
	}

	for k, v := range config.Data {
		obj := models.ConfigDataItem{
			Key:   k,
			Value: map[string]string{},
		}

		var object specV1.ConfigurationObject
		if strings.HasPrefix(k, common.ConfigObjectPrefix) {
			obj.Key = strings.TrimPrefix(k, common.ConfigObjectPrefix)
			err := json.Unmarshal([]byte(v), &object)
			if err != nil {
				return nil, err
			}
		}

		if object.Metadata != nil {
			delete(object.Metadata, "userID")
			obj.Value = object.Metadata
		} else {
			obj.Value = map[string]string{
				"type":  ConfigTypeKV,
				"value": v,
			}
		}

		configView.Data = append(configView.Data, obj)
	}
	return configView, nil
}

func (api *API) ToConfiguration(userID string, configView *models.ConfigurationView) (*specV1.Configuration, error) {
	config := new(specV1.Configuration)
	err := copier.Copy(config, configView)
	if err != nil {
		return nil, err
	}

	config.Data = map[string]string{}
	for _, v := range configView.Data {
		switch v.Value["type"] {
		case ConfigTypeKV:
			config.Data[v.Key] = v.Value["value"]
		case ConfigTypeFunction, ConfigTypeObject:
			object := &specV1.ConfigurationObject{
				URL:      v.Value["url"],
				MD5:      v.Value["md5"],
				Unpack:   v.Value["unpack"],
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
