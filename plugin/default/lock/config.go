package lock

type CloudConfig struct {
	DefaultLocker struct {
		Storage string `yaml:"storage" json:"storage" default:"database"`
	} `yaml:"defaultlocker" json:"defaultlocker"`
}
