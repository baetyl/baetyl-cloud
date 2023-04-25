package localcache

type FastCacheConfig struct {
}

type CloudConfig struct {
	FastCacheConfig struct {
		MaxBytes int `json:"maxBytes" default:"58720256" yaml:"maxBytes"` //default 56 * 1024 * 1024  56m
	} `yaml:"fastCacheConfig" json:"fastCacheConfig" default:"{}"`
}
