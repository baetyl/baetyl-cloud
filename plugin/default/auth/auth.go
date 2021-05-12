package auth

import (
	"crypto/rsa"
	"io/ioutil"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/common/util"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

type defaultAuth struct {
	cfg  CloudConfig
	priv *rsa.PrivateKey
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
	key, err := ioutil.ReadFile(cfg.DefaultAuth.KeyFile)
	if err != nil {
		return nil, err
	}
	priv, err := util.BytesToPrivateKey(key)
	if err != nil {
		return nil, err
	}
	return &defaultAuth{
		cfg:  cfg,
		priv: priv,
	}, nil
}

func (d *defaultAuth) Authenticate(c *common.Context) error {
	c.SetNamespace(d.cfg.DefaultAuth.Namespace)
	return nil
}

func (d *defaultAuth) SignToken(meta []byte) ([]byte, error) {
	return util.SignPKCS1v15(meta, d.priv)
}

func (d *defaultAuth) VerifyToken(meta, sign []byte) bool {
	return util.VerifyPKCS1v15(meta, sign, &d.priv.PublicKey)
}

func (d *defaultAuth) AuthAndVerify(c *common.Context, pr *plugin.PermissionRequest) error {
	return nil
}

func (d *defaultAuth) Verify(c *common.Context, pr *plugin.PermissionRequest) error {
	return nil
}

// Close Close
func (d *defaultAuth) Close() error {
	return nil
}
