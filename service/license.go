package service

import (
	"github.com/baetyl/baetyl-cloud/common"
	"github.com/baetyl/baetyl-cloud/config"
	"github.com/baetyl/baetyl-cloud/plugin"
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
	if p := common.GetEasyPack(); p != "" {
		// overwrite the config if easypack is enabled
		config.Plugin.License = p
	}
	l, err := plugin.GetPlugin(config.Plugin.License)
	if err != nil {
		return nil, err
	}

	return &licenseService{
		l.(plugin.License),
	}, nil
}
