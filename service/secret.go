package service

import (
	"strings"

	"github.com/baetyl/baetyl-cloud/common"
	"github.com/baetyl/baetyl-cloud/config"
	"github.com/baetyl/baetyl-cloud/models"
	"github.com/baetyl/baetyl-cloud/plugin"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
)

//go:generate mockgen -destination=../mock/service/secret.go -package=plugin github.com/baetyl/baetyl-cloud/service SecretService

// SecretService SecretService
type SecretService interface {
	Get(namespace, name, version string) (*specV1.Secret, error)
	List(namespace string, listOptions *models.ListOptions) (*models.SecretList, error)
	Create(namespace string, secret *specV1.Secret) (*specV1.Secret, error)
	Update(namespace string, secret *specV1.Secret) (*specV1.Secret, error)
	Delete(namespace, name string) error
}

type secretService struct {
	storage plugin.ModelStorage
}

// NewSecretService NewSecretService
func NewSecretService(config *config.CloudConfig) (SecretService, error) {
	ms, err := plugin.GetPlugin(config.Plugin.ModelStorage)
	if err != nil {
		return nil, err
	}
	return &secretService{
		storage: ms.(plugin.ModelStorage),
	}, nil
}

// Get get a Secret
func (s *secretService) Get(namespace, name, version string) (*specV1.Secret, error) {
	res, err := s.storage.GetSecret(namespace, name, version)
	if err != nil && strings.Contains(err.Error(), "not found") {
		return nil, common.Error(common.ErrResourceNotFound, common.Field("type", "secret"),
			common.Field("name", name))
	}
	return res, err
}

// List get list Secret
func (s *secretService) List(namespace string, listOptions *models.ListOptions) (*models.SecretList, error) {
	return s.storage.ListSecret(namespace, listOptions)
}

// Create Create a Secret
func (s *secretService) Create(namespace string, secret *specV1.Secret) (*specV1.Secret, error) {
	return s.storage.CreateSecret(namespace, secret)

}

// Update update a Secret
func (s *secretService) Update(namespace string, secret *specV1.Secret) (*specV1.Secret, error) {
	return s.storage.UpdateSecret(namespace, secret)
}

// Delete Delete a Secret
func (s *secretService) Delete(namespace, name string) error {
	return s.storage.DeleteSecret(namespace, name)
}
