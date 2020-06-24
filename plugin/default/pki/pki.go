package pki

import (
	"crypto/rand"
	"crypto/x509"
	"github.com/baetyl/baetyl-cloud/common"
	"github.com/baetyl/baetyl-cloud/common/util"
	"github.com/baetyl/baetyl-cloud/models"
	"github.com/baetyl/baetyl-cloud/plugin"
	"io/ioutil"
	"math/big"
	"time"
)

const (
	// TypeRootCA root ca, more details please refer to: https://www.cnblogs.com/sunsky303/p/11194801.html
	TypeRootCA = "RootCA"
	// TypeClientCert client which is signed by root ca
	TypeClientCert = "ClientCertificate"
	// TypeIntermediateCA is an intermediate ca which is signed by root ca
	TypeIntermediateCA = "IntermediateCA"
	// TypeIntermediateClientCert is a client cert which is signed by intermediate ca
	TypeIntermediateClientCert = "IntermediateClientCertificate"
	// TypeIntermediateServerCert is a server cert which is signed by intermediate ca
	TypeIntermediateServerCert = "IntermediateServerCertificate"
	// TypeIssuingCA is an issuing ca which is signed by intermediate ca and can be used for signing certificate for module of current node
	TypeIssuingCA = "IssuingCA"
	// TypeIssuingClientCert is an issuing client cert which is signed by issuing ca
	TypeIssuingClientCert = "IssuingClientCertificate"
	// TypeIssuingServerCert is an issuing server cert which is signed by issuing ca
	TypeIssuingServerCert = "IssuingServerCertificate"
)

type defaultPkiClient struct {
	cfg   CloudConfig
	caKey []byte
	caPem []byte
	pvc   PVC
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

	key, err := ioutil.ReadFile(cfg.PKI.RootCAKeyFile)
	if err != nil {
		return nil, err
	}
	pem, err := ioutil.ReadFile(cfg.PKI.RootCAFile)
	if err != nil {
		return nil, err
	}
	pvc, err := NewPVC(cfg.PKI.Persistent)
	if err != nil {
		return nil, err
	}
	return &defaultPkiClient{
		cfg:   cfg,
		caKey: key,
		caPem: pem,
		pvc:   pvc,
	}, nil
}

// root cert
func (p *defaultPkiClient) CreateRootCert(info *x509.CertificateRequest, parentId string) (string, error) {
	// get parent cert
	var caKeyByte []byte
	var caPemByte []byte
	if len(parentId) == 0 {
		caKeyByte = p.caKey
		caPemByte = p.caPem
	} else {
		parentCert, err := p.pvc.GetCert(parentId)
		if err != nil {
			return "", err
		}
		priv, err := util.GenCertPrivateKey(util.DefaultDSA, util.DefaultRSABits)
		if err != nil {
			return "", err
		}
		caKeyByte, err = util.EncodeCertPrivateKey(priv)
		if err != nil {
			return "", err
		}
		caPemByte = parentCert.Content
	}

	// generate cert
	caKey, err := util.ParseCertPrivateKey(caKeyByte)
	if err != nil {
		return "", err
	}
	caCert, err := util.ParseCertificates(caPemByte)
	if err != nil {
		return "", err
	}

	csr, err := x509.CreateCertificateRequest(rand.Reader, info, caKey.Key)
	if err != nil {
		return "", err
	}

	csrInfo, err := x509.ParseCertificateRequest(csr)
	if err != nil {
		return "", err
	}

	keyUsage := x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign

	certInfo := &x509.Certificate{
		IsCA:                  true,
		Subject:               info.Subject,
		SerialNumber:          big.NewInt(time.Now().UnixNano()),
		NotBefore:             time.Now().UTC(),
		NotAfter:              time.Now().AddDate(0, 0, plugin.DefaultRootDuration).UTC(),
		EmailAddresses:        info.EmailAddresses,
		IPAddresses:           info.IPAddresses,
		URIs:                  info.URIs,
		DNSNames:              info.DNSNames,
		BasicConstraintsValid: true,
		SignatureAlgorithm:    util.SigAlgorithmType(caKey),
		KeyUsage:              keyUsage,
	}

	cert, err := x509.CreateCertificate(rand.Reader, certInfo, caCert[0], csrInfo.PublicKey, caKey.Key)
	if err != nil {
		return "", err
	}

	// save cert
	certView := models.Cert{
		CertId:     common.UUIDPrune(),
		ParentId:   parentId,
		Type:       TypeIssuingCA,
		CommonName: info.Subject.CommonName,
		Csr:        []byte(util.EncodeByteToPem(csr, util.CertificateRequestBlockType)),
		Content:    []byte(util.EncodeByteToPem(cert, util.CertificateBlockType)),
		Priv:       caKeyByte,
		NotBefore:  certInfo.NotBefore,
		NotAfter:   certInfo.NotAfter,
	}
	err = p.pvc.CreateCert(certView)
	if err != nil {
		return "", err
	}

	return certView.CertId, nil
}

func (p *defaultPkiClient) GetRootCert(rootId string) ([]byte, error) {
	cert, err := p.pvc.GetCert(rootId)
	if err != nil {
		return nil, err
	}
	return cert.Content, nil
}

func (p *defaultPkiClient) DeleteRootCert(rootId string) error {
	count, err := p.pvc.CountCertByParentId(rootId)
	if err != nil {
		return err
	}
	if count > 0 {
		return common.Error(common.ErrResourceHasBeenUsed,
			common.Field("type", "certificate"),
			common.Field("name", rootId))
	}
	return p.pvc.DeleteCert(rootId)
}

// server cert
func (p *defaultPkiClient) CreateServerCert(csr []byte, rootId string) (string, error) {
	keyUsage := x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment
	extKeyUsage := []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
	return p.createSubCert(csr, rootId, keyUsage, extKeyUsage, true)
}

func (p *defaultPkiClient) GetServerCert(certId string) ([]byte, error) {
	cert, err := p.pvc.GetCert(certId)
	if err != nil {
		return nil, err
	}
	return cert.Content, nil
}

func (p *defaultPkiClient) DeleteServerCert(certId string) error {
	return p.pvc.DeleteCert(certId)
}

// client cert
func (p *defaultPkiClient) CreateClientCert(csr []byte, rootId string) (string, error) {
	keyUsage := x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment
	extKeyUsage := []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}
	return p.createSubCert(csr, rootId, keyUsage, extKeyUsage, false)
}

func (p *defaultPkiClient) GetClientCert(certId string) ([]byte, error) {
	cert, err := p.pvc.GetCert(certId)
	if err != nil {
		return nil, err
	}
	return cert.Content, nil
}

func (p *defaultPkiClient) DeleteClientCert(certId string) error {
	return p.pvc.DeleteCert(certId)
}

func (p *defaultPkiClient) Close() error {
	return p.pvc.Close()
}

func (p *defaultPkiClient) createSubCert(csr []byte, rootId string,
	keyUsage x509.KeyUsage, extKeyUsage []x509.ExtKeyUsage, isServer bool) (string, error) {
	// get ca cert
	ca, err := p.pvc.GetCert(rootId)
	if err != nil {
		return "", common.Error(common.ErrDatabase, common.Field("error", err))
	}
	if ca == nil {
		return "", common.Error(common.ErrResourceNotFound,
			common.Field("type", "certificate"),
			common.Field("name", rootId))
	}

	// parse ca cert
	caKey, err := util.ParseCertPrivateKey(ca.Priv)
	if err != nil {
		return "", err
	}
	caCert, err := util.ParseCertificates(ca.Content)
	if err != nil {
		return "", err
	}

	// create server data
	csrInfo, err := x509.ParseCertificateRequest(csr)
	if err != nil {
		return "", err
	}

	certInfo := &x509.Certificate{
		IsCA:                  false,
		SerialNumber:          big.NewInt(time.Now().UnixNano()),
		Subject:               csrInfo.Subject,
		NotBefore:             time.Now().UTC(),
		NotAfter:              caCert[0].NotAfter,
		EmailAddresses:        csrInfo.EmailAddresses,
		IPAddresses:           csrInfo.IPAddresses,
		URIs:                  csrInfo.URIs,
		DNSNames:              csrInfo.DNSNames,
		BasicConstraintsValid: true,
		SignatureAlgorithm:    util.SigAlgorithmType(caKey),
		KeyUsage:              keyUsage,
		ExtKeyUsage:           extKeyUsage,
	}

	cert, err := x509.CreateCertificate(rand.Reader, certInfo, caCert[0], csrInfo.PublicKey, caKey.Key)
	if err != nil {
		return "", err
	}

	// save cert
	certView := models.Cert{
		CertId:     common.UUIDPrune(),
		ParentId:   rootId,
		Type:       TypeIssuingClientCert,
		CommonName: certInfo.Subject.CommonName,
		Csr:        []byte(util.EncodeByteToPem(csr, util.CertificateRequestBlockType)),
		Content:    []byte(util.EncodeByteToPem(cert, util.CertificateBlockType)),
		NotBefore:  certInfo.NotBefore,
		NotAfter:   certInfo.NotAfter,
	}
	if isServer {
		certView.Type = TypeIssuingServerCert
	}
	err = p.pvc.CreateCert(certView)
	if err != nil {
		return "", err
	}

	return certView.CertId, nil
}
