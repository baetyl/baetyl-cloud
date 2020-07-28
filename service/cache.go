package service

import (
	"time"

	"github.com/baetyl/baetyl-cloud/models"
	"github.com/gin-contrib/cache/persistence"
)

//go:generate mockgen -destination=../mock/service/cache.go -package=plugin github.com/baetyl/baetyl-cloud/service CacheService

var ExpireDuration = time.Minute * 10

type CacheService interface {
	Get(key string) (string, error)
	Set(key, value string) error
}

type cacheService struct {
	cache           persistence.CacheStore
	propertyService PropertyService
}

func NewCacheService(propertyService PropertyService) (CacheService, error) {
	return &cacheService{
		cache:           persistence.NewInMemoryStore(ExpireDuration),
		propertyService: propertyService,
	}, nil
}

func (c *cacheService) get(GetProperty func(PropertyService, string) (*models.Property, error), propertyService PropertyService, key string) (string, error) {
	var value string
	err := c.cache.Get(key, &value)
	if err != persistence.ErrCacheMiss {
		return value, nil
	}

	property, err := GetProperty(propertyService, key)
	if err != nil {
		return "", err
	}
	err = c.Set(key, property.Value)
	return property.Value, err
}

func (c *cacheService) Get(k string) (string, error) {
	return c.get(PropertyService.GetProperty, c.propertyService, k)
}

func (c *cacheService) Set(key, value string) error {
	return c.cache.Set(key, value, ExpireDuration)
}
