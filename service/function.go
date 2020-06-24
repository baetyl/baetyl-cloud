package service

import (
	"github.com/baetyl/baetyl-cloud/common"
	"github.com/baetyl/baetyl-cloud/config"
	"github.com/baetyl/baetyl-cloud/models"
	"github.com/baetyl/baetyl-cloud/plugin"
)

//go:generate mockgen -destination=../mock/service/function.go -package=plugin github.com/baetyl/baetyl-cloud/service FunctionService

type FunctionService interface {
	List(userID, source string) ([]models.Function, error)
	ListFunctionVersions(userID, name, source string) ([]models.Function, error)
	ListSources() []models.FunctionSource
	GetFunction(userID, name, version, source string) (*models.Function, error)
}

type functionService struct {
	functions map[string]plugin.Function
}

// NewFunctionService NewFunctionService
func NewFunctionService(config *config.CloudConfig) (FunctionService, error) {
	functions := make(map[string]plugin.Function)
	for _, v := range config.Plugin.Functions {
		cs, err := plugin.GetPlugin(v)
		if err != nil {
			return nil, err
		}
		functions[v] = cs.(plugin.Function)
	}
	return &functionService{
		functions: functions,
	}, nil
}

// List list functions
func (c *functionService) List(userID string, source string) ([]models.Function, error) {
	functionPlugin, ok := c.functions[source]
	if !ok {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "the source is not supported"))
	}
	return functionPlugin.List(userID)
}

//ListVersions List all versions of a function
func (c *functionService) ListFunctionVersions(userID, name string, source string) ([]models.Function, error) {
	functionPlugin, ok := c.functions[source]
	if !ok {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "the source is not supported"))
	}
	return functionPlugin.ListFunctionVersions(userID, name)
}

func (c *functionService) ListSources() []models.FunctionSource {
	sources := []models.FunctionSource{}
	for name := range c.functions {
		source := models.FunctionSource{
			Name: name,
		}
		sources = append(sources, source)
	}
	return sources
}

func (c *functionService) GetFunction(userID, name, version, source string) (*models.Function, error) {
	functionPlugin, ok := c.functions[source]
	if !ok {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "the source is not supported"))
	}

	return functionPlugin.Get(userID, name, version)
}
