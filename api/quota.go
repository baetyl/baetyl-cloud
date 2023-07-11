package api

import (
	"github.com/baetyl/baetyl-go/v2/log"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
)

// GetQuota  for admin api
func (api *API) GetQuota(c *common.Context) (interface{}, error) {
	ns := c.GetNamespace()
	quotas, err := api.License.GetQuota(ns)
	return quotas, err
}

// GetQuota for mis server api
//   - param namespace string
func (api *API) GetQuotaForMis(c *common.Context) (interface{}, error) {
	ns := c.Param(common.KeyContextNamespace)
	quota := &models.Quota{
		Namespace: ns,
	}

	quotas, err := api.License.GetQuota(quota.Namespace)

	if err != nil {
		return nil, err
	}

	var quotaDatas []models.Quota
	for k, v := range quotas {
		quotaDatas = append(quotaDatas, models.Quota{
			Namespace: ns,
			QuotaName: k,
			Quota:     v,
		})
	}

	return models.MisData{
		Count: len(quotaDatas),
		Rows:  quotaDatas,
	}, nil
}

// CreateQuota for mis server api
//   - param namespace string
//   - param quotaName string
//   - param quota int
func (api *API) CreateQuota(c *common.Context) (interface{}, error) {
	quota := &models.Quota{}
	if err := c.LoadBody(quota); err != nil {
		return nil, err
	}
	quotas := map[string]int{
		quota.QuotaName: quota.Quota,
	}
	err := api.CreateQuotas(quota.Namespace, quotas)
	return nil, err
}

// UpdateQuota for mis server api
//   - param namespace string
//   - param quotaName string
//   - param quota int
func (api *API) UpdateQuota(c *common.Context) (interface{}, error) {
	quota := &models.Quota{}
	if err := c.LoadBody(quota); err != nil {
		return nil, err
	}
	err := api.License.UpdateQuota(quota.Namespace, quota.QuotaName, quota.Quota)
	return nil, err
}

// DeleteQuota for mis server api
//   - param namespace string
//   - param quotaName string
func (api *API) DeleteQuota(c *common.Context) (interface{}, error) {
	quota := &models.Quota{}
	if err := c.LoadBody(quota); err != nil {
		return nil, err
	}
	err := api.License.DeleteQuota(quota.Namespace, quota.QuotaName)
	return nil, err
}

// InitQuotas
//   - param namespace string
func (api *API) InitQuotas(namespace string) error {
	quotas, err := api.License.GetDefaultQuotas(namespace)
	if err != nil {
		return err
	}
	return api.CreateQuotas(namespace, quotas)
}

// DeleteQuotaByNamespace
//   - param namespace string
func (api *API) DeleteQuotaByNamespace(namespace string) error {
	if err := api.License.DeleteQuotaByNamespace(namespace); err != nil {
		common.LogDirtyData(err,
			log.Any("type", "DeleteQuotaByNamespace"),
			log.Any(common.KeyContextNamespace, namespace))
		return err
	}

	return nil
}

func (api *API) ReleaseQuota(namespace, quotaName string, number int) error {
	if err := api.License.ReleaseQuota(namespace, quotaName, number); err != nil {
		common.LogDirtyData(err,
			log.Any("type", "QuotaRelease"),
			log.Any(common.KeyContextNamespace, namespace),
			log.Any("name", quotaName),
			log.Any("quota", number))
		return err
	}

	return nil
}

func (api *API) CreateQuotas(namespace string, quotas map[string]int) error {
	if err := api.License.CreateQuota(namespace, quotas); err != nil {
		common.LogDirtyData(err,
			log.Any("type", "CreateQuotas"),
			log.Any(common.KeyContextNamespace, namespace),
			log.Any("quotas", quotas))
		return err
	}

	return nil
}
