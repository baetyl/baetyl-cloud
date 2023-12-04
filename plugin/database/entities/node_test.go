// Package entities 数据库存储基本结构与方法
package entities

import (
	"testing"
	"time"

	"github.com/baetyl/baetyl-cloud/v2/common"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/stretchr/testify/assert"
)

func TestFromNodeModel(t *testing.T) {
	node := &Node{
		Name:      "testApp",
		Namespace: "namespace",
		Labels:    `{"baetyl-cloud-system":"true"}`,
	}
	mNode := &specV1.Node{
		Namespace: "namespace",
		Name:      "testApp",
		Labels: map[string]string{
			common.LabelSystem: "true",
		},
	}
	modelNode, err := FromNodeModel("namespace", mNode)
	assert.NoError(t, err)
	assert.Equal(t, node.Name, modelNode.Name)
	assert.Equal(t, node.Namespace, modelNode.Namespace)
	assert.Equal(t, node.Labels, modelNode.Labels)
}
func TestToNodeModel(t *testing.T) {
	node := &Node{
		ID:          123,
		Name:        "testApp",
		Namespace:   "namespace",
		Labels:      "{\"baetyl-cloud-system\":\"true\"}",
		Attributes:  "{\"abc\":\"123\"}",
		Annotations: "{\"abc\":\"123\"}",
	}
	modelApp, err := ToNodeModel(node)
	assert.NoError(t, err)
	assert.Equal(t, node.Name, modelApp.Name)
	assert.Equal(t, node.Namespace, modelApp.Namespace)
	assert.Equal(t, map[string]string{common.LabelSystem: "true"}, modelApp.Labels)
	node = &Node{
		ID:        123,
		Name:      "testApp",
		Namespace: "namespace",
	}
	_, err = ToNodeModel(node)
	assert.NotNil(t, err)
}

func TestConvert(t *testing.T) {
	spNode := &specV1.Node{
		Namespace:   "default",
		Name:        "node01",
		Accelerator: "nv",
		Mode:        specV1.CloudMode,
		Cluster:     true,
		SysApps:     []string{"app123"},
		Labels:      map[string]string{"1": "2"},
		Annotations: map[string]string{"3": "4"},
		Attributes:  map[string]interface{}{"5": "6"},
		Description: "desc",
	}
	n, err := FromNodeModel("default", spNode)
	assert.NoError(t, err)
	res, err := ToNodeModel(n)
	assert.NoError(t, err)
	spNode.Version = res.Version
	spNode.CreationTimestamp = res.CreationTimestamp
	spNode.Attributes[specV1.KeyAccelerator] = spNode.Accelerator
	spNode.Attributes[specV1.KeyCluster] = spNode.Cluster
	spNode.Attributes[specV1.BaetylCoreFrequency] = common.DefaultCoreFrequency
	spNode.Attributes[specV1.BaetylCoreAPIPort] = common.DefaultCoreAPIPort
	spNode.Attributes[specV1.KeySyncMode] = specV1.CloudMode

	assert.EqualValues(t, spNode, res)
}

func TestEqualNode(t *testing.T) {
	// case 1
	node1 := &specV1.Node{
		Name:              "29987d6a2b8f11eabc62186590da6863",
		Namespace:         "default",
		Version:           "",
		Description:       "desc",
		CreationTimestamp: time.Unix(1000, 1000),
		Labels:            map[string]string{"a": "b"},
		Annotations:       map[string]string{"aa": "bb"},
		Attributes:        map[string]interface{}{"123": "1"},
	}
	node2 := &specV1.Node{
		Name:              "29987d6a2b8f11eabc62186590da6863",
		Namespace:         "default",
		Version:           "",
		Description:       "desc",
		CreationTimestamp: time.Unix(1000, 1000),
		Labels:            map[string]string{"a": "b"},
		Annotations:       map[string]string{"aa": "bb"},
		Attributes:        map[string]interface{}{"123": "1"},
	}
	flag := EqualNode(node1, node2)
	assert.True(t, flag)

	// case 2
	node1 = &specV1.Node{
		Name:              "29987d6a2b8f11eabc62186590da6863",
		Namespace:         "default",
		Version:           "",
		Description:       "desc",
		CreationTimestamp: time.Unix(1000, 1000),
		Labels:            map[string]string{"a": "b"},
		Annotations:       map[string]string{"aa": "bb"},
		Attributes:        map[string]interface{}{"123": "1"},
	}
	node2 = &specV1.Node{
		Name:              "29987d6a2b8f11eabc62186590da6863",
		Namespace:         "default",
		Version:           "",
		Description:       "desc",
		CreationTimestamp: time.Unix(1000, 1000),
		Labels:            map[string]string{"a": "b"},
		Annotations:       map[string]string{"aa": "bb"},
		Attributes:        map[string]interface{}{"123": "100"},
	}
	flag = EqualNode(node1, node2)
	assert.False(t, flag)

	// case 3
	node1 = &specV1.Node{
		Name:              "29987d6a2b8f11eabc62186590da6863",
		Namespace:         "default",
		Version:           "",
		Description:       "desc",
		CreationTimestamp: time.Unix(1000, 1000),
		Labels:            map[string]string{"a": "b"},
		Annotations:       map[string]string{"aa": "bb"},
		Attributes:        map[string]interface{}{"123": "1"},
	}
	node2 = &specV1.Node{
		Name:              "29987d6a2b8f11eabc62186590da6863",
		Namespace:         "default",
		Version:           "",
		Description:       "desc",
		CreationTimestamp: time.Unix(1000, 1000),
		Labels:            map[string]string{"a": "b"},
		Annotations:       map[string]string{"aa": "bbcccc"},
		Attributes:        map[string]interface{}{"123": "1"},
	}
	flag = EqualNode(node1, node2)
	assert.False(t, flag)

	// case 4
	node1 = &specV1.Node{
		Name:              "29987d6a2b8f11eabc62186590da6863",
		Namespace:         "default",
		Version:           "",
		Description:       "desc",
		CreationTimestamp: time.Unix(1000, 1000),
		Labels:            map[string]string{"a": "b"},
		Annotations:       map[string]string{"aa": "bb"},
		Attributes:        map[string]interface{}{"123": "1"},
	}
	node2 = &specV1.Node{
		Name:              "29987d6a2b8f11eabc62186590da6863",
		Namespace:         "default",
		Version:           "",
		Description:       "desc",
		CreationTimestamp: time.Unix(1000, 1000),
		Labels:            map[string]string{"a": "bbbbbb"},
		Annotations:       map[string]string{"aa": "bb"},
		Attributes:        map[string]interface{}{"123": "1"},
	}
	flag = EqualNode(node1, node2)
	assert.False(t, flag)
}
