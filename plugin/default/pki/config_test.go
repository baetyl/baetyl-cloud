package pki

import (
	"testing"
	"time"

	"github.com/baetyl/baetyl-go/v2/utils"
	"github.com/stretchr/testify/assert"
)

func TestPKI_Config(t *testing.T) {
	exp := CloudConfig{}
	exp.PKI.RootCAFile = "etc/config/cloud/ca.pem"
	exp.PKI.RootCAKeyFile = "etc/config/cloud/ca.key"
	exp.PKI.SubDuration = 20 * 365 * 24 * time.Hour
	exp.PKI.RootDuration = 50 * 365 * 24 * time.Hour
	exp.PKI.Persistent = "database"

	in := `
defaultpki:
  rootCAFile: "etc/config/cloud/ca.pem"
  rootCAKeyFile: "etc/config/cloud/ca.key"
  persistent: "database"
`
	cfg := CloudConfig{}
	err := utils.UnmarshalYAML([]byte(in), &cfg)
	assert.NoError(t, err)
	assert.EqualValues(t, exp, cfg)
}
