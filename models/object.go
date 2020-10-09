package models

import (
	"io"
	"time"
)

type BucketsView struct {
	Buckets []Bucket `json:"buckets"`
}

type ObjectStorageSourceViewV2 struct {
	Sources map[string]ObjectStorageSourceV2 `json:"sources"`
}

type Bucket struct {
	Name         string    `json:"name,omitempty"`
	CreationDate time.Time `json:"createTime,omitempty"`
}

type ObjectStorageSourceV2 struct {
	AccountEnabled bool `json:"accountEnabled,omitempty"`
}

type ObjectParams struct {
	Marker    string
	MaxKeys   int64
	Prefix    string
	Delimiter string
}

type ObjectSummaryType struct {
	ETag         string
	Key          string
	LastModified time.Time
	Size         int64
	StorageClass string
}

type PrefixType struct {
	Prefix string
}

type ListObjectsResult struct {
	Name           string
	Prefix         string
	Delimiter      string
	Marker         string
	NextMarker     string
	MaxKeys        int64
	IsTruncated    bool
	Contents       []ObjectSummaryType
	CommonPrefixes []PrefixType
}

type ObjectsView struct {
	Objects []ObjectView `json:"objects"`
}

type Object struct {
	ObjectMeta
	Body io.ReadCloser
}

type ObjectView struct {
	Name string `json:"name,omitempty"`
}

type ObjectMeta struct {
	AcceptRanges       string
	CacheControl       string
	ContentDisposition string
	ContentEncoding    string
	ContentLength      int64
	ContentType        string
	ETag               string
	Expires            string
	LastModified       time.Time
	StorageClass       string
}

type ObjectURL struct {
	URL   string
	MD5   string
	Token string
}

type ObjectRequestParams struct {
	Source             string             `json:"source,omitempty"`
	Bucket             string             `json:"bucket,omitempty"`
	Account            string             `form:"account,omitempty"`
	ExternalObjectInfo ExternalObjectInfo `form:",inline"`
}

type ExternalObjectInfo struct {
	Endpoint string `form:"endpoint,omitempty"`
	Ak       string `form:"ak,omitempty"`
	Sk       string `form:"sk,omitempty"`
}
