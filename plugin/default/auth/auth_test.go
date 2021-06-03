package auth

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

const (
	confData = `
defaultauth:
  namespace: testns
`
)

func setUp() *config.CloudConfig {
	conf := &config.CloudConfig{}
	conf.Plugin.Auth = "defaultauth"
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

func TestDefaultAuth_Authenticate(t *testing.T) {
	err := genConfig("etc/baetyl")
	assert.NoError(t, err)
	defer os.RemoveAll(path.Dir("etc/baetyl"))

	iam, err := plugin.GetPlugin("defaultauth")
	assert.NoError(t, err)
	auth := iam.(plugin.Auth)

	ctx := common.NewContextEmpty()
	err = auth.Authenticate(ctx)

	assert.NoError(t, err)
	assert.Equal(t, "testns", ctx.GetNamespace())
}
