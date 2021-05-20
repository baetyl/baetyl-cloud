package auth

type CloudConfig struct {
	DefaultAuth struct {
		Namespace string `yaml:"namespace" json:"namespace" default:"baetyl-cloud"`
	} `yaml:"defaultauth" json:"defaultauth"`
}
