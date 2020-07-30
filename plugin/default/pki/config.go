package pki

import "time"

// CloudConfig baetyl-cloud config
type CloudConfig struct {
	PKI struct {
		RootCAFile    string        `yaml:"rootCAFile" json:"rootCAFile" validate:"nonzero"`
		RootCAKeyFile string        `yaml:"rootCAKeyFile" json:"rootCAKeyFile" validate:"nonzero"`
		RootCertId    string        `yaml:"rootCertId" json:"rootCertId" validate:"nonzero"`
		SubDuration   time.Duration `yaml:"subDuration" json:"subDuration" default:"175200h"`   // 20*365*24
		RootDuration  time.Duration `yaml:"rootDuration" json:"rootDuration" default:"438000h"` // 50*365*24
		Persistent    Persistent    `yaml:"persistent" json:"persistent" validate:"nonzero"`
	} `yaml:"defaultpki" json:"defaultpki"`
}

type Persistent struct {
	Kind     string `yaml:"kind" json:"kind" validate:"nonzero"`
	Database struct {
		Type            string `yaml:"type" json:"type" validate:"nonzero"`
		URL             string `yaml:"url" json:"url" validate:"nonzero"`
		MaxConns        int    `yaml:"maxConns" json:"maxConns" default:"20"`
		MaxIdleConns    int    `yaml:"maxIdleConns" json:"maxIdleConns" default:"5"`
		ConnMaxLifetime int    `yaml:"connMaxLifetime" json:"connMaxLifetime" default:"150"`
	} `yaml:"database" json:"database" default:"{}"`
}
