package api

import (
	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/service"
)

//go:generate mockgen -destination=../mock/api/active.go -package=api github.com/baetyl/baetyl-cloud/v2/api ActiveAPI

type InitAPI interface {
	GetResource(c *common.Context) (interface{}, error)
}

type InitAPIImpl struct {
	initService service.InitService
}

func NewInitAPI(cfg *config.CloudConfig) (InitAPI, error) {
	initService, err := service.NewInitService(cfg)
	if err != nil {
		return nil, err
	}
	return &InitAPIImpl{
		initService:    initService,
	}, nil
}

func (api *InitAPIImpl) GetResource(c *common.Context) (interface{}, error) {
	resourceName := c.Param("resource")
	query := &struct {
		Token string `form:"token,omitempty"`
		Node  string `form:"node,omitempty"`
	}{}
	err := c.Bind(query)
	if err != nil {
		return nil, common.Error(
			common.ErrRequestParamInvalid,
			common.Field("error", err))
	}
	return api.initService.GetResource(resourceName, query.Node, query.Token)
}

