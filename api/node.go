package api

import (
	"github.com/baetyl/baetyl-cloud/v2/plugin"
	"time"

	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/baetyl/baetyl-go/v2/log"
	"github.com/baetyl/baetyl-go/v2/spec/v1"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/service"
)

const (
	OfflineDuration = 40 * time.Second
	NodeNumber      = 1
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

	view.Desire = nil
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
	nodeList, err := api.Node.List(ns, api.parseListOptions(c))
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
	node, err := api.Node.Create(n.Namespace, n)
	if err != nil {
		if e := api.ReleaseQuota(ns, plugin.QuotaNode, NodeNumber); e != nil {
			log.L().Error("ReleaseQuota error", log.Error(e))
		}
		return nil, err
	}

	apps, err := api.Init.GenApps(n.Namespace, n.Name)
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

					if vv, ok := secret.Labels[v1.SecretLabel]; ok && vv == v1.SecretCertificate {
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
	params := map[string]interface{}{
		"InitApplyYaml": "baetyl-init-deployment.yml",
		"mode":          mode,
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
