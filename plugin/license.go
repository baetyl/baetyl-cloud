package plugin

import "io"

//go:generate mockgen -destination=../mock/plugin/license.go -package=plugin github.com/baetyl/baetyl-cloud/v2/plugin License

const (
	QuotaNode  = "maxNodeCount"
	QuotaBatch = "maxBatchCount"
)

type QuotaCollector func(namespace string) (map[string]int, error)

type License interface {
	ProtectCode() error
	CheckLicense() error
	GetQuota(namespace string) (map[string]int, error)
	io.Closer
}
