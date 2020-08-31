package api

import (
	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/service"
)

type ActiveAPI interface {
	GetResource(c *common.Context) (interface{}, error)
}

type ActiveAPIImpl struct {
	service.ActiveService
}

func NewActiveAPI(cfg *config.CloudConfig) (ActiveAPI, error) {
	activeService, err := service.NewActiveService(cfg)
	if err != nil {
		return nil, err
	}
	return &ActiveAPIImpl{
		activeService,
	}, nil
}

func (api *ActiveAPIImpl) GetResource(c *common.Context) (interface{}, error) {
	return api.ActiveService.GetResource(c)
}

