package plugin

import (
	v1 "github.com/baetyl/baetyl-go/v2/spec/v1"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

//go:generate mockgen -destination=../mock/plugin/application.go -package=plugin github.com/baetyl/baetyl-cloud/v2/plugin Application

type Application interface {
	GetApplication(namespace, name, version string) (*v1.Application, error)
	CreateApplication(tx interface{}, namespace string, application *v1.Application) (*v1.Application, error)
	UpdateApplication(tx interface{}, namespace string, application *v1.Application) (*v1.Application, error)
	DeleteApplication(tx interface{}, namespace, name string) error
	ListApplication(tx interface{}, namespace string, listOptions *models.ListOptions) (*models.ApplicationList, error)
}
