package plugin

import (
	"io"
)

//go:generate mockgen -destination=../mock/plugin/cache.go -package=plugin github.com/baetyl/baetyl-cloud/v2/plugin DataCache

type DataCache interface {
	SetByte(key string, value []byte) error
	GetByte(key string) ([]byte, error)
	SetString(key string, value string) error
	GetString(key string) (string, error)
	Delete(key string) error
	Exist(key string) (bool, error)
	io.Closer
}
