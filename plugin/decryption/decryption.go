package decryption

import (
	"encoding/hex"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/baetyl/baetyl-go/v2/log"
	"github.com/tjfoc/gmsm/sm4"
)

type Decryption struct {
	cfg CloudConfig
	log *log.Logger
}

func init() {
	plugin.RegisterFactory("decryption", New)
}

func New() (plugin.Plugin, error) {
	var cfg CloudConfig
	if err := common.LoadConfig(&cfg); err != nil {
		return nil, errors.Trace(err)
	}
	return &Decryption{
		cfg: cfg,
		log: log.With(log.Any("plugin", "decryption")),
	}, nil
}

func (d *Decryption) Decrypt(cipherText string) (string, error) {
	var url string
	var err error
	switch d.cfg.Decryption.Type {
	case "sm4":
		url, err = d.sm4Decryption(cipherText, d.cfg.Decryption.Sm4Key, d.cfg.Decryption.IV)
	default:
		url = cipherText
	}
	return url, errors.Trace(err)
}

// 注意：sm4解密需提供密文及密钥，以hex string形式提供
func (d *Decryption) sm4Decryption(cipherText, sm4Key, iv string) (string, error) {
	decCipherText, err := hex.DecodeString(cipherText)
	if err != nil {
		d.log.Error("fail to decode cipherText", log.Any("cipherText", cipherText), log.Any("error", err))
		return "", errors.Trace(err)
	}
	decKey, err := hex.DecodeString(sm4Key)
	if err != nil {
		d.log.Error("fail to decode sm4Key", log.Any("sm4Key", sm4Key), log.Any("error", err))
		return "", errors.Trace(err)
	}

	decIV, err := hex.DecodeString(iv)
	_ = sm4.SetIV(decIV)

	plaintext, err := sm4.Sm4Cbc(decKey, decCipherText, false)
	if err != nil {
		d.log.Error("fail to decrypt cipherText", log.Any("decCipherText", decCipherText), log.Any("error", err))
		return "", errors.Trace(err)
	}
	return string(plaintext), nil
}

func (d *Decryption) Close() (err error) {
	return nil
}
