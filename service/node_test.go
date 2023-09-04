package service

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/baetyl/baetyl-go/v2/log"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	v1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/cachemsg"
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

	mockObject.shadow.EXPECT().Get(nil, node.Namespace, node.Name).Return(shadow, nil).AnyTimes()
	cs, err := NewNodeService(mockObject.conf)
	mockObject.node.EXPECT().GetNode(nil, node.Namespace, node.Name).Return(node, nil)
	assert.NoError(t, err)
	_, err = cs.Get(nil, node.Namespace, node.Name)
	assert.NoError(t, err)

	mockObject.node.EXPECT().GetNode(nil, node.Namespace, node.Name).Return(nil, fmt.Errorf("node not found"))
	n, err := cs.Get(nil, node.Namespace, node.Name)
	assert.Error(t, err)
	assert.Nil(t, n)

	mockObject.node.EXPECT().GetNode(nil, node.Namespace, node.Name).Return(nil, fmt.Errorf("err"))
	_, err = cs.Get(nil, node.Namespace, node.Name)
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
				Attributes: map[string]interface{}{
					specV1.BaetylCoreFrequency: "10",
				},
			},
		},
		Total: 1,
	}

	nsvc := NodeServiceImpl{
		Shadow: mockObject.shadow,
		Node:   mockObject.node,
		Cache:  mockObject.cache,
		logger: log.With(log.Any("service", "node")),
	}

	mockObject.node.EXPECT().ListNode(nil, ns, s).Return(nil, fmt.Errorf("error"))
	_, err := nsvc.List(ns, s)
	assert.Error(t, err)

	mockObject.node.EXPECT().ListNode(nil, ns, s).Return(list, nil)
	mockObject.cache.EXPECT().Exist(gomock.Any()).Return(true, nil).AnyTimes()
	mockObject.cache.EXPECT().GetByte(cachemsg.GetShadowReportCacheKey(list.Items[0].Namespace, list.Items[0].Name)).Return([]byte(`{"apps":[],"sysapps":[]}`), nil).AnyTimes()
	mockObject.cache.EXPECT().GetByte(cachemsg.GetShadowReportTimeCacheKey("default")).Return([]byte("{\"node01\":\""+time.Now().Format(time.RFC3339Nano)+"\"}"), nil).AnyTimes()
	res, err := nsvc.List(ns, s)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(res.Items))
	assert.Equal(t, ns, res.Items[0].Namespace)

	list.Items[0].Attributes = map[string]interface{}{}
	mockObject.node.EXPECT().ListNode(nil, ns, s).Return(list, nil)
	_, err = nsvc.List(ns, s)
	assert.Error(t, err)
}

func TestFilterNodeService_List(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	ns, s := "default", &models.ListOptions{
		NodeOptions: models.NodeOptions{
			Cluster:    "single",
			Ready:      "online",
			CreateSort: "desc",
		},
	}
	list := genNodeList(t, ns)

	nsvc := NodeServiceImpl{
		Shadow: mockObject.shadow,
		Node:   mockObject.node,
		Cache:  mockObject.cache,
		logger: log.With(log.Any("service", "node")),
	}

	mockObject.node.EXPECT().ListNode(nil, ns, s).Return(nil, fmt.Errorf("error"))
	_, err := nsvc.List(ns, s)
	assert.Error(t, err)

	mockObject.node.EXPECT().ListNode(nil, ns, gomock.Any()).Return(&list, nil).AnyTimes()
	mockObject.cache.EXPECT().Exist(gomock.Any()).Return(true, nil).AnyTimes()

	mockObject.cache.EXPECT().GetByte(cachemsg.GetShadowReportCacheKey(list.Items[0].Namespace, list.Items[0].Name)).Return([]byte(`{"apps":[],"sysapps":[]}`), nil).AnyTimes()
	mockObject.cache.EXPECT().GetByte(cachemsg.GetShadowReportCacheKey(list.Items[0].Namespace, list.Items[1].Name)).Return([]byte(`{"apps":[],"sysapps":[]}`), nil).AnyTimes()
	mockObject.cache.EXPECT().GetByte(cachemsg.GetShadowReportTimeCacheKey("default")).Return([]byte("{\"node01\":\""+time.Now().Format(time.RFC3339Nano)+"\",\"node02\":\""+time.Now().Add(-100*time.Second).Format(time.RFC3339Nano)+"\"}"), nil).AnyTimes()

	s = &models.ListOptions{
		NodeOptions: models.NodeOptions{
			Cluster:    "single",
			Ready:      "",
			CreateSort: "",
		},
	}
	res, err := nsvc.List(ns, s)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(res.Items))
	assert.Equal(t, "node02", res.Items[0].Name)

	s = &models.ListOptions{
		NodeOptions: models.NodeOptions{
			Cluster:    "",
			Ready:      "online",
			CreateSort: "",
		},
	}

	list = genNodeList(t, ns)
	res, err = nsvc.List(ns, s)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(res.Items))
	assert.Equal(t, "node01", res.Items[0].Name)

	s = &models.ListOptions{
		NodeOptions: models.NodeOptions{
			Cluster:    "",
			Ready:      "",
			CreateSort: "asc",
		},
	}
	list = genNodeList(t, ns)

	res, err = nsvc.List(ns, s)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(res.Items))
	assert.Equal(t, "node02", res.Items[0].Name)

}

func genNodeList(t *testing.T, ns string) models.NodeList {
	timeCreate, err := time.Parse("2006-01-02 03:04:05", "2023-04-03 03:04:05")
	assert.NoError(t, err)
	return models.NodeList{
		Items: []specV1.Node{
			{
				Name:              "node01",
				Namespace:         ns,
				Cluster:           true,
				CreationTimestamp: timeCreate.Add(100 * time.Second),
				Attributes: map[string]interface{}{
					specV1.BaetylCoreFrequency: "10",
				},
			},
			{
				Name:              "node02",
				Namespace:         ns,
				CreationTimestamp: timeCreate,
				Cluster:           false,
				Attributes: map[string]interface{}{
					specV1.BaetylCoreFrequency: "10",
				},
			},
		},

		Total: 2,
	}
}

func TestDefaultNodeService_Delete(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	mockIndexService := ms.NewMockIndexService(mockObject.ctl)
	cs := NodeServiceImpl{
		IndexService: mockIndexService,
		Shadow:       mockObject.shadow,
		Node:         mockObject.node,
	}

	node := genNodeTestCase()
	mockObject.shadow.EXPECT().Delete(nil, node.Namespace, node.Name).Return(nil).AnyTimes()

	mockObject.node.EXPECT().DeleteNode(nil, node.Namespace, node.Name).Return(fmt.Errorf("error"))
	err := cs.Delete(nil, node.Namespace, node)
	assert.Error(t, err)

	mockObject.node.EXPECT().DeleteNode(nil, node.Namespace, node.Name).Return(nil)
	mockIndexService.EXPECT().RefreshAppsIndexByNode(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("error"))
	err = cs.Delete(nil, node.Namespace, node)
	assert.NoError(t, err)

	mockObject.node.EXPECT().DeleteNode(nil, node.Namespace, node.Name).Return(nil)
	mockIndexService.EXPECT().RefreshAppsIndexByNode(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	err = cs.Delete(nil, node.Namespace, node)
	assert.NoError(t, err)
}

func TestDefaultNodeService_Create(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	mockIndexService := ms.NewMockIndexService(mockObject.ctl)
	mockSysAppService := ms.NewMockSystemAppService(mockObject.ctl)
	ns := NodeServiceImpl{
		IndexService:  mockIndexService,
		SysAppService: mockSysAppService,
		Shadow:        mockObject.shadow,
		Node:          mockObject.node,
		App:           mockObject.app,
		logger:        log.With(log.Any("service", "node")),
	}
	node := genNodeTestCase()
	apps := &models.ApplicationList{
		Items: []models.AppItem{
			{Namespace: node.Namespace, Name: "app01", Version: "1", Selector: "test=example"},
		},
	}
	app1 := &specV1.Application{
		Name:      "baetyl-core",
		Namespace: node.Namespace,
		Selector:  "abc",
	}

	mockObject.node.EXPECT().CreateNode(nil, node.Namespace, node).Return(nil, fmt.Errorf("error"))
	_, err := ns.Create(nil, node.Namespace, node)
	assert.NotNil(t, err)

	mockObject.node.EXPECT().CreateNode(nil, node.Namespace, node).Return(node, nil).AnyTimes()
	mockSysAppService.EXPECT().GenApps(nil, node.Namespace, gomock.Any()).Return(nil, fmt.Errorf("error"))
	_, err = ns.Create(nil, node.Namespace, node)
	assert.NotNil(t, err)

	mockSysAppService.EXPECT().GenApps(nil, node.Namespace, gomock.Any()).Return([]*specV1.Application{app1}, nil).AnyTimes()
	mockObject.app.EXPECT().ListApplication(nil, node.Namespace, gomock.Any()).Return(nil, fmt.Errorf("error"))
	_, err = ns.Create(nil, node.Namespace, node)
	assert.NotNil(t, err)

	mockObject.app.EXPECT().ListApplication(nil, node.Namespace, gomock.Any()).Return(apps, nil).AnyTimes()
	mockObject.shadow.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
	_, err = ns.Create(nil, node.Namespace, node)
	assert.NotNil(t, err)

	mockObject.shadow.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
	mockIndexService.EXPECT().RefreshAppsIndexByNode(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("error"))
	_, err = ns.Create(nil, node.Namespace, node)
	assert.NotNil(t, err)

	mockIndexService.EXPECT().RefreshAppsIndexByNode(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	_, err = ns.Create(nil, node.Namespace, node)
	assert.NoError(t, err)
}

func TestDefaultNodeService_Update(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	mockIndexService := ms.NewMockIndexService(mockObject.ctl)
	ns := NodeServiceImpl{
		IndexService: mockIndexService,
		Shadow:       mockObject.shadow,
		Node:         mockObject.node,
		App:          mockObject.app,
		logger:       log.With(log.Any("service", "node")),
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

	mockObject.shadow.EXPECT().UpdateDesire(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockObject.shadow.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(shadow, nil).AnyTimes()

	mockObject.node.EXPECT().UpdateNode(nil, node.Namespace, []*specV1.Node{node}).Return(nil, fmt.Errorf("error"))
	_, err := ns.Update(node.Namespace, node)
	assert.NotNil(t, err)

	mockObject.node.EXPECT().UpdateNode(nil, node.Namespace, []*specV1.Node{node}).Return([]*specV1.Node{node}, nil).AnyTimes()
	mockIndexService.EXPECT().RefreshAppsIndexByNode(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("error"))
	_, err = ns.Update(node.Namespace, node)
	assert.NotNil(t, err)

	//mockObject.matcher.EXPECT().IsLabelMatch(gomock.Any(), gomock.Any()).Return(true, nil).AnyTimes()
	mockIndexService.EXPECT().RefreshAppsIndexByNode(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockObject.app.EXPECT().ListApplication(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
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
	mockObject.app.EXPECT().ListApplication(gomock.Any(), gomock.Any(), gomock.Any()).Return(appList, nil).AnyTimes()
	mockObject.shadow.EXPECT().UpdateDesire(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	shad, err := ns.Update(node.Namespace, node)
	assert.NoError(t, err)
	assert.Equal(t, node.Name, shad.Name)
}

func TestUpdateNodeAppVersion(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	mockIndexService := ms.NewMockIndexService(mockObject.ctl)
	ss := NodeServiceImpl{
		IndexService: mockIndexService,
		Shadow:       mockObject.shadow,
		Node:         mockObject.node,
		App:          mockObject.app,
	}
	app := &specV1.Application{
		Name:    "appTest",
		Version: "1234",
	}
	node := genNodeTestCase()

	_, err := ss.UpdateNodeAppVersion(nil, node.Namespace, app)
	assert.NoError(t, err)
	app.Selector = "test=example"
	mockObject.node.EXPECT().ListNode(nil, node.Namespace, gomock.Any()).Return(nil, fmt.Errorf("error"))
	_, err = ss.UpdateNodeAppVersion(nil, node.Namespace, app)
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

	var shadows []*models.Shadow
	for i := 0; i < shadowList.Total; i++ {
		shadows = append(shadows, &shadowList.Items[i])
	}
	mockObject.node.EXPECT().ListNode(nil, node.Namespace, gomock.Any()).Return(nodeList, nil).AnyTimes()
	mockObject.shadow.EXPECT().ListShadowByNames(gomock.Any(), gomock.Any(), gomock.Any()).Return(shadows, nil).AnyTimes()
	mockObject.shadow.EXPECT().UpdateDesires(gomock.Any(), gomock.Any()).Return(fmt.Errorf("update error"))
	_, err = ss.UpdateNodeAppVersion(nil, node.Namespace, app)
	assert.NotNil(t, err)

	mockObject.shadow.EXPECT().UpdateDesires(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	app.Labels = map[string]string{
		common.LabelSystem: app.Name,
	}
	_, err = ss.UpdateNodeAppVersion(nil, node.Namespace, app)
	assert.NoError(t, err)
}

func TestDeleteNodeAppVersion(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	mockIndexService := ms.NewMockIndexService(mockObject.ctl)
	ss := NodeServiceImpl{
		IndexService: mockIndexService,
		Shadow:       mockObject.shadow,
		Node:         mockObject.node,
		App:          mockObject.app,
	}
	app := &specV1.Application{
		Name:    "appTest",
		Version: "1234",
	}
	node := genNodeTestCase()

	_, err := ss.DeleteNodeAppVersion(nil, node.Namespace, app)
	assert.NoError(t, err)

	app.Selector = "test=dev"

	mockObject.node.EXPECT().ListNode(nil, node.Namespace, gomock.Any()).Return(nil, fmt.Errorf("error")).Times(1)
	_, err = ss.DeleteNodeAppVersion(nil, node.Namespace, app)
	assert.Equal(t, fmt.Errorf("error"), err)

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

	var shadows []*models.Shadow
	for i := 0; i < shadowList.Total; i++ {
		shadows = append(shadows, &shadowList.Items[i])
	}
	mockObject.node.EXPECT().ListNode(nil, node.Namespace, gomock.Any()).Return(nodeList, nil).AnyTimes()
	mockObject.shadow.EXPECT().ListShadowByNames(gomock.Any(), gomock.Any(), gomock.Any()).Return(shadows, nil).AnyTimes()
	mockObject.shadow.EXPECT().UpdateDesires(gomock.Any(), shadows).Return(fmt.Errorf("error"))
	_, err = ss.DeleteNodeAppVersion(nil, node.Namespace, app)
	assert.Equal(t, fmt.Errorf("error"), err)

	app.Labels = map[string]string{
		common.LabelSystem: app.Name,
	}
	mockObject.node.EXPECT().ListNode(nil, node.Namespace, gomock.Any()).Return(nodeList, nil).AnyTimes()
	mockObject.shadow.EXPECT().UpdateDesires(gomock.Any(), gomock.Any()).Return(nil)
	_, err = ss.DeleteNodeAppVersion(nil, node.Namespace, app)
	assert.NoError(t, err)
}

func TestUpdateReport(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	ss := NodeServiceImpl{
		Shadow: mockObject.shadow,
		Node:   mockObject.node,
		App:    mockObject.app,
	}

	node := &specV1.Node{
		Name:      "node01",
		Namespace: "test",
		Report: specV1.Report{
			common.DesiredApplications: []specV1.AppInfo{
				{
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

	mockObject.shadow.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
	//mockObject.dbStorage.EXPECT().Create(gomock.Any()).Return(nil, nil)

	mockObject.node.EXPECT().GetNode(nil, node.Namespace, node.Name).Return(nil, fmt.Errorf("error"))
	_, err := ss.UpdateReport(node.Namespace, node.Name, node.Report)
	assert.NotNil(t, err)

	mockObject.shadow.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(shadow, nil)
	//mockObject.dbStorage.EXPECT().Create(gomock.Any()).Return(Shadow, nil)
	mockObject.node.EXPECT().GetNode(nil, node.Namespace, node.Name).Return(node, nil)
	mockObject.node.EXPECT().UpdateNode(nil, node.Namespace, gomock.Any()).Return([]*specV1.Node{node}, nil)
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
		Shadow: mockObject.shadow,
		Node:   mockObject.node,
		App:    mockObject.app,
	}

	namespace := "test"
	names := []string{"node01"}
	listErr := errors.New("failed to Found")

	shadows := []*models.Shadow{genShadowTestCase()}
	shadows[0].Name = names[0]
	shadows[0].Namespace = namespace
	app, _ := genAppTestCase()

	mockObject.shadow.EXPECT().ListShadowByNames(gomock.Any(), namespace, names).Return(nil, listErr)

	err := ns.UpdateDesire(nil, namespace, names, app, RefreshNodeDesireByApp)
	assert.Error(t, err, listErr)

	mockObject.shadow.EXPECT().ListShadowByNames(gomock.Any(), namespace, names).Return(shadows, nil)
	mockObject.shadow.EXPECT().UpdateDesires(gomock.Any(), shadows).Return(nil)

	err = ns.UpdateDesire(nil, namespace, names, app, RefreshNodeDesireByApp)
	assert.NoError(t, err)
}

func TestRematchApplicationForNode(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	ns := NodeServiceImpl{
		Shadow: mockObject.shadow,
		Node:   mockObject.node,
		App:    mockObject.app,
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
		Node:   mockObject.node,
		Shadow: mockObject.shadow,
		logger: log.With(log.Any("service", "node")),
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
	mockObject.node.EXPECT().GetNode(nil, gomock.Any(), gomock.Any()).Return(node, nil)
	mockObject.shadow.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(shadow, nil)
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

	mockObject.node.EXPECT().GetNode(nil, gomock.Any(), gomock.Any()).Return(nil, errors.New("failed to get node"))
	_, err = ns.GetNodeProperties("default", "abc")
	assert.Error(t, err)

	mockObject.node.EXPECT().GetNode(nil, gomock.Any(), gomock.Any()).Return(node, nil)
	mockObject.shadow.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("failed to get Shadow"))
	_, err = ns.GetNodeProperties("default", "abc")
	assert.Error(t, err)
}

func TestUpdateNodeProperties(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	ns := NodeServiceImpl{
		Node:   mockObject.node,
		Shadow: mockObject.shadow,
		logger: log.With(log.Any("service", "node")),
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

	mockObject.node.EXPECT().GetNode(nil, gomock.Any(), gomock.Any()).Return(node, nil)
	mockObject.shadow.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(shadow, nil)
	mockObject.shadow.EXPECT().UpdateDesire(gomock.Any(), gomock.Any()).Return(nil)
	mockObject.node.EXPECT().UpdateNode(nil, gomock.Any(), gomock.Any()).Return(nil, nil)
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
	mockObject.node.EXPECT().GetNode(nil, gomock.Any(), gomock.Any()).Return(node, nil)
	mockObject.shadow.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(shadow, nil)
	mockObject.shadow.EXPECT().UpdateDesire(gomock.Any(), gomock.Any()).Return(nil)
	mockObject.node.EXPECT().UpdateNode(nil, gomock.Any(), gomock.Any()).Return(nil, nil)
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

	mockObject.node.EXPECT().GetNode(nil, gomock.Any(), gomock.Any()).Return(nil, errors.New("failed to get node"))
	_, err = ns.UpdateNodeProperties("default", "abc", nodeProps)
	assert.Error(t, err)

	mockObject.node.EXPECT().GetNode(nil, gomock.Any(), gomock.Any()).Return(node, nil)
	mockObject.shadow.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("failed to get Shadow"))
	_, err = ns.UpdateNodeProperties("default", "abc", nodeProps)
	assert.Error(t, err)

	mockObject.node.EXPECT().GetNode(nil, gomock.Any(), gomock.Any()).Return(node, nil)
	mockObject.shadow.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(shadow, nil)
	mockObject.shadow.EXPECT().UpdateDesire(gomock.Any(), gomock.Any()).Return(errors.New("failed to update desire"))
	//mockObject.node.EXPECT().UpdateNode(nil, gomock.Any(), gomock.Any()).Return(nil, nil)
	_, err = ns.UpdateNodeProperties("default", "abc", nodeProps)
	assert.Error(t, err)

	mockObject.node.EXPECT().GetNode(nil, gomock.Any(), gomock.Any()).Return(node, nil)
	mockObject.shadow.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(shadow, nil)
	mockObject.shadow.EXPECT().UpdateDesire(gomock.Any(), gomock.Any()).Return(nil)
	mockObject.node.EXPECT().UpdateNode(nil, gomock.Any(), gomock.Any()).Return(nil, errors.New("failed to update node"))
	_, err = ns.UpdateNodeProperties("default", "abc", nodeProps)
	assert.Error(t, err)
}

func TestUpdateNodeMode(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	node := &v1.Node{}
	mockObject.node.EXPECT().GetNode(nil, gomock.Any(), gomock.Any()).Return(node, nil)
	mockObject.node.EXPECT().UpdateNode(nil, gomock.Any(), gomock.Any()).Return(nil, nil)

	ns := NodeServiceImpl{
		Node:   mockObject.node,
		logger: log.With(log.Any("service", "node")),
	}
	err := ns.UpdateNodeMode("default", "abc", "cloud")
	assert.NoError(t, err)

	mockObject.node.EXPECT().GetNode(nil, gomock.Any(), gomock.Any()).Return(nil, errors.New("failed to get node"))
	err = ns.UpdateNodeMode("default", "abc", "cloud")
	assert.Error(t, err)

	mockObject.node.EXPECT().GetNode(nil, gomock.Any(), gomock.Any()).Return(node, nil)
	mockObject.node.EXPECT().UpdateNode(nil, gomock.Any(), gomock.Any()).Return(nil, errors.New("failed to update node"))
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
