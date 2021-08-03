package plugin

import "io"

//go:generate mockgen -destination=../mock/plugin/decryption.go -package=plugin github.com/baetyl/baetyl-cloud/v2/plugin Decrypt

type Decrypt interface {
	Decrypt(cipherText string) (string, error)

	io.Closer
}
