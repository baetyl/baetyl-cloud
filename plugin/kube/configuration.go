package kube

import (
	"fmt"
	"time"

	"github.com/baetyl/baetyl-go/v2/utils"

	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/jinzhu/copier"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin/kube/apis/cloud/v1alpha1"
)

func toConfigurationModel(config *v1alpha1.Configuration) *specV1.Configuration {
	res := &specV1.Configuration{Version: config.ObjectMeta.ResourceVersion}
	err := copier.Copy(res, config)
	if err != nil {
		panic(fmt.Sprintf("copier exception: %s", err.Error()))
	}
	if desc, ok := config.Annotations[common.AnnotationDescription]; ok {
		res.Description = desc
	}
	if us, ok := config.Annotations[common.AnnotationUpdateTimestamp]; ok {
		res.UpdateTimestamp, _ = time.Parse(common.TimeFormat, us)
	}
	res.CreationTimestamp = config.CreationTimestamp.Time.UTC()

	return res
}

func toConfigurationListModel(configList *v1alpha1.ConfigurationList) *models.ConfigurationList {
	res := &models.ConfigurationList{
		Items: make([]specV1.Configuration, 0),
	}
	for _, item := range configList.Items {
		ptr := toConfigurationModel(&item)
		res.Items = append(res.Items, *ptr)
	}
	res.Total = len(configList.Items)
	return res
}

func fromConfigurationModel(config *specV1.Configuration) *v1alpha1.Configuration {
	res := &v1alpha1.Configuration{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Configuration",
			APIVersion: v1alpha1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:            config.Name,
			Namespace:       config.Namespace,
			Annotations:     map[string]string{},
			ResourceVersion: config.Version,
		},
	}
	err := copier.Copy(res, config)
	res.Annotations[common.AnnotationDescription] = config.Description
	res.Annotations[common.AnnotationUpdateTimestamp] = config.UpdateTimestamp.UTC().Format(common.TimeFormat)

	if err != nil {
		panic(fmt.Sprintf("copier exception: %s", err.Error()))
	}
	return res
}

func (c *client) GetConfig(tx interface{}, namespace, name, version string) (*specV1.Configuration, error) {
	options := metav1.GetOptions{ResourceVersion: version}
	defer utils.Trace(c.log.Debug, "GetConfig")()
	config, err := c.customClient.CloudV1alpha1().Configurations(namespace).Get(c.ctx, name, options)
	if err != nil {
		return nil, err
	}
	return toConfigurationModel(config), nil
}

func (c *client) CreateConfig(tx interface{}, namespace string, configModel *specV1.Configuration) (*specV1.Configuration, error) {
	configModel.UpdateTimestamp = time.Now()
	defer utils.Trace(c.log.Debug, "CreateConfig")()
	config, err := c.customClient.CloudV1alpha1().
		Configurations(namespace).
		Create(c.ctx, fromConfigurationModel(configModel), metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return toConfigurationModel(config), err
}

func (c *client) UpdateConfig(tx interface{}, namespace string, configurationModel *specV1.Configuration) (*specV1.Configuration, error) {
	defer utils.Trace(c.log.Debug, "UpdateConfig")()
	configuration, err := c.customClient.CloudV1alpha1().
		Configurations(namespace).
		Update(c.ctx, fromConfigurationModel(configurationModel), metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}
	return toConfigurationModel(configuration), err
}

func (c *client) DeleteConfig(tx interface{}, namespace, name string) error {
	defer utils.Trace(c.log.Debug, "DeleteConfig")()
	return c.customClient.CloudV1alpha1().Configurations(namespace).Delete(c.ctx, name, metav1.DeleteOptions{})
}

func (c *client) ListConfig(namespace string, listOptions *models.ListOptions) (*models.ConfigurationList, error) {
	defer utils.Trace(c.log.Debug, "ListConfig")()
	list, err := c.customClient.CloudV1alpha1().Configurations(namespace).List(c.ctx, *fromListOptionsModel(listOptions))
	if err != nil {
		return nil, err
	}
	res := toConfigurationListModel(list)
	res.ListOptions = listOptions
	return res, err
}
