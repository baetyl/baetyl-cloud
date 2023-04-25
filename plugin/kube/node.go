package kube

import (
	"fmt"

	"github.com/baetyl/baetyl-go/v2/json"
	"github.com/baetyl/baetyl-go/v2/log"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/baetyl/baetyl-go/v2/utils"
	"github.com/jinzhu/copier"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin/kube/apis/cloud/v1alpha1"
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

	if val, ok := n.Attributes[specV1.KeyOptionalSysApps]; ok && val != nil {
		ss, ok := val.([]interface{})
		if !ok {
			log.L().Error(common.ErrConvertConflict, log.Any("name", specV1.KeyOptionalSysApps), log.Any("error", "failed to interface{} to []interface{}`"))
		}

		for _, d := range ss {
			s, ok := d.(string)
			if !ok {
				log.L().Error(common.ErrConvertConflict, log.Any("name", specV1.KeyOptionalSysApps), log.Any("error", "failed to interface{} to string`"))
			}
			n.SysApps = append(n.SysApps, s)
		}
	}
	if val, ok := n.Attributes[specV1.KeyAccelerator]; ok {
		if accelerator, ok := val.(string); ok {
			n.Accelerator = accelerator
		}
	}
	if val, ok := n.Attributes[specV1.KeySyncMode]; ok {
		if mode, ok := val.(string); ok {
			n.Mode = specV1.SyncMode(mode)
		}
	}
	if val, ok := n.Attributes[specV1.KeyNodeMode]; ok {
		if nodeMode, ok := val.(string); ok {
			n.NodeMode = nodeMode
		}
	}
	if val, ok := n.Attributes[specV1.KeyLink].(string); ok {
		n.Link = val
	}
	if val, ok := n.Attributes[specV1.KeyCoreId].(string); ok {
		n.CoreId = val
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

	if node.Attributes == nil {
		node.Attributes = make(map[string]interface{})
	}
	node.Attributes[specV1.KeyOptionalSysApps] = node.SysApps
	node.Attributes[specV1.KeyAccelerator] = node.Accelerator
	node.Attributes[specV1.KeyNodeMode] = node.NodeMode
	if node.Link != "" {
		node.Attributes[specV1.KeyLink] = node.Link
	}
	if node.CoreId != "" {
		node.Attributes[specV1.KeyCoreId] = node.CoreId
	}
	// set default value
	if _, ok := node.Attributes[specV1.BaetylCoreFrequency]; !ok {
		node.Attributes[specV1.BaetylCoreFrequency] = common.DefaultCoreFrequency
	}
	if _, ok := node.Attributes[specV1.BaetylCoreAPIPort]; !ok {
		node.Attributes[specV1.BaetylCoreAPIPort] = common.DefaultCoreAPIPort
	}
	if _, ok := node.Attributes[specV1.BaetylAgentPort]; !ok {
		node.Attributes[specV1.BaetylAgentPort] = common.DefaultAgentPort
	}
	if _, ok := node.Attributes[specV1.KeySyncMode]; !ok {
		node.Attributes[specV1.KeySyncMode] = specV1.CloudMode
	}

	err := copier.Copy(&n.Spec, node)
	if err != nil {
		return nil, err
	}

	n.Spec.DesireRef = &v1.LocalObjectReference{Name: node.Name}
	n.Spec.ReportRef = &v1.LocalObjectReference{Name: node.Name}
	return n, nil
}

func (c *client) GetNode(tx interface{}, namespace, name string) (*specV1.Node, error) {
	defer utils.Trace(c.log.Debug, "GetNode")()
	node, err := c.customClient.CloudV1alpha1().Nodes(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return toNodeModel(node), nil
}

func (c *client) CreateNode(tx interface{}, namespace string, node *specV1.Node) (*specV1.Node, error) {
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

func (c *client) UpdateNode(tx interface{}, namespace string, nodes []*specV1.Node) ([]*specV1.Node, error) {
	defer utils.Trace(c.log.Debug, "UpdateNode")()

	res := []*specV1.Node{}
	if len(nodes) < 1 {
		return res, nil
	}

	for _, node := range nodes {
		n, err := fromNodeModel(node)
		if err != nil {
			return nil, err
		}

		n, err = c.customClient.CloudV1alpha1().Nodes(namespace).Update(n)
		if err != nil {
			log.L().Error("update node error", log.Error(err))
			return nil, err
		}
		res = append(res, toNodeModel(n))
	}
	return res, nil
}

func (c *client) DeleteNode(tx interface{}, namespace, name string) error {
	defer utils.Trace(c.log.Debug, "DeleteNode")()
	return c.customClient.CloudV1alpha1().Nodes(namespace).Delete(name, &metav1.DeleteOptions{})
}

func (c *client) ListNode(tx interface{}, namespace string, listOptions *models.ListOptions) (*models.NodeList, error) {
	defer utils.Trace(c.log.Debug, "ListNode")()
	list, err := c.customClient.CloudV1alpha1().Nodes(namespace).List(*fromListOptionsModel(listOptions))
	if err != nil {
		return nil, err
	}
	listOptions.Continue = list.Continue
	res := toNodeListModel(list)
	res.ListOptions = listOptions
	return res, nil
}

func (c *client) CountAllNode(tx interface{}) (int, error) {
	nsList, err := c.coreV1.Namespaces().List(metav1.ListOptions{})
	if err != nil {
		return 0, err
	}
	total := 0
	for _, ns := range nsList.Items {
		list, err := c.customClient.CloudV1alpha1().Nodes(ns.Name).List(metav1.ListOptions{})
		if err != nil {
			return 0, err
		}
		total += len(list.Items)
	}
	return total, nil
}
