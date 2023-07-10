package localcache

type FastCacheConfig struct {
}

type CloudConfig struct {
	FreeCacheConfig struct {
		MaxBytes int `json:"maxBytes" default:"10485760" yaml:"maxBytes"` //default 10 * 1024 * 1024  56m
	} `yaml:"freeCacheConfig" json:"fastCacheConfig" default:"{}"`
}
