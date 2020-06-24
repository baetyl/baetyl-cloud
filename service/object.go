package service

import (
	"github.com/baetyl/baetyl-cloud/common"
	"github.com/baetyl/baetyl-cloud/config"
	"github.com/baetyl/baetyl-cloud/models"
	"github.com/baetyl/baetyl-cloud/plugin"
)

//go:generate mockgen -destination=../mock/service/object.go -package=plugin github.com/baetyl/baetyl-cloud/service ObjectService

type ObjectService interface {
	ListSources() []models.ObjectStorageSource
	ListBuckets(userID, source string) ([]models.Bucket, error)
	ListBucketObjects(userID, bucket, source string) (*models.ListObjectsResult, error)
	CreateBucketIfNotExist(userID, bucket, permission, source string) (*models.Bucket, error)
	PutObjectFromURLIfNotExist(userID, bucket, name, url, source string) error
	GenObjectURL(userID string, param models.ConfigObjectItem) (*models.ObjectURL, error)
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

// ListBuckets ListBuckets
func (c *objectService) ListBuckets(userID, source string) ([]models.Bucket, error) {
	objectPlugin, ok := c.objects[source]
	if !ok {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "the source is not supported"))
	}
	return objectPlugin.ListBuckets(userID)
}

//ListBucketObjects ListBucketObjects
func (c *objectService) ListBucketObjects(userID, bucket, source string) (*models.ListObjectsResult, error) {
	objectPlugin, ok := c.objects[source]
	if !ok {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "the source is not supported"))
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
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "the source is not supported"))
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
func (c *objectService) ListSources() []models.ObjectStorageSource {
	sources := []models.ObjectStorageSource{}
	for name := range c.objects {
		source := models.ObjectStorageSource{
			Name: name,
		}
		sources = append(sources, source)
	}
	return sources
}

//ListSource ListSource
func (c *objectService) PutObjectFromURLIfNotExist(userID, bucket, name, url, source string) error {
	objectPlugin, ok := c.objects[source]
	if !ok {
		return common.Error(common.ErrRequestParamInvalid, common.Field("error", "the source is not supported"))
	}
	if _, err := objectPlugin.HeadObject(userID, bucket, name); err == nil {
		return nil
	}
	return objectPlugin.PutObjectFromURL(userID, bucket, name, url)
}

// GetObjectURL GetObjectURL
func (c *objectService) GenObjectURL(userID string, param models.ConfigObjectItem) (*models.ObjectURL, error) {
	if _, ok := c.objects[param.Source]; !ok {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "the source is not supported"))
	}
	if param.Endpoint != "" {
		return c.objects[param.Source].GenExternalObjectURL(userID, param)
	}
	return c.objects[param.Source].GenObjectURL(userID, param.Bucket, param.Object)
}
