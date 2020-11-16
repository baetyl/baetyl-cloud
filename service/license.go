package service

import (
	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

//go:generate mockgen -destination=../mock/service/license.go -package=service github.com/baetyl/baetyl-cloud/v2/service LicenseService
type LicenseService interface {
	ProtectCode() error
	CheckLicense() error
	CheckQuota(namespace string, collector plugin.QuotaCollector) error
	GetQuota(namespace string) (map[string]int, error)
	GetDefaultQuotas(namespace string) (map[string]int, error)
	CreateQuota(namespace string, quotas map[string]int) error
	UpdateQuota(namespace, quotaName string, quota int) error
	AcquireQuota(namespace, quotaName string, number int) error
	ReleaseQuota(namespace, quotaName string, number int) error
	DeleteQuota(namespace, quotaName string) error
	DeleteQuotaByNamespace(namespace string) error
}

type licenseService struct {
	license plugin.License
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

func (l *licenseService) ProtectCode() error {
	return l.license.ProtectCode()
}

func (l *licenseService) CheckLicense() error {
	return l.license.CheckLicense()
}

func (l *licenseService) GetQuota(namespace string) (map[string]int, error) {
	return l.license.GetQuota(namespace)
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

func (l *licenseService) GetDefaultQuotas(namespace string) (map[string]int, error) {
	return l.license.GetDefaultQuotas(namespace)
}

func (l *licenseService) CreateQuota(namespace string, quotas map[string]int) error {
	return l.license.CreateQuota(namespace, quotas)
}

func (l *licenseService) UpdateQuota(namespace, quotaName string, quota int) error {
	return l.license.UpdateQuota(namespace, quotaName, quota)
}

func (l *licenseService) AcquireQuota(namespace, quotaName string, number int) error {
	return l.license.AcquireQuota(namespace, quotaName, number)
}

func (l *licenseService) ReleaseQuota(namespace, quotaName string, number int) error {
	return l.license.ReleaseQuota(namespace, quotaName, number)
}

func (l *licenseService) DeleteQuota(namespace, quotaName string) error {
	return l.license.DeleteQuota(namespace, quotaName)
}

func (l *licenseService) DeleteQuotaByNamespace(namespace string) error {
	return l.license.DeleteQuotaByNamespace(namespace)
}
