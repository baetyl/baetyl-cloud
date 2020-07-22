package plugin

import (
	"github.com/baetyl/baetyl-cloud/models"
	"io"
)

//go:generate mockgen -destination=../mock/plugin/cache.go -package=plugin github.com/baetyl/baetyl-cloud/plugin Property

type Property interface {
	GetProperty(key string) (*models.Property, error)
	CreateProperty(property *models.Property) (*models.Property, error)
	UpdateProperty(property *models.Property) (*models.Property, error)

	DeleteProperty(key string) error
	ListProperty(page *models.Filter) (*models.AmisListView, error) //Pagination

	io.Closer
}
