// Package database 数据库存储实现
package database

import (
	"fmt"
	"testing"
	"time"

	"github.com/baetyl/baetyl-go/v2/log"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

var (
	nodetabales = []string{
		`
CREATE TABLE baetyl_node
(
	id                integer       PRIMARY KEY AUTOINCREMENT,
    name              varchar(128)  NOT NULL DEFAULT '',
    namespace         varchar(64)   NOT NULL DEFAULT '',
	version           varchar(36)   NOT NULL DEFAULT '',
    node_mode         varchar(36)   NOT NULL DEFAULT '',
	core_version      varchar(36)   NOT NULL DEFAULT '',
	labels            varchar(2048) NOT NULL DEFAULT '{}',
	annotations       varchar(2048) NOT NULL DEFAULT '',
	attributes        varchar(2048) NOT NULL DEFAULT '{}',
	description       varchar(1024) NOT NULL DEFAULT '',
    create_time       timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time       timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`,
	}
)

func (d *BaetylCloudDB) MockCreateNodeTable() {
	for _, sql := range nodetabales {
		_, err := d.Exec(nil, sql)
		if err != nil {
			panic(fmt.Sprintf("create node exception: %s", err.Error()))
		}
	}
}

func TestNode(t *testing.T) {
	node := &specV1.Node{
		Name:              "29987d6a2b8f11eabc62186590da6863",
		Namespace:         "default",
		Version:           "",
		Description:       "desc",
		CreationTimestamp: time.Unix(1000, 1000),
		Labels:            map[string]string{"a": "b"},
		Annotations:       map[string]string{"aa": "bb"},
		Attributes:        map[string]interface{}{"123": "1"},
	}
	listOptions := &models.ListOptions{}
	log.L().Info("Test node", log.Any("node", node))

	db, err := MockNewDB()
	if err != nil {
		fmt.Printf("get mock sqlite3 error = %s", err.Error())
		t.Fail()
		return
	}
	db.MockCreateNodeTable()

	res, err := db.CreateNode(nil, "default", node)
	assert.NoError(t, err)
	checkNode(t, node, res)

	node2 := &specV1.Node{Name: "tx", Namespace: "tx"}
	tx, err := db.BeginTx()
	assert.NoError(t, err)
	res, err = db.CreateNode(tx, "tx", node2)
	assert.NoError(t, err)
	assert.Equal(t, res.Namespace, "tx")
	assert.NoError(t, tx.Commit())

	res, err = db.GetNode(nil, node.Namespace, node.Name)
	assert.NoError(t, err)
	checkNode(t, node, res)

	node.Labels = map[string]string{"b": "b"}
	list, err := db.UpdateNode(nil, "default", []*specV1.Node{node})
	assert.NoError(t, err)
	checkNode(t, node, list[0])

	res, err = db.GetNode(nil, node.Namespace, node.Name)
	assert.NoError(t, err)
	checkNode(t, node, res)

	resList, err := db.ListNode(nil, node.Namespace, listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 1)
	checkNode(t, node, &resList.Items[0])

	err = db.DeleteNode(nil, node.Namespace, node.Name)
	assert.NoError(t, err)

	res, err = db.GetNode(nil, node.Namespace, node.Name)
	assert.Nil(t, res)
}

func TestListNode(t *testing.T) {
	node1 := &specV1.Node{
		Name:              "node_123",
		Namespace:         "default",
		Description:       "desc",
		CreationTimestamp: time.Unix(1000, 1000),
		Labels:            map[string]string{"label": "aaa"},
		Annotations:       map[string]string{"annotation": "aaa"},
		Attributes:        map[string]interface{}{"attr": "aaa"},
	}
	node2 := &specV1.Node{
		Name:              "node_abc",
		Namespace:         "default",
		Description:       "desc",
		CreationTimestamp: time.Unix(1000, 1000),
		Labels:            map[string]string{"label": "aaa"},
		Annotations:       map[string]string{"annotation": "aaa"},
		Attributes:        map[string]interface{}{"attr": "aaa"},
	}
	node3 := &specV1.Node{
		Name:              "node_test",
		Namespace:         "default",
		Description:       "desc",
		CreationTimestamp: time.Unix(1000, 1000),
		Labels:            map[string]string{"label": "bbb"},
		Annotations:       map[string]string{"annotation": "bbb"},
		Attributes:        map[string]interface{}{"attr": "bbb"},
	}
	node4 := &specV1.Node{
		Name:              "node_testabc",
		Namespace:         "default",
		Description:       "desc",
		CreationTimestamp: time.Unix(1000, 1000),
		Labels:            map[string]string{"label": "bbb"},
		Annotations:       map[string]string{"annotation": "bbb"},
		Attributes:        map[string]interface{}{"attr": "bbb"},
	}

	listOptions := &models.ListOptions{}

	db, err := MockNewDB()
	if err != nil {
		fmt.Printf("get mock sqlite3 error = %s", err.Error())
		t.Fail()
		return
	}
	defer db.Close()
	db.MockCreateNodeTable()

	res, err := db.CreateNode(nil, "default", node1)
	assert.NoError(t, err)
	checkNode(t, node1, res)

	res, err = db.CreateNode(nil, "default", node2)
	assert.NoError(t, err)
	checkNode(t, node2, res)

	res, err = db.CreateNode(nil, "default", node3)
	assert.NoError(t, err)
	checkNode(t, node3, res)

	res, err = db.CreateNode(nil, "default", node4)
	assert.NoError(t, err)
	checkNode(t, node4, res)

	// list option nil, return all nodes
	resList, err := db.ListNode(nil, "default", listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 4)
	checkNode(t, node1, &resList.Items[0])
	checkNode(t, node2, &resList.Items[1])
	checkNode(t, node3, &resList.Items[2])
	checkNode(t, node4, &resList.Items[3])
	// page 1 num 2
	listOptions.PageNo = 1
	listOptions.PageSize = 2
	resList, err = db.ListNode(nil, "default", listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 4)
	checkNode(t, node1, &resList.Items[0])
	checkNode(t, node2, &resList.Items[1])
	// page 2 num 2
	listOptions.PageNo = 2
	listOptions.PageSize = 2
	resList, err = db.ListNode(nil, "default", listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 4)
	checkNode(t, node1, &resList.Items[0])
	checkNode(t, node2, &resList.Items[1])
	// page 3 num 0
	listOptions.PageNo = 3
	listOptions.PageSize = 2
	resList, err = db.ListNode(nil, "default", listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 4)
	// page 1 num 2 name like node
	listOptions.PageNo = 1
	listOptions.PageSize = 2
	listOptions.Name = "node"
	resList, err = db.ListNode(nil, "default", listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 4)
	checkNode(t, node1, &resList.Items[0])
	checkNode(t, node2, &resList.Items[1])
	// page 1 num 2 name like abc
	listOptions.PageNo = 1
	listOptions.PageSize = 2
	listOptions.Name = "abc"
	resList, err = db.ListNode(nil, "default", listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 2)
	checkNode(t, node2, &resList.Items[0])
	checkNode(t, node4, &resList.Items[1])
	// page 1 num2 label : aaa
	listOptions.PageNo = 1
	listOptions.PageSize = 4
	listOptions.Name = ""
	listOptions.LabelSelector = "label=aaa"
	resList, err = db.ListNode(nil, "default", listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 2)
	checkNode(t, node1, &resList.Items[0])
	checkNode(t, node2, &resList.Items[1])

	listOptions.PageNo = 1
	listOptions.PageSize = 4
	listOptions.Name = "abc"
	listOptions.LabelSelector = "label=aaa"
	resList, err = db.ListNode(nil, "default", listOptions)
	assert.NoError(t, err)
	assert.Equal(t, resList.Total, 1)
	checkNode(t, node2, &resList.Items[0])

	total, err := db.CountAllNode(nil)
	assert.NoError(t, err)
	assert.Equal(t, 4, total)

	err = db.DeleteNode(nil, "default", node1.Name)
	assert.NoError(t, err)
	err = db.DeleteNode(nil, "default", node2.Name)
	assert.NoError(t, err)
	err = db.DeleteNode(nil, "default", node3.Name)
	assert.NoError(t, err)
	err = db.DeleteNode(nil, "default", node4.Name)
	assert.NoError(t, err)

	res, err = db.GetNode(nil, "default", node1.Name)
	assert.Nil(t, res)
	res, err = db.GetNode(nil, "default", node2.Name)
	assert.Nil(t, res)
	res, err = db.GetNode(nil, "default", node3.Name)
	assert.Nil(t, res)
	res, err = db.GetNode(nil, "default", node4.Name)
	assert.Nil(t, res)

	total, err = db.CountAllNode(nil)
	assert.NoError(t, err)
	assert.Equal(t, 0, total)
}

func Test_GetNode(t *testing.T) {
	node := &specV1.Node{
		Namespace: "c14db79154f9461a90d633a1e3ece1b0",
		Name:      "gpu_instance",
		Version:   "16157784868jib2c",
		Labels:    map[string]string{"baetyl-node-name": "gpu-instance", "group": "dcell-iot", "helmet-poc1": "1", "system_group": "system"},
		Attributes: map[string]interface{}{
			"BaetylCoreAPIPort":   "30050",
			"BaetylCoreFrequency": "20",
			"accelerator":         "",
			"cluster":             false,
			"desireMeta":          map[string]interface{}{},
			"reportMeta": map[string]interface{}{
				"arc-poc1-2.arc-detection5.apiUrl":                 "2021-03-12T07:41:29.50850158Z",
				"arc-poc1-2.arc-detection5.tps":                    "2021-03-12T07:44:29.448501195Z",
				"arc-poc1-2.media-server6.fps":                     "2021-03-12T07:13:29.457143348Z",
				"arc-poc1-2.rtspreader3.fps":                       "2021-03-12T07:12:49.52962395Z",
				"arc-poc1-2.rtspreader3.isResize":                  "2021-03-12T07:05:49.448152417Z",
				"arc-poc1-2.rtspreader3.isVideo":                   "2021-03-12T07:20:49.475630299Z",
				"electric-fence-xfj-1.electronic-fence41.fence":    "2021-03-08T13:25:28.538174665Z",
				"electric-fence-xfj-1.person-detection39.apiUrl":   "2021-03-04T07:17:48.537850055Z",
				"electric-fence-xfj-1.rtspreader38.streamUri":      "2021-03-08T12:20:28.481160928Z",
				"handrail-xfj-1.handrail-analysis11.detectFence":   "2021-03-09T03:18:08.484238992Z",
				"handrail-xfj-1.handrail-analysis11.handrailFence": "2021-03-09T03:18:48.489688726Z",
				"handrail-xfj-1.rtspreader8.location":              "2021-03-09T03:16:28.535516402Z",
				"handrail-xfj-1.rtspreader8.streamUri":             "2021-03-09T03:16:28.535516402Z",
				"helmet-poc1-1.alarm-report4.alarmInterval":        "2021-03-11T03:22:28.52631121Z",
				"helmet-poc1-1.media-server5.fps":                  "2021-03-12T03:34:49.533772036Z",
				"helmet-poc1-1.person-detection3.tps":              "2021-03-12T03:35:09.501820206Z",
				"helmet-poc1-1.rtspreader0.fps":                    "2021-03-12T03:34:49.533772036Z",
				"helmet-poc1-1.rtspreader0.isVideo":                "2021-03-12T03:26:09.518983926Z",
				"helmet-poc1-1.rtspreader0.streamUri":              "2021-03-12T03:23:49.445369466Z",
				"smoking-poc1-1.media-server4.fps":                 "2021-03-11T06:21:28.556360057Z",
				"smoking-poc1-1.rtspreader0.fps":                   "2021-03-11T06:21:08.490948416Z",
				"suit-xfj-1.person-detection30.apiUrl":             "2021-03-04T07:13:08.528497056Z",
				"suit-xfj-1.person-detection30.httpUser":           "2021-03-05T01:28:28.533009333Z",
				"suit-xfj-1.rtspreader29.streamUri":                "2021-03-04T11:26:08.505080494Z",
				"suit-xfj-1.track31.DetectArea":                    "2021-03-04T11:26:48.574173717Z",
				"suit-xfj-1.track31.apiUrl":                        "2021-03-04T07:27:48.495648479Z",
				"wa-test-2.wash-detection1.modelName":              "2021-03-03T08:38:08.507942431Z",
				"wa-test-2.wash-detection1.modelVersion":           "2021-03-03T08:38:08.507942431Z",
				"wa-test-2.wash-detection1.threshold":              "2021-03-03T08:56:28.523869806Z",
				"wa-test-2.wash-detection1.tps":                    "2021-03-03T08:37:28.466464136Z",
				"wa-test-2.wash-detection1.width":                  "2021-03-03T08:38:48.495056947Z",
				"wash-xfj-2.person-detection3.apiUrl":              "2021-03-04T07:18:08.501784868Z",
				"wash-xfj-2.person-detection3.httpUser":            "2021-03-05T01:47:28.487963935Z",
				"wash-xfj-2.rtspreader2.isVideo":                   "2021-03-05T07:36:08.493789755Z",
				"wash-xfj-2.rtspreader2.streamUri":                 "2021-03-08T06:08:48.514735189Z",
				"wash-xfj-2.rtspreader4.streamUri":                 "2021-03-08T06:34:48.571668867Z",
				"wash-xfj-2.rtspreader5.streamUri":                 "2021-03-08T06:35:08.535238269Z",
				"wash-xfj-2.rtspreader6.streamUri":                 "2021-03-08T06:34:48.571668867Z",
				"wash-xfj-2.rtspreader7.streamUri":                 "2021-03-08T06:35:28.494872411Z",
				"wash-xfj-2.track15.DetectArea":                    "2021-03-08T06:07:08.504747407Z",
				"wash-xfj-2.track15.apiUrl":                        "2021-03-04T07:28:08.511233599Z",
				"wash-xfj-2.track15.iouThresh":                     "2021-03-08T05:58:48.48039961Z",
				"wash-xfj-2.track15.trackTimeout":                  "2021-03-05T08:16:28.564055795Z",
				"wash-xfj-2.wash-detection10.DetectArea":           "2021-03-04T07:35:48.525762824Z",
				"wash-xfj-2.wash-detection10.apiUrl":               "2021-03-04T07:21:28.58135481Z",
				"wash-xfj-2.wash-detection10.handsink":             "2021-03-04T07:37:48.500756897Z",
				"wash-xfj-2.wash-detection13.DetectArea":           "2021-03-04T07:35:48.525762824Z",
				"wash-xfj-2.wash-detection13.apiUrl":               "2021-03-04T07:21:28.58135481Z",
				"wash-xfj-2.wash-detection13.handsink":             "2021-03-04T07:38:08.563656127Z",
				"wash-xfj-2.wash-detection14.DetectArea":           "2021-03-04T07:35:28.558892444Z",
				"wash-xfj-2.wash-detection14.apiUrl":               "2021-03-04T07:21:48.568227163Z",
				"wash-xfj-2.wash-detection14.handsink":             "2021-03-04T07:38:48.48352313Z",
				"wash-xfj-2.wash-detection14.top":                  "2021-03-04T07:30:48.496813806Z",
				"wash-xfj-2.wash-detection9.DetectArea":            "2021-03-04T07:35:28.558892444Z",
				"wash-xfj-2.wash-detection9.apiUrl":                "2021-03-04T07:21:28.58135481Z",
				"wash-xfj-2.wash-detection9.handsink":              "2021-03-04T07:36:48.568415616Z",
			},
			"syncMode": specV1.CloudMode,
		},
		Description: "desc",
	}

	db, err := MockNewDB()
	if err != nil {
		fmt.Printf("get mock sqlite3 error = %s", err.Error())
		t.Fail()
		return
	}
	defer db.Close()
	db.MockCreateNodeTable()

	res, err := db.CreateNode(nil, node.Namespace, node)
	assert.NoError(t, err)
	checkNode(t, node, res)

	res, err = db.GetNode(nil, node.Namespace, node.Name)
	assert.NoError(t, err)
	checkNode(t, node, res)

}

func checkNode(t *testing.T, expect, actual *specV1.Node) {
	assert.Equal(t, expect.Name, actual.Name)
	assert.Equal(t, expect.Namespace, actual.Namespace)
	assert.Equal(t, expect.Description, actual.Description)
	assert.EqualValues(t, expect.Labels, actual.Labels)
	assert.EqualValues(t, expect.Attributes, actual.Attributes)
	assert.EqualValues(t, expect.Annotations, actual.Annotations)
}

func TestGetNodeByNames(t *testing.T) {
	node1 := &specV1.Node{
		Name:              "node_123",
		Namespace:         "default",
		Description:       "desc",
		CreationTimestamp: time.Unix(1000, 1000),
		Labels:            map[string]string{"label": "aaa"},
		Annotations:       map[string]string{"annotation": "aaa"},
		Attributes:        map[string]interface{}{"attr": "aaa"},
	}
	node2 := &specV1.Node{
		Name:              "node_abc",
		Namespace:         "default",
		Description:       "desc",
		CreationTimestamp: time.Unix(1000, 1000),
		Labels:            map[string]string{"label": "aaa"},
		Annotations:       map[string]string{"annotation": "aaa"},
		Attributes:        map[string]interface{}{"attr": "aaa"},
	}
	node3 := &specV1.Node{
		Name:              "node_test",
		Namespace:         "default",
		Description:       "desc",
		CreationTimestamp: time.Unix(1000, 1000),
		Labels:            map[string]string{"label": "bbb"},
		Annotations:       map[string]string{"annotation": "bbb"},
		Attributes:        map[string]interface{}{"attr": "bbb"},
	}
	node4 := &specV1.Node{
		Name:              "node_testabc",
		Namespace:         "default",
		Description:       "desc",
		CreationTimestamp: time.Unix(1000, 1000),
		Labels:            map[string]string{"label": "bbb"},
		Annotations:       map[string]string{"annotation": "bbb"},
		Attributes:        map[string]interface{}{"attr": "bbb"},
	}

	db, err := MockNewDB()
	if err != nil {
		fmt.Printf("get mock sqlite3 error = %s", err.Error())
		t.Fail()
		return
	}
	defer db.Close()
	db.MockCreateNodeTable()

	res, err := db.CreateNode(nil, "default", node1)
	assert.NoError(t, err)
	checkNode(t, node1, res)

	res, err = db.CreateNode(nil, "default", node2)
	assert.NoError(t, err)
	checkNode(t, node2, res)

	res, err = db.CreateNode(nil, "default", node3)
	assert.NoError(t, err)
	checkNode(t, node3, res)

	res, err = db.CreateNode(nil, "default", node4)
	assert.NoError(t, err)
	checkNode(t, node4, res)

	data, err := db.GetNodeByNames(nil, "default", []string{"node_123", "node_abc"})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(data))
	assert.Equal(t, data[0].Name, "node_123")
	assert.Equal(t, data[1].Name, "node_abc")

}
