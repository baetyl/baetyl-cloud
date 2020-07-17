package service

import (
	"github.com/baetyl/baetyl-cloud/common"
	"github.com/baetyl/baetyl-cloud/config"
	"github.com/baetyl/baetyl-cloud/models"
	"github.com/baetyl/baetyl-cloud/plugin"
)

//go:generate mockgen -destination=../mock/service/cache.go -package=plugin github.com/baetyl/baetyl-cloud/service CacheService

type CacheService interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{})(interface{}, error)

	GetSystemConfig(key string) (*models.SystemConfig, error)
	ListSystemConfig(page *models.Filter) (*models.ListView, error) //分页
	CreateSystemConfig(sysConfig *models.SystemConfig) (*models.SystemConfig, error)
	UpdateSystemConfig(sysConfig *models.SystemConfig) (*models.SystemConfig, error)
	DeleteSystemConfig(key string) error
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

func (s *cacheService) Get(key string) (interface{}, error) {
	return s.dbStorage.GetSystemConfig(key)
}

func (s *cacheService) Set(key string,value interface{}) (interface{}, error){
	_, err := s.dbStorage.UpdateSystemConfig(value.(*models.SystemConfig))
	if err != nil {
		return nil, err
	}
	res, err := s.dbStorage.GetSystemConfig(key)
	if err != nil {
		return nil, common.Error(common.ErrDatabase, common.Field("error", err))
	}
	return res, nil
}

func (s *cacheService) GetSystemConfig(key string) (*models.SystemConfig, error) {
	return s.dbStorage.GetSystemConfig(key)
}

func (s *cacheService) ListSystemConfig(page *models.Filter) (*models.ListView, error) {
	systemConfigs, err := s.dbStorage.ListSystemConfig(page.Name, page.PageNo, page.PageSize)
	if err != nil {
		return nil, common.Error(common.ErrDatabase, common.Field("error", err))
	}
	count, err := s.dbStorage.CountSystemConfig(page.Name)
	if err != nil {
		return nil, common.Error(common.ErrDatabase, common.Field("error", err))
	}
	return &models.ListView{
		Total:    count,
		PageNo:   page.PageNo,
		PageSize: page.PageSize,
		Items:    systemConfigs,
	}, nil
}

func (s *cacheService) CreateSystemConfig(systemConfig *models.SystemConfig) (*models.SystemConfig, error) {
	_, err := s.dbStorage.CreateSystemConfig(systemConfig)
	if err != nil {
		return nil, err
	}
	res, err := s.dbStorage.GetSystemConfig(systemConfig.Key)
	if err != nil {
		return nil, common.Error(common.ErrDatabase, common.Field("error", err))
	}
	return res, nil
}

func (s *cacheService) DeleteSystemConfig(key string) error {
	_, err := s.dbStorage.DeleteSystemConfig(key)
	if err != nil {
		return err
	}
	return nil
}

func (s *cacheService) UpdateSystemConfig(systemConfig *models.SystemConfig) (*models.SystemConfig, error) {
	_, err := s.dbStorage.UpdateSystemConfig(systemConfig)
	if err != nil {
		return nil, err
	}
	res, err := s.dbStorage.GetSystemConfig(systemConfig.Key)
	if err != nil {
		return nil, common.Error(common.ErrDatabase, common.Field("error", err))
	}
	return res, nil
}
