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
	return nil, api.cacheService.Set(cache.Key, cache.Value)
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
	res, err := api.cacheService.List(params)
	if err != nil {
		return nil, err
	}
	return models.AmisListView{
		Status: "0",
		Msg:    "ok",
		Data: models.AmisData{
			Count: res.Total,
			Rows:  res.Items,
		},
	}, nil
}

func (api *API) UpdateCache(c *common.Context) (interface{}, error) {
	cache := &models.Cache{
		UpdateTime: time.Now(),
	}
	err := c.LoadBody(cache)
	if err != nil {
		return nil, err
	}
	err = api.cacheService.Set(c.Param("key"), cache.Value)
	return nil, err
}
