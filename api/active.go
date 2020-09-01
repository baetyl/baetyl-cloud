package api

import (
	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/service"
)

//go:generate mockgen -destination=../mock/api/active.go -package=api github.com/baetyl/baetyl-cloud/v2/api ActiveAPI

type ActiveAPI interface {
	GetResource(c *common.Context) (interface{}, error)
}

type ActiveAPIImpl struct {
	activeService service.ActiveService
}

func NewActiveAPI(cfg *config.CloudConfig) (ActiveAPI, error) {
	activeService, err := service.NewActiveService(cfg)
	if err != nil {
		return nil, err
	}
	return &ActiveAPIImpl{
		activeService:    activeService,
	}, nil
}

func (api *ActiveAPIImpl) GetResource(c *common.Context) (interface{}, error) {
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
	return api.activeService.GetResource(resourceName, query.Node, query.Token)
}

