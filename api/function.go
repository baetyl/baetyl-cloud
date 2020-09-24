package api

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/baetyl/baetyl-go/v2/errors"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
)

// ListFunctionSources ListFunctionSources
func (api *API) ListFunctionSources(c *common.Context) (interface{}, error) {
	runtimes, err := api.Func.ListRuntimes()
	if err != nil {
		return nil, errors.Trace(err)
	}
	res := api.Func.ListSources()
	return &models.FunctionSourceView{Sources: res, Runtimes: runtimes}, nil
}

// ListFunctions list functions
func (api *API) ListFunctions(c *common.Context) (interface{}, error) {
	res, err := api.Func.List(c.GetUser().ID, c.Param("source"))
	if err != nil {
		return nil, err
	}
	runtimes, err := api.Func.ListRuntimes()
	if err != nil {
		return nil, err
	}
	var filter []models.Function
	for _, v := range res {
		for runtime, _ := range runtimes {
			if strings.ToLower(v.Runtime) == runtime {
				filter = append(filter, v)
				break
			}
		}
	}
	return &models.FunctionView{Functions: filter}, nil
}

// ListFunctionVersions list versions of a function
func (api *API) ListFunctionVersions(c *common.Context) (interface{}, error) {
	id, n, source := c.GetUser().ID, c.Param("name"), c.Param("source")
	res, err := api.Func.ListFunctionVersions(id, n, source)
	if err != nil {
		return nil, err
	}
	return &models.FunctionView{Functions: res}, nil
}

// ImportFunction ImportFunction
func (api *API) ImportFunction(c *common.Context) (interface{}, error) {
	id, name, version, source := c.GetUser().ID, c.Param("name"), c.Param("version"), c.Param("source")

	functionObj, err := api.Func.GetFunction(id, name, version, source)
	if err != nil {
		return nil, err
	}

	objectNamePrefix, err := base64ToHex(functionObj.Code.Sha256)
	if err != nil {
		return nil, err
	}
	bucketName := fmt.Sprintf("%s-%s", common.BaetylCloud, id)
	objectName := fmt.Sprintf("%s/%s.%s", objectNamePrefix, functionObj.Name, common.UnpackTypeZip)

	objectSource, err := api.getDefaultObjectSource()
	if err != nil {
		return nil, err
	}

	_, err = api.Obj.CreateInternalBucketIfNotExist(id, bucketName, common.AWSS3PrivatePermission, objectSource)
	if err != nil {
		return nil, err
	}

	err = api.Obj.PutInternalObjectFromURLIfNotExist(id, bucketName, objectName, functionObj.Code.Location, objectSource)
	if err != nil {
		return nil, err
	}

	return &models.ConfigFunctionItem{
		Function: functionObj.Name,
		Version:  functionObj.Version,
		Runtime:  functionObj.Runtime,
		Handler:  functionObj.Handler,
		ConfigObjectItem: models.ConfigObjectItem{
			Source: objectSource,
			Bucket: bucketName,
			Object: objectName,
			Unpack: common.UnpackTypeZip,
		},
	}, nil
}

func base64ToHex(s string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
