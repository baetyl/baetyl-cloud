package csrf

type CloudConfig struct {
	CsrfConfig struct {
		CookieName string   `yaml:"cookieName" json:"cookieName" default:"csrftoken"`
		HeaderName string   `yaml:"headerName" json:"headerName" default:"csrftoken"`
		Whitelist  []string `yaml:"whitelist" json:"whitelist"`
	} `yaml:"defaultcsrf" json:"defaultcsrf"`
}
