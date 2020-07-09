package kube

import (
	"fmt"
	"github.com/baetyl/baetyl-go/log"
	"time"

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
	beforeRequest := time.Now().UnixNano()
	n, err := c.coreV1.Namespaces().Get(namespace, metav1.GetOptions{})
	afterRequest := time.Now().UnixNano()
	log.L().Debug("kube GetNamespace", log.Any("cost time (ns)", afterRequest-beforeRequest))
	if err != nil {
		return nil, err
	}
	return toNamespaceModel(n), nil
}

func (c *client) CreateNamespace(namespace *models.Namespace) (*models.Namespace, error) {
	beforeRequest := time.Now().UnixNano()
	n, err := c.coreV1.Namespaces().Create(fromNamespaceModel(namespace))
	afterRequest := time.Now().UnixNano()
	log.L().Debug("kube CreateNamespace", log.Any("cost time (ns)", afterRequest-beforeRequest))
	if err != nil {
		return nil, err
	}
	return toNamespaceModel(n), nil
}

func (c *client) DeleteNamespace(namespace *models.Namespace) error {
	beforeRequest := time.Now().UnixNano()
	err := c.coreV1.Namespaces().Delete(namespace.Name, &metav1.DeleteOptions{})
	afterRequest := time.Now().UnixNano()
	log.L().Debug("kube DeleteNamespace Delete", log.Any("cost time (ns)", afterRequest-beforeRequest))
	if err != nil {
		return err
	}
	beforeRequest = time.Now().UnixNano()
	n, _ := c.coreV1.Namespaces().Get(namespace.Name, metav1.GetOptions{})
	afterRequest = time.Now().UnixNano()
	log.L().Debug("kube DeleteNamespace Get", log.Any("cost time (ns)", afterRequest-beforeRequest))
	if n != nil {
		beforeRequest = time.Now().UnixNano()
		_, err = c.coreV1.Namespaces().Finalize(fromNamespaceModel(namespace))
		afterRequest = time.Now().UnixNano()
		log.L().Debug("kube DeleteNamespace Finalize", log.Any("cost time (ns)", afterRequest-beforeRequest))
	}
	return err
}
