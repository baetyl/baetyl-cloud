package auth

type CloudConfig struct {
	DefaultAuth struct {
		Namespace string `yaml:"namespace" json:"namespace" default:"baetyl-cloud"`
		KeyFile   string `yaml:"keyFile" json:"keyFile" default:"etc/baetyl/token.key"`
	} `yaml:"defaultauth" json:"defaultauth"`
}
