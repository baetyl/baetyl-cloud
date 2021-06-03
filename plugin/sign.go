package plugin

import (
	"io"
)

//go:generate mockgen -destination=../mock/plugin/sign.go -package=plugin github.com/baetyl/baetyl-cloud/v2/plugin Sign

// Sign interfaces of Sign
type Sign interface {
	Signature(meta []byte) ([]byte, error)
	Verify(meta, sign []byte) bool
	io.Closer
}
