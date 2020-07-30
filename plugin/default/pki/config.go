package pki

import "time"

// CloudConfig baetyl-cloud config
type CloudConfig struct {
	PKI struct {
		RootCAFile    string        `yaml:"rootCAFile" json:"rootCAFile" validate:"nonzero"`
		RootCAKeyFile string        `yaml:"rootCAKeyFile" json:"rootCAKeyFile" validate:"nonzero"`
		SubDuration   time.Duration `yaml:"subDuration" json:"subDuration" default:"175200h"`   // 20*365*24
		RootDuration  time.Duration `yaml:"rootDuration" json:"rootDuration" default:"438000h"` // 50*365*24
		Persistent    string        `yaml:"persistent" json:"persistent" default:"database"`
	} `yaml:"defaultpki" json:"defaultpki"`
}
