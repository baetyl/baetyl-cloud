package sign

import (
	"crypto/rsa"
	"io/ioutil"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/common/util"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

type rsaSign struct {
	cfg  CloudConfig
	priv *rsa.PrivateKey
}

func init() {
	plugin.RegisterFactory("rsasign", New)
}

// New New
func New() (plugin.Plugin, error) {
	var cfg CloudConfig
	if err := common.LoadConfig(&cfg); err != nil {
		return nil, err
	}
	key, err := ioutil.ReadFile(cfg.RSASign.KeyFile)
	if err != nil {
		return nil, err
	}
	priv, err := util.BytesToPrivateKey(key)
	if err != nil {
		return nil, err
	}
	return &rsaSign{
		cfg:  cfg,
		priv: priv,
	}, nil
}

func (r *rsaSign) Signature(meta []byte) ([]byte, error) {
	return util.SignPKCS1v15(meta, r.priv)
}

func (r *rsaSign) Verify(meta, sign []byte) bool {
	return util.VerifyPKCS1v15(meta, sign, &r.priv.PublicKey)
}

// Close Close
func (r *rsaSign) Close() error {
	return nil
}
