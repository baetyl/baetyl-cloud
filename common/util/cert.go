package util

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/baetyl/baetyl-cloud/models"
)

const (
	RsaPrivateKeyBlockType      = "RSA PRIVATE KEY"
	EcPrivateKeyBlockType       = "EC PRIVATE KEY"
	CertificateBlockType        = "CERTIFICATE"
	CertificateRequestBlockType = "CERTIFICATE REQUEST"

	DefaultDSA     = "P256"
	DefaultRSABits = 2048
)

func GenCertPrivateKey(dsa string, bits int) (*models.PrivateKey, error) {
	var err error
	var key interface{}
	var priv *models.PrivateKey
	switch dsa {
	case "rsa":
		key, err = rsa.GenerateKey(rand.Reader, bits)
		if err != nil {
			return nil, err
		}
		priv = &models.PrivateKey{
			Key:  key,
			Type: RsaPrivateKeyBlockType,
		}
	case "P224":
		key, err = ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
		if err != nil {
			return nil, err
		}
		priv = &models.PrivateKey{
			Key:  key,
			Type: EcPrivateKeyBlockType,
		}
	case "P256":
		key, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return nil, err
		}
		priv = &models.PrivateKey{
			Key:  key,
			Type: EcPrivateKeyBlockType,
		}
	case "P384":
		key, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
		if err != nil {
			return nil, err
		}
		priv = &models.PrivateKey{
			Key:  key,
			Type: EcPrivateKeyBlockType,
		}
	case "P521":
		key, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
		if err != nil {
			return nil, err
		}
		priv = &models.PrivateKey{
			Key:  key,
			Type: EcPrivateKeyBlockType,
		}
	default:
		return nil, fmt.Errorf("unRecognized digital signature algorithm: %s", dsa)
	}
	return priv, nil
}

// EncodeCertPrivateKey returns PEM-encoded private key data
func EncodeCertPrivateKey(priv *models.PrivateKey) ([]byte, error) {
	switch priv.Type {
	case RsaPrivateKeyBlockType:
		block := pem.Block{
			Type:  RsaPrivateKeyBlockType,
			Bytes: x509.MarshalPKCS1PrivateKey(priv.Key.(*rsa.PrivateKey)),
		}
		return pem.EncodeToMemory(&block), nil
	case EcPrivateKeyBlockType:
		bytes, err := x509.MarshalECPrivateKey(priv.Key.(*ecdsa.PrivateKey))
		if err != nil {
			return nil, fmt.Errorf("error mashaling ECPrivateKey: %s", err.Error())
		}
		block := pem.Block{
			Type:  EcPrivateKeyBlockType,
			Bytes: bytes,
		}
		return pem.EncodeToMemory(&block), nil
	default:
		return nil, fmt.Errorf("unRecognized type of PrivateKey")
	}
}

// ParseCertPrivateKey takes a key PEM byte array and returns a PrivateKey that represents
// Either an RSA or EC private key.
func ParseCertPrivateKey(key []byte) (*models.PrivateKey, error) {
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, errors.New("key is not PEM encoded")
	}
	switch block.Type {
	case EcPrivateKeyBlockType:
		k, err := x509.ParseECPrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		return &models.PrivateKey{Type: EcPrivateKeyBlockType, Key: k}, nil
	case RsaPrivateKeyBlockType:
		k, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		return &models.PrivateKey{Type: RsaPrivateKeyBlockType, Key: k}, nil
	default:
		return nil, fmt.Errorf("unSupported block type: %s", block.Type)
	}
}

// ParseCertificates takes a PEM encoded x509 certificates byte array and returns A x509 certificate and the block byte array
func ParseCertificates(pemCerts []byte) ([]*x509.Certificate, error) {
	ok := false
	certs := []*x509.Certificate{}
	for len(pemCerts) > 0 {
		var block *pem.Block
		block, pemCerts = pem.Decode(pemCerts)
		if block == nil {
			break
		}
		// Only use PEM "CERTIFICATE" blocks without extra headers
		if block.Type != CertificateBlockType || len(block.Headers) != 0 {
			continue
		}
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return certs, err
		}
		certs = append(certs, cert)
		ok = true
	}
	if !ok {
		return certs, errors.New("data does not contain any valid RSA or ECDSA certificates")
	}
	return certs, nil
}

func SigAlgorithmType(priv *models.PrivateKey) x509.SignatureAlgorithm {
	switch priv.Type {
	case RsaPrivateKeyBlockType:
		keySize := priv.Key.(*rsa.PrivateKey).N.BitLen()
		switch {
		case keySize >= 4096:
			return x509.SHA512WithRSA
		case keySize >= 3072:
			return x509.SHA384WithRSA
		default:
			return x509.SHA256WithRSA
		}
	case EcPrivateKeyBlockType:
		keySize := priv.Key.(*ecdsa.PrivateKey).D.BitLen()/8 + 8
		switch {
		case keySize >= 72:
			return x509.ECDSAWithSHA512
		case keySize >= 54:
			return x509.ECDSAWithSHA384
		default:
			return x509.ECDSAWithSHA256
		}
	default:
		return x509.UnknownSignatureAlgorithm
	}
}

// EncodeCertificates returns the PEM-encoded byte array that represents by the specified certs
func EncodeCertificates(certs ...*x509.Certificate) ([]byte, error) {
	b := bytes.Buffer{}
	for _, cert := range certs {
		if err := pem.Encode(&b, &pem.Block{Type: CertificateBlockType, Bytes: cert.Raw}); err != nil {
			return []byte{}, err
		}
	}
	return b.Bytes(), nil
}

// EncodeCertificatesRequest returns the PEM-encoded byte array that represents by the specified certs
func EncodeCertificatesRequest(csrs ...*x509.CertificateRequest) ([]byte, error) {
	b := bytes.Buffer{}
	for _, csr := range csrs {
		if err := pem.Encode(&b, &pem.Block{Type: CertificateRequestBlockType, Bytes: csr.Raw}); err != nil {
			return []byte{}, err
		}
	}
	return b.Bytes(), nil
}

func EncodeByteToPem(data []byte, tp string) string {
	src := base64.StdEncoding.EncodeToString(data)
	res := "-----BEGIN " + tp + "-----\n"
	for i := 0; i < len(src)/64; i++ {
		max := (i + 1) * 64
		if len(src) < (i+1)*64 {
			max = len(src)
		}
		res += src[i*64:max] + "\n"
	}
	if len(src)%64 != 0 {
		res += src[64*(len(src)/64):] + "\n"
	}
	res += "-----END " + tp + "-----\n"
	return res
}
