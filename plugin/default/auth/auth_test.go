package auth

import (
	"encoding/base64"
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
  keyFile: "etc/baetyl/ca.key"
`
	ca = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA41Ux2wPALzq5agrwhNOr0B+BYePAf5CrVPXtfFKmQDdOqBkB
jM5M61cXyWkcET3d/nKyewvNvWjKjXbaWxxaHAbQsen/5vcte9drW5dVvAquYi45
nQDI2gM00u0FAARydqMzYRFZZir7a3d/tsHiH9HWjkOBkXiBpvDWHXjIRBfCRdgX
3rSwnHDh1Tc0se7rX+0lo8mw7jFCUEH1l9vF/jFE3BGUHK7UHEPg2gDP7w2L8VOF
vTk2rQRGhy2jE3P4xF5x64i8KzOpJw8rbxBezL/yYNcpjuFkkrCpWVLHEgarHPCx
HLB8HkUGAnlM9Apx613icvtRf+sqW1J+F2qfSQIDAQABAoIBAFNrouTcpnxuTzXD
l+kWB5lSxlaWjcAB5W1C5YfWiF1OLlXu/yudVIqTpg3pvTvyePDzQ911QmU7/AAX
Wh9O8x4PvitbU+V8VLt6HFI64WIkhUNP9SJQ9GNUA+FWypvsBdjVIHiBNk4Qfbw8
2KfG0+SbSuFfkj9Aeks5W0jrVontgPSumKy8oI0Y4BsvqjJmwK+c5EHqQ0pZ5x4D
pluqyArQrAbHLMpfpVtaQXiAFdKnLwVTjEl+kZGFgkoJeoNbHGNQabvvH6fjzvPt
8DM44Bs4IvUppIHI6sd42XhcyfYXG/HjynUYtND8I0wWIqabOJJHxE71fQ2i8BfU
Y4qxe0UCgYEA7vUzrbr/a8q0iVrp0rG+98E6c3tuPBNB9cYX16dX20PBhxDznsFJ
6uWdef4x87texQDKYeDScsR9nMiTX7AQDQqMCLHbp65oZGZlXmHJ9LwC7Ppwo3Fc
kt8JGheS+6i6/+SqKAVPyoKMdC5pHQR/APsYmmodq8hLPEQZZfYOuyMCgYEA84u/
aAc6y7m3WdmsNCxwV4eEkqParrmlYCYHtWRqdRL5M3qg6AYSBZJUCzgPz+YtSPON
hO5/1Vq5WzNOMwdbtKT7M7AwCdI0OrInyEq3gKnfsJP3Dbzj54JUcc+JdZqR+m3E
5deBoVBuCp0x0/oVpHcgsEBe4yiUBXAB3UjEKKMCgYEAxdRImY8UATiLaJ/UrvMq
x9C4RH0ukRvcYs5CVO6dBNE+ekSlfIxHVuoMCsBQuJkp5201H/1SHWPhHpjLsc+A
KlvN/TDKSjNRB7XiPFY3LZ8tyOW5tQaX/pwZ2/kiXaieUFYOLR3gpiaYg2Mc8MIV
J0m6X7R0phAngVhbspcYMQMCgYAuRv6u4LjOX1K0swTiwRLzvt91EceK7eG7vF44
nIUSC/HoU0Ph8s1X2682loeCpKU0OHtKqBsISn3wE3angZ1uXO8Sqkbmhte/03x1
taTawOytW+BU7vCLXBt5qMrg2uckI9mHJwUNxv+x6p6+PcYBA1Xlx8V/+oTt55Oj
HaGQawKBgAFjOYIs9tZaTqeHFQkw1zs2Da4e54mMp153XGFcby+2N9XQfs8BQitG
CjpiLI/QVvYXOShRZK/cnMJBmqBhec8534D5GxwTeg/xUNlTq+G/9goLMtWKlBZp
rjBUmKcoSNjVQ8UgIBiSXG0XlCNv+rC23eE5gg0OLLBsPIEOuuSf
-----END RSA PRIVATE KEY-----
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
	if err := ioutil.WriteFile(path.Join(workspace, "ca.key"), []byte(ca), 0755); err != nil {
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

func TestDefaultAuth_Sign_Verify(t *testing.T) {
	err := genConfig("etc/baetyl")
	assert.NoError(t, err)
	defer os.RemoveAll(path.Dir("etc/baetyl"))

	cfg := setUp()
	iam, err := plugin.GetPlugin(cfg.Plugin.Auth)
	assert.NoError(t, err)
	auth := iam.(plugin.Auth)

	meta := []byte("test")
	sign, err := auth.SignToken(meta)
	assert.Nil(t, err)
	res := "o3wtjYemI18SszHjtL3Hzn/qpGcCxCMN4Xw8DeRIw0GTeOFG/b8LJFRPh069ga74Atg+uf91vl1RaT9BAaUNpZh0ISIg5CzKWRIjdkWZUdFMRXQ2+b+qS2gvE4hFiY0+M/JxF65CfP/He6D0f98QVks49a8XkEwv25lIugDS8sd7ayZrxL6acU93HTYU76ee65YOq3DJJkb3JOt+3NfiXtI3dGss3RgrhFpL4rGcat9wxt2EBCjhZwRw5BHYbJcuu8IKmSqQGLMX+RGCbXzdJIPnnApQPRZ5Xgvlad98VyKfuDEznmNNZmctE+iCi4ObL9gY6RZxp6pasoMVAVwyug=="
	assert.Equal(t, res, base64.StdEncoding.EncodeToString(sign))

	inc := auth.VerifyToken(meta, sign)
	assert.Equal(t, true, inc)
}
