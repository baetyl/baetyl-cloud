package api

import (
	"github.com/baetyl/baetyl-cloud/common"
	"github.com/baetyl/baetyl-cloud/models"
)

// TODO: optimize this layer, general abstraction

func (api *API) CreateProperty(c *common.Context) (interface{}, error) {
	property := &models.Property{}
	err := c.LoadBody(property)
	if err != nil {
		return nil, err
	}
	return api.propertyService.CreateProperty(property)
}

func (api *API) DeleteProperty(c *common.Context) (interface{}, error) {
	return nil, api.propertyService.DeleteProperty(c.Param("key"))
}

func (api *API) GetProperty(c *common.Context) (interface{}, error) {
	return api.propertyService.GetProperty(c.Param("key"))
}

func (api *API) ListProperty(c *common.Context) (interface{}, error) {
	params := &models.Filter{}
	if err := c.Bind(params); err != nil {
		return nil, err
	}
	params.Format()
	return api.propertyService.ListProperty(params)
}

func (api *API) UpdateProperty(c *common.Context) (interface{}, error) {
	property := &models.Property{
		Key: c.Param("key"),
	}
	err := c.LoadBody(property)
	if err != nil {
		return nil, err
	}
	return api.propertyService.UpdateProperty(property)
}
