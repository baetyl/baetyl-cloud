package service

import (
	"strings"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

//go:generate mockgen -destination=../mock/service/namespace.go -package=service github.com/baetyl/baetyl-cloud/v2/service NamespaceService

// NamespaceService NamespaceService
type NamespaceService interface {
	Get(namespace string) (*models.Namespace, error)
	Create(namespace *models.Namespace) (*models.Namespace, error)
	Delete(namespace *models.Namespace) error
}

type namespaceService struct {
	namespace plugin.Namespace
}

// NewNamespaceService NewNamespaceService
func NewNamespaceService(config *config.CloudConfig) (NamespaceService, error) {
	ms, err := plugin.GetPlugin(config.Plugin.Namespace)
	if err != nil {
		return nil, err
	}
	return &namespaceService{
		namespace: ms.(plugin.Namespace),
	}, nil
}

// Get get a namespace
func (s *namespaceService) Get(namespace string) (*models.Namespace, error) {
	res, err := s.namespace.GetNamespace(namespace)
	if err != nil && strings.Contains(err.Error(), "not found") {
		return nil, common.Error(common.ErrResourceNotFound, common.Field("type", "namespace"),
			common.Field("namespace", namespace))
	}
	return res, err
}

// Create Create a namespace
func (s *namespaceService) Create(namespace *models.Namespace) (*models.Namespace, error) {
	return s.namespace.CreateNamespace(namespace)
}

// Delete Delete the namespace
func (s *namespaceService) Delete(namespace *models.Namespace) error {
	return s.namespace.DeleteNamespace(namespace)
}
