package service

import (
	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

//go:generate mockgen -destination=../mock/service/index.go -package=service github.com/baetyl/baetyl-cloud/v2/service IndexService

type IndexService interface {
	// comment
	RefreshIndex(namespace string, keyA, keyB common.Resource, valueA string, valueBs []string) error
	ListIndex(namespace string, keyA, byKeyB common.Resource, valueB string) ([]string, error)
	// app and config
	RefreshAppIndexByConfig(namespace, config string, apps []string) error
	RefreshConfigIndexByApp(namespace, app string, configs []string) error
	ListAppIndexByConfig(namespace, config string) ([]string, error)
	ListConfigIndexByApp(namespace, app string) ([]string, error)

	ListNodesByApp(namespace, app string) ([]string, error)
	ListAppsByNode(namespace, node string) ([]string, error)
	ListAppIndexBySecret(namespace, secret string) ([]string, error)

	ListNodesByNamespace(namespace string) ([]string, error)
	ListAppsByNamespace(namespace string) ([]string, error)

	// app and secret
	RefreshSecretIndexByApp(namespace, app string, secrets []string) error
	RefreshNodesIndexByApp(namespace, appName string, nodes []string) error
	RefreshAppsIndexByNode(namespace, node string, apps []string) error

	DeleteNodesAndAppsIndexByNamespace(namespace string) error
	DeleteConfigsAndAppsIndexByNamespace(namespace string) error
	DeleteSecretsAndAppsIndexByNamespace(namespace string) error
}

type indexService struct {
	index plugin.Index
}

// NewIndexService New Index Service
func NewIndexService(config *config.CloudConfig) (IndexService, error) {
	index, err := plugin.GetPlugin(config.Plugin.Index)
	if err != nil {
		return nil, err
	}
	return &indexService{index: index.(plugin.Index)}, nil
}

func (i *indexService) RefreshIndex(namespace string, keyA, keyB common.Resource, valueA string, valueBs []string) error {
	return i.index.RefreshIndex(namespace, keyA, keyB, valueA, valueBs)
}

func (i *indexService) ListIndex(namespace string, keyA, byKeyB common.Resource, valueB string) ([]string, error) {
	return i.index.ListIndex(namespace, keyA, byKeyB, valueB)
}

// helper
// app and config
func (i *indexService) RefreshAppIndexByConfig(namespace, config string, apps []string) error {
	return i.RefreshIndex(namespace, common.Config, common.Application, config, apps)
}

func (i *indexService) RefreshConfigIndexByApp(namespace, app string, configs []string) error {
	return i.RefreshIndex(namespace, common.Application, common.Config, app, configs)
}

func (i *indexService) ListAppIndexByConfig(namespace, config string) ([]string, error) {
	return i.ListIndex(namespace, common.Application, common.Config, config)
}

func (i *indexService) ListConfigIndexByApp(namespace, app string) ([]string, error) {
	return i.ListIndex(namespace, common.Config, common.Application, app)
}

func (i *indexService) RefreshNodesIndexByApp(namespace, appName string, nodes []string) error {
	return i.RefreshIndex(namespace, common.Application, common.Node, appName, nodes)
}

func (i *indexService) RefreshAppsIndexByNode(namespace, node string, apps []string) error {
	return i.RefreshIndex(namespace, common.Node, common.Application, node, apps)
}

func (i *indexService) ListNodesByApp(namespace, app string) ([]string, error) {
	return i.ListIndex(namespace, common.Node, common.Application, app)
}

func (i *indexService) ListAppsByNode(namespace, node string) ([]string, error) {
	return i.ListIndex(namespace, common.Application, common.Node, node)
}

// secret && apps
func (i *indexService) RefreshSecretIndexByApp(namespace, app string, secrets []string) error {
	return i.RefreshIndex(namespace, common.Application, common.Secret, app, secrets)
}

func (i *indexService) ListAppIndexBySecret(namespace, secret string) ([]string, error) {
	return i.ListIndex(namespace, common.Application, common.Secret, secret)
}

func (i *indexService) DeleteNodesAndAppsIndexByNamespace(namespace string) error {
	_, err := i.index.DeleteIndexByNamespace(namespace, common.Node, common.Application)
	return err
}
func (i *indexService) DeleteConfigsAndAppsIndexByNamespace(namespace string) error {
	_, err := i.index.DeleteIndexByNamespace(namespace, common.Config, common.Application)
	return err
}
func (i *indexService) DeleteSecretsAndAppsIndexByNamespace(namespace string) error {
	_, err := i.index.DeleteIndexByNamespace(namespace, common.Secret, common.Application)
	return err
}

func (i *indexService) ListNodesByNamespace(namespace string) ([]string, error) {
	return i.index.ListResourcesByNamespace(namespace, common.Node, common.Application)
}

func (i *indexService) ListAppsByNamespace(namespace string) ([]string, error) {
	return i.index.ListResourcesByNamespace(namespace, common.Application, common.Config)
}
