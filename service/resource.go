package service

import (
	"github.com/baetyl/baetyl-cloud/v2/config"
)

// AppCombinedService is a combined service contains application, configuration and secret services.
type AppCombinedService struct {
	App    ApplicationService
	Config ConfigService
	Secret SecretService
}

func NewAppCombinedService(cfg *config.CloudConfig) (*AppCombinedService, error) {
	configService, err := NewConfigService(cfg)
	if err != nil {
		return nil, err
	}
	secretService, err := NewSecretService(cfg)
	if err != nil {
		return nil, err
	}
	appService, err := NewApplicationService(cfg)
	if err != nil {
		return nil, err
	}
	return &AppCombinedService{
		App:    appService,
		Config: configService,
		Secret: secretService,
	}, nil
}
