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
	exp.PKI.RootCertId = "98ec3bc552f0478298aa1c6702a95427"
	exp.PKI.SubDuration = 20 * 365 * 24 * time.Hour
	exp.PKI.RootDuration = 50 * 365 * 24 * time.Hour
	exp.PKI.Persistent = Persistent{Kind: "database"}
	exp.PKI.Persistent.Database.Type = "mysql"
	exp.PKI.Persistent.Database.URL = "root:12345678@(127.0.0.1:3306)/baetyl_cloud?charset=utf8&loc=Asia%2FShanghai&parseTime=true"
	exp.PKI.Persistent.Database.MaxConns = 20
	exp.PKI.Persistent.Database.MaxIdleConns = 5
	exp.PKI.Persistent.Database.ConnMaxLifetime = 150

	in := `
defaultpki:
  rootCAFile: "etc/config/cloud/ca.pem"
  rootCAKeyFile: "etc/config/cloud/ca.key"
  rootCertId: "98ec3bc552f0478298aa1c6702a95427"
  persistent:
    kind: "database"
    database:
      type: "mysql"
      url: "root:12345678@(127.0.0.1:3306)/baetyl_cloud?charset=utf8&loc=Asia%2FShanghai&parseTime=true"
`
	cfg := CloudConfig{}
	err := utils.UnmarshalYAML([]byte(in), &cfg)
	assert.NoError(t, err)
	assert.EqualValues(t, exp, cfg)
}
