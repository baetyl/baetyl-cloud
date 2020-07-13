package service

import (
	"strings"
	"time"

	"github.com/baetyl/baetyl-cloud/common"
	"github.com/baetyl/baetyl-cloud/config"
	"github.com/baetyl/baetyl-cloud/models"
	"github.com/baetyl/baetyl-cloud/plugin"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
)

//go:generate mockgen -destination=../mock/service/config.go -package=plugin github.com/baetyl/baetyl-cloud/service ConfigService

// ConfigService ConfigService
type ConfigService interface {
	Get(namespace, name, version string) (*specV1.Configuration, error)
	List(namespace string, listOptions *models.ListOptions) (*models.ConfigurationList, error)
	Create(namespace string, config *specV1.Configuration) (*specV1.Configuration, error)
	Update(namespace string, config *specV1.Configuration) (*specV1.Configuration, error)
	Upsert(namespace string, config *specV1.Configuration) (*specV1.Configuration, error)
	Delete(namespace, name string) error
}

type configService struct {
	storage plugin.ModelStorage
}

// NewConfigService NewConfigService
func NewConfigService(config *config.CloudConfig) (ConfigService, error) {
	ms, err := plugin.GetPlugin(config.Plugin.ModelStorage)
	if err != nil {
		return nil, err
	}
	return &configService{
		storage: ms.(plugin.ModelStorage),
	}, nil
}

// Get get a config
func (s *configService) Get(namespace, name, version string) (*specV1.Configuration, error) {
	res, err := s.storage.GetConfig(namespace, name, version)
	if err != nil && strings.Contains(err.Error(), "not found") {
		return nil, common.Error(common.ErrResourceNotFound, common.Field("type", "config"),
			common.Field("name", name))
	}
	return res, err
}

// List get list config
func (s *configService) List(namespace string, listOptions *models.ListOptions) (*models.ConfigurationList, error) {
	return s.storage.ListConfig(namespace, listOptions)
}

// Create Create a config
func (s *configService) Create(namespace string, config *specV1.Configuration) (*specV1.Configuration, error) {
	return s.storage.CreateConfig(namespace, config)
}

// Update update a config
func (s *configService) Update(namespace string, config *specV1.Configuration) (*specV1.Configuration, error) {
	return s.storage.UpdateConfig(namespace, config)
}

// Upsert update a config or create a config if not exist
func (s *configService) Upsert(namespace string, config *specV1.Configuration) (*specV1.Configuration, error) {
	res, err := s.storage.GetConfig(namespace, config.Name, "")
	if err != nil {
		return s.storage.CreateConfig(namespace, config)
	}

	if models.EqualConfig(res, config) {
		return res, nil
	}

	config.Version = res.Version
	config.UpdateTimestamp = time.Now()
	return s.storage.UpdateConfig(namespace, config)
}

// Delete Delete a config
func (s *configService) Delete(namespace, name string) error {
	return s.storage.DeleteConfig(namespace, name)
}
