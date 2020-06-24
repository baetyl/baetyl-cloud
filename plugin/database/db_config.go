package database

// CloudConfig baetyl-cloud config
type CloudConfig struct {
	Database struct {
		Type            string `yaml:"type" json:"type" validate:"nonzero"`
		URL             string `yaml:"url" json:"url" validate:"nonzero"`
		MaxConns        int    `yaml:"maxConns" json:"maxConns" default:20`
		MaxIdleConns    int    `yaml:"maxIdleConns" json:"maxIdleConns" default:5`
		ConnMaxLifetime int    `yaml:"connMaxLifetime" json:"connMaxLifetime" default:150`
	} `yaml:"database" json:"database" default:"{}"`
}
