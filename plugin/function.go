package plugin

import (
	"io"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

//go:generate mockgen -destination=../mock/plugin/function.go -package=plugin github.com/baetyl/baetyl-cloud/v2/plugin Function

// Function interface of Function
type Function interface {
	List(userID string) ([]models.Function, error)
	ListFunctionVersions(userID, name string) ([]models.Function, error)
	Get(userID, name, version string) (*models.Function, error)
	io.Closer
}
