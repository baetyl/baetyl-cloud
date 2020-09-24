package awss3

import (
	"bytes"
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
	return session.NewSession(s3Config)
}

func (c *awss3Storage) checkInternalSupported() {
	if !c.IsInternalEnabled() {
		panic("plugin awss3 doesn't support internal operating causing it's not configured")
	}
}

func (c *awss3Storage) IsInternalEnabled() bool {
	return c.cfg != nil
}

// CreateBucket CreateBucket
func (c *awss3Storage) CreateBucket(_, bucket, permission string) error {
	c.checkInternalSupported()

	return createBucket(c.s3Client, bucket, permission)
}

func createBucket(cli *s3.S3, bucket, permission string) error {
	//TODO: minio not implement acl completely: NotImplemented: A header you provided implies functionality that is not implemented
	input := &s3.CreateBucketInput{
		ACL:    nil,
		Bucket: &bucket,
	}
	_, err := cli.CreateBucket(input)
	if err != nil {
		return err
	}
	aclInput := &s3.PutBucketAclInput{
		ACL:    &permission,
		Bucket: &bucket,
	}
	_, err = cli.PutBucketAcl(aclInput)
	return err
}

// ListBuckets ListBuckets
func (c *awss3Storage) ListBuckets(_ string) ([]models.Bucket, error) {
	c.checkInternalSupported()

	return listBuckets(c.s3Client)
}

// HeadBucket HeadBucket
func (c *awss3Storage) HeadBucket(_, bucket string) error {
	c.checkInternalSupported()

	return headBucket(c.s3Client, bucket)
}

// ListBucketObjects ListBucketObjects
func (c *awss3Storage) ListBucketObjects(_, bucket string, params *models.ObjectParams) (*models.ListObjectsResult, error) {
	c.checkInternalSupported()

	return listBucketObjects(c.s3Client, bucket, params)
}

// PutObject PutObject
func (c *awss3Storage) PutObject(_, bucket, name string, b []byte) (err error) {
	c.checkInternalSupported()

	_, err = c.s3Client.PutObject(&s3.PutObjectInput{
		Body:   bytes.NewReader(b),
		Bucket: aws.String(bucket),
		Key:    aws.String(name),
	})
	return
}

// PutObjectFromUrl PutObjectFromUrl
func (c *awss3Storage) PutObjectFromURL(_, bucket, name, url string) error {
	c.checkInternalSupported()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	upParams := &s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(bucket),
		Body:   resp.Body,
	}
	_, err = c.uploader.Upload(upParams)
	return err
}

// GetObject GetObject
func (c *awss3Storage) GetObject(_, bucket, name string) (*models.Object, error) {
	c.checkInternalSupported()

	resp, err := c.s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(name),
	})
	if err != nil {
		return nil, err
	}
	var res models.Object
	err = copier.Copy(&res, resp)
	return &res, err
}

// HeadObject HeadObject
func (c *awss3Storage) HeadObject(_, bucket, name string) (*models.ObjectMeta, error) {
	c.checkInternalSupported()

	resp, err := c.s3Client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(name),
	})
	if err != nil {
		return nil, err
	}
	var res models.ObjectMeta
	err = copier.Copy(&res, resp)
	if err != nil {
		return nil, err
	}
	res.ETag, _ = strconv.Unquote(res.ETag)
	return &res, nil
}

// DeleteObject DeleteObject
func (c *awss3Storage) DeleteObject(_, bucket, name string) (err error) {
	c.checkInternalSupported()

	_, err = c.s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(name),
	})
	return
}

// GenObjectURL GenObjectURL
func (c *awss3Storage) GenObjectURL(_, bucket, object string) (*models.ObjectURL, error) {
	c.checkInternalSupported()

	return genObjectURL(c.s3Client, bucket, object, c.cfg.Expiration)
}

func (c *awss3Storage) ListExternalBuckets(info models.ExternalObjectInfo) ([]models.Bucket, error) {
	sessionProvider, err := newS3Session(info.Endpoint, info.Ak, info.Sk, "")
	if err != nil {
		return nil, err
	}

	return listBuckets(s3.New(sessionProvider))
}

func (c *awss3Storage) HeadExternalBucket(info models.ExternalObjectInfo, bucket string) error {
	sessionProvider, err := newS3Session(info.Endpoint, info.Ak, info.Sk, "")
	if err != nil {
		return err
	}

	return headBucket(s3.New(sessionProvider), bucket)
}

func (c *awss3Storage) ListExternalBucketObjects(info models.ExternalObjectInfo, bucket string, params *models.ObjectParams) (*models.ListObjectsResult, error) {
	sessionProvider, err := newS3Session(info.Endpoint, info.Ak, info.Sk, "")
	if err != nil {
		return nil, err
	}

	return listBucketObjects(s3.New(sessionProvider), bucket, params)
}

func (c *awss3Storage) GenExternalObjectURL(info models.ExternalObjectInfo, bucket, object string) (*models.ObjectURL, error) {
	sessionProvider, err := newS3Session(info.Endpoint, info.Ak, info.Sk, "")
	if err != nil {
		return nil, err
	}

	return genObjectURL(s3.New(sessionProvider), bucket, object, time.Hour)
}

func listBuckets(cli *s3.S3) ([]models.Bucket, error) {
	buckets, err := cli.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		return nil, err
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
	return err
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
		return nil, err
	}
	return toObjectList(objectsResult)
}

func toObjectList(objectsResult *s3.ListObjectsOutput) (*models.ListObjectsResult, error) {
	res := new(models.ListObjectsResult)
	err := copier.Copy(res, objectsResult)
	if err != nil {
		panic(fmt.Sprintf("copier exception: %s", err.Error()))
	}
	for i, _ := range res.Contents {
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
		return nil, err
	}
	return &models.ObjectURL{
		URL: url,
	}, err
}

// Close Close
func (c *awss3Storage) Close() error {
	return nil
}
