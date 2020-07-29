package api

import (
	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
)

func (api *API) CreateCallback(c *common.Context) (interface{}, error) {
	ns := c.GetNamespace()
	callback := &models.Callback{
		Namespace: ns,
		Params:    map[string]string{},
		Header:    map[string]string{},
		Body:      map[string]string{},
	}
	err := api.parseCallback(callback, c)
	if err != nil {
		return nil, err
	}
	return api.callbackService.Create(callback)
}

func (api *API) UpdateCallback(c *common.Context) (interface{}, error) {
	ns := c.GetNamespace()
	callbackName := c.Param("callbackName")
	callback, err := api.callbackService.Get(callbackName, ns)
	if err != nil {
		return nil, err
	}
	if callback == nil {
		return nil, common.Error(common.ErrResourceNotFound, common.Field("type", "callback"), common.Field("name", callbackName))
	}
	err = api.parseAndUpdateCallback(callback, c)
	if err != nil {
		return nil, err
	}
	return api.callbackService.Update(callback)
}

func (api *API) DeleteCallback(c *common.Context) (interface{}, error) {
	ns := c.GetNamespace()
	callbackName := c.Param("callbackName")
	return nil, api.callbackService.Delete(callbackName, ns)
}

func (api *API) GetCallback(c *common.Context) (interface{}, error) {
	ns := c.GetNamespace()
	callbackName := c.Param("callbackName")
	res, err := api.callbackService.Get(callbackName, ns)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, common.Error(common.ErrResourceNotFound, common.Field("type", "callback"), common.Field("name", callbackName))
	}
	return res, nil
}

func (api *API) parseCallback(callback *models.Callback, c *common.Context) error {
	err := c.LoadBody(callback)
	if err != nil {
		return common.Error(common.ErrRequestParamInvalid, common.Field("error", err.Error()))
	}

	return err
}

func (api *API) parseAndUpdateCallback(callback *models.Callback, c *common.Context) error {
	// todo optimization
	header := callback.Header
	callback.Header = nil
	body := callback.Body
	callback.Body = nil
	params := callback.Params
	callback.Params = nil

	err := api.parseCallback(callback, c)
	if err != nil {
		return err
	}
	if callback.Header == nil {
		callback.Header = header
	}
	if callback.Body == nil {
		callback.Body = body
	}
	if callback.Params == nil {
		callback.Params = params
	}
	return nil
}
