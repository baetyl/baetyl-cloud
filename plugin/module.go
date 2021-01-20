package plugin

import (
	"io"

	"github.com/jmoiron/sqlx"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
)

//go:generate mockgen -destination=../mock/plugin/module.go -package=plugin github.com/baetyl/baetyl-cloud/v2/plugin Module

type Module interface {
	GetModules(name string) ([]models.Module, error)
	GetModuleByVersion(name, version string) (*models.Module, error)
	GetModuleByImage(name, image string) (*models.Module, error)
	GetLatestModule(name string) (*models.Module, error)
	CreateModule(module *models.Module) (*models.Module, error)
	UpdateModuleByVersion(module *models.Module) (*models.Module, error)
	DeleteModules(name string) error
	DeleteModuleByVersion(name, version string) error
	ListModules(filter *models.Filter) ([]models.Module, error)
	ListModulesByType(tp common.ModuleType, filter *models.Filter) ([]models.Module, error)

	GetModuleTx(tx *sqlx.Tx, name string) ([]models.Module, error)
	GetModuleByVersionTx(tx *sqlx.Tx, name, version string) (*models.Module, error)
	GetModuleByImageTx(tx *sqlx.Tx, name, image string) (*models.Module, error)
	GetLatestModuleTx(tx *sqlx.Tx, name string) (*models.Module, error)
	CreateModuleTx(tx *sqlx.Tx, module *models.Module) error
	UpdateModuleByVersionTx(tx *sqlx.Tx, module *models.Module) error
	DeleteModulesTx(tx *sqlx.Tx, name string) error
	DeleteModuleByVersionTx(tx *sqlx.Tx, name, version string) error
	ListModulesTx(tx *sqlx.Tx, filter *models.Filter) ([]models.Module, error)
	ListModulesByTypeTx(tx *sqlx.Tx, tp common.ModuleType, filter *models.Filter) ([]models.Module, error)

	// close
	io.Closer
}
