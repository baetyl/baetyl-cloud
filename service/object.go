package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
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
		return common.Error(common.ErrRequestParamInvalid, common.Field("error", fmt.Sprintf("the source (%s) is not supported", source)))
	}
	if _, err := objectPlugin.HeadObject(userID, bucket, name); err == nil {
		return nil
	}
	return objectPlugin.PutObjectFromURL(userID, bucket, name, url)
}

// GetObjectURL GetObjectURL
func (c *objectService) GenObjectURL(userID string, param models.ConfigObjectItem) (*models.ObjectURL, error) {
	if param.Endpoint != "" {
		return c.GenAwss3ObjectURL(param)
	}

	if _, ok := c.objects[param.Source]; !ok {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", fmt.Sprintf("the source (%s) is not supported", param.Source)))
	}
	return c.objects[param.Source].GenObjectURL(userID, param.Bucket, param.Object)
}

// GetObjectURL GetObjectURL
func (c *objectService) GenAwss3ObjectURL(param models.ConfigObjectItem) (*models.ObjectURL, error) {
	if param.Ak == "" && param.Sk == "" {
		return &models.ObjectURL{
			URL: fmt.Sprintf("%s/%s/%s", param.Endpoint, param.Bucket, param.Object),
		}, nil
	}
	newSession, err := session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(param.Ak, param.Sk, ""),
		Endpoint:         aws.String(param.Endpoint),
		Region:           aws.String("us-east-1"),
		DisableSSL:       aws.Bool(!strings.HasPrefix(param.Endpoint, "https")),
		S3ForcePathStyle: aws.Bool(true),
	})
	if err != nil {
		return nil, err
	}

	cli := s3.New(newSession)
	req, _ := cli.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(param.Bucket),
		Key:    aws.String(param.Object),
	})
	url, err := req.Presign(7 * time.Hour)
	if err != nil {
		return nil, err
	}

	// TODO: Awss3 Etag is not all md5, fix this bug in the future
	return &models.ObjectURL{
		URL: url,
	}, err
}
