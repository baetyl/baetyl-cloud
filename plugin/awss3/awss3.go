package awss3

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
	"github.com/jinzhu/copier"
)

type awss3Storage struct {
	s3Client   *s3.S3
	cfg        CloudConfig
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
}

func init() {
	plugin.RegisterFactory("minio", NewMinio)
}

func NewMinio() (plugin.Plugin, error) {
	var cfg CloudConfig
	if err := common.LoadConfig(&cfg); err != nil {
		return nil, err
	}

	// Configure to use S3 Server
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(cfg.Minio.Ak, cfg.Minio.Sk, ""),
		Endpoint:         aws.String(cfg.Minio.Endpoint),
		Region:           aws.String(cfg.Minio.Region),
		DisableSSL:       aws.Bool(!strings.HasPrefix(cfg.Minio.Endpoint, "https")),
		S3ForcePathStyle: aws.Bool(true),
	}
	newSession, err := session.NewSession(s3Config)
	if err != nil {
		return nil, err
	}
	s3Client := s3.New(newSession)
	uploader := s3manager.NewUploader(newSession)
	downloader := s3manager.NewDownloader(newSession)
	cli := &awss3Storage{
		s3Client:   s3Client,
		uploader:   uploader,
		downloader: downloader,
		cfg:        cfg,
	}
	return cli, nil
}

// CreateBucket CreateBucket
func (c *awss3Storage) CreateBucket(_, bucket, permission string) error {
	//TODO: minio not implement acl completely: NotImplemented: A header you provided implies functionality that is not implemented
	input := &s3.CreateBucketInput{
		ACL:    nil,
		Bucket: &bucket,
	}
	_, err := c.s3Client.CreateBucket(input)
	if err != nil {
		return err
	}
	aclInput := &s3.PutBucketAclInput{
		ACL:    &permission,
		Bucket: &bucket,
	}
	_, err = c.s3Client.PutBucketAcl(aclInput)
	return err
}

// ListBuckets ListBuckets
func (c *awss3Storage) ListBuckets(_ string) ([]models.Bucket, error) {
	buckets, err := c.s3Client.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}
	var res []models.Bucket
	err = copier.Copy(&res, buckets.Buckets)
	return res, err
}

// HeadBucket HeadBucket
func (c *awss3Storage) HeadBucket(_, bucket string) error {
	input := &s3.HeadBucketInput{
		Bucket: &bucket,
	}
	_, err := c.s3Client.HeadBucket(input)
	return err
}

// ListBucketObjects ListBucketObjects
func (c *awss3Storage) ListBucketObjects(_, bucket string, params *models.ObjectParams) (*models.ListObjectsResult, error) {
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
	objectsResult, err := c.s3Client.ListObjects(input)
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

// PutObject PutObject
func (c *awss3Storage) PutObject(_, bucket, name string, b []byte) (err error) {
	_, err = c.s3Client.PutObject(&s3.PutObjectInput{
		Body:   bytes.NewReader(b),
		Bucket: aws.String(bucket),
		Key:    aws.String(name),
	})
	return
}

// PutObjectFromUrl PutObjectFromUrl
func (c *awss3Storage) PutObjectFromURL(_, bucket, name, url string) error {
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
	_, err = c.s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(name),
	})
	return
}

// GenObjectURL GenObjectURL
func (c *awss3Storage) GenObjectURL(_, bucket, name string) (*models.ObjectURL, error) {
	req, _ := c.s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(name),
	})
	url, err := req.Presign(c.cfg.Minio.Expiration)
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
