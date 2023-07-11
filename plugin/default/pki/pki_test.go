package pki

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/baetyl/baetyl-go/v2/pki"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/common"
	mockPKI "github.com/baetyl/baetyl-cloud/v2/mock/plugin"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

var (
	configYaml = `
defaultpki:
  rootCAFile: "{{.CA_PATH}}/ca.pem"
  rootCAKeyFile: "{{.CA_PATH}}/ca.key"
  persistent: "database"
`
	caPem = `-----BEGIN CERTIFICATE-----
MIID7DCCAtSgAwIBAgIDAYagMA0GCSqGSIb3DQEBCwUAMIGlMQswCQYDVQQGEwJD
TjEQMA4GA1UECBMHQmVpamluZzEZMBcGA1UEBxMQSGFpZGlhbiBEaXN0cmljdDEV
MBMGA1UECRMMQmFpZHUgQ2FtcHVzMQ8wDQYDVQQREwYxMDAwOTMxHjAcBgNVBAoT
FUxpbnV4IEZvdW5kYXRpb24gRWRnZTEPMA0GA1UECxMGQkFFVFlMMRAwDgYDVQQD
Ewdyb290LmNhMB4XDTIwMDMyNjAzMzE1MVoXDTMwMDMyNjAzMzE1MVowgaUxCzAJ
BgNVBAYTAkNOMRAwDgYDVQQIEwdCZWlqaW5nMRkwFwYDVQQHExBIYWlkaWFuIERp
c3RyaWN0MRUwEwYDVQQJEwxCYWlkdSBDYW1wdXMxDzANBgNVBBETBjEwMDA5MzEe
MBwGA1UEChMVTGludXggRm91bmRhdGlvbiBFZGdlMQ8wDQYDVQQLEwZCQUVUWUwx
EDAOBgNVBAMTB3Jvb3QuY2EwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIB
AQCye0EWM/owq1OXEKZdOOy6hLjXk4LlOeLIHoHWkidA2C+OvJhBg4eu0laHwlcb
0dlb4O0tZ0pDlFlNit8vfBzciFOTIQDXcRlSE7rs1USilX5YvRyoSmBAw34nuyq4
GobdQtAmMlwLds/h1MIskH6WeMApnFL2TqDHBUdPHhBdSS7fi9uC+zH+otjK7R7y
v89pPWc9mwaDQTreZcgswCKm7bZT4C73m0lgBSEOLHkQ4wa6nlQEOZMadovioBYJ
ihswoVB86++kkJ/6C2WeMebMb+ha3ExRORY15rUjWm6/M7otpoL51bcnyAhKl4Ee
UDJEjCkmrhHtYK1djaQJ1J53AgMBAAGjIzAhMA4GA1UdDwEB/wQEAwIBhjAPBgNV
HRMBAf8EBTADAQH/MA0GCSqGSIb3DQEBCwUAA4IBAQBpF9V+LSxAhxYAUsaltSJu
VFk+CVVOkDSh1i5BSkjvca9SnPk8ukjPWsq7Ru74HHiZ4ZsjfrtVtloyoijUXPji
piZOm30+kHtlaVi10T0r0E633x6345/yHYXTawVXgUOsMG9HPu2LnW2sy9DDmMYA
DHG83CZle1WWBFYE6FmUwugQ2IKUo0MYV/xIhulcYQPUNJlqnyvvJAWi4xL61jD6
MH5XrLyAGLEIkSgmrgcD/B0LkBviLNhAqNmP0GbzcrdtjmKZF1ERUpVt1ko7lgvZ
3EOyHbdBJOerlXUHI+/uEWUDPiuu59PoREZ9tuuJMO7UQlkP0NmPbqkwdvUKkaNp
-----END CERTIFICATE-----
`
	caKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAsntBFjP6MKtTlxCmXTjsuoS415OC5TniyB6B1pInQNgvjryY
QYOHrtJWh8JXG9HZW+DtLWdKQ5RZTYrfL3wc3IhTkyEA13EZUhO67NVEopV+WL0c
qEpgQMN+J7squBqG3ULQJjJcC3bP4dTCLJB+lnjAKZxS9k6gxwVHTx4QXUku34vb
gvsx/qLYyu0e8r/PaT1nPZsGg0E63mXILMAipu22U+Au95tJYAUhDix5EOMGup5U
BDmTGnaL4qAWCYobMKFQfOvvpJCf+gtlnjHmzG/oWtxMUTkWNea1I1puvzO6LaaC
+dW3J8gISpeBHlAyRIwpJq4R7WCtXY2kCdSedwIDAQABAoIBAGbyMsuEtXVnDLLg
lqTElb7LmPY3DlP7PHRjLE7AREXhrCSvYT7Ah/1tMx3hGW9hbfbR2NvMbQhnw863
IB56fwcw1svRSHP7tzghSzsZlBoXEiZLBgGHzNbuK5DtIynHmyx6QicV+wNdx3Ah
0NH1kh5mjagyk6OgHJpO0B+xXoz/FxcLEmjmdp9/H8ByGLgVXlolHPkt9CBZqCe9
SFGDQv4FIWJ8Y2KB/avOpKeASR8pIiOtKo+rttk0aWM8kc1rtPa3xabasUoeUNT7
CA7cmxV4iLdE59olKLzPM1AMLxWDNYWb0qr/lor7BmuhsObxqlnJo459j3gMiTZi
dGFIlQECgYEAx+SkqF1M9kxZwAxI4NptUn8tS7mKTiCRBuuIODgEUK1en4uM4mh3
FxbJzpnxTulQT7pKKiUMXTUVTyBoya9GI5lyfxm7tQ+w1JYFosR+8rpFF1bgod2n
eanIiifLzG4mexI0yDOQnkn8lSWLuKjZiFbZv/A9OGnCMmfetynkS4ECgYEA5JQS
QUtzW6REwQFuXGgZctjg4Gl21L7xL8K+yd/IcnhIEmPP9TJz+tJ9oGp1hjGpgGrz
LPz3U5fXiDDlaEErr5IqmKUS511dQHfmcF8FXbD8uspJGzk7aFi64NZ2PUMlbiDp
+FRG1X3at/ecH9gqg5LoTXxnnM2UU6yQoKk2RfcCgYBDKOPVmXtZKS/iYX4+5cRj
Ok16qrz4IOL5IztiQBfbD1TCX/2WuCiC/moRWxGDRMpx7xIp9MahrksZibcLRDNZ
lJ2ubHPvknUEB9+e30wTu1epTswsNi+lpdC18kb7yWpuYSCQvxpwxETzy2iVQ03L
C/sfDNVU1dukWdevTIjigQKBgQCUQUHx3cktmEcL1CzLfK184xRAGcd8R3hR3QM4
FpCBRmignOKGC7pT5fCbelFNv6pL45JkDJMyQdsGt4gj7ZkzIB/Gr9KqA9F2/g2V
ttvZH/FcCdYO9TkF/f7/07oPFB0T5/85FRh4Yk/ZYJ1/vgodGszXbSga+PAKsXOA
8R+FkwKBgBZsgv4DyFLBxgI6qPZbQ5ancjOFl2p/oCfwqdLq0iG8UaFACUBA9PTQ
ITHWRBk8fdWmDHREbrYeym3sTdIKP5HN24WkVm9A3CZ6ZJPeFfkn83S87baZ6Rmv
w3xQdGBSx9ae6exKX6qVqsjQDv5X443H8yHcU0EQ8DUnth+jwK7H
-----END RSA PRIVATE KEY-----
`
	base64CSR = "MIIBaDCCAQ8CAQAwgawxCzAJBgNVBAYTAkNOMRAwDgYDVQQIEwdC" +
		"ZWlqaW5nMRkwFwYDVQQHExBIYWlkaWFuIERpc3RyaWN0MRUwEwYDVQQJEwxCY" +
		"WlkdSBDYW1wdXMxDzANBgNVBBETBjEwMDA5MzEeMBwGA1UEChMVTGludXggRm" +
		"91bmRhdGlvbiBFZGdlMQ8wDQYDVQQLEwZCQUVUWUwxFzAVBgNVBAMTDmRlZmF" +
		"1bHQuMDYwMTA4MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEzQrp8J1rTNQj" +
		"4scxTt8ncJ0Ww2xFw2m8nkxaQTBKLfxyX+TICMhmWyGFxearqHzv5o+aEm3qd" +
		"gR1N3bt1wvU4KAAMAoGCCqGSM49BAMCA0cAMEQCIHsF8ac5nEEd4b3eDUs2d1" +
		"jvEcq5O01SZIbgK8hKj6C0AiAe/V6Ya7pnWtnlslb0qrMiDQlh9ltZ4hJLWbG" +
		"8ZNE45g=="
)

func genPKIConf(t *testing.T) string {
	tempDir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)

	configYaml = strings.ReplaceAll(configYaml, "{{.CA_PATH}}", tempDir)
	err = ioutil.WriteFile(path.Join(tempDir, "config.yml"), []byte(configYaml), 777)
	assert.NoError(t, err)
	err = ioutil.WriteFile(path.Join(tempDir, "ca.pem"), []byte(caPem), 777)
	assert.NoError(t, err)
	err = ioutil.WriteFile(path.Join(tempDir, "ca.key"), []byte(caKey), 777)
	assert.NoError(t, err)
	return tempDir
}

func genDefaultPkiClient(t *testing.T) (*defaultPkiClient, *mockPKI.MockPKIStorage) {
	mockCtl := gomock.NewController(t)
	mockSto := mockPKI.NewMockPKIStorage(mockCtl)

	cfg := &CloudConfig{}
	common.SetConfFile(path.Join(genPKIConf(t), "config.yml"))
	err := common.LoadConfig(cfg)
	assert.NoError(t, err)

	pkiClient, err := pki.NewPKIClient()
	assert.NoError(t, err)

	return &defaultPkiClient{
		cfg:       *cfg,
		sto:       mockSto,
		pkiClient: pkiClient,
	}, mockSto
}

func genRootCAView() *plugin.Cert {
	return &plugin.Cert{
		CertId:     "98ec3bc552f0478298aa1c6702a95427",
		Type:       TypeIssuingCA,
		CommonName: "root.ca",
		Content:    base64.StdEncoding.EncodeToString([]byte(caPem)),
		PrivateKey: base64.StdEncoding.EncodeToString([]byte(caKey)),
	}
}

func TestNewPKI(t *testing.T) {
	// bad case 0
	_, err := NewPKI()
	assert.Error(t, err)

	// good case
	dir := genPKIConf(t)
	common.SetConfFile(path.Join(dir, "config.yml"))
	_, err = NewPKI()
	assert.Error(t, err)
}

func TestDefaultPkiClient_CheckRootCertId(t *testing.T) {
	p, s := genDefaultPkiClient(t)

	// good case 0
	s.EXPECT().GetCert(RootCertId).Return(nil, nil).Times(1)
	s.EXPECT().DeleteCert(RootCertId).Return(nil).Times(1)
	s.EXPECT().CreateCert(gomock.Any()).Return(nil).Times(2)
	err := p.checkRootCA()
	assert.NoError(t, err)

	// good case 1
	s.EXPECT().GetCert(RootCertId).Return(nil, os.ErrNotExist).Times(1)
	err = p.checkRootCA()
	assert.NoError(t, err)
}

func TestDefaultPkiClient_GetRootCertId(t *testing.T) {
	p, _ := genDefaultPkiClient(t)
	res := p.GetRootCertID()
	assert.Equal(t, RootCertId, res)
}

func TestDefaultPkiClient_CreateRootCert(t *testing.T) {
	p, s := genDefaultPkiClient(t)

	csrInfo := &x509.CertificateRequest{
		Subject: pkix.Name{
			Country:            []string{"CN"},
			Organization:       []string{"Linux Foundation Edge"},
			OrganizationalUnit: []string{"BAETYL"},
			Locality:           []string{"Haidian District"},
			Province:           []string{"Beijing"},
			StreetAddress:      []string{"Baidu Campus"},
			PostalCode:         []string{"100093"},
			CommonName:         "test",
		},
		EmailAddresses: []string{"baetyl@lists.lfedge.org"},
	}

	s.EXPECT().GetCert(RootCertId).Return(genRootCAView(), nil).Times(1)
	s.EXPECT().CreateCert(gomock.Any()).Return(nil).Times(1)

	res, err := p.CreateRootCert(csrInfo, "")
	assert.NoError(t, err)
	assert.NotEqual(t, "", res)
}

func TestDefaultPkiClient_GetRootCert(t *testing.T) {
	p, s := genDefaultPkiClient(t)
	s.EXPECT().GetCert(RootCertId).Return(genRootCAView(), nil).Times(1)
	res, err := p.GetRootCert(RootCertId)
	assert.NoError(t, err)
	assert.Equal(t, caPem, string(res))
}

func TestDefaultPkiClient_DeleteRootCert(t *testing.T) {
	p, s := genDefaultPkiClient(t)
	// good case
	s.EXPECT().CountCertByParentId(RootCertId).Return(0, nil).Times(1)
	s.EXPECT().DeleteCert(RootCertId).Return(nil).Times(1)
	err := p.DeleteRootCert(RootCertId)
	assert.NoError(t, err)

	// bad case
	s.EXPECT().CountCertByParentId(RootCertId).Return(1, nil).Times(1)
	err = p.DeleteRootCert(RootCertId)
	assert.Error(t, err, ErrCertInUsed.Error())
}

func TestDefaultPkiClient_CreateServerCert(t *testing.T) {
	p, s := genDefaultPkiClient(t)
	s.EXPECT().GetCert(RootCertId).Return(genRootCAView(), nil).Times(1)
	s.EXPECT().CreateCert(gomock.Any()).Return(nil).Times(1)

	csr, err := base64.StdEncoding.DecodeString(base64CSR)
	assert.NoError(t, err)
	res, err := p.CreateServerCert(csr, RootCertId)
	assert.NoError(t, err)
	assert.NotEqual(t, "", res)
}

func TestDefaultPkiClient_GetServerCert(t *testing.T) {
	p, s := genDefaultPkiClient(t)
	s.EXPECT().GetCert(RootCertId).Return(genRootCAView(), nil).Times(1)
	res, err := p.GetServerCert(RootCertId)
	assert.NoError(t, err)
	assert.Equal(t, caPem, string(res))
}

func TestDefaultPkiClient_DeleteServerCert(t *testing.T) {
	p, s := genDefaultPkiClient(t)
	s.EXPECT().DeleteCert(RootCertId).Return(nil).Times(1)
	err := p.DeleteServerCert(RootCertId)
	assert.NoError(t, err)
}

func TestDefaultPkiClient_CreateClientCert(t *testing.T) {
	p, s := genDefaultPkiClient(t)
	s.EXPECT().GetCert(RootCertId).Return(genRootCAView(), nil).Times(1)
	s.EXPECT().CreateCert(gomock.Any()).Return(nil).Times(1)

	csr, err := base64.StdEncoding.DecodeString(base64CSR)
	assert.NoError(t, err)
	res, err := p.CreateClientCert(csr, RootCertId)
	assert.NoError(t, err)
	assert.NotEqual(t, "", res)
}

func TestDefaultPkiClient_GetClientCert(t *testing.T) {
	p, s := genDefaultPkiClient(t)
	s.EXPECT().GetCert(RootCertId).Return(genRootCAView(), nil).Times(1)
	res, err := p.GetClientCert(RootCertId)
	assert.NoError(t, err)
	assert.Equal(t, caPem, string(res))
}

func TestDefaultPkiClient_DeleteClientCert(t *testing.T) {
	p, s := genDefaultPkiClient(t)
	s.EXPECT().DeleteCert(RootCertId).Return(nil).Times(1)
	err := p.DeleteClientCert(RootCertId)
	assert.NoError(t, err)
}

func TestDefaultPkiClient_Close(t *testing.T) {
	p, s := genDefaultPkiClient(t)
	s.EXPECT().Close().Return(nil).Times(1)
	err := p.Close()
	assert.NoError(t, err)
}
