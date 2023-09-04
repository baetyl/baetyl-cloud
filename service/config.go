package service

import (
	"strings"
	"time"

	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

//go:generate mockgen -destination=../mock/service/config.go -package=service github.com/baetyl/baetyl-cloud/v2/service ConfigService

// ConfigService ConfigService
type ConfigService interface {
	Get(tx interface{}, namespace, name, version string) (*specV1.Configuration, error)
	List(namespace string, listOptions *models.ListOptions) (*models.ConfigurationList, error)
	Create(tx interface{}, namespace string, config *specV1.Configuration) (*specV1.Configuration, error)
	Update(tx interface{}, namespace string, config *specV1.Configuration) (*specV1.Configuration, error)
	Upsert(tx interface{}, namespace string, config *specV1.Configuration) (*specV1.Configuration, error)
	Delete(tx interface{}, namespace, name string) error
}

type configService struct {
	config plugin.Configuration
}

// NewConfigService NewConfigService
func NewConfigService(config *config.CloudConfig) (ConfigService, error) {
	cfg, err := plugin.GetPlugin(config.Plugin.Resource)
	if err != nil {
		return nil, err
	}
	return &configService{
		config: cfg.(plugin.Configuration),
	}, nil
}

// Get get a config
func (s *configService) Get(tx interface{}, namespace, name, version string) (*specV1.Configuration, error) {
	res, err := s.config.GetConfig(tx, namespace, name, version)
	if err != nil && strings.Contains(err.Error(), "not found") {
		return nil, common.Error(common.ErrResourceNotFound, common.Field("type", "config"),
			common.Field("name", name))
	}
	return res, err
}

// List get list config
func (s *configService) List(namespace string, listOptions *models.ListOptions) (*models.ConfigurationList, error) {
	return s.config.ListConfig(namespace, listOptions)
}

// Create Create a config
func (s *configService) Create(tx interface{}, namespace string, config *specV1.Configuration) (*specV1.Configuration, error) {
	return s.config.CreateConfig(tx, namespace, config)
}

// Update update a config
func (s *configService) Update(tx interface{}, namespace string, config *specV1.Configuration) (*specV1.Configuration, error) {
	return s.config.UpdateConfig(tx, namespace, config)
}

// Upsert update a config or create a config if not exist
func (s *configService) Upsert(tx interface{}, namespace string, config *specV1.Configuration) (*specV1.Configuration, error) {
	res, err := s.config.GetConfig(tx, namespace, config.Name, "")
	if err != nil {
		return s.config.CreateConfig(tx, namespace, config)
	}

	if models.EqualConfig(res, config) {
		return res, nil
	}

	config.Version = res.Version
	config.UpdateTimestamp = time.Now()
	return s.config.UpdateConfig(tx, namespace, config)
}

// Delete Delete a config
func (s *configService) Delete(tx interface{}, namespace, name string) error {
	return s.config.DeleteConfig(tx, namespace, name)
}
