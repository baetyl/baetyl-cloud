package plugin

import (
	"io"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

//go:generate mockgen -destination=../mock/plugin/shadow.go -package=plugin github.com/baetyl/baetyl-cloud/v2/plugin Shadow

// Shadow
type Shadow interface {
	Get(tx interface{}, namespace, name string) (*models.Shadow, error)
	ListShadowByNames(tx interface{}, namespace string, names []string) ([]*models.Shadow, error)
	Create(tx interface{}, shadow *models.Shadow) (*models.Shadow, error)
	Delete(namespace, name string) error
	UpdateDesire(tx interface{}, shadow *models.Shadow) error
	UpdateDesires(tx interface{}, shadows []*models.Shadow) error
	UpdateReport(shadow *models.Shadow) (*models.Shadow, error)
	List(namespace string, nodeList *models.NodeList) (*models.ShadowList, error)
	io.Closer
}
