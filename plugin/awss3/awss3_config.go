package awss3

import "time"

// CloudConfig baetyl-cloud config
type CloudConfig struct {
	AWSS3 *S3Config `yaml:"awss3" json:"awss3"`
}

type S3Config struct {
	Endpoint   string        `yaml:"endpoint" json:"endpoint"`
	Ak         string        `yaml:"ak" json:"ak" validate:"nonzero"`
	Sk         string        `yaml:"sk" json:"sk" validate:"nonzero"`
	Region     string        `yaml:"region" json:"region" default:"us-east-1"`
	Expiration time.Duration `yaml:"expiration" json:"expiration" default:"1h"`
}
