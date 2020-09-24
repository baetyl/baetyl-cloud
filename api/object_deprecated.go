package api

import (
	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
)

// Deprecated
// ListObjectSources ListObjectSources
func (api *API) ListObjectSources(c *common.Context) (interface{}, error) {
	res2 := api.Obj.ListSources()
	res := []models.ObjectStorageSource{}
	for k := range res2 {
		res = append(res, models.ObjectStorageSource{
			Name: k,
		})
	}
	return &models.ObjectStorageSourceView{Sources: res}, nil
}

// Deprecated
// ListBuckets ListBuckets
func (api *API) ListBuckets(c *common.Context) (interface{}, error) {
	res, err := api.Obj.ListInternalBuckets(c.GetUser().ID, c.Param("source"))
	if err != nil {
		return nil, err
	}
	return &models.BucketsView{Buckets: res}, err
}

// Deprecated
// ListBucketObjects ListBucketObjects
func (api *API) ListBucketObjects(c *common.Context) (interface{}, error) {
	id, bucket, source := c.GetUser().ID, c.Param("bucket"), c.Param("source")
	res, err := api.Obj.ListInternalBucketObjects(id, bucket, source)
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
