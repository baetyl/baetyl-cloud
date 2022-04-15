package api

import (
	"net/url"

	"github.com/baetyl/baetyl-go/v2/errors"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
)

const (
	CurrentAccount = "current"
	OtherAccount   = "other"
	PathStyle      = "pathStyle"
)

// ListObjectSourcesV2 ListObjectSourcesV2
func (api *API) ListObjectSourcesV2(c *common.Context) (interface{}, error) {
	res := api.Obj.ListSources()
	return &models.ObjectStorageSourceViewV2{Sources: res}, nil
}

// ListBucketsV2 ListBucketsV2
func (api *API) ListBucketsV2(c *common.Context) (interface{}, error) {
	params, err := api.parseObject(c)
	if err != nil {
		return nil, err
	}

	var res []models.Bucket
	if params.Account == OtherAccount {
		res, err = api.Obj.ListExternalBuckets(params.ExternalObjectInfo, params.Source)
	} else {
		res, err = api.Obj.ListInternalBuckets(c.GetUser().ID, params.Source)
	}
	if err != nil {
		return nil, err
	}
	return &models.BucketsView{Buckets: res}, nil
}

// ListBucketObjectsV2 ListBucketObjectsV2
func (api *API) ListBucketObjectsV2(c *common.Context) (interface{}, error) {
	params, err := api.parseObject(c)
	if err != nil {
		return nil, err
	}

	res := new(models.ListObjectsResult)
	if params.Account == OtherAccount {
		res, err = api.Obj.ListExternalBucketObjects(params.ExternalObjectInfo, params.Bucket, params.Source)
	} else {
		res, err = api.Obj.ListInternalBucketObjects(c.GetUser().ID, params.Bucket, params.Source)
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

func (api *API) GetObjectPathV2(c *common.Context) (interface{}, error) {
	params, err := api.parseObject(c)
	if err != nil {
		return nil, errors.Trace(err)
	}
	object, err := url.PathUnescape(c.Param("object"))
	if err != nil {
		return nil, errors.Trace(err)
	}

	var res *models.ObjectURL
	if params.Account == OtherAccount {
		res, err = api.Obj.GenExternalObjectURL(params.ExternalObjectInfo, params.Bucket, object, params.Source)
	} else {
		res, err = api.Obj.GenInternalObjectURL(c.GetUser().ID, params.Bucket, object, params.Source)
	}
	if err != nil {
		return nil, errors.Trace(err)
	}
	return res, nil
}

func (api *API) parseObject(c *common.Context) (*models.ObjectRequestParams, error) {
	params := &models.ObjectRequestParams{}
	params.ExternalObjectInfo.AddressFormat = PathStyle
	if err := c.Bind(params); err != nil {
		return nil, err
	}
	params.Source = c.Param("source")
	params.Bucket = c.Param("bucket")

	if params.Account != OtherAccount {
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
