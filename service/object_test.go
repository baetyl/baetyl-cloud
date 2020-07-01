package service

import (
	"errors"
	"testing"
	"time"

	"github.com/baetyl/baetyl-cloud/config"
	"github.com/baetyl/baetyl-cloud/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_NewObjectService(t *testing.T) {
	conf := &config.CloudConfig{}
	conf.Plugin.Objects = []string{}
	cs, err := NewObjectService(conf)
	assert.NoError(t, err)
	assert.Len(t, cs.(*objectService).objects, 0)
}

func TestObjectService_ListBuckets(t *testing.T) {
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
	mockObject.objectStorage.EXPECT().ListBuckets(namespace).Return(buckets, nil).Times(1)

	cs, err := NewObjectService(mockObject.conf)
	assert.NoError(t, err)
	res, err := cs.ListBuckets(namespace, mockObject.conf.Plugin.Objects[0])
	assert.NoError(t, err)
	assert.Len(t, res, 3)
	assert.Equal(t, res[0].Name, "test1")
	assert.Equal(t, res[1].Name, "test2")
	assert.Equal(t, res[2].Name, "test3")

	namespace2 := "default2"
	mockObject.objectStorage.EXPECT().ListBuckets(namespace2).Return(nil, errors.New("error"))
	_, err2 := cs.ListBuckets(namespace2, mockObject.conf.Plugin.Objects[0])
	assert.Error(t, err2)
}

func TestObjectService_ListObjects(t *testing.T) {
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
	mockObject.objectStorage.EXPECT().HeadBucket(namespace, bucket).Return(nil).Times(1)
	mockObject.objectStorage.EXPECT().ListBucketObjects(namespace, bucket, gomock.Any()).Return(objectsResult, nil).Times(1)
	cs, err := NewObjectService(mockObject.conf)
	assert.NoError(t, err)
	res, err := cs.ListBucketObjects(namespace, bucket, mockObject.conf.Plugin.Objects[0])
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, res, objectsResult)

	namespace2 := "default2"
	bucket2 := "test2"
	mockObject.objectStorage.EXPECT().HeadBucket(namespace2, bucket2).Return(errors.New("error")).Times(1)
	_, err2 := cs.ListBucketObjects(namespace2, bucket2, mockObject.conf.Plugin.Objects[0])
	assert.Error(t, err2)
}

func TestObjectService_ListSources(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	cs, err := NewObjectService(mockObject.conf)
	assert.NoError(t, err)
	res := cs.ListSources()
	assert.NotNil(t, res)
	assert.Equal(t, len(res), 1)
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

func TestObjectService_GenObjectURL(t *testing.T) {
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
	mockObject.objectStorage.EXPECT().GenObjectURL(ns, bucket, name).Return(urlObj, nil).Times(1)

	object := models.ConfigObjectItem{
		Source: mockObject.conf.Plugin.Objects[0],
		Bucket: bucket,
		Object: name,
	}
	res, err := cs.GenObjectURL(ns, object)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, res, urlObj)

	object = models.ConfigObjectItem{
		Source: "unknown",
		Bucket: bucket,
		Object: name,
	}
	_, err = cs.GenObjectURL(ns, object)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "The request parameter is invalid. (the source (unknown) is not supported)")
}

func TestObjectService_CreateBucketIfNotExist(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	cs, err := NewObjectService(mockObject.conf)
	assert.NoError(t, err)

	ns, bucket := "ns1", "bucket1"
	mockObject.objectStorage.EXPECT().HeadBucket(ns, bucket).Return(nil).Times(1)

	res, err := cs.CreateBucketIfNotExist(ns, bucket, "private", mockObject.conf.Plugin.Objects[0])
	assert.NoError(t, err)
	assert.NotNil(t, res)

	mockObject.objectStorage.EXPECT().HeadBucket(ns, bucket).Return(errors.New("err")).Times(1)
	mockObject.objectStorage.EXPECT().CreateBucket(ns, bucket, "public").Return(nil).Times(1)
	res, err = cs.CreateBucketIfNotExist(ns, bucket, "public", mockObject.conf.Plugin.Objects[0])
	assert.NoError(t, err)

	_, err = cs.CreateBucketIfNotExist(ns, bucket, "public", "default")
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "The request parameter is invalid. (the source (default) is not supported)")
}

func TestObjectService_PutObjectFromURL(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	cs, err := NewObjectService(mockObject.conf)
	assert.NoError(t, err)

	ns, bucket, name, url := "ns1", "bucket1", "object1", "http://test.com/a.zip"
	mockObject.objectStorage.EXPECT().HeadObject(ns, bucket, name).Return(nil, errors.New("err")).Times(1)
	mockObject.objectStorage.EXPECT().PutObjectFromURL(ns, bucket, name, url).Return(nil).Times(1)

	err = cs.PutObjectFromURLIfNotExist(ns, bucket, name, url, mockObject.conf.Plugin.Objects[0])
	assert.NoError(t, err)

	mockObject.objectStorage.EXPECT().HeadObject(ns, bucket, name).Return(nil, nil).Times(1)
	err = cs.PutObjectFromURLIfNotExist(ns, bucket, name, url, mockObject.conf.Plugin.Objects[0])
	assert.NoError(t, err)

	err = cs.PutObjectFromURLIfNotExist(ns, bucket, name, url, "default")
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "The request parameter is invalid. (the source (default) is not supported)")
}
