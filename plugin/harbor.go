package plugin

import (
	"io"

	"github.com/baetyl/baetyl-cloud/v2/common"
)

//go:generate mockgen -destination=../mock/plugin/harbor.go -package=plugin github.com/baetyl/baetyl-cloud/v2/plugin Harbor

type Harbor interface {
	GetImageDigest(c *common.Context) (string, error)
	io.Closer
}
