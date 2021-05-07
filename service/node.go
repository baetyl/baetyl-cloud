package service

import (
	"strings"
	"time"

	"github.com/baetyl/baetyl-go/v2/log"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/baetyl/baetyl-go/v2/utils"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

const KeyCheckResourceDependency = "checkResourceDependency"

type CheckResourceDependencyFunc func(ns, nodeName string) error

//go:generate mockgen -destination=../mock/service/node.go -package=service github.com/baetyl/baetyl-cloud/v2/service NodeService

const casRetryTimes = 3

// NodeService NodeService
type NodeService interface {
	Get(tx interface{}, namespace, name string) (*specV1.Node, error)
	List(namespace string, listOptions *models.ListOptions) (*models.NodeList, error)
	Count(namespace string) (map[string]int, error)
	Create(tx interface{}, amespace string, node *specV1.Node) (*specV1.Node, error)
	Update(namespace string, node *specV1.Node) (*specV1.Node, error)
	Delete(namespace, name string) error

	UpdateReport(namespace, name string, report specV1.Report) (*models.Shadow, error)
	UpdateDesire(namespace, name string, app *specV1.Application, f func(*models.Shadow, *specV1.Application)) (*models.Shadow, error)

	GetDesire(namespace, name string) (*specV1.Desire, error)

	UpdateNodeAppVersion(namespace string, app *specV1.Application) ([]string, error)
	DeleteNodeAppVersion(namespace string, app *specV1.Application) ([]string, error)

	GetNodeProperties(ns, name string) (*models.NodeProperties, error)
	UpdateNodeProperties(ns, name string, props *models.NodeProperties) (*models.NodeProperties, error)
	UpdateNodeMode(ns, name, mode string) error
}

type NodeServiceImpl struct {
	indexService  IndexService
	node          plugin.Node
	app           plugin.Application
	Shadow        plugin.Shadow
	SysAppService SystemAppService
	Hooks         map[string]interface{}
}

// NewNodeService NewNodeService
func NewNodeService(config *config.CloudConfig) (NodeService, error) {
	node, err := plugin.GetPlugin(config.Plugin.Resource)
	if err != nil {
		return nil, err
	}

	shadow, err := plugin.GetPlugin(config.Plugin.Shadow)
	if err != nil {
		return nil, err
	}

	app, err := plugin.GetPlugin(config.Plugin.Resource)
	if err != nil {
		return nil, err
	}

	system, err := NewSystemAppService(config)
	if err != nil {
		return nil, err
	}

	is, err := NewIndexService(config)
	if err != nil {
		return nil, err
	}

	return &NodeServiceImpl{
		indexService:  is,
		SysAppService: system,
		node:          node.(plugin.Node),
		Shadow:        shadow.(plugin.Shadow),
		app:           app.(plugin.Application),
		Hooks:         make(map[string]interface{}),
	}, nil
}

// Get get the node
func (n *NodeServiceImpl) Get(tx interface{}, namespace, name string) (*specV1.Node, error) {
	node, err := n.node.GetNode(tx, namespace, name)
	if err != nil && strings.Contains(err.Error(), "not found") {
		return nil, common.Error(common.ErrResourceNotFound, common.Field("type", "node"),
			common.Field("name", name))
	} else if err != nil {
		log.L().Error("get node failed", log.Error(err))
		return nil, err
	}

	shadow, err := n.Shadow.Get(tx, namespace, name)
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
func (n *NodeServiceImpl) Create(tx interface{}, namespace string, node *specV1.Node) (*specV1.Node, error) {
	res, err := n.node.CreateNode(tx, namespace, node)
	if err != nil {
		log.L().Error("create node failed", log.Error(err))
		return nil, err
	}

	_, err = n.SysAppService.GenApps(tx, namespace, node)
	if err != nil {
		return nil, err
	}

	if err = n.insertOrUpdateNodeAndAppIndex(tx, namespace, res, models.NewShadowFromNode(res), true); err != nil {
		return nil, err
	}
	return res, err
}

// Update update node
func (n *NodeServiceImpl) Update(namespace string, node *specV1.Node) (*specV1.Node, error) {
	res, err := n.node.UpdateNode(namespace, node)
	if err != nil {
		return nil, err
	}

	shadow, err := n.Shadow.Get(nil, namespace, node.Name)
	if err != nil {
		return nil, err
	}

	// delete indexes for node and apps
	if err := n.indexService.RefreshAppsIndexByNode(nil, namespace, res.Name, []string{}); err != nil {
		return nil, err
	}

	if err = n.insertOrUpdateNodeAndAppIndex(nil, namespace, res, shadow, false); err != nil {
		return nil, err
	}
	return res, nil
}

// List get list node
func (n *NodeServiceImpl) List(namespace string, listOptions *models.ListOptions) (*models.NodeList, error) {
	pageSize := listOptions.GetLimitNumber()
	if listOptions.NodeSelector != "" && pageSize > 0 {
		// in order to get all node data
		listOptions.PageSize = 0
	}
	list, err := n.node.ListNode(namespace, listOptions)
	if err != nil {
		return nil, err
	}
	shadowList, err := n.Shadow.List(namespace, list)
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
	if listOptions.NodeSelector != "" {
		listOptions.PageSize = pageSize
		return filterNodeListByNodeSelector(list), nil
	}
	return list, nil
}

// Count get current node number
func (n *NodeServiceImpl) Count(namespace string) (map[string]int, error) {
	list, err := n.List(namespace, &models.ListOptions{})
	if err != nil {
		return nil, err
	}
	return map[string]int{
		plugin.QuotaNode: len(list.Items),
	}, nil
}

// Delete delete node
func (n *NodeServiceImpl) Delete(namespace, name string) error {
	if check, ok := n.Hooks[KeyCheckResourceDependency].(CheckResourceDependencyFunc); ok {
		if err := check(namespace, name); err != nil {
			return err
		}
	}

	if err := n.node.DeleteNode(namespace, name); err != nil {
		return err
	}

	if err := n.Shadow.Delete(namespace, name); err != nil {
		common.LogDirtyData(err,
			log.Any("type", common.Shadow),
			log.Any("namespace", namespace),
			log.Any("name", name),
			log.Any("operation", "delete"))
	}

	if err := n.indexService.RefreshAppsIndexByNode(nil, namespace, name, []string{}); err != nil {
		common.LogDirtyData(err,
			log.Any("type", "app node index"),
			log.Any("namespace", namespace),
			log.Any("name", name),
			log.Any("operation", "delete"))
	}
	return nil
}

// UpdateReport Update Report
func (n *NodeServiceImpl) UpdateReport(namespace, name string, report specV1.Report) (*models.Shadow, error) {
	shadow, err := n.Shadow.Get(nil, namespace, name)
	if err != nil {
		return nil, err
	}

	if report != nil {
		report["time"] = time.Now().UTC()
	}

	if shadow == nil {
		_, err = n.node.GetNode(nil, namespace, name)
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
		// TODO refactor merge and remove this
		// since merge won't delete exist key-val, node info and stats should override
		if node, ok := report[common.NodeInfo]; ok {
			shadow.Report[common.NodeInfo] = node
		}
		if nodeStats, ok := report[common.NodeStats]; ok {
			shadow.Report[common.NodeStats] = nodeStats
		}
	}
	if err := n.updateReportNodeProperties(namespace, name, report, shadow); err != nil {
		return nil, err
	}
	return n.Shadow.UpdateReport(shadow)
}

func (n *NodeServiceImpl) updateReportNodeProperties(ns, name string, report specV1.Report, shad *models.Shadow) error {
	node, err := n.node.GetNode(nil, ns, name)
	if err != nil {
		return err
	}
	newProps := map[string]interface{}{}
	if props, ok := report[common.NodeProps].(map[string]interface{}); ok {
		newProps = props
	}
	oldProps := map[string]interface{}{}
	if props, ok := shad.Report[common.NodeProps].(map[string]interface{}); ok {
		oldProps = props
	}
	diff, err := specV1.Desire(newProps).DiffWithNil(oldProps)
	meta := getNodePropertiesMeta(node)
	now := time.Now().UTC()
	for key, val := range diff {
		if val != nil {
			meta.ReportMeta[key] = now
		} else {
			delete(meta.ReportMeta, key)
		}
	}
	updateNodePropertiesMeta(node, meta)
	if _, err := n.node.UpdateNode(ns, node); err != nil {
		return err
	}
	// since merge won't delete exist key-val, node props should override
	if len(newProps) == 0 {
		delete(shad.Report, common.NodeProps)
	} else {
		shad.Report[common.NodeProps] = newProps
	}
	return nil
}

// UpdateDesire Update Desire
// Parameter f can be RefreshNodeDesireByApp or DeleteNodeDesireByApp
func (n *NodeServiceImpl) UpdateDesire(namespace, name string, app *specV1.Application, f func(*models.Shadow, *specV1.Application)) (*models.Shadow, error) {
	// Retry times
	var count = 0
	for {
		newShadow, err := n.Shadow.Get(nil, namespace, name)
		if err != nil {
			return nil, err
		}

		if newShadow == nil {
			newShadow, err = n.createShadow(namespace, name, specV1.Desire{}, nil)
		}

		// Refresh desire in Shadow by app
		f(newShadow, app)
		updatedShadow, err := n.Shadow.UpdateDesire(nil, newShadow)
		if err == nil || err.Error() != common.ErrUpdateCas {
			return updatedShadow, err
		}
		count++
		if count >= casRetryTimes {
			break
		}
	}
	return nil, common.Error(common.ErrResourceConflict, common.Field("node", name), common.Field("type", "Shadow"))
}

func (n *NodeServiceImpl) updateDesire(tx interface{}, shadow *models.Shadow, desire specV1.Desire) (*models.Shadow, error) {
	if shadow.Desire == nil {
		shadow.Desire = desire
	} else {
		err := shadow.Desire.Merge(desire)
		if err != nil {
			return nil, err
		}
	}

	return n.Shadow.UpdateDesire(tx, shadow)
}

func (n *NodeServiceImpl) GetDesire(namespace, name string) (*specV1.Desire, error) {
	shadow, _ := n.Shadow.Get(nil, namespace, name)
	if shadow == nil {
		return nil, common.Error(
			common.ErrResourceNotFound,
			common.Field("type", "Shadow"),
			common.Field("name", name),
			common.Field("namespace", namespace))
	}
	return &shadow.Desire, nil
}

func (n *NodeServiceImpl) insertOrUpdateNodeAndAppIndex(tx interface{}, namespace string, node *specV1.Node, shadow *models.Shadow, flag bool) error {
	apps, err := n.app.ListApplication(tx, namespace, &models.ListOptions{})
	if err != nil {
		log.L().Error("list application error", log.Error(err))
		return err
	}

	desire, appNames := n.rematchApplicationsForNode(apps, node.Labels)

	node.Desire = desire

	if flag {
		shadow.Desire = desire
		if _, err = n.Shadow.Create(tx, shadow); err != nil {
			return err
		}
	} else {
		if _, err = n.updateDesire(tx, shadow, desire); err != nil {
			log.L().Error("update node desired node failed", log.Error(err))
			return err
		}
	}

	if err = n.indexService.RefreshAppsIndexByNode(tx, namespace, node.Name, appNames); err != nil {
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
func (n *NodeServiceImpl) rematchApplicationsForNode(apps *models.ApplicationList, labels map[string]string) (specV1.Desire, []string) {
	desireApps := make([]specV1.AppInfo, 0)
	sysApps := make([]specV1.AppInfo, 0)

	appNames := make([]string, 0)
	for _, app := range apps.Items {
		if app.Selector == "" {
			continue
		}

		if ok, err := utils.IsLabelMatch(app.Selector, labels); err == nil && ok {
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
func (n *NodeServiceImpl) UpdateNodeAppVersion(namespace string, app *specV1.Application) ([]string, error) {
	if app.Selector == "" {
		return nil, nil
	}

	// list nodes
	nodeList, err := n.node.ListNode(namespace, &models.ListOptions{LabelSelector: app.Selector})
	if err != nil {
		return nil, err
	}
	// update nodes
	var nodes []string
	for idx := range nodeList.Items {
		node := &nodeList.Items[idx]
		nodes = append(nodes, node.Name)
		_, err := n.UpdateDesire(node.Namespace, node.Name, app, RefreshNodeDesireByApp)
		if err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (n *NodeServiceImpl) createShadow(namespace, name string, desire specV1.Desire, report specV1.Report) (*models.Shadow, error) {
	shadow := models.NewShadow(namespace, name)

	if desire != nil {
		shadow.Desire = desire
	}

	if report != nil {
		shadow.Report = report
	}

	return n.Shadow.Create(nil, shadow)
}

// DeleteNodeAppVersion delete the node desire's appVersion for app deleted
func (n *NodeServiceImpl) DeleteNodeAppVersion(namespace string, app *specV1.Application) ([]string, error) {
	if app.Selector == "" {
		return nil, nil
	}

	// list nodes
	nodeList, err := n.node.ListNode(namespace, &models.ListOptions{LabelSelector: app.Selector})
	if err != nil {
		return nil, err
	}

	// update nodes
	var nodes []string
	for idx := range nodeList.Items {
		node := &nodeList.Items[idx]
		nodes = append(nodes, node.Name)

		_, err = n.UpdateDesire(node.Namespace, node.Name, app, DeleteNodeDesireByApp)
		if err != nil {
			return nil, err
		}
	}

	return nodes, nil
}

func RefreshNodeDesireByApp(shadow *models.Shadow, app *specV1.Application) {

	if shadow.Desire == nil {
		shadow.Desire = specV1.Desire{}
	}

	apps := shadow.Desire.AppInfos(app.System)

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

	shadow.Desire.SetAppInfos(app.System, apps)

}

func DeleteNodeDesireByApp(shadow *models.Shadow, app *specV1.Application) {
	if shadow.Desire == nil {
		return
	}
	appInfos := make([]specV1.AppInfo, 0)
	apps := shadow.Desire.AppInfos(app.System)

	for _, a := range apps {
		if a.Name != app.Name {
			appInfos = append(appInfos, a)
		}
	}
	shadow.Desire.SetAppInfos(app.System, appInfos)
}

func toShadowMap(shadowList *models.ShadowList) map[string]*models.Shadow {
	shadowMap := make(map[string]*models.Shadow)
	for idx := range shadowList.Items {
		shadow := shadowList.Items[idx]
		shadowMap[shadow.Name] = &shadow
	}

	return shadowMap
}

func (n *NodeServiceImpl) GetNodeProperties(namespace, name string) (*models.NodeProperties, error) {
	node, err := n.node.GetNode(nil, namespace, name)
	if err != nil && strings.Contains(err.Error(), "not found") {
		return nil, common.Error(common.ErrResourceNotFound, common.Field("type", "node"),
			common.Field("name", name))
	} else if err != nil {
		log.L().Error("get node failed", log.Error(err))
		return nil, err
	}
	shadow, err := n.Shadow.Get(nil, namespace, name)
	if err != nil {
		return nil, err
	}
	report := map[string]interface{}{}
	if props, ok := shadow.Report[common.NodeProps].(map[string]interface{}); ok {
		report = props
	}
	desire := map[string]interface{}{}
	if props, ok := shadow.Desire[common.NodeProps].(map[string]interface{}); ok {
		desire = props
	}
	meta := getNodePropertiesMeta(node)
	nodeProps := &models.NodeProperties{
		State: models.NodePropertiesState{
			Report: report,
			Desire: desire,
		},
		Meta: models.NodePropertiesMetadata{
			ReportMeta: meta.ReportMeta,
			DesireMeta: meta.DesireMeta,
		},
	}
	return nodeProps, nil
}

// UpdateNodeProperties update desire of node properties
// and can not update report of node properties
func (n *NodeServiceImpl) UpdateNodeProperties(namespace, name string, props *models.NodeProperties) (*models.NodeProperties, error) {
	node, err := n.node.GetNode(nil, namespace, name)
	if err != nil && strings.Contains(err.Error(), "not found") {
		return nil, common.Error(common.ErrResourceNotFound, common.Field("type", "node"),
			common.Field("name", name))
	} else if err != nil {
		log.L().Error("get node failed", log.Error(err))
		return nil, err
	}
	shadow, err := n.Shadow.Get(nil, namespace, name)
	if err != nil {
		return nil, err
	}
	oldDesire := map[string]interface{}{}
	if props, ok := shadow.Desire[common.NodeProps].(map[string]interface{}); ok {
		oldDesire = props
	}
	var newDesire specV1.Desire = props.State.Desire
	diff, err := newDesire.DiffWithNil(oldDesire)
	if err != nil {
		return nil, err
	}
	meta := getNodePropertiesMeta(node)
	now := time.Now().UTC()
	for key, val := range diff {
		meta.DesireMeta[key] = now
		if val == nil {
			delete(meta.DesireMeta, key)
		}
	}
	report := map[string]interface{}{}
	if props, ok := shadow.Report[common.NodeProps].(map[string]interface{}); ok {
		report = props
	}
	props.State.Report = report
	// cast to map[string]interface{} should not omit
	props.State.Desire = map[string]interface{}(newDesire)
	props.Meta.ReportMeta = meta.ReportMeta
	props.Meta.DesireMeta = meta.DesireMeta
	// cast to map[string]interface{} should not omit
	shadow.Desire[common.NodeProps] = map[string]interface{}(newDesire)
	_, err = n.Shadow.UpdateDesire(nil, shadow)
	if err != nil {
		return nil, err
	}
	updateNodePropertiesMeta(node, meta)
	if _, err := n.node.UpdateNode(namespace, node); err != nil {
		return nil, err
	}
	return props, nil
}

func (n *NodeServiceImpl) UpdateNodeMode(ns, name, mode string) error {
	node, err := n.node.GetNode(nil, ns, name)
	if err != nil && strings.Contains(err.Error(), "not found") {
		return common.Error(common.ErrResourceNotFound, common.Field("type", "node"),
			common.Field("name", name))
	} else if err != nil {
		log.L().Error("get node failed", log.Error(err))
		return err
	}
	if node.Attributes == nil {
		node.Attributes = map[string]interface{}{}
	}
	node.Attributes[specV1.KeySyncMode] = specV1.SyncMode(mode)
	_, err = n.node.UpdateNode(ns, node)
	if err != nil {
		return err
	}
	return nil
}

func getNodePropertiesMeta(node *specV1.Node) *models.NodePropertiesMetadata {
	propsMeta := &models.NodePropertiesMetadata{
		ReportMeta: make(map[string]interface{}),
		DesireMeta: make(map[string]interface{}),
	}
	if node == nil || node.Attributes == nil {
		return propsMeta
	}
	if meta, ok := node.Attributes[common.ReportMeta].(map[string]interface{}); ok {
		propsMeta.ReportMeta = meta
	}
	if meta, ok := node.Attributes[common.DesireMeta].(map[string]interface{}); ok {
		propsMeta.DesireMeta = meta
	}
	return propsMeta
}

func updateNodePropertiesMeta(node *specV1.Node, meta *models.NodePropertiesMetadata) {
	if meta == nil || node == nil {
		return
	}
	if node.Attributes == nil {
		node.Attributes = make(map[string]interface{})
	}
	if meta.ReportMeta != nil {
		node.Attributes[common.ReportMeta] = meta.ReportMeta
	}
	if meta.DesireMeta != nil {
		node.Attributes[common.DesireMeta] = meta.DesireMeta
	}
}

func filterNodeListByNodeSelector(list *models.NodeList) *models.NodeList {
	// filter nodes according to nodeSelector
	items := []specV1.Node{}
	for _, item := range list.Items {
		if item.Report == nil || len(item.Report) == 0 {
			continue
		}
		clusterVar, ok := item.Report[common.NodeInfo]
		if !ok {
			continue
		}
		cluster, ok := clusterVar.(map[string]interface{})
		if !ok {
			continue
		}
		for _, nodeVar := range cluster {
			node, ok := nodeVar.(map[string]interface{})
			if !ok {
				continue
			}
			labelVar, ok := node["labels"]
			if !ok {
				continue
			}
			labels, ok := labelVar.(map[string]interface{})
			if !ok {
				continue
			}
			ls := map[string]string{}
			for k, v := range labels {
				ls[k] = v.(string)
			}
			if ok, err := utils.IsLabelMatch(list.ListOptions.NodeSelector, ls); err != nil || !ok {
				continue
			}
			items = append(items, item)
			break
		}
	}
	start, end := models.GetPagingParam(list.ListOptions, len(items))
	return &models.NodeList{
		Total:       len(items),
		ListOptions: list.ListOptions,
		Items:       items[start:end],
	}
}
