package api

import (
	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
)

// ListObjectSources ListObjectSources
func (api *API) ListObjectSources(c *common.Context) (interface{}, error) {
	res := api.Obj.ListSources()
	return &models.ObjectStorageSourceViewV1{Sources: res}, nil
}

// ListBuckets ListBuckets
func (api *API) ListBuckets(c *common.Context) (interface{}, error) {
	res, err := api.Obj.ListInternalBuckets(c.GetUser().ID, c.Param("source"))
	if err != nil {
		return nil, err
	}
	return &models.BucketsView{Buckets: res}, err
}

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

// ListObjectSourcesV2 ListObjectSourcesV2
func (api *API) ListObjectSourcesV2(c *common.Context) (interface{}, error) {
	res := api.Obj.ListSourcesV2()
	return &models.ObjectStorageSourceView{Sources: res}, nil
}

// ListBucketsV2 ListBucketsV2
func (api *API) ListBucketsV2(c *common.Context) (interface{}, error) {
	params, err := api.parseObject(c)
	if err != nil {
		return nil, err
	}

	var res []models.Bucket
	if params.Internal {
		res, err = api.Obj.ListInternalBuckets(c.GetUser().ID, params.Source)
	} else {
		res, err = api.Obj.ListExternalBuckets(params.ExternalObjectInfo, params.Source)
	}
	if err != nil {
		return nil, err
	}
	return &models.BucketsView{Buckets: res}, err
}

// ListBucketObjectsV2 ListBucketObjectsV2
func (api *API) ListBucketObjectsV2(c *common.Context) (interface{}, error) {
	params, err := api.parseObject(c)
	if err != nil {
		return nil, err
	}

	res := new(models.ListObjectsResult)
	if params.Internal {
		res, err = api.Obj.ListInternalBucketObjects(c.GetUser().ID, params.Bucket, params.Source)
	} else {
		res, err = api.Obj.ListExternalBucketObjects(params.ExternalObjectInfo, params.Bucket, params.Source)
	}
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

func (api *API) parseObject(c *common.Context) (*models.ObjectRequestParams, error) {
	params := &models.ObjectRequestParams{}
	if err := c.Bind(params); err != nil {
		return nil, err
	}
	params.Source = c.Param("source")
	params.Bucket = c.Param("bucket")

	if params.Internal {
		return params, nil
	}

	if params.ExternalObjectInfo.Endpoint == "" {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "the parameter 'endpoint' is required for external object"))
	}

	if params.ExternalObjectInfo.Ak == "" {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "the parameter 'ak' is required for external object"))
	}

	if params.ExternalObjectInfo.Sk == "" {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "the parameter 'sk' is required for external object"))
	}

	return params, nil
}

func (api *API) getDefaultObjectSource() (string, error) {
	os, err := api.Prop.GetPropertyValue("object-source")
	if err != nil {
		return "", err
	}
	return os, nil
}
