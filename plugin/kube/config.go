package kube

// CloudConfig baetyl-cloud config
type CloudConfig struct {
	Kube struct {
		OutCluster bool   `yaml:"outCluster" json:"outCluster"`
		ConfigPath string `yaml:"configPath" json:"configPath" default:"~/.kube/config"`
		// TODO Remove from plugin
		AES struct {
			Key string `yaml:"key" json:"key" default:"baetyl2020202020"`
		} `yaml:"aes" json:"aes" default:"{}"`
	} `yaml:"kube" json:"kube"`
}
