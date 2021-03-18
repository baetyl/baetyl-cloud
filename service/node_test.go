package service

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	v1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/common"
	ms "github.com/baetyl/baetyl-cloud/v2/mock/service"
	"github.com/baetyl/baetyl-cloud/v2/models"
)

func genNodeTestCase() *specV1.Node {
	node := &specV1.Node{
		Namespace: "default",
		Name:      "abc",
		Labels: map[string]string{
			"test": "example",
		},
		Desire: map[string]interface{}{
			"sysapps": []specV1.AppInfo{{
				Name:    "baetyl-core-abc",
				Version: "123",
			}},
		},
	}
	return node
}

func genShadowTestCase() *models.Shadow {
	shadow := &models.Shadow{
		Namespace: "default",
		Name:      "node01",
		Desire: map[string]interface{}{
			"sysapps": []specV1.AppInfo{{
				Name:    "baetyl-core-node01",
				Version: "123",
			}},
		},
	}
	return shadow
}

func TestDefaultNodeService_Get(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	node := genNodeTestCase()
	shadow := genShadowTestCase()

	mockObject.shadow.EXPECT().Get(node.Namespace, node.Name).Return(shadow, nil).AnyTimes()
	cs, err := NewNodeService(mockObject.conf)
	mockObject.node.EXPECT().GetNode(node.Namespace, node.Name).Return(node, nil)
	assert.NoError(t, err)
	_, err = cs.Get(node.Namespace, node.Name)
	assert.NoError(t, err)

	mockObject.node.EXPECT().GetNode(node.Namespace, node.Name).Return(nil, fmt.Errorf("node not found"))
	n, err := cs.Get(node.Namespace, node.Name)
	assert.Error(t, err)
	assert.Nil(t, n)

	mockObject.node.EXPECT().GetNode(node.Namespace, node.Name).Return(nil, fmt.Errorf("err"))
	_, err = cs.Get(node.Namespace, node.Name)
	assert.Error(t, err)
}

func TestDefaultNodeService_List(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	ns, s := "default", &models.ListOptions{}
	list := &models.NodeList{
		Items: []specV1.Node{
			{
				Name:      "node01",
				Namespace: ns,
			},
		},
	}

	shadowList := &models.ShadowList{
		Items: []models.Shadow{},
	}

	nsvc := NodeServiceImpl{
		shadow: mockObject.shadow,
		node:   mockObject.node,
	}

	mockObject.node.EXPECT().ListNode(ns, s).Return(list, nil)
	mockObject.shadow.EXPECT().List(ns, gomock.Any()).Return(shadowList, nil)
	res, err := nsvc.List(ns, s)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(res.Items))
	assert.Equal(t, ns, res.Items[0].Namespace)

	mockObject.node.EXPECT().ListNode(ns, s).Return(nil, fmt.Errorf("error"))
	_, err = nsvc.List(ns, s)
	assert.Error(t, err)

	mockObject.node.EXPECT().ListNode(ns, s).Return(list, nil)
	mockObject.shadow.EXPECT().List(ns, gomock.Any()).Return(nil, fmt.Errorf("error"))
	_, err = nsvc.List(ns, s)
	assert.Error(t, err)
}

func TestDefaultNodeService_Delete(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	mockIndexService := ms.NewMockIndexService(mockObject.ctl)
	cs := NodeServiceImpl{
		indexService: mockIndexService,
		shadow:       mockObject.shadow,
		node:         mockObject.node,
	}

	node := genNodeTestCase()
	mockObject.shadow.EXPECT().Delete(node.Namespace, node.Name).Return(nil).AnyTimes()

	mockObject.node.EXPECT().DeleteNode(node.Namespace, node.Name).Return(fmt.Errorf("error"))
	err := cs.Delete(node.Namespace, node.Name)
	assert.Error(t, err)

	mockObject.node.EXPECT().DeleteNode(node.Namespace, node.Name).Return(nil)
	mockIndexService.EXPECT().RefreshAppsIndexByNode(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("error"))
	err = cs.Delete(node.Namespace, node.Name)
	assert.NoError(t, err)

	mockObject.node.EXPECT().DeleteNode(node.Namespace, node.Name).Return(nil)
	mockIndexService.EXPECT().RefreshAppsIndexByNode(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	err = cs.Delete(node.Namespace, node.Name)
	assert.NoError(t, err)
}

func TestDefaultNodeService_Create(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	mockIndexService := ms.NewMockIndexService(mockObject.ctl)
	ns := NodeServiceImpl{
		indexService: mockIndexService,
		shadow:       mockObject.shadow,
		node:         mockObject.node,
		app:          mockObject.app,
	}
	node := genNodeTestCase()
	shadow := genShadowTestCase()

	mockObject.shadow.EXPECT().Create(gomock.Any()).Return(shadow, nil).AnyTimes()

	mockObject.shadow.EXPECT().Get(gomock.Any(), gomock.Any()).Return(shadow, nil).AnyTimes()

	mockObject.node.EXPECT().CreateNode(node.Namespace, node).Return(nil, fmt.Errorf("error"))
	_, err := ns.Create(node.Namespace, node)
	assert.NotNil(t, err)

	mockObject.node.EXPECT().CreateNode(node.Namespace, node).Return(node, nil)
	mockObject.app.EXPECT().ListApplication(node.Namespace, gomock.Any()).Return(nil, fmt.Errorf("error"))
	_, err = ns.Create(node.Namespace, node)
	assert.NotNil(t, err)

	apps := &models.ApplicationList{
		Items: []models.AppItem{
			{Namespace: node.Namespace, Name: "app01", Version: "1", Selector: "test=example"},
		},
	}

	//mockObject.matcher.EXPECT().IsLabelMatch(gomock.Any(), gomock.Any()).Return(true, nil).AnyTimes()
	mockObject.node.EXPECT().CreateNode(node.Namespace, node).Return(node, nil)
	mockObject.app.EXPECT().ListApplication(node.Namespace, gomock.Any()).Return(apps, nil)
	mockObject.shadow.EXPECT().UpdateDesire(gomock.Any()).Return(nil, fmt.Errorf("error"))
	_, err = ns.Create(node.Namespace, node)
	assert.NotNil(t, err)

	mockObject.node.EXPECT().CreateNode(node.Namespace, node).Return(node, nil)
	mockObject.app.EXPECT().ListApplication(node.Namespace, gomock.Any()).Return(apps, nil)
	mockObject.shadow.EXPECT().UpdateDesire(gomock.Any()).Return(nil, nil)
	mockIndexService.EXPECT().RefreshAppsIndexByNode(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("error"))
	_, err = ns.Create(node.Namespace, node)
	assert.NotNil(t, err)

	mockObject.node.EXPECT().CreateNode(node.Namespace, node).Return(node, nil)
	mockObject.app.EXPECT().ListApplication(node.Namespace, gomock.Any()).Return(apps, nil)
	mockObject.shadow.EXPECT().UpdateDesire(gomock.Any()).Return(nil, nil)
	mockIndexService.EXPECT().RefreshAppsIndexByNode(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	_, err = ns.Create(node.Namespace, node)
	assert.NoError(t, err)
}

func TestDefaultNodeService_Update(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	mockIndexService := ms.NewMockIndexService(mockObject.ctl)
	ns := NodeServiceImpl{
		indexService: mockIndexService,
		shadow:       mockObject.shadow,
		node:         mockObject.node,
		app:          mockObject.app,
	}
	app := &specV1.Application{
		Name:    "appTest",
		Version: "1234",
	}

	node := &specV1.Node{
		Name:      "node01",
		Namespace: "test",
	}

	shadow := genShadowTestCase()

	mockObject.shadow.EXPECT().UpdateDesire(gomock.Any()).Return(shadow, nil).AnyTimes()
	mockObject.shadow.EXPECT().Get(gomock.Any(), gomock.Any()).Return(shadow, nil).AnyTimes()

	mockObject.node.EXPECT().UpdateNode(node.Namespace, node).Return(nil, fmt.Errorf("error"))
	_, err := ns.Update(node.Namespace, node)
	assert.NotNil(t, err)

	mockObject.node.EXPECT().UpdateNode(node.Namespace, node).Return(node, nil).AnyTimes()
	mockIndexService.EXPECT().RefreshAppsIndexByNode(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("error"))
	_, err = ns.Update(node.Namespace, node)
	assert.NotNil(t, err)

	//mockObject.matcher.EXPECT().IsLabelMatch(gomock.Any(), gomock.Any()).Return(true, nil).AnyTimes()
	mockIndexService.EXPECT().RefreshAppsIndexByNode(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockObject.app.EXPECT().ListApplication(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
	_, err = ns.Update(node.Namespace, node)
	assert.NotNil(t, err)

	appList := &models.ApplicationList{Items: []models.AppItem{
		{
			Name:      app.Name,
			Namespace: app.Namespace,
			Version:   app.Version,
			Selector:  "env=test",
			Labels: map[string]string{
				common.LabelSystem: app.Name,
			},
		},
		{
			Name:      "app01",
			Namespace: app.Namespace,
			Version:   app.Version,
			Selector:  "env=test",
		},
		{
			Name:      "app02",
			Namespace: app.Namespace,
			Version:   app.Version,
		},
	}}
	mockObject.app.EXPECT().ListApplication(gomock.Any(), gomock.Any()).Return(appList, nil).AnyTimes()
	mockObject.shadow.EXPECT().UpdateDesire(gomock.Any()).Return(nil, nil).AnyTimes()
	shad, err := ns.Update(node.Namespace, node)
	assert.NoError(t, err)
	assert.Equal(t, node.Name, shad.Name)
}

func TestUpdateNodeAppVersion(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	mockIndexService := ms.NewMockIndexService(mockObject.ctl)
	ss := NodeServiceImpl{
		indexService: mockIndexService,
		shadow:       mockObject.shadow,
		node:         mockObject.node,
		app:          mockObject.app,
	}
	app := &specV1.Application{
		Name:    "appTest",
		Version: "1234",
	}
	node := genNodeTestCase()
	shadow := genShadowTestCase()

	_, err := ss.UpdateNodeAppVersion(node.Namespace, app)
	assert.NoError(t, err)
	app.Selector = "test=example"
	mockObject.node.EXPECT().ListNode(node.Namespace, gomock.Any()).Return(nil, fmt.Errorf("error"))
	_, err = ss.UpdateNodeAppVersion(node.Namespace, app)
	assert.NotNil(t, err)

	nodeList := &models.NodeList{
		Items: []specV1.Node{
			{
				Name: "test01",
			},
			{
				Name:   "test02",
				Desire: map[string]interface{}{},
			},
			{
				Name: "test03",
				Desire: map[string]interface{}{
					common.DesiredApplications: []specV1.AppInfo{
						{
							"appTest",
							"123",
						},
					},
				},
			},
			{
				Name: "test04",
				Desire: map[string]interface{}{
					common.DesiredApplications: []specV1.AppInfo{
						{
							"appTest",
							"1245",
						},
					},
				},
			},
			{
				Name: "test0t",
				Desire: map[string]interface{}{
					common.DesiredApplications: []specV1.AppInfo{
						{
							"test0t",
							"1245",
						},
					},
				},
			},
		},
	}

	shadowList := &models.ShadowList{
		Items: []models.Shadow{
			{
				Name: "test01",
			},
			{
				Name:   "test02",
				Desire: map[string]interface{}{},
			},
			{
				Name: "test03",
				Desire: map[string]interface{}{
					common.DesiredApplications: []specV1.AppInfo{
						{
							"appTest",
							"123",
						},
					},
				},
			},
			{
				Name: "test04",
				Desire: map[string]interface{}{
					common.DesiredApplications: []specV1.AppInfo{
						{
							"appTest",
							"1245",
						},
					},
				},
			},
			{
				Name: "test0t",
				Desire: map[string]interface{}{
					common.DesiredApplications: []specV1.AppInfo{
						{
							"test0t",
							"1245",
						},
					},
				},
			},
		},
	}

	mockObject.shadow.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
	mockObject.node.EXPECT().ListNode(node.Namespace, gomock.Any()).Return(nodeList, nil)
	mockObject.node.EXPECT().UpdateNode(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
	mockObject.shadow.EXPECT().UpdateDesire(gomock.Any()).Return(nil, nil).AnyTimes()
	_, err = ss.UpdateNodeAppVersion(node.Namespace, app)
	assert.NotNil(t, err)

	mockObject.shadow.EXPECT().Get(gomock.Any(), gomock.Any()).Return(&shadowList.Items[2], nil).AnyTimes()
	mockObject.shadow.EXPECT().Get(gomock.Any(), gomock.Any()).Return(shadow, nil).AnyTimes()
	mockObject.node.EXPECT().ListNode(node.Namespace, gomock.Any()).Return(nodeList, nil).AnyTimes()
	mockObject.node.EXPECT().GetNode(gomock.Any(), gomock.Any()).Return(&nodeList.Items[0], nil).AnyTimes()
	_, err = ss.UpdateNodeAppVersion(node.Namespace, app)
	assert.NoError(t, err)

	app.Labels = map[string]string{
		common.LabelSystem: app.Name,
	}
	_, err = ss.UpdateNodeAppVersion(node.Namespace, app)
	assert.NoError(t, err)
}

func TestDeleteNodeAppVersion(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	mockIndexService := ms.NewMockIndexService(mockObject.ctl)
	ss := NodeServiceImpl{
		indexService: mockIndexService,
		shadow:       mockObject.shadow,
		node:         mockObject.node,
		app:          mockObject.app,
	}
	app := &specV1.Application{
		Name:    "appTest",
		Version: "1234",
	}
	node := genNodeTestCase()
	shadow := genShadowTestCase()

	_, err := ss.DeleteNodeAppVersion(node.Namespace, app)
	assert.NoError(t, err)

	app.Selector = "test=dev"

	mockObject.node.EXPECT().ListNode(node.Namespace, gomock.Any()).Return(nil, fmt.Errorf("error")).Times(1)
	_, err = ss.DeleteNodeAppVersion(node.Namespace, app)
	assert.Equal(t, fmt.Errorf("error"), err)

	mockObject.node.EXPECT().ListNode(node.Namespace, gomock.Any()).Return(&models.NodeList{}, nil).Times(1)
	_, err = ss.DeleteNodeAppVersion(node.Namespace, app)
	assert.NoError(t, err)

	nodeList := &models.NodeList{
		Items: []specV1.Node{
			{
				Name: "test01",
			},
			{
				Name:   "test02",
				Desire: map[string]interface{}{},
			},
			{
				Name: "test03",
				Desire: map[string]interface{}{
					common.DesiredApplications: []specV1.AppInfo{
						{
							"appTest",
							"123",
						},
					},
				},
			},
			{
				Name: "test04",
				Desire: map[string]interface{}{
					common.DesiredApplications: []specV1.AppInfo{
						{
							"appTest",
							"1245",
						},
					},
				},
			},
			{
				Name: "test05",
				Desire: map[string]interface{}{
					common.DesiredApplications: []specV1.AppInfo{{
						"test05",
						"1245",
					},
					},
				},
			},
		},
	}
	shadowList := &models.ShadowList{
		Items: []models.Shadow{
			{
				Name: "test01",
			},
			{
				Name:   "test02",
				Desire: map[string]interface{}{},
			},
			{
				Name: "test03",
				Desire: map[string]interface{}{
					common.DesiredApplications: []specV1.AppInfo{
						{
							"appTest",
							"123",
						},
					},
				},
			},
			{
				Name: "test04",
				Desire: map[string]interface{}{
					common.DesiredApplications: []specV1.AppInfo{
						{
							"appTest",
							"1245",
						},
					},
				},
			},
			{
				Name: "test0t",
				Desire: map[string]interface{}{
					common.DesiredApplications: []specV1.AppInfo{
						{
							"test0t",
							"1245",
						},
					},
				},
			},
		},
	}

	mockObject.shadow.EXPECT().List(node.Namespace, gomock.Any()).Return(shadowList, nil).AnyTimes()

	mockObject.shadow.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
	mockObject.node.EXPECT().ListNode(node.Namespace, gomock.Any()).Return(nodeList, nil)
	mockObject.node.EXPECT().UpdateNode(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
	mockObject.shadow.EXPECT().UpdateDesire(gomock.Any()).Return(nil, nil).AnyTimes()
	_, err = ss.DeleteNodeAppVersion(node.Namespace, app)
	assert.Equal(t, fmt.Errorf("error"), err)

	mockObject.shadow.EXPECT().Get(gomock.Any(), gomock.Any()).Return(shadow, nil).AnyTimes()
	mockObject.shadow.EXPECT().Get(gomock.Any(), gomock.Any()).Return(&shadowList.Items[2], nil).AnyTimes()

	mockObject.node.EXPECT().ListNode(node.Namespace, gomock.Any()).Return(nodeList, nil).AnyTimes()
	mockObject.shadow.EXPECT().UpdateDesire(gomock.Any()).Return(nil, nil).AnyTimes()
	_, err = ss.DeleteNodeAppVersion(node.Namespace, app)
	assert.NoError(t, err)

	app.Labels = map[string]string{
		common.LabelSystem: app.Name,
	}
	mockObject.node.EXPECT().ListNode(node.Namespace, gomock.Any()).Return(nodeList, nil).AnyTimes()
	mockObject.shadow.EXPECT().UpdateDesire(gomock.Any()).Return(nil, nil).AnyTimes()
	_, err = ss.DeleteNodeAppVersion(node.Namespace, app)
	assert.NoError(t, err)
}

func TestUpdateReport(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	ss := NodeServiceImpl{
		shadow: mockObject.shadow,
		node:   mockObject.node,
		app:    mockObject.app,
	}

	node := &specV1.Node{
		Name:      "node01",
		Namespace: "test",
		Report: specV1.Report{
			common.DesiredApplications: []specV1.AppInfo{
				specV1.AppInfo{
					Name:    "appTest",
					Version: "1245",
				},
			},
		},
	}

	report := specV1.Report{
		common.DesiredApplications: []specV1.AppInfo{
			{
				Name:    "appTest-1",
				Version: "123",
			},
		},
	}

	shadow := genShadowTestCase()

	mockObject.shadow.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, nil)
	//mockObject.dbStorage.EXPECT().Create(gomock.Any()).Return(nil, nil)

	mockObject.node.EXPECT().GetNode(node.Namespace, node.Name).Return(nil, fmt.Errorf("error"))
	_, err := ss.UpdateReport(node.Namespace, node.Name, node.Report)
	assert.NotNil(t, err)

	mockObject.shadow.EXPECT().Get(gomock.Any(), gomock.Any()).Return(shadow, nil)
	//mockObject.dbStorage.EXPECT().Create(gomock.Any()).Return(shadow, nil)
	mockObject.node.EXPECT().GetNode(node.Namespace, node.Name).Return(node, nil)
	mockObject.node.EXPECT().UpdateNode(node.Namespace, gomock.Any()).Return(node, nil)
	mockObject.shadow.EXPECT().UpdateReport(gomock.Any()).Return(shadow, nil)
	shad, err := ss.UpdateReport(node.Namespace, node.Name, report)
	assert.NoError(t, err)
	assert.Equal(t, node.Name, shad.Name)
	assert.Equal(t, "appTest-1", shad.Report["apps"].([]specV1.AppInfo)[0].Name)
}

func TestNodeMerge(t *testing.T) {
	report1 := specV1.Report{
		"apps": []specV1.AppInfo{
			{
				Name:    "core-zx-033110",
				Version: "690366",
			},
		},
		"appstats": []specV1.AppStats{
			{

				AppInfo: specV1.AppInfo{
					Name:    "core-zx-033110",
					Version: "690366",
				},
				InstanceStats: map[string]specV1.InstanceStats{
					"core-zx-033110": {
						Name: "core-zx-033110",
						Usage: map[string]string{
							"cpu":    "1160103n",
							"memory": "8420Ki",
						},
						Status: "Running",
					},
				},
			},
			{
				AppInfo: specV1.AppInfo{
					Name:    "function-zx-033110",
					Version: "690371",
				},
				InstanceStats: map[string]specV1.InstanceStats{
					"function-zx-033110": {
						Name: "function-zx-033110",
						Usage: map[string]string{
							"cpu":    "19710n",
							"memory": "1696Ki",
						},
						Status: "Running",
					},
				},
			},
		},
		"node": specV1.NodeInfo{
			Hostname:         "docker-desktop",
			Address:          "192.168.65.3",
			Arch:             "amd64",
			KernelVersion:    "4.19.76-linuxkit",
			OS:               "linux",
			ContainerRuntime: "docker://19.3.5",
			MachineID:        "b49d5b1b-1c0a-42a9-9ee5-5cf69f9f8070",
			BootID:           "d2cd79ae-e825-4a31-bf19-ab6e68e300f7",
			SystemUUID:       "dabd4f62-0000-0000-95e1-f0f38b9e9135",
			OSImage:          "Docker Desktop",
		},
		"nodestats": specV1.NodeStats{
			Usage: map[string]string{
				"cpu":    "211466235n",
				"memory": "1112896Ki",
			},
			Capacity: map[string]string{
				"cpu":    "2",
				"memory": "2037620Ki",
			},
		},
	}

	report2 := specV1.Report{
		"apps": []specV1.AppInfo{
			{
				Name:    "core-zx-033110-1",
				Version: "690366",
			},
			{
				Name:    "function-zx-033110-2",
				Version: "690371",
			},
		},
		"appstats": []specV1.AppStats{
			{

				AppInfo: specV1.AppInfo{
					Name:    "core-zx-033110",
					Version: "690366",
				},
				InstanceStats: map[string]specV1.InstanceStats{
					"core-zx-033110": {
						Name: "core-zx-033110",
						Usage: map[string]string{
							"cpu":    "1160103n",
							"memory": "8420Ki",
						},
						Status: "Running",
					},
				},
			},
			{
				AppInfo: specV1.AppInfo{
					Name:    "function-zx-033110",
					Version: "690371",
				},
				InstanceStats: map[string]specV1.InstanceStats{
					"function-zx-033110": {
						Name: "function-zx-033110",
						Usage: map[string]string{
							"cpu":    "19710n",
							"memory": "1696Ki",
						},
						Status: "Running",
					},
				},
			},
		},
		"node": specV1.NodeInfo{
			Hostname:         "docker-desktop",
			Address:          "192.168.65.3",
			Arch:             "amd64",
			KernelVersion:    "4.19.76-linuxkit",
			OS:               "linux",
			ContainerRuntime: "docker://19.3.5",
			MachineID:        "b49d5b1b-1c0a-42a9-9ee5-5cf69f9f8070",
			BootID:           "d2cd79ae-e825-4a31-bf19-ab6e68e300f7",
			SystemUUID:       "dabd4f62-0000-0000-95e1-f0f38b9e9135",
			OSImage:          "Docker Desktop",
		},
		"nodestats": specV1.NodeStats{
			Usage: map[string]string{
				"cpu":    "211466235n",
				"memory": "1112896Ki",
			},
			Capacity: map[string]string{
				"cpu":    "2",
				"memory": "2037620Ki",
			},
		},
	}

	err := report1.Merge(report2)
	assert.NoError(t, err)
	apps := report1["apps"].([]specV1.AppInfo)
	assert.Equal(t, 2, len(apps))
	assert.Equal(t, "core-zx-033110-1", apps[0].Name)
	assert.Equal(t, "function-zx-033110-2", apps[1].Name)
}

func TestUpdateDesired(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	ns := NodeServiceImpl{
		shadow: mockObject.shadow,
		node:   mockObject.node,
		app:    mockObject.app,
	}

	namespace := "test"
	name := "node01"

	shadow := genShadowTestCase()
	shadow.Name = name
	shadow.Namespace = namespace
	app, _ := genAppTestCase()

	mockObject.shadow.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, nil)
	mockObject.shadow.EXPECT().Create(gomock.Any()).Return(shadow, nil)
	mockObject.shadow.EXPECT().UpdateDesire(gomock.Any()).Return(shadow, nil)

	shd, err := ns.UpdateDesire(namespace, name, app, RefreshNodeDesireByApp)
	assert.NoError(t, err)
	apps := shd.Desire[common.DesiredApplications].([]specV1.AppInfo)
	assert.Equal(t, 1, len(apps))
	assert.Equal(t, "abc", apps[0].Name)

	mockObject.shadow.EXPECT().Get(gomock.Any(), gomock.Any()).Return(shadow, nil).AnyTimes()
	mockObject.shadow.EXPECT().UpdateDesire(gomock.Any()).Return(shadow, nil)

	shd, err = ns.UpdateDesire(namespace, name, app, RefreshNodeDesireByApp)
	assert.NoError(t, err)
	apps = shd.Desire[common.DesiredApplications].([]specV1.AppInfo)
	assert.Equal(t, 1, len(apps))
	assert.Equal(t, "abc", apps[0].Name)

	mockObject.shadow.EXPECT().UpdateDesire(gomock.Any()).Return(nil, errors.New("unknown error"))
	shd, err = ns.UpdateDesire(namespace, name, app, RefreshNodeDesireByApp)
	assert.Error(t, err, "unknown error")

	mockObject.shadow.EXPECT().UpdateDesire(gomock.Any()).Return(nil, errors.New(common.ErrUpdateCas)).Times(3)
	shd, err = ns.UpdateDesire(namespace, name, app, RefreshNodeDesireByApp)
	assert.NotEqual(t, err, nil)
}

func TestRematchApplicationForNode(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	ns := NodeServiceImpl{
		shadow: mockObject.shadow,
		node:   mockObject.node,
		app:    mockObject.app,
	}

	apps := &models.ApplicationList{
		Items: []models.AppItem{
			{
				Name:     "app01",
				Selector: "env=test",
			},
			{
				Name:     "app02",
				Selector: "env=dev",
				System:   true,
				Version:  "1",
			},
			{
				Name:     "app03",
				Selector: "env=dev",
				Version:  "2",
			},
			{
				Name: "app04",
			},
		},
	}

	expect := specV1.Desire{
		common.DesiredSysApplications: []v1.AppInfo{
			{
				"app02",
				"1",
			},
		},
		common.DesiredApplications: []v1.AppInfo{
			{
				"app03",
				"2",
			},
		},
	}

	labels := map[string]string{"env": "dev"}

	names := []string{"app02", "app03"}
	//mockObject.matcher.EXPECT().IsLabelMatch("env=dev", labels).Return(true, nil).Times(2)
	//mockObject.matcher.EXPECT().IsLabelMatch("env=test", labels).Return(false, nil)
	desire, appNames := ns.rematchApplicationsForNode(apps, labels)
	assert.Equal(t, expect, desire)
	assert.Equal(t, names, appNames)

}

func TestGetNodeProperties(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	ns := NodeServiceImpl{
		node:   mockObject.node,
		shadow: mockObject.shadow,
	}

	node := &v1.Node{Attributes: map[string]interface{}{
		common.ReportMeta: map[string]interface{}{
			"a": "reportTime",
		},
		common.DesireMeta: map[string]interface{}{
			"a": "desireTime",
		},
	}}
	shadow := &models.Shadow{
		Report: map[string]interface{}{
			common.NodeProps: map[string]interface{}{"a": "1"},
		},
		Desire: map[string]interface{}{
			common.NodeProps: map[string]interface{}{"a": "2"},
		},
	}
	mockObject.node.EXPECT().GetNode(gomock.Any(), gomock.Any()).Return(node, nil)
	mockObject.shadow.EXPECT().Get(gomock.Any(), gomock.Any()).Return(shadow, nil)
	res, err := ns.GetNodeProperties("default", "abc")
	assert.NoError(t, err)
	expect := &models.NodeProperties{
		State: models.NodePropertiesState{
			Report: map[string]interface{}{"a": "1"},
			Desire: map[string]interface{}{"a": "2"},
		},
		Meta: models.NodePropertiesMetadata{
			ReportMeta: map[string]interface{}{"a": "reportTime"},
			DesireMeta: map[string]interface{}{"a": "desireTime"},
		},
	}
	assert.Equal(t, expect, res)

	mockObject.node.EXPECT().GetNode(gomock.Any(), gomock.Any()).Return(nil, errors.New("failed to get node"))
	_, err = ns.GetNodeProperties("default", "abc")
	assert.Error(t, err)

	mockObject.node.EXPECT().GetNode(gomock.Any(), gomock.Any()).Return(node, nil)
	mockObject.shadow.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, errors.New("failed to get shadow"))
	_, err = ns.GetNodeProperties("default", "abc")
	assert.Error(t, err)
}

func TestUpdateNodeProperties(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	ns := NodeServiceImpl{
		node:   mockObject.node,
		shadow: mockObject.shadow,
	}

	node := &v1.Node{Attributes: map[string]interface{}{
		common.ReportMeta: map[string]interface{}{
			"a": "reportTime",
		},
		common.DesireMeta: map[string]interface{}{
			"a": "desireTime",
		},
	}}
	shadow := &models.Shadow{
		Report: map[string]interface{}{
			common.NodeProps: map[string]interface{}{"a": "1"},
		},
		Desire: map[string]interface{}{
			common.NodeProps: map[string]interface{}{"a": "1"},
		},
	}

	nodeProps := &models.NodeProperties{
		State: models.NodePropertiesState{
			Desire: map[string]interface{}{
				"a": "2",
			},
		},
	}

	mockObject.node.EXPECT().GetNode(gomock.Any(), gomock.Any()).Return(node, nil)
	mockObject.shadow.EXPECT().Get(gomock.Any(), gomock.Any()).Return(shadow, nil)
	mockObject.shadow.EXPECT().UpdateDesire(gomock.Any()).Return(nil, nil)
	mockObject.node.EXPECT().UpdateNode(gomock.Any(), gomock.Any()).Return(nil, nil)
	res, err := ns.UpdateNodeProperties("default", "abc", nodeProps)
	assert.NoError(t, err)
	expect := &models.NodeProperties{
		State: models.NodePropertiesState{
			Report: map[string]interface{}{
				"a": "1",
			},
			Desire: map[string]interface{}{
				"a": "2",
			},
		},
		Meta: models.NodePropertiesMetadata{
			ReportMeta: map[string]interface{}{
				"a": "reportTime",
			},
			DesireMeta: map[string]interface{}{},
		},
	}
	// replace desire meta of expect data with desire meta of result data
	expect.Meta.DesireMeta["a"] = res.Meta.DesireMeta["a"]
	assert.Equal(t, expect, res)

	nodeProps = &models.NodeProperties{
		State: models.NodePropertiesState{
			Desire: map[string]interface{}{},
		},
	}
	mockObject.node.EXPECT().GetNode(gomock.Any(), gomock.Any()).Return(node, nil)
	mockObject.shadow.EXPECT().Get(gomock.Any(), gomock.Any()).Return(shadow, nil)
	mockObject.shadow.EXPECT().UpdateDesire(gomock.Any()).Return(nil, nil)
	mockObject.node.EXPECT().UpdateNode(gomock.Any(), gomock.Any()).Return(nil, nil)
	res, err = ns.UpdateNodeProperties("default", "abc", nodeProps)
	assert.NoError(t, err)
	expect = &models.NodeProperties{
		State: models.NodePropertiesState{
			Report: map[string]interface{}{
				"a": "1",
			},
			Desire: map[string]interface{}{},
		},
		Meta: models.NodePropertiesMetadata{
			ReportMeta: map[string]interface{}{
				"a": "reportTime",
			},
			DesireMeta: map[string]interface{}{},
		},
	}
	assert.Equal(t, expect, res)

	mockObject.node.EXPECT().GetNode(gomock.Any(), gomock.Any()).Return(nil, errors.New("failed to get node"))
	_, err = ns.UpdateNodeProperties("default", "abc", nodeProps)
	assert.Error(t, err)

	mockObject.node.EXPECT().GetNode(gomock.Any(), gomock.Any()).Return(node, nil)
	mockObject.shadow.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, errors.New("failed to get shadow"))
	_, err = ns.UpdateNodeProperties("default", "abc", nodeProps)
	assert.Error(t, err)

	mockObject.node.EXPECT().GetNode(gomock.Any(), gomock.Any()).Return(node, nil)
	mockObject.shadow.EXPECT().Get(gomock.Any(), gomock.Any()).Return(shadow, nil)
	mockObject.shadow.EXPECT().UpdateDesire(gomock.Any()).Return(nil, errors.New("failed to update desire"))
	//mockObject.node.EXPECT().UpdateNode(gomock.Any(), gomock.Any()).Return(nil, nil)
	_, err = ns.UpdateNodeProperties("default", "abc", nodeProps)
	assert.Error(t, err)

	mockObject.node.EXPECT().GetNode(gomock.Any(), gomock.Any()).Return(node, nil)
	mockObject.shadow.EXPECT().Get(gomock.Any(), gomock.Any()).Return(shadow, nil)
	mockObject.shadow.EXPECT().UpdateDesire(gomock.Any()).Return(nil, nil)
	mockObject.node.EXPECT().UpdateNode(gomock.Any(), gomock.Any()).Return(nil, errors.New("failed to update node"))
	_, err = ns.UpdateNodeProperties("default", "abc", nodeProps)
	assert.Error(t, err)
}

func TestUpdateNodeMode(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	node := &v1.Node{}
	mockObject.node.EXPECT().GetNode(gomock.Any(), gomock.Any()).Return(node, nil)
	mockObject.node.EXPECT().UpdateNode(gomock.Any(), gomock.Any()).Return(nil, nil)

	ns := NodeServiceImpl{
		node: mockObject.node,
	}
	err := ns.UpdateNodeMode("default", "abc", "cloud")
	assert.NoError(t, err)

	mockObject.node.EXPECT().GetNode(gomock.Any(), gomock.Any()).Return(nil, errors.New("failed to get node"))
	err = ns.UpdateNodeMode("default", "abc", "cloud")
	assert.Error(t, err)

	mockObject.node.EXPECT().GetNode(gomock.Any(), gomock.Any()).Return(node, nil)
	mockObject.node.EXPECT().UpdateNode(gomock.Any(), gomock.Any()).Return(nil, errors.New("failed to update node"))
	err = ns.UpdateNodeMode("default", "abc", "cloud")
	assert.Error(t, err)
}

func copyDesire(src *v1.Desire, dst *v1.Desire) {
	var apps []specV1.AppInfo
	dstApps := src.AppInfos(false)
	if len(dstApps) == 0 {
		return
	}
	for index := range dstApps {
		apps = append(apps, specV1.AppInfo{Name: dstApps[index].Name, Version: dstApps[index].Version})
	}
	dst.SetAppInfos(false, apps)
}

func TestUpdateNodeAppVersion2(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	mockIndexService := ms.NewMockIndexService(mockObject.ctl)
	ss := NodeServiceImpl{
		indexService: mockIndexService,
		shadow:       mockObject.shadow,
		node:         mockObject.node,
		app:          mockObject.app,
	}
	app1 := &specV1.Application{
		Name:     "app01",
		Version:  "2",
		Selector: "test",
	}
	app2 := &specV1.Application{
		Name:     "app02",
		Version:  "2",
		Selector: "test",
	}
	namespace := "test"

	shadowList := &models.ShadowList{
		Items: []models.Shadow{
			{
				Name: "test01",
				Desire: map[string]interface{}{
					common.DesiredApplications: []specV1.AppInfo{
						{
							"app01",
							"1",
						},
						{
							"app02",
							"1",
						},
					},
				},
				DesireVersion: "1",
			},
		},
	}

	var desireLock sync.Mutex
	mockObject.node.EXPECT().ListNode(namespace, gomock.Any()).DoAndReturn(
		func(namespace string, listOptions *models.ListOptions) (*models.NodeList, error) {
			nodeList := &models.NodeList{
				Items: []specV1.Node{
					{
						Name: "test01",
					},
				},
			}
			return nodeList, nil
		}).AnyTimes()

	mockObject.shadow.EXPECT().Get(gomock.Any(), gomock.Any()).DoAndReturn(
		func(namespace, name string) (*models.Shadow, error) {
			newShadow := &models.Shadow{
				Name:   "test01",
				Desire: map[string]interface{}{},
			}
			desireLock.Lock()
			copyDesire(&shadowList.Items[0].Desire, &newShadow.Desire)
			newShadow.DesireVersion = shadowList.Items[0].DesireVersion
			desireLock.Unlock()
			return newShadow, nil
		}).AnyTimes()

	mockObject.shadow.EXPECT().UpdateDesire(gomock.Any()).DoAndReturn(
		func(shadow *models.Shadow) (*models.Shadow, error) {
			defer desireLock.Unlock()
			time.Sleep(time.Duration(10000))
			desireLock.Lock()
			copyDesire(&shadow.Desire, &shadowList.Items[0].Desire)
			shadowList.Items[0].DesireVersion = "2"
			return &shadowList.Items[0], nil
		})
	mockObject.shadow.EXPECT().UpdateDesire(gomock.Any()).DoAndReturn(
		func(shadow *models.Shadow) (*models.Shadow, error) {
			defer desireLock.Unlock()
			time.Sleep(time.Duration(100000))
			desireLock.Lock()
			if shadow.DesireVersion != shadowList.Items[0].DesireVersion {
				return nil, errors.New(common.ErrUpdateCas)
			}
			copyDesire(&shadow.Desire, &shadowList.Items[0].Desire)
			return &shadowList.Items[0], nil
		}).AnyTimes()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		_, err := ss.UpdateNodeAppVersion(namespace, app1)
		assert.NoError(t, err)
	}()

	go func() {
		defer wg.Done()
		_, err := ss.UpdateNodeAppVersion(namespace, app2)
		assert.NoError(t, err)
	}()
	wg.Wait()
	assert.Equal(t, shadowList.Items[0].Desire.AppInfos(false)[0].Version, "2")
	assert.Equal(t, shadowList.Items[0].Desire.AppInfos(false)[1].Version, "2")
}

func TestFilterNodeListByNodeSelector(t *testing.T) {
	list := &models.NodeList{
		Total: 8,
		ListOptions: &models.ListOptions{
			NodeSelector: "1=1",
			Filter: models.Filter{
				PageNo:   1,
				PageSize: 20,
			},
		},
		Items: []specV1.Node{
			{
				Name: "n0",
				Report: map[string]interface{}{
					"node": map[string]interface{}{
						"master": map[string]interface{}{
							"labels": map[string]interface{}{
								"1": "1",
							},
						},
					},
				},
			},
			{
				Name: "n1",
			},
			{
				Name: "n2",
				Report: map[string]interface{}{
					"time": "2021-03-04T18:07:02.958761Z",
				},
			},
			{
				Name: "n3",
				Report: map[string]interface{}{
					"node": "err",
				},
			},
			{
				Name: "n4",
				Report: map[string]interface{}{
					"node": map[string]interface{}{
						"master": "err",
					},
				},
			},
			{
				Name: "n5",
				Report: map[string]interface{}{
					"node": map[string]interface{}{
						"master": map[string]interface{}{
							"arch": "amd64",
						},
					},
				},
			},
			{
				Name: "n6",
				Report: map[string]interface{}{
					"node": map[string]interface{}{
						"master": map[string]interface{}{
							"labels": "err",
						},
					},
				},
			},
			{
				Name: "n7",
				Report: map[string]interface{}{
					"node": map[string]interface{}{
						"master": map[string]interface{}{
							"labels": map[string]interface{}{
								"2": "2",
							},
						},
					},
				},
			},
		},
	}
	expect := &models.NodeList{
		Total: 1,
		ListOptions: &models.ListOptions{
			NodeSelector: "1=1",
			Filter: models.Filter{
				PageNo:   1,
				PageSize: 20,
			},
		},
		Items: []specV1.Node{
			{
				Name: "n0",
				Report: map[string]interface{}{
					"node": map[string]interface{}{
						"master": map[string]interface{}{
							"labels": map[string]interface{}{
								"1": "1",
							},
						},
					},
				},
			},
		},
	}
	res := filterNodeListByNodeSelector(list)
	assert.EqualValues(t, expect, res)
}
