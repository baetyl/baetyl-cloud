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
	ListSources() map[string]models.ObjectStorageSourceV2

	ListInternalBuckets(userID, source string) ([]models.Bucket, error)
	ListInternalBucketObjects(userID, bucket, source string) (*models.ListObjectsResult, error)
	CreateInternalBucketIfNotExist(userID, bucket, permission, source string) (*models.Bucket, error)
	PutInternalObjectFromURLIfNotExist(userID, bucket, object, url, source string) error
	GenInternalObjectURL(userID string, bucket, object, source string) (*models.ObjectURL, error)

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
func (c *objectService) ListSources() map[string]models.ObjectStorageSourceV2 {
	sources := map[string]models.ObjectStorageSourceV2{}
	for name, object := range c.objects {
		sources[name] = models.ObjectStorageSourceV2{
			AccountEnabled: object.IsAccountEnabled(),
		}
	}
	return sources
}

// ListInternalBuckets ListInternalBuckets
func (c *objectService) ListInternalBuckets(userID, source string) ([]models.Bucket, error) {
	objectPlugin, ok := c.objects[source]
	if !ok {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", fmt.Sprintf("the source (%s) is not supported", source)))
	}
	return objectPlugin.ListInternalBuckets(userID)
}

//ListInternalBucketObjects ListInternalBucketObjects
func (c *objectService) ListInternalBucketObjects(userID, bucket, source string) (*models.ListObjectsResult, error) {
	objectPlugin, ok := c.objects[source]
	if !ok {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", fmt.Sprintf("the source (%s) is not supported", source)))
	}
	return objectPlugin.ListInternalBucketObjects(userID, bucket, &models.ObjectParams{})
}

// CreateInternalBucketIfNotExist CreateInternalBucketIfNotExist
func (c *objectService) CreateInternalBucketIfNotExist(userID, bucket, permission, source string) (*models.Bucket, error) {
	objectPlugin, ok := c.objects[source]
	if !ok {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", fmt.Sprintf("the source (%s) is not supported", source)))
	}
	err := objectPlugin.HeadInternalBucket(userID, bucket)
	if err == nil {
		return &models.Bucket{
			Name: bucket,
		}, nil
	}
	return &models.Bucket{
		Name: bucket,
	}, objectPlugin.CreateInternalBucket(userID, bucket, permission)
}

// PutInternalObjectFromURLIfNotExist PutInternalObjectFromURLIfNotExist
func (c *objectService) PutInternalObjectFromURLIfNotExist(userID, bucket, name, url, source string) error {
	objectPlugin, ok := c.objects[source]
	if !ok {
		return common.Error(common.ErrRequestParamInvalid, common.Field("error", fmt.Sprintf("the source (%s) is not supported", source)))
	}
	if _, err := objectPlugin.HeadInternalObject(userID, bucket, name); err == nil {
		return nil
	}
	return objectPlugin.PutInternalObjectFromURL(userID, bucket, name, url)
}

// GenInternalObjectURL GenInternalObjectURL
func (c *objectService) GenInternalObjectURL(userID string, bucket, object, source string) (*models.ObjectURL, error) {
	objectPlugin, ok := c.objects[source]
	if !ok {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", fmt.Sprintf("the source (%s) is not supported", source)))
	}

	return objectPlugin.GenInternalObjectURL(userID, bucket, object)
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
