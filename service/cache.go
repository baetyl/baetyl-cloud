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

	List(page *models.Filter) (*models.ListView, error) //分页
	Delete(key string) error
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
	//读完值后判断是否需要扩展
	return s.dbStorage.GetCache(key)
}
func (s *cacheService) Set(key string,value interface{}) (interface{}, error){
	oldCache, err := s.dbStorage.GetCache(key)
	if err != nil {
		return nil, err
	}
	if oldCache != nil{
		_, err := s.dbStorage.ReplaceCache(key, value.(string))
		if err != nil {
			return nil, common.Error(common.ErrDatabase, common.Field("error", err))
		}
	}else{
		_, err := s.dbStorage.AddCache(key, value.(string))
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (s *cacheService) List(page *models.Filter) (*models.ListView, error) {
	systemConfigs, err := s.dbStorage.ListCache(page.Name, page.PageNo, page.PageSize)
	if err != nil {
		return nil, common.Error(common.ErrDatabase, common.Field("error", err))
	}
	count, err := s.dbStorage.CountCache(page.Name)
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

func (s *cacheService) Delete(key string) error {
	_, err := s.dbStorage.DeleteCache(key)
	return err
}
