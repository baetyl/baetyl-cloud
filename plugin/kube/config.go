package kube

// CloudConfig baetyl-cloud config
type CloudConfig struct {
	Kubernetes struct {
		InCluster  bool   `yaml:"inCluster" json:"inCluster" default:"false"`
		ConfigPath string `yaml:"configPath" json:"configPath" default:"etc/baetyl/kube.yml"`
		// TODO Remove from plugin
		AES struct {
			Key string `yaml:"key" json:"key" default:"baetyl2020202020"`
		} `yaml:"aes" json:"aes" default:"{}"`
	} `yaml:"kubernetes" json:"kubernetes"`
}
