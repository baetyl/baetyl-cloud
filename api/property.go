package api

import (
	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
)

func (api *API) CreateProperty(c *common.Context) (interface{}, error) {
	property := &models.Property{}
	err := c.LoadBody(property)
	if err != nil {
		return nil, err
	}
	return nil, api.propertyService.CreateProperty(property)
}

func (api *API) DeleteProperty(c *common.Context) (interface{}, error) {
	return nil, api.propertyService.DeleteProperty(c.Param("name"))
}

func (api *API) ListProperty(c *common.Context) (interface{}, error) {
	params := &models.Filter{}
	if err := c.Bind(params); err != nil {
		return nil, err
	}
	params.Format()
	properties, err := api.propertyService.ListProperty(params)
	if err != nil {
		return nil, err
	}
	count, err := api.propertyService.CountProperty(params.Name)
	if err != nil {
		return nil, err
	}
	return models.MisData{
		Count: count,
		Rows:  properties,
	}, nil
}

func (api *API) UpdateProperty(c *common.Context) (interface{}, error) {
	property := &models.Property{
		Name: c.Param("name"),
	}
	err := c.LoadBody(property)
	if err != nil {
		return nil, err
	}
	return nil, api.propertyService.UpdateProperty(property)
}
