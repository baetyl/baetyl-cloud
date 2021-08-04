package plugin

import (
	"io"

	"github.com/baetyl/baetyl-cloud/v2/common"
)

//go:generate mockgen -destination=../mock/plugin/csrf.go -package=plugin github.com/baetyl/baetyl-cloud/v2/plugin CsrfValidator

type CsrfValidator interface {
	Verify(c *common.Context) error
	io.Closer
}
