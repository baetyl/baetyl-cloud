package localcache

type CloudConfig struct {
	FreeCacheConfig struct {
		MaxBytes int `json:"maxBytes" default:"10485760" yaml:"maxBytes"` //default 10 * 1024 * 1024  10m
	} `yaml:"freeCacheConfig" json:"freeCacheConfig" default:"{}"`
}
