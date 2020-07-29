package service

import (
	"fmt"
	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
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
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", fmt.Sprintf("the source (%s) is not supported", source)))
	}
	return functionPlugin.List(userID)
}

//ListVersions List all versions of a function
func (c *functionService) ListFunctionVersions(userID, name string, source string) ([]models.Function, error) {
	functionPlugin, ok := c.functions[source]
	if !ok {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", fmt.Sprintf("the source (%s) is not supported", source)))
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
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", fmt.Sprintf("the source (%s) is not supported", source)))
	}

	return functionPlugin.Get(userID, name, version)
}
