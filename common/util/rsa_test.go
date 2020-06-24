package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	privateKey = `-----BEGIN RSA PRIVATE KEY-----
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

func TestGenerateKeyPair(t *testing.T) {
	_, _, err := GenerateKeyPair(1024)
	assert.Nil(t, err)
	_, _, err = GenerateKeyPair(-1)
	assert.NotNil(t, err)
}

func TestRsa(t *testing.T) {
	errText := []byte("error")
	_, err := BytesToPrivateKey(errText)
	assert.Equal(t, ErrRsaPemDecode, err)
	_, err = BytesToPublicKey(errText)
	assert.Equal(t, ErrRsaPemDecode, err)

	privByte := []byte(privateKey)
	priv, err := BytesToPrivateKey(privByte)
	assert.Nil(t, err)

	privByteConvert := PrivateKeyToBytes(priv)
	assert.EqualValues(t, privByte, privByteConvert)

	pubByte, err := PublicKeyToBytes(&priv.PublicKey)
	assert.Nil(t, err)

	pub, err := BytesToPublicKey(pubByte)
	assert.Nil(t, err)
	assert.EqualValues(t, priv.PublicKey, *pub)

	// sign & verify
	text := []byte("to be sign")

	sign, err := SignPKCS1v15(text, priv)
	assert.Nil(t, err)

	res := VerifyPKCS1v15(text, sign, pub)
	assert.Equal(t, true, res)

	res = VerifyPKCS1v15(errText, sign, pub)
	assert.Equal(t, false, res)
}
