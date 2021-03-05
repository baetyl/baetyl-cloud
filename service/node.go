package service

import (
	"strings"
	"time"

	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/baetyl/baetyl-go/v2/log"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/baetyl/baetyl-go/v2/utils"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

//go:generate mockgen -destination=../mock/service/node.go -package=service github.com/baetyl/baetyl-cloud/v2/service NodeService

const casRetryTimes = 3

// NodeService NodeService
type NodeService interface {
	Get(namespace, name string) (*specV1.Node, error)
	List(namespace string, listOptions *models.ListOptions) (*models.NodeList, error)
	Count(namespace string) (map[string]int, error)
	Create(namespace string, node *specV1.Node) (*specV1.Node, error)
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

type nodeService struct {
	indexService IndexService
	node         plugin.Node
	shadow       plugin.Shadow
	app          plugin.Application
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

	is, err := NewIndexService(config)
	if err != nil {
		return nil, err
	}

	return &nodeService{
		indexService: is,
		node:         node.(plugin.Node),
		shadow:       shadow.(plugin.Shadow),
		app:          app.(plugin.Application),
	}, nil
}

// Get get the node
func (n *nodeService) Get(namespace, name string) (*specV1.Node, error) {
	node, err := n.node.GetNode(namespace, name)
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
	res, err := n.node.CreateNode(namespace, node)
	if err != nil {
		log.L().Error("create node failed", log.Error(err))
		return nil, err
	}

	shadow, err := n.shadow.Create(models.NewShadowFromNode(res))
	if err != nil {
		return nil, err
	}

	if err = n.updateNodeAndAppIndex(namespace, res, shadow); err != nil {
		return nil, err
	}
	return res, err
}

// Update update node
func (n *nodeService) Update(namespace string, node *specV1.Node) (*specV1.Node, error) {
	res, err := n.node.UpdateNode(namespace, node)
	if err != nil {
		return nil, err
	}

	shadow, err := n.shadow.Get(namespace, node.Name)
	if err != nil {
		return nil, err
	}

	// delete indexes for node and apps
	if err := n.indexService.RefreshAppsIndexByNode(namespace, res.Name, []string{}); err != nil {
		return nil, err
	}

	if err = n.updateNodeAndAppIndex(namespace, res, shadow); err != nil {
		return nil, err
	}
	return res, nil
}

// List get list node
func (n *nodeService) List(namespace string, listOptions *models.ListOptions) (*models.NodeList, error) {
	list, err := n.node.ListNode(namespace, listOptions)
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

// Count get current node number
func (n *nodeService) Count(namespace string) (map[string]int, error) {
	list, err := n.List(namespace, &models.ListOptions{})
	if err != nil {
		return nil, err
	}
	return map[string]int{
		plugin.QuotaNode: len(list.Items),
	}, nil
}

// Delete delete node
func (n *nodeService) Delete(namespace, name string) error {
	if err := n.node.DeleteNode(namespace, name); err != nil {
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
		_, err = n.node.GetNode(namespace, name)
		if err != nil {
			return nil, err
		}
		return n.createShadow(namespace, name, nil, report)
	}

	// update node props meta
	node, err := n.node.GetNode(namespace, name)
	if err != nil {
		return nil, err
	}
	meta, err := getNodePropertiesMeta(node)
	if err != nil {
		return nil, err
	}
	newPropsReport := map[string]interface{}{}
	if props, ok := report[common.NodeProps]; ok && props != nil {
		newPropsReport, ok = props.(map[string]interface{})
		if !ok {
			return nil, errors.New("invalid report of node properties")
		}
	}
	oldPropsReport := map[string]interface{}{}
	if props, ok := shadow.Report[common.NodeProps]; ok && props != nil {
		oldPropsReport, ok = props.(map[string]interface{})
		if !ok {
			return nil, errors.New("invalid report of node properties")
		}
	}
	diff, err := specV1.Desire(newPropsReport).DiffWithNil(oldPropsReport)
	now := time.Now().UTC()
	for key, val := range diff {
		meta.ReportMeta[key] = now
		if val == nil {
			delete(meta.ReportMeta, key)
		}
	}
	updateNodePropertiesMeta(node, meta)
	if _, err := n.node.UpdateNode(namespace, node); err != nil {
		return nil, err
	}

	if shadow.Report == nil {
		shadow.Report = report
	} else {
		err = shadow.Report.Merge(report)
		if err != nil {
			return nil, err
		}
		// TODO refactor merge and remove this
		// since merge won't delete exist key-val, node props should override
		if len(newPropsReport) == 0 {
			delete(shadow.Report, common.NodeProps)
		} else {
			shadow.Report[common.NodeProps] = newPropsReport
		}
	}
	return n.shadow.UpdateReport(shadow)
}

// UpdateDesire Update Desire
// Parameter f can be RefreshNodeDesireByApp or DeleteNodeDesireByApp
func (n *nodeService) UpdateDesire(namespace, name string, app *specV1.Application, f func(*models.Shadow, *specV1.Application)) (*models.Shadow, error) {
	// Retry times
	var count = 0
	for {
		newShadow, err := n.shadow.Get(namespace, name)
		if err != nil {
			return nil, err
		}

		if newShadow == nil {
			newShadow, err = n.createShadow(namespace, name, specV1.Desire{}, nil)
		}

		// Refresh desire in shadow by app
		f(newShadow, app)
		updatedShadow, err := n.shadow.UpdateDesire(newShadow)
		if err == nil || err.Error() != common.ErrUpdateCas {
			return updatedShadow, err
		}
		count++
		if count >= casRetryTimes {
			break
		}
	}
	return nil, common.Error(common.ErrResourceConflict, common.Field("node", name), common.Field("type", "shadow"))
}

func (n *nodeService) updateDesire(shadow *models.Shadow, desire specV1.Desire) (*models.Shadow, error) {
	if shadow.Desire == nil {
		shadow.Desire = desire
	} else {
		err := shadow.Desire.Merge(desire)
		if err != nil {
			return nil, err
		}
	}

	return n.shadow.UpdateDesire(shadow)
}

func (n *nodeService) GetDesire(namespace, name string) (*specV1.Desire, error) {
	shadow, _ := n.shadow.Get(namespace, name)
	if shadow == nil {
		return nil, common.Error(
			common.ErrResourceNotFound,
			common.Field("type", "shadow"),
			common.Field("name", name),
			common.Field("namespace", namespace))
	}
	return &shadow.Desire, nil
}

func (n *nodeService) updateNodeAndAppIndex(namespace string, node *specV1.Node, shadow *models.Shadow) error {
	apps, err := n.app.ListApplication(namespace, &models.ListOptions{})
	if err != nil {
		log.L().Error("list application error", log.Error(err))
		return err
	}

	desire, appNames := n.rematchApplicationsForNode(apps, node.Labels)

	node.Desire = desire

	if _, err = n.updateDesire(shadow, desire); err != nil {
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
func (n *nodeService) UpdateNodeAppVersion(namespace string, app *specV1.Application) ([]string, error) {
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

func (n *nodeService) GetNodeProperties(namespace, name string) (*models.NodeProperties, error) {
	node, err := n.node.GetNode(namespace, name)
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
	meta, err := getNodePropertiesMeta(node)
	if err != nil {
		return nil, err
	}
	report := map[string]interface{}{}
	props, ok := shadow.Report[common.NodeProps]
	if ok && props != nil {
		report, ok = props.(map[string]interface{})
		if !ok {
			return nil, errors.New("invalid report shadow")
		}
	}
	desire := map[string]interface{}{}
	props, ok = shadow.Desire[common.NodeProps]
	if ok && props != nil {
		desire, ok = props.(map[string]interface{})
		if !ok {
			return nil, errors.New("invalid report shadow")
		}
	}
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
func (n *nodeService) UpdateNodeProperties(namespace, name string, props *models.NodeProperties) (*models.NodeProperties, error) {
	node, err := n.node.GetNode(namespace, name)
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
	meta, err := getNodePropertiesMeta(node)
	if err != nil {
		return nil, common.Error(common.ErrResourceNotFound, common.Field("node", node.Name), common.Field("meta", "shadow meta"))
	}
	oldDesire := map[string]interface{}{}
	if props, ok := shadow.Desire[common.NodeProps]; ok && props != nil {
		oldDesire, ok = props.(map[string]interface{})
		if !ok {
			return nil, errors.New("invalid desire of node properties")
		}
	}
	var newDesire specV1.Desire = props.State.Desire
	diff, err := newDesire.DiffWithNil(oldDesire)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	for key, val := range diff {
		meta.DesireMeta[key] = now
		if val == nil {
			delete(meta.DesireMeta, key)
		}
	}
	report := map[string]interface{}{}
	if props, ok := shadow.Report[common.NodeProps]; ok && props != nil {
		report, ok = props.(map[string]interface{})
		if !ok {
			return nil, errors.New("invalid report of node properties")
		}
	}
	props.State.Report = report
	// cast to map[string]interface{} should not omit
	props.State.Desire = map[string]interface{}(newDesire)
	props.Meta.ReportMeta = meta.ReportMeta
	props.Meta.DesireMeta = meta.DesireMeta
	// cast to map[string]interface{} should not omit
	shadow.Desire[common.NodeProps] = map[string]interface{}(newDesire)
	_, err = n.shadow.UpdateDesire(shadow)
	if err != nil {
		return nil, err
	}
	updateNodePropertiesMeta(node, meta)
	if _, err := n.node.UpdateNode(namespace, node); err != nil {
		return nil, err
	}
	return props, nil
}

func (n *nodeService) UpdateNodeMode(ns, name, mode string) error {
	node, err := n.node.GetNode(ns, name)
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

func getNodePropertiesMeta(node *specV1.Node) (*models.NodePropertiesMetadata, error) {
	propsMeta := &models.NodePropertiesMetadata{
		ReportMeta: map[string]interface{}{},
		DesireMeta: map[string]interface{}{},
	}
	meta, ok := node.Attributes[common.ReportMeta]
	if ok && meta != nil {
		propsMeta.ReportMeta, ok = meta.(map[string]interface{})
		if !ok {
			return nil, errors.New("invalid report meta")
		}
	}
	meta, ok = node.Attributes[common.DesireMeta]
	if ok && meta != nil {
		propsMeta.DesireMeta, ok = meta.(map[string]interface{})
		if !ok {
			return nil, errors.New("invalid desire meta")
		}
	}
	return propsMeta, nil
}

func updateNodePropertiesMeta(node *specV1.Node, meta *models.NodePropertiesMetadata) {
	if meta == nil {
		return
	}
	if node.Attributes == nil {
		node.Attributes = map[string]interface{}{}
	}
	if meta.ReportMeta != nil {
		node.Attributes[common.ReportMeta] = meta.ReportMeta
	}
	if meta.DesireMeta != nil {
		node.Attributes[common.DesireMeta] = meta.DesireMeta
	}
}
