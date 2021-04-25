package service

import (
	"strings"

	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

//go:generate mockgen -destination=../mock/service/secret.go -package=service github.com/baetyl/baetyl-cloud/v2/service SecretService

// SecretService SecretService
type SecretService interface {
	Get(namespace, name, version string) (*specV1.Secret, error)
	List(namespace string, listOptions *models.ListOptions) (*models.SecretList, error)
	Create(tx interface{}, namespace string, secret *specV1.Secret) (*specV1.Secret, error)
	Update(namespace string, secret *specV1.Secret) (*specV1.Secret, error)
	Delete(namespace, name string) error
}

type secretService struct {
	secret plugin.Secret
}

// NewSecretService NewSecretService
func NewSecretService(config *config.CloudConfig) (SecretService, error) {
	secret, err := plugin.GetPlugin(config.Plugin.Resource)
	if err != nil {
		return nil, err
	}
	return &secretService{
		secret: secret.(plugin.Secret),
	}, nil
}

// Get get a Secret
func (s *secretService) Get(namespace, name, version string) (*specV1.Secret, error) {
	res, err := s.secret.GetSecret(nil, namespace, name, version)
	if err != nil && strings.Contains(err.Error(), "not found") {
		return nil, common.Error(common.ErrResourceNotFound, common.Field("type", "secret"),
			common.Field("name", name))
	}
	return res, err
}

// List get list Secret
func (s *secretService) List(namespace string, listOptions *models.ListOptions) (*models.SecretList, error) {
	return s.secret.ListSecret(namespace, listOptions)
}

// Create Create a Secret
func (s *secretService) Create(tx interface{}, namespace string, secret *specV1.Secret) (*specV1.Secret, error) {
	return s.secret.CreateSecret(tx, namespace, secret)

}

// Update update a Secret
func (s *secretService) Update(namespace string, secret *specV1.Secret) (*specV1.Secret, error) {
	return s.secret.UpdateSecret(namespace, secret)
}

// Delete Delete a Secret
func (s *secretService) Delete(namespace, name string) error {
	return s.secret.DeleteSecret(namespace, name)
}
