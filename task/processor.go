package task

import (
	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
	"github.com/baetyl/baetyl-cloud/v2/service"
	"strings"
)

const (
	DeleteNamespace             = "delete_namespace"
	DeleteNode                  = "delete_node"
	DeleteApp                   = "delete_app"
	DeleteSecret                = "delete_secret"
	DeleteConfig                = "delete_config"
	DeleteQuotaByNamespace      = "delete_quota_by_namespace"
	DeleteIndexByNamespace      = "delete_index_by_namespace"
	DeleteAppHistoryByNamespace = "delete_app_his_by_namespace"
)

func RegisterNamespaceProcessor(cfg *config.CloudConfig) error {
	processor, err := NewNamespaceProcessor(cfg)
	if err != nil {
		return err
	}

	plugin.TaskRegister.AddTask(common.TaskNamespaceDelete, DeleteNamespace, processor.DeleteNamespace)
	plugin.TaskRegister.AddTask(common.TaskNamespaceDelete, DeleteNode, processor.DeleteNodesByNamespace)
	plugin.TaskRegister.AddTask(common.TaskNamespaceDelete, DeleteApp, processor.DeleteAppsByNamespace)
	plugin.TaskRegister.AddTask(common.TaskNamespaceDelete, DeleteSecret, processor.DeleteSecretsByNamespace)
	plugin.TaskRegister.AddTask(common.TaskNamespaceDelete, DeleteConfig, processor.DeleteConfigsByNamespace)
	plugin.TaskRegister.AddTask(common.TaskNamespaceDelete, DeleteQuotaByNamespace, processor.DeleteQuotaByNamespace)

	// must before DeleteIndexByNamespace
	plugin.TaskRegister.AddTask(common.TaskNamespaceDelete, DeleteAppHistoryByNamespace, processor.DeleteAppsHisByNamespace)

	plugin.TaskRegister.AddTask(common.TaskNamespaceDelete, DeleteIndexByNamespace, processor.DeleteIndexByNamespace)

	return nil
}

type namespaceProcessor struct {
	indexService     service.IndexService
	namespaceService service.NamespaceService
	resourceService  *service.AppCombinedService
	lisenceService   service.LicenseService
	nodeService      service.NodeService
}

func NewNamespaceProcessor(cfg *config.CloudConfig) (*namespaceProcessor, error) {
	is, err := service.NewIndexService(cfg)

	if err != nil {
		return nil, err
	}

	rs, err := service.NewAppCombinedService(cfg)
	if err != nil {
		return nil, err
	}

	ns, err := service.NewNamespaceService(cfg)
	if err != nil {
		return nil, err
	}

	ls, err := service.NewLicenseService(cfg)
	if err != nil {
		return nil, err
	}

	nodeSvc, err := service.NewNodeService(cfg)
	if err != nil {
		return nil, err
	}
	return &namespaceProcessor{
		indexService:     is,
		resourceService:  rs,
		namespaceService: ns,
		lisenceService:   ls,
		nodeService:      nodeSvc,
	}, nil
}

func (n *namespaceProcessor) DeleteNamespace(task *models.Task) error {
	return n.namespaceService.Delete(&models.Namespace{
		Name: task.Namespace,
	})
}

func (n *namespaceProcessor) DeleteQuotaByNamespace(task *models.Task) error {
	return n.lisenceService.DeleteQuotaByNamespace(task.Namespace)
}

func (n *namespaceProcessor) DeleteIndexByNamespace(task *models.Task) error {

	err := n.indexService.DeleteConfigsAndAppsIndexByNamespace(task.Namespace)
	if err != nil {
		return err
	}

	err = n.indexService.DeleteNodesAndAppsIndexByNamespace(task.Namespace)
	if err != nil {
		return err
	}

	err = n.indexService.DeleteSecretsAndAppsIndexByNamespace(task.Namespace)
	return err
}

func (n *namespaceProcessor) DeleteNodesByNamespace(task *models.Task) error {
	list, err := n.nodeService.List(task.Namespace, &models.ListOptions{})
	if err != nil {
		return err
	}

	for _, node := range list.Items {
		err = n.nodeService.Delete(node.Namespace, node.Name)

		if err != nil && !strings.Contains(err.Error(), "not found") {
			return err
		}

	}

	return nil
}

func (n *namespaceProcessor) DeleteAppsByNamespace(task *models.Task) error {
	list, err := n.resourceService.App.List(task.Namespace, &models.ListOptions{})
	if err != nil {
		return err
	}

	for _, app := range list.Items {
		err = n.resourceService.App.Delete(app.Namespace, app.Name, "")

		if err != nil && !strings.Contains(err.Error(), "not found") {
			return err
		}

	}

	return nil
}

func (n *namespaceProcessor) DeleteSecretsByNamespace(task *models.Task) error {
	list, err := n.resourceService.Secret.List(task.Namespace, &models.ListOptions{})
	if err != nil {
		return err
	}

	for _, secret := range list.Items {
		err = n.resourceService.Secret.Delete(secret.Namespace, secret.Name)

		if err != nil && !strings.Contains(err.Error(), "not found") {
			return err
		}
	}

	return nil
}

func (n *namespaceProcessor) DeleteConfigsByNamespace(task *models.Task) error {
	list, err := n.resourceService.Config.List(task.Namespace, &models.ListOptions{})
	if err != nil {
		return err
	}

	for _, config := range list.Items {
		err = n.resourceService.Config.Delete(config.Namespace, config.Name)

		if err != nil && !strings.Contains(err.Error(), "not found") {
			return err
		}
	}

	return nil
}

func (n *namespaceProcessor) DeleteAppsHisByNamespace(task *models.Task) error {
	appNames, err := n.indexService.ListAppsByNamespace(task.Namespace)
	if err != nil {
		return err
	}

	for _, appName := range appNames {
		err = n.resourceService.App.DeleteAppHis(task.Namespace, appName)

		if err != nil {
			return err
		}

	}
	return nil
}
