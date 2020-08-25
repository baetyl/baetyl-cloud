package api

import (
	"time"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/baetyl/baetyl-go/v2/log"
	"github.com/baetyl/baetyl-go/v2/spec/v1"
)

const offlineDuration = 40 * time.Second

var (
	CmdExpirationInSeconds = int64(60 * 60)
)

// GetNode get a node
func (api *API) GetNode(c *common.Context) (interface{}, error) {
	ns, n := c.GetNamespace(), c.GetNameFromParam()
	node, err := api.nodeService.Get(ns, n)
	if err != nil {
		return nil, err
	}

	view, err := node.View(offlineDuration)
	if err != nil {
		return nil, err
	}

	view.Desire = nil
	return view, nil
}

func (api *API) GetNodes(c *common.Context) (interface{}, error) {
	ns := c.GetNamespace()
	nodeNames := &models.NodeNames{}
	err := c.LoadBody(nodeNames)
	if err != nil {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", err.Error()))
	}
	var nodesView = []*v1.NodeView{}
	for _, name := range nodeNames.Names{
		node, err := api.nodeService.Get(ns, name)
		if err != nil {
			continue
		}
		view, err := node.View(offlineDuration)
		if err != nil {
			continue
		}
		view.Desire = nil
		nodesView = append(nodesView, view)
	}
	return nodesView, nil
}

// GetNodeStats get a node stats
func (api *API) GetNodeStats(c *common.Context) (interface{}, error) {
	ns, n := c.GetNamespace(), c.GetNameFromParam()

	node, err := api.nodeService.Get(ns, n)
	if err != nil {
		return nil, err
	}

	view, err := node.View(offlineDuration)
	if err != nil {
		return nil, err
	}

	view.Desire = nil
	return view, nil

}

// ListNode list node
func (api *API) ListNode(c *common.Context) (interface{}, error) {
	ns := c.GetNamespace()
	nodeList, err := api.nodeService.List(ns, api.parseListOptions(c))
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
		view, err = n.View(offlineDuration)
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
	n, err := api.parseAndCheckNode(c)
	if err != nil {
		return nil, err
	}
	ns, name := c.GetNamespace(), n.Name

	n.Labels = common.AddSystemLabel(n.Labels, map[string]string{
		common.LabelNodeName: name,
	})

	oldNode, err := api.nodeService.Get(ns, name)
	if err != nil {
		if e, ok := err.(errors.Coder); !ok || e.Code() != common.ErrResourceNotFound {
			return nil, err
		}
	}

	if oldNode != nil {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "this name is already in use"))
	}

	node, err := api.nodeService.Create(ns, n)
	if err != nil {
		return nil, err
	}

	// generate system applications
	list := []common.SystemApplication{
		common.BaetylCore,
		common.BaetylFunction,
	}
	_, err = api.GenSysApp(name, ns, list)
	if err != nil {
		return nil, err
	}

	view, err := node.View(offlineDuration)
	if err != nil {
		return nil, err
	}

	view.Desire = nil
	view.Report = nil
	return view, nil
}

// UpdateNode update the node
func (api *API) UpdateNode(c *common.Context) (interface{}, error) {
	node, err := api.parseAndCheckNode(c)
	if err != nil {
		return nil, err
	}
	ns, n := c.GetNamespace(), c.GetNameFromParam()
	oldNode, err := api.nodeService.Get(ns, n)
	if err != nil {
		return nil, err
	}

	node.Labels = common.AddSystemLabel(node.Labels, map[string]string{
		common.LabelNodeName: node.Name,
	})
	node.Version = oldNode.Version
	node, err = api.nodeService.Update(c.GetNamespace(), node)

	if err != nil {
		return nil, err
	}

	view, err := node.View(offlineDuration)
	if err != nil {
		return nil, err
	}

	view.Desire = nil
	return view, nil
}

// DeleteNode delete the node
func (api *API) DeleteNode(c *common.Context) (interface{}, error) {
	ns, n := c.GetNamespace(), c.GetNameFromParam()
	node, err := api.nodeService.Get(ns, n)
	if err != nil {
		return nil, err
	}

	// Delete Node
	if err := api.nodeService.Delete(c.GetNamespace(), c.GetNameFromParam()); err != nil {
		return nil, err
	}

	sysAppInfos := node.Desire.AppInfos(true)
	for _, ai := range sysAppInfos {
		// Clean APP
		app, err := api.applicationService.Get(ns, ai.Name, "")
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
					if err := api.configService.Delete(ns, v.Config.Name); err != nil {
						common.LogDirtyData(err,
							log.Any("type", common.Config),
							log.Any("namespace", ns),
							log.Any("name", v.Config.Name))
					}
				}
				// Clean Secret
				if v.Secret != nil {
					secret, err := api.secretService.Get(ns, v.Secret.Name, "")
					if err != nil {
						common.LogDirtyData(err,
							log.Any("type", common.Secret),
							log.Any(common.KeyContextNamespace, ns),
							log.Any("name", v.Secret.Name))
						continue
					}

					if vv, ok := secret.Labels[v1.SecretLabel]; ok && vv == v1.SecretCertificate {
						if certID, _ok := secret.Annotations[common.AnnotationPkiCertID]; _ok {
							if err := api.pkiService.DeleteClientCertificate(certID); err != nil {
								common.LogDirtyData(err,
									log.Any("type", "pki"),
									log.Any(common.KeyContextNamespace, ns),
									log.Any(common.AnnotationPkiCertID, certID))
							}
						} else {
							log.L().Warn("failed to get "+common.AnnotationPkiCertID+" of certificate secret", log.Any(common.KeyContextNamespace, ns), log.Any("name", v.Secret.Name))
						}
					}
					if err := api.secretService.Delete(ns, v.Secret.Name); err != nil {
						common.LogDirtyData(err,
							log.Any("type", common.Secret),
							log.Any(common.KeyContextNamespace, ns),
							log.Any("name", v.Secret.Name))
					}
				}
			}
		}
		if err := api.applicationService.Delete(ns, ai.Name, ai.Version); err != nil {
			common.LogDirtyData(err,
				log.Any("type", common.Application),
				log.Any(common.KeyContextNamespace, ns),
				log.Any("name", ai.Name))
		}
		if err := api.indexService.RefreshNodesIndexByApp(ns, ai.Name, make([]string, 0)); err != nil {
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

	node, err := api.nodeService.Get(ns, n)
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
	_, err := api.nodeService.Get(ns, name)
	if err != nil {
		return nil, err
	}
	cmd, err := api.genCmd(string(common.Node), ns, name)
	if err != nil {
		return nil, err
	}
	return map[string]string{"cmd": cmd}, nil
}

// GetNodeDeployHistory list node // TODO will support later
func (api *API) GetNodeDeployHistory(c *common.Context) (interface{}, error) {
	return nil, nil
}

func (api *API) parseAndCheckNode(c *common.Context) (*v1.Node, error) {
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

func (api *API) NodeNumberCollector(namespace string) (map[string]int, error) {
	list, err := api.nodeService.List(namespace, &models.ListOptions{})
	if err != nil {
		return nil, err
	}
	return map[string]int{
		plugin.QuotaNode: len(list.Items),
	}, nil
}
