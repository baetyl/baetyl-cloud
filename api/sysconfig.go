package api

import (
	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
)

// TODO: to use property service with cache

func (api *API) GetSysConfig(c *common.Context) (interface{}, error) {
	tp, key := c.Param("type"), c.Param("key")
	return api.sysConfigService.GetSysConfig(tp, key)
}

func (api *API) ListSysConfig(c *common.Context) (interface{}, error) {
	res, err := api.sysConfigService.ListSysConfigAll(c.Param("type"))
	if err != nil {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", err.Error()))
	}
	return &models.SysConfigView{
		SysConfigs: res,
	}, nil
}
