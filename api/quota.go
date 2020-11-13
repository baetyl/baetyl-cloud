package api

import (
	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
	"github.com/baetyl/baetyl-go/v2/log"
)

func (api *API) GetQuota(c *common.Context) (interface{}, error) {
	ns := c.GetNamespace()
	quotas, err := api.License.GetQuota(ns)
	return quotas, err
}

func (api *API) CreateQuota(c *common.Context) (interface{}, error) {
	quota := &models.Quota{}
	if err := c.LoadBody(quota); err != nil {
		return nil, err
	}
	quotas := map[string]int{
		quota.QuotaName: quota.Quota,
	}
	err := api.License.CreateQuota(quota.Namespace, quotas)
	return nil, err
}

func (api *API) UpdateQuota(c *common.Context) (interface{}, error) {
	quota := &models.Quota{}
	if err := c.LoadBody(quota); err != nil {
		return nil, err
	}
	err := api.License.UpdateQuota(quota.Namespace, quota.QuotaName, quota.Quota)
	return nil, err
}

func (api *API) DeleteQuota(c *common.Context) (interface{}, error) {
	quota := &models.Quota{}
	if err := c.LoadBody(quota); err != nil {
		return nil, err
	}
	err := api.License.DeleteQuota(quota.Namespace, quota.QuotaName)
	return nil, err
}

func (api *API) InitQuotas(namespace string) error {
	quotas, err := api.License.GetDefaultQuotas(namespace)
	if err != nil {
		return err
	}
	if err := api.License.CreateQuota(namespace, quotas); err != nil {
		common.LogDirtyData(err,
			log.Any("type", "InitQuotas"),
			log.Any(common.KeyContextNamespace, namespace),
			log.Any("quotas", quotas))
		return err
	}

	return nil
}

func (api *API) DeleteAllQuotas(namespace string) error {
	if err := api.License.DeleteQuotaByNamespace(namespace); err != nil {
		common.LogDirtyData(err,
			log.Any("type", "DeleteAllQuotas"),
			log.Any(common.KeyContextNamespace, namespace))
		return err
	}

	return nil
}

func (api *API) realseNodeQuota(namespace string, number int) error {
	if err := api.License.ReleaseQuota(namespace, plugin.QuotaNode, number); err != nil {
		common.LogDirtyData(err,
			log.Any("type", "QuotaRelease"),
			log.Any(common.KeyContextNamespace, namespace),
			log.Any("name", plugin.QuotaNode),
			log.Any("quota", number))
		return err
	}

	return nil
}

func (api *API) createQuotas(namespace string, quotas map[string]int) error {
	if err := api.License.CreateQuota(namespace, quotas); err != nil {
		common.LogDirtyData(err,
			log.Any("type", "CreateQuotas"),
			log.Any(common.KeyContextNamespace, namespace),
			log.Any("quotas", quotas))
		return err
	}

	return nil
}
