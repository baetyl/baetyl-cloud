package models

type Activation struct {
	FingerprintValue string            `json:"fingerprintValue,omitempty" db:"fingerprint_value"`
	PenetrateData    map[string]string `json:"penetrateData,omitempty" db:"penetrate_data"`
}

type PackageParam struct {
	Platform string `form:"platform"`
}

type Package struct {
	Url string `json:"url,omitempty"`
	MD5 string `json:"md5,omitempty"`
	Cmd string `json:"cmd,omitempty"`
}
