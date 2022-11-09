package license

import (
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

type license struct {
}

func init() {
	plugin.RegisterFactory("defaultlicense", New)
}

func New() (plugin.Plugin, error) {
	return &license{}, nil
}

var _ plugin.License = &license{}

func (l *license) ProtectCode() error {
	return nil
}

func (l *license) CheckLicense() error {
	return nil
}

func (l *license) Close() error {
	return nil
}
