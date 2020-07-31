package plugin

import (
	"github.com/baetyl/baetyl-cloud/v2/models"
	"io"
)

//go:generate mockgen -destination=../mock/plugin/shadow.go -package=plugin github.com/baetyl/baetyl-cloud/v2/plugin Shadow

// Shadow
type Shadow interface {
	Get(namespace, name string) (*models.Shadow, error)
	Create(shadow *models.Shadow) (*models.Shadow, error)
	Delete(namespace, name string) error
	UpdateDesire(shadow *models.Shadow) (*models.Shadow, error)
	UpdateReport(shadow *models.Shadow) (*models.Shadow, error)
	List(namespace string, nodeList *models.NodeList) (*models.ShadowList, error)
	io.Closer
}
