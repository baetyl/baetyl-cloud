package kube

import (
	"fmt"
	"strings"
	"time"

	"github.com/baetyl/baetyl-go/v2/log"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/baetyl/baetyl-go/v2/utils"
	"github.com/jinzhu/copier"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/common/util"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin/kube/apis/cloud/v1alpha1"
)

func (c *client) toSecretModel(secret *v1alpha1.Secret) *specV1.Secret {
	res := &specV1.Secret{Version: secret.ObjectMeta.ResourceVersion}
	err := copier.Copy(res, secret)
	if err != nil {
		panic(fmt.Sprintf("copier exception: %s", err.Error()))
	}
	res.Data, err = util.DecryptMap(res.Data, c.aesKey)
	if err != nil {
		log.L().Error("decrypt exception", log.Error(err))
	}
	if desc, ok := secret.Annotations[common.AnnotationDescription]; ok {
		res.Description = desc
	}
	if us, ok := secret.Annotations[common.AnnotationUpdateTimestamp]; ok {
		res.UpdateTimestamp, _ = time.Parse(common.TimeFormat, us)
	}
	res.CreationTimestamp = secret.CreationTimestamp.Time.UTC()
	res.Annotations = secret.Annotations

	return res
}

func (c *client) toSecretListModel(secretList *v1alpha1.SecretList) *models.SecretList {
	res := &models.SecretList{
		Items: make([]specV1.Secret, 0),
	}
	for _, item := range secretList.Items {
		ptr := c.toSecretModel(&item)
		res.Items = append(res.Items, *ptr)
	}
	res.Total = len(secretList.Items)
	return res
}

func (c *client) fromSecretModel(secret *specV1.Secret) (*v1alpha1.Secret, error) {
	res := &v1alpha1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "secret",
			APIVersion: v1alpha1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			ResourceVersion: secret.Version,
		},
	}
	err := copier.Copy(res, secret)
	if err != nil {
		panic(fmt.Sprintf("copier exception: %s", err.Error()))
	}
	res.Data, err = util.EncryptMap(res.Data, c.aesKey)
	if err != nil {
		log.L().Error("encrypt exception", log.Error(err))
	}
	res.Annotations = map[string]string{}
	if secret.Annotations != nil {
		res.Annotations = secret.Annotations
	}
	res.Annotations[common.AnnotationDescription] = secret.Description
	res.Annotations[common.AnnotationUpdateTimestamp] = secret.UpdateTimestamp.UTC().Format(common.TimeFormat)

	return res, nil
}

func (c *client) GetSecret(tx interface{}, namespace, name, version string) (*specV1.Secret, error) {
	options := metav1.GetOptions{ResourceVersion: version}
	defer utils.Trace(c.log.Debug, "GetSecret")()
	Secret, err := c.customClient.CloudV1alpha1().Secrets(namespace).Get(c.ctx, name, options)
	if err != nil {
		return nil, err
	}
	return c.toSecretModel(Secret), nil
}

func (c *client) CreateSecret(tx interface{}, namespace string, secretModel *specV1.Secret) (*specV1.Secret, error) {
	secretModel.UpdateTimestamp = time.Now()

	model, err := c.fromSecretModel(secretModel)
	if err != nil {
		return nil, err
	}

	defer utils.Trace(c.log.Debug, "CreateSecret")()
	Secret, err := c.customClient.CloudV1alpha1().
		Secrets(namespace).
		Create(c.ctx, model, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return c.toSecretModel(Secret), err
}

func (c *client) UpdateSecret(namespace string, secretMapModel *specV1.Secret) (*specV1.Secret, error) {
	model, err := c.fromSecretModel(secretMapModel)
	if err != nil {
		return nil, err
	}
	defer utils.Trace(c.log.Debug, "UpdateSecret")()
	SecretMap, err := c.customClient.CloudV1alpha1().
		Secrets(namespace).
		Update(c.ctx, model, metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}
	return c.toSecretModel(SecretMap), err
}

func (c *client) DeleteSecret(_ any, namespace, name string) error {
	defer utils.Trace(c.log.Debug, "DeleteSecret")()
	err := c.customClient.CloudV1alpha1().Secrets(namespace).Delete(c.ctx, name, metav1.DeleteOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil
		}
		log.L().Error("delete secret failed", log.Error(err))
	}
	return err
}

func (c *client) ListSecret(namespace string, listOptions *models.ListOptions) (*models.SecretList, error) {
	defer utils.Trace(c.log.Debug, "ListSecret")()
	list, err := c.customClient.CloudV1alpha1().Secrets(namespace).List(c.ctx, *fromListOptionsModel(listOptions))
	if err != nil {
		return nil, err
	}
	res := c.toSecretListModel(list)
	res.ListOptions = listOptions
	return res, err
}
