package service

import (
	"time"

	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/gin-contrib/cache/persistence"
)

//go:generate mockgen -destination=../mock/service/cache.go -package=plugin github.com/baetyl/baetyl-cloud/service CacheService

type CacheService struct {
	cache          persistence.CacheStore
	expireDuration time.Duration
}

func NewCacheService(config *config.CloudConfig) (CacheService, error) {
	return CacheService{
		cache:          persistence.NewInMemoryStore(config.Cache.ExpirationDuration),
		expireDuration: config.Cache.ExpirationDuration,
	}, nil
}

func (c *CacheService) Get(key string, load func(string) (string, error)) (string, error) {
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
