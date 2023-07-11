package plugin

import "io"

//go:generate mockgen -destination=../mock/plugin/quota.go -package=plugin github.com/baetyl/baetyl-cloud/v2/plugin Quota

const (
	QuotaNode  = "maxNodeCount"
	QuotaBatch = "maxBatchCount"
	MenuEnable = "menuEnable"
)

type QuotaCollector func(namespace string) (map[string]int, error)

type Quota interface {
	GetQuota(namespace string) (map[string]int, error)
	GetDefaultQuotas(namespace string) (map[string]int, error)
	CreateQuota(namespace string, quotas map[string]int) error
	UpdateQuota(namespace, quotaName string, quota int) error
	AcquireQuota(namespace, quotaName string, number int) error
	ReleaseQuota(namespace, quotaName string, number int) error
	DeleteQuota(namespace, quotaName string) error
	DeleteQuotaByNamespace(namespace string) error
	io.Closer
}
