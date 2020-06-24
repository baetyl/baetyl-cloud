package pki

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"fmt"
	"github.com/baetyl/baetyl-cloud/common"
	mockPVC "github.com/baetyl/baetyl-cloud/mock/plugin/default"
	"github.com/baetyl/baetyl-cloud/models"
	"github.com/baetyl/baetyl-cloud/plugin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
)

var (
	configYaml = `
defaultpki:
  rootCAFile: "{{.CA_PATH}}/ca.pem"
  rootCAKeyFile: "{{.CA_PATH}}/ca.key"
  persistent:
    kind: "database"
    database:
      type: "sqlite3"
      url: ":memory:"
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
)

func GenPKIConf(t *testing.T) string {
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

func TestNewPKI(t *testing.T) {
	// bad case 0
	_, err := NewPKI()
	assert.Error(t, err)

	// good case
	dir := GenPKIConf(t)
	common.SetConfFile(path.Join(dir, "config.yml"))
	p, err := NewPKI()
	assert.NoError(t, err)
	_, ok := p.(plugin.PKI)
	assert.True(t, ok)

	// bad case 1
	err = os.Remove(path.Join(dir, "ca.pem"))
	_, err = NewPKI()
	assert.Error(t, err)

	// bad case 2
	err = os.Remove(path.Join(dir, "ca.key"))
	_, err = NewPKI()
	assert.Error(t, err)
}

func TestCreateRootCert(t *testing.T) {
	mockCtl := gomock.NewController(t)
	mockPVC := mockPVC.NewMockPVC(mockCtl)

	parentId := "12345678"

	cli := &defaultPkiClient{
		cfg:   CloudConfig{},
		caKey: []byte(caKey),
		caPem: []byte(caPem),
		pvc:   mockPVC,
	}

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

	// good case 0
	mockPVC.EXPECT().CreateCert(gomock.Any()).Return(nil).Times(3)
	_, err := cli.CreateRootCert(csrInfo, "")
	assert.NoError(t, err)

	// good case 1
	cert := &models.Cert{
		Content: []byte(caPem),
	}
	mockPVC.EXPECT().GetCert(parentId).Return(cert, nil).Times(1)
	_, err = cli.CreateRootCert(csrInfo, parentId)
	assert.NoError(t, err)

	// bad case 0
	mockPVC.EXPECT().GetCert(parentId).Return(nil, fmt.Errorf("get cert err")).Times(1)
	_, err = cli.CreateRootCert(csrInfo, parentId)
	assert.Error(t, err)
}

func TestDeleteRootCert(t *testing.T) {
	mockCtl := gomock.NewController(t)
	mockPVC := mockPVC.NewMockPVC(mockCtl)

	parentId := "12345678"

	cli := &defaultPkiClient{
		cfg:   CloudConfig{},
		caKey: []byte(caKey),
		caPem: []byte(caPem),
		pvc:   mockPVC,
	}

	// bad case 0
	mockPVC.EXPECT().CountCertByParentId(parentId).Return(0, fmt.Errorf("count err")).Times(1)
	err := cli.DeleteRootCert(parentId)
	assert.Error(t, err)

	// bad case 1
	mockPVC.EXPECT().CountCertByParentId(parentId).Return(1, nil).Times(1)
	err = cli.DeleteRootCert(parentId)
	assert.Error(t, err, common.Error(common.ErrResourceHasBeenUsed,
		common.Field("type", "certificate"),
		common.Field("name", parentId)))

	// good case 0
	mockPVC.EXPECT().CountCertByParentId(parentId).Return(0, nil).Times(1)
	mockPVC.EXPECT().DeleteCert(parentId).Return(nil).Times(1)
	err = cli.DeleteRootCert(parentId)
	assert.NoError(t, err)
}

func TestCreateSubCert(t *testing.T) {
	mockCtl := gomock.NewController(t)
	mockPVC := mockPVC.NewMockPVC(mockCtl)

	parentId := "12345678"
	base64CSR := "MIIBaDCCAQ8CAQAwgawxCzAJBgNVBAYTAkNOMRAwDgYDVQQIEwdCZWlqaW5nMRkwFwYDVQQHExBIYWlkaWFuIERpc3RyaWN0MRUwEwYDVQQJEwxCYWlkdSBDYW1wdXMxDzANBgNVBBETBjEwMDA5MzEeMBwGA1UEChMVTGludXggRm91bmRhdGlvbiBFZGdlMQ8wDQYDVQQLEwZCQUVUWUwxFzAVBgNVBAMTDmRlZmF1bHQuMDYwMTA4MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEzQrp8J1rTNQj4scxTt8ncJ0Ww2xFw2m8nkxaQTBKLfxyX+TICMhmWyGFxearqHzv5o+aEm3qdgR1N3bt1wvU4KAAMAoGCCqGSM49BAMCA0cAMEQCIHsF8ac5nEEd4b3eDUs2d1jvEcq5O01SZIbgK8hKj6C0AiAe/V6Ya7pnWtnlslb0qrMiDQlh9ltZ4hJLWbG8ZNE45g=="

	cert := &models.Cert{
		Priv:    []byte(caKey),
		Content: []byte(caPem),
	}

	cli := &defaultPkiClient{
		cfg:   CloudConfig{},
		caKey: []byte(caKey),
		caPem: []byte(caPem),
		pvc:   mockPVC,
	}

	csr, err := base64.StdEncoding.DecodeString(base64CSR)
	assert.NoError(t, err)

	// bad case 0
	mockPVC.EXPECT().GetCert(parentId).Return(nil, fmt.Errorf("err")).Times(1)
	_, err = cli.CreateServerCert(csr, parentId)
	assert.Error(t, err)

	// bad case 1
	mockPVC.EXPECT().GetCert(parentId).Return(nil, nil).Times(1)
	_, err = cli.CreateServerCert(csr, parentId)
	assert.Error(t, common.Error(common.ErrResourceNotFound,
		common.Field("type", "certificate"),
		common.Field("name", parentId)))

	// good case
	mockPVC.EXPECT().GetCert(parentId).Return(cert, nil).Times(2)
	mockPVC.EXPECT().CreateCert(gomock.Any()).Return(nil).Times(2)
	_, err = cli.CreateServerCert(csr, parentId)
	assert.NoError(t, err)
	_, err = cli.CreateClientCert(csr, parentId)
	assert.NoError(t, err)
}

func TestGetCert(t *testing.T) {
	mockCtl := gomock.NewController(t)
	mockPVC := mockPVC.NewMockPVC(mockCtl)

	parentId := "12345678"

	cli := &defaultPkiClient{
		cfg:   CloudConfig{},
		caKey: []byte(caKey),
		caPem: []byte(caPem),
		pvc:   mockPVC,
	}

	// bad case
	mockPVC.EXPECT().GetCert(parentId).Return(nil, fmt.Errorf("get cert err")).Times(3)
	_, err := cli.GetRootCert(parentId)
	assert.Error(t, err)
	_, err = cli.GetClientCert(parentId)
	assert.Error(t, err)
	_, err = cli.GetServerCert(parentId)
	assert.Error(t, err)

	// good case
	mockPVC.EXPECT().GetCert(parentId).Return(&models.Cert{}, nil).Times(3)
	_, err = cli.GetRootCert(parentId)
	assert.NoError(t, err)
	_, err = cli.GetClientCert(parentId)
	assert.NoError(t, err)
	_, err = cli.GetServerCert(parentId)
	assert.NoError(t, err)
}

func TestClose(t *testing.T) {
	mockCtl := gomock.NewController(t)
	mockPVC := mockPVC.NewMockPVC(mockCtl)

	cli := &defaultPkiClient{
		cfg:   CloudConfig{},
		caKey: []byte(caKey),
		caPem: []byte(caPem),
		pvc:   mockPVC,
	}

	mockPVC.EXPECT().Close().Return(nil).Times(1)
	err := cli.Close()
	assert.NoError(t, err)
}

func TestDeleteSubCert(t *testing.T) {
	mockCtl := gomock.NewController(t)
	mockPVC := mockPVC.NewMockPVC(mockCtl)

	cli := &defaultPkiClient{
		cfg:   CloudConfig{},
		caKey: []byte(caKey),
		caPem: []byte(caPem),
		pvc:   mockPVC,
	}

	certId := "12345678"
	mockPVC.EXPECT().DeleteCert(certId).Return(nil).Times(2)
	err := cli.DeleteServerCert(certId)
	assert.NoError(t, err)
	err = cli.DeleteClientCert(certId)
	assert.NoError(t, err)
}
