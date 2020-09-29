package awss3

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

func TestNew(t *testing.T) {
	p, err := New()
	assert.Error(t, err)
	assert.Nil(t, p)
	assert.EqualError(t, err, "open etc/baetyl/cloud.yml: no such file or directory")
}

func TestNotConfigureInternal(t *testing.T) {
	conf := `
minio:
  endpoint: http://106.12.34.129:30900/
  ak: xx
  sk: xx
`
	filename := "cloud.yml"
	err := ioutil.WriteFile(filename, []byte(conf), 0644)
	defer os.Remove(filename)
	common.SetConfFile(filename)

	p, err := New()
	assert.NoError(t, err)
	assert.NotNil(t, p)

	namespace := "default"
	aws3 := p.(plugin.Object)

	bucket := common.RandString(6)
	err = aws3.CreateInternalBucket(namespace, bucket, common.AWSS3PrivatePermission)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "plugin awss3 doesn't support internal object caused it's not configured")
}

func TestAwss3(t *testing.T) {
	t.Skip(t.Name())

	fmt.Println("------------------> Internal <----------------------")
	conf := `
awss3:
 endpoint: http://106.12.34.129:30900
 ak: xx
 sk: xx
`
	filename := "cloud.yml"
	err := ioutil.WriteFile(filename, []byte(conf), 0644)
	defer os.Remove(filename)
	common.SetConfFile(filename)

	p, err := New()
	assert.NoError(t, err)
	assert.NotNil(t, p)

	namespace := "default"
	aws3 := p.(plugin.Object)

	bucket := common.RandString(6)
	err = aws3.CreateInternalBucket(namespace, bucket, common.AWSS3PrivatePermission)
	assert.NoError(t, err)

	err = aws3.HeadInternalBucket(namespace, bucket)
	assert.NoError(t, err)

	err = aws3.HeadInternalBucket(namespace, "xx")
	assert.Contains(t, err.Error(), "The (bucket) resource (xx) is not found")

	buckets, err := aws3.ListInternalBuckets(namespace)
	assert.NoError(t, err)
	assert.NotNil(t, buckets)
	// Length is not predicted of a unknown repo

	err = aws3.PutInternalObject(namespace, bucket, "a", []byte("test"))
	assert.NoError(t, err)

	_, err = aws3.GetInternalObject(namespace, "xx", "a")
	assert.Contains(t, err.Error(), "The (bucket) resource (xx) is not found")

	object, err := aws3.GetInternalObject(namespace, bucket, "a")
	assert.NoError(t, err)
	assert.NotNil(t, object)
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, object.Body)
	assert.NoError(t, err)
	assert.Equal(t, []byte("test"), buf.Bytes())

	_, err = aws3.HeadInternalObject(namespace, "xx", "a")
	assert.Contains(t, err.Error(), "The (bucket) resource (xx) is not found")

	objectMeta, err := aws3.HeadInternalObject(namespace, bucket, "a")
	assert.NoError(t, err)
	assert.Equal(t, int64(4), objectMeta.ContentLength)
	assert.Equal(t, "098f6bcd4621d373cade4e832627b4f6", objectMeta.ETag)

	err = aws3.PutInternalObject(namespace, bucket, "b", []byte("test2"))
	assert.NoError(t, err)

	object, err = aws3.GetInternalObject(namespace, bucket, "b")
	assert.NoError(t, err)
	assert.NotNil(t, object)
	buf = new(bytes.Buffer)
	_, err = io.Copy(buf, object.Body)
	assert.NoError(t, err)
	assert.Equal(t, []byte("test2"), buf.Bytes())

	_, err = aws3.GenInternalObjectURL(namespace, "xx", "b")
	assert.Contains(t, err.Error(), "The (bucket) resource (xx) is not found")

	url, err := aws3.GenInternalObjectURL(namespace, bucket, "b")
	assert.NoError(t, err)
	assert.NotNil(t, url)

	err = aws3.PutInternalObjectFromURL(namespace, bucket, "c", url.URL)
	assert.NoError(t, err)

	object, err = aws3.GetInternalObject(namespace, bucket, "c")
	assert.NoError(t, err)
	assert.NotNil(t, object)
	buf = new(bytes.Buffer)
	_, err = io.Copy(buf, object.Body)
	assert.NoError(t, err)
	assert.Equal(t, []byte("test2"), buf.Bytes())

	_, err = aws3.ListInternalBucketObjects(namespace, "xx", &models.ObjectParams{})
	assert.Contains(t, err.Error(), "The (bucket) resource (xx) is not found")

	objects, err := aws3.ListInternalBucketObjects(namespace, bucket, &models.ObjectParams{})
	assert.NoError(t, err)
	assert.NotNil(t, objects)
	assert.Len(t, objects.Contents, 3)

	objects, err = aws3.ListInternalBucketObjects(namespace, bucket, &models.ObjectParams{
		Prefix: "c",
	})
	assert.NoError(t, err)
	assert.NotNil(t, objects)
	assert.Len(t, objects.Contents, 1)

	objects, err = aws3.ListInternalBucketObjects(namespace, bucket, &models.ObjectParams{
		MaxKeys: 2,
	})
	assert.NoError(t, err)
	assert.NotNil(t, objects)
	assert.Len(t, objects.Contents, 2)

	objects, err = aws3.ListInternalBucketObjects(namespace, bucket, &models.ObjectParams{
		Delimiter: "/",
	})
	assert.NoError(t, err)
	assert.NotNil(t, objects)
	assert.Len(t, objects.Contents, 3)

	fmt.Println("------------------> External <----------------------")

	externalInfo := models.ExternalObjectInfo{
		Endpoint: "http://106.12.34.129:30900/",
		Ak:       "xx",
		Sk:       "xx",
	}

	buckets, err = aws3.ListExternalBuckets(externalInfo)
	assert.NoError(t, err)
	assert.NotNil(t, buckets)

	err = aws3.HeadExternalBucket(externalInfo, "xx")
	assert.Contains(t, err.Error(), "The (bucket) resource (xx) is not found")

	err = aws3.HeadExternalBucket(externalInfo, bucket)
	assert.NoError(t, err)

	_, err = aws3.ListExternalBucketObjects(externalInfo, "xx", &models.ObjectParams{})
	assert.Contains(t, err.Error(), "The (bucket) resource (xx) is not found")

	objects, err = aws3.ListExternalBucketObjects(externalInfo, bucket, &models.ObjectParams{})
	assert.NoError(t, err)
	assert.NotNil(t, objects)
	assert.Len(t, objects.Contents, 3)

	_, err = aws3.GenExternalObjectURL(externalInfo, "xx", "b")
	assert.Contains(t, err.Error(), "The (bucket) resource (xx) is not found")

	url, err = aws3.GenExternalObjectURL(externalInfo, bucket, "b")
	assert.NoError(t, err)
	assert.NotNil(t, url)

	err = aws3.PutInternalObjectFromURL(namespace, "xx", "d", url.URL)
	assert.Contains(t, err.Error(), "The (bucket) resource (xx) is not found")

	err = aws3.PutInternalObjectFromURL(namespace, bucket, "d", url.URL)
	assert.NoError(t, err)

	object, err = aws3.GetInternalObject(namespace, bucket, "d")
	assert.NoError(t, err)
	assert.NotNil(t, object)
	buf = new(bytes.Buffer)
	_, err = io.Copy(buf, object.Body)
	assert.NoError(t, err)
	assert.Equal(t, []byte("test2"), buf.Bytes())

	fmt.Println("------------------> Delete <----------------------")

	err = aws3.DeleteInternalObject(namespace, bucket, "a")
	assert.NoError(t, err)

	err = aws3.DeleteInternalObject(namespace, bucket, "b")
	assert.NoError(t, err)

	err = aws3.DeleteInternalObject(namespace, bucket, "c")
	assert.NoError(t, err)

	err = aws3.DeleteInternalObject(namespace, bucket, "d")
	assert.NoError(t, err)

	objects, err = aws3.ListInternalBucketObjects(namespace, bucket, &models.ObjectParams{})
	assert.NoError(t, err)
	assert.NotNil(t, objects)
	assert.Len(t, objects.Contents, 0)
}

func TestCheckResourceNotFound(t *testing.T) {
	err := errors.New("[Code: NotFound; Message: 404 Not Found; RequestId: xxx]")
	assert.True(t, checkResourceNotFound(err))
	err = errors.New("NotFound: Not Found\n        \t            \t\tstatus code: 404, request id: 1639362962512D92, host id: ")
	assert.True(t, checkResourceNotFound(err))
	err = errors.New("NoSuchBucket: The specified bucket does not exist")
	assert.True(t, checkResourceNotFound(err))
	err = errors.New("RequestError")
	assert.False(t, checkResourceNotFound(err))
}
