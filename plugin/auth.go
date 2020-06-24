package plugin

import (
	"github.com/baetyl/baetyl-cloud/common"
	"io"
)

//go:generate mockgen -destination=../mock/plugin/auth.go -package=plugin github.com/baetyl/baetyl-cloud/plugin Auth

// Auth interfaces of auth
type Auth interface {
	Authenticate(c *common.Context) error
	SignToken(meta []byte) ([]byte, error)
	VerifyToken(meta, sign []byte) bool
	io.Closer
}
