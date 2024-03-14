package decryption

import (
	"encoding/hex"

	"github.com/ZZMarquis/gm/sm4"
	"github.com/ZZMarquis/gm/util"
	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/baetyl/baetyl-go/v2/log"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
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
		url, err = d.sm4Decryption(cipherText, d.cfg.Decryption.Sm4EncKey, d.cfg.Decryption.IV, d.cfg.Decryption.Sm4ProtectKey)
	default:
		url = cipherText
	}
	return url, errors.Trace(err)
}

// 注意：sm4解密需提供密文、IV分量及密钥保护分量，密钥加密分量，以hex string形式提供
func (d *Decryption) sm4Decryption(cipherText, sm4EncKey, iv string, sm4ProtectKey []string) (string, error) {
	if cipherText == "" || sm4EncKey == "" || iv == "" || sm4ProtectKey == nil {
		d.log.Error("failed to decrypt, not enough param provided")
		return "", errors.Trace(errors.New("failed to decrypt, not enough param provided"))
	}
	decCipherText, err := hex.DecodeString(cipherText)
	if err != nil {
		d.log.Error("failed to decode cipherText", log.Any("cipherText", cipherText), log.Any("error", err))
		return "", errors.Trace(err)
	}

	decIV, err := hex.DecodeString(iv)
	if err != nil {
		d.log.Error("failed to decode iv", log.Any("iv", iv), log.Any("error", err))
		return "", errors.Trace(err)
	}
	var proKey []byte
	for i, p := range sm4ProtectKey {
		decP, err := hex.DecodeString(p)
		if err != nil {
			d.log.Error("failed to decode sm4ProtectKey", log.Any("sm4ProtectKey", sm4ProtectKey), log.Any("error", err))
			return "", errors.Trace(err)
		}
		if i == 0 {
			proKey = decP
		} else {
			proKey = XORBytes(proKey, decP)
			if proKey == nil {
				d.log.Error("failed to xor sm4ProtectKey", log.Any("error", err))
				return "", errors.Trace(err)
			}
		}
	}
	encKey, err := hex.DecodeString(sm4EncKey)
	if err != nil {
		d.log.Error("failed to decode sm4EncKey", log.Any("sm4Key", sm4EncKey), log.Any("error", err))
		return "", errors.Trace(err)
	}

	decKey, err := sm4.ECBDecrypt(proKey, encKey)
	if err != nil {
		d.log.Error("failed to decode proKey", log.Any("proKey", proKey), log.Any("error", err))
		return "", errors.Trace(err)
	}

	plaintext, err := sm4.CBCDecrypt(decKey, decIV, decCipherText)
	if err != nil {
		d.log.Error("failed to decrypt cipherText", log.Any("decCipherText", decCipherText), log.Any("error", err))
		return "", errors.Trace(err)
	}
	return string(util.PKCS5UnPadding(plaintext)), nil
}

func XORBytes(a, b []byte) []byte {
	if len(a) != len(b) {
		return nil
	}
	buf := make([]byte, len(a))
	for i := range a {
		buf[i] = a[i] ^ b[i]
	}
	return buf
}

func (d *Decryption) Close() (err error) {
	return nil
}
