package service

import (
	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

//go:generate mockgen -destination=../mock/service/module.go -package=service github.com/baetyl/baetyl-cloud/v2/service ModuleService

type ModuleService interface {
	GetModules(name string) ([]models.Module, error)
	GetModuleByVersion(name, version string) (*models.Module, error)
	GetModuleByImage(name, image string) (*models.Module, error)
	GetLatestModule(name string) (*models.Module, error)
	CreateModule(module *models.Module) (*models.Module, error)
	UpdateModuleByVersion(module *models.Module) (*models.Module, error)
	DeleteModules(name string) error
	DeleteModuleByVersion(name, version string) error
	ListModules(filter *models.Filter, tp common.ModuleType) ([]models.Module, error)

	GetLatestModuleImage(name string) (string, error)
	GetLatestModuleProgram(name, platform string) (string, error)
}

// NewModuleService
func NewModuleService(config *config.CloudConfig) (ModuleService, error) {
	ds, err := plugin.GetPlugin(config.Plugin.Module)
	if err != nil {
		return nil, err
	}
	return ds.(plugin.Module), nil
}
