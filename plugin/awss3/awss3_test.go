package awss3

import (
	"bytes"
	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	p, err := NewMinio()
	assert.Error(t, err)
	assert.Nil(t, p)
	assert.EqualError(t, err, "open etc/baetyl/cloud.yml: no such file or directory")
}

func TestNewLocal(t *testing.T) {
	t.Skip(t.Name())
	conf := `
minio:
  endpoint: xx
  ak: xx
  sk: xx
`
	filename := "cloud.yml"
	err := ioutil.WriteFile(filename, []byte(conf), 0644)
	defer os.Remove(filename)
	common.SetConfFile(filename)

	p, err := NewMinio()
	assert.NoError(t, err)
	assert.NotNil(t, p)

	namespace := "default"
	aws3 := p.(plugin.Object)

	err = aws3.PutObjectFromURL(namespace, "baetyl-test", "csd", "http://sdfdsfsdfdsf.bj-bos-sandbox.baidu-int.com/api.json")
	assert.NoError(t, err)

	buckets, err := aws3.ListBuckets(namespace)
	assert.NoError(t, err)
	assert.NotNil(t, buckets)

	err = aws3.CreateBucket(namespace, "test-xxx-"+namespace, common.AWSS3PrivatePermission)
	assert.NoError(t, err)

	err = aws3.PutObject(namespace, "baetyl-test-test", "a", []byte("test"))
	assert.NoError(t, err)

	object, err := aws3.GetObject(namespace, "baetyl-test-test", "a")
	assert.NoError(t, err)
	assert.NotNil(t, object)
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, object.Body)
	assert.NoError(t, err)
	assert.Equal(t, []byte("test"), buf.Bytes())

	objectMeta, err := aws3.HeadObject(namespace, "baetyl-test-test", "a")
	assert.Equal(t, int64(4), objectMeta.ContentLength)
	assert.Equal(t, "098f6bcd4621d373cade4e832627b4f6", objectMeta.ETag)

	url, err := aws3.GenObjectURL(namespace, "baetyl-test-test", "a")
	assert.NoError(t, err)
	assert.NotNil(t, url)
	assert.Equal(t, url.Token, "")

	err = aws3.PutObjectFromURL(namespace, "baetyl-test-test", "b", url.URL)
	assert.NoError(t, err)

	err = aws3.PutObject(namespace, "baetyl-test-test", "c", []byte("ccc"))
	assert.NoError(t, err)

	err = aws3.PutObject(namespace, "baetyl-test-test", "d/d", []byte("ddd"))
	assert.NoError(t, err)

	objects, err := aws3.ListBucketObjects(namespace, "baetyl-test-test", &models.ObjectParams{})
	assert.NoError(t, err)
	assert.NotNil(t, objects)
	assert.Len(t, objects.Contents, 4)

	objects, err = aws3.ListBucketObjects(namespace, "baetyl-test-test", &models.ObjectParams{
		Prefix: "c",
	})
	assert.NoError(t, err)
	assert.NotNil(t, objects)
	assert.Len(t, objects.Contents, 1)

	objects, err = aws3.ListBucketObjects(namespace, "baetyl-test-test", &models.ObjectParams{
		MaxKeys: 2,
	})
	assert.NoError(t, err)
	assert.NotNil(t, objects)
	assert.Len(t, objects.Contents, 2)

	objects, err = aws3.ListBucketObjects(namespace, "baetyl-test-test", &models.ObjectParams{
		Delimiter: "/",
	})
	assert.NoError(t, err)
	assert.NotNil(t, objects)
	assert.Len(t, objects.Contents, 3)
	assert.Len(t, objects.CommonPrefixes, 1)

	objects, err = aws3.ListBucketObjects(namespace, "baetyl-test-test", &models.ObjectParams{
		Marker: "b",
	})
	assert.NoError(t, err)
	assert.NotNil(t, objects)
	assert.Len(t, objects.Contents, 2)

	objects, err = aws3.ListBucketObjects(namespace, "baetyl-test-test", &models.ObjectParams{
		Marker: "c",
	})
	assert.NoError(t, err)
	assert.NotNil(t, objects)
	assert.Len(t, objects.Contents, 1)

	err = aws3.DeleteObject(namespace, "baetyl-test-test", "a")
	assert.NoError(t, err)

	err = aws3.DeleteObject(namespace, "baetyl-test-test", "b")
	assert.NoError(t, err)

	err = aws3.DeleteObject(namespace, "baetyl-test-test", "c")
	assert.NoError(t, err)

	err = aws3.DeleteObject(namespace, "baetyl-test-test", "d/d")
	assert.NoError(t, err)

	objects, err = aws3.ListBucketObjects(namespace, "baetyl-test-test", &models.ObjectParams{})
	assert.NoError(t, err)
	assert.NotNil(t, objects)
	assert.Len(t, objects.Contents, 0)
}
