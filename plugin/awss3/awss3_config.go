package awss3

import "time"

// CloudConfig baetyl-cloud config
type CloudConfig struct {
	AWSS3 *S3Config `yaml:"awss3" json:"awss3"`
}

type S3Config struct {
	Endpoint      string        `yaml:"endpoint" json:"endpoint"`
	Ak            string        `yaml:"ak" json:"ak" binding:"nonzero"`
	Sk            string        `yaml:"sk" json:"sk" binding:"nonzero"`
	Region        string        `yaml:"region" json:"region" default:"us-east-1"`
	AddressFormat string        `yaml:"addressFormat" json:"addressFormat" default:"pathStyle"`
	Expiration    time.Duration `yaml:"expiration" json:"expiration" default:"1h"`
}
