package plugin

import (
	"io"

	"github.com/baetyl/baetyl-cloud/models"
)

//go:generate mockgen -destination=../mock/plugin/object.go -package=plugin github.com/baetyl/baetyl-cloud/plugin Object

// Object Object
//TODO: userID doesn't belong to Object, should in the metedata
type Object interface {
	ListBuckets(userID string) ([]models.Bucket, error)
	HeadBucket(userID, bucket string) error
	CreateBucket(userID, bucket, permission string) error
	ListBucketObjects(userID, bucket string, params *models.ObjectParams) (*models.ListObjectsResult, error)
	PutObject(userID, bucket, name string, b []byte) error
	PutObjectFromURL(userID, bucket, name, url string) error
	GetObject(userID, bucket, name string) (*models.Object, error)
	HeadObject(userID, bucket, name string) (*models.ObjectMeta, error)
	DeleteObject(userID, bucket, name string) error
	GenObjectURL(userID, bucket, name string) (*models.ObjectURL, error)
	GenExternalObjectURL(userID string, param models.ConfigObjectItem) (*models.ObjectURL, error)
	io.Closer
}
