package service

import (
	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

//go:generate mockgen -destination=../mock/service/license.go -package=service github.com/baetyl/baetyl-cloud/v2/service LicenseService
type LicenseService interface {
	plugin.License
	CheckQuota(namespace string, collector plugin.QuotaCollector) error
}

type licenseService struct {
	plugin.License
}

func NewLicenseService(config *config.CloudConfig) (LicenseService, error) {
	l, err := plugin.GetPlugin(config.Plugin.License)
	if err != nil {
		return nil, err
	}

	return &licenseService{
		l.(plugin.License),
	}, nil
}

func (l *licenseService) CheckQuota(namespace string, collector plugin.QuotaCollector) error {
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
