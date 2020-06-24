package plugin

import (
	"crypto/x509"
	"io"
)

//go:generate mockgen -destination=../mock/plugin/pki.go -package=plugin github.com/baetyl/baetyl-cloud/plugin PKI

const (
	// 证书有效期，以天为单位 [1, 50*365]
	DefaultRootDuration = 50 * 365
	DefaultCertDuration = 20 * 365
)

type PKI interface {
	// root cert
	// info : 生成根证书的相关信息   parentId : 上一级根证书id，可为空
	CreateRootCert(info *x509.CertificateRequest, parentId string) (string, error)
	GetRootCert(rootId string) ([]byte, error)
	DeleteRootCert(rootId string) error

	// server cert
	CreateServerCert(csr []byte, rootId string) (string, error)
	GetServerCert(certId string) ([]byte, error)
	DeleteServerCert(certId string) error

	// client cert
	CreateClientCert(csr []byte, rootId string) (string, error)
	GetClientCert(certId string) ([]byte, error)
	DeleteClientCert(certId string) error

	// close
	io.Closer
}
