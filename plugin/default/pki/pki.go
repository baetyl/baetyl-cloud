package pki

import (
	"crypto/x509"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
	"github.com/baetyl/baetyl-go/v2/pki"
)

type defaultPkiClient struct {
	cfg       CloudConfig
	pkiClient pki.PKI
}

func init() {
	plugin.RegisterFactory("defaultpki", NewPKI)
}

// NewPKI new
func NewPKI() (plugin.Plugin, error) {
	var cfg CloudConfig
	if err := common.LoadConfig(&cfg); err != nil {
		return nil, err
	}

	sto, err := NewStorage(cfg.PKI.Persistent)
	if err != nil {
		return nil, err
	}
	pkiClient, err := pki.NewPKIClient(cfg.PKI.RootCAKeyFile, cfg.PKI.RootCAFile, sto)
	if err != nil {
		return nil, err
	}
	return &defaultPkiClient{
		cfg:       cfg,
		pkiClient: pkiClient,
	}, nil
}

// root cert
func (p *defaultPkiClient) CreateRootCert(info *x509.CertificateRequest, parentId string) (string, error) {
	return p.pkiClient.CreateRootCert(info, p.cfg.PKI.RootDuration, parentId)
}

func (p *defaultPkiClient) GetRootCert(certId string) ([]byte, error) {
	ca, err := p.pkiClient.GetRootCert(certId)
	if err != nil {
		return nil, err
	}
	return ca.Crt, nil
}

func (p *defaultPkiClient) DeleteRootCert(rootId string) error {
	return p.pkiClient.DeleteRootCert(rootId)
}

// server cert
func (p *defaultPkiClient) CreateServerCert(csr []byte, rootId string) (string, error) {
	return p.pkiClient.CreateSubCert(csr, p.cfg.PKI.SubDuration, rootId)
}

func (p *defaultPkiClient) GetServerCert(certId string) ([]byte, error) {
	return p.pkiClient.GetSubCert(certId)
}

func (p *defaultPkiClient) DeleteServerCert(certId string) error {
	return p.pkiClient.DeleteSubCert(certId)
}

// client cert
func (p *defaultPkiClient) CreateClientCert(csr []byte, rootId string) (string, error) {
	return p.pkiClient.CreateSubCert(csr, p.cfg.PKI.SubDuration, rootId)
}

func (p *defaultPkiClient) GetClientCert(certId string) ([]byte, error) {
	return p.pkiClient.GetSubCert(certId)
}

func (p *defaultPkiClient) DeleteClientCert(certId string) error {
	return p.pkiClient.DeleteSubCert(certId)
}

func (p *defaultPkiClient) Close() error {
	return p.pkiClient.Close()
}
