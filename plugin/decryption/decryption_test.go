package decryption

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func MockNewDP() (*Decryption, error) {
	var cfg CloudConfig
	cfg.Decryption.Type = "sm4"
	cfg.Decryption.Sm4EncKey = "AEF55CE747A502A9C571A1AC6DACE2EA"
	cfg.Decryption.IV = "00000000000000000000000000000000"
	cfg.Decryption.Sm4ProtectKey = make([]string, 3)
	cfg.Decryption.Sm4ProtectKey[0] = "e6af2c278c595cc8a00d4be45fdb0c40"
	cfg.Decryption.Sm4ProtectKey[1] = "927eafc3e53dc9898b2da46a6d6be336"
	cfg.Decryption.Sm4ProtectKey[2] = "e909ee2ddb3b31a12ab4b04557e0494d"

	return &Decryption{cfg: cfg}, nil
}

func TestDecryp(t *testing.T) {
	dp, err := MockNewDP()
	assert.NoError(t, err)

	decode, err := dp.Decrypt("c43b88c3b5ae15eee7b99dd9188c5ebd45be72bb74a22091eafc3beca75a47a5")
	assert.NoError(t, err)
	assert.Equal(t, decode, "SymmetricCipher Test data222")

	decode2, err := dp.Decrypt("ecc28664636fa967c8e87a9b4ba0581f")
	assert.NoError(t, err)
	assert.Equal(t, decode2, "secretpassword")
}
