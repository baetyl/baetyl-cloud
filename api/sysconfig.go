package api

import (
	"bytes"
	"text/template"
	"github.com/baetyl/baetyl-cloud/common"
	"github.com/baetyl/baetyl-cloud/models"
	"strconv"
	"time"
)

// TODO: optimize this layer, general abstraction

// GetSysConfig get a system config
func (api *API) GetSysConfig(c *common.Context) (interface{}, error) {
	tp, key := c.Param("type"), c.Param("key")
	return api.sysConfigService.GetSysConfig(tp, key)
}


func (api *API) ParseTemplate(key string, data map[string]string) ([]byte, error) {
	tl, err := api.initService.GetResource(key)
	if err != nil {
		return nil, err
	}
	t, err := template.New(key).Parse(tl)
	if err != nil {
		return nil, err
	}
	buf := &bytes.Buffer{}
	err = t.Execute(buf, data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (api *API) ListSysConfigAll(c *common.Context) (interface{}, error) {
	res, err := api.sysConfigService.ListSysConfigAll(c.Param("type"))
	if err != nil {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", err.Error()))
	}
	return &models.SysConfigView{
		SysConfigs: res,
	}, nil
}


func (api *API) ListSysConfig(c *common.Context) (interface{}, error) {
	pageNo, _ :=  strconv.Atoi(c.Query("pageNo"))
	pageSize, _ :=  strconv.Atoi(c.Query("pageSize"))
	res, err := api.sysConfigService.ListSysConfig(c.Query("type"), pageNo, pageSize)
	if err != nil {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", err.Error()))
	}
	return &models.SysConfigView{
		SysConfigs: res,
	}, nil
}

//// CreateSysConfig create a system config
func (api *API) CreateSysConfig(c *common.Context) (interface{}, error){
	sysConfig := &(models.SysConfig{
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	})
	err := api.ParseSysConfig(sysConfig, c)
	if err != nil {
		return nil, err
	}

	// avoid inserting data with duplicate primary keys
	oldSysConfig, err := api.sysConfigService.GetSysConfig(sysConfig.Type, sysConfig.Key)
	if oldSysConfig!= nil {
		return nil, common.Error(common.ErrResourceHasBeenUsed,
			common.Field("error", "system config ["+oldSysConfig.Type+","+oldSysConfig.Key+"] is already exists"))
	}
	return api.sysConfigService.CreateSysConfig(sysConfig)
}

func (api *API) ParseSysConfig(sysConfig *models.SysConfig, c *common.Context) error {
	err := c.LoadBody(sysConfig)
	if err != nil {
		return common.Error(common.ErrRequestParamInvalid, common.Field("error", err.Error()))
	}

	return err
}

func (api *API) DeleteSysConfig(c *common.Context) (interface{}, error) {
	tp, key := c.Param("type"), c.Param("key")
	return api.sysConfigService.DeleteSysConfig(tp, key)
}

func (api *API) UpdateSysConfig(c *common.Context) (interface{}, error) {
	tp, key := c.Param("type"), c.Param("key")
	oldSysConfig, err := api.sysConfigService.GetSysConfig(tp, key)
	// ensure that the modified data exists
	if err != nil {
		return nil, common.Error(common.ErrResourceHasBeenUsed,
			common.Field("error", "this system config does not exist"))
	}
	oldSysConfig.UpdateTime = time.Now()

	err = c.LoadBody(oldSysConfig)
	if err != nil {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", err.Error()))
	}

	return api.sysConfigService.UpdateSysConfig(oldSysConfig)
}
