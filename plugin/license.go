package plugin

import "io"

//go:generate mockgen -destination=../mock/plugin/license.go -package=plugin github.com/baetyl/baetyl-cloud/v2/plugin License

type License interface {
	ProtectCode() error
	CheckLicense() error
	io.Closer
}
