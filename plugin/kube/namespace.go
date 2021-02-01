package kube

import (
	"fmt"
	"strings"

	"github.com/baetyl/baetyl-go/v2/utils"
	"github.com/jinzhu/copier"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/baetyl/baetyl-cloud/v2/models"
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
	defer utils.Trace(c.log.Debug, "GetNamespace")()
	n, err := c.coreV1.Namespaces().Get(namespace, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return toNamespaceModel(n), nil
}

func (c *client) CreateNamespace(namespace *models.Namespace) (*models.Namespace, error) {
	defer utils.Trace(c.log.Debug, "CreateNamespace")()
	n, err := c.coreV1.Namespaces().Create(fromNamespaceModel(namespace))
	if err != nil {
		return nil, err
	}
	return toNamespaceModel(n), nil
}

func (c *client) ListNamespace(listOptions *models.ListOptions) (*models.NamespaceList, error) {
	defer utils.Trace(c.log.Debug, "ListNamespace")()
	list, err := c.coreV1.Namespaces().List(*fromListOptionsModel(listOptions))
	if err != nil {
		return nil, err
	}
	listOptions.Continue = list.Continue
	res := toNamespaceListModel(list)
	res.ListOptions = listOptions
	return res, nil
}

func (c *client) DeleteNamespace(namespace *models.Namespace) error {
	defer utils.Trace(c.log.Debug, "DeleteNamespace")()
	err := c.coreV1.Namespaces().Delete(namespace.Name, &metav1.DeleteOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil
		}
		return err
	}
	if n, _ := c.coreV1.Namespaces().Get(namespace.Name, metav1.GetOptions{}); n != nil {
		_, err = c.coreV1.Namespaces().Finalize(fromNamespaceModel(namespace))
	}
	return err
}

func toNamespaceListModel(list *v1.NamespaceList) *models.NamespaceList {
	res := &models.NamespaceList{
		Items: make([]models.Namespace, 0),
	}
	for _, ns := range list.Items {
		n := toNamespaceModel(&ns)
		res.Items = append(res.Items, *n)
	}
	res.Total = len(list.Items)
	return res
}
