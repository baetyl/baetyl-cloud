package quota

import (
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

type quota struct {
}

func init() {
	plugin.RegisterFactory("defaultquota", New)
}

func New() (plugin.Plugin, error) {
	return &quota{}, nil
}

var _ plugin.Quota = &quota{}

func (l *quota) GetQuota(namespace string) (map[string]int, error) {
	return map[string]int{}, nil
}

func (l *quota) GetDefaultQuotas(namespace string) (map[string]int, error) {
	return map[string]int{}, nil
}

func (l *quota) CreateQuota(namespace string, quotas map[string]int) error {
	return nil
}

func (l *quota) UpdateQuota(namespace, quotaName string, quota int) error {
	return nil
}

func (l *quota) AcquireQuota(namespace, quotaName string, number int) error {
	return nil
}

func (l *quota) ReleaseQuota(namespace, quotaName string, number int) error {
	return nil
}

func (l *quota) DeleteQuota(namespace, quotaName string) error {
	return nil
}

func (l *quota) DeleteQuotaByNamespace(namespace string) error {
	return nil
}

func (l *quota) Close() error {
	return nil
}
