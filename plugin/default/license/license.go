package license

import (
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

type license struct {
}

func init() {
	plugin.RegisterFactory("defaultlicense", New)
}

func New() (plugin.Plugin, error) {
	return &license{}, nil
}

var _ plugin.License = &license{}

func (l *license) ProtectCode() error {
	return nil
}
func (l *license) CheckLicense() error {
	return nil
}

func (l *license) GetQuota(namespace string) (map[string]int, error) {
	return map[string]int{}, nil
}

func (l *license) GetDefaultQuotas(namespace string) (map[string]int, error) {
	return map[string]int{}, nil
}

func (l *license) CreateQuota(namespace string, quotas map[string]int) error {
	return nil
}

func (l *license) UpdateQuota(namespace, quotaName string, quota int) error {
	return nil
}

func (l *license) AcquireQuota(namespace, quotaName string, number int) error {
	return nil
}

func (l *license) ReleaseQuota(namespace, quotaName string, number int) error {
	return nil
}

func (l *license) DeleteQuota(namespace, quotaName string) error {
	return nil
}

func (l *license) DeleteQuotaByNamespace(namespace string) error {
	return nil
}
func (l *license) Close() error {
	return nil
}