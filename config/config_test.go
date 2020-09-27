package config

import (
	"os"
	"testing"
	"time"

	"github.com/baetyl/baetyl-go/v2/utils"
	"github.com/stretchr/testify/assert"
)

func TestSetPortFromEnv(t *testing.T) {
	cfg := &CloudConfig{
		InitServer:  Server{Port: ":1"},
		AdminServer: Server{Port: ":2"},
		MisServer:   MisServer{},
	}
	cfg.MisServer.Port = ":4"
	// no env
	SetPortFromEnv(cfg)
	assert.Equal(t, ":1", cfg.InitServer.Port)
	assert.Equal(t, ":2", cfg.AdminServer.Port)
	assert.Equal(t, ":4", cfg.MisServer.Port)

	// env
	err := os.Setenv(InitServerPort, "4")
	assert.NoError(t, err)
	err = os.Setenv(AdminServerPort, "5")
	assert.NoError(t, err)
	err = os.Setenv(MisServerPort, "7")
	assert.NoError(t, err)

	SetPortFromEnv(cfg)
	assert.Equal(t, ":4", cfg.InitServer.Port)
	assert.Equal(t, ":5", cfg.AdminServer.Port)
	assert.Equal(t, ":7", cfg.MisServer.Port)
}

func TestDefaultValue(t *testing.T) {
	expect := &CloudConfig{}
	expect.InitServer.Port = ":9003"
	expect.InitServer.WriteTimeout = time.Second * 30
	expect.InitServer.ReadTimeout = time.Second * 30
	expect.InitServer.ShutdownTime = time.Second * 3

	expect.AdminServer.Port = ":9004"
	expect.AdminServer.WriteTimeout = time.Second * 30
	expect.AdminServer.ReadTimeout = time.Second * 30
	expect.AdminServer.ShutdownTime = time.Second * 3

	expect.MisServer.Port = ":9006"
	expect.MisServer.WriteTimeout = time.Second * 30
	expect.MisServer.ReadTimeout = time.Second * 30
	expect.MisServer.ShutdownTime = time.Second * 3
	expect.MisServer.AuthToken = "baetyl-cloud-token"
	expect.MisServer.TokenHeader = "baetyl-cloud-token"
	expect.MisServer.UserHeader = "baetyl-cloud-user"

	expect.LogInfo.Level = "info"
	expect.LogInfo.MaxAge = 15
	expect.LogInfo.MaxSize = 50
	expect.LogInfo.MaxBackups = 15
	expect.LogInfo.Encoding = "json"

	expect.Plugin.PKI = "defaultpki"
	expect.Plugin.Auth = "defaultauth"
	expect.Plugin.License = "defaultlicense"
	expect.Plugin.DatabaseStorage = "database"
	expect.Plugin.ModelStorage = "kube"
	expect.Plugin.Shadow = "database"
	expect.Plugin.Functions = []string{}
	expect.Plugin.Objects = []string{}
	expect.Plugin.Property = "database"
	expect.Plugin.SyncLinks = []string{"httplink"}
	expect.Plugin.MQ = "defaultmq"

	expect.Template.Path = "/etc/baetyl/templates"

	expect.Cache.ExpirationDuration = time.Minute * 10
	// case 0
	cfg := &CloudConfig{}
	err := utils.UnmarshalYAML(nil, cfg)
	assert.NoError(t, err)
	assert.EqualValues(t, expect, cfg)

	// case 1
	cfg = &CloudConfig{}
	in := `
adminServer:
  port: ":9993"

initServer:
  port: ":9995"

misServer:
  port: ":9996"
`
	expect.AdminServer.Port = ":9993"
	expect.InitServer.Port = ":9995"
	expect.MisServer.Port = ":9996"
	err = utils.UnmarshalYAML([]byte(in), cfg)
	assert.NoError(t, err)

	assert.EqualValues(t, expect, cfg)
}
