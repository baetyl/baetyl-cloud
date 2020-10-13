package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/common"
	ms "github.com/baetyl/baetyl-cloud/v2/mock/service"
	"github.com/baetyl/baetyl-cloud/v2/models"
)

func initObjectV2API(t *testing.T) (*API, *gin.Engine, *gomock.Controller) {
	api := &API{}
	router := gin.Default()
	mockCtl := gomock.NewController(t)
	mockIM := func(c *gin.Context) {
		common.NewContext(c).SetNamespace("default")
		common.NewContext(c).SetUser(common.User{ID: "default"})
	}
	v1 := router.Group("v2")
	{
		objects := v1.Group("/objects")
		objects.GET("", mockIM, common.Wrapper(api.ListObjectSourcesV2))
		objects.GET("/:source/buckets", mockIM, common.Wrapper(api.ListBucketsV2))
		objects.GET("/:source/buckets/:bucket/objects", mockIM, common.Wrapper(api.ListBucketObjectsV2))
	}
	return api, router, mockCtl
}

func TestListObjectSourcesV2(t *testing.T) {
	api, router, mockCtl := initObjectV2API(t)
	defer mockCtl.Finish()
	mkObjectService := ms.NewMockObjectService(mockCtl)
	api.Obj = mkObjectService

	sources := map[string]models.ObjectStorageSourceV2{
		"baidubos": {
			AccountEnabled: true,
		},
		"awss3": {
			AccountEnabled: false,
		},
	}
	// 200
	mkObjectService.EXPECT().ListSources().Return(sources).Times(1)
	req0, _ := http.NewRequest(http.MethodGet, "/v2/objects", nil)
	w0 := httptest.NewRecorder()
	router.ServeHTTP(w0, req0)
	assert.Equal(t, http.StatusOK, w0.Code)
	bytes := w0.Body.Bytes()
	var resSource models.ObjectStorageSourceViewV2
	err := json.Unmarshal(bytes, &resSource)
	assert.NoError(t, err)
	assert.Len(t, resSource.Sources, 2)
}

func TestListBucketsV2(t *testing.T) {
	api, router, mockCtl := initObjectV2API(t)
	defer mockCtl.Finish()
	mkObjectService := ms.NewMockObjectService(mockCtl)
	api.Obj = mkObjectService

	buckets := []models.Bucket{
		{
			Name: "test1",
		},
		{
			Name: "test2",
		},
	}
	mkObjectService.EXPECT().ListInternalBuckets("default", "baidubos").Return(buckets, nil).Times(1)
	req, _ := http.NewRequest(http.MethodGet, "/v2/objects/baidubos/buckets?account=current", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	info := models.ExternalObjectInfo{
		Endpoint: "x",
		Ak:       "xx",
		Sk:       "xxx",
	}
	mkObjectService.EXPECT().ListExternalBuckets(info, "baidubos").Return(buckets, nil).Times(1)
	req, _ = http.NewRequest(http.MethodGet, "/v2/objects/baidubos/buckets?account=other&endpoint=x&ak=xx&sk=xxx", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest(http.MethodGet, "/v2/objects/baidubos/buckets?account=other&endpoint=x", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	req, _ = http.NewRequest(http.MethodGet, "/v2/objects/baidubos/buckets?account=other&endpoint=x&ak=xx", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mkObjectService.EXPECT().ListExternalBuckets(info, "test").Return(nil, errors.New("error")).Times(1)
	req2, _ := http.NewRequest(http.MethodGet, "/v2/objects/test/buckets?account=other&endpoint=x&ak=xx&sk=xxx", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusInternalServerError, w2.Code)
}

func TestListBucketObjectsV2(t *testing.T) {
	api, router, mockCtl := initObjectV2API(t)
	defer mockCtl.Finish()
	mkObjectService := ms.NewMockObjectService(mockCtl)
	api.Obj = mkObjectService

	objectsResult := &models.ListObjectsResult{
		Name:        "o1",
		Prefix:      "/",
		Delimiter:   "/",
		Marker:      "a",
		NextMarker:  "b",
		MaxKeys:     10,
		IsTruncated: false,
		Contents: []models.ObjectSummaryType{
			{
				ETag:         "098f6bcd4621d373cade4e832627b4f6",
				Key:          "a.txt",
				LastModified: time.Now(),
				Size:         2,
				StorageClass: "COLD",
			},
			{
				ETag:         "032f6bcd4621d373cade4e832627b4f6",
				Key:          "d/d.txt",
				LastModified: time.Now(),
				Size:         3,
				StorageClass: "STANDARD",
			},
		},
		CommonPrefixes: []models.PrefixType{
			{
				Prefix: "test1",
			},
			{
				Prefix: "test2",
			},
		},
	}
	mkObjectService.EXPECT().ListInternalBucketObjects("default", "baetyl-test", "baidubos").Return(objectsResult, nil).Times(1)

	// 200
	req, _ := http.NewRequest(http.MethodGet, "/v2/objects/baidubos/buckets/baetyl-test/objects?account=current", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest(http.MethodGet, "/v2/objects/baidubos/buckets/baetyl-test/objects", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mkObjectService.EXPECT().ListInternalBucketObjects("default", "unknown", "baidubos").Return(nil, errors.New("error")).Times(1)
	// 404
	req, _ = http.NewRequest(http.MethodGet, "/v2/objects/baidubos/buckets/unknown/objects?account=current", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req)
	assert.Equal(t, http.StatusInternalServerError, w2.Code)

	mkObjectService.EXPECT().ListInternalBucketObjects("default", "unknown2", "baidubos").Return(nil, errors.New("error")).Times(1)
	// 500
	req, _ = http.NewRequest(http.MethodGet, "/v2/objects/baidubos/buckets/unknown2/objects?account=current", nil)
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req)
	assert.Equal(t, http.StatusInternalServerError, w3.Code)

	mkObjectService.EXPECT().ListInternalBucketObjects("default", "unknown3", "baidubos").Return(nil, common.Error(common.ErrResourceNotFound, common.Field("type", "object"))).Times(1)
	// 500
	req, _ = http.NewRequest(http.MethodGet, "/v2/objects/baidubos/buckets/unknown3/objects?account=current", nil)
	w4 := httptest.NewRecorder()
	router.ServeHTTP(w4, req)
	assert.Equal(t, http.StatusNotFound, w4.Code)

	info := models.ExternalObjectInfo{
		Endpoint: "x",
		Ak:       "xx",
		Sk:       "xxx",
	}
	mkObjectService.EXPECT().ListExternalBucketObjects(info, "baetyl-test", "baidubos").Return(objectsResult, nil).Times(1)
	req, _ = http.NewRequest(http.MethodGet, "/v2/objects/baidubos/buckets/baetyl-test/objects?account=other&endpoint=x&ak=xx&sk=xxx", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest(http.MethodGet, "/v2/objects/baidubos/buckets/baetyl-test/objects?account=other&endpoint=x", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	req, _ = http.NewRequest(http.MethodGet, "/v2/objects/baidubos/buckets/baetyl-test/objects?account=other&endpoint=x&ak=xx", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mkObjectService.EXPECT().ListExternalBucketObjects(info, "baetyl-test", "baidubos").Return(nil, errors.New("error")).Times(1)
	req, _ = http.NewRequest(http.MethodGet, "/v2/objects/baidubos/buckets/baetyl-test/objects?account=other&endpoint=x&ak=xx&sk=xxx", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
