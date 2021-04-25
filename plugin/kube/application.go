package kube

import (
	"fmt"

	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/baetyl/baetyl-go/v2/utils"
	"github.com/jinzhu/copier"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin/kube/apis/cloud/v1alpha1"
)

func toAppModel(app *v1alpha1.Application) *specV1.Application {
	description, _ := app.Annotations[common.AnnotationDescription]
	nodeSelector, _ := app.Annotations[common.AnnotationNodeSelector]
	res := &specV1.Application{
		Name:         app.ObjectMeta.Name,
		Namespace:    app.ObjectMeta.Namespace,
		Version:      app.ObjectMeta.ResourceVersion,
		Description:  description,
		NodeSelector: nodeSelector,
		Labels:       app.ObjectMeta.Labels,
	}

	err := copier.Copy(res, &app.Spec)
	if err != nil {
		panic(fmt.Sprintf("copier exception: %s", err.Error()))
	}
	res.CreationTimestamp = app.CreationTimestamp.Time.UTC()
	return res
}

func toAppListModel(list *v1alpha1.ApplicationList) *models.ApplicationList {
	res := &models.ApplicationList{
		Items: make([]models.AppItem, 0),
	}
	for _, item := range list.Items {
		description, _ := item.Annotations[common.AnnotationDescription]
		nodeSelector, _ := item.Annotations[common.AnnotationNodeSelector]
		res.Items = append(res.Items, models.AppItem{
			Name:              item.ObjectMeta.Name,
			Type:              item.Spec.Type,
			Namespace:         item.ObjectMeta.Namespace,
			Version:           item.ObjectMeta.ResourceVersion,
			Labels:            item.ObjectMeta.Labels,
			Selector:          item.Spec.Selector,
			CreationTimestamp: item.CreationTimestamp.Time.UTC(),
			Description:       description,
			NodeSelector:      nodeSelector,
			System:            item.Spec.System,
		})
	}

	res.Total = len(list.Items)
	return res
}

func fromAppModel(namespace string, app *specV1.Application) *v1alpha1.Application {
	res := &v1alpha1.Application{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Application",
			APIVersion: v1alpha1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:            app.Name,
			Namespace:       namespace,
			ResourceVersion: app.Version,
			Labels:          app.Labels,
			Annotations:     map[string]string{},
		},
	}

	if app.Description != "" {
		res.Annotations[common.AnnotationDescription] = app.Description
	}

	if app.NodeSelector != "" {
		res.Annotations[common.AnnotationNodeSelector] = app.NodeSelector
	}

	err := copier.Copy(&res.Spec, app)
	if err != nil {
		panic(fmt.Sprintf("copier exception: %s", err.Error()))
	}
	return res
}

func fromListOptionsModel(listOptions *models.ListOptions) *metav1.ListOptions {
	res := &metav1.ListOptions{}
	err := copier.Copy(res, listOptions)
	if err != nil {
		panic(fmt.Sprintf("copier exception: %s", err.Error()))
	}
	return res
}

func (c *client) GetApplication(namespace, name, version string) (*specV1.Application, error) {
	defer utils.Trace(c.log.Debug, "GetApplication")()
	options := metav1.GetOptions{ResourceVersion: version}
	app, err := c.customClient.CloudV1alpha1().Applications(namespace).Get(name, options)
	if err != nil {
		return nil, err
	}
	return toAppModel(app), nil
}

func (c *client) CreateApplication(tx interface{}, namespace string, application *specV1.Application) (*specV1.Application, error) {
	app := fromAppModel(namespace, application)
	defer utils.Trace(c.log.Debug, "CreateApplication")()
	app, err := c.customClient.CloudV1alpha1().Applications(namespace).Create(app)
	if err != nil {
		return nil, err
	}
	res := toAppModel(app)
	return res, nil
}

func (c *client) UpdateApplication(namespace string, application *specV1.Application) (*specV1.Application, error) {
	app := fromAppModel(namespace, application)
	defer utils.Trace(c.log.Debug, "UpdateApplication")()
	app, err := c.customClient.CloudV1alpha1().Applications(namespace).Update(app)
	if err != nil {
		return nil, err
	}
	return toAppModel(app), nil
}

func (c *client) DeleteApplication(namespace, name string) error {
	defer utils.Trace(c.log.Debug, "DeleteApplication")()
	err := c.customClient.CloudV1alpha1().Applications(namespace).Delete(name, &metav1.DeleteOptions{})
	return err
}

func (c *client) ListApplication(tx interface{}, namespace string, listOptions *models.ListOptions) (*models.ApplicationList, error) {
	defer utils.Trace(c.log.Debug, "ListApplication")()
	list, err := c.customClient.CloudV1alpha1().Applications(namespace).List(*fromListOptionsModel(listOptions))
	listOptions.Continue = list.Continue
	if err != nil {
		return nil, err
	}
	res := toAppListModel(list)
	res.ListOptions = listOptions
	return res, err
}
