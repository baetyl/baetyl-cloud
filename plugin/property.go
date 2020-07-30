package plugin

import (
	"github.com/baetyl/baetyl-cloud/v2/models"
	"io"
)

//go:generate mockgen -destination=../mock/plugin/cache.go -package=plugin github.com/baetyl/baetyl-cloud/v2/plugin Property

type Property interface {
	CreateProperty(property *models.Property) error
	DeleteProperty(key string) error
	ListProperty(page *models.Filter) ([]models.Property, error) //Pagination
	CountProperty(key string) (int, error)
	UpdateProperty(property *models.Property) error

	GetPropertyValue(key string) (string, error)

	io.Closer
}
