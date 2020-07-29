package service

import (
	"github.com/baetyl/baetyl-cloud/config"
	"github.com/baetyl/baetyl-cloud/models"
	"github.com/baetyl/baetyl-cloud/plugin"
)

//go:generate mockgen -destination=../mock/service/property.go -package=plugin github.com/baetyl/baetyl-cloud/service PropertyService

type PropertyService interface {
	CreateProperty(property *models.Property) error
	DeleteProperty(key string) error
	GetProperty(key string) (interface{}, error)
	ListProperty(page *models.Filter) ([]models.Property, error) //Pagination
	CountProperty(key string) (int, error)
	UpdateProperty(property *models.Property) error
}

// NewPropertyService
func NewPropertyService(config *config.CloudConfig) (PropertyService, error) {
	ds, err := plugin.GetPlugin(config.Plugin.Property)
	if err != nil {
		return nil, err
	}
	return ds.(plugin.Property), nil
}
