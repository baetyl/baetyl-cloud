package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
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
