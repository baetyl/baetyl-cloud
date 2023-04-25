package redis

import "time"

type TypeRedis string

const (
	SingleRedis  = "single"
	ClusterRedis = "cluster"
)

type CloudConfig struct {
	CacheRedisConfig struct {
		MaxIdle     int           `json:"maxIdle" yaml:"maxIdle"`
		MaxActive   int           `json:"maxActive" yaml:"maxActive"`
		IdleTimeout time.Duration `json:"idleTimeout" yaml:"idleTimeout"`
		TypeRedis   TypeRedis     `json:"typeRedis" yaml:"typeRedis"`
		Addr        string        `json:"addr" yaml:"addr"`
		ClusterAddr []string      `json:"clusterAddr" yaml:"clusterAddr"`
		Password    string        `json:"password" yaml:"password"`
		Db          string        `json:"db" yaml:"db" default:"0"`
	} `yaml:"cacheRedisConfig" json:"cacheRedisConfig" default:"{}"`
}
