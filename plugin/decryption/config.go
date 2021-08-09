package decryption

type CloudConfig struct {
	// CipherText 密码密文  IV sm4 cbc解密iv分量  Sm4ProtectKey sm4解密密钥保护密钥分量 Sm4EncKey sm4解密密钥加密分量
	Decryption struct {
		Type          string   `yaml:"type" json:"type"`
		CipherText    string   `yaml:"cipherText" json:"cipherText"`
		IV            string   `yaml:"iv" json:"iv"`
		Sm4ProtectKey []string `yaml:"sm4ProtectKey" json:"sm4ProtectKey"`
		Sm4EncKey     string   `yaml:"sm4EncKey" json:"sm4EncKey"`
	} `yaml:"decryption" json:"decryption" default:"{}"`
}
