package plugin

import (
	"github.com/baetyl/baetyl-cloud/models"
	"io"
)

//go:generate mockgen -destination=../mock/plugin/cache.go -package=plugin github.com/baetyl/baetyl-cloud/plugin Property

type Property interface {
	CreateProperty(property *models.Property) error
	DeleteProperty(key string) error
	GetProperty(key string) (*models.Property, error)
	ListProperty(page *models.Filter) ([]models.Property, error) //Pagination
	CountProperty(key string) (int, error)
	UpdateProperty(property *models.Property) error

	io.Closer
}
