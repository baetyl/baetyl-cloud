package plugin

import (
	"github.com/jmoiron/sqlx"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

//go:generate mockgen -destination=../mock/plugin/node.go -package=plugin github.com/baetyl/baetyl-cloud/v2/plugin Node

type Module interface {
	GetModule(name string) (*models.Module, error)
	GetModuleByVersion(name, version string) (*models.Module, error)
	CreateModule(module *models.Module) error
	UpdateModule(module *models.Module) error
	DeleteModule(name string) error
	ListModule(filter *models.Filter) ([]models.Module, error)
	ListModuleWithOptions(tp string, hidden bool, filter *models.Filter) ([]models.Module, error)

	GetModuleTx(tx *sqlx.Tx, name string) (*models.Module, error)
	GetModuleByVersionTx(tx *sqlx.Tx, name, version string) (*models.Module, error)
	CreateModuleTx(tx *sqlx.Tx, module *models.Module) error
	UpdateModuleTx(tx *sqlx.Tx, module *models.Module) error
	DeleteModuleTx(tx *sqlx.Tx, name string) error
	ListModuleTx(tx *sqlx.Tx, filter *models.Filter) ([]models.Module, error)
	ListModuleWithOptionsTx(tx *sqlx.Tx, tp string, hidden bool, filter *models.Filter) ([]models.Module, error)
}
