package sign

import (
	"bytes"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

type defaultSign struct {
	cfg CloudConfig
}

func init() {
	plugin.RegisterFactory("defaultsign", New)
}

// New New
func New() (plugin.Plugin, error) {
	var cfg CloudConfig
	if err := common.LoadConfig(&cfg); err != nil {
		return nil, err
	}
	return &defaultSign{
		cfg: cfg,
	}, nil
}

func (d *defaultSign) Signature(meta []byte) ([]byte, error) {
	return meta, nil
}

func (d *defaultSign) Verify(meta, sign []byte) bool {
	return bytes.Equal(meta, sign)
}

// Close Close
func (d *defaultSign) Close() error {
	return nil
}
