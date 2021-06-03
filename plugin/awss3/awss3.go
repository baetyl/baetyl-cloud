package awss3

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/jinzhu/copier"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

type awss3Storage struct {
	s3Client *s3.S3
	cfg      *S3Config
	uploader *s3manager.Uploader
}

func init() {
	plugin.RegisterFactory("awss3", New)
}

func New() (plugin.Plugin, error) {
	var cfg CloudConfig
	if err := common.LoadConfig(&cfg); err != nil {
		return nil, err
	}

	if cfg.AWSS3 == nil {
		return new(awss3Storage), nil
	}

	sessionProvider, err := newS3Session(cfg.AWSS3.Endpoint, cfg.AWSS3.Ak, cfg.AWSS3.Sk, cfg.AWSS3.Region)
	if err != nil {
		return nil, err
	}

	return &awss3Storage{
		s3Client: s3.New(sessionProvider),
		cfg:      cfg.AWSS3,
		uploader: s3manager.NewUploader(sessionProvider),
	}, nil
}

func (c *awss3Storage) IsAccountEnabled() bool {
	return c.cfg != nil
}

// ListInternalBuckets ListInternalBuckets
func (c *awss3Storage) ListInternalBuckets(_ string) ([]models.Bucket, error) {
	err := c.checkInternalSupported()
	if err != nil {
		return nil, err
	}

	return listBuckets(c.s3Client)
}

// HeadInternalBucket HeadInternalBucket
func (c *awss3Storage) HeadInternalBucket(_, bucket string) error {
	err := c.checkInternalSupported()
	if err != nil {
		return err
	}

	return headBucket(c.s3Client, bucket)
}

// CreateInternalBucket CreateInternalBucket
func (c *awss3Storage) CreateInternalBucket(_, bucket, permission string) error {
	err := c.checkInternalSupported()
	if err != nil {
		return err
	}

	//TODO: minio not implement acl completely: NotImplemented: A header you provided implies functionality that is not implemented
	input := &s3.CreateBucketInput{
		ACL:    nil,
		Bucket: &bucket,
	}
	_, err = c.s3Client.CreateBucket(input)
	if err != nil {
		return common.Error(common.ErrObjectOperationException, common.Field("error", err.Error()), common.Field("source", "awss3"))
	}
	aclInput := &s3.PutBucketAclInput{
		ACL:    &permission,
		Bucket: &bucket,
	}
	_, err = c.s3Client.PutBucketAcl(aclInput)
	if err != nil {
		return common.Error(common.ErrObjectOperationException, common.Field("error", err.Error()), common.Field("source", "awss3"))
	}
	return nil
}

// ListInternalBucketObjects ListInternalBucketObjects
func (c *awss3Storage) ListInternalBucketObjects(_, bucket string, params *models.ObjectParams) (*models.ListObjectsResult, error) {
	err := c.checkInternalSupported()
	if err != nil {
		return nil, err
	}

	err = headBucket(c.s3Client, bucket)
	if err != nil {
		return nil, err
	}

	return listBucketObjects(c.s3Client, bucket, params)
}

// PutInternalObject PutInternalObject
func (c *awss3Storage) PutInternalObject(_, bucket, name string, b []byte) error {
	err := c.checkInternalSupported()
	if err != nil {
		return err
	}

	err = headBucket(c.s3Client, bucket)
	if err != nil {
		return err
	}

	_, err = c.s3Client.PutObject(&s3.PutObjectInput{
		Body:   bytes.NewReader(b),
		Bucket: aws.String(bucket),
		Key:    aws.String(name),
	})
	if err != nil {
		return common.Error(common.ErrObjectOperationException, common.Field("error", err.Error()), common.Field("source", "awss3"))
	}
	return nil
}

// PutInternalObjectFromURL PutInternalObjectFromURL
func (c *awss3Storage) PutInternalObjectFromURL(_, bucket, name, url string) error {
	err := c.checkInternalSupported()
	if err != nil {
		return err
	}

	err = headBucket(c.s3Client, bucket)
	if err != nil {
		return err
	}

	resp, err := http.Get(url)
	if err != nil {
		return common.Error(common.ErrObjectOperationException, common.Field("error", err.Error()), common.Field("source", "awss3"))
	}
	defer resp.Body.Close()

	upParams := &s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(name),
		Body:   resp.Body,
	}
	_, err = c.uploader.Upload(upParams)
	if err != nil {
		return common.Error(common.ErrObjectOperationException, common.Field("error", err.Error()), common.Field("source", "awss3"))
	}
	return nil
}

// GetInternalObject GetInternalObject
func (c *awss3Storage) GetInternalObject(_, bucket, name string) (*models.Object, error) {
	err := c.checkInternalSupported()
	if err != nil {
		return nil, err
	}

	err = headBucket(c.s3Client, bucket)
	if err != nil {
		return nil, err
	}

	resp, err := c.s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(name),
	})
	if err != nil {
		if checkResourceNotFound(err) {
			return nil, common.Error(common.ErrResourceNotFound, common.Field("type", "object"), common.Field("name", name))
		}
		return nil, common.Error(common.ErrObjectOperationException, common.Field("error", err.Error()), common.Field("source", "awss3"))
	}
	var res models.Object
	err = copier.Copy(&res, resp)
	return &res, err
}

// HeadInternalObject HeadInternalObject
func (c *awss3Storage) HeadInternalObject(_, bucket, name string) (*models.ObjectMeta, error) {
	err := c.checkInternalSupported()
	if err != nil {
		return nil, err
	}

	err = headBucket(c.s3Client, bucket)
	if err != nil {
		return nil, err
	}

	resp, err := c.s3Client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(name),
	})
	if err != nil {
		if checkResourceNotFound(err) {
			return nil, common.Error(common.ErrResourceNotFound, common.Field("type", "object"), common.Field("name", name))
		}
		return nil, common.Error(common.ErrObjectOperationException, common.Field("error", err.Error()), common.Field("source", "awss3"))
	}
	var res models.ObjectMeta
	err = copier.Copy(&res, resp)
	if err != nil {
		return nil, err
	}
	res.ETag, _ = strconv.Unquote(res.ETag)
	return &res, nil
}

// DeleteInternalObject DeleteInternalObject
func (c *awss3Storage) DeleteInternalObject(_, bucket, name string) (err error) {
	err = c.checkInternalSupported()
	if err != nil {
		return
	}

	_, err = c.s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(name),
	})
	return
}

// GenInternalObjectURL GenInternalObjectURL
func (c *awss3Storage) GenInternalObjectURL(_, bucket, object string) (*models.ObjectURL, error) {
	err := c.checkInternalSupported()
	if err != nil {
		return nil, err
	}

	err = headBucket(c.s3Client, bucket)
	if err != nil {
		return nil, err
	}

	return genObjectURL(c.s3Client, bucket, object, c.cfg.Expiration)
}

// ListExternalBuckets ListExternalBuckets
func (c *awss3Storage) ListExternalBuckets(info models.ExternalObjectInfo) ([]models.Bucket, error) {
	sessionProvider, err := newS3Session(info.Endpoint, info.Ak, info.Sk, "")
	if err != nil {
		return nil, err
	}

	return listBuckets(s3.New(sessionProvider))
}

// HeadExternalBucket HeadExternalBucket
func (c *awss3Storage) HeadExternalBucket(info models.ExternalObjectInfo, bucket string) error {
	sessionProvider, err := newS3Session(info.Endpoint, info.Ak, info.Sk, "")
	if err != nil {
		return err
	}

	return headBucket(s3.New(sessionProvider), bucket)
}

// ListExternalBucketObjects ListExternalBucketObjects
func (c *awss3Storage) ListExternalBucketObjects(info models.ExternalObjectInfo, bucket string, params *models.ObjectParams) (*models.ListObjectsResult, error) {
	sessionProvider, err := newS3Session(info.Endpoint, info.Ak, info.Sk, "")
	if err != nil {
		return nil, err
	}

	cli := s3.New(sessionProvider)

	err = headBucket(cli, bucket)
	if err != nil {
		return nil, err
	}

	return listBucketObjects(cli, bucket, params)
}

// GenExternalObjectURL GenExternalObjectURL
func (c *awss3Storage) GenExternalObjectURL(info models.ExternalObjectInfo, bucket, object string) (*models.ObjectURL, error) {
	sessionProvider, err := newS3Session(info.Endpoint, info.Ak, info.Sk, "")
	if err != nil {
		return nil, err
	}

	cli := s3.New(sessionProvider)

	err = headBucket(cli, bucket)
	if err != nil {
		return nil, err
	}

	return genObjectURL(cli, bucket, object, time.Hour)
}

// Close Close
func (c *awss3Storage) Close() error {
	return nil
}

func (c *awss3Storage) checkInternalSupported() error {
	if !c.IsAccountEnabled() {
		return errors.New("plugin awss3 doesn't support internal object caused it's not configured")
	}
	return nil
}

func newS3Session(endpoint, ak, sk, region string) (*session.Session, error) {
	if region == "" {
		region = "us-east-1"
	}

	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(ak, sk, ""),
		Endpoint:         aws.String(endpoint),
		Region:           aws.String(region),
		DisableSSL:       aws.Bool(!strings.HasPrefix(endpoint, "https")),
		S3ForcePathStyle: aws.Bool(true),
	}
	s, err := session.NewSession(s3Config)
	if err != nil {
		return nil, common.Error(common.ErrObjectOperationException, common.Field("error", err.Error()), common.Field("source", "awss3"))
	}
	return s, nil
}

func listBuckets(cli *s3.S3) ([]models.Bucket, error) {
	buckets, err := cli.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		return nil, common.Error(common.ErrObjectOperationException, common.Field("error", err.Error()), common.Field("source", "awss3"))
	}
	var res []models.Bucket
	err = copier.Copy(&res, buckets.Buckets)
	return res, err
}

func headBucket(cli *s3.S3, bucket string) error {
	input := &s3.HeadBucketInput{
		Bucket: &bucket,
	}
	_, err := cli.HeadBucket(input)
	if err == nil {
		return nil
	}

	if checkResourceNotFound(err) {
		return common.Error(common.ErrResourceNotFound, common.Field("type", "bucket"), common.Field("name", bucket))
	}
	return common.Error(common.ErrObjectOperationException, common.Field("error", err.Error()), common.Field("source", "awss3"))
}

func listBucketObjects(cli *s3.S3, bucket string, params *models.ObjectParams) (*models.ListObjectsResult, error) {
	input := &s3.ListObjectsInput{
		Bucket: &bucket,
	}
	if params.Delimiter != "" {
		input.Delimiter = &params.Delimiter
	}
	if params.Marker != "" {
		input.Marker = &params.Marker
	}
	if params.MaxKeys > 0 {
		input.MaxKeys = &params.MaxKeys
	}
	if params.Prefix != "" {
		input.Prefix = &params.Prefix
	}
	objectsResult, err := cli.ListObjects(input)
	if err != nil {
		return nil, common.Error(common.ErrObjectOperationException, common.Field("error", err.Error()), common.Field("source", "awss3"))
	}
	return toObjectList(objectsResult)
}

func toObjectList(objectsResult *s3.ListObjectsOutput) (*models.ListObjectsResult, error) {
	res := new(models.ListObjectsResult)
	err := copier.Copy(res, objectsResult)
	if err != nil {
		panic(fmt.Sprintf("copier exception: %s", err.Error()))
	}
	for i := range res.Contents {
		res.Contents[i].ETag, _ = strconv.Unquote(res.Contents[i].ETag)
	}
	return res, nil
}

func genObjectURL(cli *s3.S3, bucket, name string, expiration time.Duration) (*models.ObjectURL, error) {
	req, _ := cli.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(name),
	})
	url, err := req.Presign(expiration)
	if err != nil {
		return nil, common.Error(common.ErrObjectOperationException, common.Field("error", err.Error()), common.Field("source", "awss3"))
	}
	return &models.ObjectURL{
		URL: url,
	}, nil
}

func checkResourceNotFound(err error) bool {
	msg := err.Error()
	if strings.Contains(msg, "404") ||
		strings.Contains(strings.ToUpper(msg), "NOTFOUND") ||
		strings.Contains(strings.ToUpper(msg), "NOT FOUND") ||
		strings.Contains(strings.ToUpper(msg), "NOTEXIST") ||
		strings.Contains(strings.ToUpper(msg), "NOT EXIST") {
		return true
	}
	return false
}
