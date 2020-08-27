package service

import (
	"bytes"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path"
	"text/template"

	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-go/v2/errors"
)

const (
	TemplateCoreConf = "baetyl-core-conf.yml"
	TemplateCoreApp  = "baetyl-core-app.yml"
	TemplateFuncConf = "baetyl-function-conf.yml"
	TemplateFuncApp  = "baetyl-function-app.yml"
)

type TemplateService struct {
	path  string
	cache *CacheService
}

func NewTemplateService(cfg *config.CloudConfig) (*TemplateService, error) {
	cacheService, err := NewCacheService(cfg)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return &TemplateService{path: cfg.Template.Path, cache: cacheService}, nil
}

func (s *TemplateService) ParseSystemTemplate(name string, data map[string]string, out interface{}) error {
	tl, err := s.getSystemTemplate(name)
	if err != nil {
		return errors.Trace(err)
	}
	t, err := template.New(name).Parse(tl)
	if err != nil {
		return errors.Trace(err)
	}
	buf := &bytes.Buffer{}
	err = t.Execute(buf, data)
	if err != nil {
		return errors.Trace(err)
	}
	return yaml.Unmarshal(buf.Bytes(), out)
}

func (s *TemplateService) getSystemTemplate(name string) (string, error) {
	return s.cache.Get(name, func(key string) (string, error) {
		data, err := ioutil.ReadFile(path.Join(s.path, key))
		if err != nil {
			return "", errors.Trace(err)
		}
		return string(data), nil
	})
}
