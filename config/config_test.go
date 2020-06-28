package config

import (
	"github.com/baetyl/baetyl-go/utils"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestSetPortFromEnv(t *testing.T) {
	cfg := &CloudConfig{
		ActiveServer: Server{Port: ":1"},
		AdminServer:  Server{Port: ":2"},
		NodeServer:   Server{Port: ":3"},
	}
	// no env
	SetPortFromEnv(cfg)
	assert.Equal(t, ":1", cfg.ActiveServer.Port)
	assert.Equal(t, ":2", cfg.AdminServer.Port)
	assert.Equal(t, ":3", cfg.NodeServer.Port)

	// env
	err := os.Setenv(ActiveServerPort, "4")
	assert.NoError(t, err)
	err = os.Setenv(AdminServerPort, "5")
	assert.NoError(t, err)
	err = os.Setenv(NodeServerPort, "6")
	assert.NoError(t, err)

	SetPortFromEnv(cfg)
	assert.Equal(t, ":4", cfg.ActiveServer.Port)
	assert.Equal(t, ":5", cfg.AdminServer.Port)
	assert.Equal(t, ":6", cfg.NodeServer.Port)
}

func TestDefaultValue(t *testing.T) {
	expect := &CloudConfig{}
	expect.ActiveServer.Port = ":9003"
	expect.ActiveServer.WriteTimeout = time.Second * 30
	expect.ActiveServer.ReadTimeout = time.Second * 30
	expect.ActiveServer.ShutdownTime = time.Second * 3

	expect.AdminServer.Port = ":9004"
	expect.AdminServer.WriteTimeout = time.Second * 30
	expect.AdminServer.ReadTimeout = time.Second * 30
	expect.AdminServer.ShutdownTime = time.Second * 3

	expect.NodeServer.Port = ":9005"
	expect.NodeServer.WriteTimeout = time.Second * 30
	expect.NodeServer.ReadTimeout = time.Second * 30
	expect.NodeServer.ShutdownTime = time.Second * 3

	expect.LogInfo.Level = "info"
	expect.LogInfo.MaxAge = 15
	expect.LogInfo.MaxSize = 50
	expect.LogInfo.MaxBackups = 15
	expect.LogInfo.Encoding = "json"

	expect.Plugin.PKI = "defaultpki"
	expect.Plugin.Auth = "defaultauth"
	expect.Plugin.License = "defaultlicense"
	expect.Plugin.DatabaseStorage = "database"
	expect.Plugin.ModelStorage = "kubernetes"
	expect.Plugin.Shadow = "database"
	expect.Plugin.Functions = []string{}
	expect.Plugin.Objects = []string{}

	// case 0
	cfg := &CloudConfig{}
	err := utils.UnmarshalYAML(nil, cfg)
	assert.NoError(t, err)
	assert.EqualValues(t, expect, cfg)

	// case 1
	cfg = &CloudConfig{}
	in := `
nodeServer:
  port: ":9994"

activeServer:
  port: ":9995"

adminServer:
  port: ":9993"
`
	expect.AdminServer.Port = ":9993"
	expect.ActiveServer.Port = ":9995"
	expect.NodeServer.Port = ":9994"
	err = utils.UnmarshalYAML([]byte(in), cfg)
	assert.NoError(t, err)

	assert.EqualValues(t, expect, cfg)
}
