package plugin

import (
	"github.com/baetyl/baetyl-cloud/v2/models"
	v1 "github.com/baetyl/baetyl-go/v2/spec/v1"
)

//go:generate mockgen -destination=../mock/plugin/application.go -package=plugin github.com/baetyl/baetyl-cloud/v2/plugin Application

type Application interface {
	GetApplication(namespace, name, version string) (*v1.Application, error)
	CreateApplication(namespace string, application *v1.Application) (*v1.Application, error)
	UpdateApplication(namespace string, application *v1.Application) (*v1.Application, error)
	DeleteApplication(namespace, name string) error
	ListApplication(namespace string, listOptions *models.ListOptions) (*models.ApplicationList, error)
}
