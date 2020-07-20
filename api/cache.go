package api

import (
	"github.com/baetyl/baetyl-cloud/common"
	"github.com/baetyl/baetyl-cloud/models"
)

// TODO: optimize this layer, general abstraction

// GetCache get a systemtem config
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

//// AddCache create a systemtem config
func (api *API) AddCache(c *common.Context) (interface{}, error) {
	var cache models.Cache
	err := c.LoadBody(cache)
	if err != nil {
		return nil, err
	}

	return api.cacheService.Set(cache.Key,cache.Value)
}

func (api *API) DeleteCache(c *common.Context) (interface{}, error) {
	key := c.Param("key")
	return nil, api.cacheService.Delete(key)
}

func (api *API) ReplaceCache(c *common.Context) (interface{}, error) {
	key := c.Param("key")
	var cache models.Cache
	err := c.LoadBody(cache)
	if err != nil {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", err.Error()))
	}

	return api.cacheService.Set(key,cache.Value)
}
