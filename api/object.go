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
	params, err := api.parseObject(c)
	if err != nil {
		return nil, err
	}

	var res []models.Bucket
	if params.Internal {
		res, err = api.Obj.ListBuckets(c.GetUser().ID, params.Source)
	} else {
		res, err = api.Obj.ListExternalBuckets(params.ExternalObjectInfo, params.Source)
	}
	if err != nil {
		return nil, err
	}
	return &models.BucketsView{Buckets: res}, err
}

// ListBucketObjects ListBucketObjects
func (api *API) ListBucketObjects(c *common.Context) (interface{}, error) {
	params, err := api.parseObject(c)
	if err != nil {
		return nil, err
	}

	res := new(models.ListObjectsResult)
	if params.Internal {
		res, err = api.Obj.ListBucketObjects(c.GetUser().ID, params.Bucket, params.Source)
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
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "endpoint field is required when internal is false"))
	}

	if params.ExternalObjectInfo.Ak == "" {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "ak field is required when internal is false"))
	}

	if params.ExternalObjectInfo.Sk == "" {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "sk field is required when internal is false"))
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
