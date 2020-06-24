package util

import (
	"crypto/x509"
	"encoding/base64"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	errPEMEncodedKeyData = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAvzq/jxpdJlIdJDsuhZd6hrTo+Tv7NNPyE2iv8IIxZtLDpMEG
tqtltCGeZEbQwpPjntgsnIlzC7aRBW5kLx3NU9o4YSKe+1kdVUqnWKd8NkRcCke0
mbSuMt4efQ8/jgpRl75wRABciRA3kyYNmxDll7rZJ7gvM4HSPNvYUT0AtqXfrGFL
U2oFEwMUoE/Y3oNNiSVl4W4yy4EXBIaUfpukLN56vc+6T/PqwhcJS5SJ20rWFP9o
V5jgNPoUQfAbwt5fqVCvgoojshJe500WKbNM85TkltrEyCBacX2q4PfRBtMIip2N
99ic2/JHDrBpGokUtii2jwgNpEGd58aTrJFcCwIDAQABAoIBAQCPQmrJwT1xfOi6
AOAsUrwG+LbkKIxcGj/rTG0qJ87ssg7BeZ4W8TaDCQCJLQguDO7kTBy3tL0MVFw2
jmndZl0xaXT9SBuEP1GbWQ8fenGykenGBcwFrncmQoLiu66oMZDItnGScBbi09Re
l45v0eu+jMssV259Ds+6qhRXX/UN3s8XCuOmFq8R1cNwnMoDIPDO2EIfyCFxBD/L
VIkbUSirypnQV/EG5o3TYA92nzAax166prUzZ53Kyd3Zh5KqUsUq6A==
-----END RSA PRIVATE KEY-----`

	unSupportedBlockTypeKeyData = `-----BEGIN OPENSSH PRIVATE KEY-----
MIIEpAIBAAKCAQEA0Vrs3CL9ojFjMXdyNWM/y7aMHgXAtxIAAN19Ew18TqYNTO5K
KdiXahgwbIREV9ZQGegabl3I5QmtrACi6Z+2gH6MG53Eb0kE+087hifUO/+SNqhb
Jkksqd7dao6H/t5CUBjEpdRIa1AscBdz67tL0InkCL3Q9aCu8GK9EW/+HdNDh0sU
VPEtWat01UrD3ktXCUALbjef4piEpR0GIQ1y/erZ1AanfFSr0j54cv1Zyne60TAF
wPhHqMx6Pfaap3sT7OvUExCNVc6KXV608vF6/56CXQvXVkZB3zv0FdIFpkwhf4cz
rrRYY+boGUrJHQmi85gwnC+LNyZ3abb+zSFyZwIDAQABAoIBAFTEpDNeV6RcqvVU
kAHd9e7eM03Utntp5mZzSDl2tGaEBc6ojY70DBsBQFowFBwcwsI6oLkfcECM8q06
dLxz1smgc8qazvbgcgvvwQJJDj5c9S78bCvMZTFC9BQ5MgeYpvEXlkgu9EO/ar7Y
QC+q3r/JlXOUqA9MyIi88iElX5djP9FOW+qRF9asbYxMjyW63+vdRMnzBez8uWU2
9y/K+3VvhNjl5uWKC9nSGmqVorgRBTaUmVBk0yF1lk3yganiYrxqTdN/5CHvKuQ4
WGP0+Cbl0pqZ7+l5FExNABWisxtRWEuzl8DIxVnzzwviRAua3r4i9ZW8YSQBPAIT
CpZNJAECgYEA8EWlkuEyEnPQplA6Zsc6jWT5oYHxE6aEagDcW+pDyAVOtqju89EM
NX2FvBPo+jydKKfWYhQUrY908QVyGL35Hu85U/IGLOFqBEQ8JaxLNlQKPtezHcmH
O37egl2pRA7oRAZ0w0yt3w73yBk/kHgux3WCwCaYPTZpNBlx/9T/XicCgYEA3w8w
h8WVj1yq/iG3HMdSLiPM7dP3XTuqxX7RVN3vTKfGkDeRpTQcYU0ZyPvAMqiZOi6o
cHVC9c2fi6gIQR4G4/I57TP4oRzDe80Bj2OKIsBv0jgf5b1voUxdBBnXBHJg+YRh
NF+M77Jium/JAMmaKNC5o3en6KoLmUIT8V96McECgYAKCM6WaMM/lAiluXoG6tEe
MJZgUV3xFSY4ixqo2ArGorob1MhN9HAPF9Pq++Xh9YAWv5Oreu02JmSa4EBYmi56
RUFeqR/q5esYjIT6icyGU1IuN7HqT41PRcgjJ6g3CGxY0vAza9NjGmvstmk6Llq9
x8GTJsl63PfdziY9qfaURwKBgQC2HU2vHDdOlAcbc2VwPqAvAZW3+x5z5Vo44qCA
HK7atARfDK+B5PjizDMoL7qs4ZAwu5VUM7jWvOns+OS8XYqcotB+hLcSu0wzEJ6c
dlV6qAjj5mTMiozQcWtkBMDTZZsdPOKsAvMrZEZNFyVR2kdd2YQnHXNedy7/Er77
i8tVQQKBgQDWJaJO6jVuorSu6N2AfvUc+BUIicWi6HeN6azow+0R5T6dVhLhj18W
vX8tNbOEnt3k5yqW8AZ8d8x8H/JZ6IWyXnKK/IV7+7R1Z6DRvjo/UuOZOxG54S0f
uEKnjHV/9NaDbxDkJBgcJhGexi2+Xxhk38gOsqnj1MFGneY/Yfb/4g==
-----END OPENSSH PRIVATE KEY-----`

	errBlockTypeCert = `-----BEGIN RSA CERTIFICATE-----
MIIEJjCCAw6gAwIBAgIEALs3ETANBgkqhkiG9w0BAQsFADCBpjELMAkGA1UEBhMC
Q04xEDAOBgNVBAgTB0JlaWppbmcxGTAXBgNVBAcTEEhhaWRpYW4gRGlzdHJpY3Qx
FTATBgNVBAkTDEJhaWR1IENhbXB1czEPMA0GA1UEERMGMTAwMDkzMR4wHAYDVQQK
ExVMaW51eCBGb3VuZGF0aW9uIEVkZ2UxDzANBgNVBAsTBkJBRVRZTDERMA8GA1UE
AxMIbm9kZTEuY2EwHhcNMjAwMzE5MTgxMTE0WhcNMzAwMzE5MTgxMTE0WjCBujEL
MAkGA1UEBhMCQ04xEDAOBgNVBAgTB0JlaWppbmcxGTAXBgNVBAcTEEhhaWRpYW4g
RGlzdHJpY3QxFTATBgNVBAkTDEJhaWR1IENhbXB1czEPMA0GA1UEERMGMTAwMDkz
MR4wHAYDVQQKExVMaW51eCBGb3VuZGF0aW9uIEVkZ2UxDzANBgNVBAsTBkJBRVRZ
TDElMCMGA1UEAxMcbm9kZTEuYnJva2VyLnNlcnZlci5pbnRlcm5hbDCCASIwDQYJ
KoZIhvcNAQEBBQADggEPADCCAQoCggEBAL8FXx14YYu5ieO1WvR4yZ+oDR3fZcx/
BeGOSvvmGV8mzYpVQ2VXVOQ5AdNqfRVbmedNcahQPF5r6hebI3IoDcqEa49FUShc
OOC1HnqLpEF4Q7vKQOhnfeBymm72cPZURQtRgXm5oa6BWaD4pHbbrWxrJOKaVRxr
U/k0DCyrFAYYtSug0eM7IX6KQ9zvWbfd+myMn0BJcf2Csr8myqdEP7s2DUJZ07P5
BABhjYKnd993Q2jbSYTX5gMNm+lhfYJFiITlWRXSi6abvt6iRpVKTW/adlwzL9YA
p6jkSWUjqeSyDqaE9AHtQ84IldUAClT4mkF46bp7S1fC8BqZZ0ul428CAwEAAaNG
MEQwDgYDVR0PAQH/BAQDAgWgMBMGA1UdJQQMMAoGCCsGAQUFBwMBMAwGA1UdEwEB
/wQCMAAwDwYDVR0RBAgwBocEfwAAATANBgkqhkiG9w0BAQsFAAOCAQEAlmHzWiL8
vNK6MK9ynRtWAf+PKQAEbltaAtEW1u64DReouM5zPDnyMKZhpdZ43FqE5fR7Nw1+
Iv4PxLCEzr1krdPAxYFzOjiboD6a3l5K7UVpgLbmGsi7KSHu1iKKv5Ey0tDtYL96
iSo4wgrN5pkHgSCOtXjQ0iin1o105Z2GBmy0k6fG145hEFjTvfFFLO2rV82SS4uL
6xlAzlIa7LpS7U8PsORJaNMiySl+QvM00XTiUhhGQ3U/FiftDc9weg4otxH5nXOB
fxwQi/QxxhGCUSylUDefv+wwUytBH6jr2/oJcas8aAkwIMWxXDX7mSGoRZF7Eh9z
Z5N93gnTd5qEWw==
-----END RSA CERTIFICATE-----`

	errCertData = `-----BEGIN CERTIFICATE-----
MIIEJjCCAw6gAwIBAgIEALs3ETANBgkqhkiG9w0BAQsFADCBpjELMAkGA1UEBhMC
Q04xEDAOBgNVBAgTB0JlaWppbmcxGTAXBgNVBAcTEEhhaWRpYW4gRGlzdHJpY3Qx
FTATBgNVBAkTDEJhaWR1IENhbXB1czEPMA0GA1UEERMGMTAwMDkzMR4wHAYDVQQK
ExVMaW51eCBGb3VuZGF0aW9uIEVkZ2UxDzANBgNVBAsTBkJBRVRZTDERMA8GA1UE
AxMIbm9kZTEuY2EwHhcNMjAwMzE5MTgxMTE0WhcNMzAwMzE5MTgxMTE0WjCBujEL
/wQCMAAwDwYDVR0RBAgwBocEfwAAATANBgkqhkiG9w0BAQsFAAOCAQEAlmHzWiL8
vNK6MK9ynRtWAf+PKQAEbltaAtEW1u64DReouM5zPDnyMKZhpdZ43FqE5fR7Nw1+
Iv4PxLCEzr1krdPAxYFzOjiboD6a3l5K7UVpgLbmGsi7KSHu1iKKv5Ey0tDtYL96
iSo4wgrN5pkHgSCOtXjQ0iin1o105Z2GBmy0k6fG145hEFjTvfFFLO2rV82SS4uL
6xlAzlIa7LpS7U8PsORJaNMiySl+QvM00XTiUhhGQ3U/FiftDc9weg4otxH5nXOB
fxwQi/QxxhGCUSylUDefv+wwUytBH6jr2/oJcas8aAkwIMWxXDX7mSGoRZF7Eh9z
Z5N93gnTd5qEWw==
-----END CERTIFICATE-----`

	testCertData = `-----BEGIN CERTIFICATE-----
MIIEIzCCAwugAwIBAgIEALs3EDANBgkqhkiG9w0BAQsFADCBpjELMAkGA1UEBhMC
Q04xEDAOBgNVBAgTB0JlaWppbmcxGTAXBgNVBAcTEEhhaWRpYW4gRGlzdHJpY3Qx
FTATBgNVBAkTDEJhaWR1IENhbXB1czEPMA0GA1UEERMGMTAwMDkzMR4wHAYDVQQK
ExVMaW51eCBGb3VuZGF0aW9uIEVkZ2UxDzANBgNVBAsTBkJBRVRZTDERMA8GA1UE
AxMIbm9kZTEuY2EwHhcNMjAwMzE5MTgxMTE0WhcNMzAwMzE5MTgxMTE0WjCBtzEL
MAkGA1UEBhMCQ04xEDAOBgNVBAgTB0JlaWppbmcxGTAXBgNVBAcTEEhhaWRpYW4g
RGlzdHJpY3QxFTATBgNVBAkTDEJhaWR1IENhbXB1czEPMA0GA1UEERMGMTAwMDkz
MR4wHAYDVQQKExVMaW51eCBGb3VuZGF0aW9uIEVkZ2UxDzANBgNVBAsTBkJBRVRZ
TDEiMCAGA1UEAxMZbm9kZTEuYnJva2VyLnRpbWVyLmNsaWVudDCCASIwDQYJKoZI
hvcNAQEBBQADggEPADCCAQoCggEBANFa7Nwi/aIxYzF3cjVjP8u2jB4FwLcSAADd
fRMNfE6mDUzuSinYl2oYMGyERFfWUBnoGm5dyOUJrawAoumftoB+jBudxG9JBPtP
O4Yn1Dv/kjaoWyZJLKne3WqOh/7eQlAYxKXUSGtQLHAXc+u7S9CJ5Ai90PWgrvBi
vRFv/h3TQ4dLFFTxLVmrdNVKw95LVwlAC243n+KYhKUdBiENcv3q2dQGp3xUq9I+
eHL9Wcp3utEwBcD4R6jMej32mqd7E+zr1BMQjVXOil1etPLxev+egl0L11ZGQd87
9BXSBaZMIX+HM660WGPm6BlKyR0JovOYMJwvizcmd2m2/s0hcmcCAwEAAaNGMEQw
DgYDVR0PAQH/BAQDAgWgMBMGA1UdJQQMMAoGCCsGAQUFBwMCMAwGA1UdEwEB/wQC
MAAwDwYDVR0RBAgwBocEfwAAATANBgkqhkiG9w0BAQsFAAOCAQEArHFPtQnyL89/
fKF5Az5LT5S2BffiXI/BhKywYJcGEU6j6P0FabesRAu/cU9UjeAuTEvp1gJQCtg0
EOh0Z+UKwzTXzmn4Sorw5EFdWJ8VkEQ99kVgO3pcP0FoKbG9+56tztFzN1Fw/vSb
drALeJVdCDYXNUUcTc9InX25G3w5g5oUbcJs4LVfL8wmwsqyy5D8Q7eXqpxsQi+O
jUgPK8WBrQDqwtKdY0DnJOAVCBewGpyhru6fJSFfGIfpimJ0NOR+hUWXgU+p5Vfa
t/LUSwDtPIZcwjL3FAOjFi9qGliye8HzfhesVLVn3cSCVkk6uyI8JwMi3GeLabLC
iY1Gr3Xc6w==
-----END CERTIFICATE-----`
)

func TestPrivateKey(t *testing.T) {
	tests := []struct {
		dsa  string
		bits int
	}{
		{
			dsa:  "P224",
			bits: 0,
		},
		{
			dsa:  "P256",
			bits: 0,
		},
		{
			dsa:  "P384",
			bits: 0,
		},
		{
			dsa:  "P521",
			bits: 0,
		},
		{
			dsa:  "rsa",
			bits: 2048,
		},
		{
			dsa:  "rsa",
			bits: 3072,
		},
		{
			dsa:  "rsa",
			bits: 4096,
		},
	}
	target := []string{
		EcPrivateKeyBlockType,
		EcPrivateKeyBlockType,
		EcPrivateKeyBlockType,
		EcPrivateKeyBlockType,
		RsaPrivateKeyBlockType,
		RsaPrivateKeyBlockType,
		RsaPrivateKeyBlockType,
	}
	sigAlgorithms := []x509.SignatureAlgorithm{
		x509.ECDSAWithSHA256,
		x509.ECDSAWithSHA256,
		x509.ECDSAWithSHA384,
		x509.ECDSAWithSHA512,
		x509.SHA256WithRSA,
		x509.SHA384WithRSA,
		x509.SHA512WithRSA,
	}
	for k, tt := range tests {
		t.Run("test", func(t *testing.T) {
			priv, err := GenCertPrivateKey(tt.dsa, tt.bits)
			assert.NoError(t, err)
			assert.NotNil(t, priv)
			assert.Equal(t, priv.Type, target[k])

			// signature algorithm validate
			sigAlgorithm := SigAlgorithmType(priv)
			assert.Equal(t, sigAlgorithms[k], sigAlgorithm)
			// encode private key
			keyPEM, err := EncodeCertPrivateKey(priv)
			assert.NoError(t, err)
			assert.NotNil(t, keyPEM)

			// parse private key
			newPriv, err := ParseCertPrivateKey(keyPEM)
			assert.NoError(t, err)
			assert.NotNil(t, newPriv)
			assert.Equal(t, priv, newPriv)
		})
	}

	// test unRecognized
	testConf := []struct {
		dsa  string
		bits int
	}{
		{
			dsa:  "test",
			bits: 2048,
		},
	}
	priv, err := GenCertPrivateKey(testConf[0].dsa, testConf[0].bits)
	assert.Nil(t, priv)
	assert.NotNil(t, err)
	assert.Equal(t, "unRecognized digital signature algorithm: test", err.Error())

	// test unRecognized encode
	testConf[0].dsa = "P256"
	priv, err = GenCertPrivateKey(testConf[0].dsa, testConf[0].bits)
	assert.NoError(t, err)
	assert.NotNil(t, priv)
	assert.Equal(t, EcPrivateKeyBlockType, priv.Type)
	priv.Type = "test"
	keyPEM, err := EncodeCertPrivateKey(priv)
	assert.Nil(t, keyPEM)
	assert.NotNil(t, err)
	assert.Equal(t, "unRecognized type of PrivateKey", err.Error())
}

func TestParseCertificates(t *testing.T) {
	// error cert data
	certInfo, err := ParseCertificates([]byte(errCertData))
	assert.Error(t, err)

	// normal cert data
	certInfo, err = ParseCertificates([]byte(testCertData))
	assert.NoError(t, err)
	assert.NotNil(t, certInfo)
}

func TestEncodeCertificateByteToPem(t *testing.T) {
	in := "1234567812345678123456781234567812345678123456781234567812345678"
	out := `-----BEGIN CERTIFICATE-----
MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4
MTIzNDU2NzgxMjM0NTY3OA==
-----END CERTIFICATE-----
`
	res := EncodeByteToPem([]byte(in), "CERTIFICATE")
	assert.Equal(t, out, res)
}

func TestEncodeCertificatesAndCertRequest(t *testing.T) {
	certBase64 := "MIIDPzCCAiegAwIBAgIIFhQmuxIIlYAwDQYJKoZIhvcNAQELBQAwgaUxCzAJBgNVBAYTAkNOMRAwDgYDVQQIEwdCZWlqaW5nMRkwFwYDVQQHExBIYWlkaWFuIERpc3RyaWN0MRUwEwYDVQQJEwxCYWlkdSBDYW1wdXMxDzANBgNVBBETBjEwMDA5MzEeMBwGA1UEChMVTGludXggRm91bmRhdGlvbiBFZGdlMQ8wDQYDVQQLEwZCQUVUWUwxEDAOBgNVBAMTB3Jvb3QuY2EwHhcNMjAwNTMxMTUzMjMzWhcNMzAwMzI2MDMzMTUxWjCBrDELMAkGA1UEBhMCQ04xEDAOBgNVBAgTB0JlaWppbmcxGTAXBgNVBAcTEEhhaWRpYW4gRGlzdHJpY3QxFTATBgNVBAkTDEJhaWR1IENhbXB1czEPMA0GA1UEERMGMTAwMDkzMR4wHAYDVQQKExVMaW51eCBGb3VuZGF0aW9uIEVkZ2UxDzANBgNVBAsTBkJBRVRZTDEXMBUGA1UEAxMOZGVmYXVsdC4wNTMxMTIwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAARSRL0W0GGg7sUeWOEs56E6PWE3DFCsgb0eCD70HE2XOHT/BTUzMw36K9W5qHpej2QaLfsvynCfUUDTqtNzh8JfozUwMzAOBgNVHQ8BAf8EBAMCBaAwEwYDVR0lBAwwCgYIKwYBBQUHAwIwDAYDVR0TAQH/BAIwADANBgkqhkiG9w0BAQsFAAOCAQEAmCIzun0CqkM3vMlr4LquDuOGemSlTUq8HJngJPUWL3OALLrOjW+fusl3f8YbQiIIcLsVTx3xUKTcqDAibWgrvi1UdBTCM354uHBD50aCTs1oZEye0eMExYrvsOVbcAZ1Kj3Ywc5zuhtgd9d0NHMdp52nqFp3u/RqODZRPqyzsKbGlA2HaJPj+Bedc0s8K7hkdJY4F5byfgaCk2PBZsXZBo6xIoAjvRg5MIrd2AxHH3k1RMZhDAa40lq2Yt65r36iSC6nb1ewvSNZf3HgUKTtb3pCfaEj7l2fr3k80e3Y9IcfsVtEknyMl9JIxtgjmzzw//wwff4WllatLHYViRZKGw=="
	expectCert := `-----BEGIN CERTIFICATE-----
MIIDPzCCAiegAwIBAgIIFhQmuxIIlYAwDQYJKoZIhvcNAQELBQAwgaUxCzAJBgNV
BAYTAkNOMRAwDgYDVQQIEwdCZWlqaW5nMRkwFwYDVQQHExBIYWlkaWFuIERpc3Ry
aWN0MRUwEwYDVQQJEwxCYWlkdSBDYW1wdXMxDzANBgNVBBETBjEwMDA5MzEeMBwG
A1UEChMVTGludXggRm91bmRhdGlvbiBFZGdlMQ8wDQYDVQQLEwZCQUVUWUwxEDAO
BgNVBAMTB3Jvb3QuY2EwHhcNMjAwNTMxMTUzMjMzWhcNMzAwMzI2MDMzMTUxWjCB
rDELMAkGA1UEBhMCQ04xEDAOBgNVBAgTB0JlaWppbmcxGTAXBgNVBAcTEEhhaWRp
YW4gRGlzdHJpY3QxFTATBgNVBAkTDEJhaWR1IENhbXB1czEPMA0GA1UEERMGMTAw
MDkzMR4wHAYDVQQKExVMaW51eCBGb3VuZGF0aW9uIEVkZ2UxDzANBgNVBAsTBkJB
RVRZTDEXMBUGA1UEAxMOZGVmYXVsdC4wNTMxMTIwWTATBgcqhkjOPQIBBggqhkjO
PQMBBwNCAARSRL0W0GGg7sUeWOEs56E6PWE3DFCsgb0eCD70HE2XOHT/BTUzMw36
K9W5qHpej2QaLfsvynCfUUDTqtNzh8JfozUwMzAOBgNVHQ8BAf8EBAMCBaAwEwYD
VR0lBAwwCgYIKwYBBQUHAwIwDAYDVR0TAQH/BAIwADANBgkqhkiG9w0BAQsFAAOC
AQEAmCIzun0CqkM3vMlr4LquDuOGemSlTUq8HJngJPUWL3OALLrOjW+fusl3f8Yb
QiIIcLsVTx3xUKTcqDAibWgrvi1UdBTCM354uHBD50aCTs1oZEye0eMExYrvsOVb
cAZ1Kj3Ywc5zuhtgd9d0NHMdp52nqFp3u/RqODZRPqyzsKbGlA2HaJPj+Bedc0s8
K7hkdJY4F5byfgaCk2PBZsXZBo6xIoAjvRg5MIrd2AxHH3k1RMZhDAa40lq2Yt65
r36iSC6nb1ewvSNZf3HgUKTtb3pCfaEj7l2fr3k80e3Y9IcfsVtEknyMl9JIxtgj
mzzw//wwff4WllatLHYViRZKGw==
-----END CERTIFICATE-----
`

	// test cert
	data, err := base64.StdEncoding.DecodeString(certBase64)
	assert.NoError(t, err)
	res, err := x509.ParseCertificate(data)
	assert.NoError(t, err)
	certPemByte, err := EncodeCertificates(res)
	assert.NoError(t, err)
	assert.Equal(t, expectCert, string(certPemByte))

	expectCertCSR := "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURSBSRVFVRVNULS0tLS0KTUlJRFB6Q0NBaWVnQXdJQkFnSUlGaFFtdXhJSWxZQXdEUVlKS29aSWh2Y05BUUVMQlFBd2dhVXhDekFKQmdOVgpCQVlUQWtOT01SQXdEZ1lEVlFRSUV3ZENaV2xxYVc1bk1Sa3dGd1lEVlFRSEV4QklZV2xrYVdGdUlFUnBjM1J5CmFXTjBNUlV3RXdZRFZRUUpFd3hDWVdsa2RTQkRZVzF3ZFhNeER6QU5CZ05WQkJFVEJqRXdNREE1TXpFZU1Cd0cKQTFVRUNoTVZUR2x1ZFhnZ1JtOTFibVJoZEdsdmJpQkZaR2RsTVE4d0RRWURWUVFMRXdaQ1FVVlVXVXd4RURBTwpCZ05WQkFNVEIzSnZiM1F1WTJFd0hoY05NakF3TlRNeE1UVXpNak16V2hjTk16QXdNekkyTURNek1UVXhXakNCCnJERUxNQWtHQTFVRUJoTUNRMDR4RURBT0JnTlZCQWdUQjBKbGFXcHBibWN4R1RBWEJnTlZCQWNURUVoaGFXUnAKWVc0Z1JHbHpkSEpwWTNReEZUQVRCZ05WQkFrVERFSmhhV1IxSUVOaGJYQjFjekVQTUEwR0ExVUVFUk1HTVRBdwpNRGt6TVI0d0hBWURWUVFLRXhWTWFXNTFlQ0JHYjNWdVpHRjBhVzl1SUVWa1oyVXhEekFOQmdOVkJBc1RCa0pCClJWUlpUREVYTUJVR0ExVUVBeE1PWkdWbVlYVnNkQzR3TlRNeE1USXdXVEFUQmdjcWhrak9QUUlCQmdncWhrak8KUFFNQkJ3TkNBQVJTUkwwVzBHR2c3c1VlV09FczU2RTZQV0UzREZDc2diMGVDRDcwSEUyWE9IVC9CVFV6TXczNgpLOVc1cUhwZWoyUWFMZnN2eW5DZlVVRFRxdE56aDhKZm96VXdNekFPQmdOVkhROEJBZjhFQkFNQ0JhQXdFd1lEClZSMGxCQXd3Q2dZSUt3WUJCUVVIQXdJd0RBWURWUjBUQVFIL0JBSXdBREFOQmdrcWhraUc5dzBCQVFzRkFBT0MKQVFFQW1DSXp1bjBDcWtNM3ZNbHI0THF1RHVPR2VtU2xUVXE4SEpuZ0pQVVdMM09BTExyT2pXK2Z1c2wzZjhZYgpRaUlJY0xzVlR4M3hVS1RjcURBaWJXZ3J2aTFVZEJUQ00zNTR1SEJENTBhQ1RzMW9aRXllMGVNRXhZcnZzT1ZiCmNBWjFLajNZd2M1enVodGdkOWQwTkhNZHA1Mm5xRnAzdS9ScU9EWlJQcXl6c0tiR2xBMkhhSlBqK0JlZGMwczgKSzdoa2RKWTRGNWJ5ZmdhQ2syUEJac1haQm82eElvQWp2Umc1TUlyZDJBeEhIM2sxUk1aaERBYTQwbHEyWXQ2NQpyMzZpU0M2bmIxZXd2U05aZjNIZ1VLVHRiM3BDZmFFajdsMmZyM2s4MGUzWTlJY2ZzVnRFa255TWw5Skl4dGdqCm16encvL3d3ZmY0V2xsYXRMSFlWaVJaS0d3PT0KLS0tLS1FTkQgQ0VSVElGSUNBVEUgUkVRVUVTVC0tLS0tCg=="
	// test cert request
	csr := &x509.CertificateRequest{
		Raw:                      res.Raw,
		RawTBSCertificateRequest: res.RawTBSCertificate,
		RawSubjectPublicKeyInfo:  res.RawSubjectPublicKeyInfo,
		RawSubject:               res.RawSubject,
		Version:                  res.Version,
		Signature:                res.Signature,
		SignatureAlgorithm:       res.SignatureAlgorithm,
		PublicKeyAlgorithm:       res.PublicKeyAlgorithm,
		PublicKey:                res.PublicKey,
		Subject:                  res.Subject,
		Extensions:               res.Extensions,
		ExtraExtensions:          res.ExtraExtensions,
		DNSNames:                 res.DNSNames,
		EmailAddresses:           res.EmailAddresses,
		IPAddresses:              res.IPAddresses,
		URIs:                     res.URIs,
	}
	csrdata, err := EncodeCertificatesRequest(csr)
	assert.NoError(t, err)
	assert.Equal(t, expectCertCSR, base64.StdEncoding.EncodeToString(csrdata))
}
