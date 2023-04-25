package redis

import (
	"time"

	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/gomodule/redigo/redis"
	"github.com/mna/redisc"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

type cacheRedis struct {
	redisType     TypeRedis
	clusterClient *redisc.Cluster
	singleClient  *redis.Pool
	password      string
}

func init() {
	plugin.RegisterFactory("rediscache", New)
}

func New() (plugin.Plugin, error) {
	var cfg CloudConfig
	if err := common.LoadConfig(&cfg); err != nil {
		return nil, errors.Trace(err)
	}
	redisImpl := &cacheRedis{
		redisType: cfg.CacheRedisConfig.TypeRedis,
	}
	switch cfg.CacheRedisConfig.TypeRedis {
	case ClusterRedis:
		clusterClient, err := CreateClusterRedis(cfg)
		if err != nil {
			return nil, errors.Trace(err)
		}
		redisImpl.clusterClient = clusterClient
	case SingleRedis:
		singleClient, err := CreateSingleRedis(cfg)
		if err != nil {
			return nil, errors.Trace(err)
		}
		redisImpl.singleClient = singleClient
	default:
		return nil, errors.Trace(errors.New("config typeRedis not in [single,cluster]"))
	}
	return redisImpl, nil
}

func CreateClusterRedis(cfg CloudConfig) (*redisc.Cluster, error) {
	cluster := &redisc.Cluster{
		StartupNodes: cfg.CacheRedisConfig.ClusterAddr,
		DialOptions:  []redis.DialOption{redis.DialConnectTimeout(cfg.CacheRedisConfig.IdleTimeout * time.Second)},
		CreatePool: func(addr string, opts ...redis.DialOption) (*redis.Pool, error) {
			return &redis.Pool{
				MaxIdle:     cfg.CacheRedisConfig.MaxIdle,
				MaxActive:   cfg.CacheRedisConfig.MaxActive,
				IdleTimeout: cfg.CacheRedisConfig.IdleTimeout * time.Second,
				Dial: func() (redis.Conn, error) {
					c, err := redis.Dial("tcp", cfg.CacheRedisConfig.Addr)
					if err != nil {
						return nil, err
					}
					if cfg.CacheRedisConfig.Password != "" {
						if _, err := c.Do("AUTH", cfg.CacheRedisConfig.Password); err != nil {
							c.Close()
							return nil, err
						}
					}
					if _, err := c.Do("SELECT", cfg.CacheRedisConfig.Db); err != nil {
						c.Close()
						return nil, err
					}
					return c, nil
				},
				TestOnBorrow: func(c redis.Conn, t time.Time) error {
					_, err := c.Do("PING")
					return err
				},
			}, nil
		},
	}

	return cluster, nil
}

func CreateSingleRedis(cfg CloudConfig) (*redis.Pool, error) {
	return &redis.Pool{
		MaxIdle:     cfg.CacheRedisConfig.MaxIdle,
		MaxActive:   cfg.CacheRedisConfig.MaxActive,
		IdleTimeout: cfg.CacheRedisConfig.IdleTimeout * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", cfg.CacheRedisConfig.Addr)
			if err != nil {
				return nil, err
			}
			if cfg.CacheRedisConfig.Password != "" {
				if _, err := c.Do("AUTH", cfg.CacheRedisConfig.Password); err != nil {
					c.Close()
					return nil, err
				}
			}
			if _, err := c.Do("SELECT", cfg.CacheRedisConfig.Db); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}, nil
}

func (c cacheRedis) getClient() (client redis.Conn, err error) {
	switch c.redisType {
	case SingleRedis:
		client = c.singleClient.Get()
	case ClusterRedis:
		client = c.clusterClient.Get()
	default:
		return client, errors.New("redisType get client error")
	}

	return client, err
}

func (c cacheRedis) Set(key string, value string) error {
	client, err := c.getClient()
	if err != nil {
		return err
	}
	defer func() {
		_ = client.Close()
	}()
	_, err = client.Do("SET", key, value)
	return err
}

func (c cacheRedis) Get(key string) (string, error) {
	client, err := c.getClient()
	if err != nil {
		return "", err
	}
	defer func() {
		_ = client.Close()
	}()
	return redis.String(client.Do("GET", key))

}

func (c cacheRedis) Delete(key string) error {
	client, err := c.getClient()
	if err != nil {
		return err
	}

	defer func() {
		_ = client.Close()
	}()
	_, err = client.Do("DEL", key)
	return err
}

func (c cacheRedis) Exist(key string) (bool, error) {
	client, err := c.getClient()
	if err != nil {
		return false, err
	}
	defer func() {
		_ = client.Close()
	}()
	return redis.Bool(client.Do("EXISTS", key))
}

func (c cacheRedis) Close() error {
	return nil
}
