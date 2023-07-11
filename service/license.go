package service

import (
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

//go:generate mockgen -destination=../mock/service/license.go -package=service github.com/baetyl/baetyl-cloud/v2/service LicenseService
type LicenseService interface {
	plugin.License
}

type LicenseServiceImpl struct {
	plugin.License
}

func NewLicenseService(config *config.CloudConfig) (LicenseService, error) {
	l, err := plugin.GetPlugin(config.Plugin.License)
	if err != nil {
		return nil, err
	}

	return &LicenseServiceImpl{
		l.(plugin.License),
	}, nil
}
