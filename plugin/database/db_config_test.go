package database

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/common"
)

const (
	confData = `
database:
  type: "mysql"
  url: "xxx"
`
)

func genConf(t *testing.T) string {
	tempDir := t.TempDir()

	err := os.WriteFile(path.Join(tempDir, "config.yml"), []byte(confData), 777)
	assert.NoError(t, err)
	return tempDir
}

func TestDB_Config(t *testing.T) {
	common.SetConfFile(path.Join(genConf(t), "config.yml"))
	var cfg CloudConfig
	err := common.LoadConfig(&cfg)
	assert.NoError(t, err)

	exp := CloudConfig{}

	exp.Database.DecryptionPlugin = "decryption"
	exp.Database.Type = "mysql"
	exp.Database.URL = "xxx"
	exp.Database.MaxConns = 20
	exp.Database.MaxIdleConns = 5
	exp.Database.ConnMaxLifetime = 150

	assert.EqualValues(t, exp, cfg)
}
