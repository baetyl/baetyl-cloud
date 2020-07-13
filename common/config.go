package common

import (
	"github.com/baetyl/baetyl-go/v2/utils"
)

func LoadConfig(cfg interface{}, files ...string) error {
	f := GetConfFile()
	if len(files) > 0 && len(files[0]) > 0 {
		f = files[0]
	}
	return utils.LoadYAML(f, cfg)
}
