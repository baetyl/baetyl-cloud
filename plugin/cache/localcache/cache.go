// Package localcache  freecache implements a local cache for files
package localcache

import (
	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/baetyl/baetyl-go/v2/log"
	"github.com/coocood/freecache"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

type localFreeCache struct {
	c *freecache.Cache
}

func init() {
	plugin.RegisterFactory("freecache", New)
}

func New() (plugin.Plugin, error) {
	var cfg CloudConfig
	if err := common.LoadConfig(&cfg); err != nil {
		return nil, errors.Trace(err)
	}
	cache := freecache.NewCache(cfg.FreeCacheConfig.MaxBytes)
	log.L().Info("cache set maxBytesSize", log.Any("size", cfg.FreeCacheConfig.MaxBytes))
	return &localFreeCache{
		c: cache,
	}, nil
}

func (f localFreeCache) SetString(key string, value string) error {
	return f.c.Set([]byte(key), []byte(value), 0)
}

func (f localFreeCache) Exist(key string) (bool, error) {
	_, err := f.c.Get([]byte(key))
	if err != nil {
		if err == freecache.ErrNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (f localFreeCache) GetString(key string) (string, error) {
	getData, err := f.c.Get([]byte(key))
	if err != nil {
		return "", err
	}
	return string(getData), err
}

func (f localFreeCache) SetByte(key string, value []byte) error {
	return f.c.Set([]byte(key), value, 0)
}

func (f localFreeCache) GetByte(key string) ([]byte, error) {
	return f.c.Get([]byte(key))
}

func (f localFreeCache) Delete(key string) error {
	f.c.Del([]byte(key))
	return nil
}

func (f localFreeCache) Close() error {
	return nil
}
