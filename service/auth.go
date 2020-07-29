package service

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

//go:generate mockgen -destination=../mock/service/auth.go -package=plugin github.com/baetyl/baetyl-cloud/service AuthService

type AuthService interface {
	Authenticate(c *common.Context) error
	SignToken(meta []byte) ([]byte, error)
	VerifyToken(meta, sign []byte) bool
	GenToken(map[string]interface{}) (string, error)
}

type authService struct {
	plugin.Auth
}

func NewAuthService(config *config.CloudConfig) (AuthService, error) {
	auth, err := plugin.GetPlugin(config.Plugin.Auth)
	if err != nil {
		return nil, err
	}
	return &authService{auth.(plugin.Auth)}, nil
}

func (a *authService) GenToken(data map[string]interface{}) (string, error) {
	signData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	dataStr := hex.EncodeToString(signData)
	sign, err := a.SignToken(signData)
	if err != nil {
		return "", err
	}
	hashed := md5.Sum(sign)
	signStr := hex.EncodeToString(hashed[:])
	return fmt.Sprintf("%s%s", signStr[:10], dataStr), nil
}
