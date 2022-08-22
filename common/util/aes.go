package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

func EncryptMap(data map[string][]byte, key []byte) (result map[string][]byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err, _ = r.(error)
		}
	}()
	result = map[string][]byte{}
	for k, v := range data {
		b, err := Encrypt(v, key)
		if err != nil {
			return nil, err
		}
		result[k] = b
	}
	return result, nil
}

func DecryptMap(data map[string][]byte, key []byte) (result map[string][]byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err, _ = r.(error)
		}
	}()
	result = map[string][]byte{}
	for k, v := range data {
		b, err := Decrypt(v, key)
		if err != nil {
			return nil, err
		}
		result[k] = b
	}
	return result, nil
}

func Encrypt(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	plaintext = PKCS7Padding(plaintext, blockSize)
	aesgcm, err := cipher.NewGCMWithNonceSize(block, blockSize)
	if err != nil {
		return nil, err
	}
	crypted := aesgcm.Seal(nil, key[:blockSize], plaintext, nil)
	return crypted, nil
}

func Decrypt(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	aesgcm, err := cipher.NewGCMWithNonceSize(block, blockSize)
	if err != nil {
		return nil, err
	}
	origData, err := aesgcm.Open(nil, key[:blockSize], ciphertext, nil)
	if err != nil {
		return nil, err
	}
	origData = PKCS7UnPadding(origData)
	return origData, nil
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
