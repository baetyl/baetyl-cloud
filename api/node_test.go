package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/baetyl/baetyl-go/v2/context"
	"github.com/baetyl/baetyl-go/v2/log"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	mf "github.com/baetyl/baetyl-cloud/v2/mock/facade"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	mockPlugin "github.com/baetyl/baetyl-cloud/v2/mock/plugin"
	ms "github.com/baetyl/baetyl-cloud/v2/mock/service"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
	_ "github.com/baetyl/baetyl-cloud/v2/plugin/default/transaction"
	"github.com/baetyl/baetyl-cloud/v2/service"
)

func getMockNode() *specV1.Node {
	mNode := &specV1.Node{
		Namespace: "default",
		Name:      "abc",
		Labels: map[string]string{
			"tag":                "baidu",
			common.LabelNodeName: "abc",
		},
		Attributes: map[string]interface{}{
			specV1.BaetylCoreFrequency: common.DefaultCoreFrequency,
			specV1.KeyAccelerator:      "",
			specV1.KeyOptionalSysApps:  interface{}([]interface{}{"a"}),
			specV1.KeySyncMode:         specV1.CloudMode,
		},
		SysApps: []string{"a"},
		Mode:    specV1.CloudMode,
	}
	return mNode
}

func getMockNode2() *specV1.Node {
	mNode := &specV1.Node{
		Namespace: "default",
		Name:      "abc",
		Labels: map[string]string{
			"tag":                "baidu",
			common.LabelNodeName: "abc",
		},
		Attributes: map[string]interface{}{
			specV1.BaetylCoreFrequency: common.DefaultCoreFrequency,
			specV1.KeyAccelerator:      "",
		},
	}
	return mNode
}

func initNodeAPI(t *testing.T) (*API, *gin.Engine, *gomock.Controller) {
	api := &API{}
	api.log = log.L().With(log.Any("test", "api"))
	api.AppCombinedService = &service.AppCombinedService{}
	router := gin.Default()
	mockCtl := gomock.NewController(t)
	mockIM := func(c *gin.Context) { common.NewContext(c).SetNamespace("default") }
	v1 := router.Group("v1")
	{
		nodes := v1.Group("/nodes")
		nodes.GET("/:name", mockIM, common.Wrapper(api.GetNode))
		nodes.PUT("", mockIM, common.Wrapper(api.GetNodes))
		nodes.GET("/:name/stats", mockIM, common.Wrapper(api.GetNodeStats))
		nodes.GET("/:name/apps", mockIM, common.Wrapper(api.GetAppByNode))
		nodes.PUT("/:name", mockIM, common.Wrapper(api.UpdateNode))
		nodes.DELETE("/:name", mockIM, common.Wrapper(api.DeleteNode))
		nodes.GET("/:name/init", mockIM, common.Wrapper(api.GenInitCmdFromNode))
		nodes.POST("", mockIM, common.Wrapper(api.CreateNode))
		nodes.GET("", mockIM, common.Wrapper(api.ListNode))
		nodes.GET("/:name/deploys", mockIM, common.Wrapper(api.GetNodeDeployHistory))
		nodes.GET("/:name/properties", mockIM, common.Wrapper(api.GetNodeProperties))
		nodes.PUT("/:name/properties", mockIM, common.Wrapper(api.UpdateNodeProperties))
		nodes.PUT("/:name/mode", mockIM, common.Wrapper(api.UpdateNodeMode))
		nodes.PUT("/:name/core/configs", mockIM, common.Wrapper(api.UpdateCoreApp))
		nodes.GET("/:name/core/configs", mockIM, common.Wrapper(api.GetCoreAppConfigs))
		nodes.GET("/:name/core/versions", mockIM, common.Wrapper(api.GetCoreAppVersions))
	}
	return api, router, mockCtl
}

func TestNewAPI(t *testing.T) {
	c := &config.CloudConfig{}
	c.Plugin.Auth = common.RandString(9)
	c.Plugin.PKI = common.RandString(9)
	c.Plugin.Functions = []string{common.RandString(9), common.RandString(9)}
	c.Plugin.Shadow = common.RandString(9)
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	mockObjectStorage := mockPlugin.NewMockObject(mockCtl)
	for _, v := range c.Plugin.Objects {
		plugin.RegisterFactory(v, func() (plugin.Plugin, error) {
			return mockObjectStorage, nil
		})
	}

	mockFunction := mockPlugin.NewMockFunction(mockCtl)
	for _, v := range c.Plugin.Functions {
		plugin.RegisterFactory(v, func() (plugin.Plugin, error) {
			return mockFunction, nil
		})
	}
	mockAuth := mockPlugin.NewMockAuth(mockCtl)
	plugin.RegisterFactory(c.Plugin.Auth, func() (plugin.Plugin, error) {
		return mockAuth, nil
	})
	mockPKI := mockPlugin.NewMockPKI(mockCtl)
	plugin.RegisterFactory(c.Plugin.PKI, func() (plugin.Plugin, error) {
		return mockPKI, nil
	})
	NewAPI(c)
}

func TestGetNode(t *testing.T) {
	api, router, mockCtl := initNodeAPI(t)
	defer mockCtl.Finish()
	sNode := ms.NewMockNodeService(mockCtl)
	api.Node = sNode

	mNode := getMockNode()

	sNode.EXPECT().Get(nil, mNode.Namespace, mNode.Name).Return(mNode, nil)

	// 200
	req, _ := http.NewRequest(http.MethodGet, "/v1/nodes/abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	bytes := w.Body.Bytes()
	assert.Equal(t, string(bytes), "{\"namespace\":\"default\",\"name\":\"abc\",\"createTime\":\"0001-01-01T00:00:00Z\",\"labels\":{\"baetyl-node-name\":\"abc\",\"tag\":\"baidu\"},\"sysApps\":[\"a\"],\"cluster\":false,\"ready\":0,\"mode\":\"cloud\"}\n")

	sNode.EXPECT().Get(nil, mNode.Namespace, mNode.Name).Return(nil, common.Error(common.ErrResourceNotFound))
	// 404
	req, _ = http.NewRequest(http.MethodGet, "/v1/nodes/abc", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req)
	assert.Equal(t, http.StatusNotFound, w2.Code)

	sNode.EXPECT().Get(nil, mNode.Namespace, mNode.Name).Return(nil, fmt.Errorf("error"))
	// 500
	req, _ = http.NewRequest(http.MethodGet, "/v1/nodes/abc", nil)
	w2 = httptest.NewRecorder()
	router.ServeHTTP(w2, req)
	assert.Equal(t, http.StatusInternalServerError, w2.Code)
}

func TestGetNodes(t *testing.T) {
	api, router, mockCtl := initNodeAPI(t)
	defer mockCtl.Finish()
	sNode := ms.NewMockNodeService(mockCtl)
	api.Node = sNode

	mNode := getMockNode()
	mNode2 := getMockNode()
	mNode2.Name = "abc2"
	mNode2.Labels[common.LabelNodeName] = "abc2"

	// 200
	sNode.EXPECT().Get(nil, mNode.Namespace, mNode.Name).Return(mNode, nil).Times(1)
	sNode.EXPECT().Get(nil, mNode.Namespace, mNode2.Name).Return(mNode2, nil).Times(1)
	nodeNames := &models.NodeNames{
		Names: []string{"abc", "abc2"},
	}
	body, err := json.Marshal(nodeNames)
	assert.NoError(t, err)
	req, _ := http.NewRequest(http.MethodPut, "/v1/nodes?batch", bytes.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, string(w.Body.Bytes()), "{\"total\":2,\"items\":[{\"namespace\":\"default\",\"name\":\"abc\",\"createTime\":\"0001-01-01T00:00:00Z\",\"labels\":{\"baetyl-node-name\":\"abc\",\"tag\":\"baidu\"},\"sysApps\":[\"a\"],\"cluster\":false,\"ready\":0,\"mode\":\"cloud\"},{\"namespace\":\"default\",\"name\":\"abc2\",\"createTime\":\"0001-01-01T00:00:00Z\",\"labels\":{\"baetyl-node-name\":\"abc2\",\"tag\":\"baidu\"},\"sysApps\":[\"a\"],\"cluster\":false,\"ready\":0,\"mode\":\"cloud\"}]}\n")

	// 200 ResourceNotFound
	sNode.EXPECT().Get(nil, mNode.Namespace, mNode.Name).Return(mNode, nil).Times(1)
	sNode.EXPECT().Get(nil, mNode.Namespace, "err_abc").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	sNode.EXPECT().Get(nil, mNode.Namespace, mNode2.Name).Return(mNode2, nil).Times(1)
	nodeNames = &models.NodeNames{
		Names: []string{"abc", "err_abc", "abc2"},
	}
	body, err = json.Marshal(nodeNames)
	assert.NoError(t, err)
	req, _ = http.NewRequest(http.MethodPut, "/v1/nodes?batch", bytes.NewReader(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, string(w.Body.Bytes()), "{\"total\":2,\"items\":[{\"namespace\":\"default\",\"name\":\"abc\",\"createTime\":\"0001-01-01T00:00:00Z\",\"labels\":{\"baetyl-node-name\":\"abc\",\"tag\":\"baidu\"},\"sysApps\":[\"a\"],\"cluster\":false,\"ready\":0,\"mode\":\"cloud\"},{\"namespace\":\"default\",\"name\":\"abc2\",\"createTime\":\"0001-01-01T00:00:00Z\",\"labels\":{\"baetyl-node-name\":\"abc2\",\"tag\":\"baidu\"},\"sysApps\":[\"a\"],\"cluster\":false,\"ready\":0,\"mode\":\"cloud\"}]}\n")

	// 400 validate error
	nodeNames = &models.NodeNames{}
	for i := 0; i < 21; i++ {
		nodeNames.Names = append(nodeNames.Names, "abc")
	}
	body, err = json.Marshal(nodeNames)
	assert.NoError(t, err)
	req, _ = http.NewRequest(http.MethodPut, "/v1/nodes?batch", bytes.NewReader(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 400 invalid request param
	assert.NoError(t, err)
	req, _ = http.NewRequest(http.MethodPut, "/v1/nodes", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	//500
	nodeNames = &models.NodeNames{
		Names: []string{"abc", "abc2"},
	}
	body, err = json.Marshal(nodeNames)
	assert.NoError(t, err)
	sNode.EXPECT().Get(nil, mNode.Namespace, mNode.Name).Return(nil, fmt.Errorf("error")).Times(1)
	req, _ = http.NewRequest(http.MethodPut, "/v1/nodes?batch", bytes.NewReader(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetNodeStats(t *testing.T) {
	api, router, mockCtl := initNodeAPI(t)
	defer mockCtl.Finish()
	sNode := ms.NewMockNodeService(mockCtl)
	api.Node = sNode

	mNode := getMockNode()

	sNode.EXPECT().Get(nil, mNode.Namespace, gomock.Any()).Return(mNode, nil)
	// 200
	req, _ := http.NewRequest(http.MethodGet, "/v1/nodes/abc/stats", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	sNode.EXPECT().Get(nil, mNode.Namespace, "cba").Return(nil, common.Error(common.ErrResourceNotFound))
	// 404
	req, _ = http.NewRequest(http.MethodGet, "/v1/nodes/cba/stats", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req)
	assert.Equal(t, http.StatusNotFound, w2.Code)

	sNode.EXPECT().Get(nil, mNode.Namespace, "cba").Return(nil, fmt.Errorf("error"))
	// 500
	req, _ = http.NewRequest(http.MethodGet, "/v1/nodes/cba/stats", nil)
	w2 = httptest.NewRecorder()
	router.ServeHTTP(w2, req)
	assert.Equal(t, http.StatusInternalServerError, w2.Code)

	mNode.Report = map[string]interface{}{
		"nodestats": map[string]interface{}{
			"master": map[string]interface{}{
				"usage": map[string]string{
					"cpu":    "1",
					"memory": "512Mi",
				},
				"capacity": map[string]string{
					"cpu":    "2",
					"memory": "1024Mi",
				},
			},
		},
		"time": "2020-04-13T10:07:12.267728Z",
	}
	sNode.EXPECT().Get(nil, mNode.Namespace, "abc").Return(mNode, nil)
	// 200
	req, _ = http.NewRequest(http.MethodGet, "/v1/nodes/abc/stats", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	mNode.Report = map[string]interface{}{
		"nodestats": map[string]interface{}{
			"master": map[string]interface{}{
				"usage": map[string]string{
					"cpu":    "0.5",
					"memory": "512Mi",
				},
				"capacity": map[string]string{
					"cpu":    "2.5",
					"memory": "1024Mi",
				},
			},
		},
	}
	sNode.EXPECT().Get(nil, mNode.Namespace, "abc").Return(mNode, nil)
	// 200
	req, _ = http.NewRequest(http.MethodGet, "/v1/nodes/abc/stats", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	mNode.Report = map[string]interface{}{
		"nodestats": map[string]interface{}{
			"master": map[string]interface{}{
				"usage": map[string]string{
					"cpu":    "0.5a",
					"memory": "512M",
				},
				"capacity": map[string]string{
					"cpu":    "2.5a",
					"memory": "1024M",
				},
			},
		},
	}
	sNode.EXPECT().Get(nil, mNode.Namespace, "abc").Return(mNode, nil)
	req, _ = http.NewRequest(http.MethodGet, "/v1/nodes/abc/stats", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "quantities must match the regular expression '^([+-]?[0-9.]+)([eEinumkKMGTP]*[-+]?[0-9]*)$'")

	mNode.Report = map[string]interface{}{
		"nodestats": map[string]interface{}{
			"master": map[string]interface{}{
				"usage": map[string]string{
					"cpu":    "0.5",
					"memory": "512a",
				},
				"capacity": map[string]string{
					"cpu":    "2.5",
					"memory": "1024a",
				},
			},
		},
	}
	sNode.EXPECT().Get(nil, mNode.Namespace, "abc").Return(mNode, nil)
	req, _ = http.NewRequest(http.MethodGet, "/v1/nodes/abc/stats", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "quantities must match the regular expression '^([+-]?[0-9.]+)([eEinumkKMGTP]*[-+]?[0-9]*)$'")

	mNode.Report = map[string]interface{}{
		"nodestats": map[string]interface{}{
			"master": map[string]interface{}{
				"usage": map[string]string{
					"cpu":    "0.5a",
					"memory": "512a",
				},
				"capacity": map[string]string{
					"cpu":    "2.5",
					"memory": "1024Mi",
				},
			},
		},
	}
	sNode.EXPECT().Get(nil, mNode.Namespace, "abc").Return(mNode, nil)
	req, _ = http.NewRequest(http.MethodGet, "/v1/nodes/abc/stats", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "quantities must match the regular expression '^([+-]?[0-9.]+)([eEinumkKMGTP]*[-+]?[0-9]*)$'")

	mNode.Report = map[string]interface{}{
		"core": specV1.CoreInfo{
			GoVersion:   "1",
			BinVersion:  "1",
			GitRevision: "1",
		},
		"node":      nil,
		"nodestats": nil,
		"apps":      nil,
		"sysapps":   nil,
		"appstats":  nil,
	}

	sNode.EXPECT().Get(nil, mNode.Namespace, "abc").Return(mNode, nil)
	req, _ = http.NewRequest(http.MethodGet, "/v1/nodes/abc/stats", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestListNode(t *testing.T) {
	api, router, mockCtl := initNodeAPI(t)
	defer mockCtl.Finish()
	sNode := ms.NewMockNodeService(mockCtl)
	api.Node = sNode

	mClist := &models.NodeList{
		Total: 1,
		Items: []specV1.Node{
			{
				Name: "node01",
				Attributes: map[string]interface{}{
					specV1.BaetylCoreFrequency: common.DefaultCoreFrequency,
				},
				Labels: map[string]string{
					"test": "test",
				},
			},
		},
	}

	sNode.EXPECT().List("default", &models.ListOptions{
		NodeSelector: "test=test",
	}).Return(mClist, nil)

	// 200
	req, _ := http.NewRequest(http.MethodGet, "/v1/nodes?nodeSelector=test=test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	bytes := w.Body.Bytes()
	fmt.Println(string(bytes))
	assert.Equal(t, string(bytes), "{\"total\":1,\"items\":[{\"name\":\"node01\",\"createTime\":\"0001-01-01T00:00:00Z\",\"labels\":{\"test\":\"test\"},\"cluster\":false,\"ready\":0,\"mode\":\"\"}]}\n")
	nodelist := new(models.NodeList)
	err := json.Unmarshal(bytes, nodelist)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(nodelist.Items))
	assert.Equal(t, "", nodelist.Items[0].Labels[common.LabelNodeName])
	assert.Equal(t, "node01", nodelist.Items[0].Name)
	assert.Nil(t, nodelist.Items[0].Desire)
}

func TestCreateNode(t *testing.T) {
	api, router, mockCtl := initNodeAPI(t)
	defer mockCtl.Finish()

	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sSecret := ms.NewMockSecretService(mockCtl)
	api.AppCombinedService = &service.AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}
	mLicense := ms.NewMockLicenseService(mockCtl)
	api.License = mLicense

	sNode, sIndex := ms.NewMockNodeService(mockCtl), ms.NewMockIndexService(mockCtl)
	api.Node, api.Index = sNode, sIndex

	sProp := ms.NewMockPropertyService(mockCtl)
	sPKI := ms.NewMockPKIService(mockCtl)
	api.Prop = sProp
	api.PKI = sPKI

	sModule := ms.NewMockModuleService(mockCtl)
	api.Module = sModule

	cfg := &config.CloudConfig{}
	cfg.Plugin.Tx = "defaulttx"
	wrpper, _ := service.NewWrapperService(cfg)
	api.Wrapper = wrpper

	mNode := getMockNode2()

	mLicense.EXPECT().AcquireQuota(mNode.Namespace, plugin.QuotaNode, 1).Return(nil)
	sNode.EXPECT().Get(nil, gomock.Any(), gomock.Any()).Return(nil, nil)
	sNode.EXPECT().Create(nil, mNode.Namespace, gomock.Any()).Return(mNode, nil)
	m := &models.Module{
		Name:    "baetyl",
		Version: "2.1.2",
	}
	sModule.EXPECT().GetLatestModule(gomock.Any()).Return(m, nil).AnyTimes()

	mNode2 := getMockNode2()
	w := httptest.NewRecorder()
	body, _ := json.Marshal(mNode2)
	req, _ := http.NewRequest(http.MethodPost, "/v1/nodes", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	mNode.Name = "node-test"
	mNode.Labels[common.LabelNodeName] = mNode.Name

	sNode.EXPECT().Get(nil, gomock.Any(), gomock.Any()).Return(nil, nil)
	sNode.EXPECT().Create(nil, mNode.Namespace, gomock.Any()).Return(mNode, nil)
	mLicense.EXPECT().AcquireQuota(mNode.Namespace, plugin.QuotaNode, 1).Return(nil)

	w = httptest.NewRecorder()
	body, _ = json.Marshal(mNode)
	req, _ = http.NewRequest(http.MethodPost, "/v1/nodes", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	sNode.EXPECT().Get(nil, gomock.Any(), gomock.Any()).Return(nil, nil)
	sNode.EXPECT().Create(nil, mNode.Namespace, gomock.Any()).Return(nil, fmt.Errorf("create node error"))
	mLicense.EXPECT().AcquireQuota(mNode.Namespace, plugin.QuotaNode, 1).Return(nil)
	mLicense.EXPECT().ReleaseQuota(mNode.Namespace, plugin.QuotaNode, 1).Return(nil)
	w = httptest.NewRecorder()
	body, _ = json.Marshal(mNode)
	req, _ = http.NewRequest(http.MethodPost, "/v1/nodes", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	sNode.EXPECT().Get(nil, gomock.Any(), gomock.Any()).Return(nil, nil)
	mLicense.EXPECT().AcquireQuota(mNode.Namespace, plugin.QuotaNode, 1).Return(fmt.Errorf("quota error"))

	w = httptest.NewRecorder()
	body, _ = json.Marshal(mNode)
	req, _ = http.NewRequest(http.MethodPost, "/v1/nodes", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	sNode.EXPECT().Get(nil, gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
	mNode.Name = "node-baetyl-test"
	mNode.Labels[common.LabelNodeName] = mNode.Name
	w = httptest.NewRecorder()
	body, _ = json.Marshal(mNode)
	req, _ = http.NewRequest(http.MethodPost, "/v1/nodes", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	sNode.EXPECT().Get(nil, gomock.Any(), gomock.Any()).Return(mNode, nil)
	mNode.Name = "node-baetyl-test"
	mNode.Labels[common.LabelNodeName] = mNode.Name
	w = httptest.NewRecorder()
	body, _ = json.Marshal(mNode)
	req, _ = http.NewRequest(http.MethodPost, "/v1/nodes", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mNode.Name = ""
	mNode.Labels[common.LabelNodeName] = mNode.Name
	w = httptest.NewRecorder()
	body, _ = json.Marshal(mNode)
	req, _ = http.NewRequest(http.MethodPost, "/v1/nodes", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "The request parameter is invalid. (name is required)")
}

func TestCreateNodeWithSysApps(t *testing.T) {
	api, router, mockCtl := initNodeAPI(t)
	defer mockCtl.Finish()

	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sSecret := ms.NewMockSecretService(mockCtl)
	api.AppCombinedService = &service.AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}
	mLicense := ms.NewMockLicenseService(mockCtl)
	api.License = mLicense

	sNode, sIndex := ms.NewMockNodeService(mockCtl), ms.NewMockIndexService(mockCtl)
	api.Node, api.Index = sNode, sIndex

	sProp := ms.NewMockPropertyService(mockCtl)
	sPKI := ms.NewMockPKIService(mockCtl)
	sInit := ms.NewMockInitService(mockCtl)
	api.Prop = sProp
	api.PKI = sPKI
	api.Init = sInit

	sModule := ms.NewMockModuleService(mockCtl)
	api.Module = sModule

	cfg := &config.CloudConfig{}
	cfg.Plugin.Tx = "defaulttx"
	wrpper, _ := service.NewWrapperService(cfg)
	api.Wrapper = wrpper

	mNode := &specV1.Node{
		Namespace: "default",
		Name:      "abc",
		Labels: map[string]string{
			"tag":                "baidu",
			common.LabelNodeName: "abc",
		},
		Attributes: map[string]interface{}{
			specV1.BaetylCoreFrequency: common.DefaultCoreFrequency,
			specV1.KeyAccelerator:      "",
		},
	}

	nodeList := []string{"s0", "s1", "s2"}

	sNode.EXPECT().UpdateNodeAppVersion(nil, mNode.Namespace, gomock.Any()).Return(nodeList, nil).AnyTimes()
	sIndex.EXPECT().RefreshNodesIndexByApp(nil, mNode.Namespace, gomock.Any(), nodeList).AnyTimes()
	mLicense.EXPECT().AcquireQuota(mNode.Namespace, plugin.QuotaNode, 1).Return(nil)
	sNode.EXPECT().Get(nil, gomock.Any(), gomock.Any()).Return(nil, nil)
	sNode.EXPECT().Create(nil, mNode.Namespace, gomock.Any()).Return(mNode, nil)
	m := &models.Module{
		Name:    "core",
		Version: "2.1.2",
	}
	sModule.EXPECT().GetLatestModule(gomock.Any()).Return(m, nil).AnyTimes()

	mNode2 := getMockNode2()
	w := httptest.NewRecorder()
	body, _ := json.Marshal(mNode2)
	req, _ := http.NewRequest(http.MethodPost, "/v1/nodes", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	mNode.Name = "node-test"
	mNode.Labels[common.LabelNodeName] = mNode.Name

	sNode.EXPECT().Get(nil, gomock.Any(), gomock.Any()).Return(nil, nil)
	sNode.EXPECT().Create(nil, mNode.Namespace, gomock.Any()).Return(mNode, nil)
	mLicense.EXPECT().AcquireQuota(mNode.Namespace, plugin.QuotaNode, 1).Return(nil)

	w = httptest.NewRecorder()
	body, _ = json.Marshal(mNode)
	req, _ = http.NewRequest(http.MethodPost, "/v1/nodes", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	sNode.EXPECT().Get(nil, gomock.Any(), gomock.Any()).Return(nil, nil)
	sNode.EXPECT().Create(nil, mNode.Namespace, gomock.Any()).Return(nil, fmt.Errorf("create node error"))
	mLicense.EXPECT().AcquireQuota(mNode.Namespace, plugin.QuotaNode, 1).Return(nil)
	mLicense.EXPECT().ReleaseQuota(mNode.Namespace, plugin.QuotaNode, 1).Return(nil)
	w = httptest.NewRecorder()
	body, _ = json.Marshal(mNode)
	req, _ = http.NewRequest(http.MethodPost, "/v1/nodes", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	sNode.EXPECT().Get(nil, gomock.Any(), gomock.Any()).Return(nil, nil)
	mLicense.EXPECT().AcquireQuota(mNode.Namespace, plugin.QuotaNode, 1).Return(fmt.Errorf("quota error"))

	w = httptest.NewRecorder()
	body, _ = json.Marshal(mNode)
	req, _ = http.NewRequest(http.MethodPost, "/v1/nodes", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	sNode.EXPECT().Get(nil, gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
	mNode.Name = "node-baetyl-test"
	mNode.Labels[common.LabelNodeName] = mNode.Name
	w = httptest.NewRecorder()
	body, _ = json.Marshal(mNode)
	req, _ = http.NewRequest(http.MethodPost, "/v1/nodes", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	sNode.EXPECT().Get(nil, gomock.Any(), gomock.Any()).Return(mNode, nil)
	mNode.Name = "node-baetyl-test"
	mNode.Labels[common.LabelNodeName] = mNode.Name
	w = httptest.NewRecorder()
	body, _ = json.Marshal(mNode)
	req, _ = http.NewRequest(http.MethodPost, "/v1/nodes", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mNode.Name = ""
	mNode.Labels[common.LabelNodeName] = mNode.Name
	w = httptest.NewRecorder()
	body, _ = json.Marshal(mNode)
	req, _ = http.NewRequest(http.MethodPost, "/v1/nodes", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "The request parameter is invalid. (name is required)")
}

func TestCreateNodeWithInvalidLabel(t *testing.T) {
	api, router, mockCtl := initNodeAPI(t)
	defer mockCtl.Finish()
	sNode := ms.NewMockNodeService(mockCtl)
	api.Node = sNode
	mNode1 := &specV1.Node{
		Namespace: "default",
		Name:      "abc",
		Labels: map[string]string{
			"+tag": "baidu",
		},
	}

	w := httptest.NewRecorder()
	body, _ := json.Marshal(mNode1)
	req, _ := http.NewRequest(http.MethodPost, "/v1/nodes", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "The field (Labels) must contains labels which can be an empty string or a string which is consist of no more than 63 alphanumeric characters, '-', '_', and must start and end with an alphanumeric character")

	mNode2 := &specV1.Node{
		Namespace: "default",
		Name:      "abc",
		Labels: map[string]string{
			"tag": "+baidu",
		},
	}

	w = httptest.NewRecorder()
	body, _ = json.Marshal(mNode2)
	req, _ = http.NewRequest(http.MethodPost, "/v1/nodes", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "The field (Labels) must contains labels which can be an empty string or a string which is consist of no more than 63 alphanumeric characters, '-', '_', and must start and end with an alphanumeric character")

	mNode3 := &specV1.Node{
		Namespace: "default",
		Name:      "abc",
		Labels: map[string]string{
			"tag+": "baidu",
		},
	}

	w = httptest.NewRecorder()
	body, _ = json.Marshal(mNode3)
	req, _ = http.NewRequest(http.MethodPost, "/v1/nodes", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "The field (Labels) must contains labels which can be an empty string or a string which is consist of no more than 63 alphanumeric characters, '-', '_', and must start and end with an alphanumeric character")

	mNode4 := &specV1.Node{
		Namespace: "default",
		Name:      "abc",
		Labels: map[string]string{
			"tag": "baidu+",
		},
	}

	w = httptest.NewRecorder()
	body, _ = json.Marshal(mNode4)
	req, _ = http.NewRequest(http.MethodPost, "/v1/nodes", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "The field (Labels) must contains labels which can be an empty string or a string which is consist of no more than 63 alphanumeric characters, '-', '_', and must start and end with an alphanumeric character")

	mNode5 := &specV1.Node{
		Namespace: "default",
		Name:      "abc",
		Labels: map[string]string{
			"tag": "J2U2m25qfUzJdFN3xxqiOy0MLhJ5q1vH38d0al8CNH1gqMw8LPJ71hY86S9i3c3d",
		},
	}

	w = httptest.NewRecorder()
	body, _ = json.Marshal(mNode5)
	req, _ = http.NewRequest(http.MethodPost, "/v1/nodes", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "The field (Labels) must contains labels which can be an empty string or a string which is consist of no more than 63 alphanumeric characters, '-', '_', and must start and end with an alphanumeric character")
}

func TestUpdateNodeAddSysApp(t *testing.T) {
	api, router, mockCtl := initNodeAPI(t)
	defer mockCtl.Finish()
	sNode := ms.NewMockNodeService(mockCtl)
	sInit := ms.NewMockInitService(mockCtl)
	sProp := ms.NewMockPropertyService(mockCtl)
	sIndex := ms.NewMockIndexService(mockCtl)
	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sModule := ms.NewMockModuleService(mockCtl)
	sSysApp := ms.NewMockSystemAppService(mockCtl)
	api.Node = sNode
	api.Init = sInit
	api.Prop = sProp
	api.Index = sIndex
	api.App = sApp
	api.Config = sConfig
	api.Module = sModule
	api.SysApp = sSysApp

	mNode := &specV1.Node{
		Namespace: "default",
		Name:      "abc",
		Labels: map[string]string{
			"tag":                   "baidu",
			common.LabelNodeName:    "abc",
			common.LabelNodeMode:    context.RunModeKube,
			"test":                  "test",
			common.LabelAccelerator: "",
			common.LabelCluster:     "false",
		},
		Attributes: map[string]interface{}{
			specV1.BaetylCoreFrequency: common.DefaultCoreFrequency,
			specV1.KeyAccelerator:      "",
		},
	}

	sNode.EXPECT().Get(nil, gomock.Any(), gomock.Any()).Return(mNode, nil).Times(1)
	sNode.EXPECT().Update(mNode.Namespace, mNode).Return(mNode, nil)
	// equal case
	w := httptest.NewRecorder()
	body, _ := json.Marshal(mNode)
	req, _ := http.NewRequest(http.MethodPut, "/v1/nodes/abc", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	mNode2 := &specV1.Node{
		Namespace: "default",
		Name:      "abc",
		Labels: map[string]string{
			"tag":                "baidu",
			common.LabelNodeName: "abc",
			common.LabelNodeMode: context.RunModeKube,
			"test":               "test",
		},
		Attributes: map[string]interface{}{
			specV1.BaetylCoreFrequency: common.DefaultCoreFrequency,
			specV1.KeyAccelerator:      "",
			specV1.KeyOptionalSysApps:  interface{}([]interface{}{"a"}),
		},
		SysApps: []string{"a"},
	}

	sNode.EXPECT().Get(nil, gomock.Any(), gomock.Any()).Return(mNode2, nil).Times(1)
	modules := []models.Module{
		{
			Name:  "a",
			Image: "a-image",
		},
		{
			Name:  "b",
			Image: "b-image",
		},
	}
	sModule.EXPECT().ListModules(&models.Filter{}, gomock.Any()).Return(modules, nil).Times(1)

	mNode3 := &specV1.Node{
		Namespace: "default",
		Name:      "abc",
		Labels: map[string]string{
			"tag":                   "baidu",
			common.LabelNodeName:    "abc",
			common.LabelNodeMode:    context.RunModeKube,
			"test":                  "test",
			common.LabelCluster:     "false",
			common.LabelAccelerator: "",
		},
		Attributes: map[string]interface{}{
			specV1.BaetylCoreFrequency: common.DefaultCoreFrequency,
			specV1.KeyAccelerator:      "",
			specV1.KeyOptionalSysApps:  interface{}([]interface{}{"a"}),
		},
		SysApps: []string{"a"},
	}
	sNode.EXPECT().Update(mNode.Namespace, mNode3).Return(mNode3, nil)
	// equal
	w = httptest.NewRecorder()
	body, _ = json.Marshal(mNode3)
	req, _ = http.NewRequest(http.MethodPut, "/v1/nodes/abc", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	mNode4 := &specV1.Node{
		Namespace: "default",
		Name:      "abc",
		Labels: map[string]string{
			"tag":                   "baidu",
			common.LabelNodeName:    "abc",
			common.LabelNodeMode:    context.RunModeKube,
			"test":                  "test",
			common.LabelAccelerator: "",
			common.LabelCluster:     "false",
		},
		Attributes: map[string]interface{}{
			specV1.BaetylCoreFrequency: common.DefaultCoreFrequency,
			specV1.KeyAccelerator:      "",
		},
	}

	sNode.EXPECT().Get(nil, mNode4.Namespace, mNode4.Name).Return(mNode4, nil).Times(1)
	sModule.EXPECT().ListModules(&models.Filter{}, gomock.Any()).Return(modules, nil).Times(1)
	sSysApp.EXPECT().GenOptionalApps(nil, mNode.Namespace, mNode, []string{"a"}).Times(1)
	nodeList := []string{"abc"}
	sNode.EXPECT().UpdateNodeAppVersion(nil, mNode.Namespace, gomock.Any()).Return(nodeList, nil).AnyTimes()
	sIndex.EXPECT().RefreshNodesIndexByApp(nil, mNode.Namespace, gomock.Any(), nodeList).AnyTimes()

	mNode6 := &specV1.Node{
		Namespace: "default",
		Name:      "abc",
		Labels: map[string]string{
			"tag":                   "baidu",
			common.LabelNodeName:    "abc",
			common.LabelNodeMode:    context.RunModeKube,
			"test":                  "test",
			common.LabelCluster:     "false",
			common.LabelAccelerator: "",
		},
		Attributes: map[string]interface{}{
			specV1.BaetylCoreFrequency: common.DefaultCoreFrequency,
			specV1.KeyAccelerator:      "",
		},
		SysApps: []string{"a"},
	}
	sNode.EXPECT().Update(mNode.Namespace, mNode6).Return(mNode6, nil)

	mNode5 := &specV1.Node{
		Namespace: "default",
		Name:      "abc",
		Labels: map[string]string{
			"tag":                "baidu",
			common.LabelNodeName: "abc",
			common.LabelNodeMode: context.RunModeKube,
			"test":               "test",
		},
		Attributes: map[string]interface{}{
			specV1.BaetylCoreFrequency: common.DefaultCoreFrequency,
			specV1.KeyAccelerator:      "",
		},
		SysApps: []string{"a"},
	}
	w = httptest.NewRecorder()
	body, _ = json.Marshal(mNode5)
	req, _ = http.NewRequest(http.MethodPut, "/v1/nodes/abc", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateNodeDeleteSysApp(t *testing.T) {
	api, router, mockCtl := initNodeAPI(t)
	defer mockCtl.Finish()
	sNode := ms.NewMockNodeService(mockCtl)
	sInit := ms.NewMockInitService(mockCtl)
	sProp := ms.NewMockPropertyService(mockCtl)
	sIndex := ms.NewMockIndexService(mockCtl)
	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	api.Node = sNode
	api.Init = sInit
	api.Prop = sProp
	api.Index = sIndex
	api.App = sApp
	api.Config = sConfig

	desireContent := `
	{
		"sysapps":[{
	      	"name":"rule-node12",
			"version": "12"
		}]
	}`

	desire := specV1.Desire{}
	err := json.Unmarshal([]byte(desireContent), &desire)
	assert.NoError(t, err)

	mNode7 := &specV1.Node{
		Namespace: "default",
		Name:      "abc",
		Labels: map[string]string{
			"tag":                "baidu",
			common.LabelNodeName: "abc",
			common.LabelNodeMode: context.RunModeKube,
			"test":               "test",
		},
		Attributes: map[string]interface{}{
			specV1.BaetylCoreFrequency: common.DefaultCoreFrequency,
			specV1.KeyAccelerator:      "",
		},
		SysApps: []string{"rule-node12"},
		Desire:  desire,
	}

	sNode.EXPECT().Get(nil, mNode7.Namespace, mNode7.Name).Return(mNode7, nil).Times(1)

	appRule := &specV1.Application{
		Name:      "rule-node12",
		Namespace: mNode7.Namespace,
		Version:   "12",
		Labels: map[string]string{
			common.LabelSystem: "true",
		},
		Volumes: []specV1.Volume{
			{
				Name: "config",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name:    "config1",
						Version: "1",
					},
				},
			},
		},
	}

	sApp.EXPECT().Get(mNode7.Namespace, appRule.Name, "").Return(appRule, nil).Times(1)
	sApp.EXPECT().Delete(nil, mNode7.Namespace, appRule.Name, "").Return(nil).Times(1)
	res := &specV1.Configuration{
		Labels: map[string]string{
			common.LabelSystem: "true",
		},
	}
	sConfig.EXPECT().Get(mNode7.Namespace, "config1", "").Return(res, nil).Times(1)
	sConfig.EXPECT().Delete(nil, mNode7.Namespace, appRule.Volumes[0].Config.Name).Times(1)
	sNode.EXPECT().DeleteNodeAppVersion(nil, mNode7.Namespace, gomock.Any()).Return(nil, nil).Times(1)
	sIndex.EXPECT().RefreshNodesIndexByApp(nil, mNode7.Namespace, appRule.Name, gomock.Any()).Return(nil).Times(1)

	mNode9 := &specV1.Node{
		Namespace: "default",
		Name:      "abc",
		Labels: map[string]string{
			"tag":                   "baidu",
			common.LabelNodeName:    "abc",
			common.LabelNodeMode:    context.RunModeKube,
			"test":                  "test",
			common.LabelAccelerator: "",
			common.LabelCluster:     "false",
		},
		Attributes: map[string]interface{}{
			specV1.BaetylCoreFrequency: common.DefaultCoreFrequency,
			specV1.KeyAccelerator:      "",
		},
	}
	sNode.EXPECT().Update(mNode7.Namespace, mNode9).Return(mNode9, nil)

	mNode8 := &specV1.Node{
		Namespace: "default",
		Name:      "abc",
		Labels: map[string]string{
			"tag":                "baidu",
			common.LabelNodeName: "abc",
			common.LabelNodeMode: context.RunModeKube,
			"test":               "test",
		},
		Attributes: map[string]interface{}{
			specV1.BaetylCoreFrequency: common.DefaultCoreFrequency,
			specV1.KeyAccelerator:      "",
		},
	}
	w := httptest.NewRecorder()
	body, _ := json.Marshal(mNode8)
	req, _ := http.NewRequest(http.MethodPut, "/v1/nodes/abc", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateNodeAccelerator(t *testing.T) {
	api, router, mockCtl := initNodeAPI(t)
	defer mockCtl.Finish()
	sNode := ms.NewMockNodeService(mockCtl)
	sInit := ms.NewMockInitService(mockCtl)
	sProp := ms.NewMockPropertyService(mockCtl)
	sIndex := ms.NewMockIndexService(mockCtl)
	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sModule := ms.NewMockModuleService(mockCtl)
	sSysApp := ms.NewMockSystemAppService(mockCtl)
	facade := mf.NewMockFacade(mockCtl)
	api.Node = sNode
	api.Init = sInit
	api.Prop = sProp
	api.Index = sIndex
	api.App = sApp
	api.Config = sConfig
	api.Module = sModule
	api.SysApp = sSysApp
	api.Facade = facade

	mNode := &specV1.Node{
		Namespace: "default",
		Name:      "abc",
		Labels: map[string]string{
			"tag":                "baidu",
			common.LabelNodeName: "abc",
			common.LabelNodeMode: context.RunModeKube,
			"test":               "test",
		},
		Attributes: map[string]interface{}{
			specV1.BaetylCoreFrequency: common.DefaultCoreFrequency,
			specV1.BaetylAgentPort:     common.DefaultAgentPort,
		},
		Desire: specV1.Desire{
			specV1.KeySysApps: []interface{}{
				map[string]interface{}{
					"name":    specV1.BaetylCore,
					"version": "v1",
				},
				map[string]interface{}{
					"name":    specV1.BaetylInit,
					"version": "v1",
				},
			},
		},
	}
	modules := []models.Module{
		{
			Name:  specV1.BaetylGPUMetrics,
			Image: "gpu",
		},
	}
	newNode := &specV1.Node{
		Namespace: "default",
		Name:      "abc",
		Labels: map[string]string{
			"tag":                   "baidu",
			common.LabelNodeName:    "abc",
			common.LabelNodeMode:    context.RunModeKube,
			"test":                  "test",
			common.LabelAccelerator: specV1.NVAccelerator,
			common.LabelCluster:     "false",
		},
		Accelerator: specV1.NVAccelerator,
		Attributes: map[string]interface{}{
			specV1.BaetylCoreFrequency: common.DefaultCoreFrequency,
			specV1.BaetylAgentPort:     common.DefaultAgentPort,
			BaetylCoreLogLevel:         LogLevelDebug,
		},
		Desire: specV1.Desire{
			specV1.KeySysApps: []interface{}{
				map[string]interface{}{
					"name":    specV1.BaetylCore,
					"version": "v1",
				},
				map[string]interface{}{
					"name":    specV1.BaetylInit,
					"version": "v1",
				},
				map[string]interface{}{
					"name":    specV1.BaetylGPUMetrics,
					"version": "v1",
				},
			},
		},
		SysApps: []string{specV1.BaetylGPUMetrics},
	}

	coreApp := &specV1.Application{
		Namespace: "default",
		Name:      specV1.BaetylCore,
		Volumes:   []specV1.Volume{{Name: BaetylCoreConfPrefix, VolumeSource: specV1.VolumeSource{Config: &specV1.ObjectReference{Name: BaetylCoreConfPrefix}}}},
	}
	initApp := &specV1.Application{
		Namespace: "default",
		Name:      specV1.BaetylInit,
		Volumes:   []specV1.Volume{{Name: BaetylInitConfPrefix, VolumeSource: specV1.VolumeSource{Config: &specV1.ObjectReference{Name: BaetylInitConfPrefix}}}},
	}
	coreConf := &specV1.Configuration{
		Name:      specV1.BaetylCore,
		Namespace: "default",
	}
	initConf := &specV1.Configuration{
		Name:      specV1.BaetylInit,
		Namespace: "default",
	}

	sNode.EXPECT().Get(nil, gomock.Any(), gomock.Any()).Return(mNode, nil).Times(1)
	sIndex.EXPECT().ListAppsByNode("default", newNode.Name).Return([]string{specV1.BaetylGPUMetrics, specV1.BaetylCore, specV1.BaetylInit}, nil)
	sModule.EXPECT().ListModules(&models.Filter{}, gomock.Any()).Return(modules, nil).Times(1)
	sSysApp.EXPECT().GenOptionalApps(nil, mNode.Namespace, mNode, []string{specV1.BaetylGPUMetrics}).Times(1)
	nodeList := []string{"abc"}
	sNode.EXPECT().UpdateNodeAppVersion(nil, mNode.Namespace, gomock.Any()).Return(nodeList, nil).AnyTimes()
	sIndex.EXPECT().RefreshNodesIndexByApp(nil, mNode.Namespace, gomock.Any(), nodeList).AnyTimes()

	sApp.EXPECT().Get("default", gomock.Any(), gomock.Any()).Return(coreApp, nil)
	sConfig.EXPECT().Get(coreApp.Namespace, gomock.Any(), gomock.Any()).Return(coreConf, nil)
	sApp.EXPECT().Get("default", gomock.Any(), gomock.Any()).Return(initApp, nil)
	sConfig.EXPECT().Get(initApp.Namespace, gomock.Any(), gomock.Any()).Return(initConf, nil)
	coreConfBs, _ := json.Marshal(coreConf)
	sInit.EXPECT().GetResource(gomock.Any(), mNode.Name, service.TemplateCoreConfYaml, gomock.Any()).Return(coreConfBs, nil)
	sConfig.EXPECT().Update(nil, coreConf.Namespace, coreConf).Return(coreConf, nil)
	sApp.EXPECT().Update(nil, gomock.Any(), coreApp).Return(coreApp, nil)

	initConfBs, _ := json.Marshal(initConf)
	sInit.EXPECT().GetResource(gomock.Any(), mNode.Name, service.TemplateInitConfYaml, gomock.Any()).Return(initConfBs, nil)
	sConfig.EXPECT().Update(nil, initConf.Namespace, initConf).Return(initConf, nil)
	sApp.EXPECT().Update(nil, gomock.Any(), initApp).Return(initApp, nil)

	sNode.EXPECT().Update(newNode.Namespace, newNode).Return(newNode, nil)
	// equal case
	w := httptest.NewRecorder()
	body, _ := json.Marshal(newNode)
	req, _ := http.NewRequest(http.MethodPut, "/v1/nodes/abc", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteNode(t *testing.T) {
	api, router, mockCtl := initNodeAPI(t)
	defer mockCtl.Finish()

	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sSecret := ms.NewMockSecretService(mockCtl)
	api.AppCombinedService = &service.AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}
	mLicense := ms.NewMockLicenseService(mockCtl)
	api.License = mLicense

	sNode, sIndex := ms.NewMockNodeService(mockCtl), ms.NewMockIndexService(mockCtl)
	api.Node, api.Index = sNode, sIndex

	sPKI := ms.NewMockPKIService(mockCtl)
	api.PKI = sPKI

	mNode := &specV1.Node{
		Namespace: "default",
		Name:      "abc",
		Desire:    genDesireOfSysApps(),
	}
	appCore := &specV1.Application{
		Name:      "core-node12",
		Namespace: mNode.Namespace,
		Version:   "12",
		Labels: map[string]string{
			common.LabelSystem: "true",
		},
		Volumes: []specV1.Volume{
			{
				Name: "config",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name:    "config1",
						Version: "1",
					},
				},
			},
			{
				Name: "secret",
				VolumeSource: specV1.VolumeSource{
					Secret: &specV1.ObjectReference{
						Name:    "secret1",
						Version: "1",
					},
				},
			},
		},
	}

	appFunction := &specV1.Application{
		Name:      "function-node12",
		Namespace: mNode.Namespace,
		Version:   "13",
		Labels: map[string]string{
			common.LabelSystem: "true",
		},
		Volumes: []specV1.Volume{
			{
				Name: "configf",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name:    "config1f",
						Version: "1",
					},
				},
			},
			{
				Name: "secretf",
				VolumeSource: specV1.VolumeSource{
					Secret: &specV1.ObjectReference{
						Name:    "secret1f",
						Version: "1",
					},
				},
			},
		},
	}

	secret1 := &specV1.Secret{
		Name:      "secret1",
		Namespace: mNode.Namespace,
		Labels: map[string]string{
			common.LabelSystem: "true",
			specV1.SecretLabel: specV1.SecretConfig,
		},
		Annotations: map[string]string{
			common.AnnotationPkiCertID: "certId1",
		},
	}
	secret1f := &specV1.Secret{
		Name:      "secret1f",
		Namespace: mNode.Namespace,
		Labels: map[string]string{
			common.LabelSystem: "true",
			specV1.SecretLabel: specV1.SecretConfig,
		},
		Annotations: map[string]string{
			common.AnnotationPkiCertID: "certId1f",
		},
	}

	sNode.EXPECT().Get(nil, gomock.Any(), gomock.Any()).Return(mNode, nil).Times(1)
	sNode.EXPECT().Delete(mNode.Namespace, mNode).Return(nil).Times(1)
	sApp.EXPECT().Get(mNode.Namespace, appCore.Name, "").Return(appCore, nil).Times(1)
	sApp.EXPECT().Delete(nil, mNode.Namespace, appCore.Name, "").Return(nil).Times(1)
	sIndex.EXPECT().RefreshNodesIndexByApp(nil, mNode.Namespace, appCore.Name, gomock.Any()).Return(nil).Times(1)
	sConfig.EXPECT().Delete(nil, mNode.Namespace, appCore.Volumes[0].Config.Name).Times(1)
	sSecret.EXPECT().Get(mNode.Namespace, appCore.Volumes[1].Secret.Name, "").Return(secret1, nil).Times(1)
	sPKI.EXPECT().DeleteClientCertificate("certId1").Return(nil).Times(1)
	sSecret.EXPECT().Delete(mNode.Namespace, appCore.Volumes[1].Secret.Name).Times(1)

	sApp.EXPECT().Get(mNode.Namespace, appFunction.Name, "").Return(appFunction, nil).Times(1)
	sApp.EXPECT().Delete(nil, mNode.Namespace, appFunction.Name, "").Return(nil).Times(1)
	sIndex.EXPECT().RefreshNodesIndexByApp(nil, mNode.Namespace, appFunction.Name, gomock.Any()).Return(nil).Times(1)
	sConfig.EXPECT().Delete(nil, mNode.Namespace, appFunction.Volumes[0].Config.Name).Times(1)
	sSecret.EXPECT().Get(mNode.Namespace, appFunction.Volumes[1].Secret.Name, "").Return(secret1f, nil).Times(1)
	sPKI.EXPECT().DeleteClientCertificate("certId1f").Return(nil).Times(1)
	sSecret.EXPECT().Delete(mNode.Namespace, appFunction.Volumes[1].Secret.Name).Times(1)

	mLicense.EXPECT().ReleaseQuota(mNode.Namespace, plugin.QuotaNode, 1).Return(nil).AnyTimes()

	res := &specV1.Configuration{
		Labels: map[string]string{
			common.LabelSystem: "true",
		},
	}
	sConfig.EXPECT().Get(mNode.Namespace, "config1", "").Return(res, nil).Times(1)

	sConfig.EXPECT().Get(mNode.Namespace, "config1f", "").Return(res, nil).Times(1)

	// 200
	req, _ := http.NewRequest(http.MethodDelete, "/v1/nodes/abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	sNode.EXPECT().Get(nil, gomock.Any(), gomock.Any()).Return(nil, errors.New("error")).Times(1)
	req, _ = http.NewRequest(http.MethodDelete, "/v1/nodes/abc", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	sNode.EXPECT().Get(nil, gomock.Any(), gomock.Any()).Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	req, _ = http.NewRequest(http.MethodDelete, "/v1/nodes/abc", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	sNode.EXPECT().Get(nil, gomock.Any(), gomock.Any()).Return(mNode, nil).Times(1)
	sNode.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(errors.New("error"))
	req, _ = http.NewRequest(http.MethodDelete, "/v1/nodes/abc", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeleteNodeError(t *testing.T) {
	api, router, mockCtl := initNodeAPI(t)
	defer mockCtl.Finish()

	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sSecret := ms.NewMockSecretService(mockCtl)
	api.AppCombinedService = &service.AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}
	mLicense := ms.NewMockLicenseService(mockCtl)
	api.License = mLicense

	sNode, sIndex := ms.NewMockNodeService(mockCtl), ms.NewMockIndexService(mockCtl)
	api.Node, api.Index = sNode, sIndex

	sPKI := ms.NewMockPKIService(mockCtl)
	api.PKI = sPKI

	mNode := &specV1.Node{
		Namespace: "default",
		Name:      "abc",
		Desire:    genDesireOfSysApps(),
	}
	appCore := &specV1.Application{
		Name:      "core-node12",
		Namespace: mNode.Namespace,
		Version:   "12",
		Labels: map[string]string{
			common.LabelSystem: "true",
		},
		Volumes: []specV1.Volume{
			{
				Name: "config",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name:    "config1",
						Version: "1",
					},
				},
			},
			{
				Name: "secret",
				VolumeSource: specV1.VolumeSource{
					Secret: &specV1.ObjectReference{
						Name:    "secret1",
						Version: "1",
					},
				},
			},
		},
	}

	appFunction := &specV1.Application{
		Name:      "function-node12",
		Namespace: mNode.Namespace,
		Version:   "13",
		Labels: map[string]string{
			common.LabelSystem: "true",
		},
		Volumes: []specV1.Volume{
			{
				Name: "configf",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name:    "config1f",
						Version: "1",
					},
				},
			},
			{
				Name: "secretf",
				VolumeSource: specV1.VolumeSource{
					Secret: &specV1.ObjectReference{
						Name:    "secret1f",
						Version: "1",
					},
				},
			},
		},
	}

	sNode.EXPECT().Get(nil, gomock.Any(), gomock.Any()).Return(mNode, nil).Times(1)
	sNode.EXPECT().Delete(mNode.Namespace, mNode).Return(nil).Times(1)
	sApp.EXPECT().Get(mNode.Namespace, appCore.Name, "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	sApp.EXPECT().Get(mNode.Namespace, appFunction.Name, "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	mLicense.EXPECT().ReleaseQuota(mNode.Namespace, plugin.QuotaNode, 1).Return(nil).AnyTimes()
	sIndex.EXPECT().RefreshNodesIndexByApp(nil, mNode.Namespace, appCore.Name, gomock.Any()).Return(nil).Times(1)
	sIndex.EXPECT().RefreshNodesIndexByApp(nil, mNode.Namespace, appFunction.Name, gomock.Any()).Return(nil).Times(1)

	// 200
	req, _ := http.NewRequest(http.MethodDelete, "/v1/nodes/abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	secret1 := &specV1.Secret{
		Name:      "secret1",
		Namespace: mNode.Namespace,
		Labels: map[string]string{
			specV1.SecretLabel: specV1.SecretConfig,
			common.LabelSystem: "true",
		},
		Annotations: map[string]string{
			common.AnnotationPkiCertID: "certId1",
		},
	}

	sNode.EXPECT().Get(nil, gomock.Any(), gomock.Any()).Return(mNode, nil).Times(1)
	sNode.EXPECT().Delete(mNode.Namespace, mNode).Return(nil).Times(1)
	sApp.EXPECT().Get(mNode.Namespace, appCore.Name, "").Return(appCore, nil).Times(1)
	sApp.EXPECT().Delete(nil, mNode.Namespace, appCore.Name, "").Return(errors.New("error")).Times(1)
	sIndex.EXPECT().RefreshNodesIndexByApp(nil, mNode.Namespace, appCore.Name, gomock.Any()).Return(errors.New("error")).Times(1)
	sConfig.EXPECT().Delete(nil, mNode.Namespace, appCore.Volumes[0].Config.Name).Return(errors.New("error")).Times(1)
	sSecret.EXPECT().Get(mNode.Namespace, appCore.Volumes[1].Secret.Name, "").Return(secret1, nil).Times(1)
	sPKI.EXPECT().DeleteClientCertificate("certId1").Return(errors.New("error")).Times(1)
	sSecret.EXPECT().Delete(mNode.Namespace, appCore.Volumes[1].Secret.Name).Times(1)

	sApp.EXPECT().Get(mNode.Namespace, appFunction.Name, "").Return(appFunction, nil).Times(1)
	sApp.EXPECT().Delete(nil, mNode.Namespace, appFunction.Name, "").Return(errors.New("error")).Times(1)
	sIndex.EXPECT().RefreshNodesIndexByApp(nil, mNode.Namespace, appFunction.Name, gomock.Any()).Return(errors.New("error")).Times(1)
	sConfig.EXPECT().Delete(nil, mNode.Namespace, appFunction.Volumes[0].Config.Name).Return(errors.New("error")).Times(1)
	sSecret.EXPECT().Get(mNode.Namespace, appFunction.Volumes[1].Secret.Name, "").Return(nil, errors.New("error")).Times(1)

	res := &specV1.Configuration{
		Labels: map[string]string{
			common.LabelSystem: "true",
		},
	}
	sConfig.EXPECT().Get(mNode.Namespace, "config1", "").Return(res, nil).Times(1)

	sConfig.EXPECT().Get(mNode.Namespace, "config1f", "").Return(res, nil).Times(1)

	// 200
	req2, _ := http.NewRequest(http.MethodDelete, "/v1/nodes/abc", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)

	sNode.EXPECT().Get(nil, gomock.Any(), gomock.Any()).Return(mNode, nil).Times(1)
	sNode.EXPECT().Delete(mNode.Namespace, mNode).Return(nil).Times(1)
	sApp.EXPECT().Get(mNode.Namespace, appCore.Name, "").Return(nil, errors.New("error")).Times(1)
	sIndex.EXPECT().RefreshNodesIndexByApp(nil, mNode.Namespace, appCore.Name, gomock.Any()).Return(errors.New("error")).Times(1)

	sApp.EXPECT().Get(mNode.Namespace, appFunction.Name, "").Return(appFunction, nil).Times(1)
	sApp.EXPECT().Delete(nil, mNode.Namespace, appFunction.Name, "").Return(errors.New("error")).Times(1)
	sIndex.EXPECT().RefreshNodesIndexByApp(nil, mNode.Namespace, appFunction.Name, gomock.Any()).Return(errors.New("error")).Times(1)
	sConfig.EXPECT().Delete(nil, mNode.Namespace, appFunction.Volumes[0].Config.Name).Return(errors.New("error")).Times(1)
	sSecret.EXPECT().Get(mNode.Namespace, appFunction.Volumes[1].Secret.Name, "").Return(nil, errors.New("error")).Times(1)

	sConfig.EXPECT().Get(mNode.Namespace, "config1f", "").Return(res, nil).Times(1)

	// 200
	req3, _ := http.NewRequest(http.MethodDelete, "/v1/nodes/abc", nil)
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusOK, w3.Code)
}

func genDesireOfSysApps() specV1.Desire {
	content := `
	{
		"sysapps":[{
	       	"name":"core-node12",
			"version": "12"
		},
		{
	       	"name":"function-node12",
			"version": "12"
		}]
	}`

	desire := specV1.Desire{}
	json.Unmarshal([]byte(content), &desire)
	return desire
}

func TestAPI_GetNodeDeployHistory(t *testing.T) {
	api, router, mockCtl := initNodeAPI(t)
	defer mockCtl.Finish()
	sNode := ms.NewMockNodeService(mockCtl)
	api.Node = sNode

	req, _ := http.NewRequest(http.MethodGet, "/v1/nodes/abc/deploys", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestGenInitCmdFromNode(t *testing.T) {
	api, router, mockCtl := initNodeAPI(t)
	defer mockCtl.Finish()

	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sSecret := ms.NewMockSecretService(mockCtl)
	api.AppCombinedService = &service.AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}

	sNode := ms.NewMockNodeService(mockCtl)
	api.Node = sNode

	sInit := ms.NewMockInitService(mockCtl)
	api.Init = sInit

	node := getMockNode()
	params := map[string]interface{}{
		"InitApplyYaml": "baetyl-init-deployment.yml",
		"mode":          "kube",
		"template":      service.TemplateBaetylInitCommand,
	}
	var expect interface{} = []byte("setup")
	sInit.EXPECT().GetResource("default", "abc", service.TemplateBaetylInitCommand, params).Return(expect, nil).Times(1)
	sNode.EXPECT().Get(nil, node.Namespace, node.Name).Return(node, nil).Times(1)

	req, _ := http.NewRequest(http.MethodGet, "/v1/nodes/abc/init", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGenInitCmdFromNode_ErrNode(t *testing.T) {
	api, router, mockCtl := initNodeAPI(t)
	defer mockCtl.Finish()
	nMock := ms.NewMockNodeService(mockCtl)
	api.Node = nMock

	node := getMockNode()

	nMock.EXPECT().Get(nil, node.Namespace, node.Name).Return(nil,
		common.Error(common.ErrResourceNotFound,
			common.Field("type", "nodes"),
			common.Field("name", node.Name))).Times(1)

	req, _ := http.NewRequest(http.MethodGet, "/v1/nodes/abc/init", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetAppByNode(t *testing.T) {
	api, router, mockCtl := initNodeAPI(t)
	defer mockCtl.Finish()

	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sSecret := ms.NewMockSecretService(mockCtl)
	api.AppCombinedService = &service.AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}

	sNode, sIndex := ms.NewMockNodeService(mockCtl), ms.NewMockIndexService(mockCtl)
	api.Node, api.Index = sNode, sIndex

	sInit := ms.NewMockInitService(mockCtl)
	api.Init = sInit

	appNames := []string{"app1", "app2", "app3"}
	sysAppNames := []string{"sysapp1", "sysapp2", "sysapp3"}
	appinfos := []specV1.AppInfo{
		{
			Name:    appNames[0],
			Version: "v1",
		},
		{
			Name:    appNames[1],
			Version: "v1",
		},
		{
			Name:    appNames[2],
			Version: "v1",
		},
	}

	apps := []*specV1.Application{
		{
			Name:    appNames[0],
			Version: "v1",
		},
		{
			Name:    appNames[1],
			Version: "v1",
		},
		{
			Name:    appNames[2],
			Version: "v1",
		},
	}

	sysappinfos := []specV1.AppInfo{
		{
			Name:    sysAppNames[0],
			Version: "v1",
		},
		{
			Name:    sysAppNames[1],
			Version: "v1",
		},
		{
			Name:    sysAppNames[2],
			Version: "v1",
		},
	}

	sysapps := []*specV1.Application{
		{
			Name:    sysAppNames[0],
			Version: "v1",
		},
		{
			Name:    sysAppNames[1],
			Version: "v1",
		},
		{
			Name:    sysAppNames[2],
			Version: "v1",
		},
	}
	node := &specV1.Node{
		Namespace:   "default",
		Name:        "abc",
		Description: "haha",
		Desire:      specV1.Desire{},
	}

	sNode.EXPECT().Get(nil, gomock.Any(), gomock.Any()).Return(node, nil).AnyTimes()

	w4 := httptest.NewRecorder()
	req4, _ := http.NewRequest(http.MethodGet, "/v1/nodes/abc/apps", nil)
	router.ServeHTTP(w4, req4)
	assert.Equal(t, http.StatusOK, w4.Code)
	list := &models.ApplicationList{}
	json.Unmarshal(w4.Body.Bytes(), list)
	assert.Equal(t, 0, list.Total)

	node.Desire.SetAppInfos(true, sysappinfos)
	node.Desire.SetAppInfos(false, appinfos)

	sApp.EXPECT().Get(node.Namespace, appNames[0], "").Return(apps[0], nil).AnyTimes()
	sApp.EXPECT().Get(node.Namespace, appNames[1], "").Return(apps[1], nil).AnyTimes()
	sApp.EXPECT().Get(node.Namespace, appNames[2], "").Return(nil, common.Error(common.ErrResourceNotFound)).AnyTimes()
	sApp.EXPECT().Get(node.Namespace, sysAppNames[0], "").Return(sysapps[0], nil).AnyTimes()
	sApp.EXPECT().Get(node.Namespace, sysAppNames[1], "").Return(sysapps[1], nil).AnyTimes()
	sApp.EXPECT().Get(node.Namespace, sysAppNames[2], "").Return(nil, common.Error(common.ErrResourceNotFound)).AnyTimes()

	w4 = httptest.NewRecorder()
	req4, _ = http.NewRequest(http.MethodGet, "/v1/nodes/abc/apps", nil)
	router.ServeHTTP(w4, req4)
	assert.Equal(t, http.StatusOK, w4.Code)
	json.Unmarshal(w4.Body.Bytes(), list)
	assert.Equal(t, 4, list.Total)

	w4 = httptest.NewRecorder()
	req4, _ = http.NewRequest(http.MethodGet, "/v1/nodes/abc/apps?selector="+common.LabelSystem+"=true", nil)
	router.ServeHTTP(w4, req4)
	assert.Equal(t, http.StatusOK, w4.Code)
	json.Unmarshal(w4.Body.Bytes(), list)
	assert.Equal(t, 4, list.Total)
}

func TestAPI_NodeNumberCollector(t *testing.T) {
	api, _, mockCtl := initNodeAPI(t)
	defer mockCtl.Finish()
	sNode := ms.NewMockNodeService(mockCtl)
	api.Node = sNode
	namespace := "iotCoreId"

	sNode.EXPECT().Count(namespace).Return(nil, errors.New("error"))
	_, err := api.NodeNumberCollector(namespace)
	assert.Error(t, err)

	list := map[string]int{
		plugin.QuotaNode: 12,
	}
	sNode.EXPECT().Count(namespace).Return(list, nil)
	res, err := api.NodeNumberCollector(namespace)
	assert.NoError(t, err)
	assert.Equal(t, 12, res[plugin.QuotaNode])
}

func TestGetNodeProperties(t *testing.T) {
	api, router, mockCtl := initNodeAPI(t)
	defer mockCtl.Finish()

	sNode := ms.NewMockNodeService(mockCtl)
	api.Node = sNode

	nodeProps := &models.NodeProperties{
		State: models.NodePropertiesState{
			Report: map[string]interface{}{"a": "1"},
			Desire: map[string]interface{}{"b": "2"},
		},
	}
	sNode.EXPECT().GetNodeProperties(gomock.Any(), gomock.Any()).Return(nodeProps, nil).AnyTimes()

	w4 := httptest.NewRecorder()
	req4, _ := http.NewRequest(http.MethodGet, "/v1/nodes/abc/properties", nil)
	router.ServeHTTP(w4, req4)
	assert.Equal(t, http.StatusOK, w4.Code)

	var res models.NodeProperties
	err := json.Unmarshal(w4.Body.Bytes(), &res)
	assert.NoError(t, err)
	expect := models.NodeProperties{
		State: models.NodePropertiesState{
			Report: map[string]interface{}{"a": "1"},
			Desire: map[string]interface{}{"b": "2"},
		},
	}
	assert.Equal(t, expect, res)
}

func TestUpdateNodeProperties(t *testing.T) {
	api, router, mockCtl := initNodeAPI(t)
	defer mockCtl.Finish()

	sNode := ms.NewMockNodeService(mockCtl)
	api.Node = sNode

	nodeProps := &models.NodeProperties{
		State: models.NodePropertiesState{
			Report: map[string]interface{}{"a": "1"},
			Desire: map[string]interface{}{"b": "2"},
		},
	}
	sNode.EXPECT().UpdateNodeProperties(gomock.Any(), gomock.Any(), gomock.Any()).Return(nodeProps, nil).AnyTimes()

	reqNodeProps := &models.NodeProperties{}
	data, err := json.Marshal(reqNodeProps)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/v1/nodes/abc/properties", bytes.NewReader(data))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var res models.NodeProperties
	err = json.Unmarshal(w.Body.Bytes(), &res)
	assert.NoError(t, err)
	expect := models.NodeProperties{
		State: models.NodePropertiesState{
			Report: map[string]interface{}{"a": "1"},
			Desire: map[string]interface{}{"b": "2"},
		},
	}
	assert.Equal(t, expect, res)

	// invalid request params
	reqNodeProps = &models.NodeProperties{
		State: models.NodePropertiesState{
			Desire: map[string]interface{}{"a": 1},
		},
	}
	data, err = json.Marshal(reqNodeProps)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPut, "/v1/nodes/abc/properties", bytes.NewReader(data))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateNodeMode(t *testing.T) {
	api, router, mockCtl := initNodeAPI(t)
	defer mockCtl.Finish()

	sNode := ms.NewMockNodeService(mockCtl)
	api.Node = sNode

	sNode.EXPECT().UpdateNodeMode(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	mode := models.NodeMode{Mode: "cloud"}
	data, err := json.Marshal(mode)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/v1/nodes/abc/mode", bytes.NewReader(data))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// update node mode failed
	sNode.EXPECT().UpdateNodeMode(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("failed to update mode"))
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPut, "/v1/nodes/abc/mode", bytes.NewReader(data))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// invalid request param
	mode = models.NodeMode{Mode: "invalid"}
	data, err = json.Marshal(mode)
	assert.NoError(t, err)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPut, "/v1/nodes/abc/mode", bytes.NewReader(data))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAPI_UpdateCoreApp(t *testing.T) {
	api, router, mockCtl := initNodeAPI(t)
	defer mockCtl.Finish()

	n := "test"
	ns := "default"

	mockNode := ms.NewMockNodeService(mockCtl)
	mockIndex := ms.NewMockIndexService(mockCtl)
	mockApp := ms.NewMockApplicationService(mockCtl)
	mockProp := ms.NewMockPropertyService(mockCtl)
	mockConfig := ms.NewMockConfigService(mockCtl)
	mockInit := ms.NewMockInitService(mockCtl)
	mockModule := ms.NewMockModuleService(mockCtl)
	facade := mf.NewMockFacade(mockCtl)
	api.Init = mockInit
	api.Node = mockNode
	api.Index = mockIndex
	api.App = mockApp
	api.Prop = mockProp
	api.Config = mockConfig
	api.Module = mockModule
	api.Facade = facade

	node := &specV1.Node{
		Namespace: ns,
		Name:      n,
		Version:   "0",
		NodeMode:  context.RunModeKube,
		Attributes: map[string]interface{}{
			specV1.BaetylCoreFrequency: common.DefaultCoreFrequency,
			specV1.BaetylCoreAPIPort:   common.DefaultCoreAPIPort,
			specV1.BaetylAgentPort:     common.DefaultAgentPort,
		},
		Report: map[string]interface{}{"1": "1"},
		Desire: map[string]interface{}{"2": "2"},
	}

	mockNode.EXPECT().Get(nil, ns, n).Return(node, nil).Times(1)

	appList := []string{
		"baetyl-core-1",
		"baetyl-function-2",
	}
	mockIndex.EXPECT().ListAppsByNode(ns, n).Return(appList, nil).Times(1)

	coreApp := &specV1.Application{
		Name:      "baetyl-core-1",
		Type:      "kube",
		Namespace: ns,
		Version:   "0",
		Services: []specV1.Service{
			{
				Name:  "baetyl-core",
				Image: "baetyl-core:v2.0.0",
				Ports: []specV1.ContainerPort{
					{
						HostPort:      30050,
						ContainerPort: 80,
					},
				},
			},
		},
		Volumes: []specV1.Volume{
			{
				Name: "core-conf",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name:    "baetyl-core-conf-ialplsycd",
						Version: "879303",
					},
				},
			},
		},
		System: true,
	}
	mockApp.EXPECT().Get(ns, "baetyl-core-1", "").Return(coreApp, nil).Times(1)
	module := &models.Module{
		Name:    "baetyl-core",
		Version: "v2.0.0",
		Image:   "baetyl-core:v2.0.0",
	}
	mockModule.EXPECT().GetModuleByImage(BaetylModule, "baetyl-core:v2.0.0").Return(module, nil).Times(1)

	// equal case
	coreConfig := models.NodeCoreConfigs{
		Version:   "v2.0.0",
		Frequency: 20,
		APIPort:   30050,
		AgentPort: 30080,
		LogLevel:  "error",
	}
	data, err := json.Marshal(coreConfig)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/v1/nodes/test/core/configs", bytes.NewReader(data))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	mockNode.EXPECT().Get(nil, ns, n).Return(node, nil).Times(1)
	mockIndex.EXPECT().ListAppsByNode(ns, n).Return(appList, nil).Times(1)
	mockApp.EXPECT().Get(ns, "baetyl-core-1", "").Return(coreApp, nil).Times(1)
	mockModule.EXPECT().GetModuleByImage(BaetylModule, "baetyl-core:v2.0.0").Return(module, nil).Times(1)
	mockModule.EXPECT().GetModuleByVersion(BaetylModule, "v2.0.0").Return(module, nil).Times(1)

	cconfig := &specV1.Configuration{
		Name:      "baetyl-core-conf-ialplsycd",
		Namespace: ns,
		Data: map[string]string{
			common.DefaultMasterConfFile: "conf",
		},
	}
	mockConfig.EXPECT().Get(ns, "baetyl-core-conf-ialplsycd", "").Return(cconfig, nil).Times(1)

	pparams := map[string]interface{}{
		"CoreAppName":      "baetyl-core-1",
		"CoreConfName":     "baetyl-core-conf-ialplsycd",
		"CoreFrequency":    "40s",
		"NodeMode":         "kube",
		"AgentPort":        "30080",
		"GPUStats":         true,
		"DiskNetStats":     true,
		"QPSStats":         true,
		BaetylCoreLogLevel: LogLevelDebug,
	}

	confData, err := json.Marshal(cconfig)
	assert.NoError(t, err)
	mockInit.EXPECT().GetResource(ns, node.Name, service.TemplateCoreConfYaml, pparams).Return(confData, nil).Times(1)
	mockConfig.EXPECT().Update(nil, ns, cconfig).Return(cconfig, nil).Times(1)

	mockApp.EXPECT().Update(nil, ns, coreApp).Return(coreApp, nil).Times(1)
	mockNode.EXPECT().UpdateNodeAppVersion(nil, ns, coreApp).Return(appList, nil).Times(1)
	mockNode.EXPECT().Update(ns, node).Return(node, nil).Times(1)

	coreConfig = models.NodeCoreConfigs{
		Version:   "v2.0.0",
		Frequency: 40,
		APIPort:   30000,
		AgentPort: 30080,
	}
	data, err = json.Marshal(coreConfig)
	assert.NoError(t, err)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPut, "/v1/nodes/test/core/configs", bytes.NewReader(data))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	res := specV1.Application{}
	err = json.Unmarshal(w.Body.Bytes(), &res)
	assert.NoError(t, err)
	assert.EqualValues(t, "baetyl-core:v2.0.0", res.Services[0].Image)
	assert.EqualValues(t, int32(30000), res.Services[0].Ports[0].HostPort)

	mockNode.EXPECT().Get(nil, ns, n).Return(node, nil).Times(1)
	mockIndex.EXPECT().ListAppsByNode(ns, n).Return(appList, nil).Times(1)
	mockApp.EXPECT().Get(ns, "baetyl-core-1", "").Return(coreApp, nil).Times(1)
	mockModule.EXPECT().GetModuleByImage(BaetylModule, "baetyl-core:v2.0.0").Return(nil, errors.New("err")).Times(1)

	coreConfig = models.NodeCoreConfigs{
		Version:   "v2.0.0",
		Frequency: 40,
		APIPort:   30000,
		AgentPort: 30080,
	}
	data, err = json.Marshal(coreConfig)
	assert.NoError(t, err)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPut, "/v1/nodes/test/core/configs", bytes.NewReader(data))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	coreConfig = models.NodeCoreConfigs{
		Version:   "v2.0.0",
		Frequency: 30,
		APIPort:   30000,
		AgentPort: 30080,
	}
	data, err = json.Marshal(coreConfig)
	assert.NoError(t, err)

	mockNode.EXPECT().Get(nil, ns, n).Return(node, nil).Times(1)
	mockIndex.EXPECT().ListAppsByNode(ns, n).Return(appList, nil).Times(1)
	mockApp.EXPECT().Get(ns, "baetyl-core-1", "").Return(coreApp, nil).Times(1)
	mockModule.EXPECT().GetModuleByImage(BaetylModule, "baetyl-core:v2.0.0").Return(module, nil).Times(1)
	mockModule.EXPECT().GetModuleByVersion(BaetylModule, "v2.0.0").Return(nil, errors.New("err")).Times(1)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPut, "/v1/nodes/test/core/configs", bytes.NewReader(data))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockNode.EXPECT().Get(nil, ns, n).Return(node, nil).Times(1)
	mockIndex.EXPECT().ListAppsByNode(ns, n).Return(appList, nil).Times(1)
	mockApp.EXPECT().Get(ns, "baetyl-core-1", "").Return(coreApp, nil).Times(1)
	mockModule.EXPECT().GetModuleByImage(BaetylModule, "baetyl-core:v2.0.0").Return(module, nil).Times(1)
	mockModule.EXPECT().GetModuleByVersion(BaetylModule, "v2.0.0").Return(module, nil).Times(1)
	ccconfig := &specV1.Configuration{
		Name:      "baetyl-core-conf-ialplsycd",
		Namespace: ns,
		Data: map[string]string{
			common.DefaultMasterConfFile: "conf",
		},
	}
	mockConfig.EXPECT().Get(ns, "baetyl-core-conf-ialplsycd", "").Return(ccconfig, nil).Times(1)
	mockInit.EXPECT().GetResource(ns, node.Name, service.TemplateCoreConfYaml, pparams).Return(nil, errors.New("err")).Times(1)

	coreConfig = models.NodeCoreConfigs{
		Version:   "v2.0.0",
		Frequency: 40,
		APIPort:   20000,
		AgentPort: 30080,
	}
	data, err = json.Marshal(coreConfig)
	assert.NoError(t, err)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPut, "/v1/nodes/test/core/configs", bytes.NewReader(data))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestAPI_GetCoreAppConfigs(t *testing.T) {
	api, router, mockCtl := initNodeAPI(t)
	defer mockCtl.Finish()

	n := "test"
	ns := "default"

	mockNode := ms.NewMockNodeService(mockCtl)
	mockIndex := ms.NewMockIndexService(mockCtl)
	mockApp := ms.NewMockApplicationService(mockCtl)
	mockProp := ms.NewMockPropertyService(mockCtl)
	mockConfig := ms.NewMockConfigService(mockCtl)
	mockInit := ms.NewMockInitService(mockCtl)
	mockModule := ms.NewMockModuleService(mockCtl)
	api.Init = mockInit
	api.Node = mockNode
	api.Index = mockIndex
	api.App = mockApp
	api.Prop = mockProp
	api.Config = mockConfig
	api.Module = mockModule

	node := &specV1.Node{
		Namespace: ns,
		Name:      n,
		Version:   "0",
		Attributes: map[string]interface{}{
			specV1.BaetylCoreFrequency: common.DefaultCoreFrequency,
			specV1.BaetylCoreAPIPort:   common.DefaultCoreAPIPort,
			specV1.BaetylAgentPort:     common.DefaultAgentPort,
		},
		Report: map[string]interface{}{"1": "1"},
		Desire: map[string]interface{}{"2": "2"},
	}

	mockNode.EXPECT().Get(nil, ns, n).Return(node, nil).Times(1)
	module := &models.Module{
		Name:    "baetyl-core",
		Version: "v2.0.0",
		Image:   "baetyl-core:v2.0.0",
	}
	mockModule.EXPECT().GetModuleByImage(BaetylModule, "baetyl-core:v2.0.0").Return(module, nil).Times(1)

	appList := []string{
		"baetyl-core-1",
		"baetyl-function-2",
	}
	mockIndex.EXPECT().ListAppsByNode(ns, n).Return(appList, nil).Times(1)

	coreApp := &specV1.Application{
		Name:      "baetyl-core-1",
		Type:      "kube",
		Namespace: ns,
		Version:   "0",
		Services: []specV1.Service{
			{
				Name:  "baetyl-core",
				Image: "baetyl-core:v2.0.0",
				Ports: []specV1.ContainerPort{
					{
						HostPort:      30050,
						ContainerPort: 80,
					},
				},
			},
		},
		Volumes: []specV1.Volume{
			{
				Name: "core-conf",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name:    "baetyl-core-conf-ialplsycd",
						Version: "879303",
					},
				},
			},
		},
		System: true,
	}
	mockApp.EXPECT().Get(ns, "baetyl-core-1", "").Return(coreApp, nil).Times(1)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/v1/nodes/test/core/configs", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	mockNode.EXPECT().Get(nil, ns, n).Return(nil, errors.New("err")).Times(1)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/v1/nodes/test/core/configs", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockNode.EXPECT().Get(nil, ns, n).Return(node, nil).Times(1)
	mockIndex.EXPECT().ListAppsByNode(ns, n).Return(nil, errors.New("err")).Times(1)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/v1/nodes/test/core/configs", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestAPI_GetCoreAppVersions(t *testing.T) {
	api, router, mockCtl := initNodeAPI(t)
	defer mockCtl.Finish()

	n := "test"
	ns := "default"

	mockNode := ms.NewMockNodeService(mockCtl)
	mockIndex := ms.NewMockIndexService(mockCtl)
	mockApp := ms.NewMockApplicationService(mockCtl)
	mockProp := ms.NewMockPropertyService(mockCtl)
	mockConfig := ms.NewMockConfigService(mockCtl)
	mockInit := ms.NewMockInitService(mockCtl)
	mockModule := ms.NewMockModuleService(mockCtl)
	api.Init = mockInit
	api.Node = mockNode
	api.Index = mockIndex
	api.App = mockApp
	api.Prop = mockProp
	api.Config = mockConfig
	api.Module = mockModule

	node := &specV1.Node{
		Namespace: ns,
		Name:      n,
		Version:   "0",
		Attributes: map[string]interface{}{
			specV1.BaetylCoreFrequency: common.DefaultCoreFrequency,
		},
		Report: map[string]interface{}{"1": "1"},
		Desire: map[string]interface{}{"2": "2"},
	}

	mockNode.EXPECT().Get(nil, ns, n).Return(node, nil).Times(1)

	appList := []string{
		"baetyl-core-1",
		"baetyl-function-2",
	}
	mockIndex.EXPECT().ListAppsByNode(ns, n).Return(appList, nil).Times(1)

	module := &models.Module{
		Name:    "baetyl-core",
		Version: "v2.0.0",
		Image:   "baetyl-core:v2.0.0",
	}
	mockModule.EXPECT().GetModuleByImage(BaetylModule, "baetyl-core:v2.0.0").Return(module, nil).Times(1)

	coreApp := &specV1.Application{
		Name:      "baetyl-core-1",
		Type:      "kube",
		Namespace: ns,
		Version:   "0",
		Services: []specV1.Service{
			{
				Name:  "baetyl-core",
				Image: "baetyl-core:v2.0.0",
				Ports: []specV1.ContainerPort{
					{
						HostPort:      30050,
						ContainerPort: 80,
					},
				},
			},
		},
		Volumes: []specV1.Volume{
			{
				Name: "core-conf",
				VolumeSource: specV1.VolumeSource{
					Config: &specV1.ObjectReference{
						Name:    "baetyl-core-conf-ialplsycd",
						Version: "879303",
					},
				},
			},
		},
		System: true,
	}
	mockApp.EXPECT().Get(ns, "baetyl-core-1", "").Return(coreApp, nil).Times(1)

	module2 := &models.Module{
		Name:    "v2.1.0",
		Version: "v2.1.0",
		Image:   "baetyl-core:v2.1.0",
	}
	mockModule.EXPECT().GetLatestModule(BaetylModule).Return(module2, nil).Times(1)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/v1/nodes/test/core/versions", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	mockNode.EXPECT().Get(nil, ns, n).Return(nil, errors.New("err")).Times(1)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/v1/nodes/test/core/versions", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockNode.EXPECT().Get(nil, ns, n).Return(node, nil).Times(1)
	mockIndex.EXPECT().ListAppsByNode(ns, n).Return(nil, errors.New("err")).Times(1)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/v1/nodes/test/core/versions", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
