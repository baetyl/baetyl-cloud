package api

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
)

// ListFunctionSources ListFunctionSources
func (api *API) ListFunctionSources(c *common.Context) (interface{}, error) {
	res := api.functionService.ListSources()
	return &models.FunctionSourceView{Sources: res}, nil
}

// ListFunctions list functions
func (api *API) ListFunctions(c *common.Context) (interface{}, error) {
	res, err := api.functionService.List(c.GetUser().ID, c.Param("source"))
	if err != nil {
		return nil, err
	}
	runtimes, err := api.getFunctionRuntimes()
	if err != nil {
		return nil, err
	}
	filter := []models.Function{}
	for _, v := range res {
		for _, runtime := range runtimes {
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
	res, err := api.functionService.ListFunctionVersions(id, n, source)
	if err != nil {
		return nil, err
	}
	return &models.FunctionView{Functions: res}, nil
}

// ImportFunction ImportFunction
func (api *API) ImportFunction(c *common.Context) (interface{}, error) {
	id, name, version, source := c.GetUser().ID, c.Param("name"), c.Param("version"), c.Param("source")

	functionObj, err := api.functionService.GetFunction(id, name, version, source)
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

	_, err = api.objectService.CreateBucketIfNotExist(id, bucketName, common.AWSS3PrivatePermission, objectSource)
	if err != nil {
		return nil, err
	}

	err = api.objectService.PutObjectFromURLIfNotExist(id, bucketName, objectName, functionObj.Code.Location, objectSource)
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

func (api *API) getFunctionRuntimes() ([]string, error) {
	_type := common.BaetylFunctionRuntime
	res, err := api.sysConfigService.ListSysConfigAll(_type)
	if err != nil {
		return nil, err
	}
	var runtimes []string
	for _, v := range res {
		runtimes = append(runtimes, strings.ToLower(v.Key))
	}
	return runtimes, nil
}

func (api *API) getFunctionImageByRuntime(runtime string) (string, error) {
	_type := common.BaetylFunctionRuntime
	res, err := api.sysConfigService.GetSysConfig(_type, runtime)
	if err != nil {
		return "", err
	}
	return res.Value, nil
}

func base64ToHex(s string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
