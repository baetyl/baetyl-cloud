package service

import (
	"fmt"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

//go:generate mockgen -destination=../mock/service/object.go -package=service github.com/baetyl/baetyl-cloud/v2/service ObjectService

type ObjectService interface {
	ListSources() map[string]models.ObjectStorageSource

	ListBuckets(userID, source string) ([]models.Bucket, error)
	ListBucketObjects(userID, bucket, source string) (*models.ListObjectsResult, error)
	CreateBucketIfNotExist(userID, bucket, permission, source string) (*models.Bucket, error)
	PutObjectFromURLIfNotExist(userID, bucket, object, url, source string) error
	GenObjectURL(userID string, bucket, object, source string) (*models.ObjectURL, error)

	ListExternalBuckets(info models.ExternalObjectInfo, source string) ([]models.Bucket, error)
	ListExternalBucketObjects(info models.ExternalObjectInfo, bucket, source string) (*models.ListObjectsResult, error)
	GenExternalObjectURL(info models.ExternalObjectInfo, bucket, object, source string) (*models.ObjectURL, error)
}

type objectService struct {
	objects map[string]plugin.Object
}

// NewObjectService NewObjectService
func NewObjectService(config *config.CloudConfig) (ObjectService, error) {
	objects := make(map[string]plugin.Object)
	for _, v := range config.Plugin.Objects {
		cs, err := plugin.GetPlugin(v)
		if err != nil {
			return nil, err
		}
		objects[v] = cs.(plugin.Object)
	}
	return &objectService{
		objects: objects,
	}, nil
}

//ListSource ListSource
func (c *objectService) ListSources() map[string]models.ObjectStorageSource {
	sources := map[string]models.ObjectStorageSource{}
	for name, object := range c.objects {
		sources[name] = models.ObjectStorageSource{
			InternalEnabled: object.IsInternalEnabled(),
		}
	}
	return sources
}

// ListBuckets ListBuckets
func (c *objectService) ListBuckets(userID, source string) ([]models.Bucket, error) {
	objectPlugin, ok := c.objects[source]
	if !ok {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", fmt.Sprintf("the source (%s) is not supported", source)))
	}
	return objectPlugin.ListBuckets(userID)
}

//ListBucketObjects ListBucketObjects
func (c *objectService) ListBucketObjects(userID, bucket, source string) (*models.ListObjectsResult, error) {
	objectPlugin, ok := c.objects[source]
	if !ok {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", fmt.Sprintf("the source (%s) is not supported", source)))
	}
	err := objectPlugin.HeadBucket(userID, bucket)
	if err != nil {
		return nil, common.Error(common.ErrResourceNotFound, common.Field("type", "objects"), common.Field("name", bucket))
	}
	return objectPlugin.ListBucketObjects(userID, bucket, &models.ObjectParams{})
}

func (c *objectService) CreateBucketIfNotExist(userID, bucket, permission, source string) (*models.Bucket, error) {
	objectPlugin, ok := c.objects[source]
	if !ok {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", fmt.Sprintf("the source (%s) is not supported", source)))
	}
	err := objectPlugin.HeadBucket(userID, bucket)
	if err == nil {
		return &models.Bucket{
			Name: bucket,
		}, nil
	}
	return &models.Bucket{
		Name: bucket,
	}, objectPlugin.CreateBucket(userID, bucket, permission)
}

//ListSource ListSource
func (c *objectService) PutObjectFromURLIfNotExist(userID, bucket, name, url, source string) error {
	objectPlugin, ok := c.objects[source]
	if !ok {
		return common.Error(common.ErrRequestParamInvalid, common.Field("error", fmt.Sprintf("the source (%s) is not supported", source)))
	}
	if _, err := objectPlugin.HeadObject(userID, bucket, name); err == nil {
		return nil
	}
	return objectPlugin.PutObjectFromURL(userID, bucket, name, url)
}

// GetObjectURL GetObjectURL
func (c *objectService) GenObjectURL(userID string, bucket, object, source string) (*models.ObjectURL, error) {
	objectPlugin, ok := c.objects[source]
	if !ok {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", fmt.Sprintf("the source (%s) is not supported", source)))
	}

	return objectPlugin.GenObjectURL(userID, bucket, object)
}

// ListExternalBuckets ListExternalBuckets
func (c *objectService) ListExternalBuckets(info models.ExternalObjectInfo, source string) ([]models.Bucket, error) {
	objectPlugin, ok := c.objects[source]
	if !ok {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", fmt.Sprintf("the source (%s) is not supported", source)))
	}
	return objectPlugin.ListExternalBuckets(info)
}

// ListExternalBucketObjects ListExternalBucketObjects
func (c *objectService) ListExternalBucketObjects(info models.ExternalObjectInfo, bucket, source string) (*models.ListObjectsResult, error) {
	objectPlugin, ok := c.objects[source]
	if !ok {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", fmt.Sprintf("the source (%s) is not supported", source)))
	}
	err := objectPlugin.HeadExternalBucket(info, bucket)
	if err != nil {
		return nil, common.Error(common.ErrResourceNotFound, common.Field("type", "objects"), common.Field("name", bucket))
	}
	return objectPlugin.ListExternalBucketObjects(info, bucket, &models.ObjectParams{})
}

// GenExternalObjectURL GenExternalObjectURL
func (c *objectService) GenExternalObjectURL(info models.ExternalObjectInfo, bucket, object, source string) (*models.ObjectURL, error) {
	objectPlugin, ok := c.objects[source]
	if !ok {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", fmt.Sprintf("the source (%s) is not supported", source)))
	}
	return objectPlugin.GenExternalObjectURL(info, bucket, object)
}
