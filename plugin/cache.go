package plugin

import (
	"github.com/baetyl/baetyl-cloud/models"
	"io"
)

//go:generate mockgen -destination=../mock/plugin/cache.go -package=plugin github.com/baetyl/baetyl-cloud/plugin CacheStorage

type CacheStorage interface {
	GetCache(key string) (string, error)
	SetCache(key, value string) error

	DeleteCache(key string) error
	ListCache(page *models.Filter) (*models.AmisListView, error) //Pagination

	io.Closer
}
