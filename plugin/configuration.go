package plugin

import (
	v1 "github.com/baetyl/baetyl-go/v2/spec/v1"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

//go:generate mockgen -destination=../mock/plugin/configuration.go -package=plugin github.com/baetyl/baetyl-cloud/v2/plugin Configuration

type Configuration interface {
	GetConfig(tx interface{}, namespace, name, version string) (*v1.Configuration, error)
	CreateConfig(tx interface{}, namespace string, configModel *v1.Configuration) (*v1.Configuration, error)
	UpdateConfig(namespace string, configurationModel *v1.Configuration) (*v1.Configuration, error)
	DeleteConfig(namespace, name string) error
	ListConfig(namespace string, listOptions *models.ListOptions) (*models.ConfigurationList, error)
}
