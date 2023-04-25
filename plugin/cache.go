package plugin

import (
	"io"
)

//go:generate mockgen -destination=../mock/plugin/cache.go -package=plugin github.com/baetyl/baetyl-cloud/v2/plugin DataCache

type DataCache interface {
	Set(key string, value string) error
	Get(key string) (string, error)
	Delete(key string) error
	Exist(key string) (bool, error)
	io.Closer
}
