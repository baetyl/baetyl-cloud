package plugin

import (
	"io"

	"github.com/baetyl/baetyl-cloud/v2/common"
)

type Harbor interface {
	GetImageDigest(c *common.Context) (string, error)
	io.Closer
}
