package plugin

import (
	"github.com/baetyl/baetyl-cloud/v2/models"
	v1 "github.com/baetyl/baetyl-go/v2/spec/v1"
)

//go:generate mockgen -destination=../mock/plugin/secret.go -package=plugin github.com/baetyl/baetyl-cloud/v2/plugin Secret

type Secret interface {
	GetSecret(namespace, name, version string) (*v1.Secret, error)
	CreateSecret(namespace string, SecretModel *v1.Secret) (*v1.Secret, error)
	UpdateSecret(namespace string, SecretMapModel *v1.Secret) (*v1.Secret, error)
	DeleteSecret(namespace, name string) error
	ListSecret(namespace string, listOptions *models.ListOptions) (*models.SecretList, error)
}
