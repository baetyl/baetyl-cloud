package service

import (
	"fmt"

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
	UpdateModule(module *models.Module) (*models.Module, error)
	DeleteModules(name string) error
	DeleteModuleByVersion(name, version string) error
	ListModules(filter *models.Filter) ([]models.Module, error)
	ListOptionalSysModules(filter *models.Filter) ([]models.Module, error)
	ListRuntimeModules(filter *models.Filter) ([]models.Module, error)

	GetLatestModuleImage(name string) (string, error)
	GetLatestModuleProgram(name, platform string) (string, error)
}

type moduleService struct {
	module plugin.Module
}

// NewModuleService
func NewModuleService(config *config.CloudConfig) (ModuleService, error) {
	ds, err := plugin.GetPlugin(config.Plugin.Module)
	if err != nil {
		return nil, err
	}

	p := &moduleService{
		module: ds.(plugin.Module),
	}

	return p, nil
}

func (m *moduleService) GetModules(name string) ([]models.Module, error) {
	return m.module.GetModules(name)
}

func (m *moduleService) GetModuleByVersion(name, version string) (*models.Module, error) {
	return m.module.GetModuleByVersion(name, version)
}

func (m *moduleService) GetModuleByImage(name, image string) (*models.Module, error) {
	return m.module.GetModuleByImage(name, image)
}

func (m *moduleService) GetLatestModule(name string) (*models.Module, error) {
	return m.module.GetLatestModule(name)
}

func (m *moduleService) CreateModule(module *models.Module) (*models.Module, error) {
	return m.module.CreateModule(module)
}

func (m *moduleService) UpdateModule(module *models.Module) (*models.Module, error) {
	return m.module.UpdateModuleByVersion(module)
}

func (m *moduleService) DeleteModules(name string) error {
	return m.module.DeleteModules(name)
}

func (m *moduleService) DeleteModuleByVersion(name, version string) error {
	return m.module.DeleteModuleByVersion(name, version)
}

func (m *moduleService) ListModules(filter *models.Filter) ([]models.Module, error) {
	return m.module.ListModules(filter)
}

func (m *moduleService) ListOptionalSysModules(filter *models.Filter) ([]models.Module, error) {
	t := common.Type_System_Optional
	return m.module.ListModulesByType(t, filter)
}

func (m *moduleService) ListRuntimeModules(filter *models.Filter) ([]models.Module, error) {
	t := common.Type_User_RUNTIME
	return m.module.ListModulesByType(t, filter)
}

func (m *moduleService) GetLatestModuleImage(name string) (string, error) {
	module, err := m.module.GetLatestModule(name)
	if err != nil {
		return "", err
	}
	return module.Image, nil
}

func (m *moduleService) GetLatestModuleProgram(name, platform string) (string, error) {
	module, err := m.module.GetLatestModule(name)
	if err != nil {
		return "", err
	}
	for k, v := range module.Programs {
		if k == platform {
			return v, nil
		}
	}
	return "", common.Error(common.ErrResourceNotFound,
		common.Field("type", "program"),
		common.Field("name", fmt.Sprintf("%s-%s", name, platform)))
}
