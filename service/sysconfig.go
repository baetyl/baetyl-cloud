package service

import (
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

//go:generate mockgen -destination=../mock/service/sysconfig.go -package=plugin github.com/baetyl/baetyl-cloud/service SysConfigService

type SysConfigService interface {
	GetSysConfig(tp, key string) (*models.SysConfig, error)
	ListSysConfigAll(tp string) ([]models.SysConfig, error)
}

type sysConfigService struct {
	cfg       *config.CloudConfig
	dbStorage plugin.DBStorage
}

// NewSysConfigService
func NewSysConfigService(config *config.CloudConfig) (SysConfigService, error) {
	ds, err := plugin.GetPlugin(config.Plugin.DatabaseStorage)
	if err != nil {
		return nil, err
	}
	return &sysConfigService{
		cfg:       config,
		dbStorage: ds.(plugin.DBStorage),
	}, nil
}

func (s *sysConfigService) GetSysConfig(tp, key string) (*models.SysConfig, error) {
	return s.dbStorage.GetSysConfig(tp, key)
}

func (s *sysConfigService) ListSysConfigAll(tp string) ([]models.SysConfig, error) {
	return s.dbStorage.ListSysConfigAll(tp)
}
