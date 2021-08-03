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
	Log *log.Logger
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
		Log: log.With(log.Any("plugin", "decryption")),
	}, nil
}

func (d *Decryption) Decrypt(cipherText string) (string, error) {
	var url string
	var err error
	switch d.cfg.Decryption.Type {
	case "sm4":
		url, err = sm4Decryption(cipherText, d.cfg.Decryption.Sm4Key)
	default:
		url = cipherText
	}
	return url, err
}

// 注意：sm4解密需提供密文及密钥，以hex string形式提供
func sm4Decryption(cipherText, sm4Key string) (string, error) {
	decCipherText, err := hex.DecodeString(cipherText)
	if err != nil {
		log.L().Error("fail to decode cipherText", log.Any("cipherText", cipherText), log.Any("error", err))
		return "", err
	}
	decKey, err := hex.DecodeString(sm4Key)
	if err != nil {
		log.L().Error("fail to decode sm4Key", log.Any("sm4Key", sm4Key), log.Any("error", err))
		return "", err
	}

	iv := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	_ = sm4.SetIV(iv)

	plaintext, err := sm4.Sm4Cbc(decKey, decCipherText, false)
	if err != nil {
		log.L().Error("fail to decrypt cipherText", log.Any("decCipherText", decCipherText), log.Any("error", err))
		return "", err
	}
	return string(plaintext), nil
}

func (d *Decryption) Close() (err error) {
	return nil
}
