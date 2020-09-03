package models

import (
	"bytes"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"net"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/baetyl/baetyl-go/v2/errors"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/jinzhu/copier"
)

// AltNames contains the domain names and IP addresses that will be added
// to the API Server's x509 certificate SubAltNames field. The values will
// be passed directly to the x509.Certificate object.
type AltNames struct {
	DNSNames []string   `json:"dnsNames,omitempty"`
	IPs      []net.IP   `json:"ips,omitempty"`
	Emails   []string   `json:"emails,omitempty"`
	URIs     []*url.URL `json:"uris,omitempty"`
}

// PEMCredential holds a certificate, private key pem data
type PEMCredential struct {
	CertPEM []byte
	KeyPEM  []byte
	CertId  string
}

// CertStorage contains certName and keyName which can be used to
// storage certificate and private key pem data to secret.
type CertStorage struct {
	CertName string
	KeyName  string
}

//TODO: move to baetyl-go
const SecretCertificate = "custom-certificate"

// Certificate Certificate
type Certificate struct {
	Name               string              `json:"name,omitempty" validate:"omitempty,resourceName,nonBaetyl"`
	Namespace          string              `json:"namespace,omitempty"`
	SignatureAlgorithm string              `json:"signatureAlgorithm,omitempty"`
	EffectiveTime      string              `json:"effectiveTime,omitempty"`
	ExpiredTime        string              `json:"expiredTime,omitempty"`
	SerialNumber       string              `json:"serialNumber,omitempty"`
	Issuer             string              `json:"issuer,omitempty"`
	FingerPrint        string              `json:"fingerPrint,omitempty"`
	Data               CertificateDataItem `json:"data,omitempty"`
	CreationTimestamp  time.Time           `json:"createTime,omitempty"`
	UpdateTimestamp    time.Time           `json:"updateTime,omitempty"`
	Description        string              `json:"description"`
	Version            string              `json:"version,omitempty"`
}

type CertificateDataItem struct {
	Key  string `json:"key,omitempty"`
	Cert string `json:"cert,omitempty"`
}

// CertificateList Certificate List
type CertificateList struct {
	Total       int           `json:"total"`
	ListOptions *ListOptions  `json:"listOptions"`
	Items       []Certificate `json:"items"`
}

func (r *Certificate) Equal(target *Certificate) bool {
	return reflect.DeepEqual(r.Data, target.Data) &&
		reflect.DeepEqual(r.Description, target.Description)
}

func (r *Certificate) ToSecret() *specV1.Secret {
	res := &specV1.Secret{
		Labels: map[string]string{
			specV1.SecretLabel: SecretCertificate,
		},
	}
	err := copier.Copy(res, r)
	if err != nil {
		panic(fmt.Sprintf("copier exception: %s", err.Error()))
	}
	res.Data = map[string][]byte{
		"key":                []byte(r.Data.Key),
		"cert":               []byte(r.Data.Cert),
		"signatureAlgorithm": []byte(r.SignatureAlgorithm),
		"effectiveTime":      []byte(r.EffectiveTime),
		"expiredTime":        []byte(r.ExpiredTime),
		"serialNumber":       []byte(r.SerialNumber),
		"issuer":             []byte(r.Issuer),
		"fingerPrint":        []byte(r.FingerPrint),
	}
	return res
}

func (r *Certificate) ParseCertInfo() error {
	_, err := tls.X509KeyPair([]byte(r.Data.Cert), []byte(r.Data.Key))
	if err != nil {
		return err
	}

	var block *pem.Block
	rest := []byte(r.Data.Cert)
	block, rest = pem.Decode(rest)
	if block == nil {
		return errors.New("failed to find cert that matched private key")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return errors.Errorf("failed to parse certificate, err: %s", err)
	}

	r.FingerPrint = fingerprint(block.Bytes)
	r.Issuer = cert.Issuer.CommonName
	r.ExpiredTime = cert.NotAfter.String()
	r.EffectiveTime = cert.NotBefore.String()
	r.SerialNumber = cert.SerialNumber.String()
	r.SignatureAlgorithm = cert.SignatureAlgorithm.String()

	return nil
}

func fingerprint(data []byte) string {
	digest := sha1.Sum(data)
	buf := &bytes.Buffer{}
	for i := 0; i < len(digest); i++ {
		if buf.Len() > 0 {
			buf.WriteString(":")
		}
		buf.WriteString(strings.ToUpper(hex.EncodeToString(digest[i : i+1])))
	}
	return buf.String()
}
