package api

import (
	"bytes"
	"text/template"

	"github.com/baetyl/baetyl-cloud/common"
	"github.com/baetyl/baetyl-cloud/models"
)

// TODO: optimize this layer, general abstraction

// GetSysConfig get a system config
func (api *API) GetSysConfig(c *common.Context) (interface{}, error) {
	tp, key := c.Param("type"), c.Param("key")
	return api.sysConfigService.GetSysConfig(tp, key)
}

func (api *API) ListSysConfig(c *common.Context) (interface{}, error) {
	res, err := api.sysConfigService.ListSysConfigAll(c.Param("type"))
	if err != nil {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", err.Error()))
	}
	return &models.SysConfigView{
		SysConfigs: res,
	}, nil
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
