package api

import (
	"github.com/baetyl/baetyl-cloud/common"
	"github.com/baetyl/baetyl-cloud/models"
)

func (api *API) CreateProperty(c *common.Context) error {
	property := &models.Property{}
	err := c.LoadBody(property)
	if err != nil {
		return err
	}
	return api.propertyService.CreateProperty(property)
}

func (api *API) DeleteProperty(c *common.Context) error {
	return api.propertyService.DeleteProperty(c.Param("key"))
}

func (api *API) ListProperty(c *common.Context) error {
	params := &models.Filter{}
	if err := c.Bind(params); err != nil {
		return err
	}
	params.Format()
	properties, err := api.propertyService.ListProperty(params)
	if err != nil {
		return err
	}
	count, err := api.propertyService.CountProperty(params.Name)
	if err != nil {
		return err
	}
	res := models.MisData{
		Count: count,
		Rows:  properties,
	}
	c.PureJSON(common.PackageMisResponse(res))
	return nil
}

func (api *API) UpdateProperty(c *common.Context) error {
	property := &models.Property{
		Key: c.Param("key"),
	}
	err := c.LoadBody(property)
	if err != nil {
		return err
	}
	return api.propertyService.UpdateProperty(property)
}
