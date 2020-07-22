package api

import (
	"github.com/baetyl/baetyl-cloud/common"
	"github.com/baetyl/baetyl-cloud/models"
	"strings"
)

// TODO: optimize this layer, general abstraction

func (api *API) amisAccessAuth(c *common.Context) (interface{}, error) {
	misAccessToken := "test_token_123"
	//misAccessToken, err := api.propertyService.Get("bce.iot.hubble.showx.access.token")
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

func (api *API) CreateProperty(c *common.Context) (interface{}, error) {
	_, err := api.amisAccessAuth(c)
	if err != nil {
		return nil, err
	}
	property := &models.Property{}
	err = c.LoadBody(property)
	if err != nil {
		return nil, err
	}
	return api.propertyService.CreateProperty(property)
}

func (api *API) DeleteProperty(c *common.Context) (interface{}, error) {
	_, err := api.amisAccessAuth(c)
	if err != nil {
		return nil, err
	}
	return nil, api.propertyService.DeleteProperty(c.Param("key"))
}

func (api *API) GetProperty(c *common.Context) (interface{}, error) {
	_, err := api.amisAccessAuth(c)
	if err != nil {
		return nil, err
	}
	return api.propertyService.GetProperty(c.Param("key"))
}

func (api *API) ListProperty(c *common.Context) (interface{}, error) {
	_, err := api.amisAccessAuth(c)
	if err != nil {
		return nil, err
	}
	params := &models.Filter{}
	if err := c.Bind(params); err != nil {
		return nil, err
	}
	params.Format()
	return api.propertyService.ListProperty(params)
}

func (api *API) UpdateProperty(c *common.Context) (interface{}, error) {
	_, err := api.amisAccessAuth(c)
	if err != nil {
		return nil, err
	}
	property := &models.Property{
		Key: c.Param("key"),
	}
	err = c.LoadBody(property)
	if err != nil {
		return nil, err
	}
	return api.propertyService.UpdateProperty(property)
}
