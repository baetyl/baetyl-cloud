package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/baetyl/baetyl-go/v2/json"

	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

//go:generate mockgen -destination=../mock/service/sign.go -package=service github.com/baetyl/baetyl-cloud/v2/service SignService

type SignService interface {
	plugin.Sign
	GenToken(map[string]interface{}) (string, error)
}

type signService struct {
	plugin.Sign
}

func NewSignService(config *config.CloudConfig) (SignService, error) {
	s, err := plugin.GetPlugin(config.Plugin.Sign)
	if err != nil {
		return nil, err
	}
	return &signService{s.(plugin.Sign)}, nil
}

func (s *signService) GenToken(data map[string]interface{}) (string, error) {
	signData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	dataStr := hex.EncodeToString(signData)
	sign, err := s.Signature(signData)
	if err != nil {
		return "", err
	}
	hashed := sha256.Sum256(sign)
	signStr := hex.EncodeToString(hashed[:])
	return fmt.Sprintf("%s%s", signStr[:10], dataStr), nil
}
