package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	mockPlugin "github.com/baetyl/baetyl-cloud/v2/mock/plugin"
	ms "github.com/baetyl/baetyl-cloud/v2/mock/service"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func getMockNode() *specV1.Node {
	mNode := &specV1.Node{
		Namespace: "default",
		Name:      "abc",
		Labels: map[string]string{
			"tag":                "baidu",
			common.LabelNodeName: "abc",
		},
	}
	return mNode
}

func initNodeAPI(t *testing.T) (*API, *gin.Engine, *gomock.Controller) {
	api := &API{}
	router := gin.Default()
	mockCtl := gomock.NewController(t)
	mockIM := func(c *gin.Context) { common.NewContext(c).SetNamespace("default") }
	v1 := router.Group("v1")
	{
		configs := v1.Group("/nodes")
		configs.GET("/:name", mockIM, common.Wrapper(api.GetNode))
		configs.PUT("", mockIM, common.Wrapper(api.GetNodes))
		configs.GET("/:name/stats", mockIM, common.Wrapper(api.GetNodeStats))
		configs.GET("/:name/apps", mockIM, common.Wrapper(api.GetAppByNode))
		configs.PUT("/:name", mockIM, common.Wrapper(api.UpdateNode))
		configs.DELETE("/:name", mockIM, common.Wrapper(api.DeleteNode))
		configs.GET("/:name/init", mockIM, common.Wrapper(api.GenInitCmdFromNode))
		configs.POST("", mockIM, common.Wrapper(api.CreateNode))
		configs.GET("", mockIM, common.Wrapper(api.ListNode))
		configs.GET("/:name/deploys", mockIM, common.Wrapper(api.GetNodeDeployHistory))
	}
	return api, router, mockCtl
}

func TestNewAPI(t *testing.T) {
	c := &config.CloudConfig{}
	c.Plugin.Auth = common.RandString(9)
	c.Plugin.ModelStorage = common.RandString(9)
	c.Plugin.DatabaseStorage = common.RandString(9)
	c.Plugin.PKI = common.RandString(9)
	c.Plugin.Functions = []string{common.RandString(9), common.RandString(9)}
	c.Plugin.Shadow = c.Plugin.ModelStorage
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	mockModelStorage := mockPlugin.NewMockModelStorage(mockCtl)
	plugin.RegisterFactory(c.Plugin.ModelStorage, func() (plugin.Plugin, error) {
		return mockModelStorage, nil
	})
	mockDBStorage := mockPlugin.NewMockDBStorage(mockCtl)
	plugin.RegisterFactory(c.Plugin.DatabaseStorage, func() (plugin.Plugin, error) {
		return mockDBStorage, nil
	})
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
	mkNodeService := ms.NewMockNodeService(mockCtl)
	api.nodeService = mkNodeService

	mNode := getMockNode()

	mkNodeService.EXPECT().Get(mNode.Namespace, mNode.Name).Return(mNode, nil)

	// 200
	req, _ := http.NewRequest(http.MethodGet, "/v1/nodes/abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	bytes := w.Body.Bytes()
	assert.Equal(t, string(bytes), "{\"namespace\":\"default\",\"name\":\"abc\",\"createTime\":\"0001-01-01T00:00:00Z\",\"labels\":{\"baetyl-node-name\":\"abc\",\"tag\":\"baidu\"},\"ready\":false}\n")

	mkNodeService.EXPECT().Get(mNode.Namespace, mNode.Name).Return(nil, common.Error(common.ErrResourceNotFound))
	// 404
	req, _ = http.NewRequest(http.MethodGet, "/v1/nodes/abc", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req)
	assert.Equal(t, http.StatusNotFound, w2.Code)

	mkNodeService.EXPECT().Get(mNode.Namespace, mNode.Name).Return(nil, fmt.Errorf("error"))
	// 500
	req, _ = http.NewRequest(http.MethodGet, "/v1/nodes/abc", nil)
	w2 = httptest.NewRecorder()
	router.ServeHTTP(w2, req)
	assert.Equal(t, http.StatusInternalServerError, w2.Code)
}

func TestGetNodes(t *testing.T) {
	api, router, mockCtl := initNodeAPI(t)
	defer mockCtl.Finish()
	mkNodeService := ms.NewMockNodeService(mockCtl)
	api.nodeService = mkNodeService

	mNode := getMockNode()
	mNode2 := getMockNode()
	mNode2.Name = "abc2"
	mNode2.Labels[common.LabelNodeName] = "abc2"

	// 200
	mkNodeService.EXPECT().Get(mNode.Namespace, mNode.Name).Return(mNode, nil).Times(1)
	mkNodeService.EXPECT().Get(mNode.Namespace, mNode2.Name).Return(mNode2, nil).Times(1)
	nodeNames := &models.NodeNames{
		Names: []string{"abc", "abc2"},
	}
	body, err := json.Marshal(nodeNames)
	assert.NoError(t, err)
	req, _ := http.NewRequest(http.MethodPut, "/v1/nodes?batch", bytes.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, string(w.Body.Bytes()), "[{\"namespace\":\"default\",\"name\":\"abc\",\"createTime\":\"0001-01-01T00:00:00Z\",\"labels\":{\"baetyl-node-name\":\"abc\",\"tag\":\"baidu\"},\"ready\":false},{\"namespace\":\"default\",\"name\":\"abc2\",\"createTime\":\"0001-01-01T00:00:00Z\",\"labels\":{\"baetyl-node-name\":\"abc2\",\"tag\":\"baidu\"},\"ready\":false}]\n")

	// 200 ResourceNotFound
	mkNodeService.EXPECT().Get(mNode.Namespace, mNode.Name).Return(mNode, nil).Times(1)
	mkNodeService.EXPECT().Get(mNode.Namespace, "err_abc").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	mkNodeService.EXPECT().Get(mNode.Namespace, mNode2.Name).Return(mNode2, nil).Times(1)
	nodeNames = &models.NodeNames{
		Names: []string{"abc", "err_abc", "abc2"},
	}
	body, err = json.Marshal(nodeNames)
	assert.NoError(t, err)
	req, _ = http.NewRequest(http.MethodPut, "/v1/nodes?batch", bytes.NewReader(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, string(w.Body.Bytes()), "[{\"namespace\":\"default\",\"name\":\"abc\",\"createTime\":\"0001-01-01T00:00:00Z\",\"labels\":{\"baetyl-node-name\":\"abc\",\"tag\":\"baidu\"},\"ready\":false},{\"namespace\":\"default\",\"name\":\"abc2\",\"createTime\":\"0001-01-01T00:00:00Z\",\"labels\":{\"baetyl-node-name\":\"abc2\",\"tag\":\"baidu\"},\"ready\":false}]\n")

	nodeNames = &models.NodeNames{}
	// 400 validate error
	for i := 0; i < 21; i++ {
		nodeNames.Names = append(nodeNames.Names, "abc")
	}
	body, err = json.Marshal(nodeNames)
	assert.NoError(t, err)
	req, _ = http.NewRequest(http.MethodPut, "/v1/nodes?batch", bytes.NewReader(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	//500
	nodeNames = &models.NodeNames{
		Names: []string{"abc", "abc2"},
	}
	body, err = json.Marshal(nodeNames)
	assert.NoError(t, err)
	mkNodeService.EXPECT().Get(mNode.Namespace, mNode.Name).Return(nil, fmt.Errorf("error")).Times(1)
	req, _ = http.NewRequest(http.MethodPut, "/v1/nodes?batch", bytes.NewReader(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetNodeStats(t *testing.T) {
	api, router, mockCtl := initNodeAPI(t)
	defer mockCtl.Finish()
	mkNodeService := ms.NewMockNodeService(mockCtl)
	api.nodeService = mkNodeService

	mNode := getMockNode()

	mkNodeService.EXPECT().Get(mNode.Namespace, gomock.Any()).Return(mNode, nil)
	// 200
	req, _ := http.NewRequest(http.MethodGet, "/v1/nodes/abc/stats", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	mkNodeService.EXPECT().Get(mNode.Namespace, "cba").Return(nil, common.Error(common.ErrResourceNotFound))
	// 404
	req, _ = http.NewRequest(http.MethodGet, "/v1/nodes/cba/stats", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req)
	assert.Equal(t, http.StatusNotFound, w2.Code)

	mkNodeService.EXPECT().Get(mNode.Namespace, "cba").Return(nil, fmt.Errorf("error"))
	// 500
	req, _ = http.NewRequest(http.MethodGet, "/v1/nodes/cba/stats", nil)
	w2 = httptest.NewRecorder()
	router.ServeHTTP(w2, req)
	assert.Equal(t, http.StatusInternalServerError, w2.Code)

	mNode.Report = map[string]interface{}{
		"nodestats": map[string]interface{}{
			"usage": map[string]string{
				"cpu":    "1",
				"memory": "512Mi",
			},
			"capacity": map[string]string{
				"cpu":    "2",
				"memory": "1024Mi",
			},
		},
		"time": "2020-04-13T10:07:12.267728Z",
	}
	mkNodeService.EXPECT().Get(mNode.Namespace, "abc").Return(mNode, nil)
	// 200
	req, _ = http.NewRequest(http.MethodGet, "/v1/nodes/abc/stats", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	mNode.Report = map[string]interface{}{
		"nodestats": map[string]interface{}{
			"usage": map[string]string{
				"cpu":    "0.5",
				"memory": "512Mi",
			},
			"capacity": map[string]string{
				"cpu":    "2.5",
				"memory": "1024Mi",
			},
		},
	}
	mkNodeService.EXPECT().Get(mNode.Namespace, "abc").Return(mNode, nil)
	// 200
	req, _ = http.NewRequest(http.MethodGet, "/v1/nodes/abc/stats", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	mNode.Report = map[string]interface{}{
		"nodestats": map[string]interface{}{
			"usage": map[string]string{
				"cpu":    "0.5a",
				"memory": "512M",
			},
			"capacity": map[string]string{
				"cpu":    "2.5a",
				"memory": "1024M",
			},
		},
	}
	mkNodeService.EXPECT().Get(mNode.Namespace, "abc").Return(mNode, nil)
	req, _ = http.NewRequest(http.MethodGet, "/v1/nodes/abc/stats", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "quantities must match the regular expression '^([+-]?[0-9.]+)([eEinumkKMGTP]*[-+]?[0-9]*)$'")

	mNode.Report = map[string]interface{}{
		"nodestats": map[string]interface{}{
			"usage": map[string]string{
				"cpu":    "0.5",
				"memory": "512a",
			},
			"capacity": map[string]string{
				"cpu":    "2.5",
				"memory": "1024a",
			},
		},
	}
	mkNodeService.EXPECT().Get(mNode.Namespace, "abc").Return(mNode, nil)
	req, _ = http.NewRequest(http.MethodGet, "/v1/nodes/abc/stats", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "quantities must match the regular expression '^([+-]?[0-9.]+)([eEinumkKMGTP]*[-+]?[0-9]*)$'")

	mNode.Report = map[string]interface{}{
		"nodestats": map[string]interface{}{
			"usage": map[string]string{
				"cpu":    "0.5a",
				"memory": "512a",
			},
			"capacity": map[string]string{
				"cpu":    "2.5",
				"memory": "1024Mi",
			},
		},
	}
	mkNodeService.EXPECT().Get(mNode.Namespace, "abc").Return(mNode, nil)
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

	mkNodeService.EXPECT().Get(mNode.Namespace, "abc").Return(mNode, nil)
	req, _ = http.NewRequest(http.MethodGet, "/v1/nodes/abc/stats", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestListNode(t *testing.T) {
	api, router, mockCtl := initNodeAPI(t)
	defer mockCtl.Finish()
	mkNodeService := ms.NewMockNodeService(mockCtl)
	api.nodeService = mkNodeService

	mClist := &models.NodeList{
		Items: []specV1.Node{
			{
				Name: "node01",
			},
		},
	}

	mkNodeService.EXPECT().List("default", &models.ListOptions{}).Return(mClist, nil)

	// 200
	req, _ := http.NewRequest(http.MethodGet, "/v1/nodes", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	bytes := w.Body.Bytes()
	fmt.Println(string(bytes))
	assert.Equal(t, string(bytes), "{\"total\":0,\"listOptions\":null,\"items\":[{\"name\":\"node01\",\"createTime\":\"0001-01-01T00:00:00Z\",\"ready\":false}]}\n")
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
	mkNodeService := ms.NewMockNodeService(mockCtl)
	api.nodeService = mkNodeService
	mNode := getMockNode()

	// sysapp
	as := ms.NewMockApplicationService(mockCtl)
	cs := ms.NewMockConfigService(mockCtl)
	ss := ms.NewMockSecretService(mockCtl)
	scs := ms.NewMockSysConfigService(mockCtl)
	pki := ms.NewMockPKIService(mockCtl)
	is := ms.NewMockIndexService(mockCtl)
	init := ms.NewMockInitializeService(mockCtl)

	api.applicationService = as
	api.configService = cs
	api.secretService = ss
	api.sysConfigService = scs
	api.pkiService = pki
	api.indexService = is
	api.initService = init

	conf := &specV1.Configuration{
		Name:      "testConf",
		Namespace: mNode.Namespace,
		Data:      nil,
	}
	app := &specV1.Application{
		Name:      "testApp",
		Namespace: mNode.Namespace,
	}
	nodeList := []string{"s0", "s1", "s2"}
	sysConf := &models.SysConfig{
		Type:  "baetyl-edge",
		Key:   "test",
		Value: "123",
	}
	certPEM := &models.PEMCredential{
		CertPEM: []byte("test"),
		KeyPEM:  []byte("test"),
	}
	certMap := map[string][]byte{
		"client.pem": certPEM.CertPEM,
		"client.key": certPEM.KeyPEM,
		"ca.pem":     []byte("test"),
	}
	secret := &specV1.Secret{
		Name:      "sync-" + mNode.Name + "-core-8d4djspg",
		Namespace: mNode.Namespace,
		Data:      certMap,
		Version:   "123",
	}
	cs.EXPECT().Create(mNode.Namespace, gomock.Any()).Return(conf, nil).AnyTimes()
	as.EXPECT().Create(mNode.Namespace, gomock.Any()).Return(app, nil).AnyTimes()
	ss.EXPECT().Get(mNode.Namespace, gomock.Any(), "").Return(secret, nil).AnyTimes()
	scs.EXPECT().GetSysConfig(gomock.Any(), gomock.Any()).Return(sysConf, nil).AnyTimes()
	pki.EXPECT().SignClientCertificate(gomock.Any(), gomock.Any()).Return(certPEM, nil).AnyTimes()
	pki.EXPECT().GetCA().Return([]byte("test"), nil).AnyTimes()
	ss.EXPECT().Create(mNode.Namespace, gomock.Any()).Return(secret, nil).AnyTimes()
	mkNodeService.EXPECT().UpdateNodeAppVersion(mNode.Namespace, gomock.Any()).Return(nodeList, nil).AnyTimes()
	is.EXPECT().RefreshNodesIndexByApp(mNode.Namespace, gomock.Any(), nodeList).AnyTimes()
	init.EXPECT().GetResource(gomock.Any()).Return("{}", nil).AnyTimes()

	mkNodeService.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, nil)
	mkNodeService.EXPECT().Create(mNode.Namespace, mNode).Return(mNode, nil)
	w := httptest.NewRecorder()
	body, _ := json.Marshal(mNode)
	req, _ := http.NewRequest(http.MethodPost, "/v1/nodes", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	mNode.Name = "node-test"
	mNode.Labels[common.LabelNodeName] = mNode.Name

	mkNodeService.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, nil)
	mkNodeService.EXPECT().Create(mNode.Namespace, mNode).Return(mNode, nil)

	w = httptest.NewRecorder()
	body, _ = json.Marshal(mNode)
	req, _ = http.NewRequest(http.MethodPost, "/v1/nodes", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	mkNodeService.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
	mNode.Name = "node-baetyl-test"
	mNode.Labels[common.LabelNodeName] = mNode.Name
	w = httptest.NewRecorder()
	body, _ = json.Marshal(mNode)
	req, _ = http.NewRequest(http.MethodPost, "/v1/nodes", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mkNodeService.EXPECT().Get(gomock.Any(), gomock.Any()).Return(mNode, nil)
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
	mkNodeService := ms.NewMockNodeService(mockCtl)
	api.nodeService = mkNodeService
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
	assert.Contains(t, w.Body.String(), "The request parameter is invalid. (The field (Labels) must contains labels which can be an empty string or a string which is consist of no more than 63 alphanumeric characters, '-', '_' or '.', and must start and end with an alphanumeric character")

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
	assert.Contains(t, w.Body.String(), "The request parameter is invalid. (The field (Labels) must contains labels which can be an empty string or a string which is consist of no more than 63 alphanumeric characters, '-', '_' or '.', and must start and end with an alphanumeric character")

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
	assert.Contains(t, w.Body.String(), "The request parameter is invalid. (The field (Labels) must contains labels which can be an empty string or a string which is consist of no more than 63 alphanumeric characters, '-', '_' or '.', and must start and end with an alphanumeric character")

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
	assert.Contains(t, w.Body.String(), "The request parameter is invalid. (The field (Labels) must contains labels which can be an empty string or a string which is consist of no more than 63 alphanumeric characters, '-', '_' or '.', and must start and end with an alphanumeric character")

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
	assert.Contains(t, w.Body.String(), "The request parameter is invalid. (The field (Labels) must contains labels which can be an empty string or a string which is consist of no more than 63 alphanumeric characters, '-', '_' or '.', and must start and end with an alphanumeric character")
}

func TestUpdateNode(t *testing.T) {
	api, router, mockCtl := initNodeAPI(t)
	defer mockCtl.Finish()
	mkNodeService := ms.NewMockNodeService(mockCtl)
	api.nodeService = mkNodeService

	mApp := getMockNode()

	mkNodeService.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)
	w := httptest.NewRecorder()
	body, _ := json.Marshal(mApp)
	req, _ := http.NewRequest(http.MethodPut, "/v1/nodes/abc", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mApp2 := getMockNode()
	mApp2.Labels["test"] = "test"
	mkNodeService.EXPECT().Get(gomock.Any(), gomock.Any()).Return(mApp2, nil).AnyTimes()
	mkNodeService.EXPECT().Update(mApp.Namespace, mApp).Return(nil, fmt.Errorf("error"))
	w = httptest.NewRecorder()
	body, _ = json.Marshal(mApp)
	req, _ = http.NewRequest(http.MethodPut, "/v1/nodes/abc", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mkNodeService.EXPECT().Get(gomock.Any(), gomock.Any()).Return(mApp, nil).AnyTimes()
	mkNodeService.EXPECT().Update(mApp.Namespace, mApp).Return(mApp2, nil)
	w = httptest.NewRecorder()
	body, _ = json.Marshal(mApp)
	req, _ = http.NewRequest(http.MethodPut, "/v1/nodes/abc", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteNode(t *testing.T) {
	api, router, mockCtl := initNodeAPI(t)
	defer mockCtl.Finish()
	mkNodeService := ms.NewMockNodeService(mockCtl)
	mkApplicationService := ms.NewMockApplicationService(mockCtl)
	mkIndexService := ms.NewMockIndexService(mockCtl)
	mkConfigService := ms.NewMockConfigService(mockCtl)
	mkSecretService := ms.NewMockSecretService(mockCtl)
	mkPkiService := ms.NewMockPKIService(mockCtl)

	api.nodeService = mkNodeService
	api.applicationService = mkApplicationService
	api.indexService = mkIndexService
	api.configService = mkConfigService
	api.secretService = mkSecretService
	api.pkiService = mkPkiService

	mNode := &specV1.Node{
		Namespace: "default",
		Name:      "abc",
		Desire:    genDesireOfSysApps(),
	}
	appCore := &specV1.Application{
		Name:      "core-node12",
		Namespace: mNode.Namespace,
		Version:   "12",
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
			specV1.SecretLabel: specV1.SecretCertificate,
		},
		Annotations: map[string]string{
			common.AnnotationPkiCertID: "certId1",
		},
	}
	secret1f := &specV1.Secret{
		Name:      "secret1f",
		Namespace: mNode.Namespace,
		Labels: map[string]string{
			specV1.SecretLabel: specV1.SecretCertificate,
		},
		Annotations: map[string]string{
			common.AnnotationPkiCertID: "certId1f",
		},
	}

	mkNodeService.EXPECT().Get(gomock.Any(), gomock.Any()).Return(mNode, nil).Times(1)
	mkNodeService.EXPECT().Delete(mNode.Namespace, mNode.Name).Return(nil).Times(1)
	mkApplicationService.EXPECT().Get(mNode.Namespace, appCore.Name, "").Return(appCore, nil).Times(1)
	mkApplicationService.EXPECT().Delete(mNode.Namespace, appCore.Name, appCore.Version).Return(nil).Times(1)
	mkIndexService.EXPECT().RefreshNodesIndexByApp(mNode.Namespace, appCore.Name, gomock.Any()).Return(nil).Times(1)
	mkConfigService.EXPECT().Delete(mNode.Namespace, appCore.Volumes[0].Config.Name).Times(1)
	mkSecretService.EXPECT().Get(mNode.Namespace, appCore.Volumes[1].Secret.Name, "").Return(secret1, nil).Times(1)
	mkPkiService.EXPECT().DeleteClientCertificate("certId1").Return(nil).Times(1)
	mkSecretService.EXPECT().Delete(mNode.Namespace, appCore.Volumes[1].Secret.Name).Times(1)

	mkApplicationService.EXPECT().Get(mNode.Namespace, appFunction.Name, "").Return(appFunction, nil).Times(1)
	mkApplicationService.EXPECT().Delete(mNode.Namespace, appFunction.Name, appFunction.Version).Return(nil).Times(1)
	mkIndexService.EXPECT().RefreshNodesIndexByApp(mNode.Namespace, appFunction.Name, gomock.Any()).Return(nil).Times(1)
	mkConfigService.EXPECT().Delete(mNode.Namespace, appFunction.Volumes[0].Config.Name).Times(1)
	mkSecretService.EXPECT().Get(mNode.Namespace, appFunction.Volumes[1].Secret.Name, "").Return(secret1f, nil).Times(1)
	mkPkiService.EXPECT().DeleteClientCertificate("certId1f").Return(nil).Times(1)
	mkSecretService.EXPECT().Delete(mNode.Namespace, appFunction.Volumes[1].Secret.Name).Times(1)

	// 200
	req, _ := http.NewRequest(http.MethodDelete, "/v1/nodes/abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	mkNodeService.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, errors.New("error")).Times(1)
	req, _ = http.NewRequest(http.MethodDelete, "/v1/nodes/abc", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mkNodeService.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	req, _ = http.NewRequest(http.MethodDelete, "/v1/nodes/abc", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	mkNodeService.EXPECT().Get(gomock.Any(), gomock.Any()).Return(mNode, nil).Times(1)
	mkNodeService.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(errors.New("error"))
	req, _ = http.NewRequest(http.MethodDelete, "/v1/nodes/abc", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeleteNodeError(t *testing.T) {
	api, router, mockCtl := initNodeAPI(t)
	defer mockCtl.Finish()
	mkNodeService := ms.NewMockNodeService(mockCtl)
	mkApplicationService := ms.NewMockApplicationService(mockCtl)
	mkIndexService := ms.NewMockIndexService(mockCtl)
	mkConfigService := ms.NewMockConfigService(mockCtl)
	mkSecretService := ms.NewMockSecretService(mockCtl)
	mkPkiService := ms.NewMockPKIService(mockCtl)

	api.nodeService = mkNodeService
	api.applicationService = mkApplicationService
	api.indexService = mkIndexService
	api.configService = mkConfigService
	api.secretService = mkSecretService
	api.pkiService = mkPkiService

	mNode := &specV1.Node{
		Namespace: "default",
		Name:      "abc",
		Desire:    genDesireOfSysApps(),
	}
	appCore := &specV1.Application{
		Name:      "core-node12",
		Namespace: mNode.Namespace,
		Version:   "12",
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
			specV1.SecretLabel: specV1.SecretCertificate,
		},
		Annotations: map[string]string{
			common.AnnotationPkiCertID: "certId1",
		},
	}

	mkNodeService.EXPECT().Get(gomock.Any(), gomock.Any()).Return(mNode, nil).Times(1)
	mkNodeService.EXPECT().Delete(mNode.Namespace, mNode.Name).Return(nil).Times(1)
	mkApplicationService.EXPECT().Get(mNode.Namespace, appCore.Name, "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	mkApplicationService.EXPECT().Get(mNode.Namespace, appFunction.Name, "").Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)

	// 200
	req, _ := http.NewRequest(http.MethodDelete, "/v1/nodes/abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	mkNodeService.EXPECT().Get(gomock.Any(), gomock.Any()).Return(mNode, nil).Times(1)
	mkNodeService.EXPECT().Delete(mNode.Namespace, mNode.Name).Return(nil).Times(1)
	mkApplicationService.EXPECT().Get(mNode.Namespace, appCore.Name, "").Return(appCore, nil).Times(1)
	mkApplicationService.EXPECT().Delete(mNode.Namespace, appCore.Name, appCore.Version).Return(errors.New("error")).Times(1)
	mkIndexService.EXPECT().RefreshNodesIndexByApp(mNode.Namespace, appCore.Name, gomock.Any()).Return(errors.New("error")).Times(1)
	mkConfigService.EXPECT().Delete(mNode.Namespace, appCore.Volumes[0].Config.Name).Return(errors.New("error")).Times(1)
	mkSecretService.EXPECT().Get(mNode.Namespace, appCore.Volumes[1].Secret.Name, "").Return(secret1, nil).Times(1)
	mkPkiService.EXPECT().DeleteClientCertificate("certId1").Return(errors.New("error")).Times(1)
	mkSecretService.EXPECT().Delete(mNode.Namespace, appCore.Volumes[1].Secret.Name).Times(1)

	mkApplicationService.EXPECT().Get(mNode.Namespace, appFunction.Name, "").Return(appFunction, nil).Times(1)
	mkApplicationService.EXPECT().Delete(mNode.Namespace, appFunction.Name, appFunction.Version).Return(errors.New("error")).Times(1)
	mkIndexService.EXPECT().RefreshNodesIndexByApp(mNode.Namespace, appFunction.Name, gomock.Any()).Return(errors.New("error")).Times(1)
	mkConfigService.EXPECT().Delete(mNode.Namespace, appFunction.Volumes[0].Config.Name).Return(errors.New("error")).Times(1)
	mkSecretService.EXPECT().Get(mNode.Namespace, appFunction.Volumes[1].Secret.Name, "").Return(nil, errors.New("error")).Times(1)

	// 200
	req2, _ := http.NewRequest(http.MethodDelete, "/v1/nodes/abc", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)

	mkNodeService.EXPECT().Get(gomock.Any(), gomock.Any()).Return(mNode, nil).Times(1)
	mkNodeService.EXPECT().Delete(mNode.Namespace, mNode.Name).Return(nil).Times(1)
	mkApplicationService.EXPECT().Get(mNode.Namespace, appCore.Name, "").Return(nil, errors.New("error")).Times(1)
	mkApplicationService.EXPECT().Delete(mNode.Namespace, appCore.Name, appCore.Version).Return(errors.New("error")).Times(1)
	mkIndexService.EXPECT().RefreshNodesIndexByApp(mNode.Namespace, appCore.Name, gomock.Any()).Return(errors.New("error")).Times(1)

	mkApplicationService.EXPECT().Get(mNode.Namespace, appFunction.Name, "").Return(appFunction, nil).Times(1)
	mkApplicationService.EXPECT().Delete(mNode.Namespace, appFunction.Name, appFunction.Version).Return(errors.New("error")).Times(1)
	mkIndexService.EXPECT().RefreshNodesIndexByApp(mNode.Namespace, appFunction.Name, gomock.Any()).Return(errors.New("error")).Times(1)
	mkConfigService.EXPECT().Delete(mNode.Namespace, appFunction.Volumes[0].Config.Name).Return(errors.New("error")).Times(1)
	mkSecretService.EXPECT().Get(mNode.Namespace, appFunction.Volumes[1].Secret.Name, "").Return(nil, errors.New("error")).Times(1)

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
	mkNodeService := ms.NewMockNodeService(mockCtl)
	api.nodeService = mkNodeService

	req, _ := http.NewRequest(http.MethodGet, "/v1/nodes/abc/deploys", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestGenInitCmdFromNode(t *testing.T) {
	api, router, mockCtl := initNodeAPI(t)
	defer mockCtl.Finish()
	nMock := ms.NewMockNodeService(mockCtl)
	scs := ms.NewMockSysConfigService(mockCtl)
	auth := ms.NewMockAuthService(mockCtl)
	api.nodeService = nMock
	api.sysConfigService = scs
	api.authService = auth

	node := getMockNode()

	sc := &models.SysConfig{Key: common.AddressActive, Type: "address", Value: "baetyl.com"}

	scs.EXPECT().GetSysConfig(sc.Type, sc.Key).Return(sc, nil).Times(1)
	auth.EXPECT().GenToken(gomock.Any()).Return("token", nil).Times(1)
	nMock.EXPECT().Get(node.Namespace, node.Name).Return(node, nil).Times(1)

	req, _ := http.NewRequest(http.MethodGet, "/v1/nodes/abc/init", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGenInitCmdFromNode_ErrNode(t *testing.T) {
	api, router, mockCtl := initNodeAPI(t)
	defer mockCtl.Finish()
	nMock := ms.NewMockNodeService(mockCtl)
	api.nodeService = nMock

	node := getMockNode()

	nMock.EXPECT().Get(node.Namespace, node.Name).Return(nil,
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
	mkNodeService, mkIndexService := ms.NewMockNodeService(mockCtl), ms.NewMockIndexService(mockCtl)
	mkAppService := ms.NewMockApplicationService(mockCtl)
	api.applicationService = mkAppService
	api.nodeService, api.indexService = mkNodeService, mkIndexService

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

	mkNodeService.EXPECT().Get(gomock.Any(), gomock.Any()).Return(node, nil).AnyTimes()

	w4 := httptest.NewRecorder()
	req4, _ := http.NewRequest(http.MethodGet, "/v1/nodes/abc/apps", nil)
	router.ServeHTTP(w4, req4)
	assert.Equal(t, http.StatusOK, w4.Code)
	list := &models.ApplicationList{}
	json.Unmarshal(w4.Body.Bytes(), list)
	assert.Equal(t, 0, list.Total)

	node.Desire.SetAppInfos(true, sysappinfos)
	node.Desire.SetAppInfos(false, appinfos)

	mkAppService.EXPECT().Get(node.Namespace, appNames[0], "").Return(apps[0], nil).AnyTimes()
	mkAppService.EXPECT().Get(node.Namespace, appNames[1], "").Return(apps[1], nil).AnyTimes()
	mkAppService.EXPECT().Get(node.Namespace, appNames[2], "").Return(nil, common.Error(common.ErrResourceNotFound)).AnyTimes()
	mkAppService.EXPECT().Get(node.Namespace, sysAppNames[0], "").Return(sysapps[0], nil).AnyTimes()
	mkAppService.EXPECT().Get(node.Namespace, sysAppNames[1], "").Return(sysapps[1], nil).AnyTimes()
	mkAppService.EXPECT().Get(node.Namespace, sysAppNames[2], "").Return(nil, common.Error(common.ErrResourceNotFound)).AnyTimes()

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
	mkNodeService := ms.NewMockNodeService(mockCtl)
	api.nodeService = mkNodeService
	namespace := "iotCoreId"

	mkNodeService.EXPECT().List(namespace, gomock.Any()).Return(nil, errors.New("error"))
	_, err := api.NodeNumberCollector(namespace)
	assert.Error(t, err)

	list := &models.NodeList{
		Items: []specV1.Node{
			{
				Namespace: namespace,
				Name:      "test1",
			},
		},
	}
	mkNodeService.EXPECT().List(namespace, gomock.Any()).Return(list, nil)
	res, err := api.NodeNumberCollector(namespace)
	assert.NoError(t, err)
	assert.Equal(t, 1, res[plugin.QuotaNode])
}
