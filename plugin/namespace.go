package plugin

import (
	"io"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

//go:generate mockgen -destination=../mock/plugin/namespace.go -package=plugin github.com/baetyl/baetyl-cloud/v2/plugin Namespace

type Namespace interface {
	GetNamespace(namespace string) (*models.Namespace, error)
	CreateNamespace(namespace *models.Namespace) (*models.Namespace, error)
	DeleteNamespace(namespace *models.Namespace) error
	io.Closer
}
