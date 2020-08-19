package kube

import (
	"encoding/json"
	"fmt"
	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin/kube/apis/cloud/v1alpha1"
	"github.com/baetyl/baetyl-go/v2/log"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/baetyl/baetyl-go/v2/utils"
	"github.com/jinzhu/copier"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func toNodeModel(node *v1alpha1.Node) *specV1.Node {
	n := &specV1.Node{
		Name:      node.ObjectMeta.Name,
		Namespace: node.ObjectMeta.Namespace,
		Version:   node.ObjectMeta.ResourceVersion,
		Labels:    node.ObjectMeta.Labels,
	}
	err := copier.Copy(n, &node.Spec)
	if err != nil {
		panic(fmt.Sprintf("copier exception: %s", err.Error()))
	}
	n.CreationTimestamp = node.CreationTimestamp.Time.UTC()

	if desc, ok := node.Annotations[common.AnnotationDescription]; ok {
		n.Description = desc
	}

	if metadata, ok := node.Annotations[common.AnnotationMetadata]; ok {
		n.Annotations = map[string]string{}
		err = json.Unmarshal([]byte(metadata), &n.Annotations)
		if err != nil {
			log.L().Error("report unmarshal exception", log.Error(err))
		}
	}
	return n
}

func toNodeListModel(list *v1alpha1.NodeList) *models.NodeList {
	res := &models.NodeList{
		Items: make([]specV1.Node, 0),
	}
	for _, node := range list.Items {
		n := toNodeModel(&node)
		res.Items = append(res.Items, *n)
	}
	res.Total = len(list.Items)
	return res
}

func fromNodeModel(node *specV1.Node) (*v1alpha1.Node, error) {
	n := &v1alpha1.Node{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Node",
			APIVersion: v1alpha1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:            node.Name,
			Namespace:       node.Namespace,
			Labels:          node.Labels,
			ResourceVersion: node.Version,
			Annotations:     map[string]string{},
		},
	}
	n.Annotations[common.AnnotationDescription] = node.Description

	if node.Annotations != nil {
		metadata, err := json.Marshal(node.Annotations)
		if err != nil {
			log.L().Error("node desire marshal exception", log.Error(err))
			return nil, err
		}
		n.Annotations[common.AnnotationMetadata] = string(metadata)
	}

	err := copier.Copy(&n.Spec, node)
	if err != nil {
		return nil, err
	}

	n.Spec.DesireRef = &v1.LocalObjectReference{Name: node.Name}
	n.Spec.ReportRef = &v1.LocalObjectReference{Name: node.Name}
	return n, nil
}

func (c *client) GetNode(namespace, name string) (*specV1.Node, error) {
	defer utils.Trace(c.log.Debug, "GetNode")()
	node, err := c.customClient.CloudV1alpha1().Nodes(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return toNodeModel(node), nil
}

func (c *client) CreateNode(namespace string, node *specV1.Node) (*specV1.Node, error) {
	n, err := fromNodeModel(node)
	if err != nil {
		return nil, err
	}

	defer utils.Trace(c.log.Debug, "CreateNode")()
	n, err = c.customClient.CloudV1alpha1().Nodes(namespace).Create(n)
	if err != nil {
		return nil, err
	}

	return toNodeModel(n), nil
}

func (c *client) UpdateNode(namespace string, node *specV1.Node) (*specV1.Node, error) {
	n, err := fromNodeModel(node)
	if err != nil {
		return nil, err
	}
	defer utils.Trace(c.log.Debug, "UpdateNode")()
	n, err = c.customClient.CloudV1alpha1().Nodes(namespace).Update(n)
	if err != nil {
		log.L().Error("update node error", log.Error(err))
		return nil, err
	}
	return toNodeModel(n), nil
}

func (c *client) DeleteNode(namespace, name string) error {
	defer utils.Trace(c.log.Debug, "DeleteNode")()
	return c.customClient.CloudV1alpha1().Nodes(namespace).Delete(name, &metav1.DeleteOptions{})
}

func (c *client) ListNode(namespace string, listOptions *models.ListOptions) (*models.NodeList, error) {
	defer utils.Trace(c.log.Debug, "ListNode")()
	list, err := c.customClient.CloudV1alpha1().Nodes(namespace).List(*fromListOptionsModel(listOptions))
	if err != nil {
		return nil, err
	}
	listOptions.Continue = list.Continue
	res, err := toNodeListModel(list), nil
	if err != nil {
		return nil, err
	}
	res.ListOptions = listOptions
	return res, nil
}
