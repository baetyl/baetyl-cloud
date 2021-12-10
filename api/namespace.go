package api

import (
	"github.com/baetyl/baetyl-go/v2/log"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
)

// CreateNamespace create one namespace
func (api *API) CreateNamespace(c *common.Context) (interface{}, error) {
	ns, err := api.NS.Create(&models.Namespace{
		Name: c.GetNamespace(),
	})
	if err != nil {
		return ns, err
	}
	if e := api.InitQuotas(ns.Name); e != nil {
		log.L().Error("InitQuotas error", log.Error(e))
	}
	return ns, nil
}

// GetNamespace get one namespace
func (api *API) GetNamespace(c *common.Context) (interface{}, error) {
	res, err := api.NS.Get(c.GetNamespace())
	if res == nil {
		return nil, common.Error(common.ErrResourceNotFound,
			common.Field("type", "namespace"),
			common.Field("name", c.GetNamespace()))
	}
	return res, err
}

func (api *API) DeleteNamespace(c *common.Context) (interface{}, error) {
	ns := c.GetNamespace()
	_, err := api.NS.Get(ns)
	if err != nil {
		return nil, err
	}
	_, err = api.Task.AddTaskWithKey("DeleteNamespaceTask", map[string]interface{}{"ns": ns})

	return nil, err
}
