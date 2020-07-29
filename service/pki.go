package service

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"github.com/baetyl/baetyl-go/v2/pki"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

//go:generate mockgen -destination=../mock/service/pki.go -package=plugin github.com/baetyl/baetyl-cloud/service PKIService

type PKIService interface {
	// GetCA get ca
	GetCA() ([]byte, error)
	// SignServerCertificate sign a certificate which can be used to connect to cloud
	SignServerCertificate(cn string, altNames models.AltNames) (*models.PEMCredential, error)
	// SignNodeCertificate sign a certificate which can be used to connect to cloud
	SignClientCertificate(cn string, altNames models.AltNames) (*models.PEMCredential, error)
	// DeleteServerCertificate delete a server certificate by certId
	DeleteServerCertificate(certId string) error
	// DeleteClientCertificate delete a server certificate by certId
	DeleteClientCertificate(certId string) error
}

const (
	Certificate = "certificate"
	CertRoot    = "baetyl.ca"
)

type pkiService struct {
	pki plugin.PKI
	db  plugin.DBStorage
}

// NewPKIService create a certificate service
func NewPKIService(config *config.CloudConfig) (PKIService, error) {
	pk, err := plugin.GetPlugin(config.Plugin.PKI)
	if err != nil {
		return nil, err
	}

	ds, err := plugin.GetPlugin(config.Plugin.DatabaseStorage)
	if err != nil {
		return nil, err
	}

	p := &pkiService{
		pki: pk.(plugin.PKI),
		db:  ds.(plugin.DBStorage),
	}

	return p, nil
}

func (p *pkiService) GetCA() ([]byte, error) {
	// determine whether it exists in the system configuration
	sysConf, err := p.db.GetSysConfig(Certificate, CertRoot)
	if err == nil {
		return p.pki.GetRootCert(sysConf.Value)
	}
	// if it does not exist in the system configuration, apply to the pki service to create a root certificate
	certInfo := p.genDefaultCSR("baetyl.ca")
	rootId, err := p.pki.CreateRootCert(certInfo, "")
	if err != nil {
		return nil, err
	}
	// save the applied root certificate to the system configuration
	rootCertConfig := &models.SysConfig{
		Type:  Certificate,
		Key:   CertRoot,
		Value: rootId,
	}
	_, err = p.db.CreateSysConfig(rootCertConfig)
	if err != nil {
		return nil, err
	}
	return p.pki.GetRootCert(rootId)
}

func (p *pkiService) DeleteServerCertificate(certId string) error {
	return p.pki.DeleteServerCert(certId)
}

func (p *pkiService) DeleteClientCertificate(certId string) error {
	return p.pki.DeleteClientCert(certId)
}

func (p *pkiService) SignServerCertificate(cn string, altNames models.AltNames) (*models.PEMCredential, error) {
	return p.signCertificate(cn, altNames, p.pki.CreateServerCert, p.pki.GetServerCert)
}

func (p *pkiService) SignClientCertificate(cn string, altNames models.AltNames) (*models.PEMCredential, error) {
	return p.signCertificate(cn, altNames, p.pki.CreateClientCert, p.pki.GetClientCert)
}

func (p *pkiService) signCertificate(cn string, altNames models.AltNames, create func(csr []byte, rootId string) (string, error), get func(certId string) ([]byte, error)) (*models.PEMCredential, error) {
	csrInfo := p.genDefaultCSR(cn)
	csrInfo.DNSNames = altNames.DNSNames
	csrInfo.EmailAddresses = altNames.Emails
	csrInfo.IPAddresses = altNames.IPs
	csrInfo.URIs = altNames.URIs

	priv, err := pki.GenCertPrivateKey(pki.DefaultDSA, pki.DefaultRSABits)
	if err != nil {
		return nil, err
	}

	keyPem, err := pki.EncodeCertPrivateKey(priv)
	if err != nil {
		return nil, err
	}

	csr, err := x509.CreateCertificateRequest(rand.Reader, csrInfo, priv.Key)
	if err != nil {
		return nil, err
	}

	sysConf, _ := p.db.GetSysConfig(Certificate, CertRoot)
	if sysConf == nil {
		return nil, common.Error(
			common.ErrResourceNotFound,
			common.Field("type", Certificate),
			common.Field("name", CertRoot))
	}

	certId, err := create(csr, sysConf.Value)
	if err != nil {
		return nil, err
	}
	certPem, err := get(certId)
	if err != nil {
		return nil, err
	}

	return &models.PEMCredential{
		CertPEM: certPem,
		KeyPEM:  keyPem,
		CertId:  certId,
	}, nil
}

func (p *pkiService) genDefaultCSR(cn string) *x509.CertificateRequest {
	return &x509.CertificateRequest{
		Subject: pkix.Name{
			Country:            []string{"CN"},
			Organization:       []string{"Linux Foundation Edge"},
			OrganizationalUnit: []string{"BAETYL"},
			Locality:           []string{"Haidian District"},
			Province:           []string{"Beijing"},
			StreetAddress:      []string{"Baidu Campus"},
			PostalCode:         []string{"100093"},
			CommonName:         cn,
		},
		EmailAddresses: []string{"baetyl@lists.lfedge.org"},
	}
}
