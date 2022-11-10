package service

import (
	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

//go:generate mockgen -destination=../mock/service/quota.go -package=service github.com/baetyl/baetyl-cloud/v2/service QuotaService
type QuotaService interface {
	plugin.Quota
	CheckQuota(namespace string, collector plugin.QuotaCollector) error
}

type QuotaServiceImpl struct {
	plugin.Quota
}

func NewQuotaService(config *config.CloudConfig) (QuotaService, error) {
	l, err := plugin.GetPlugin(config.Plugin.Quota)
	if err != nil {
		return nil, err
	}

	return &QuotaServiceImpl{
		l.(plugin.Quota),
	}, nil
}

func (l *QuotaServiceImpl) CheckQuota(namespace string, collector plugin.QuotaCollector) error {
	limits, err := l.GetQuota(namespace)
	if err != nil {
		return err
	}

	counts, err := collector(namespace)
	if err != nil {
		return err
	}

	if counts == nil || limits == nil {
		return nil
	}

	for k, v := range counts {
		if limits[k] != 0 && v >= limits[k] {
			return common.Error(
				common.ErrLicenseQuota,
				common.Field("name", k),
				common.Field("limit", limits[k]))
		}
	}
	return nil
}
