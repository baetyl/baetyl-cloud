package api

import (
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
	// remove deprecated baetyl-gpu-metrics and dmp
	for i := 0; i < len(res); i++ {
		if res[i].Name == DeprecatedGPUMetrics || res[i].Name == DeprecatedDmp {
			res = append(res[:i], res[i+1:]...)
			i--
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
	return nil
}
