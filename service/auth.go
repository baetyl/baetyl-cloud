package service

import (
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

//go:generate mockgen -destination=../mock/service/auth.go -package=service github.com/baetyl/baetyl-cloud/v2/service AuthService

type AuthService interface {
	plugin.Auth
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
