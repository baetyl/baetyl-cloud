package models

type NodeDeviceCacheView struct {
	FilePath            string `yaml:"filePath" json:"filePath" default:"/var/lib/baetyl/store"`
	Enable              bool   `yaml:"enable" json:"enable"`
	CacheDuration       string `yaml:"cacheDuration" json:"cacheDuration" default:"24h"`
	CacheType           string `yaml:"cacheType" json:"cacheType" default:"disk"`
	Clear               bool   `yaml:"clear" json:"clear"`
	MaxMemoryUse        string `yaml:"maxMemoryUse" json:"maxMemoryUse" default:"100MB"`
	SendMessageInterval string `yaml:"sendMessageInterval" json:"sendMessageInterval" default:"3s"`
	MaxInterval         string `yaml:"maxInterval" json:"maxInterval" default:"60s"`
	Interval            string `yaml:"interval" json:"interval" default:"3s"`
	MinInterval         string `yaml:"minInterval" json:"minInterval" default:"2s"`
	FailChanSize        int    `yaml:"failChanSize" json:"failChanSize" default:"100"`
}
