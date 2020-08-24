package httplink

import (
	"github.com/baetyl/baetyl-cloud/v2/config"
)

type CloudConfig struct {
	HTTPLink SyncServer `yaml:"httpLink" json:"httpLink" default:"{\"port\":\":9005\",\"readTimeout\":30000000000,\"writeTimeout\":30000000000,\"shutdownTime\":3000000000,\"commonName\":\"common-name\"}"`
}

type SyncServer struct {
	config.Server `yaml:",inline" json:",inline"`
	CommonName    string `yaml:"commonName" json:"commonName" default:"common-name"`
}
