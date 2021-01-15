package service

import (
	"encoding/json"
	"errors"
	"fmt"

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

	GetSysApp(name string) (*models.NodeSysAppInfo, error)
	CreateSysApp(property *models.NodeSysAppInfo) error
	UpdateSysApp(info *models.NodeSysAppInfo) error
	DeleteSysApp(name string) error
	ListSysApps() ([]models.NodeSysAppInfo, error)
	ListOptionalSysApps() ([]models.NodeSysAppInfo, error)
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

func (p *propertyService) GetSysApp(name string) (*models.NodeSysAppInfo, error) {
	var app models.NodeSysAppInfo
	var err error

	app.Name = name
	app.Image, err = p.property.GetPropertyValue(genSysAppImageKey(name))
	if err != nil {
		return nil, err
	}

	app.Description, err = p.property.GetPropertyValue(genSysAppDescriptionKey(name))
	if err != nil {
		return nil, err
	}

	prefix := genSysAppProgramPrefixKey(name)
	params := &models.Filter{
		Name: prefix,
	}
	res, err := p.property.ListProperty(params)
	if err != nil {
		return nil, err
	}
	programs := make(map[string]string)
	for _, item := range res {
		programs[item.Name[len(prefix):]] = item.Value
	}
	app.Programs = programs

	return &app, nil
}

func (p *propertyService) CreateSysApp(info *models.NodeSysAppInfo) error {
	property := &models.Property{
		Name:  genSysAppImageKey(info.Name),
		Value: info.Image,
	}
	err := p.property.CreateProperty(property)
	if err != nil {
		return err
	}

	property = &models.Property{
		Name:  genSysAppDescriptionKey(info.Name),
		Value: info.Description,
	}
	err = p.property.CreateProperty(property)
	if err != nil {
		return err
	}

	// update platform programs
	for platform, url := range info.Programs {
		property = &models.Property{
			Name:  genSysAppProgramKey(info.Name, platform),
			Value: url,
		}
		err = p.property.CreateProperty(property)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *propertyService) UpdateSysApp(info *models.NodeSysAppInfo) error {
	// update image
	property := &models.Property{
		Name:  genSysAppImageKey(info.Name),
		Value: info.Image,
	}
	err := p.property.UpdateProperty(property)
	if err != nil {
		return err
	}

	// update description
	property = &models.Property{
		Name:  genSysAppDescriptionKey(info.Name),
		Value: info.Description,
	}
	err = p.property.UpdateProperty(property)
	if err != nil {
		return err
	}

	// update platform programs
	for platform, url := range info.Programs {
		property = &models.Property{
			Name:  genSysAppProgramKey(info.Name, platform),
			Value: url,
		}
		err = p.property.UpdateProperty(property)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *propertyService) DeleteSysApp(name string) error {
	err := p.property.DeleteProperty(genSysAppImageKey(name))
	if err != nil {
		return err
	}

	err = p.property.DeleteProperty(genSysAppDescriptionKey(name))
	if err != nil {
		return err
	}

	prefix := genSysAppProgramPrefixKey(name)
	params := &models.Filter{
		Name: prefix,
	}
	res, err := p.property.ListProperty(params)
	if err != nil {
		return err
	}

	for _, item := range res {
		err = p.property.DeleteProperty(item.Name)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *propertyService) ListSysApps() ([]models.NodeSysAppInfo, error) {
	// TODO: how to support list sys apps, like ListOptionalSysApps() do ?
	return nil, errors.New("this method is not supported")
}

func (p *propertyService) ListOptionalSysApps() ([]models.NodeSysAppInfo, error) {
	var appNames []string
	res, err := p.property.GetPropertyValue(genOptionalSysAppKey())
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(res), &appNames)
	if err != nil {
		return nil, err
	}

	var infos []models.NodeSysAppInfo
	for _, name := range appNames {
		info, err := p.GetSysApp(name)
		if err != nil {
			return nil, err
		}
		infos = append(infos, *info)
	}
	return infos, nil
}

func genSysAppImageKey(app string) string {
	return fmt.Sprintf("%s-image", app)
}

func genSysAppDescriptionKey(app string) string {
	return fmt.Sprintf("%s-description", app)
}

func genSysAppProgramPrefixKey(app string) string {
	return fmt.Sprintf("%s-program", app)
}

func genSysAppProgramKey(app, platfrom string) string {
	return fmt.Sprintf("%s-program-%s", app, "platfrom")
}

func genOptionalSysAppKey() string {
	return fmt.Sprintf("baetyl-optional-sysapps")
}
