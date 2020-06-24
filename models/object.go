package models

import (
	"io"
	"time"
)

type BucketsView struct {
	Buckets []Bucket `json:"buckets"`
}

type ObjectStorageSourceView struct {
	Sources []ObjectStorageSource `json:"sources"`
}

type Bucket struct {
	Name         string    `json:"name,omitempty"`
	CreationDate time.Time `json:"createTime,omitempty"`
}

type ObjectStorageSource struct {
	Name string `json:"name,omitempty"`
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
