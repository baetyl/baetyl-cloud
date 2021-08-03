package decryption

type CloudConfig struct {
	Decryption struct {
		Type       string `yaml:"type" json:"type"`
		CipherText string `yaml:"cipherText" json:"cipherText"`
		IV         string `yaml:"iv" json:"iv"`
		Sm4Key     string `yaml:"sm4Key" json:"sm4Key"`
	} `yaml:"decryption" json:"decryption" default:"{}"`
}
