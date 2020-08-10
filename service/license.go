package service

import (
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

type LicenseService interface {
	ProtectCode() error
	CheckLicense() error
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
