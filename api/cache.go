package api

import (
	"github.com/baetyl/baetyl-cloud/common"
	"github.com/baetyl/baetyl-cloud/models"
	"strings"
	"time"
)

// TODO: optimize this layer, general abstraction

func (api *API) amisAccessAuth(c *common.Context) (interface{}, error) {
	misAccessToken := "test_token_123"
	//misAccessToken, err := api.cacheService.Get("bce.iot.hubble.showx.access.token")
	token := c.Request.Header.Get("amis_token")
	if strings.Compare(token, misAccessToken) != 0 {
		return nil, common.Error(common.ErrAMisTokenForbidden, common.Field("error", common.Code(common.ErrAMisTokenForbidden)))
	}
	user := c.Request.Header.Get("amis_user")
	if len(user) == 0 {
		return nil, common.Error(common.ErrAMisUserNotFound, common.Field("error", common.Code(common.ErrAMisUserNotFound)))
	}
	return nil, nil
}

func (api *API) CreateCache(c *common.Context) (interface{}, error) {
	_, err := api.amisAccessAuth(c)
	if err != nil {
		return nil, err
	}
	cache := &models.Cache{}
	err = c.LoadBody(cache)
	if err != nil {
		return nil, err
	}
	return nil, api.cacheService.Set(cache.Key, cache.Value)
}

func (api *API) DeleteCache(c *common.Context) (interface{}, error) {
	_, err := api.amisAccessAuth(c)
	if err != nil {
		return nil, err
	}
	key := c.Param("key")
	return nil, api.cacheService.Delete(key)
}

func (api *API) GetCache(c *common.Context) (interface{}, error) {
	_, err := api.amisAccessAuth(c)
	if err != nil {
		return nil, err
	}
	key := c.Param("key")
	return api.cacheService.Get(key)
}

func (api *API) ListCache(c *common.Context) (interface{}, error) {
	_, err := api.amisAccessAuth(c)
	if err != nil {
		return nil, err
	}
	params := &models.Filter{}
	if err := c.Bind(params); err != nil {
		return nil, err
	}
	params.Format()
	return api.cacheService.List(params)
}

func (api *API) UpdateCache(c *common.Context) (interface{}, error) {
	_, err := api.amisAccessAuth(c)
	if err != nil {
		return nil, err
	}
	cache := &models.Cache{
		UpdateTime: time.Now(),
	}
	err = c.LoadBody(cache)
	if err != nil {
		return nil, err
	}
	err = api.cacheService.Set(c.Param("key"), cache.Value)
	return nil, err
}
