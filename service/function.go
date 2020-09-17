package service

import (
	"fmt"

	"github.com/baetyl/baetyl-go/v2/errors"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

//go:generate mockgen -destination=../mock/service/function.go -package=service github.com/baetyl/baetyl-cloud/v2/service FunctionService

const functionRuntimePrefix = "baetyl-function-runtime-"

type FunctionService interface {
	List(userID, source string) ([]models.Function, error)
	ListFunctionVersions(userID, name, source string) ([]models.Function, error)
	ListSources() []models.FunctionSource
	ListRuntimes() (map[string]string, error)
	GetFunction(userID, name, version, source string) (*models.Function, error)
}

type functionService struct {
	prop      PropertyService
	functions map[string]plugin.Function
}

// NewFunctionService NewFunctionService
func NewFunctionService(cfg *config.CloudConfig) (FunctionService, error) {
	sProp, err := NewPropertyService(cfg)
	if err != nil {
		return nil, errors.Trace(err)
	}
	functions := make(map[string]plugin.Function)
	for _, v := range cfg.Plugin.Functions {
		cs, err := plugin.GetPlugin(v)
		if err != nil {
			return nil, err
		}
		functions[v] = cs.(plugin.Function)
	}
	return &functionService{
		prop:      sProp,
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

func (c *functionService) ListRuntimes() (map[string]string, error) {
	res, err := c.prop.ListProperty(&models.Filter{Name: functionRuntimePrefix})
	if err != nil {
		return nil, errors.Trace(err)
	}
	runtimes := make(map[string]string)
	for _, item := range res {
		runtimes[item.Name[len(functionRuntimePrefix):]] = item.Value
	}
	return runtimes, nil
}

func (c *functionService) GetFunction(userID, name, version, source string) (*models.Function, error) {
	functionPlugin, ok := c.functions[source]
	if !ok {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", fmt.Sprintf("the source (%s) is not supported", source)))
	}

	return functionPlugin.Get(userID, name, version)
}
