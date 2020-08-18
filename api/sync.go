package api

import (
	"github.com/baetyl/baetyl-cloud/v2/common"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
)

// Report for node report
func (api *API) Report(c *common.Context) (interface{}, error) {
	ns, n := c.GetNamespace(), c.GetName()
	var report specV1.Report
	err := c.BindJSON(&report)
	if ns == "" || n == "" {
		return nil, common.Error(common.ErrRequestParamInvalid)
	}
	if err != nil {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", err.Error()))
	}
	return api.syncService.Report(ns, n, report)
}

// Desire for node synchronize desire info
func (api *API) Desire(c *common.Context) (interface{}, error) {
	ns := c.GetNamespace()
	var request specV1.DesireRequest
	err := c.BindJSON(&request)
	if ns == "" {
		return nil, common.Error(common.ErrRequestParamInvalid)
	}
	if err != nil {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", err.Error()))
	}
	res, err := api.syncService.Desire(ns, c.GetHeader("platform"), request.Infos)
	if err != nil {
		return nil, err
	}
	return specV1.DesireResponse{Values: res}, nil
}
