package service

import (
	"fmt"
	"testing"

	"github.com/baetyl/baetyl-cloud/v2/common"
	ms "github.com/baetyl/baetyl-cloud/v2/mock/service"
	"github.com/baetyl/baetyl-cloud/v2/models"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	v1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
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

	nsvc := nodeService{
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
	cs := nodeService{
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
	ns := nodeService{
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
	ns := nodeService{
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
	ss := nodeService{
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

	mockObject.shadow.EXPECT().List(node.Namespace, gomock.Any()).Return(shadowList, nil).AnyTimes()

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
	ss := nodeService{
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
	mockObject.shadow.EXPECT().List(node.Namespace, gomock.Any()).Return(&models.ShadowList{}, nil)
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

	ss := nodeService{
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

	ns := nodeService{
		shadow: mockObject.shadow,
		node:   mockObject.node,
		app:    mockObject.app,
	}

	namespace := "test"
	name := "node01"
	desire := specV1.Desire{
		common.DesiredApplications: []specV1.AppInfo{
			{
				"app01",
				"1245",
			},
		},
	}

	shadow := genShadowTestCase()
	mockObject.shadow.EXPECT().Get(gomock.Any(), gomock.Any()).Return(shadow, nil)
	mockObject.shadow.EXPECT().UpdateDesire(gomock.Any()).Return(shadow, nil)

	shd, err := ns.UpdateDesire(namespace, name, desire)
	assert.NoError(t, err)
	apps := shd.Desire[common.DesiredApplications].([]specV1.AppInfo)
	assert.Equal(t, 1, len(apps))
	assert.Equal(t, "app01", apps[0].Name)

	mockObject.shadow.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, nil)
	mockObject.shadow.EXPECT().Create(gomock.Any()).Return(shadow, nil)
	//mockObject.modelStorage.EXPECT().UpdateDesire(gomock.Any()).Return(shadow, nil)

	shd, err = ns.UpdateDesire(namespace, name, desire)
	assert.NoError(t, err)
	apps = shd.Desire[common.DesiredApplications].([]specV1.AppInfo)
	assert.Equal(t, 1, len(apps))
	assert.Equal(t, "app01", apps[0].Name)
}

func TestRematchApplicationForNode(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	ns := nodeService{
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
