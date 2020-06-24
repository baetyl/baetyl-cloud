package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/baetyl/baetyl-cloud/common"
	ms "github.com/baetyl/baetyl-cloud/mock/service"
	"github.com/baetyl/baetyl-cloud/models"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func initFunctionAPI(t *testing.T) (*API, *gin.Engine, *gomock.Controller) {
	api := &API{}
	router := gin.Default()
	mockCtl := gomock.NewController(t)
	mockIM := func(c *gin.Context) {
		common.NewContext(c).SetNamespace("default")
		common.NewContext(c).SetUser(common.User{ID: "default"})
	}
	v1 := router.Group("v1")
	{
		function := v1.Group("/functions")
		function.GET("", mockIM, common.Wrapper(api.ListFunctionSources))
		function.GET("/:source/functions", mockIM, common.Wrapper(api.ListFunctions))
		function.GET("/:source/functions/:name/versions", mockIM, common.Wrapper(api.ListFunctionVersions))
		function.POST("/:source/functions/:name/versions/:version", mockIM, common.Wrapper(api.ImportFunction))
	}
	return api, router, mockCtl
}

func TestListFunctionSources(t *testing.T) {
	api, router, mockCtl := initFunctionAPI(t)
	defer mockCtl.Finish()
	mkFunctionService := ms.NewMockFunctionService(mockCtl)
	api.functionService = mkFunctionService

	sources := []models.FunctionSource{
		{
			Name: "test1",
		},
		{
			Name: "test2",
		},
	}
	// 200
	mkFunctionService.EXPECT().ListSources().Return(sources).Times(1)

	req, _ := http.NewRequest(http.MethodGet, "/v1/functions", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	bytes := w.Body.Bytes()
	var resSource models.FunctionSourceView
	err := json.Unmarshal(bytes, &resSource)
	assert.NoError(t, err)
	assert.Len(t, resSource.Sources, 2)
}

func TestListFunctions(t *testing.T) {
	api, router, mockCtl := initFunctionAPI(t)
	defer mockCtl.Finish()
	mkFunctionService := ms.NewMockFunctionService(mockCtl)
	mkSysConfigService := ms.NewMockSysConfigService(mockCtl)
	api.functionService = mkFunctionService
	api.sysConfigService = mkSysConfigService

	functions := []models.Function{
		{
			Name:    "name1",
			Handler: "handler",
			Version: "version",
			Runtime: "python36",
		},
		{
			Name:    "name2",
			Handler: "handler",
			Version: "version",
			Runtime: "node10",
		},
		{
			Name:    "name",
			Handler: "handler",
			Version: "version",
			Runtime: "other",
		},
	}

	runtimes := []models.SysConfig{
		{
			Key: "python36",
		},
		{
			Key: "node10",
		},
	}

	// 200
	mkFunctionService.EXPECT().List("default", "baiducfc").Return(functions, nil).Times(1)
	mkSysConfigService.EXPECT().ListSysConfigAll(common.BaetylFunctionRuntime).Return(runtimes, nil).Times(1)
	req1, _ := http.NewRequest(http.MethodGet, "/v1/functions/baiducfc/functions", nil)
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)
	var res models.FunctionView
	err := json.Unmarshal(w1.Body.Bytes(), &res)
	assert.NoError(t, err)
	assert.Len(t, res.Functions, 2)
	assert.Equal(t, res.Functions[0].Name, functions[0].Name)
	assert.Equal(t, res.Functions[1].Name, functions[1].Name)

	// 500
	mkFunctionService.EXPECT().List("default", "unknown").Return(nil, errors.New("err")).Times(1)
	req2, _ := http.NewRequest(http.MethodGet, "/v1/functions/unknown/functions", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusInternalServerError, w2.Code)
}

func TestListFunctionVersions(t *testing.T) {
	api, router, mockCtl := initFunctionAPI(t)
	defer mockCtl.Finish()
	mkPluginService := ms.NewMockFunctionService(mockCtl)
	api.functionService = mkPluginService

	functions := []models.Function{
		{
			Name:    "test1",
			Version: "v1",
		},
		{
			Name:    "test1",
			Version: "v2",
		},
	}

	// 200
	mkPluginService.EXPECT().ListFunctionVersions("default", "abc", "baiducfc").Return(functions, nil).Times(1)
	req, _ := http.NewRequest(http.MethodGet, "/v1/functions/baiducfc/functions/abc/versions", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 500
	mkPluginService.EXPECT().ListFunctionVersions("default", "cba", "baiducfc").Return(nil, errors.New("error")).Times(1)
	req, _ = http.NewRequest(http.MethodGet, "/v1/functions/baiducfc/functions/cba/versions", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req)
	assert.Equal(t, http.StatusInternalServerError, w2.Code)
}

func TestImportFunction(t *testing.T) {
	api, router, mockCtl := initFunctionAPI(t)
	defer mockCtl.Finish()
	mkFunctionService := ms.NewMockFunctionService(mockCtl)
	mkObjectService := ms.NewMockObjectService(mockCtl)
	mkSysConfigService := ms.NewMockSysConfigService(mockCtl)
	api.functionService = mkFunctionService
	api.objectService = mkObjectService
	api.sysConfigService = mkSysConfigService

	function := &models.Function{
		Name:    "name1",
		Handler: "handler1",
		Version: "version1",
		Runtime: "runtime1",
		Code: models.FunctionCode{
			Size:     120,
			Sha256:   "nwJRg4SsziinnzTflN8XBilgUzeGIUZS/mxjwnQkzM8=",
			Location: "bj",
		},
	}
	namespace := "default"
	mkFunctionService.EXPECT().GetFunction(namespace, function.Name,
		function.Version, "baiducfc").Return(function, nil).Times(1)

	sysConfig := &models.SysConfig{
		Type:  "object",
		Key:   common.ObjectSource,
		Value: "awss3",
	}
	mkSysConfigService.EXPECT().GetSysConfig(sysConfig.Type, sysConfig.Key).Return(sysConfig, nil).Times(1)

	bucket := &models.Bucket{
		Name: fmt.Sprintf("%s-%s", common.BaetylCloud, namespace),
	}
	mkObjectService.EXPECT().CreateBucketIfNotExist(namespace, bucket.Name, common.AWSS3PrivatePermission, sysConfig.Value).Return(bucket, nil).Times(1)
	mkObjectService.EXPECT().PutObjectFromURLIfNotExist(namespace, bucket.Name, gomock.Any(), function.Code.Location, sysConfig.Value).Return(nil).Times(1)

	// 200
	url := fmt.Sprintf("/v1/functions/baiducfc/functions/%s/versions/%s", function.Name, function.Version)
	req, _ := http.NewRequest(http.MethodPost, url, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	fmt.Println(w.Body.String())
	var res models.ConfigFunctionItem
	err := json.Unmarshal(w.Body.Bytes(), &res)
	assert.NoError(t, err)
	assert.Equal(t, "name1", res.Function)
	assert.Equal(t, "version1", res.Version)
	assert.Equal(t, "runtime1", res.Runtime)
	assert.Equal(t, "handler1", res.Handler)
	assert.Equal(t, "baetyl-cloud-default", res.Bucket)
	assert.Equal(t, "9f02518384acce28a79f34df94df17062960533786214652fe6c63c27424cccf/name1.zip", res.Object)

	mkFunctionService.EXPECT().GetFunction(namespace, function.Name,
		function.Version, "baiducfc").Return(nil, errors.New("err")).Times(1)

	// 500
	url = fmt.Sprintf("/v1/functions/baiducfc/functions/%s/versions/%s", function.Name, function.Version)
	req, _ = http.NewRequest(http.MethodPost, url, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mkFunctionService.EXPECT().GetFunction(namespace, function.Name,
		function.Version, "baiducfc").Return(function, nil).Times(1)
	mkSysConfigService.EXPECT().GetSysConfig(sysConfig.Type, sysConfig.Key).Return(nil, errors.New("err")).Times(1)

	// 500
	url = fmt.Sprintf("/v1/functions/baiducfc/functions/%s/versions/%s", function.Name, function.Version)
	req, _ = http.NewRequest(http.MethodPost, url, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mkFunctionService.EXPECT().GetFunction(namespace, function.Name,
		function.Version, "baiducfc").Return(function, nil).Times(1)
	mkSysConfigService.EXPECT().GetSysConfig(sysConfig.Type, sysConfig.Key).Return(sysConfig, nil).Times(1)
	mkObjectService.EXPECT().CreateBucketIfNotExist(namespace, bucket.Name, common.AWSS3PrivatePermission, sysConfig.Value).Return(nil, errors.New("err")).Times(1)

	// 500
	url = fmt.Sprintf("/v1/functions/baiducfc/functions/%s/versions/%s", function.Name, function.Version)
	req, _ = http.NewRequest(http.MethodPost, url, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mkFunctionService.EXPECT().GetFunction(namespace, function.Name,
		function.Version, "baiducfc").Return(function, nil).Times(1)
	mkSysConfigService.EXPECT().GetSysConfig(sysConfig.Type, sysConfig.Key).Return(sysConfig, nil).Times(1)
	mkObjectService.EXPECT().CreateBucketIfNotExist(namespace, bucket.Name, common.AWSS3PrivatePermission, sysConfig.Value).Return(bucket, nil).Times(1)
	mkObjectService.EXPECT().PutObjectFromURLIfNotExist(namespace, bucket.Name, gomock.Any(), function.Code.Location, sysConfig.Value).Return(errors.New("err")).Times(1)

	// 500
	url = fmt.Sprintf("/v1/functions/baiducfc/functions/%s/versions/%s", function.Name, function.Version)
	req, _ = http.NewRequest(http.MethodPost, url, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
