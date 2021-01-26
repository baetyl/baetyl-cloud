package service

import (
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

//go:generate mockgen -destination=../mock/service/property.go -package=service github.com/baetyl/baetyl-cloud/v2/service PropertyService

type PropertyService interface {
	GetProperty(name string) (*models.Property, error)
	CreateProperty(property *models.Property) error
	DeleteProperty(name string) error
	ListProperty(page *models.Filter) ([]models.Property, error) //Pagination
	CountProperty(name string) (int, error)
	UpdateProperty(property *models.Property) error

	GetPropertyValue(name string) (string, error)
}

type propertyService struct {
	property plugin.Property
}

// NewPropertyService
func NewPropertyService(config *config.CloudConfig) (PropertyService, error) {
	ds, err := plugin.GetPlugin(config.Plugin.Property)
	if err != nil {
		return nil, err
	}

	p := &propertyService{
		property: ds.(plugin.Property),
	}

	return p, nil
}

func (p *propertyService) GetProperty(name string) (*models.Property, error) {
	return p.property.GetProperty(name)
}

func (p *propertyService) CreateProperty(property *models.Property) error {
	return p.property.CreateProperty(property)
}

func (p *propertyService) DeleteProperty(name string) error {
	return p.property.DeleteProperty(name)
}

func (p *propertyService) ListProperty(page *models.Filter) ([]models.Property, error) {
	return p.property.ListProperty(page)
}

func (p *propertyService) CountProperty(name string) (int, error) {
	return p.property.CountProperty(name)
}

func (p *propertyService) UpdateProperty(property *models.Property) error {
	return p.property.UpdateProperty(property)
}

func (p *propertyService) GetPropertyValue(name string) (string, error) {
	return p.property.GetPropertyValue(name)
}
