package service

import (
	"strings"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
	"github.com/baetyl/baetyl-go/v2/log"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"time"
)

//go:generate mockgen -destination=../mock/service/node.go -package=service github.com/baetyl/baetyl-cloud/v2/service NodeService

// NodeService NodeService
type NodeService interface {
	Get(namespace, name string) (*specV1.Node, error)
	List(namespace string, listOptions *models.ListOptions) (*models.NodeList, error)
	Create(namespace string, node *specV1.Node) (*specV1.Node, error)
	Update(namespace string, node *specV1.Node) (*specV1.Node, error)
	Delete(namespace, name string) error

	UpdateReport(namespace, name string, report specV1.Report) (*models.Shadow, error)
	UpdateDesire(namespace, name string, desire specV1.Desire) (*models.Shadow, error)

	UpdateNodeAppVersion(namespace string, app *specV1.Application) ([]string, error)
	DeleteNodeAppVersion(namespace string, app *specV1.Application) ([]string, error)
}

type nodeService struct {
	storage      plugin.ModelStorage
	indexService IndexService
	shadow       plugin.Shadow
}

// NewNodeService NewNodeService
func NewNodeService(config *config.CloudConfig) (NodeService, error) {
	ms, err := plugin.GetPlugin(config.Plugin.ModelStorage)
	if err != nil {
		log.L().Error("get storage plugin failed", log.Error(err))
		return nil, err
	}

	shadow, err := plugin.GetPlugin(config.Plugin.Shadow)
	if err != nil {
		return nil, err
	}

	is, err := NewIndexService(config)
	if err != nil {
		return nil, err
	}

	return &nodeService{
		storage:      ms.(plugin.ModelStorage),
		indexService: is,
		shadow:       shadow.(plugin.Shadow),
	}, nil
}

// Get get the node
func (n *nodeService) Get(namespace, name string) (*specV1.Node, error) {
	node, err := n.storage.GetNode(namespace, name)
	if err != nil && strings.Contains(err.Error(), "not found") {
		return nil, common.Error(common.ErrResourceNotFound, common.Field("type", "node"),
			common.Field("name", name))
	} else if err != nil {
		log.L().Error("get node failed", log.Error(err))
		return nil, err
	}

	shadow, err := n.shadow.Get(namespace, name)
	if err != nil {
		return nil, err
	}
	if shadow != nil {
		node.Report = shadow.Report
		node.Desire = shadow.Desire
	}

	return node, nil
}

// Create create a node
func (n *nodeService) Create(namespace string, node *specV1.Node) (*specV1.Node, error) {
	res, err := n.storage.CreateNode(namespace, node)
	if err != nil {
		log.L().Error("create node failed", log.Error(err))
		return nil, err
	}

	_, err = n.shadow.Create(models.NewShadowFromNode(res))
	if err != nil {
		return nil, err
	}

	if err = n.updateNodeAndAppIndex(namespace, res); err != nil {
		return nil, err
	}
	return res, err
}

// Update update node
func (n *nodeService) Update(namespace string, node *specV1.Node) (*specV1.Node, error) {
	res, err := n.storage.UpdateNode(namespace, node)
	if err != nil {
		return nil, err
	}

	// delete indexes for node and apps
	if err := n.indexService.RefreshAppsIndexByNode(namespace, res.Name, []string{}); err != nil {
		return nil, err
	}

	if err = n.updateNodeAndAppIndex(namespace, res); err != nil {
		return nil, err
	}
	return res, nil
}

// List get list node
func (n *nodeService) List(namespace string, listOptions *models.ListOptions) (*models.NodeList, error) {
	list, err := n.storage.ListNode(namespace, listOptions)
	if err != nil {
		return nil, err
	}
	shadowList, err := n.shadow.List(namespace, list)
	if err != nil {
		return nil, err
	}
	shadowMap := toShadowMap(shadowList)

	for idx := range list.Items {
		node := &list.Items[idx]
		if shadow, ok := shadowMap[node.Name]; ok {
			node.Desire = shadow.Desire
			node.Report = shadow.Report
		}
	}

	return list, nil
}

// Delete delete node
func (n *nodeService) Delete(namespace, name string) error {
	if err := n.storage.DeleteNode(namespace, name); err != nil {
		return err
	}

	if err := n.shadow.Delete(namespace, name); err != nil {
		common.LogDirtyData(err,
			log.Any("type", common.Shadow),
			log.Any("namespace", namespace),
			log.Any("name", name),
			log.Any("operation", "delete"))
	}

	if err := n.indexService.RefreshAppsIndexByNode(namespace, name, []string{}); err != nil {
		common.LogDirtyData(err,
			log.Any("type", "app node index"),
			log.Any("namespace", namespace),
			log.Any("name", name),
			log.Any("operation", "delete"))
	}
	return nil
}

// UpdateReport Update Report
func (n *nodeService) UpdateReport(namespace, name string, report specV1.Report) (*models.Shadow, error) {
	shadow, err := n.shadow.Get(namespace, name)
	if err != nil {
		return nil, err
	}

	if report != nil {
		report["time"] = time.Now().UTC()
	}

	if shadow == nil {
		_, err = n.storage.GetNode(namespace, name)
		if err != nil {
			return nil, err
		}
		return n.createShadow(namespace, name, nil, report)
	}

	if shadow.Report == nil {
		shadow.Report = report
	} else {
		err = shadow.Report.Merge(report)
		if err != nil {
			return nil, err
		}
	}

	return n.shadow.UpdateReport(shadow)
}

// UpdateDesire Update Desire
func (n *nodeService) UpdateDesire(namespace, name string, desire specV1.Desire) (*models.Shadow, error) {
	shadow, err := n.shadow.Get(namespace, name)
	if err != nil {
		return nil, err
	}

	if shadow == nil {
		return n.createShadow(namespace, name, desire, nil)
	}

	if shadow.Desire == nil {
		shadow.Desire = desire
	} else {
		err = shadow.Desire.Merge(desire)
		if err != nil {
			return nil, err
		}
	}

	return n.shadow.UpdateDesire(shadow)
}

func (n *nodeService) updateNodeAndAppIndex(namespace string, node *specV1.Node) error {
	apps, err := n.storage.ListApplication(namespace, &models.ListOptions{})
	if err != nil {
		log.L().Error("list application error", log.Error(err))
		return err
	}

	desire, appNames := n.rematchApplicationsForNode(apps, node.Labels)

	node.Desire = desire

	if _, err = n.UpdateDesire(node.Namespace, node.Name, desire); err != nil {
		log.L().Error("update node desired node failed", log.Error(err))
		return err
	}

	if err = n.indexService.RefreshAppsIndexByNode(namespace, node.Name, appNames); err != nil {
		log.L().Error("refresh app index by node failed", log.Error(err))
		return err
	}
	return nil
}

// rematchApplicationsForNode rematch applications for node
//   - param apps: all applications for the namespace
//   - param nodeLabels: the labels of node
//   - return desire: the node's desire
//   - return appNames: matched application names
func (n *nodeService) rematchApplicationsForNode(apps *models.ApplicationList, labels map[string]string) (specV1.Desire, []string) {
	desireApps := make([]specV1.AppInfo, 0)
	sysApps := make([]specV1.AppInfo, 0)

	appNames := make([]string, 0)
	for _, app := range apps.Items {
		if app.Selector == "" {
			continue
		}

		if ok, err := n.storage.IsLabelMatch(app.Selector, labels); err == nil && ok {
			if app.System {
				sysApps = append(sysApps, specV1.AppInfo{
					Name:    app.Name,
					Version: app.Version,
				})
			} else {
				desireApps = append(desireApps, specV1.AppInfo{
					Name:    app.Name,
					Version: app.Version,
				})
			}
			appNames = append(appNames, app.Name)
		}
	}

	return specV1.Desire{
		common.DesiredApplications:    desireApps,
		common.DesiredSysApplications: sysApps,
	}, appNames

}

// UpdateNodeAppVersion update the node desire's appVersion for app changed
func (n *nodeService) UpdateNodeAppVersion(namespace string, app *specV1.Application) ([]string, error) {
	if app.Selector == "" {
		return nil, nil
	}

	// list nodes
	nodeList, err := n.List(namespace, &models.ListOptions{LabelSelector: app.Selector})
	if err != nil {
		return nil, err
	}
	// update nodes
	var nodes []string
	for idx := range nodeList.Items {
		node := &nodeList.Items[idx]
		nodes = append(nodes, node.Name)
		refreshNodeDesireByApp(node, app)
		_, err := n.UpdateDesire(namespace, node.Name, node.Desire)
		if err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (n *nodeService) createShadow(namespace, name string, desire specV1.Desire, report specV1.Report) (*models.Shadow, error) {
	shadow := models.NewShadow(namespace, name)

	if desire != nil {
		shadow.Desire = desire
	}

	if report != nil {
		shadow.Report = report
	}

	return n.shadow.Create(shadow)
}

// DeleteNodeAppVersion delete the node desire's appVersion for app deleted
func (n *nodeService) DeleteNodeAppVersion(namespace string, app *specV1.Application) ([]string, error) {
	if app.Selector == "" {
		return nil, nil
	}

	// list nodes
	nodeList, err := n.List(namespace, &models.ListOptions{LabelSelector: app.Selector})
	if err != nil {
		return nil, err
	}

	// update nodes
	var nodes []string
	for idx := range nodeList.Items {
		node := &nodeList.Items[idx]
		nodes = append(nodes, node.Name)
		appInfos := make([]specV1.AppInfo, 0)

		if node.Desire != nil {
			apps := node.Desire.AppInfos(app.System)

			for _, a := range apps {
				if a.Name != app.Name {
					appInfos = append(appInfos, a)
				}
			}
			node.Desire.SetAppInfos(app.System, appInfos)

			_, err = n.UpdateDesire(namespace, node.Name, node.Desire)
			if err != nil {
				return nil, err
			}
		}
	}

	return nodes, nil
}

func refreshNodeDesireByApp(node *specV1.Node, app *specV1.Application) {

	if node.Desire == nil {
		node.Desire = specV1.Desire{}
	}

	apps := node.Desire.AppInfos(app.System)

	if apps == nil {
		apps = make([]specV1.AppInfo, 0)
	}

	idx := -1
	for i, a := range apps {
		if a.Name == app.Name {
			idx = i
		}
	}

	// add new app
	if idx == -1 {
		apps = append(apps, specV1.AppInfo{
			Name:    app.Name,
			Version: app.Version,
		})
	} else {
		// modified the old app
		apps[idx].Version = app.Version
	}

	node.Desire.SetAppInfos(app.System, apps)

}

func toShadowMap(shadowList *models.ShadowList) map[string]*models.Shadow {
	shadowMap := make(map[string]*models.Shadow)
	for idx := range shadowList.Items {
		shadow := shadowList.Items[idx]
		shadowMap[shadow.Name] = &shadow
	}

	return shadowMap
}
