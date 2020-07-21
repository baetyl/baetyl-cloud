package api

import (
	"github.com/baetyl/baetyl-cloud/common"
	"github.com/baetyl/baetyl-cloud/models"
	"time"
)

// TODO: optimize this layer, general abstraction

func (api *API) CreateCache(c *common.Context) (interface{}, error) {
	cache := &models.Cache{}
	err := c.LoadBody(cache)
	if err != nil {
		return nil, err
	}
	return nil, api.cacheService.Set(cache.Key,cache.Value)
}

func (api *API) DeleteCache(c *common.Context) (interface{}, error) {
	key := c.Param("key")
	return nil, api.cacheService.Delete(key)
}

func (api *API) GetCache(c *common.Context) (interface{}, error) {
	key := c.Param("key")
	return api.cacheService.Get(key)
}

func (api *API) ListCache(c *common.Context) (interface{}, error) {
	params := &models.Filter{}
	if err := c.Bind(params); err != nil {
		return nil, err
	}
	params.Format()
	return api.cacheService.List(params)
}

func (api *API) UpdateCache(c *common.Context) (interface{}, error) {
	key := c.Param("key")
	_, err := api.cacheService.Get(key)
	if err != nil{
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", err.Error()))
	}
	cache := &models.Cache{
		UpdateTime: time.Now(),
	}
	err = c.LoadBody(cache)
	if err != nil {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", err.Error()))
	}
	err = api.cacheService.Set(key, cache.Value)
	if err != nil {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", err.Error()))
	}
	return nil, nil
}
