package sign

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

const (
	confData = `
defaultsign:
`
)

func setUp() *config.CloudConfig {
	conf := &config.CloudConfig{}
	conf.Plugin.Sign = "defaultsign"
	return conf
}

func genConfig(workspace string) error {
	if err := os.MkdirAll(workspace, 0755); err != nil {
		return err
	}
	if err := ioutil.WriteFile(path.Join(workspace, "cloud.yml"), []byte(confData), 0755); err != nil {
		return err
	}
	return nil
}

func TestDefaultSign_Sign_Verify(t *testing.T) {
	err := genConfig("etc/baetyl")
	assert.NoError(t, err)
	defer os.RemoveAll(path.Dir("etc/baetyl"))

	cfg := setUp()
	iam, err := plugin.GetPlugin(cfg.Plugin.Sign)
	assert.NoError(t, err)
	auth := iam.(plugin.Sign)

	meta := []byte("test")
	sign, err := auth.Signature(meta)
	assert.Nil(t, err)
	res := "dGVzdA=="
	assert.Equal(t, res, base64.StdEncoding.EncodeToString(sign))

	inc := auth.Verify(meta, sign)
	assert.Equal(t, true, inc)
}
