package task

import (
	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/service"
)

const (
	DeleteNamespace        = "delete_namespace"
	DeleteQuotaByNamespace = "delete_quota_by_namespace"
	DeleteIndexByNamespace = "delete_index_by_namespace"
)

func RegisterNamespaceProcessor(cfg *config.CloudConfig) error {
	processor, err := NewNamespaceProcessor(cfg)
	if err != nil {
		return err
	}

	TaskRegister.Register(common.TaskNamespaceDelete, DeleteNamespace, processor.DeleteNamespace)
	TaskRegister.Register(common.TaskNamespaceDelete, DeleteQuotaByNamespace, processor.DeleteQuotaByNamespace)
	TaskRegister.Register(common.TaskNamespaceDelete, DeleteIndexByNamespace, processor.DeleteIndexByNamespace)

	return nil
}

type namespaceProcessor struct {
	indexService     service.IndexService
	namespaceService service.NamespaceService
	resourceService  *service.AppCombinedService
	licenceService   service.LicenseService
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
	return &namespaceProcessor{
		indexService:     is,
		resourceService:  rs,
		namespaceService: ns,
		licenceService:   ls,
	}, nil
}

func (n *namespaceProcessor) DeleteNamespace(task *models.Task) error {
	return n.namespaceService.Delete(&models.Namespace{
		Name: task.Namespace,
	})
}

func (n *namespaceProcessor) DeleteQuotaByNamespace(task *models.Task) error {
	return n.licenceService.DeleteQuotaByNamespace(task.Namespace)
}

func (n *namespaceProcessor) DeleteIndexByNamespace(task *models.Task) error {
	apps, err := n.resourceService.App.List(task.Namespace, &models.ListOptions{})
	if err != nil {
		return err
	}

	for _, v := range apps.Items {
		err = n.indexService.RefreshConfigIndexByApp(task.Namespace, v.Name, []string{})
		if err != nil {
			return err
		}

		err = n.indexService.RefreshNodesIndexByApp(task.Namespace, v.Name, []string{})
		if err != nil {
			return err
		}

		err = n.indexService.RefreshSecretIndexByApp(task.Namespace, v.Name, []string{})
		if err != nil {
			return err
		}
	}
	return nil
}
