package service

import (
	"time"

	"github.com/gin-contrib/cache/persistence"

	"github.com/baetyl/baetyl-cloud/v2/config"
)

//go:generate mockgen -destination=../mock/service/cache.go -package=service github.com/baetyl/baetyl-cloud/v2/service CacheService

type CacheService interface {
	Get(key string, load func(string) (string, error)) (string, error)
}

type CacheServiceImpl struct {
	cache          persistence.CacheStore
	expireDuration time.Duration
}

func NewCacheService(config *config.CloudConfig) (CacheService, error) {
	return &CacheServiceImpl{
		cache:          persistence.NewInMemoryStore(config.Cache.ExpirationDuration),
		expireDuration: config.Cache.ExpirationDuration,
	}, nil
}

func (c *CacheServiceImpl) Get(key string, load func(string) (string, error)) (string, error) {
	var value string
	if err := c.cache.Get(key, &value); err == nil {
		return value, nil
	}
	value, err := load(key)
	if err != nil {
		return "", err
	}
	return value, nil
}
