package api

import (
	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
)

// CreateNamespace create one namespace
func (api *API) CreateNamespace(c *common.Context) (interface{}, error) {
	return api.namespaceService.Create(&models.Namespace{
		Name: c.GetNamespace(),
	})
}

// GetNamespace get one namespace
func (api *API) GetNamespace(c *common.Context) (interface{}, error) {
	res, err := api.namespaceService.Get(c.GetNamespace())
	if res == nil {
		return nil, common.Error(common.ErrResourceNotFound,
			common.Field("type", "namespace"),
			common.Field("name", c.GetNamespace()))
	}
	return res, err
}

func (api *API) DeleteNamespace(c *common.Context) (interface{}, error) {
	ns := c.GetNamespace()
	return nil, api.namespaceService.Delete(&models.Namespace{Name: ns})

}
