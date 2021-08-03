package decryption

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func MockNewDP() (*Decryption, error) {
	var cfg CloudConfig
	cfg.Decryption.Type = "sm4"
	cfg.Decryption.Sm4Key = "d58d8effe6d50d8e5a28fa142c5c012b"

	return &Decryption{cfg: cfg}, nil
}

func TestDecryp(t *testing.T) {
	dp, err := MockNewDP()
	assert.NoError(t, err)

	decode, err := dp.Decrypt("c43b88c3b5ae15eee7b99dd9188c5ebd45be72bb74a22091eafc3beca75a47a5")
	assert.NoError(t, err)
	assert.Equal(t, decode, "SymmetricCipher Test data222")
}
