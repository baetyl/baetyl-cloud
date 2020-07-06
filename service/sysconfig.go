package service

import (
	"database/sql"
	"github.com/baetyl/baetyl-cloud/config"
	"github.com/baetyl/baetyl-cloud/models"
	"github.com/baetyl/baetyl-cloud/plugin"
)

//go:generate mockgen -destination=../mock/service/sysconfig.go -package=plugin github.com/baetyl/baetyl-cloud/service SysConfigService

type SysConfigService interface {
	GetSysConfig(tp, key string) (*models.SysConfig, error)
	ListSysConfigAll(tp string) ([]models.SysConfig, error)

	ListSysConfig(tp string, page, size int) ([]models.SysConfig, error)
	CreateSysConfig(sysConfig *models.SysConfig) (sql.Result, error)
	DeleteSysConfig(tp, key string) (sql.Result, error)
	UpdateSysConfig(sysConfig *models.SysConfig) (sql.Result, error)
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

func (s *sysConfigService) ListSysConfig(tp string, page, size int) ([]models.SysConfig, error){
	return s.dbStorage.ListSysConfig(tp, page, size)
}

func (s *sysConfigService) CreateSysConfig(sysConfig *models.SysConfig) (sql.Result, error){
	return s.dbStorage.CreateSysConfig(sysConfig)
}

func (s *sysConfigService) DeleteSysConfig(tp, key string) (sql.Result, error) {
	return s.dbStorage.DeleteSysConfig(tp, key)
}

func (s *sysConfigService) UpdateSysConfig(sysConfig *models.SysConfig) (sql.Result, error) {
	return s.dbStorage.UpdateSysConfig(sysConfig)
}
