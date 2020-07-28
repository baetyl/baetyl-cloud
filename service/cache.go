package service

import (
	"github.com/baetyl/baetyl-cloud/config"
	"time"

	"github.com/baetyl/baetyl-cloud/models"
	"github.com/gin-contrib/cache/persistence"
)

//go:generate mockgen -destination=../mock/service/cache.go -package=plugin github.com/baetyl/baetyl-cloud/service CacheService

type CacheService struct {
	cache          persistence.CacheStore
	expireDuration time.Duration
}

func NewCacheService(config *config.CloudConfig) (CacheService, error) {
	return CacheService{
		cache:          persistence.NewInMemoryStore(config.CacheExpirationDuration),
		expireDuration: config.CacheExpirationDuration,
	}, nil
}

func (c *CacheService) Get(key string, load func(string) (*models.Property, error)) (string, error) {
	var value string
	if err := c.cache.Get(key, &value); err == nil {
		return value, nil
	}
	property, err := load(key)
	if err != nil {
		return "", err
	}
	return property.Value, c.Set(key, property.Value)
}
func (c *CacheService) Set(key, value string) error {
	return c.cache.Set(key, value, c.expireDuration)
}
