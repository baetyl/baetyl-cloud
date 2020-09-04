package service

import (
	"io/ioutil"
	"time"

	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/gin-contrib/cache/persistence"

	"github.com/baetyl/baetyl-cloud/v2/config"
)

//go:generate mockgen -destination=../mock/service/cache.go -package=service github.com/baetyl/baetyl-cloud/v2/service CacheService

type CacheService interface {
	Get(key string, load func(string) (string, error)) (string, error)
	GetProperty(key string) (string, error)
	GetFileData(file string) (string, error)
}

type CacheServiceImpl struct {
	expireDuration time.Duration
	cache          persistence.CacheStore

	prop PropertyService // default backend
}

func NewCacheService(cfg *config.CloudConfig) (CacheService, error) {
	propertyService, err := NewPropertyService(cfg)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return &CacheServiceImpl{
		expireDuration: cfg.Cache.ExpirationDuration,
		cache:          persistence.NewInMemoryStore(cfg.Cache.ExpirationDuration),
		prop:           propertyService,
	}, nil
}

func (s *CacheServiceImpl) Get(key string, load func(string) (string, error)) (string, error) {
	var value string
	if err := s.cache.Get(key, &value); err == nil {
		return value, nil
	}
	value, err := load(key)
	if err != nil {
		return "", errors.Trace(err)
	}
	return value, nil
}

func (s *CacheServiceImpl) GetProperty(key string) (string, error) {
	return s.Get(key, s.prop.GetPropertyValue)
}

func (s *CacheServiceImpl) GetFileData(file string) (string, error) {
	return s.Get(file, func(key string) (string, error) {
		data, err := ioutil.ReadFile(key)
		if err != nil {
			return "", errors.Trace(err)
		}
		return string(data), nil
	})
}
