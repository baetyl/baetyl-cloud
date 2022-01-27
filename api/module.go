package api

import (
	"fmt"

	"github.com/baetyl/baetyl-go/v2/log"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
)

func (api *API) GetModules(c *common.Context) (interface{}, error) {
	res, err := api.Module.GetModules(c.Param("name"))
	if err != nil {
		return nil, err
	}
	return models.ListView{
		Total: len(res),
		Items: res,
	}, nil
}

func (api *API) GetModuleByVersion(c *common.Context) (interface{}, error) {
	res, err := api.Module.GetModuleByVersion(c.Param("name"), c.Param("version"))
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (api *API) GetLatestModule(c *common.Context) (interface{}, error) {
	res, err := api.Module.GetLatestModule(c.Param("name"))
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (api *API) CreateModule(c *common.Context) (interface{}, error) {
	var module models.Module
	err := api.parseAndCheckModule(&module, c)
	if err != nil {
		return nil, err
	}
	res, err := api.Module.CreateModule(&module)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (api *API) UpdateModule(c *common.Context) (interface{}, error) {
	name, version := c.GetNameFromParam(), c.Param("version")
	module, err := api.Module.GetModuleByVersion(name, version)
	if err != nil {
		return nil, err
	}

	err = api.parseAndCheckModule(module, c)
	if err != nil {
		return nil, err
	}
	module.Name = name
	module.Version = version
	res, err := api.Module.UpdateModuleByVersion(module)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (api *API) DeleteModules(c *common.Context) (interface{}, error) {
	name, version := c.GetNameFromParam(), c.Param("version")
	var err error
	if version == "" {
		err = api.Module.DeleteModules(name)
	} else {
		err = api.Module.DeleteModuleByVersion(name, version)
	}
	if err != nil {
		log.L().Error("failed to delete modules", log.Any("module", c.GetNameFromParam()), log.Any("version", version), log.Error(err))
	}
	return nil, nil
}

func (api *API) ListModules(c *common.Context) (interface{}, error) {
	tp := c.Query("type")
	params := &models.Filter{}
	if err := c.Bind(params); err != nil {
		return nil, err
	}
	res, err := api.Module.ListModules(params, common.ModuleType(tp))
	if err != nil {
		return nil, err
	}
	// remove deprecated baetyl-gpu-metrics
	if tp == string(common.TypeSystemKube) {
		index := -1
		for i, app := range res {
			if app.Name == DeprecatedGPUMetrics {
				index = i
				break
			}
		}
		if index != -1 {
			res = append(res[:index], res[index+1:]...)
		}
	}

	return models.ListView{
		Total:    len(res),
		PageNo:   params.PageNo,
		PageSize: params.PageSize,
		Items:    res,
	}, nil
}

func (api *API) parseAndCheckModule(module *models.Module, c *common.Context) error {
	err := c.LoadBody(module)
	if err != nil {
		return common.Error(common.ErrRequestParamInvalid, common.Field("error", err.Error()))
	}
	if module.Name == "" {
		return common.Error(common.ErrRequestParamInvalid, common.Field("error", "name is required"))
	}
	if module.Version == "" {
		return common.Error(common.ErrRequestParamInvalid, common.Field("error", "version is required"))
	}
	if module.Type == string(common.TypeSystemOptional) {
		supportSysApps := api.SysApp.GetOptionalApps()
		var ok bool
		for _, v := range supportSysApps {
			if v == module.Name {
				ok = true
			}
		}
		if !ok {
			return common.Error(common.ErrRequestParamInvalid, common.Field("error", fmt.Sprintf("the module (%s) isn't optional system module", module.Name)))
		}
	}
	return nil
}
