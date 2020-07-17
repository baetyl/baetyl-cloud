package config

import (
	"os"
	"time"

	"github.com/baetyl/baetyl-go/v2/log"
	"github.com/baetyl/baetyl-go/v2/utils"
)

const (
	AdminServerPort  = "ADMIN_PORT"
	NodeServerPort   = "NODE_PORT"
	ActiveServerPort = "ACTIVE_PORT"
	MisServerPort    = "MIS_PORT"
)

// CloudConfig baetyl-cloud config
type CloudConfig struct {
	ActiveServer Server     `yaml:"activeServer" json:"activeServer" default:"{\"port\":\":9003\",\"readTimeout\":30000000000,\"writeTimeout\":30000000000,\"shutdownTime\":3000000000}"`
	AdminServer  Server     `yaml:"adminServer" json:"adminServer" default:"{\"port\":\":9004\",\"readTimeout\":30000000000,\"writeTimeout\":30000000000,\"shutdownTime\":3000000000}"`
	NodeServer   NodeServer `yaml:"nodeServer" json:"nodeServer" default:"{\"port\":\":9005\",\"readTimeout\":30000000000,\"writeTimeout\":30000000000,\"shutdownTime\":3000000000,\"commonName\":\"common-name\"}"`
	MisServer    Server     `yaml:"misServer" json:"misServer" default:"{\"port\":\":9006\",\"readTimeout\":30000000000,\"writeTimeout\":30000000000,\"shutdownTime\":3000000000}"`
	LogInfo      log.Config `yaml:"logger" json:"logger"`
	Plugin       struct {
		PKI       string   `yaml:"pki" json:"pki" default:"defaultpki"`
		Auth      string   `yaml:"auth" json:"auth" default:"defaultauth"`
		License   string   `yaml:"license" json:"license" default:"defaultlicense"`
		Shadow    string   `yaml:"shadow" json:"shadow" default:"database"`
		Objects   []string `yaml:"objects" json:"objects" default:"[]"`
		Functions []string `yaml:"functions" json:"functions" default:"[]"`

		// TODO: deprecated
		CacheStorage string `yaml:"cacheStorage" json:"cacheStorage" default:"database"`
		ModelStorage    string `yaml:"modelStorage" json:"modelStorage" default:"kubernetes"`
		DatabaseStorage string `yaml:"databaseStorage" json:"databaseStorage" default:"database"`
	} `yaml:"plugin" json:"plugin"`
}

type NodeServer struct {
	Server     `yaml:",inline" json:",inline"`
	CommonName string `yaml:"commonName" json:"commonName" default:"common-name"`
}

// Server server config
type Server struct {
	Port         string            `yaml:"port" json:"port"`
	ReadTimeout  time.Duration     `yaml:"readTimeout" json:"readTimeout" default:"30s"`
	WriteTimeout time.Duration     `yaml:"writeTimeout" json:"writeTimeout" default:"30s"`
	ShutdownTime time.Duration     `yaml:"shutdownTime" json:"shutdownTime" default:"3s"`
	Certificate  utils.Certificate `yaml:",inline" json:",inline"`
}

func SetPortFromEnv(cfg *CloudConfig) {
	adminPort := os.Getenv(AdminServerPort)
	if adminPort != "" {
		cfg.AdminServer.Port = ":" + adminPort
	}
	activePort := os.Getenv(ActiveServerPort)
	if activePort != "" {
		cfg.ActiveServer.Port = ":" + activePort
	}
	nodePort := os.Getenv(NodeServerPort)
	if nodePort != "" {
		cfg.NodeServer.Port = ":" + nodePort
	}
	misPort := os.Getenv(MisServerPort)
	if misPort != "" {
		cfg.MisServer.Port = ":" + misPort
	}
}
