package models

import (
	"net"
	"net/url"
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
