package plugin

import (
	"crypto/x509"
	"io"
	"time"
)

//go:generate mockgen -destination=../mock/plugin/pki.go -package=plugin -source=pki.go

type Cert struct {
	CertId      string    `db:"cert_id"`
	ParentId    string    `db:"parent_id"`
	Type        string    `db:"type"`
	CommonName  string    `db:"common_name"`
	Csr         string    `db:"csr"`         // base64
	Content     string    `db:"content"`     // base64
	PrivateKey  string    `db:"private_key"` // base64
	Description string    `db:"description"`
	NotBefore   time.Time `db:"not_before"`
	NotAfter    time.Time `db:"not_after"`
}

type PKI interface {
	// root cert
	GetRootCertID() string
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

type PKIStorage interface {
	CreateCert(cert Cert) error
	DeleteCert(certId string) error
	UpdateCert(cert Cert) error
	GetCert(certId string) (*Cert, error)
	CountCertByParentId(parentId string) (int, error)
	io.Closer
}
