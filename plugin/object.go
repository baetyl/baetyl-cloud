package plugin

import (
	"io"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

//go:generate mockgen -destination=../mock/plugin/object.go -package=plugin github.com/baetyl/baetyl-cloud/v2/plugin Object

// Object Object
//TODO: userID doesn't belong to Object, should in the metedata
type Object interface {
	IsInternalEnabled() bool

	ListInternalBuckets(userID string) ([]models.Bucket, error)
	HeadInternalBucket(userID, bucket string) error
	CreateInternalBucket(userID, bucket, permission string) error
	ListInternalBucketObjects(userID, bucket string, params *models.ObjectParams) (*models.ListObjectsResult, error)
	PutInternalObject(userID, bucket, name string, b []byte) error
	PutInternalObjectFromURL(userID, bucket, name, url string) error
	GetInternalObject(userID, bucket, name string) (*models.Object, error)
	HeadInternalObject(userID, bucket, name string) (*models.ObjectMeta, error)
	DeleteInternalObject(userID, bucket, name string) error
	GenInternalObjectURL(userID, bucket, name string) (*models.ObjectURL, error)

	ListExternalBuckets(info models.ExternalObjectInfo) ([]models.Bucket, error)
	HeadExternalBucket(info models.ExternalObjectInfo, bucket string) error
	ListExternalBucketObjects(info models.ExternalObjectInfo, bucket string, params *models.ObjectParams) (*models.ListObjectsResult, error)
	GenExternalObjectURL(info models.ExternalObjectInfo, bucket, name string) (*models.ObjectURL, error)

	io.Closer
}
