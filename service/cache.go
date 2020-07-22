package service

import (
	"github.com/baetyl/baetyl-cloud/config"
	"github.com/baetyl/baetyl-cloud/models"
	"github.com/baetyl/baetyl-cloud/plugin"
)

//go:generate mockgen -destination=../mock/service/cache.go -package=plugin github.com/baetyl/baetyl-cloud/service CacheService

type CacheService interface {
	Get(key string) (string, error)
	Set(key, value string) error

	Delete(key string) error
	List(page *models.Filter) (*models.AmisListView, error) //分页
}

type cacheService struct {
	cfg       *config.CloudConfig
	dbStorage plugin.CacheStorage
}

// NewCacheService
func NewCacheService(config *config.CloudConfig) (CacheService, error) {
	ds, err := plugin.GetPlugin(config.Plugin.CacheStorage)
	if err != nil {
		return nil, err
	}
	return &cacheService{
		cfg:       config,
		dbStorage: ds.(plugin.CacheStorage),
	}, nil
}

func (s *cacheService) Get(key string) (string, error) {
	return s.dbStorage.GetCache(key)
}
func (s *cacheService) Set(key string, value string) error {
	return s.dbStorage.SetCache(key, value)
}

func (s *cacheService) List(page *models.Filter) (*models.AmisListView, error) {
	return s.dbStorage.ListCache(page)
}

func (s *cacheService) Delete(key string) error {
	return s.dbStorage.DeleteCache(key)
}
