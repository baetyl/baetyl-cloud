package kube

import (
	"github.com/baetyl/baetyl-go/v2/log"
	"testing"

	"encoding/json"
	"github.com/baetyl/baetyl-cloud/common"
	"github.com/baetyl/baetyl-cloud/models"
	"github.com/baetyl/baetyl-cloud/plugin/kube/apis/cloud/v1alpha1"
	"github.com/baetyl/baetyl-cloud/plugin/kube/client/clientset/versioned/fake"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/stretchr/testify/assert"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func genNodeRuntime() []runtime.Object {
	rs := []runtime.Object{
		&v1alpha1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-get",
				Namespace: "default",
				Labels: map[string]string{
					"node": "get",
				},
				Annotations: map[string]string{
					common.AnnotationMetadata: `{"brand":"ZTC"}`,
				},
			},
			Spec: v1alpha1.NodeSpec{
				DesireRef: &v1.LocalObjectReference{
					Name: "test-get",
				},
				ReportRef: &v1.LocalObjectReference{
					Name: "test-get",
				},
			},
		},
		&v1alpha1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-get01",
				Namespace: "default",
				Annotations: map[string]string{
					common.AnnotationMetadata: `{"brand":"ZTC"}`,
				},
			},
			Spec: v1alpha1.NodeSpec{
				DesireRef: &v1.LocalObjectReference{
					Name: "test-get01",
				},
				ReportRef: &v1.LocalObjectReference{
					Name: "test-get01",
				},
			},
		},
		&v1alpha1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-get02",
				Namespace: "default",
				Annotations: map[string]string{
					common.AnnotationMetadata: `{"brand":"ZTC"}`,
				},
			},
			Spec: v1alpha1.NodeSpec{
				DesireRef: &v1.LocalObjectReference{
					Name: "test-get02",
				},
				ReportRef: &v1.LocalObjectReference{
					Name: "test-get02",
				},
			},
		},
		&v1alpha1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-update",
				Namespace: "default",
			},
		},
		&v1alpha1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-create",
				Namespace: "default",
			},
		},
		&v1alpha1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-delete",
				Namespace: "default",
			},
			Spec: v1alpha1.NodeSpec{
				DesireRef: &v1.LocalObjectReference{
					Name: "test-delete",
				},
				ReportRef: &v1.LocalObjectReference{
					Name: "test-delete",
				},
			},
		},
		&v1alpha1.NodeDesire{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-get",
				Namespace: "default",
			},
		},
		&v1alpha1.NodeReport{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-get",
				Namespace: "default",
			},
		},
		&v1alpha1.NodeDesire{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-delete",
				Namespace: "default",
			},
		},
		&v1alpha1.NodeReport{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-delete",
				Namespace: "default",
			},
		},
		&v1alpha1.NodeDesire{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-get02",
				Namespace: "default",
			},
		},
		&v1alpha1.NodeDesire{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-create01",
				Namespace: "default",
			},
		},
		&v1alpha1.NodeReport{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-create01",
				Namespace: "default",
			},
		},
		&v1alpha1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-delete01",
				Namespace: "default",
			},
			Spec: v1alpha1.NodeSpec{
				DesireRef: &v1.LocalObjectReference{
					Name: "test-delete01",
				},
				ReportRef: &v1.LocalObjectReference{
					Name: "test-delete01",
				},
			},
		},
		&v1alpha1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-delete02",
				Namespace: "default",
			},
			Spec: v1alpha1.NodeSpec{
				DesireRef: &v1.LocalObjectReference{
					Name: "test-delete02",
				},
				ReportRef: &v1.LocalObjectReference{
					Name: "test-delete02",
				},
			},
		},
		&v1alpha1.NodeDesire{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-delete02",
				Namespace: "default",
			},
		},
		&v1alpha1.NodeReport{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-delete01",
				Namespace: "default",
			},
		},
	}
	return rs
}

func initNodeClient() *client {
	fc := fake.NewSimpleClientset(genNodeRuntime()...)
	return &client{
		customClient: fc,
		log:          log.With(log.Any("plugin", "kube")),
	}
}

func TestGetNode(t *testing.T) {
	c := initNodeClient()
	_, err := c.GetNode("default", "test")
	assert.NotNil(t, err)

	cfg, err := c.GetNode("default", "test-get")
	assert.Equal(t, cfg.Name, "test-get")

	_, err = c.GetNode("default", "test-get01")
	assert.Nil(t, err)

	_, err = c.GetNode("default", "test-get02")
	assert.Nil(t, err)
}

func TestCreateNode(t *testing.T) {
	c := initNodeClient()
	node := &specV1.Node{
		Name:      "test-create",
		Namespace: "default",
	}
	_, err := c.CreateNode(node.Namespace, node)
	assert.NotNil(t, err)
	node.Name += "new"
	node2, err := c.CreateNode(node.Namespace, node)
	assert.NoError(t, err)
	assert.Equal(t, node.Name, node2.Name)

	node = &specV1.Node{
		Name:      "test-create01",
		Namespace: "default",
		Annotations: map[string]string{
			"brand": "ZTC",
		},
	}
	node3, err := c.CreateNode(node.Namespace, node)
	assert.Nil(t, err)
	assert.Equal(t, node.Name, node3.Name)
}

func TestUpdateNode(t *testing.T) {
	c := initNodeClient()
	cfg := &specV1.Node{
		Name:      "test-update",
		Namespace: "default",
		Labels: map[string]string{
			"tag": "test",
		},
	}
	cfg2, err := c.UpdateNode(cfg.Namespace, cfg)
	assert.NoError(t, err)
	assert.Equal(t, cfg.Name, cfg2.Name)
	v, _ := cfg2.Labels["tag"]
	assert.Equal(t, v, "test")

	cfg.Name = cfg.Name + "NULL"
	_, err = c.UpdateNode(cfg.Namespace, cfg)
	assert.NotNil(t, err)
}

func TestDeleteNode(t *testing.T) {
	c := initNodeClient()
	err := c.DeleteNode("default", "test-delete")
	assert.NoError(t, err)

	err = c.DeleteNode("default", "test-delete03")
	assert.Error(t, err)

	err = c.DeleteNode("default", "test-delete01")
	assert.Nil(t, err)

	err = c.DeleteNode("default", "test-delete02")
	assert.Nil(t, err)
}

func TestListNode(t *testing.T) {
	c := initNodeClient()
	list, err := c.ListNode("default", &models.ListOptions{})
	assert.NoError(t, err)
	assert.True(t, true, len(list.Items) > 0)

	_, err = c.ListNode("baetyl-cloud", &models.ListOptions{LabelSelector: "node=test"})
	assert.NoError(t, err)
}

func TestUpdateNodeDesire(t *testing.T) {
	c := initNodeClient()
	namespace := "default"
	name := "test_name"
	shadow := &models.Shadow{
		Namespace: namespace,
		Name:      name,
		Desire:    genDesire(),
	}
	_, err := c.UpdateDesire(shadow)
	assert.Error(t, err)

	shadow = &models.Shadow{
		Namespace: namespace,
		Name:      "test-update",
		Desire:    genDesire(),
	}
	_, err = c.UpdateDesire(shadow)
	assert.Error(t, err)

	shadow = &models.Shadow{
		Namespace: namespace,
		Name:      "test-create01",
		Desire:    genDesire(),
	}
	n, err := c.UpdateDesire(shadow)
	assert.NoError(t, err)
	assert.Equal(t, shadow.Name, n.Name)

}

func TestUpdateNodeReport(t *testing.T) {
	c := initNodeClient()
	namespace := "default"
	name := "test_name"
	shadow := &models.Shadow{
		Namespace: namespace,
		Name:      name,
		Report:    genReport(),
	}
	_, err := c.UpdateReport(shadow)
	assert.Error(t, err)

	shadow = &models.Shadow{
		Namespace: namespace,
		Name:      "test-update",
		Desire:    genDesire(),
	}
	_, err = c.UpdateReport(shadow)
	assert.Error(t, err)

	shadow = &models.Shadow{
		Namespace: namespace,
		Name:      "test-create01",
		Desire:    genDesire(),
	}
	s, err := c.UpdateReport(shadow)
	assert.NoError(t, err)
	assert.Equal(t, shadow.Name, s.Name)
}

func genReport() specV1.Report {
	content := `
{
  "time":"2019-12-23T09:13:32.190611958Z",
  "software":{
      "os":"linux",
      "arch":"amd64",
      "pwd":"/usr/local",
      "mode":"docker",
      "go_version":"go1.12.5",
      "bin_version":"0.1.5",
      "conf_version":"V21"
  }
}`

	report := specV1.Report{}
	json.Unmarshal([]byte(content), report)
	return report
}

func genDesire() specV1.Desire {
	content := `
{
  "apps":{
      "app01":"001"
	}
}`

	desire := specV1.Desire{}
	json.Unmarshal([]byte(content), &desire)
	return desire
}
