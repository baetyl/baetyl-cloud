package util

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

var (
	ErrRsaPemDecode = errors.New("rsa pem decode error")
	ErrRsaConvert   = errors.New("rsa convert not ok")
)

func GenerateKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privkey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	return privkey, &privkey.PublicKey, nil
}

func PrivateKeyToBytes(priv *rsa.PrivateKey) []byte {
	return pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(priv),
		},
	)
}

func PublicKeyToBytes(pub *rsa.PublicKey) ([]byte, error) {
	pubASN1, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil, err
	}
	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})
	return pubBytes, nil
}

func BytesToPrivateKey(priv []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(priv)
	if block == nil {
		return nil, ErrRsaPemDecode
	}
	b := block.Bytes
	if x509.IsEncryptedPEMBlock(block) {
		res, err := x509.DecryptPEMBlock(block, nil)
		if err != nil {
			return nil, err
		}
		b = res
	}
	key, err := x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func BytesToPublicKey(pub []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pub)
	if block == nil {
		return nil, ErrRsaPemDecode
	}
	b := block.Bytes
	if x509.IsEncryptedPEMBlock(block) {
		res, err := x509.DecryptPEMBlock(block, nil)
		if err != nil {
			return nil, err
		}
		b = res
	}
	ifc, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		return nil, err
	}
	key, ok := ifc.(*rsa.PublicKey)
	if !ok {
		return nil, ErrRsaConvert
	}
	return key, nil
}

func SignPKCS1v15(text []byte, priv *rsa.PrivateKey) ([]byte, error) {
	rng := rand.Reader
	hashed := sha512.Sum512(text)
	return rsa.SignPKCS1v15(rng, priv, crypto.SHA512, hashed[:])
}

func VerifyPKCS1v15(text, signature []byte, pub *rsa.PublicKey) bool {
	hashed := sha512.Sum512(text)
	return nil == rsa.VerifyPKCS1v15(pub, crypto.SHA512, hashed[:], signature)
}
