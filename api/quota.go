package api

import (
	"github.com/baetyl/baetyl-cloud/v2/common"
)

func (api *API) GetQuota(c *common.Context) (interface{}, error) {
	ns := c.GetNamespace()
	quotas, err := api.License.GetQuota(ns)
	return quotas, err
}
