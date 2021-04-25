package plugin

import (
	v1 "github.com/baetyl/baetyl-go/v2/spec/v1"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

//go:generate mockgen -destination=../mock/plugin/secret.go -package=plugin github.com/baetyl/baetyl-cloud/v2/plugin Secret

type Secret interface {
	GetSecret(tx interface{}, namespace, name, version string) (*v1.Secret, error)
	CreateSecret(tx interface{}, namespace string, secretModel *v1.Secret) (*v1.Secret, error)
	UpdateSecret(namespace string, secretMapModel *v1.Secret) (*v1.Secret, error)
	DeleteSecret(namespace, name string) error
	ListSecret(namespace string, listOptions *models.ListOptions) (*models.SecretList, error)
}
