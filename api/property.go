package api

import (
	"github.com/baetyl/baetyl-cloud/common"
	"github.com/baetyl/baetyl-cloud/models"
)

func (api *API) CreateProperty(c *common.Context) (interface{}, error) {
	property := &models.Property{}
	err := c.LoadBody(property)
	if err == nil {
		err = api.propertyService.CreateProperty(property)
	}
	return api.packageMisResponse(err)
}

func (api *API) DeleteProperty(c *common.Context) (interface{}, error) {
	err := api.propertyService.DeleteProperty(c.Param("key"))
	return api.packageMisResponse(err)
}

func (api *API) ListProperty(c *common.Context) (interface{}, error) {
	params := &models.Filter{}
	if err := c.Bind(params); err != nil {
		return nil, err
	}
	params.Format()
	properties, err := api.propertyService.ListProperty(params)
	if err != nil {
		return api.packageMisResponse(err)
	}
	count, err := api.propertyService.CountProperty(params.Name)
	if err != nil {
		return api.packageMisResponse(err)
	}
	return models.MisResponse{
		Status: "0",
		Msg:    "ok",
		Data: models.MisData{
			Count: count,
			Rows:  properties,
		},
	}, nil
}

func (api *API) UpdateProperty(c *common.Context) (interface{}, error) {
	property := &models.Property{
		Key: c.Param("key"),
	}
	err := c.LoadBody(property)
	if err == nil {
		err = api.propertyService.UpdateProperty(property)
	}
	return api.packageMisResponse(err)
}

func (api *API) packageMisResponse(err error) (interface{}, error) {
	var res models.MisResponse
	if err != nil {
		res = models.MisResponseFailure
		res.Msg = err.Error()
	}else{
		res = models.MisResponseSuccess
	}
	return res, nil
}
