package plugin

import (
	"io"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

//go:generate mockgen -destination=../mock/plugin/property.go -package=plugin github.com/baetyl/baetyl-cloud/v2/plugin Property

type Property interface {
	GetProperty(name string) (*models.Property, error)
	CreateProperty(property *models.Property) error
	DeleteProperty(name string) error
	ListProperty(page *models.Filter) ([]models.Property, error) //Pagination
	CountProperty(name string) (int, error)
	UpdateProperty(property *models.Property) error

	GetPropertyValue(name string) (string, error)

	io.Closer
}
