package localcache

import (
	"github.com/VictoriaMetrics/fastcache"
	"github.com/baetyl/baetyl-go/v2/errors"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

const BigSize = 65535

type localFastCache struct {
	c *fastcache.Cache
}

func init() {
	plugin.RegisterFactory("fastcache", New)
}

func New() (plugin.Plugin, error) {
	var cfg CloudConfig
	if err := common.LoadConfig(&cfg); err != nil {
		return nil, errors.Trace(err)
	}
	cache := fastcache.New(cfg.FastCacheConfig.MaxBytes)
	return &localFastCache{
		c: cache,
	}, nil
}

func (f localFastCache) SetString(key string, value string) error {
	if len(value) > BigSize {
		f.c.SetBig([]byte(key), []byte(value))
	} else {
		f.c.Set([]byte(key), []byte(value))
	}
	return nil
}

func (f localFastCache) Exist(key string) (bool, error) {
	return f.c.Has([]byte(key)), nil
}

func (f localFastCache) GetString(key string) (string, error) {
	getData := f.c.GetBig(nil, []byte(key))
	if getData == nil {
		getData = f.c.Get(nil, []byte(key))
	}
	return string(getData), nil
}

func (f localFastCache) SetByte(key string, value []byte) error {
	if len(value) > BigSize {
		f.c.SetBig([]byte(key), value)
	} else {
		f.c.Set([]byte(key), value)
	}
	return nil
}

func (f localFastCache) GetByte(key string) ([]byte, error) {
	getData := f.c.GetBig(nil, []byte(key))
	if getData == nil {
		getData = f.c.Get(nil, []byte(key))
	}
	return getData, nil
}

func (f localFastCache) Delete(key string) error {
	f.c.Del([]byte(key))
	return nil
}

func (f localFastCache) Close() error {
	return nil
}
