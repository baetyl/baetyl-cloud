package service

import (
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

//go:generate mockgen -destination=../mock/service/property.go -package=service github.com/baetyl/baetyl-cloud/v2/service PropertyService

type PropertyService interface {
	CreateProperty(property *models.Property) error
	DeleteProperty(key string) error
	ListProperty(page *models.Filter) ([]models.Property, error) //Pagination
	CountProperty(key string) (int, error)
	UpdateProperty(property *models.Property) error

	GetPropertyValue(key string) (string, error)
}

// NewPropertyService
func NewPropertyService(config *config.CloudConfig) (PropertyService, error) {
	ds, err := plugin.GetPlugin(config.Plugin.Property)
	if err != nil {
		return nil, err
	}
	return ds.(plugin.Property), nil
}
