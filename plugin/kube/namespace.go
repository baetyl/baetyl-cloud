package kube

import (
	"fmt"

	"github.com/baetyl/baetyl-cloud/models"
	"github.com/jinzhu/copier"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func toNamespaceModel(namespace *v1.Namespace) *models.Namespace {
	res := &models.Namespace{}
	err := copier.Copy(res, namespace)
	if err != nil {
		panic(fmt.Sprintf("copier exception: %s", err.Error()))
	}
	return res
}

func fromNamespaceModel(namespace *models.Namespace) *v1.Namespace {
	res := &v1.Namespace{}
	err := copier.Copy(res, namespace)
	if err != nil {
		panic(fmt.Sprintf("copier exception: %s", err.Error()))
	}
	return res
}

func (c *client) GetNamespace(namespace string) (*models.Namespace, error) {
	n, err := c.coreV1.Namespaces().Get(namespace, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return toNamespaceModel(n), nil
}

func (c *client) CreateNamespace(namespace *models.Namespace) (*models.Namespace, error) {
	n, err := c.coreV1.Namespaces().Create(fromNamespaceModel(namespace))
	if err != nil {
		return nil, err
	}
	return toNamespaceModel(n), nil
}

func (c *client) DeleteNamespace(namespace *models.Namespace) error {
	err := c.coreV1.Namespaces().Delete(namespace.Name, &metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	n, _ := c.coreV1.Namespaces().Get(namespace.Name, metav1.GetOptions{})
	if n != nil {
		_, err = c.coreV1.Namespaces().Finalize(fromNamespaceModel(namespace))
	}
	return err
}
