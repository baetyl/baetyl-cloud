package api

import (
	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/service"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
)

type SyncAPI struct {
	SyncService service.SyncService
}

func NewSyncAPI(cfg *config.CloudConfig) (*SyncAPI, error) {
	syncService, err := service.NewSyncService(cfg)
	if err != nil {
		return nil, err
	}
	return &SyncAPI{
		SyncService: syncService,
	}, nil
}

// Report for node report
func (s *SyncAPI) Report(c *common.Context) (interface{}, error) {
	ns, n := c.GetNamespace(), c.GetName()
	var report specV1.Report
	err := c.BindJSON(&report)
	if ns == "" || n == "" {
		return nil, common.Error(common.ErrRequestParamInvalid)
	}
	if err != nil {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", err.Error()))
	}
	return s.SyncService.Report(ns, n, report)
}

// Desire for node synchronize desire info
func (s *SyncAPI) Desire(c *common.Context) (interface{}, error) {
	ns := c.GetNamespace()
	var request specV1.DesireRequest
	err := c.BindJSON(&request)
	if ns == "" {
		return nil, common.Error(common.ErrRequestParamInvalid)
	}
	if err != nil {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", err.Error()))
	}
	metadata := map[string]string{}
	for k := range c.Request.Header {
		metadata[k] = c.GetHeader(k)
	}
	res, err := s.SyncService.Desire(ns, request.Infos, metadata)
	if err != nil {
		return nil, err
	}
	return specV1.DesireResponse{Values: res}, nil
}
