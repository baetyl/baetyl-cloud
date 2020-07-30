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
	storage plugin.ModelStorage
}

// NewNamespaceService NewNamespaceService
func NewNamespaceService(config *config.CloudConfig) (NamespaceService, error) {
	ms, err := plugin.GetPlugin(config.Plugin.ModelStorage)
	if err != nil {
		return nil, err
	}
	return &namespaceService{
		storage: ms.(plugin.ModelStorage),
	}, nil
}

// Get get a namespace
func (s *namespaceService) Get(namespace string) (*models.Namespace, error) {
	res, err := s.storage.GetNamespace(namespace)
	if err != nil && strings.Contains(err.Error(), "not found") {
		return nil, common.Error(common.ErrResourceNotFound, common.Field("type", "namespace"),
			common.Field("namespace", namespace))
	}
	return res, err
}

// Create Create a namespace
func (s *namespaceService) Create(namespace *models.Namespace) (*models.Namespace, error) {
	return s.storage.CreateNamespace(namespace)
}

// Delete Delete the namespace
func (s *namespaceService) Delete(namespace *models.Namespace) error {
	return s.storage.DeleteNamespace(namespace)
}
