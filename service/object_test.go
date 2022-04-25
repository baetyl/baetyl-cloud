package service

import (
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/models"
)

func Test_NewObjectService(t *testing.T) {
	conf := &config.CloudConfig{}
	conf.Plugin.Objects = []string{}
	cs, err := NewObjectService(conf)
	assert.NoError(t, err)
	assert.Len(t, cs.(*objectService).objects, 0)
}

func TestObjectService_ListSourcesWithEmptySource(t *testing.T) {
	mockObject := InitEmptyMockEnvironment(t)
	defer mockObject.Close()
	cs, err := NewObjectService(mockObject.conf)
	assert.NoError(t, err)
	res := cs.ListSources()
	assert.NotNil(t, res)
	assert.Equal(t, len(res), 0)
}

func TestObjectService_ListSources(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	mockObject.objectStorage.EXPECT().IsAccountEnabled().Return(true).Times(1)

	cs, err := NewObjectService(mockObject.conf)
	assert.NoError(t, err)
	res := cs.ListSources()
	assert.NotNil(t, res)
	assert.Equal(t, len(res), 1)
}

func TestObjectService_ListInternalBuckets(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	buckets := []models.Bucket{
		{
			Name: "test1",
		},
		{
			Name: "test2",
		},
		{
			Name: "test3",
		},
	}

	namespace := "default"
	mockObject.objectStorage.EXPECT().ListInternalBuckets(namespace).Return(buckets, nil).Times(1)

	cs, err := NewObjectService(mockObject.conf)
	assert.NoError(t, err)
	res, err := cs.ListInternalBuckets(namespace, mockObject.conf.Plugin.Objects[0])
	assert.NoError(t, err)
	assert.Len(t, res, 3)
	assert.Equal(t, res[0].Name, "test1")
	assert.Equal(t, res[1].Name, "test2")
	assert.Equal(t, res[2].Name, "test3")

	namespace2 := "default2"
	mockObject.objectStorage.EXPECT().ListInternalBuckets(namespace2).Return(nil, errors.New("error"))
	_, err2 := cs.ListInternalBuckets(namespace2, mockObject.conf.Plugin.Objects[0])
	assert.Error(t, err2)
}

func TestObjectService_ListInternalObjects(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

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

	namespace := "default"
	bucket := "test1"
	mockObject.objectStorage.EXPECT().ListInternalBucketObjects(namespace, bucket, gomock.Any()).Return(objectsResult, nil).Times(1)
	cs, err := NewObjectService(mockObject.conf)
	assert.NoError(t, err)
	res, err := cs.ListInternalBucketObjects(namespace, bucket, mockObject.conf.Plugin.Objects[0])
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, res, objectsResult)
}

func TestObjectService_GenInternalObjectURL(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	cs, err := NewObjectService(mockObject.conf)
	assert.NoError(t, err)

	urlObj := &models.ObjectURL{
		URL:   "url1",
		MD5:   "md5",
		Token: "token1",
	}
	ns, bucket, name := "ns1", "bucket1", "object1"
	mockObject.objectStorage.EXPECT().GenInternalObjectURL(ns, bucket, name).Return(urlObj, nil).Times(1)

	res, err := cs.GenInternalObjectURL(ns, bucket, name, mockObject.conf.Plugin.Objects[0])
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, res, urlObj)

	_, err = cs.GenInternalObjectURL(ns, bucket, name, "unknown")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "The request parameter is invalid. (the source (unknown) is not supported)")

	mockObject.objectStorage.EXPECT().GenInternalPutObjectURL(ns, bucket, name).Return(urlObj, nil).Times(1)
	res, err = cs.GenInternalObjectPutURL(ns, bucket, name, mockObject.conf.Plugin.Objects[0])
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, res, urlObj)

	_, err = cs.GenInternalObjectPutURL(ns, bucket, name, "unknown")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "The request parameter is invalid. (the source (unknown) is not supported)")
}

func TestObjectService_CreateInternalBucketIfNotExist(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	cs, err := NewObjectService(mockObject.conf)
	assert.NoError(t, err)

	ns, bucket := "ns1", "bucket1"
	mockObject.objectStorage.EXPECT().HeadInternalBucket(ns, bucket).Return(nil).Times(1)

	res, err := cs.CreateInternalBucketIfNotExist(ns, bucket, "private", mockObject.conf.Plugin.Objects[0])
	assert.NoError(t, err)
	assert.NotNil(t, res)

	mockObject.objectStorage.EXPECT().HeadInternalBucket(ns, bucket).Return(errors.New("err")).Times(1)
	mockObject.objectStorage.EXPECT().CreateInternalBucket(ns, bucket, "public").Return(nil).Times(1)
	res, err = cs.CreateInternalBucketIfNotExist(ns, bucket, "public", mockObject.conf.Plugin.Objects[0])
	assert.NoError(t, err)

	_, err = cs.CreateInternalBucketIfNotExist(ns, bucket, "public", "default")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "The request parameter is invalid. (the source (default) is not supported)")
}

func TestObjectService_PutInternalObjectFromURL(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	cs, err := NewObjectService(mockObject.conf)
	assert.NoError(t, err)

	ns, bucket, name, url := "ns1", "bucket1", "object1", "http://test.com/a.zip"
	mockObject.objectStorage.EXPECT().HeadInternalObject(ns, bucket, name).Return(nil, errors.New("err")).Times(1)
	mockObject.objectStorage.EXPECT().PutInternalObjectFromURL(ns, bucket, name, url).Return(nil).Times(1)

	err = cs.PutInternalObjectFromURLIfNotExist(ns, bucket, name, url, mockObject.conf.Plugin.Objects[0])
	assert.NoError(t, err)

	mockObject.objectStorage.EXPECT().HeadInternalObject(ns, bucket, name).Return(nil, nil).Times(1)
	err = cs.PutInternalObjectFromURLIfNotExist(ns, bucket, name, url, mockObject.conf.Plugin.Objects[0])
	assert.NoError(t, err)

	err = cs.PutInternalObjectFromURLIfNotExist(ns, bucket, name, url, "default")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "The request parameter is invalid. (the source (default) is not supported)")
}

func TestObjectService_ListExternalBuckets(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	info := models.ExternalObjectInfo{
		Endpoint: "xxx",
		Ak:       "xxx",
		Sk:       "xxx",
	}

	buckets := []models.Bucket{
		{
			Name: "test1",
		},
		{
			Name: "test2",
		},
		{
			Name: "test3",
		},
	}

	mockObject.objectStorage.EXPECT().ListExternalBuckets(info).Return(buckets, nil).Times(1)

	cs, err := NewObjectService(mockObject.conf)
	assert.NoError(t, err)
	res, err := cs.ListExternalBuckets(info, mockObject.conf.Plugin.Objects[0])
	assert.NoError(t, err)
	assert.Len(t, res, 3)
	assert.Equal(t, res[0].Name, "test1")
	assert.Equal(t, res[1].Name, "test2")
	assert.Equal(t, res[2].Name, "test3")

	mockObject.objectStorage.EXPECT().ListExternalBuckets(info).Return(nil, errors.New("error"))
	_, err2 := cs.ListExternalBuckets(info, mockObject.conf.Plugin.Objects[0])
	assert.Error(t, err2)
}

func TestObjectService_ListExternalObjects(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	info := models.ExternalObjectInfo{
		Endpoint: "xxx",
		Ak:       "xxx",
		Sk:       "xxx",
	}

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

	bucket := "test1"
	mockObject.objectStorage.EXPECT().ListExternalBucketObjects(info, bucket, gomock.Any()).Return(objectsResult, nil).Times(1)
	cs, err := NewObjectService(mockObject.conf)
	assert.NoError(t, err)
	res, err := cs.ListExternalBucketObjects(info, bucket, mockObject.conf.Plugin.Objects[0])
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, res, objectsResult)
}

func TestObjectService_GenExternalObjectURL(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	cs, err := NewObjectService(mockObject.conf)
	assert.NoError(t, err)

	info := models.ExternalObjectInfo{
		Endpoint: "xxx",
		Ak:       "xxx",
		Sk:       "xxx",
	}

	urlObj := &models.ObjectURL{
		URL:   "url1",
		MD5:   "md5",
		Token: "token1",
	}
	bucket, name := "bucket1", "object1"
	mockObject.objectStorage.EXPECT().GenExternalObjectURL(info, bucket, name).Return(urlObj, nil).Times(1)

	res, err := cs.GenExternalObjectURL(info, bucket, name, mockObject.conf.Plugin.Objects[0])
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, res, urlObj)

	_, err = cs.GenExternalObjectURL(info, bucket, name, "unknown")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "The request parameter is invalid. (the source (unknown) is not supported)")
}
