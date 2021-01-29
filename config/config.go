package config

import (
	"time"

	"github.com/baetyl/baetyl-go/v2/log"
	"github.com/baetyl/baetyl-go/v2/utils"
)

// CloudConfig baetyl-cloud config
type CloudConfig struct {
	InitServer  Server     `yaml:"initServer" json:"initServer" default:"{\"port\":\":9003\",\"readTimeout\":30000000000,\"writeTimeout\":30000000000,\"shutdownTime\":3000000000}"`
	AdminServer Server     `yaml:"adminServer" json:"adminServer" default:"{\"port\":\":9004\",\"readTimeout\":30000000000,\"writeTimeout\":30000000000,\"shutdownTime\":3000000000}"`
	MisServer   MisServer  `yaml:"misServer" json:"misServer" default:"{\"port\":\":9006\",\"readTimeout\":30000000000,\"writeTimeout\":30000000000,\"shutdownTime\":3000000000,\"authToken\":\"baetyl-cloud-token\",\"tokenHeader\":\"baetyl-cloud-token\",\"userHeader\":\"baetyl-cloud-user\"}"`
	LogInfo     log.Config `yaml:"logger" json:"logger"`
	Cache       struct {
		ExpirationDuration time.Duration `yaml:"expirationDuration" json:"expirationDuration" default:"10m"`
	} `yaml:"cache" json:"cache"`
	Template struct {
		Path string `yaml:"path" json:"path" default:"/etc/baetyl/templates"`
	} `yaml:"template" json:"template"`
	Plugin struct {
		Pubsub     string   `yaml:"pubsub" json:"pubsub" default:"defaultpubsub"`
		PKI        string   `yaml:"pki" json:"pki" default:"defaultpki"`
		Auth       string   `yaml:"auth" json:"auth" default:"defaultauth"`
		License    string   `yaml:"license" json:"license" default:"defaultlicense"`
		Resource   string   `yaml:"resource" json:"resource" default:"kube"`
		Shadow     string   `yaml:"shadow" json:"shadow" default:"database"`
		Index      string   `yaml:"index" json:"index" default:"database"`
		Batch      string   `yaml:"batch" json:"batch" default:"databaseext"`
		Record     string   `yaml:"record" json:"record" default:"databaseext"`
		Callback   string   `yaml:"callback" json:"callback" default:"databaseext"`
		AppHistory string   `yaml:"appHistory" json:"appHistory" default:"database"`
		Objects    []string `yaml:"objects" json:"objects" default:"[]"`
		Functions  []string `yaml:"functions" json:"functions" default:"[]"`
		Property   string   `yaml:"property" json:"property" default:"database"`
		Module     string   `yaml:"module" json:"module" default:"database"`
		SyncLinks  []string `yaml:"synclinks" json:"synclinks" default:"[\"httplink\"]"`
	} `yaml:"plugin" json:"plugin"`
}

type MisServer struct {
	Server      `yaml:",inline" json:",inline"`
	AuthToken   string `yaml:"authToken" json:"authToken" default:"baetyl-cloud-token"`
	TokenHeader string `yaml:"tokenHeader" json:"tokenHeader" default:"baetyl-cloud-token"`
	UserHeader  string `yaml:"userHeader" json:"userHeader" default:"baetyl-cloud-user"`
}

// Server server config
type Server struct {
	Port         string            `yaml:"port" json:"port"`
	ReadTimeout  time.Duration     `yaml:"readTimeout" json:"readTimeout" default:"30s"`
	WriteTimeout time.Duration     `yaml:"writeTimeout" json:"writeTimeout" default:"30s"`
	ShutdownTime time.Duration     `yaml:"shutdownTime" json:"shutdownTime" default:"3s"`
	Certificate  utils.Certificate `yaml:",inline" json:",inline"`
}
