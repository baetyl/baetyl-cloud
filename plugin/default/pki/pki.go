package pki

import (
	"crypto/x509"
	"encoding/base64"
	"errors"
	"io/ioutil"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
	"github.com/baetyl/baetyl-go/v2/pki"
)

const (
	// TypeIssuingCA is a root certificate that can be used to issue sub-certificates
	TypeIssuingCA = "IssuingCA"
	// TypeIssuingSubCert is an issuing sub cert which is signed by issuing ca
	TypeIssuingSubCert = "IssuingSubCertificate"
	// Root cert ID
	RootCertId = "baetyl-cloud-system-cert-root"
)

var (
	ErrParseCert  = errors.New("failed to parse certificate")
	ErrCertInUsed = errors.New("there are also sub-certificates issued according to this certificate in use and cannot be deleted")
	ErrPlugin     = errors.New("plugin type conversion error")
)

type defaultPkiClient struct {
	cfg       CloudConfig
	sto       plugin.PKIStorage
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

	db, err := plugin.GetPlugin(cfg.PKI.Persistent)
	if err != nil {
		return nil, err
	}
	sto, ok := db.(plugin.PKIStorage)
	if !ok {
		return nil, ErrPlugin
	}

	pkiClient, err := pki.NewPKIClient()
	if err != nil {
		return nil, err
	}
	cli := &defaultPkiClient{
		cfg:       cfg,
		sto:       sto,
		pkiClient: pkiClient,
	}
	err = cli.checkRootCA()
	if err != nil {
		return nil, err
	}
	return cli, nil
}

func (p *defaultPkiClient) GetRootCertId() string {
	return RootCertId
}

// root cert
func (p *defaultPkiClient) CreateRootCert(info *x509.CertificateRequest, parentId string) (string, error) {
	if parentId == "" {
		parentId = RootCertId
	}
	parent, err := p.getRootCA(parentId)
	if err != nil {
		return "", err
	}
	cert, err := p.pkiClient.CreateRootCert(info, (int)(p.cfg.PKI.RootDuration.Hours()/24), parent)
	if err != nil {
		return "", err
	}
	certId := common.UUIDPrune()
	err = p.saveCert(certId, cert, []byte(""))
	if err != nil {
		return "", err
	}
	return certId, nil
}

func (p *defaultPkiClient) GetRootCert(certId string) ([]byte, error) {
	ca, err := p.sto.GetCert(certId)
	if err != nil {
		return nil, err
	}
	return base64.StdEncoding.DecodeString(ca.Content)
}

func (p *defaultPkiClient) DeleteRootCert(rootId string) error {
	num, err := p.sto.CountCertByParentId(rootId)
	if err != nil {
		return err
	}
	if num != 0 {
		return ErrCertInUsed
	}
	return p.sto.DeleteCert(rootId)
}

// server cert
func (p *defaultPkiClient) CreateServerCert(csr []byte, rootId string) (string, error) {
	return p.createSubCert(csr, rootId)
}

func (p *defaultPkiClient) GetServerCert(certId string) ([]byte, error) {
	return p.getCert(certId)
}

func (p *defaultPkiClient) DeleteServerCert(certId string) error {
	return p.sto.DeleteCert(certId)
}

// client cert
func (p *defaultPkiClient) CreateClientCert(csr []byte, rootId string) (string, error) {
	return p.createSubCert(csr, rootId)
}

func (p *defaultPkiClient) GetClientCert(certId string) ([]byte, error) {
	return p.getCert(certId)
}

func (p *defaultPkiClient) DeleteClientCert(certId string) error {
	return p.sto.DeleteCert(certId)
}

func (p *defaultPkiClient) Close() error {
	return p.sto.Close()
}

func (p *defaultPkiClient) checkRootCA() error {
	_, err := p.sto.GetCert(RootCertId)
	if err == nil {
		if err := p.sto.DeleteCert(RootCertId); err != nil {
			return err
		}
	}
	crt, err := ioutil.ReadFile(p.cfg.PKI.RootCAFile)
	if err != nil {
		return err
	}
	key, err := ioutil.ReadFile(p.cfg.PKI.RootCAKeyFile)
	if err != nil {
		return err
	}
	return p.saveCert(RootCertId, &pki.CertPem{
		Crt: crt,
		Key: key,
	}, []byte(""))
}

func (p *defaultPkiClient) createSubCert(csr []byte, rootId string) (string, error) {
	parent, err := p.getRootCA(rootId)
	if err != nil {
		return "", err
	}
	crt, err := p.pkiClient.CreateSubCert(csr, (int)(p.cfg.PKI.SubDuration.Hours()/24), parent)
	if err != nil {
		return "", err
	}
	certId := common.UUIDPrune()
	err = p.saveCert(certId, &pki.CertPem{
		Crt: crt,
		Key: []byte(""),
	}, csr)
	if err != nil {
		return "", err
	}
	return certId, nil
}

func (p *defaultPkiClient) getCert(certId string) ([]byte, error) {
	cert, err := p.sto.GetCert(certId)
	if err != nil {
		return nil, err
	}
	return base64.StdEncoding.DecodeString(cert.Content)
}

func (p *defaultPkiClient) saveCert(certId string, cert *pki.CertPem, csr []byte) error {
	crtInfo, err := pki.ParseCertificates(cert.Crt)
	if err != nil {
		return err
	}
	if len(crtInfo) != 1 {
		return ErrParseCert
	}
	tp := TypeIssuingSubCert
	if crtInfo[0].IsCA {
		tp = TypeIssuingCA
	}
	certView := plugin.Cert{
		CertId:     certId,
		Type:       tp,
		CommonName: crtInfo[0].Subject.CommonName,
		Content:    base64.StdEncoding.EncodeToString(cert.Crt),
		PrivateKey: base64.StdEncoding.EncodeToString(cert.Key),
		Csr:        base64.StdEncoding.EncodeToString(csr),
		NotBefore:  crtInfo[0].NotBefore,
		NotAfter:   crtInfo[0].NotAfter,
	}
	return p.sto.CreateCert(certView)
}

func (p *defaultPkiClient) getRootCA(certId string) (*pki.CertPem, error) {
	res, err := p.sto.GetCert(certId)
	if err != nil {
		return nil, err
	}
	crt, err := base64.StdEncoding.DecodeString(res.Content)
	if err != nil {
		return nil, err
	}
	key, err := base64.StdEncoding.DecodeString(res.PrivateKey)
	if err != nil {
		return nil, err
	}
	return &pki.CertPem{
		Crt: crt,
		Key: key,
	}, nil
}
