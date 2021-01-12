package api

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/baetyl/baetyl-cloud/v2/plugin"

	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/baetyl/baetyl-go/v2/log"
	v1 "github.com/baetyl/baetyl-go/v2/spec/v1"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/service"
)

const (
	OfflineDuration         = 40 * time.Second
	NodeNumber              = 1
	BaetylCorePrevVersion   = "BaetylCorePrevVersion"
	BaetylCoreFrequency     = "BaetylCoreFrequency"
	BaetylNodeNameKey       = "baetyl-node-name"
	BaetylAppNameKey        = "baetyl-app-name"
	BaetylCoreConfPrefix    = "baetyl-core-conf"
	BaetylCoreContainerPort = 80
	BaetylVersionPrefix     = "baetyl-version-"
	DefaultMode             = "kube"
)

// GetNode get a node
func (api *API) GetNode(c *common.Context) (interface{}, error) {
	ns, n := c.GetNamespace(), c.GetNameFromParam()
	node, err := api.Node.Get(ns, n)
	if err != nil {
		return nil, err
	}

	view, err := node.View(OfflineDuration)
	if err != nil {
		return nil, err
	}

	return view, nil
}

func (api *API) GetNodes(c *common.Context) (interface{}, error) {
	nodeNames, err := api.ParseAndCheckNodeNames(c)
	if err != nil {
		return nil, err
	}
	ns := c.GetNamespace()
	nodeViewList := models.NodeViewList{
		Items: make([]v1.NodeView, 0),
	}
	for _, name := range nodeNames.Names {
		node, err := api.Node.Get(ns, name)
		if err != nil {
			if e, ok := err.(errors.Coder); ok && e.Code() == common.ErrResourceNotFound {
				continue
			}
			return nil, err
		}
		view, err := node.View(OfflineDuration)
		if err != nil {
			return nil, err
		}
		view.Desire = nil
		nodeViewList.Items = append(nodeViewList.Items, *view)
	}

	nodeViewList.Total = len(nodeViewList.Items)
	return nodeViewList, nil
}

// GetNodeStats get a node stats
func (api *API) GetNodeStats(c *common.Context) (interface{}, error) {
	ns, n := c.GetNamespace(), c.GetNameFromParam()

	node, err := api.Node.Get(ns, n)
	if err != nil {
		return nil, err
	}

	view, err := node.View(OfflineDuration)
	if err != nil {
		return nil, err
	}

	view.Desire = nil
	return view, nil
}

// ListNode list node
func (api *API) ListNode(c *common.Context) (interface{}, error) {
	ns := c.GetNamespace()
	params, err := api.parseListOptions(c)
	if err != nil {
		return nil, err
	}
	nodeList, err := api.Node.List(ns, params)
	if err != nil {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", err.Error()))
	}
	nodeViewList := models.NodeViewList{
		Total:       nodeList.Total,
		ListOptions: nodeList.ListOptions,
		Items:       make([]v1.NodeView, 0, len(nodeList.Items)),
	}

	for idx := range nodeList.Items {
		n := &nodeList.Items[idx]

		var view *v1.NodeView
		view, err = n.View(OfflineDuration)
		if err != nil {
			return nil, err
		}

		view.Desire = nil
		nodeViewList.Items = append(nodeViewList.Items, *view)
	}
	return nodeViewList, nil
}

// CreateNode create one node
func (api *API) CreateNode(c *common.Context) (interface{}, error) {
	n, err := api.ParseAndCheckNode(c)
	if err != nil {
		return nil, err
	}
	ns := c.GetNamespace()
	n.Namespace = ns

	n.Labels = common.AddSystemLabel(n.Labels, map[string]string{
		common.LabelNodeName: n.Name,
	})

	oldNode, err := api.Node.Get(n.Namespace, n.Name)
	if err != nil {
		if e, ok := err.(errors.Coder); !ok || e.Code() != common.ErrResourceNotFound {
			return nil, err
		}
	}

	if oldNode != nil {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "this name is already in use"))
	}

	err = api.License.AcquireQuota(ns, plugin.QuotaNode, NodeNumber)
	if err != nil {
		return nil, err
	}

	// set default frequency here
	if n.Attributes == nil {
		n.Attributes = map[string]interface{}{}
	}
	n.Attributes[BaetylCoreFrequency] = common.DefaultCoreFrequency
	n.Attributes[v1.KeyAccelerator] = n.Accelerator

	node, err := api.Node.Create(n.Namespace, n)
	if err != nil {
		if e := api.ReleaseQuota(ns, plugin.QuotaNode, NodeNumber); e != nil {
			log.L().Error("ReleaseQuota error", log.Error(e))
		}
		return nil, err
	}

	apps, err := api.Init.GenApps(n.Namespace, n)
	if err != nil {
		return nil, err
	}

	for _, app := range apps {
		err = api.UpdateNodeAndAppIndex(n.Namespace, app)
		if err != nil {
			return nil, errors.Trace(err)
		}
	}

	view, err := node.View(OfflineDuration)
	if err != nil {
		return nil, err
	}

	view.Desire = nil
	view.Report = nil
	return view, nil
}

// UpdateNode update the node
func (api *API) UpdateNode(c *common.Context) (interface{}, error) {
	node, err := api.ParseAndCheckNode(c)
	if err != nil {
		return nil, err
	}
	ns, n := c.GetNamespace(), c.GetNameFromParam()
	oldNode, err := api.Node.Get(ns, n)
	if err != nil {
		return nil, err
	}

	node.Labels = common.AddSystemLabel(node.Labels, map[string]string{
		common.LabelNodeName: node.Name,
	})
	node.Version = oldNode.Version

	if models.EqualNode(node, oldNode) {
		return oldNode.View(OfflineDuration)
	}

	node, err = api.Node.Update(c.GetNamespace(), node)

	if err != nil {
		return nil, err
	}

	view, err := node.View(OfflineDuration)
	if err != nil {
		return nil, err
	}

	view.Desire = nil
	return view, nil
}

// DeleteNode delete the node
func (api *API) DeleteNode(c *common.Context) (interface{}, error) {
	ns, n := c.GetNamespace(), c.GetNameFromParam()
	node, err := api.Node.Get(ns, n)
	if err != nil {
		if e, ok := err.(errors.Coder); ok && e.Code() == common.ErrResourceNotFound {
			return nil, nil
		}
		return nil, err
	}

	// Delete Node
	if err := api.Node.Delete(c.GetNamespace(), c.GetNameFromParam()); err != nil {
		return nil, err
	}
	if e := api.ReleaseQuota(ns, plugin.QuotaNode, NodeNumber); e != nil {
		log.L().Error("ReleaseQuota error", log.Error(e))
	}
	sysAppInfos := node.Desire.AppInfos(true)
	for _, ai := range sysAppInfos {
		// Clean APP
		app, err := api.App.Get(ns, ai.Name, "")
		if err != nil {
			if e, ok := err.(errors.Coder); ok && e.Code() == common.ErrResourceNotFound {
				continue
			}
			common.LogDirtyData(err,
				log.Any("type", common.Application),
				log.Any(common.KeyContextNamespace, ns),
				log.Any("name", ai.Name))
		} else {
			ai.Version = app.Version
			for _, v := range app.Volumes {
				// Clean Config
				if v.Config != nil {
					config, err := api.Config.Get(ns, v.Config.Name, "")
					if err != nil {
						common.LogDirtyData(err,
							log.Any("type", common.Config),
							log.Any(common.KeyContextNamespace, ns),
							log.Any("name", v.Config.Name))
						continue
					}

					if res := checkIsSysResources(config.Labels); !res {
						continue
					}

					if err := api.Config.Delete(ns, v.Config.Name); err != nil {
						common.LogDirtyData(err,
							log.Any("type", common.Config),
							log.Any("namespace", ns),
							log.Any("name", v.Config.Name))
					}
				}
				// Clean Secret
				if v.Secret != nil {
					secret, err := api.Secret.Get(ns, v.Secret.Name, "")
					if err != nil {
						common.LogDirtyData(err,
							log.Any("type", common.Secret),
							log.Any(common.KeyContextNamespace, ns),
							log.Any("name", v.Secret.Name))
						continue
					}

					if res := checkIsSysResources(secret.Labels); !res {
						continue
					}

					if vv, ok := secret.Labels[v1.SecretLabel]; ok && vv == v1.SecretConfig {
						if certID, _ok := secret.Annotations[common.AnnotationPkiCertID]; _ok {
							if err := api.PKI.DeleteClientCertificate(certID); err != nil {
								common.LogDirtyData(err,
									log.Any("type", "pki"),
									log.Any(common.KeyContextNamespace, ns),
									log.Any(common.AnnotationPkiCertID, certID))
							}
						} else {
							log.L().Warn("failed to get "+common.AnnotationPkiCertID+" of certificate secret", log.Any(common.KeyContextNamespace, ns), log.Any("name", v.Secret.Name))
						}
					}
					if err := api.Secret.Delete(ns, v.Secret.Name); err != nil {
						common.LogDirtyData(err,
							log.Any("type", common.Secret),
							log.Any(common.KeyContextNamespace, ns),
							log.Any("name", v.Secret.Name))
					}
				}
			}
		}
		if err := api.App.Delete(ns, ai.Name, ai.Version); err != nil {
			common.LogDirtyData(err,
				log.Any("type", common.Application),
				log.Any(common.KeyContextNamespace, ns),
				log.Any("name", ai.Name))
		}
		if err := api.Index.RefreshNodesIndexByApp(ns, ai.Name, make([]string, 0)); err != nil {
			common.LogDirtyData(err,
				log.Any("type", common.Index),
				log.Any(common.KeyContextNamespace, ns),
				log.Any("app", ai.Name))
		}
	}
	return nil, nil
}

// GetAppByNode list app
func (api *API) GetAppByNode(c *common.Context) (interface{}, error) {
	ns, n := c.GetNamespace(), c.GetNameFromParam()

	node, err := api.Node.Get(ns, n)
	if err != nil {
		return nil, err
	}

	appNames := make([]string, 0)
	if node.Desire != nil {
		// sysapp
		apps := node.Desire.AppInfos(true)
		for _, a := range apps {
			appNames = append(appNames, a.Name)
		}

		apps = node.Desire.AppInfos(false)
		for _, a := range apps {
			appNames = append(appNames, a.Name)
		}
	}

	return api.listAppByNames(ns, appNames)
}

// GenInitCmdFromNode generate install command
func (api *API) GenInitCmdFromNode(c *common.Context) (interface{}, error) {
	ns, name := c.GetNamespace(), c.Param("name")
	_, err := api.Node.Get(ns, name)
	if err != nil {
		return nil, err
	}
	mode := c.Query("mode")
	if mode == "" {
		mode = DefaultMode
	}
	params := map[string]interface{}{
		"mode": mode,
	}
	if mode == "kube" {
		params["InitApplyYaml"] = "baetyl-init-deployment.yml"
	} else if mode == "native" {
		params["InitApplyYaml"] = "baetyl-init-apply.json"
	} else {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("mode", mode))
	}

	cmd, err := api.Init.GetResource(ns, name, service.TemplateBaetylInitCommand, params)
	if err != nil {
		return nil, err
	}
	return map[string]string{"cmd": string(cmd.([]byte))}, nil
}

// GetNodeDeployHistory list node // TODO will support later
func (api *API) GetNodeDeployHistory(c *common.Context) (interface{}, error) {
	return nil, nil
}

func (api *API) ParseAndCheckNode(c *common.Context) (*v1.Node, error) {
	node := new(v1.Node)
	node.Name = c.GetNameFromParam()
	node.Namespace = c.GetNamespace()
	err := c.LoadBody(node)
	if err != nil {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", err.Error()))
	}
	if name := c.GetNameFromParam(); name != "" {
		node.Name = name
	}
	if node.Name == "" {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "name is required"))
	}

	return node, nil
}

func (api *API) ParseAndCheckNodeNames(c *common.Context) (*models.NodeNames, error) {
	_, ok := c.GetQuery("batch")
	if !ok {
		return nil, common.Error(common.ErrRequestParamInvalid)
	}
	nodeNames := &models.NodeNames{}
	err := c.LoadBody(nodeNames)
	if err != nil {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", err.Error()))
	}
	return nodeNames, nil
}

func (api *API) NodeNumberCollector(namespace string) (map[string]int, error) {
	return api.Node.Count(namespace)
}

func (api *API) GetNodeProperties(c *common.Context) (interface{}, error) {
	ns, n := c.GetNamespace(), c.GetNameFromParam()
	return api.Node.GetNodeProperties(ns, n)
}

func (api *API) getNodeSysAppSecretLikedResources(c *common.Context) (*models.SecretList, error) {
	ns, n := c.GetNamespace(), c.GetNameFromParam()
	_, err := api.App.Get(ns, n, "")
	if err != nil {
		return nil, err
	}

	ops := &models.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s", BaetylAppNameKey, n),
	}

	return api.Secret.List(ns, ops)
}

func (api *API) UpdateNodeProperties(c *common.Context) (interface{}, error) {
	ns, n := c.GetNamespace(), c.GetNameFromParam()
	props, err := api.ParseAndCheckProperties(c)
	if err != nil {
		return nil, err
	}
	return api.Node.UpdateNodeProperties(ns, n, props)
}

func (api *API) UpdateNodeMode(c *common.Context) (interface{}, error) {
	ns, n := c.GetNamespace(), c.GetNameFromParam()
	nodeMode, err := api.ParseAndCheckNodeMode(c)
	if err != nil {
		return nil, err
	}
	err = api.Node.UpdateNodeMode(ns, n, nodeMode.Mode)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (api *API) ParseAndCheckNodeMode(c *common.Context) (*models.NodeMode, error) {
	nodeMode := &models.NodeMode{}
	err := c.LoadBody(nodeMode)
	if err != nil {
		return nil, err
	}
	if nodeMode.Mode != string(v1.CloudMode) && nodeMode.Mode != string(v1.LocalMode) {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("mode", "mode should be local or cloud"))
	}
	return nodeMode, nil
}

func (api *API) ParseAndCheckProperties(c *common.Context) (*models.NodeProperties, error) {
	props := new(models.NodeProperties)
	err := c.LoadBody(props)
	if err != nil {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", err.Error()))
	}
	for _, v := range props.State.Desire {
		if _, ok := v.(string); !ok {
			return nil, common.Error(common.ErrRequestParamInvalid, common.Field("value", "desire value should be string"))
		}
	}
	return props, nil
}

func (api *API) UpdateCoreApp(c *common.Context) (interface{}, error) {
	ns, n := c.GetNamespace(), c.GetNameFromParam()

	coreConfig, err := api.parseCoreAppConfigs(c)
	if err != nil {
		return nil, err
	}
	// get node
	node, err := api.Node.Get(ns, n)
	if err != nil {
		return nil, err
	}

	// get core app
	app, err := api.getCoreAppByNodeName(ns, n)
	if err != nil {
		return nil, err
	}

	coreService, err := api.getCoreAppService(app)
	if err != nil {
		return nil, err
	}

	port, err := api.getCoreAppAPIPort(ns, coreService)
	if err != nil {
		return nil, err
	}

	version, err := api.getCoreCurrentVersionByImage(coreService.Image)
	if err != nil {
		return nil, err
	}

	// get frequency
	freq, err := api.getCoreAppFrequency(node)
	if err != nil {
		return nil, err
	}

	if coreConfig.Version == version &&
		coreConfig.Frequency == freq &&
		coreConfig.APIPort == port {
		return api.ToApplicationView(app)
	}

	api.updateCoreVersions(node, version, coreConfig.Version)

	image, err := api.getCoreImageByVersion(coreConfig.Version)
	if err != nil {
		return nil, err
	}
	coreService.Image = image

	err = api.updateCoreAppConfig(app, node, coreConfig.Frequency)
	if err != nil {
		return nil, err
	}
	node.Attributes[BaetylCoreFrequency] = fmt.Sprintf("%d", coreConfig.Frequency)

	err = api.updateCoreAppAPIPort(ns, coreService, coreConfig.APIPort)
	if err != nil {
		return nil, err
	}

	res, err := api.App.Update(ns, app)
	if err != nil {
		return nil, err
	}

	_, err = api.Node.UpdateNodeAppVersion(ns, res)
	if err != nil {
		return nil, err
	}

	_, err = api.Node.Update(ns, node)
	if err != nil {
		return nil, err
	}

	return api.ToApplicationView(res)
}

func (api *API) GetCoreAppConfigs(c *common.Context) (interface{}, error) {
	ns, n := c.GetNamespace(), c.GetNameFromParam()
	node, err := api.Node.Get(ns, n)
	if err != nil {
		return nil, err
	}

	var coreInfo models.NodeCoreConfigs
	app, err := api.getCoreAppByNodeName(ns, n)
	if err != nil {
		return nil, err
	}

	coreService, err := api.getCoreAppService(app)
	if err != nil {
		return nil, err
	}

	coreVersion, err := api.getCoreCurrentVersionByImage(coreService.Image)
	if err != nil {
		return nil, err
	}
	coreInfo.Version = coreVersion

	// get frequency
	coreInfo.Frequency, err = api.getCoreAppFrequency(node)
	if err != nil {
		return nil, err
	}

	// get api port
	coreInfo.APIPort, err = api.getCoreAppAPIPort(ns, coreService)
	if err != nil {
		return nil, err
	}
	return coreInfo, nil
}

func (api *API) GetCoreAppVersions(c *common.Context) (interface{}, error) {
	ns, n := c.GetNamespace(), c.GetNameFromParam()
	node, err := api.Node.Get(ns, n)
	if err != nil {
		return nil, err
	}

	if node.Attributes == nil {
		return nil, common.Error(common.ErrResourceNotFound, common.Field("type", "Attributes"), common.Field("namespace", ns))
	}

	var coreVersions models.NodeCoreVersions
	if v, ok := node.Attributes[BaetylCorePrevVersion]; ok {
		res, ok := v.(string)
		if !ok {
			return nil, common.Error(common.ErrConvertConflict, common.Field("name", "BaetylCorePrevVersion"), common.Field("error", "failed to convert to string`"))
		}
		coreVersions.Versions = append(coreVersions.Versions, res)
	}

	app, err := api.getCoreAppByNodeName(ns, n)
	if err != nil {
		return nil, err
	}

	coreService, err := api.getCoreAppService(app)
	if err != nil {
		return nil, err
	}

	currentVersion, err := api.getCoreCurrentVersionByImage(coreService.Image)
	if err != nil {
		return nil, err
	}
	coreVersions.Versions = append(coreVersions.Versions, currentVersion)

	latestVersion, err := api.Prop.GetPropertyValue(BaetylVersionPrefix + "latest")
	if err != nil {
		return nil, err
	}
	if latestVersion != "" && latestVersion != currentVersion {
		coreVersions.Versions = append(coreVersions.Versions, latestVersion)
	}
	return coreVersions, nil
}

// don't delete resource which doesn't belong to system
func checkIsSysResources(labels map[string]string) bool {
	v, ok := labels[common.LabelSystem]
	if !ok {
		return false
	}
	if res, _ := strconv.ParseBool(v); !res {
		return false
	}
	return true
}

func (api *API) parseCoreAppConfigs(c *common.Context) (*models.NodeCoreConfigs, error) {
	config := new(models.NodeCoreConfigs)
	err := c.LoadBody(config)
	if err != nil {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", err.Error()))
	}
	if config.Frequency < 1 || config.Frequency > 300 {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "freq must be between 1 - 300"))
	}

	if config.APIPort < 1024 || config.APIPort > 65535 {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "api port must be between 1024 - 65535"))
	}
	return config, nil
}

func (api *API) filterCoreVersionByImage(image string) (string, error) {
	params := &models.Filter{
		Name: BaetylVersionPrefix,
	}

	res, err := api.Prop.ListProperty(params)
	if err != nil {
		return "", err
	}
	for _, v := range res {
		if v.Value == image {
			return strings.TrimPrefix(v.Name, BaetylVersionPrefix), nil
		}
	}
	return "", common.Error(common.ErrResourceNotFound, common.Field("type", "BaetylCoreVersion"), common.Field("name", image))
}

func (api *API) getCoreCurrentVersionByImage(image string) (string, error) {
	version, err := api.filterCoreVersionByImage(image)
	if err != nil {
		return "", err
	}

	return version, nil
}

func (api *API) getCoreImageByVersion(version string) (string, error) {
	prop, err := api.Prop.GetProperty(BaetylVersionPrefix + version)
	if err != nil {
		return "", err
	}
	return prop.Value, nil
}

func (api *API) getCoreAppService(app *v1.Application) (*v1.Service, error) {
	for i, svr := range app.Services {
		if svr.Name != v1.BaetylCore {
			continue
		}
		return &app.Services[i], nil
	}
	return nil, common.Error(common.ErrResourceNotFound, common.Field("type", "service"), common.Field("name", v1.BaetylCore), common.Field("namespace", app.Namespace))
}

func (api *API) getCoreAppConfig(app *v1.Application) (*v1.Configuration, error) {
	for _, volume := range app.Volumes {
		if volume.Config == nil || !strings.Contains(volume.Config.Name, BaetylCoreConfPrefix) {
			continue
		}
		conf, err := api.Config.Get(app.Namespace, volume.Config.Name, "")
		if err != nil {
			return nil, err
		}

		return conf, nil
	}
	return nil, common.Error(common.ErrResourceNotFound, common.Field("type", "config"), common.Field("name", v1.BaetylCore), common.Field("namespace", app.Namespace))
}

func (api *API) updateCoreVersions(node *v1.Node, currentVersion, updateVersion string) {
	if v, ok := node.Attributes[BaetylCorePrevVersion]; ok && v.(string) == updateVersion {
		delete(node.Attributes, BaetylCorePrevVersion)
		return
	}
	if currentVersion == updateVersion {
		return
	}
	node.Attributes[BaetylCorePrevVersion] = currentVersion
}

func (api *API) updateCoreAppConfig(app *v1.Application, node *v1.Node, freq int) error {
	config, err := api.getCoreAppConfig(app)
	if err != nil {
		return err
	}
	var accelerator string
	if node.Attributes != nil {
		accelerator, _ = node.Attributes[v1.KeyAccelerator].(string)
	}
	params := map[string]interface{}{
		"CoreConfName":  config.Name,
		"CoreAppName":   app.Name,
		"CoreFrequency": fmt.Sprintf("%ds", freq),
		"GPUStats":      accelerator == v1.NVAccelerator,
	}
	res, err := api.Init.GetResource(config.Namespace, node.Name, service.TemplateCoreConfYaml, params)
	if err != nil {
		return err
	}

	var data []byte
	var ok bool
	if data, ok = res.([]byte); !ok {
		return common.Error(common.ErrConvertConflict, common.Field("name", "BaetylCoreConfig"), common.Field("error", "failed to convert to []byte`"))
	}

	var newConf v1.Configuration
	err = yaml.Unmarshal(data, &newConf)
	if err != nil {
		return common.Error(common.ErrTemplate, common.Field("error", err))
	}

	newConf.Name = config.Name
	newConf.Version = config.Version
	_, err = api.Config.Update(config.Namespace, &newConf)
	if err != nil {
		return err
	}
	return nil
}

func (api *API) getCoreAppAPIPort(ns string, service *v1.Service) (int, error) {
	for _, v := range service.Ports {
		if v.ContainerPort == int32(BaetylCoreContainerPort) {
			return int(v.HostPort), nil
		}
	}
	return 0, common.Error(common.ErrResourceNotFound, common.Field("type", "APIPort"), common.Field("name", v1.BaetylCore), common.Field("namespace", ns))
}

func (api *API) updateCoreAppAPIPort(ns string, service *v1.Service, port int) error {
	for i, v := range service.Ports {
		if v.ContainerPort == int32(BaetylCoreContainerPort) {
			service.Ports[i].HostPort = int32(port)
			return nil
		}
	}
	return common.Error(common.ErrResourceNotFound, common.Field("type", "APIPort"), common.Field("name", v1.BaetylCore), common.Field("namespace", ns))
}

func (api *API) getCoreAppByNodeName(ns, node string) (*v1.Application, error) {
	appList, err := api.Index.ListAppsByNode(ns, node)
	if err != nil {
		return nil, err
	}
	var core string
	for _, item := range appList {
		if strings.Contains(item, v1.BaetylCore) {
			core = item
			break
		}
	}
	if core == "" {
		return nil, common.Error(common.ErrResourceNotFound, common.Field("type", "app"), common.Field("name", v1.BaetylCore), common.Field("namespace", ns))
	}
	return api.App.Get(ns, core, "")
}

func (api *API) getCoreAppFrequency(node *v1.Node) (int, error) {
	if node.Attributes == nil {
		return 0, common.Error(common.ErrResourceNotFound, common.Field("type", "Attributes"), common.Field("namespace", node.Namespace))
	}
	if _, ok := node.Attributes[BaetylCoreFrequency]; !ok {
		return 0, common.Error(common.ErrResourceNotFound, common.Field("type", BaetylCoreFrequency), common.Field("namespace", node.Namespace))
	}
	freq, ok := node.Attributes[BaetylCoreFrequency].(string)
	if !ok {
		return 0, common.Error(common.ErrConvertConflict, common.Field("name", "BaetylCoreFrequency"), common.Field("error", "failed to convert to string`"))
	}
	res, err := strconv.Atoi(freq)
	if err != nil {
		return 0, errors.Trace(err)
	}
	return res, nil
}
