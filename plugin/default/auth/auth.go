package auth

import (
	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

type defaultAuth struct {
	cfg CloudConfig
}

func init() {
	plugin.RegisterFactory("defaultauth", New)
}

// New New
func New() (plugin.Plugin, error) {
	var cfg CloudConfig
	if err := common.LoadConfig(&cfg); err != nil {
		return nil, err
	}
	return &defaultAuth{
		cfg: cfg,
	}, nil
}

func (d *defaultAuth) Authenticate(c *common.Context) error {
	c.SetNamespace(d.cfg.DefaultAuth.Namespace)
	c.SetUserInfo(common.UserInfo{
		User:   common.User{ID: d.cfg.DefaultAuth.Namespace, Name: d.cfg.DefaultAuth.Namespace},
		Domain: common.Domain{ID: d.cfg.DefaultAuth.Namespace, Name: d.cfg.DefaultAuth.Namespace},
	})
	return nil
}

func (d *defaultAuth) AuthAndVerify(c *common.Context, pr *plugin.PermissionRequest) error {
	return d.Authenticate(c)
}

func (d *defaultAuth) Verify(c *common.Context, pr *plugin.PermissionRequest) error {
	return nil
}

// Close Close
func (d *defaultAuth) Close() error {
	return nil
}
