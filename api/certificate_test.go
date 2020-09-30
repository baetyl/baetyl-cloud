package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/common"
	ms "github.com/baetyl/baetyl-cloud/v2/mock/service"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/service"
)

func initCertificateAPI(t *testing.T) (*API, *gin.Engine, *gomock.Controller) {
	api := &API{}
	router := gin.Default()
	mockCtl := gomock.NewController(t)
	mockIM := func(c *gin.Context) { common.NewContext(c).SetNamespace("default") }
	v1 := router.Group("v1")
	{
		certificate := v1.Group("/certificates")
		certificate.GET("/:name", mockIM, common.Wrapper(api.GetCertificate))
		certificate.PUT("/:name", mockIM, common.Wrapper(api.UpdateCertificate))
		certificate.DELETE("/:name", mockIM, common.Wrapper(api.DeleteCertificate))
		certificate.POST("", mockIM, common.Wrapper(api.CreateCertificate))
		certificate.GET("", mockIM, common.Wrapper(api.ListCertificate))
		certificate.GET("/:name/apps", mockIM, common.Wrapper(api.GetAppByCertificate))
	}
	return api, router, mockCtl
}

func TestCreateCertificate(t *testing.T) {
	api, router, mockCtl := initCertificateAPI(t)
	defer mockCtl.Finish()
	mkSecretService := ms.NewMockSecretService(mockCtl)
	api.AppCombinedService = &service.AppCombinedService{
		Secret: mkSecretService,
	}

	ns := "ns"
	keyData1 := `-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEApd4vdBUOtuW+v4YXkqJUFyz18N3tk00Ff3X5SfT036qEC4SB
FsZDyNAX5bhXKkKMX0oAWt7XowpCsUtppChkibfAQNtQbOBOTVC9UJEtef/edEVz
PxyEb8jBq359pee3zZoV55EwxN7uWq1oS/zZ+i6rSncfqKNzI/9kLoyxda6CgnAB
rIpEwdMtqWs8b1Y/0q+s/94owc4jzypxTWSwXLr13s7p2ZDNaJ7xl3AL8zzJw9NK
Y849CImAGu4FO8J+a5zlTEioKMKz3O6jUvZZXpRLQdA6AOZKeiozwNXqKJcNwx8g
7kZTyTxqinq/OeiGiePr9cwElqe8KkwP4pIVpQIDAQABAoIBAQCFGvgZv4w/Wb7p
E0J3eazhrELxOCcevgBbeODEaL7ZfozYcUzmadSboeKLhpLsZtse3NPMGGgTfnhm
ro3oHkIQAlVVtqmjtZ0gjlpd/SLxdFOgGtuRGeFtkz1X0foi2QC3DZ/mZK0uT3gX
bHD2CcMi8bCj4VSWkBQmHxzV/jGqrUVkkEbfXjeujkXmcv9X/8Sj1UvnIIECWK4r
HT6lRoy75tWa+/wfcNq+oFEfw20VySSp+fh7O0swXohdTdtZxYXaEAXjE8ptdodW
8/sKAsPbjfnmqstkfq3iWl3vImkNcyuOxBtar0mL57kHfEOSmTQJEsNhZQNpSDFI
YNdKX79BAoGBANIGkUugCxX19uzcRJcGJuhYiDdQ+Ew7njN5nqjhMRnV023iFLiv
3xaUxBUMgDQX/iHYENC/Z1uJZr4f/gCMy4Ze9Nqm++p9QEyw3zyKDfXu486lKdyx
XSLiBZQ6kGksGBHLMqwwzZ9bL942chVjpnO3kn/RPt6X57t8iF1JREKTAoGBAMot
GmQyKs05pVMAvmk7ZPD3V0fn2PigOxLRUuGRA8IhGDNR7FBWkbuUt1kniV11RCJ+
L3+xFA937FpqqzSzKodu6mXUwtKxMK5KiPdsKQqUmMiuOnyzaPpC+iz9qe65OVA6
W+9BWk9cmN5XqyaIUDdCb5SLKOPhkMEfskDp4NHnAoGAbPm1aC0Js4Jldi8wc8Bg
bcyKGVGtFDkW9BSV64C1LneRdgGJyO6QbbIRL+7FksIkPcFTsEywP4HCysHk1Lo5
XGZm3BEqw1fsBh78JfhoGAS1NWLjnrx03AW06V2d0sRrVMg/abME7jutUbqkZU7I
bmCA5ktXOL5PIiwSwXyjq3sCgYEApEZfslg9BQI4/ieVkCXtkAo5xjhxyRtQxKqH
ILdXCW8gndqMHH8q7PMaw3tnlyPImApWB/hXZ3Y2+wS/VhPak68hEFr/bnkBKC1x
+zDMbEdvmWhQJ7ETtH2lj9cRM+MW2cSBnPdKLT/9CnTLoYSTQUNfLKCiOf+3QeTC
TxJ6VbMCgYEAi9Yp/M4aaeuA6UieSNLZcQpQLvVLDtHg8uyeD7yIH/l6KfT1X5mO
qyDPy5ZoeaNFLrcquIrjd0Z5HcZQF0G3/gVlgql7AGMiWgK2wPmCV6HgGiZEdua2
OqPxJuZbP3aANy19Vn5k+W8wE81Neh8RjO2CC8f0wkSxHarBKDJxR38=
-----END RSA PRIVATE KEY-----`
	certData1 := `-----BEGIN CERTIFICATE-----
MIICvjCCAaYCCQCrqj3SrB8ZxjANBgkqhkiG9w0BAQsFADAhMQswCQYDVQQGEwJD
TjESMBAGA1UEAwwJY2F0dGxlLWNhMB4XDTIwMDkwMjA4MTE1NloXDTMwMDgzMTA4
MTE1NlowITELMAkGA1UEBhMCQ04xEjAQBgNVBAMMCWNhdHRsZS1jYTCCASIwDQYJ
KoZIhvcNAQEBBQADggEPADCCAQoCggEBAKXeL3QVDrblvr+GF5KiVBcs9fDd7ZNN
BX91+Un09N+qhAuEgRbGQ8jQF+W4VypCjF9KAFre16MKQrFLaaQoZIm3wEDbUGzg
Tk1QvVCRLXn/3nRFcz8chG/Iwat+faXnt82aFeeRMMTe7lqtaEv82fouq0p3H6ij
cyP/ZC6MsXWugoJwAayKRMHTLalrPG9WP9KvrP/eKMHOI88qcU1ksFy69d7O6dmQ
zWie8ZdwC/M8ycPTSmPOPQiJgBruBTvCfmuc5UxIqCjCs9zuo1L2WV6US0HQOgDm
SnoqM8DV6iiXDcMfIO5GU8k8aop6vznohonj6/XMBJanvCpMD+KSFaUCAwEAATAN
BgkqhkiG9w0BAQsFAAOCAQEAPUUFgOOenmB8eexetT2o8HKFb13jeQFEiasuOyM4
gb8u1L0e3OUHLQi4CJml7szw10sdc4/x6wqtxJd8K11m1xJ8NxqmsmGGedzGtIYu
lZ6xthcgpvwSoDqVQr4FYYYjCmurNsaYFu7GVSqPCgL0SbiqNmfLK5sKf9BVkMlt
cDNsVm/zxNtmRFyTM8uGX3RoA/kltxnElAOqYJEAcePMQ4IoThfif/ql3Y1D6J5u
O0o7IM5tlD8UhYfAYVH/xyKxdv6K23zHgbuF4QWrBEnYey0eqznO+5bNV7bz0pkc
tu5nww5RdjCz4Uks08P2GNmZjLO81MgYkhR7B9wi3KDNxg==
-----END CERTIFICATE-----`

	cert1 := &models.Certificate{
		Name:        "cert1",
		Description: "desp1",
		Data: models.CertificateDataItem{
			Key:  keyData1,
			Cert: certData1,
		},
	}

	err := cert1.ParseCertInfo()
	assert.NoError(t, err)

	temp1 := cert1.ToSecret()

	res1 := &specV1.Secret{
		Name:      cert1.Name,
		Namespace: ns,
		Labels: map[string]string{
			specV1.SecretLabel: specV1.SecretCustomCertificate,
		},
		Data: map[string][]byte{
			"key":                []byte(keyData1),
			"cert":               []byte(certData1),
			"signatureAlgorithm": []byte(`SHA256-RSA`),
			"effectiveTime":      []byte("2020-09-02 08:11:56 +0000 UTC"),
			"expiredTime":        []byte("2030-08-31 08:11:56 +0000 UTC"),
			"serialNumber":       []byte("12369767301566634438"),
			"issuer":             []byte("cattle-ca"),
			"fingerPrint":        []byte("5F:D4:2B:70:C3:8F:54:FB:3B:AA:28:C8:4E:47:F3:38:FA:74:83:94"),
		},
		CreationTimestamp: time.Now(),
		UpdateTimestamp:   time.Now(),
		Description:       cert1.Description,
		Version:           "1234",
	}
	mkSecretService.EXPECT().Get(gomock.Any(), cert1.Name, gomock.Any()).Return(nil, nil).Times(1)
	mkSecretService.EXPECT().Create(gomock.Any(), temp1).Return(res1, nil).Times(1)
	w := httptest.NewRecorder()
	body, _ := json.Marshal(cert1)
	req, _ := http.NewRequest(http.MethodPost, "/v1/certificates", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	keyData2 := `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIPcuJ/l+/s8PXvAN5M6VNBZrKD4HDW6n6Y4rQCYinF5doAoGCCqGSM49
AwEHoUQDQgAEQXCQTGn4+frJYOumFk8gs8BIbgduEuiHonhYdJTFGIPLiOqPQoIv
DmICod7W0oIzYYXwTF4NfadliSryoXx9IQ==
-----END EC PRIVATE KEY-----`
	certData2 := `-----BEGIN CERTIFICATE-----
MIICojCCAkigAwIBAgIIFizowlwTTkAwCgYIKoZIzj0EAwIwgaUxCzAJBgNVBAYT
AkNOMRAwDgYDVQQIEwdCZWlqaW5nMRkwFwYDVQQHExBIYWlkaWFuIERpc3RyaWN0
MRUwEwYDVQQJEwxCYWlkdSBDYW1wdXMxDzANBgNVBBETBjEwMDA5MzEeMBwGA1UE
ChMVTGludXggRm91bmRhdGlvbiBFZGdlMQ8wDQYDVQQLEwZCQUVUWUwxEDAOBgNV
BAMTB3Jvb3QuY2EwHhcNMjAwODIwMDcxODA5WhcNNDAwODE1MDcxODA5WjCBpDEL
MAkGA1UEBhMCQ04xEDAOBgNVBAgTB0JlaWppbmcxGTAXBgNVBAcTEEhhaWRpYW4g
RGlzdHJpY3QxFTATBgNVBAkTDEJhaWR1IENhbXB1czEPMA0GA1UEERMGMTAwMDkz
MR4wHAYDVQQKExVMaW51eCBGb3VuZGF0aW9uIEVkZ2UxDzANBgNVBAsTBkJBRVRZ
TDEPMA0GA1UEAxMGc2VydmVyMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEQXCQ
TGn4+frJYOumFk8gs8BIbgduEuiHonhYdJTFGIPLiOqPQoIvDmICod7W0oIzYYXw
TF4NfadliSryoXx9IaNhMF8wDgYDVR0PAQH/BAQDAgWgMB0GA1UdJQQWMBQGCCsG
AQUFBwMCBggrBgEFBQcDATAMBgNVHRMBAf8EAjAAMCAGA1UdEQQZMBeCCWxvY2Fs
aG9zdIcEAAAAAIcEfwAAATAKBggqhkjOPQQDAgNIADBFAiB5vz8+oob7SkN54uf7
RErbE4tWT5AHtgBgIs3A+TjnyQIhAPvnL8W1dq4qdkVr0eiH5He0xNHdsQc6eWxS
RcKyjhh1
-----END CERTIFICATE-----`

	cert2 := &models.Certificate{
		Name:        "cert2",
		Description: "desp2",
		Data: models.CertificateDataItem{
			Key:  keyData2,
			Cert: certData2,
		},
	}

	err = cert2.ParseCertInfo()
	assert.NoError(t, err)

	temp2 := cert2.ToSecret()

	res2 := &specV1.Secret{
		Name:      cert2.Name,
		Namespace: ns,
		Labels: map[string]string{
			specV1.SecretLabel: specV1.SecretCustomCertificate,
		},
		Data: map[string][]byte{
			"key":                []byte(keyData2),
			"cert":               []byte(certData2),
			"signatureAlgorithm": []byte(`ECDSA-SHA256`),
			"effectiveTime":      []byte("2020-08-20 07:18:09 +0000 UTC"),
			"expiredTime":        []byte("2040-08-15 07:18:09 +0000 UTC"),
			"serialNumber":       []byte("1597907889275752000"),
			"issuer":             []byte("root.ca"),
			"fingerPrint":        []byte("8A:80:76:69:91:B0:EB:A6:45:E2:FE:88:03:F0:66:4A:A9:BB:16:96"),
		},
		CreationTimestamp: time.Now(),
		UpdateTimestamp:   time.Now(),
		Description:       cert2.Description,
		Version:           "1234",
	}
	mkSecretService.EXPECT().Get(gomock.Any(), cert2.Name, gomock.Any()).Return(nil, nil).Times(1)
	mkSecretService.EXPECT().Create(gomock.Any(), temp2).Return(res2, nil).Times(1)
	w = httptest.NewRecorder()
	body, _ = json.Marshal(cert2)
	req, _ = http.NewRequest(http.MethodPost, "/v1/certificates", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	keyData3 := ``
	certData3 := ``

	cert3 := &models.Certificate{
		Name:        "cert3",
		Description: "desp3",
		Data: models.CertificateDataItem{
			Key:  keyData3,
			Cert: certData3,
		},
	}

	err = cert3.ParseCertInfo()
	assert.Equal(t, err.Error(), "tls: failed to find any PEM data in certificate input")

	keyData4 := `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIPcuJ/l+/s8PXvAN5M6VNBZrKD4HDW6n6Y4rQCYinF5doAoGCCqGSM49
AwEHoUQDQgAEQXCQTGn4+frJYOumFk8gs8BIbgduEuiHonhYdJTFGIPLiOqPQoIv
DmICod7W0oIzYYXwTF4NfadliSryoXx9IQ==
-----END EC PRIVATE KEY-----`
	certData4 := ``

	cert4 := &models.Certificate{
		Name:        "cert4",
		Description: "desp4",
		Data: models.CertificateDataItem{
			Key:  keyData4,
			Cert: certData4,
		},
	}

	err = cert4.ParseCertInfo()
	assert.Equal(t, err.Error(), "tls: failed to find any PEM data in certificate input")

	keyData5 := ``
	certData5 := `-----BEGIN CERTIFICATE-----
MIICojCCAkigAwIBAgIIFizowlwTTkAwCgYIKoZIzj0EAwIwgaUxCzAJBgNVBAYT
AkNOMRAwDgYDVQQIEwdCZWlqaW5nMRkwFwYDVQQHExBIYWlkaWFuIERpc3RyaWN0
MRUwEwYDVQQJEwxCYWlkdSBDYW1wdXMxDzANBgNVBBETBjEwMDA5MzEeMBwGA1UE
ChMVTGludXggRm91bmRhdGlvbiBFZGdlMQ8wDQYDVQQLEwZCQUVUWUwxEDAOBgNV
BAMTB3Jvb3QuY2EwHhcNMjAwODIwMDcxODA5WhcNNDAwODE1MDcxODA5WjCBpDEL
MAkGA1UEBhMCQ04xEDAOBgNVBAgTB0JlaWppbmcxGTAXBgNVBAcTEEhhaWRpYW4g
RGlzdHJpY3QxFTATBgNVBAkTDEJhaWR1IENhbXB1czEPMA0GA1UEERMGMTAwMDkz
MR4wHAYDVQQKExVMaW51eCBGb3VuZGF0aW9uIEVkZ2UxDzANBgNVBAsTBkJBRVRZ
TDEPMA0GA1UEAxMGc2VydmVyMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEQXCQ
TGn4+frJYOumFk8gs8BIbgduEuiHonhYdJTFGIPLiOqPQoIvDmICod7W0oIzYYXw
TF4NfadliSryoXx9IaNhMF8wDgYDVR0PAQH/BAQDAgWgMB0GA1UdJQQWMBQGCCsG
AQUFBwMCBggrBgEFBQcDATAMBgNVHRMBAf8EAjAAMCAGA1UdEQQZMBeCCWxvY2Fs
aG9zdIcEAAAAAIcEfwAAATAKBggqhkjOPQQDAgNIADBFAiB5vz8+oob7SkN54uf7
RErbE4tWT5AHtgBgIs3A+TjnyQIhAPvnL8W1dq4qdkVr0eiH5He0xNHdsQc6eWxS
RcKyjhh1
-----END CERTIFICATE-----`

	cert5 := &models.Certificate{
		Name:        "cert5",
		Description: "desp5",
		Data: models.CertificateDataItem{
			Key:  keyData5,
			Cert: certData5,
		},
	}

	err = cert5.ParseCertInfo()
	assert.Equal(t, err.Error(), "tls: failed to find any PEM data in key input")

	keyData6 := `-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEApd4vdBUOtuW+v4YXkqJUFyz18N3tk00Ff3X5SfT036qEC4SB
FsZDyNAX5bhXKkKMX0oAWt7XowpCsUtppChkibfAQNtQbOBOTVC9UJEtef/edEVz
PxyEb8jBq359pee3zZoV55EwxN7uWq1oS/zZ+i6rSncfqKNzI/9kLoyxda6CgnAB
rIpEwdMtqWs8b1Y/0q+s/94owc4jzypxTWSwXLr13s7p2ZDNaJ7xl3AL8zzJw9NK
Y849CImAGu4FO8J+a5zlTEioKMKz3O6jUvZZXpRLQdA6AOZKeiozwNXqKJcNwx8g
7kZTyTxqinq/OeiGiePr9cwElqe8KkwP4pIVpQIDAQABAoIBAQCFGvgZv4w/Wb7p
E0J3eazhrELxOCcevgBbeODEaL7ZfozYcUzmadSboeKLhpLsZtse3NPMGGgTfnhm
ro3oHkIQAlVVtqmjtZ0gjlpd/SLxdFOgGtuRGeFtkz1X0foi2QC3DZ/mZK0uT3gX
bHD2CcMi8bCj4VSWkBQmHxzV/jGqrUVkkEbfXjeujkXmcv9X/8Sj1UvnIIECWK4r
HT6lRoy75tWa+/wfcNq+oFEfw20VySSp+fh7O0swXohdTdtZxYXaEAXjE8ptdodW
8/sKAsPbjfnmqstkfq3iWl3vImkNcyuOxBtar0mL57kHfEOSmTQJEsNhZQNpSDFI
YNdKX79BAoGBANIGkUugCxX19uzcRJcGJuhYiDdQ+Ew7njN5nqjhMRnV023iFLiv
3xaUxBUMgDQX/iHYENC/Z1uJZr4f/gCMy4Ze9Nqm++p9QEyw3zyKDfXu486lKdyx
XSLiBZQ6kGksGBHLMqwwzZ9bL942chVjpnO3kn/RPt6X57t8iF1JREKTAoGBAMot
GmQyKs05pVMAvmk7ZPD3V0fn2PigOxLRUuGRA8IhGDNR7FBWkbuUt1kniV11RCJ+
L3+xFA937FpqqzSzKodu6mXUwtKxMK5KiPdsKQqUmMiuOnyzaPpC+iz9qe65OVA6
W+9BWk9cmN5XqyaIUDdCb5SLKOPhkMEfskDp4NHnAoGAbPm1aC0Js4Jldi8wc8Bg
bcyKGVGtFDkW9BSV64C1LneRdgGJyO6QbbIRL+7FksIkPcFTsEywP4HCysHk1Lo5
XGZm3BEqw1fsBh78JfhoGAS1NWLjnrx03AW06V2d0sRrVMg/abME7jutUbqkZU7I
bmCA5ktXOL5PIiwSwXyjq3sCgYEApEZfslg9BQI4/ieVkCXtkAo5xjhxyRtQxKqH
ILdXCW8gndqMHH8q7PMaw3tnlyPImApWB/hXZ3Y2+wS/VhPak68hEFr/bnkBKC1x
+zDMbEdvmWhQJ7ETtH2lj9cRM+MW2cSBnPdKLT/9CnTLoYSTQUNfLKCiOf+3QeTC
TxJ6VbMCgYEAi9Yp/M4aaeuA6UieSNLZcQpQLvVLDtHg8uyeD7yIH/l6KfT1X5mO
qyDPy5ZoeaNFLrcquIrjd0Z5HcZQF0G3/gVlgql7AGMiWgK2wPmCV6HgGiZEdua2
OqPxJuZbP3aANy19Vn5k+W8wE81Neh8RjO2CC8f0wkSxHarBKDJxR38=
-----END RSA PRIVATE KEY-----`
	certData6 := `-----BEGIN CERTIFICATE-----
MIICojCCAkigAwIBAgIIFizowlwTTkAwCgYIKoZIzj0EAwIwgaUxCzAJBgNVBAYT
AkNOMRAwDgYDVQQIEwdCZWlqaW5nMRkwFwYDVQQHExBIYWlkaWFuIERpc3RyaWN0
MRUwEwYDVQQJEwxCYWlkdSBDYW1wdXMxDzANBgNVBBETBjEwMDA5MzEeMBwGA1UE
ChMVTGludXggRm91bmRhdGlvbiBFZGdlMQ8wDQYDVQQLEwZCQUVUWUwxEDAOBgNV
BAMTB3Jvb3QuY2EwHhcNMjAwODIwMDcxODA5WhcNNDAwODE1MDcxODA5WjCBpDEL
MAkGA1UEBhMCQ04xEDAOBgNVBAgTB0JlaWppbmcxGTAXBgNVBAcTEEhhaWRpYW4g
RGlzdHJpY3QxFTATBgNVBAkTDEJhaWR1IENhbXB1czEPMA0GA1UEERMGMTAwMDkz
MR4wHAYDVQQKExVMaW51eCBGb3VuZGF0aW9uIEVkZ2UxDzANBgNVBAsTBkJBRVRZ
TDEPMA0GA1UEAxMGc2VydmVyMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEQXCQ
TGn4+frJYOumFk8gs8BIbgduEuiHonhYdJTFGIPLiOqPQoIvDmICod7W0oIzYYXw
TF4NfadliSryoXx9IaNhMF8wDgYDVR0PAQH/BAQDAgWgMB0GA1UdJQQWMBQGCCsG
AQUFBwMCBggrBgEFBQcDATAMBgNVHRMBAf8EAjAAMCAGA1UdEQQZMBeCCWxvY2Fs
aG9zdIcEAAAAAIcEfwAAATAKBggqhkjOPQQDAgNIADBFAiB5vz8+oob7SkN54uf7
RErbE4tWT5AHtgBgIs3A+TjnyQIhAPvnL8W1dq4qdkVr0eiH5He0xNHdsQc6eWxS
RcKyjhh1
-----END CERTIFICATE-----`

	cert6 := &models.Certificate{
		Name:        "cert6",
		Description: "desp6",
		Data: models.CertificateDataItem{
			Key:  keyData6,
			Cert: certData6,
		},
	}

	err = cert6.ParseCertInfo()
	assert.Equal(t, err.Error(), "tls: private key type does not match public key type")

	cert7 := &models.Certificate{
		Name:        "cert7",
		Description: "desp7",
		Data:        models.CertificateDataItem{
			Key:  keyData6,
			Cert: certData6,
		},
	}

	mkSecretService.EXPECT().Get(gomock.Any(), cert7.Name, gomock.Any()).Return(nil, fmt.Errorf("error")).Times(1)
	w = httptest.NewRecorder()
	body, _ = json.Marshal(cert7)
	req, _ = http.NewRequest(http.MethodPost, "/v1/certificates", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	cert8 := &models.Certificate{
		Name:        "cert8",
		Description: "desp8",
		Data: models.CertificateDataItem{
			Key:  keyData2,
			Cert: certData2,
		},
	}

	err = cert8.ParseCertInfo()
	assert.NoError(t, err)

	temp8 := cert8.ToSecret()

	mkSecretService.EXPECT().Get(gomock.Any(), cert8.Name, gomock.Any()).Return(nil, nil).Times(1)
	mkSecretService.EXPECT().Create(gomock.Any(), temp8).Return(nil, fmt.Errorf("error")).Times(1)
	w = httptest.NewRecorder()
	body, _ = json.Marshal(cert8)
	req, _ = http.NewRequest(http.MethodPost, "/v1/certificates", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	cert9 := &models.Certificate{
		Name:        "cert9",
		Description: "desp9",
		Data:        models.CertificateDataItem{},
	}

	w = httptest.NewRecorder()
	body, _ = json.Marshal(cert9)
	req, _ = http.NewRequest(http.MethodPost, "/v1/certificates", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateCertificate(t *testing.T) {
	api, router, mockCtl := initCertificateAPI(t)
	defer mockCtl.Finish()
	mkSecretService := ms.NewMockSecretService(mockCtl)
	api.AppCombinedService = &service.AppCombinedService{
		Secret: mkSecretService,
	}

	ns := "ns"
	name := "cert"
	keyData1 := `-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEApd4vdBUOtuW+v4YXkqJUFyz18N3tk00Ff3X5SfT036qEC4SB
FsZDyNAX5bhXKkKMX0oAWt7XowpCsUtppChkibfAQNtQbOBOTVC9UJEtef/edEVz
PxyEb8jBq359pee3zZoV55EwxN7uWq1oS/zZ+i6rSncfqKNzI/9kLoyxda6CgnAB
rIpEwdMtqWs8b1Y/0q+s/94owc4jzypxTWSwXLr13s7p2ZDNaJ7xl3AL8zzJw9NK
Y849CImAGu4FO8J+a5zlTEioKMKz3O6jUvZZXpRLQdA6AOZKeiozwNXqKJcNwx8g
7kZTyTxqinq/OeiGiePr9cwElqe8KkwP4pIVpQIDAQABAoIBAQCFGvgZv4w/Wb7p
E0J3eazhrELxOCcevgBbeODEaL7ZfozYcUzmadSboeKLhpLsZtse3NPMGGgTfnhm
ro3oHkIQAlVVtqmjtZ0gjlpd/SLxdFOgGtuRGeFtkz1X0foi2QC3DZ/mZK0uT3gX
bHD2CcMi8bCj4VSWkBQmHxzV/jGqrUVkkEbfXjeujkXmcv9X/8Sj1UvnIIECWK4r
HT6lRoy75tWa+/wfcNq+oFEfw20VySSp+fh7O0swXohdTdtZxYXaEAXjE8ptdodW
8/sKAsPbjfnmqstkfq3iWl3vImkNcyuOxBtar0mL57kHfEOSmTQJEsNhZQNpSDFI
YNdKX79BAoGBANIGkUugCxX19uzcRJcGJuhYiDdQ+Ew7njN5nqjhMRnV023iFLiv
3xaUxBUMgDQX/iHYENC/Z1uJZr4f/gCMy4Ze9Nqm++p9QEyw3zyKDfXu486lKdyx
XSLiBZQ6kGksGBHLMqwwzZ9bL942chVjpnO3kn/RPt6X57t8iF1JREKTAoGBAMot
GmQyKs05pVMAvmk7ZPD3V0fn2PigOxLRUuGRA8IhGDNR7FBWkbuUt1kniV11RCJ+
L3+xFA937FpqqzSzKodu6mXUwtKxMK5KiPdsKQqUmMiuOnyzaPpC+iz9qe65OVA6
W+9BWk9cmN5XqyaIUDdCb5SLKOPhkMEfskDp4NHnAoGAbPm1aC0Js4Jldi8wc8Bg
bcyKGVGtFDkW9BSV64C1LneRdgGJyO6QbbIRL+7FksIkPcFTsEywP4HCysHk1Lo5
XGZm3BEqw1fsBh78JfhoGAS1NWLjnrx03AW06V2d0sRrVMg/abME7jutUbqkZU7I
bmCA5ktXOL5PIiwSwXyjq3sCgYEApEZfslg9BQI4/ieVkCXtkAo5xjhxyRtQxKqH
ILdXCW8gndqMHH8q7PMaw3tnlyPImApWB/hXZ3Y2+wS/VhPak68hEFr/bnkBKC1x
+zDMbEdvmWhQJ7ETtH2lj9cRM+MW2cSBnPdKLT/9CnTLoYSTQUNfLKCiOf+3QeTC
TxJ6VbMCgYEAi9Yp/M4aaeuA6UieSNLZcQpQLvVLDtHg8uyeD7yIH/l6KfT1X5mO
qyDPy5ZoeaNFLrcquIrjd0Z5HcZQF0G3/gVlgql7AGMiWgK2wPmCV6HgGiZEdua2
OqPxJuZbP3aANy19Vn5k+W8wE81Neh8RjO2CC8f0wkSxHarBKDJxR38=
-----END RSA PRIVATE KEY-----`
	certData1 := `-----BEGIN CERTIFICATE-----
MIICvjCCAaYCCQCrqj3SrB8ZxjANBgkqhkiG9w0BAQsFADAhMQswCQYDVQQGEwJD
TjESMBAGA1UEAwwJY2F0dGxlLWNhMB4XDTIwMDkwMjA4MTE1NloXDTMwMDgzMTA4
MTE1NlowITELMAkGA1UEBhMCQ04xEjAQBgNVBAMMCWNhdHRsZS1jYTCCASIwDQYJ
KoZIhvcNAQEBBQADggEPADCCAQoCggEBAKXeL3QVDrblvr+GF5KiVBcs9fDd7ZNN
BX91+Un09N+qhAuEgRbGQ8jQF+W4VypCjF9KAFre16MKQrFLaaQoZIm3wEDbUGzg
Tk1QvVCRLXn/3nRFcz8chG/Iwat+faXnt82aFeeRMMTe7lqtaEv82fouq0p3H6ij
cyP/ZC6MsXWugoJwAayKRMHTLalrPG9WP9KvrP/eKMHOI88qcU1ksFy69d7O6dmQ
zWie8ZdwC/M8ycPTSmPOPQiJgBruBTvCfmuc5UxIqCjCs9zuo1L2WV6US0HQOgDm
SnoqM8DV6iiXDcMfIO5GU8k8aop6vznohonj6/XMBJanvCpMD+KSFaUCAwEAATAN
BgkqhkiG9w0BAQsFAAOCAQEAPUUFgOOenmB8eexetT2o8HKFb13jeQFEiasuOyM4
gb8u1L0e3OUHLQi4CJml7szw10sdc4/x6wqtxJd8K11m1xJ8NxqmsmGGedzGtIYu
lZ6xthcgpvwSoDqVQr4FYYYjCmurNsaYFu7GVSqPCgL0SbiqNmfLK5sKf9BVkMlt
cDNsVm/zxNtmRFyTM8uGX3RoA/kltxnElAOqYJEAcePMQ4IoThfif/ql3Y1D6J5u
O0o7IM5tlD8UhYfAYVH/xyKxdv6K23zHgbuF4QWrBEnYey0eqznO+5bNV7bz0pkc
tu5nww5RdjCz4Uks08P2GNmZjLO81MgYkhR7B9wi3KDNxg==
-----END CERTIFICATE-----`

	cert1 := &models.Certificate{
		Name:        name,
		Description: "desp1",
		Data: models.CertificateDataItem{
			Key:  keyData1,
			Cert: certData1,
		},
	}

	err := cert1.ParseCertInfo()
	assert.NoError(t, err)

	//temp1 := cert1.ToSecret()

	res1 := &specV1.Secret{
		Name:      cert1.Name,
		Namespace: ns,
		Labels: map[string]string{
			specV1.SecretLabel: specV1.SecretCustomCertificate,
		},
		Data: map[string][]byte{
			"key":                []byte(keyData1),
			"cert":               []byte(certData1),
			"signatureAlgorithm": []byte(`SHA256-RSA`),
			"effectiveTime":      []byte("2020-09-02 08:11:56 +0000 UTC"),
			"expiredTime":        []byte("2030-08-31 08:11:56 +0000 UTC"),
			"serialNumber":       []byte("12369767301566634438"),
			"issuer":             []byte("cattle-ca"),
			"fingerPrint":        []byte("5F:D4:2B:70:C3:8F:54:FB:3B:AA:28:C8:4E:47:F3:38:FA:74:83:94"),
		},
		Description: cert1.Description,
		Version:     "1234",
	}

	keyData2 := `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIPcuJ/l+/s8PXvAN5M6VNBZrKD4HDW6n6Y4rQCYinF5doAoGCCqGSM49
AwEHoUQDQgAEQXCQTGn4+frJYOumFk8gs8BIbgduEuiHonhYdJTFGIPLiOqPQoIv
DmICod7W0oIzYYXwTF4NfadliSryoXx9IQ==
-----END EC PRIVATE KEY-----`
	certData2 := `-----BEGIN CERTIFICATE-----
MIICojCCAkigAwIBAgIIFizowlwTTkAwCgYIKoZIzj0EAwIwgaUxCzAJBgNVBAYT
AkNOMRAwDgYDVQQIEwdCZWlqaW5nMRkwFwYDVQQHExBIYWlkaWFuIERpc3RyaWN0
MRUwEwYDVQQJEwxCYWlkdSBDYW1wdXMxDzANBgNVBBETBjEwMDA5MzEeMBwGA1UE
ChMVTGludXggRm91bmRhdGlvbiBFZGdlMQ8wDQYDVQQLEwZCQUVUWUwxEDAOBgNV
BAMTB3Jvb3QuY2EwHhcNMjAwODIwMDcxODA5WhcNNDAwODE1MDcxODA5WjCBpDEL
MAkGA1UEBhMCQ04xEDAOBgNVBAgTB0JlaWppbmcxGTAXBgNVBAcTEEhhaWRpYW4g
RGlzdHJpY3QxFTATBgNVBAkTDEJhaWR1IENhbXB1czEPMA0GA1UEERMGMTAwMDkz
MR4wHAYDVQQKExVMaW51eCBGb3VuZGF0aW9uIEVkZ2UxDzANBgNVBAsTBkJBRVRZ
TDEPMA0GA1UEAxMGc2VydmVyMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEQXCQ
TGn4+frJYOumFk8gs8BIbgduEuiHonhYdJTFGIPLiOqPQoIvDmICod7W0oIzYYXw
TF4NfadliSryoXx9IaNhMF8wDgYDVR0PAQH/BAQDAgWgMB0GA1UdJQQWMBQGCCsG
AQUFBwMCBggrBgEFBQcDATAMBgNVHRMBAf8EAjAAMCAGA1UdEQQZMBeCCWxvY2Fs
aG9zdIcEAAAAAIcEfwAAATAKBggqhkjOPQQDAgNIADBFAiB5vz8+oob7SkN54uf7
RErbE4tWT5AHtgBgIs3A+TjnyQIhAPvnL8W1dq4qdkVr0eiH5He0xNHdsQc6eWxS
RcKyjhh1
-----END CERTIFICATE-----`

	res2 := &specV1.Secret{
		Name:      name,
		Namespace: ns,
		Labels: map[string]string{
			specV1.SecretLabel: specV1.SecretCustomCertificate,
		},
		Data: map[string][]byte{
			"key":                []byte(keyData2),
			"cert":               []byte(certData2),
			"signatureAlgorithm": []byte(`ECDSA-SHA256`),
			"effectiveTime":      []byte("2020-08-20 07:18:09 +0000 UTC"),
			"expiredTime":        []byte("2040-08-15 07:18:09 +0000 UTC"),
			"serialNumber":       []byte("1597907889275752000"),
			"issuer":             []byte("root.ca"),
			"fingerPrint":        []byte("8A:80:76:69:91:B0:EB:A6:45:E2:FE:88:03:F0:66:4A:A9:BB:16:96"),
		},
		CreationTimestamp: time.Now(),
		UpdateTimestamp:   time.Now(),
		Description:       "desp2",
		Version:           "1234",
	}

	mkSecretService.EXPECT().Get(gomock.Any(), name, gomock.Any()).Return(res2, nil).Times(1)
	mkSecretService.EXPECT().Update(gomock.Any(), gomock.Any()).Return(res1, nil).Times(1)
	w := httptest.NewRecorder()
	body, _ := json.Marshal(cert1)
	req, _ := http.NewRequest(http.MethodPut, "/v1/certificates/"+name, bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	mkSecretService.EXPECT().Get(gomock.Any(), name, gomock.Any()).Return(res1, nil).Times(1)
	cert1.Data.Key = ""
	w = httptest.NewRecorder()
	body, _ = json.Marshal(cert1)
	req, _ = http.NewRequest(http.MethodPut, "/v1/certificates/"+name, bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	cert1.Data.Key = keyData1
	mkSecretService.EXPECT().Get(gomock.Any(), name, gomock.Any()).Return(res2, nil).Times(1)
	mkSecretService.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
	w = httptest.NewRecorder()
	body, _ = json.Marshal(cert1)
	req, _ = http.NewRequest(http.MethodPut, "/v1/certificates/"+name, bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetCertificate(t *testing.T) {
	api, router, mockCtl := initCertificateAPI(t)
	defer mockCtl.Finish()
	mkSecretService := ms.NewMockSecretService(mockCtl)
	api.AppCombinedService = &service.AppCombinedService{
		Secret: mkSecretService,
	}

	ns := "default"
	name := "cert"
	keyData2 := `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIPcuJ/l+/s8PXvAN5M6VNBZrKD4HDW6n6Y4rQCYinF5doAoGCCqGSM49
AwEHoUQDQgAEQXCQTGn4+frJYOumFk8gs8BIbgduEuiHonhYdJTFGIPLiOqPQoIv
DmICod7W0oIzYYXwTF4NfadliSryoXx9IQ==
-----END EC PRIVATE KEY-----`
	certData2 := `-----BEGIN CERTIFICATE-----
MIICojCCAkigAwIBAgIIFizowlwTTkAwCgYIKoZIzj0EAwIwgaUxCzAJBgNVBAYT
AkNOMRAwDgYDVQQIEwdCZWlqaW5nMRkwFwYDVQQHExBIYWlkaWFuIERpc3RyaWN0
MRUwEwYDVQQJEwxCYWlkdSBDYW1wdXMxDzANBgNVBBETBjEwMDA5MzEeMBwGA1UE
ChMVTGludXggRm91bmRhdGlvbiBFZGdlMQ8wDQYDVQQLEwZCQUVUWUwxEDAOBgNV
BAMTB3Jvb3QuY2EwHhcNMjAwODIwMDcxODA5WhcNNDAwODE1MDcxODA5WjCBpDEL
MAkGA1UEBhMCQ04xEDAOBgNVBAgTB0JlaWppbmcxGTAXBgNVBAcTEEhhaWRpYW4g
RGlzdHJpY3QxFTATBgNVBAkTDEJhaWR1IENhbXB1czEPMA0GA1UEERMGMTAwMDkz
MR4wHAYDVQQKExVMaW51eCBGb3VuZGF0aW9uIEVkZ2UxDzANBgNVBAsTBkJBRVRZ
TDEPMA0GA1UEAxMGc2VydmVyMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEQXCQ
TGn4+frJYOumFk8gs8BIbgduEuiHonhYdJTFGIPLiOqPQoIvDmICod7W0oIzYYXw
TF4NfadliSryoXx9IaNhMF8wDgYDVR0PAQH/BAQDAgWgMB0GA1UdJQQWMBQGCCsG
AQUFBwMCBggrBgEFBQcDATAMBgNVHRMBAf8EAjAAMCAGA1UdEQQZMBeCCWxvY2Fs
aG9zdIcEAAAAAIcEfwAAATAKBggqhkjOPQQDAgNIADBFAiB5vz8+oob7SkN54uf7
RErbE4tWT5AHtgBgIs3A+TjnyQIhAPvnL8W1dq4qdkVr0eiH5He0xNHdsQc6eWxS
RcKyjhh1
-----END CERTIFICATE-----`
	res := &specV1.Secret{
		Name:      name,
		Namespace: ns,
		Labels: map[string]string{
			specV1.SecretLabel: specV1.SecretCustomCertificate,
		},
		Data: map[string][]byte{
			"key":                []byte(keyData2),
			"cert":               []byte(certData2),
			"signatureAlgorithm": []byte(`ECDSA-SHA256`),
			"effectiveTime":      []byte("2020-08-20 07:18:09 +0000 UTC"),
			"expiredTime":        []byte("2040-08-15 07:18:09 +0000 UTC"),
			"serialNumber":       []byte("1597907889275752000"),
			"issuer":             []byte("root.ca"),
			"fingerPrint":        []byte("8A:80:76:69:91:B0:EB:A6:45:E2:FE:88:03:F0:66:4A:A9:BB:16:96"),
		},
		CreationTimestamp: time.Now(),
		UpdateTimestamp:   time.Now(),
		Description:       "desp2",
		Version:           "1234",
	}

	mkSecretService.EXPECT().Get(ns, name, "").Return(res, nil).Times(1)
	req, _ := http.NewRequest(http.MethodGet, "/v1/certificates/"+name, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	mkSecretService.EXPECT().Get(ns, name, "").Return(nil, fmt.Errorf("error")).Times(1)
	req, _ = http.NewRequest(http.MethodGet, "/v1/certificates/"+name, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeleteCertificate(t *testing.T) {
	api, router, mockCtl := initCertificateAPI(t)
	defer mockCtl.Finish()

	mkSecretService := ms.NewMockSecretService(mockCtl)
	api.AppCombinedService = &service.AppCombinedService{
		Secret: mkSecretService,
	}

	mkIndexService := ms.NewMockIndexService(mockCtl)
	api.Index = mkIndexService

	ns := "default"
	name := "cert"

	mkSecretService.EXPECT().Delete(ns, name).Return(nil)
	mkIndexService.EXPECT().ListAppIndexBySecret(gomock.Any(), gomock.Any()).Return(nil, nil)
	// 200
	req, _ := http.NewRequest(http.MethodDelete, "/v1/certificates/"+name, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestListCertificate(t *testing.T) {
	api, router, mockCtl := initCertificateAPI(t)
	defer mockCtl.Finish()

	mkSecretService := ms.NewMockSecretService(mockCtl)
	api.AppCombinedService = &service.AppCombinedService{
		Secret: mkSecretService,
	}

	mClist := &models.SecretList{
		Total: 0,
	}

	mkSecretService.EXPECT().List("default", gomock.Any()).Return(mClist, nil)

	// 200
	req, _ := http.NewRequest(http.MethodGet, "/v1/certificates", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetAppByCertificate(t *testing.T) {
	api, router, mockCtl := initCertificateAPI(t)
	defer mockCtl.Finish()

	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sSecret := ms.NewMockSecretService(mockCtl)
	api.AppCombinedService = &service.AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}

	sNode, sIndex := ms.NewMockNodeService(mockCtl), ms.NewMockIndexService(mockCtl)
	api.Node, api.Index = sNode, sIndex

	appNames := []string{"app1", "app2", "app3"}
	apps := []*specV1.Application{
		{
			Namespace: "default",
			Name:      appNames[0],
		},
		{
			Namespace: "default",
			Name:      appNames[1],
		},
		{
			Namespace: "default",
			Name:      appNames[2],
		},
	}

	mConfSecret3 := &specV1.Secret{
		Namespace:   "default",
		Name:        "abc",
		Description: "haha",
		Version:     "5",
		Labels: map[string]string{
			specV1.SecretLabel: specV1.SecretCustomCertificate,
		},
	}

	sSecret.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(mConfSecret3, nil)

	sIndex.EXPECT().ListAppIndexBySecret(mConfSecret3.Namespace, mConfSecret3.Name).Return(appNames, nil).Times(1)
	sApp.EXPECT().Get(mConfSecret3.Namespace, appNames[0], "").Return(apps[0], nil).AnyTimes()
	sApp.EXPECT().Get(mConfSecret3.Namespace, appNames[1], "").Return(apps[1], nil).AnyTimes()
	sApp.EXPECT().Get(mConfSecret3.Namespace, appNames[2], "").Return(apps[2], nil).AnyTimes()

	w4 := httptest.NewRecorder()
	req4, _ := http.NewRequest(http.MethodGet, "/v1/certificates/abc/apps", nil)
	router.ServeHTTP(w4, req4)
	assert.Equal(t, http.StatusOK, w4.Code)
}
