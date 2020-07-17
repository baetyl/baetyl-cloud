package api

import (
	"github.com/baetyl/baetyl-cloud/common"
	"github.com/baetyl/baetyl-cloud/models"
	"time"
)

// TODO: optimize this layer, general abstraction

// GetSystemConfig get a systemtem config
func (api *API) GetSystemConfig(c *common.Context) (interface{}, error) {
	key := c.Param("key")
	return api.cacheService.GetSystemConfig(key)
}

func (api *API) ListSystemConfig(c *common.Context) (interface{}, error) {
	params := &models.Filter{}
	if err := c.Bind(params); err != nil {
		return nil, err
	}
	params.Format()
	return api.cacheService.ListSystemConfig(params)
}

//// CreateSystemConfig create a systemtem config
func (api *API) CreateSystemConfig(c *common.Context) (interface{}, error) {
	systemConfig := &models.SystemConfig{
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	err := c.LoadBody(systemConfig)
	if err != nil {
		return nil, err
	}

	return api.cacheService.CreateSystemConfig(systemConfig)
}

func (api *API) DeleteSystemConfig(c *common.Context) (interface{}, error) {
	key := c.Param("key")
	return nil, api.cacheService.DeleteSystemConfig(key)
}

func (api *API) UpdateSystemConfig(c *common.Context) (interface{}, error) {
	key := c.Param("key")
	oldSystemConfig, err := api.cacheService.GetSystemConfig(key)
	// ensure that the modified data exists
	if err != nil {
		return nil, common.Error(common.ErrResourceHasBeenUsed,
			common.Field("error", "this systemtem config does not exist"))
	}
	oldSystemConfig.UpdateTime = time.Now()

	err = c.LoadBody(oldSystemConfig)
	if err != nil {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", err.Error()))
	}

	return api.cacheService.UpdateSystemConfig(oldSystemConfig)
}
