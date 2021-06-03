package sign

type CloudConfig struct {
	RSASign struct {
		KeyFile string `yaml:"keyFile" json:"keyFile" default:"etc/baetyl/token.key"`
	} `yaml:"rsasign" json:"rsasign"`
}
