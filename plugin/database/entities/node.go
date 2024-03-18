// Package entities 数据库存储基本结构与方法
package entities

import (
	"fmt"
	"reflect"
	"time"

	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/baetyl/baetyl-go/v2/json"
	"github.com/baetyl/baetyl-go/v2/log"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"

	"github.com/baetyl/baetyl-cloud/v2/common"
)

type Node struct {
	ID          int64     `db:"id"`
	Namespace   string    `db:"namespace"`
	Name        string    `db:"name"`
	Labels      string    `db:"labels"`
	Version     string    `db:"version"`
	CoreVersion string    `db:"core_version"`
	NodeMode    string    `db:"node_mode"`
	Description string    `db:"description"`
	Annotations string    `db:"annotations"`
	Attributes  string    `db:"attributes"`
	CreateTime  time.Time `db:"create_time"`
	UpdateTime  time.Time `db:"update_time"`
}

func ToNodeModel(node *Node) (*specV1.Node, error) {
	labels := map[string]string{}
	err := json.Unmarshal([]byte(node.Labels), &labels)
	if err != nil {
		return nil, errors.Trace(err)
	}
	attributes := map[string]interface{}{}
	if node.Attributes != "" {
		err = json.Unmarshal([]byte(node.Attributes), &attributes)
		if err != nil {
			return nil, errors.Trace(err)
		}
	}

	accelerator := ""
	if res, ok := attributes[specV1.KeyAccelerator]; ok {
		if val, ok := res.(string); ok {
			accelerator = val
		}
	}
	attributes[specV1.KeyAccelerator] = accelerator

	var mode specV1.SyncMode
	if res, ok := attributes[specV1.KeySyncMode]; ok {
		if val, ok := res.(string); ok {
			mode = specV1.SyncMode(val)
		}
	}
	attributes[specV1.KeySyncMode] = mode

	cluster := false
	if res, ok := attributes[specV1.KeyCluster]; ok {
		if val, ok := res.(bool); ok {
			cluster = val
		}
	}
	attributes[specV1.KeyCluster] = cluster

	var link string
	if res, ok := attributes[specV1.KeyLink].(string); ok {
		link = res
	}
	var coreID string
	if res, ok := attributes[specV1.KeyCoreId].(string); ok {
		coreID = res
	}

	var sysApps []string
	if val, ok := attributes[specV1.KeyOptionalSysApps]; ok && val != nil {
		ss, ok := val.([]interface{})
		if !ok {
			log.L().Error(common.ErrConvertConflict, log.Any("name", specV1.KeyOptionalSysApps), log.Any("error", "failed to interface{} to []interface{}`"))
		}

		for _, d := range ss {
			s, ok := d.(string)
			if !ok {
				log.L().Error(common.ErrConvertConflict, log.Any("name", specV1.KeyOptionalSysApps), log.Any("error", "failed to interface{} to string`"))
			}
			sysApps = append(sysApps, s)
		}
	}
	attributes[specV1.KeyOptionalSysApps] = sysApps

	annotations := map[string]string{}
	if node.Annotations != "" {
		err = json.Unmarshal([]byte(node.Annotations), &annotations)
		if err != nil {
			return nil, errors.Trace(err)
		}
	}
	return &specV1.Node{
		Namespace:         node.Namespace,
		Name:              node.Name,
		Version:           node.Version,
		CreationTimestamp: node.CreateTime.UTC(),
		Accelerator:       accelerator,
		Mode:              mode,
		NodeMode:          node.NodeMode,
		Cluster:           cluster,
		Link:              link,
		CoreId:            coreID,
		Description:       node.Description,
		Annotations:       annotations,
		Attributes:        attributes,
		Labels:            labels,
		SysApps:           sysApps,
	}, nil
}

func FromNodeModel(namespace string, node *specV1.Node) (*Node, error) {
	labels, err := json.Marshal(node.Labels)
	if err != nil {
		return nil, errors.Trace(err)
	}

	if node.Attributes == nil {
		node.Attributes = map[string]interface{}{}
	}
	node.Attributes[specV1.KeyAccelerator] = node.Accelerator
	node.Attributes[specV1.KeyCluster] = node.Cluster
	node.Attributes[specV1.KeyOptionalSysApps] = node.SysApps
	node.Attributes[specV1.KeyLink] = node.Link
	node.Attributes[specV1.KeyCoreId] = node.CoreId
	// set default value
	if _, ok := node.Attributes[specV1.BaetylCoreFrequency]; !ok {
		node.Attributes[specV1.BaetylCoreFrequency] = common.DefaultCoreFrequency
	}
	if _, ok := node.Attributes[specV1.BaetylCoreAPIPort]; !ok {
		node.Attributes[specV1.BaetylCoreAPIPort] = common.DefaultCoreAPIPort
	}
	if _, ok := node.Attributes[specV1.BaetylAgentPort]; !ok {
		node.Attributes[specV1.BaetylAgentPort] = common.DefaultAgentPort
	}
	if _, ok := node.Attributes[specV1.KeySyncMode]; !ok {
		node.Attributes[specV1.KeySyncMode] = specV1.CloudMode
	}
	coreVersion := ""
	if v, ok := node.Attributes[specV1.BaetylCoreVersion]; ok {
		coreVersion = v.(string)
	}

	attributes, err := json.Marshal(node.Attributes)
	if err != nil {
		return nil, errors.Trace(err)
	}
	annotations, err := json.Marshal(node.Annotations)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return &Node{
		Name:        node.Name,
		Namespace:   namespace,
		Version:     GenResourceVersion(),
		CoreVersion: coreVersion,
		NodeMode:    node.NodeMode,
		Description: node.Description,
		Attributes:  string(attributes),
		Annotations: string(annotations),
		Labels:      string(labels),
	}, nil
}

func EqualNode(node1, node2 *specV1.Node) bool {
	return reflect.DeepEqual(node1.Labels, node2.Labels) &&
		reflect.DeepEqual(node1.Description, node2.Description) &&
		reflect.DeepEqual(node1.Annotations, node2.Annotations) &&
		reflect.DeepEqual(node1.Attributes, node2.Attributes) &&
		reflect.DeepEqual(node1.SysApps, node2.SysApps) &&
		node1.Accelerator == node2.Accelerator &&
		node1.NodeMode == node2.NodeMode
}

func GenResourceVersion() string {
	return fmt.Sprintf("%d%s", time.Now().UTC().Unix(), common.RandString(6))
}
