package service

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/baetyl/baetyl-go/v2/json"
	"github.com/baetyl/baetyl-go/v2/log"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/baetyl/baetyl-go/v2/trigger"
	"github.com/baetyl/baetyl-go/v2/utils"

	"github.com/baetyl/baetyl-cloud/v2/cachemsg"
	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
	"github.com/baetyl/baetyl-cloud/v2/triggerfunc"
)

const (
	KeyCheckResourceDependency = "checkResourceDependency"
	KeyDeleteCoreExtResource   = "deleteCoreExtResource"
)

const ReportTimeKey = "time"

type CheckResourceDependencyFunc func(ns, nodeName string) error
type DeleteCoreExtResource func(ns string, node *specV1.Node) error

//go:generate mockgen -destination=../mock/service/node.go -package=service github.com/baetyl/baetyl-cloud/v2/service NodeService

const casRetryTimes = 3

// NodeService NodeService
type NodeService interface {
	Get(tx interface{}, namespace, name string) (*specV1.Node, error)
	List(namespace string, listOptions *models.ListOptions) (*models.NodeList, error)
	Count(namespace string) (map[string]int, error)
	CountAll() (map[string]int, error)

	Create(tx interface{}, namespace string, node *specV1.Node) (*specV1.Node, error)
	Update(namespace string, node *specV1.Node) (*specV1.Node, error)
	Delete(namespace string, node *specV1.Node) error

	UpdateReport(namespace, name string, report specV1.Report) (*models.Shadow, error)
	UpdateDesire(tx interface{}, namespace string, names []string, app *specV1.Application, f func(*models.Shadow, *specV1.Application)) error

	GetDesire(namespace, name string) (*specV1.Desire, error)

	UpdateNodeAppVersion(tx interface{}, namespace string, app *specV1.Application) ([]string, error)
	DeleteNodeAppVersion(tx interface{}, namespace string, app *specV1.Application) ([]string, error)

	GetNodeProperties(ns, name string) (*models.NodeProperties, error)
	UpdateNodeProperties(ns, name string, props *models.NodeProperties) (*models.NodeProperties, error)
	UpdateNodeMode(ns, name, mode string) error
}

type NodeServiceImpl struct {
	IndexService  IndexService
	App           plugin.Application
	Node          plugin.Node
	Shadow        plugin.Shadow
	Cache         plugin.DataCache
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

	cache, err := plugin.GetPlugin(config.Plugin.Cache)
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

	err = trigger.Register(triggerfunc.ShadowCreateOrUpdateTrigger, trigger.EventFunc{
		Args:  []interface{}{cache.(plugin.DataCache)},
		Event: triggerfunc.ShadowCreateOrUpdateCacheSet,
	})
	if err != nil {
		return nil, err
	}

	err = trigger.Register(triggerfunc.ShadowDelete, trigger.EventFunc{
		Args:  []interface{}{cache.(plugin.DataCache)},
		Event: triggerfunc.ShadowDeleteCache,
	})

	if err != nil {
		return nil, err
	}

	return &NodeServiceImpl{
		IndexService:  is,
		SysAppService: system,
		Node:          node.(plugin.Node),
		Shadow:        shadow.(plugin.Shadow),
		App:           app.(plugin.Application),
		Hooks:         make(map[string]interface{}),
		Cache:         cache.(plugin.DataCache),
	}, nil
}

// Get get the node
func (n *NodeServiceImpl) Get(tx interface{}, namespace, name string) (*specV1.Node, error) {
	node, err := n.Node.GetNode(tx, namespace, name)
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
	res, err := n.Node.CreateNode(tx, namespace, node)
	if err != nil {
		log.L().Error("create node failed", log.Error(err))
		return nil, err
	}

	_, err = n.SysAppService.GenApps(tx, namespace, node)
	if err != nil {
		return nil, err
	}

	if err = n.InsertOrUpdateNodeAndAppIndex(tx, namespace, res, models.NewShadowFromNode(res), true); err != nil {
		return nil, err
	}
	return res, err
}

// Update update node
func (n *NodeServiceImpl) Update(namespace string, node *specV1.Node) (*specV1.Node, error) {
	list, err := n.Node.UpdateNode(nil, namespace, []*specV1.Node{node})
	if err != nil || len(list) < 1 {
		return nil, err
	}
	res := list[0]

	shadow, err := n.Shadow.Get(nil, namespace, node.Name)
	if err != nil {
		return nil, err
	}

	// delete indexes for node and apps
	if err := n.IndexService.RefreshAppsIndexByNode(nil, namespace, res.Name, []string{}); err != nil {
		return nil, err
	}

	if err = n.InsertOrUpdateNodeAndAppIndex(nil, namespace, res, shadow, false); err != nil {
		return nil, err
	}
	return res, nil
}

// List get list node
func (n *NodeServiceImpl) List(namespace string, listOptions *models.ListOptions) (*models.NodeList, error) {
	// get list default create desc
	list, err := n.Node.ListNode(nil, namespace, listOptions)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if len(list.Items) == 0 {
		return list, nil
	}
	var resNode []specV1.Node

	// get shadow report time form cache if exists
	shadowReportTimeMap, err := n.GetAllShadowReportTime(namespace, len(list.Items))
	if err != nil {
		return nil, errors.Trace(err)
	}
	if listOptions.CreateSort != "" || listOptions.Ready != "" || listOptions.Cluster != "" {
		// filter sort
		resNode, err = n.filterListNode(list, namespace, listOptions, shadowReportTimeMap)
		list.Total = len(resNode)
	} else {
		// default sort  online ranked first then  crateTim desc
		resNode, err = n.defaultListNode(list, namespace, shadowReportTimeMap)
	}
	if err != nil {
		return nil, errors.Trace(err)
	}
	start, end := models.GetPagingParam(listOptions, list.Total)
	list.Items = resNode[start:end]

	var names []string
	for i := range list.Items {
		names = append(names, list.Items[i].Name)
	}
	// only get need to return data report
	shadowReportMap, err := n.GetShadowReportCacheByNames(namespace, names)
	if err != nil {
		return nil, errors.Trace(err)
	}
	// set node report
	for i := range list.Items {
		report := specV1.Report{}
		data := shadowReportMap[list.Items[i].Name]
		if data != nil {
			err = json.Unmarshal(data, &report)
			if err != nil {
				return nil, errors.Trace(err)
			}
		}
		list.Items[i].Report = report
	}

	return list, nil
}

func Reverse(s []interface{}) {
	for i := 0; i < len(s)/2; i++ {
		j := len(s) - i - 1
		s[i], s[j] = s[j], s[i]
	}
}

func (n *NodeServiceImpl) filterListNode(list *models.NodeList, namespace string, listOptions *models.ListOptions, shadowReportTimeMap map[string]string) ([]specV1.Node, error) {
	var resNode []specV1.Node
	if listOptions.Ready != "" || listOptions.Cluster != "" {
		for i := range list.Items {
			node := list.Items[i]
			if listOptions.Cluster != "" {
				if node.Cluster != (listOptions.Cluster == models.NodeTypeCluster) {
					continue
				}
			}
			if listOptions.Ready != "" {
				reportTime := shadowReportTimeMap[node.Name]
				after, err := n.GetNodeAfterTime(node, namespace)
				if err != nil {
					return nil, err
				}
				t, _ := time.Parse(time.RFC3339Nano, reportTime)
				switch listOptions.Ready {
				case models.ReadyTypeOnline:
					if time.Now().UTC().After(t.Add(after)) {
						continue
					}
				case models.ReadyTypOffline:
					noTime, err := time.Parse("2006-01-02 03:04:05", "0001-01-01 00:00:00")
					if err != nil {
						return nil, err
					}
					if t == noTime || time.Now().UTC().Before(t.Add(after)) {
						continue
					}
				case models.ReadyTypeUninstall:
					noTime, err := time.Parse("2006-01-02 03:04:05", "0001-01-01 00:00:00")
					if err != nil {
						return nil, err
					}
					if t != noTime {
						continue
					}
				default:

				}
			}
			resNode = append(resNode, node)
		}
	} else {
		resNode = list.Items
	}

	if listOptions.CreateSort != "" {
		if listOptions.CreateSort == models.NodeSortAsc {
			resNode = reverseLieNodeSlice(resNode)
		}
	}
	return resNode, nil
}

func reverseLieNodeSlice(s []specV1.Node) []specV1.Node {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func (n *NodeServiceImpl) defaultListNode(list *models.NodeList, namespace string, shadowReportMap map[string]string) ([]specV1.Node, error) {
	var onlineNode, offlineNode, resNode []specV1.Node
	for idx := range list.Items {
		node := list.Items[idx]
		reportTime, ok := shadowReportMap[node.Name]
		if !ok {
			reportTime = ""
		}
		after, err := n.GetNodeAfterTime(node, namespace)
		if err != nil {
			return nil, err
		}
		onLineFlag := false
		if reportTime != "" {
			t, _ := time.Parse(time.RFC3339Nano, reportTime)
			if time.Now().UTC().Before(t.Add(after)) {
				onLineFlag = true
			}
		}
		if onLineFlag == true {
			onlineNode = append(onlineNode, node)
		} else {
			offlineNode = append(offlineNode, node)
		}
	}
	resNode = append(resNode, onlineNode...)
	resNode = append(resNode, offlineNode...)
	return resNode, nil
}

func (n *NodeServiceImpl) GetNodeAfterTime(node specV1.Node, namespace string) (time.Duration, error) {
	if node.Attributes == nil {
		return 0, common.Error(common.ErrResourceNotFound, common.Field("type", "Attributes"), common.Field("namespace", node.Namespace))
	}
	freq, ok := node.Attributes[specV1.BaetylCoreFrequency].(string)
	if !ok {
		return 0, common.Error(common.ErrResourceNotFound, common.Field("type", specV1.BaetylCoreFrequency), common.Field("namespace", node.Namespace))
	}
	freqTime, err := strconv.Atoi(freq)
	if err != nil {
		return 0, err
	}
	after := time.Duration(freqTime+20) * time.Second

	return after, err
}

// GetAllShadowReportTime get node report time from cache
func (n *NodeServiceImpl) GetAllShadowReportTime(namespace string, lenNode int) (reportTimeMap map[string]string, err error) {

	reportTimeMap = map[string]string{}

	reportTimeOk, err := n.Cache.Exist(cachemsg.AllShadowReportTimeCache)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if !reportTimeOk {
		// cache not exist set cache
		return n.SetNodeShadowCache(namespace)
	} else {
		dataReportTime, err := n.Cache.GetByte(cachemsg.AllShadowReportTimeCache)
		if err != nil {
			return nil, errors.Trace(err)
		}
		if dataReportTime != nil {
			err = json.Unmarshal(dataReportTime, &reportTimeMap)
			if err != nil {
				return nil, errors.Trace(err)
			}
		}
		// check time cache len < lenNode
		if len(reportTimeMap) < lenNode {
			return n.SetNodeShadowCache(namespace)
		}
	}
	return reportTimeMap, nil
}

func (n *NodeServiceImpl) GetShadowReportCacheByNames(namespace string, names []string) (reportMap map[string][]byte, err error) {
	reportMap = map[string][]byte{}
	for i := range names {
		report, err := n.Cache.GetByte(cachemsg.GetShadowReportCacheKey(names[i]))
		// if report not exit then init report cache
		if err != nil || report == nil {
			return n.setShadowReportCacheByNames(namespace, names)
		}
		reportMap[names[i]] = report
	}
	return reportMap, nil
}

// setShadowReportCacheByNames init report cache by names
func (n *NodeServiceImpl) setShadowReportCacheByNames(namespace string, names []string) (reportMap map[string][]byte, err error) {
	shadowList, err := n.Shadow.ListShadowByNames(nil, namespace, names)
	if err != nil {
		return nil, errors.Trace(err)
	}
	reportMap = map[string][]byte{}
	for i := range shadowList {
		reportMap[shadowList[i].Name] = []byte(shadowList[i].ReportStr)
	}
	// sync set if CacheReportSetLock ==false
	// if CacheReportSetLock == ture return data form database
	exit, err := n.Cache.Exist(cachemsg.CacheReportSetLock)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if exit {
		lockTime, err := n.Cache.GetString(cachemsg.CacheReportSetLock)
		if err != nil {
			return nil, errors.Trace(err)
		}
		// delete lock key when key set time > 30 minute
		if time.Now().Add(-30*time.Minute).Format(time.RFC3339Nano) > lockTime {
			err = n.Cache.Delete(cachemsg.CacheReportSetLock)
			if err != nil {
				log.L().Error("delete lock key cacheUpdateReportTimeLock error", log.Error(err))
			}
		}
		log.L().Info("lock data return database back")
	} else {
		err := n.Cache.SetString(cachemsg.CacheReportSetLock, time.Now().Format(time.RFC3339Nano))
		if err != nil {
			return nil, errors.Trace(err)
		}
		// set node report cache
		go n.setShadowReportCache(reportMap)
	}
	return reportMap, nil
}

// setShadowReportCache set node report cache
func (n *NodeServiceImpl) setShadowReportCache(reportMap map[string][]byte) {
	log.L().Info("start set report cache")
	defer func() {
		if p := recover(); p != nil {
			log.L().Error(fmt.Sprintf("set report cache error %s", p))
		}
	}()

	for name, value := range reportMap {
		err := n.Cache.SetByte(cachemsg.GetShadowReportCacheKey(name), value)
		if err != nil {
			log.L().Error(fmt.Sprintf("set report cache %s failed", name))
		}
	}
	err := n.Cache.Delete(cachemsg.CacheReportSetLock)
	if err != nil {
		log.L().Error(fmt.Sprintf("delete CacheReportSetLock err %s", err))
	}
}

// SetNodeShadowCache set node shadow report cache time
func (n *NodeServiceImpl) SetNodeShadowCache(namespace string) (reportTimeMap map[string]string, err error) {
	log.L().Info("set report and reportTime cache")
	reportTimeMap = map[string]string{}
	shadowList, err := n.Shadow.ListAll(namespace)
	if err != nil {
		return nil, errors.Trace(err)
	}
	for i := range shadowList.Items {
		shadowMode := shadowList.Items[i]
		reportTimeMap[shadowMode.Name] = shadowMode.Time.Format(time.RFC3339Nano)
	}
	reportTimeData, err := json.Marshal(reportTimeMap)
	if err != nil {
		return nil, errors.Trace(err)
	}
	err = n.Cache.SetByte(cachemsg.AllShadowReportTimeCache, reportTimeData)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return
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

// Count get current node number
func (n *NodeServiceImpl) CountAll() (map[string]int, error) {
	total, err := n.Node.CountAllNode(nil)
	if err != nil {
		return nil, err
	}
	return map[string]int{
		plugin.QuotaNode: total,
	}, nil
}

// Delete delete node
func (n *NodeServiceImpl) Delete(namespace string, node *specV1.Node) error {
	if check, ok := n.Hooks[KeyCheckResourceDependency].(CheckResourceDependencyFunc); ok {
		if err := check(namespace, node.Name); err != nil {
			return err
		}
	}
	if delFunc, ok := n.Hooks[KeyDeleteCoreExtResource].(DeleteCoreExtResource); ok {
		if err := delFunc(namespace, node); err != nil {
			return err
		}
	}

	if err := n.Node.DeleteNode(nil, namespace, node.Name); err != nil {
		return err
	}

	if err := n.Shadow.Delete(namespace, node.Name); err != nil {
		common.LogDirtyData(err,
			log.Any("type", common.Shadow),
			log.Any("namespace", namespace),
			log.Any("name", node.Name),
			log.Any("operation", "delete"))
	}

	if err := n.IndexService.RefreshAppsIndexByNode(nil, namespace, node.Name, []string{}); err != nil {
		common.LogDirtyData(err,
			log.Any("type", "app node index"),
			log.Any("namespace", namespace),
			log.Any("name", node.Name),
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
		_, err = n.Node.GetNode(nil, namespace, name)
		if err != nil {
			return nil, err
		}
		return n.createShadow(nil, namespace, name, nil, report)
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
	node, err := n.Node.GetNode(nil, ns, name)
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
	if _, err := n.Node.UpdateNode(nil, ns, []*specV1.Node{node}); err != nil {
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
func (n *NodeServiceImpl) UpdateDesire(tx interface{}, namespace string, names []string, app *specV1.Application, f func(*models.Shadow, *specV1.Application)) error {
	shadows, err := n.Shadow.ListShadowByNames(tx, namespace, names)
	if err != nil {
		return err
	}
	for _, shadow := range shadows {
		// Refresh desire in Shadow by app
		f(shadow, app)
	}
	return n.Shadow.UpdateDesires(tx, shadows)
}

func (n *NodeServiceImpl) updateDesire(tx interface{}, shadow *models.Shadow, desire specV1.Desire) error {
	if shadow.Desire == nil {
		shadow.Desire = desire
	} else {
		err := shadow.Desire.Merge(desire)
		if err != nil {
			return err
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

func (n *NodeServiceImpl) InsertOrUpdateNodeAndAppIndex(tx interface{}, namespace string, node *specV1.Node, shadow *models.Shadow, flag bool) error {
	apps, err := n.App.ListApplication(tx, namespace, &models.ListOptions{})
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
		if err = n.updateDesire(tx, shadow, desire); err != nil {
			log.L().Error("update node desired node failed", log.Error(err))
			return err
		}
	}

	if err = n.IndexService.RefreshAppsIndexByNode(tx, namespace, node.Name, appNames); err != nil {
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
func (n *NodeServiceImpl) UpdateNodeAppVersion(tx interface{}, namespace string, app *specV1.Application) ([]string, error) {
	if app.Selector == "" {
		return nil, nil
	}

	// list nodes
	nodeList, err := n.Node.ListNode(tx, namespace, &models.ListOptions{LabelSelector: app.Selector})
	if err != nil {
		return nil, err
	}
	// update nodes
	var nodes []string
	for idx := range nodeList.Items {
		node := &nodeList.Items[idx]
		nodes = append(nodes, node.Name)
	}
	err = n.UpdateDesire(tx, namespace, nodes, app, RefreshNodeDesireByApp)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

func (n *NodeServiceImpl) createShadow(tx interface{}, namespace, name string, desire specV1.Desire, report specV1.Report) (*models.Shadow, error) {
	shadow := models.NewShadow(namespace, name)

	if desire != nil {
		shadow.Desire = desire
	}

	if report != nil {
		shadow.Report = report
	}
	return n.Shadow.Create(tx, shadow)
}

// DeleteNodeAppVersion delete the node desire's appVersion for app deleted
func (n *NodeServiceImpl) DeleteNodeAppVersion(tx interface{}, namespace string, app *specV1.Application) ([]string, error) {
	if app.Selector == "" {
		return nil, nil
	}

	// list nodes
	nodeList, err := n.Node.ListNode(tx, namespace, &models.ListOptions{LabelSelector: app.Selector})
	if err != nil {
		return nil, err
	}

	// update nodes
	var nodes []string
	for idx := range nodeList.Items {
		node := &nodeList.Items[idx]
		nodes = append(nodes, node.Name)
	}
	err = n.UpdateDesire(tx, namespace, nodes, app, DeleteNodeDesireByApp)
	if err != nil {
		return nil, err
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
	node, err := n.Node.GetNode(nil, namespace, name)
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
	node, err := n.Node.GetNode(nil, namespace, name)
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
	err = n.Shadow.UpdateDesire(nil, shadow)
	if err != nil {
		return nil, err
	}
	updateNodePropertiesMeta(node, meta)
	if _, err := n.Node.UpdateNode(nil, namespace, []*specV1.Node{node}); err != nil {
		return nil, err
	}
	return props, nil
}

func (n *NodeServiceImpl) UpdateNodeMode(ns, name, mode string) error {
	node, err := n.Node.GetNode(nil, ns, name)
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
	_, err = n.Node.UpdateNode(nil, ns, []*specV1.Node{node})
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
