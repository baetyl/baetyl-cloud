package service

import (
	"bytes"
	"path"
	"text/template"

	"github.com/baetyl/baetyl-go/v2/errors"
	"gopkg.in/yaml.v2"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
)

//go:generate mockgen -destination=../mock/service/template.go -package=service github.com/baetyl/baetyl-cloud/v2/service TemplateService

type TemplateService interface {
	Execute(name, text string, params map[string]interface{}) ([]byte, error)
	GetTemplate(filename string) (string, error)
	ParseTemplate(filename string, params map[string]interface{}) ([]byte, error)
	UnmarshalTemplate(filename string, params map[string]interface{}, out interface{}) error
}

// TemplateServiceImpl is a service to read and parse template files.
type TemplateServiceImpl struct {
	path  string
	cache CacheService
	funcs map[string]interface{}
}

func NewTemplateService(cfg *config.CloudConfig, funcs map[string]interface{}) (TemplateService, error) {
	sCache, err := NewCacheService(cfg)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return &TemplateServiceImpl{
		path:  cfg.Template.Path,
		cache: sCache,
		funcs: funcs,
	}, nil
}

func (s *TemplateServiceImpl) Execute(name, text string, params map[string]interface{}) ([]byte, error) {
	t, err := template.New(name).Option("missingkey=error").Funcs(s.funcs).Parse(text)
	if err != nil {
		return nil, common.Error(common.ErrTemplate, common.Field("error", err))
	}
	buf := &bytes.Buffer{}
	err = t.Execute(buf, params)
	if err != nil {
		return nil, common.Error(common.ErrTemplate, common.Field("error", err))
	}
	return buf.Bytes(), nil
}

func (s *TemplateServiceImpl) GetTemplate(filename string) (string, error) {
	value, err := s.cache.GetFileData(path.Join(s.path, filename))
	if err != nil {
		return "", common.Error(common.ErrTemplate, common.Field("error", err))
	}
	return value, nil
}

func (s *TemplateServiceImpl) ParseTemplate(filename string, params map[string]interface{}) ([]byte, error) {
	tl, err := s.GetTemplate(filename)
	if err != nil {
		return nil, errors.Trace(err)
	}
	data, err := s.Execute(filename, tl, params)
	if err != nil {
		return nil, common.Error(common.ErrTemplate, common.Field("error", err))
	}
	return data, nil
}

func (s *TemplateServiceImpl) UnmarshalTemplate(filename string, params map[string]interface{}, out interface{}) error {
	tp, err := s.ParseTemplate(filename, params)
	if err != nil {
		return errors.Trace(err)
	}
	err = yaml.Unmarshal(tp, out)
	if err != nil {
		return common.Error(common.ErrTemplate, common.Field("error", err))
	}
	return nil
}