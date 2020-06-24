package plugin

import (
	"github.com/baetyl/baetyl-cloud/models"
	"github.com/baetyl/baetyl-go/spec/v1"
)

type Namespace interface {
	GetNamespace(namespace string) (*models.Namespace, error)
	CreateNamespace(namespace *models.Namespace) (*models.Namespace, error)
	DeleteNamespace(namespace *models.Namespace) error
}

type Node interface {
	GetNode(namespace, name string) (*v1.Node, error)
	CreateNode(namespace string, node *v1.Node) (*v1.Node, error)
	UpdateNode(namespace string, node *v1.Node) (*v1.Node, error)
	DeleteNode(namespace, name string) error
	ListNode(namespace string, listOptions models.ListOptions) (models.NodeList, error)
}

type Application interface {
	GetApplication(namespace, name, version string) (*v1.Application, error)
	CreateApplication(namespace string, application *v1.Application) (*v1.Application, error)
	UpdateApplication(namespace string, application *v1.Application) (*v1.Application, error)
	DeleteApplication(namespace, name string) error
	ListApplication(namespace string, listOptions *models.ListOptions) (*models.ApplicationList, error)
}

type Configuration interface {
	GetConfig(namespace, name, version string) (*v1.Configuration, error)
	CreateConfig(namespace string, configModel *v1.Configuration) (*v1.Configuration, error)
	UpdateConfig(namespace string, configurationModel *v1.Configuration) (*v1.Configuration, error)
	DeleteConfig(namespace, name string) error
	ListConfig(namespace string, listOptions *models.ListOptions) (*models.ConfigurationList, error)
}

type Secret interface {
	GetSecret(namespace, name, version string) (*v1.Secret, error)
	CreateSecret(namespace string, SecretModel *v1.Secret) (*v1.Secret, error)
	UpdateSecret(namespace string, SecretMapModel *v1.Secret) (*v1.Secret, error)
	DeleteSecret(namespace, name string) error
	ListSecret(namespace string, listOptions *models.ListOptions) (*models.SecretList, error)
}
