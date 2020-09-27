package api

import (
	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
)

func (api *API) ListObjectSources(c *common.Context) (interface{}, error) {
	res := api.Obj.ListSources()
	return &models.ObjectStorageSourceView{Sources: res}, nil
}

// ListBuckets ListBuckets
func (api *API) ListBuckets(c *common.Context) (interface{}, error) {
	res, err := api.Obj.ListBuckets(c.GetUser().ID, c.Param("source"))
	if err != nil {
		return nil, err
	}
	return &models.BucketsView{Buckets: res}, err
}

// ListBucketObjects ListBucketObjects
func (api *API) ListBucketObjects(c *common.Context) (interface{}, error) {
	id, bucket, source := c.GetUser().ID, c.Param("bucket"), c.Param("source")
	res, err := api.Obj.ListBucketObjects(id, bucket, source)
	if err != nil {
		return nil, err
	}

	var objects []models.ObjectView
	for _, v := range res.Contents {
		view := models.ObjectView{Name: v.Key}
		objects = append(objects, view)
	}
	return &models.ObjectsView{Objects: objects}, err
}

func (api *API) getDefaultObjectSource() (string, error) {
	os, err := api.Prop.GetPropertyValue("object-source")
	if err != nil {
		return "", err
	}
	return os, nil
}
