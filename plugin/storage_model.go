package plugin

import (
	"github.com/baetyl/baetyl-cloud/models"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
)

//go:generate mockgen -destination=../mock/plugin/storage_model.go -package=plugin github.com/baetyl/baetyl-cloud/plugin ModelStorage

// ModelStorage ModelStorage
type ModelStorage interface {
	GetNamespace(namespace string) (*models.Namespace, error)
	CreateNamespace(namespace *models.Namespace) (*models.Namespace, error)
	DeleteNamespace(namespace *models.Namespace) error

	GetNode(namespace, name string) (*specV1.Node, error)
	CreateNode(namespace string, node *specV1.Node) (*specV1.Node, error)
	UpdateNode(namespace string, node *specV1.Node) (*specV1.Node, error)
	DeleteNode(namespace, name string) error
	ListNode(namespace string, listOptions *models.ListOptions) (*models.NodeList, error)

	GetConfig(namespace, name, version string) (*specV1.Configuration, error)
	CreateConfig(namespace string, config *specV1.Configuration) (*specV1.Configuration, error)
	UpdateConfig(namespace string, config *specV1.Configuration) (*specV1.Configuration, error)
	DeleteConfig(namespace, name string) error
	ListConfig(namespace string, listOptions *models.ListOptions) (*models.ConfigurationList, error)

	GetApplication(namespace, name, version string) (*specV1.Application, error)
	CreateApplication(namespace string, application *specV1.Application) (*specV1.Application, error)
	UpdateApplication(namespace string, application *specV1.Application) (*specV1.Application, error)
	DeleteApplication(namespace, name string) error
	ListApplication(namespace string, listOptions *models.ListOptions) (*models.ApplicationList, error)

	GetSecret(namespace, name, version string) (*specV1.Secret, error)
	CreateSecret(namespace string, config *specV1.Secret) (*specV1.Secret, error)
	UpdateSecret(namespace string, config *specV1.Secret) (*specV1.Secret, error)
	DeleteSecret(namespace, name string) error
	ListSecret(namespace string, listOptions *models.ListOptions) (*models.SecretList, error)

	IsLabelMatch(labelSelector string, labels map[string]string) (bool, error)

	Shadow
}
