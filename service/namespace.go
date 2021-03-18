package service

import (
	"strings"

	"github.com/baetyl/baetyl-go/v2/errors"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

//go:generate mockgen -destination=../mock/service/namespace.go -package=service github.com/baetyl/baetyl-cloud/v2/service NamespaceService

const KeyCreateExtraNamespaceResources = "createExtraNamespaceResources"
const KeyDeleteExtraNamespaceResources = "deleteExtraNamespaceResources"

type CreateExtraNamespaceResourcesFunc func(namespace string) error
type DeleteExtraNamespaceResourcesFunc func(namespace string) error

// NamespaceService NamespaceService
type NamespaceService interface {
	Get(namespace string) (*models.Namespace, error)
	Create(namespace *models.Namespace) (*models.Namespace, error)
	List(listOptions *models.ListOptions) (*models.NamespaceList, error)
	Delete(namespace *models.Namespace) error
}

type NamespaceServiceImpl struct {
	namespace plugin.Namespace
	Hooks     map[string]interface{}
}

// NewNamespaceService NewNamespaceService
func NewNamespaceService(config *config.CloudConfig) (NamespaceService, error) {
	ms, err := plugin.GetPlugin(config.Plugin.Resource)
	if err != nil {
		return nil, err
	}
	return &NamespaceServiceImpl{
		namespace: ms.(plugin.Namespace),
		Hooks:     make(map[string]interface{}),
	}, nil
}

// Get get a namespace
func (s *NamespaceServiceImpl) Get(namespace string) (*models.Namespace, error) {
	res, err := s.namespace.GetNamespace(namespace)
	if err != nil && strings.Contains(err.Error(), "not found") {
		return nil, common.Error(common.ErrResourceNotFound, common.Field("type", "namespace"),
			common.Field("namespace", namespace))
	}
	return res, err
}

// Create Create a namespace
func (s *NamespaceServiceImpl) Create(namespace *models.Namespace) (*models.Namespace, error) {
	if createFunc, ok := s.Hooks[KeyCreateExtraNamespaceResources].(CreateExtraNamespaceResourcesFunc); ok {
		if err := createFunc(namespace.Name); err != nil {
			return nil, errors.Trace(err)
		}
	}
	return s.namespace.CreateNamespace(namespace)
}

// List get list namespace
func (s *NamespaceServiceImpl) List(listOptions *models.ListOptions) (*models.NamespaceList, error) {
	return s.namespace.ListNamespace(listOptions)
}

// Delete Delete the namespace
func (s *NamespaceServiceImpl) Delete(namespace *models.Namespace) error {
	if deleteFunc, ok := s.Hooks[KeyDeleteExtraNamespaceResources].(DeleteExtraNamespaceResourcesFunc); ok {
		if err := deleteFunc(namespace.Name); err != nil {
			return errors.Trace(err)
		}
	}
	return s.namespace.DeleteNamespace(namespace)
}
