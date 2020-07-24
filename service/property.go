package service

import (
	"github.com/baetyl/baetyl-cloud/config"
	"github.com/baetyl/baetyl-cloud/models"
	"github.com/baetyl/baetyl-cloud/plugin"
)

//go:generate mockgen -destination=../mock/service/property.go -package=plugin github.com/baetyl/baetyl-cloud/service CacheService

type PropertyService interface {
	CreateProperty(property *models.Property) (*models.Property, error)
	DeleteProperty(key string) error
	GetProperty(key string) (*models.Property, error)
	ListProperty(page *models.Filter) ([]models.Property, int, error) //Pagination
	UpdateProperty(property *models.Property) (*models.Property, error)
}

// NewPropertyService
func NewPropertyService(config *config.CloudConfig) (PropertyService, error) {
	ds, err := plugin.GetPlugin(config.Plugin.Property)
	if err != nil {
		return nil, err
	}
	return &propertyService{
		cfg:       config,
		dbStorage: ds.(plugin.Property),
	}, nil
}

type propertyService struct {
	cfg       *config.CloudConfig
	dbStorage plugin.Property
}

func (s *propertyService) CreateProperty(property *models.Property) (*models.Property, error) {
	return s.dbStorage.CreateProperty(property)
}

func (s *propertyService) DeleteProperty(key string) error {
	return s.dbStorage.DeleteProperty(key)
}

func (s *propertyService) GetProperty(key string) (*models.Property, error) {
	return s.dbStorage.GetProperty(key)
}

func (s *propertyService) ListProperty(page *models.Filter) ([]models.Property, int, error) {
	return s.dbStorage.ListProperty(page)
}

func (s *propertyService) UpdateProperty(property *models.Property) (*models.Property, error) {
	return s.dbStorage.UpdateProperty(property)
}
