package kube

import (
	"github.com/baetyl/baetyl-go/v2/log"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin/kube/apis/cloud/v1alpha1"
	"github.com/baetyl/baetyl-cloud/v2/plugin/kube/client/clientset/versioned/fake"
	"github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func genShadowRuntime() []runtime.Object {
	rs := []runtime.Object{
		&v1alpha1.NodeDesire{
			TypeMeta: metav1.TypeMeta{
				Kind:       "NodeDesire",
				APIVersion: v1alpha1.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "node01",
				Namespace: "default",
				Labels: map[string]string{
					common.LabelNodeName: "node01",
				},
			},
		},
		&v1alpha1.NodeReport{
			TypeMeta: metav1.TypeMeta{
				Kind:       "NodeReport",
				APIVersion: v1alpha1.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "node01",
				Namespace: "default",
				Labels: map[string]string{
					common.LabelNodeName: "node01",
				},
			},
		},
		&v1alpha1.NodeDesire{
			TypeMeta: metav1.TypeMeta{
				Kind:       "NodeDesire",
				APIVersion: v1alpha1.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "node02",
				Namespace: "default",
				Labels: map[string]string{
					common.LabelNodeName: "node02",
				},
			},
		},
		&v1alpha1.NodeReport{
			TypeMeta: metav1.TypeMeta{
				Kind:       "NodeReport",
				APIVersion: v1alpha1.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "node02",
				Namespace: "default",
				Labels: map[string]string{
					common.LabelNodeName: "node02",
				},
			},
		},
	}
	return rs
}

func initShadowClient() *client {
	fc := fake.NewSimpleClientset(genShadowRuntime()...)
	return &client{
		customClient: fc,
		log:          log.With(log.Any("plugin", "kube")),
	}
}

func TestClient_Get(t *testing.T) {
	c := initShadowClient()
	shadow, err := c.Get("default", "node01")
	assert.NoError(t, err)
	assert.Equal(t, "node01", shadow.Name)
	assert.NotNil(t, shadow.Report)
	assert.NotNil(t, shadow.Desire)
}

func TestClient_Create(t *testing.T) {
	c := initShadowClient()

	namespace := "default"
	name := "node-test"
	shadow := models.NewShadow(namespace, name)
	shd, err := c.Create(shadow)
	assert.NoError(t, err)
	assert.Equal(t, name, shd.Name)
	assert.NotNil(t, shd.Report)

	namespace = "default"
	name = "node01"
	shadow = models.NewShadow(namespace, name)
	shd, err = c.Create(shadow)
	assert.NoError(t, err)
	assert.Equal(t, name, shd.Name)
	assert.NotNil(t, shd.Report)
}

func TestClient_List(t *testing.T) {
	c := initShadowClient()

	namespace := "default"
	nodeList := &models.NodeList{
		Items: []specV1.Node{
			{
				Namespace: namespace,
				Name:      "node01",
			},
			{
				Namespace: namespace,
				Name:      "node02",
			},
		},
	}

	list, err := c.List(namespace, nodeList)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(list.Items))
	assert.Equal(t, "node01", list.Items[0].Name)
	assert.Equal(t, "node02", list.Items[1].Name)
}

func TestClient_Delete(t *testing.T) {
	c := initShadowClient()

	namespace := "default"
	name := "node01"
	err := c.Delete(namespace, name)
	assert.NoError(t, err)
	_, err = c.Get(namespace, name)
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "not found"))
}

func TestClient_UpdateDesire(t *testing.T) {
	c := initShadowClient()

	namespace := "default"
	name := "node01"
	shadow := models.NewShadow(namespace, name)
	shadow.Desire["apps"] = []v1.AppInfo{
		{
			Name:    "app01",
			Version: "1",
		},
	}
	shd, err := c.UpdateDesire(shadow)
	assert.NoError(t, err)
	assert.Equal(t, "app01", shd.Desire.AppInfos(false)[0].Name)
}

func TestClient_UpdateReport(t *testing.T) {
	c := initShadowClient()

	namespace := "default"
	name := "node01"
	shadow := models.NewShadow(namespace, name)
	shadow.Report["apps"] = []v1.AppInfo{
		{
			Name:    "app01",
			Version: "1",
		},
	}
	shd, err := c.UpdateReport(shadow)
	assert.NoError(t, err)
	assert.Equal(t, "app01", shd.Report.AppInfos(false)[0].Name)
}

func TestGeneratorLabelSelector(t *testing.T) {
	nodelist := &models.NodeList{
		Items: []v1.Node{
			{
				Name: "node01",
			},
			{
				Name: "node02",
			},
			{
				Name: "node03",
			},
		},
	}

	expectSelector := common.LabelNodeName + " in ( node01,node02,node03 )"
	assert.Equal(t, expectSelector, generatorLabelSelector(nodelist))
}
