package kube

import (
	"fmt"
	"time"

	"github.com/baetyl/baetyl-cloud/common"
	"github.com/baetyl/baetyl-cloud/models"
	"github.com/baetyl/baetyl-cloud/plugin/kube/apis/cloud/v1alpha1"
	"github.com/baetyl/baetyl-go/log"
	specV1 "github.com/baetyl/baetyl-go/spec/v1"
	"github.com/jinzhu/copier"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *client) toSecretModel(Secret *v1alpha1.Secret) *specV1.Secret {
	res := &specV1.Secret{Version: Secret.ObjectMeta.ResourceVersion}
	err := copier.Copy(res, Secret)
	if err != nil {
		panic(fmt.Sprintf("copier exception: %s", err.Error()))
	}
	res.Data, err = DecryptMap(res.Data, c.aesKey)
	if err != nil {
		log.L().Error("decrypt exception", log.Error(err))
	}
	if desc, ok := Secret.Annotations[common.AnnotationDescription]; ok {
		res.Description = desc
	}
	if us, ok := Secret.Annotations[common.AnnotationUpdateTimestamp]; ok {
		res.UpdateTimestamp, _ = time.Parse(common.TimeFormat, us)
	}
	res.CreationTimestamp = Secret.CreationTimestamp.Time.UTC()
	res.Annotations = Secret.Annotations

	return res
}

func (c *client) toSecretListModel(SecretList *v1alpha1.SecretList) *models.SecretList {
	res := &models.SecretList{
		Items: make([]specV1.Secret, 0),
	}
	for _, item := range SecretList.Items {
		ptr := c.toSecretModel(&item)
		res.Items = append(res.Items, *ptr)
	}
	res.Total = len(SecretList.Items)
	return res
}

func (c *client) fromSecretModel(Secret *specV1.Secret) (*v1alpha1.Secret, error) {
	res := &v1alpha1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: v1alpha1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			ResourceVersion: Secret.Version,
		},
	}
	err := copier.Copy(res, Secret)
	if err != nil {
		panic(fmt.Sprintf("copier exception: %s", err.Error()))
	}
	res.Data, err = EncryptMap(res.Data, c.aesKey)
	if err != nil {
		log.L().Error("encrypt exception", log.Error(err))
	}
	res.Annotations = map[string]string{}
	if Secret.Annotations != nil {
		res.Annotations = Secret.Annotations
	}
	res.Annotations[common.AnnotationDescription] = Secret.Description
	res.Annotations[common.AnnotationUpdateTimestamp] = Secret.UpdateTimestamp.UTC().Format(common.TimeFormat)

	return res, nil
}

func (c *client) GetSecret(namespace, name, version string) (*specV1.Secret, error) {
	options := metav1.GetOptions{ResourceVersion: version}
	Secret, err := c.customClient.CloudV1alpha1().Secrets(namespace).Get(name, options)
	if err != nil {
		return nil, err
	}
	return c.toSecretModel(Secret), nil
}

func (c *client) CreateSecret(namespace string, SecretModel *specV1.Secret) (*specV1.Secret, error) {
	model, err := c.fromSecretModel(SecretModel)
	if err != nil {
		return nil, err
	}
	SecretModel.UpdateTimestamp = time.Now()
	Secret, err := c.customClient.CloudV1alpha1().
		Secrets(namespace).
		Create(model)
	if err != nil {
		return nil, err
	}
	return c.toSecretModel(Secret), err
}

func (c *client) UpdateSecret(namespace string, SecretMapModel *specV1.Secret) (*specV1.Secret, error) {
	model, err := c.fromSecretModel(SecretMapModel)
	if err != nil {
		return nil, err
	}
	SecretMap, err := c.customClient.CloudV1alpha1().
		Secrets(namespace).
		Update(model)
	if err != nil {
		return nil, err
	}
	return c.toSecretModel(SecretMap), err
}

func (c *client) DeleteSecret(namespace, name string) error {
	return c.customClient.CloudV1alpha1().Secrets(namespace).Delete(name, &metav1.DeleteOptions{})
}

func (c *client) ListSecret(namespace string, listOptions *models.ListOptions) (*models.SecretList, error) {
	list, err := c.customClient.CloudV1alpha1().Secrets(namespace).List(*fromListOptionsModel(listOptions))
	if err != nil {
		return nil, err
	}
	res := c.toSecretListModel(list)
	res.ListOptions = listOptions
	return res, err
}
