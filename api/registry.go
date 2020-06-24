package api

import (
	"fmt"
	"time"

	"github.com/baetyl/baetyl-cloud/common"
	"github.com/baetyl/baetyl-cloud/models"
	"github.com/baetyl/baetyl-go/errors"
	specV1 "github.com/baetyl/baetyl-go/spec/v1"
)

// TODO: optimize this layer, general abstraction

// GetRegistry get a Registry
func (api *API) GetRegistry(c *common.Context) (interface{}, error) {
	ns, n := c.GetNamespace(), c.GetNameFromParam()
	res, err := wrapRegistry(api.secretService.Get(ns, n, ""))
	if err != nil {
		return nil, err
	}
	return hidePwd(res), nil
}

// ListRegistry list Registry
func (api *API) ListRegistry(c *common.Context) (interface{}, error) {
	ns := c.GetNamespace()
	res, err := wrapRegistryList(api.secretService.List(ns, wrapRegistryListOption(api.parseListOptions(c))))
	if err != nil {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", err.Error()))
	}
	for i := range res.Items {
		hidePwd(&res.Items[i])
	}
	return res, err
}

// CreateRegistry create one Registry
func (api *API) CreateRegistry(c *common.Context) (interface{}, error) {
	cfg, err := api.parseAndCheckRegistryModel(c)
	if err != nil {
		return nil, err
	}
	ns, name := c.GetNamespace(), cfg.Name
	sd, err := api.secretService.Get(ns, name, "")
	if err != nil {
		if e, ok := err.(errors.Coder); !ok || e.Code() != common.ErrResourceNotFound {
			return nil, err
		}
	}

	if sd != nil {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "this name is already in use"))
	}
	if err = api.validateRegistryModel(cfg); err != nil {
		return nil, err
	}
	res, err := wrapRegistry(api.secretService.Create(ns, cfg.ToSecret()))
	return hidePwd(res), err
}

// UpdateRegistry update the Registry
func (api *API) UpdateRegistry(c *common.Context) (interface{}, error) {
	cfg, err := api.parseAndCheckRegistryModel(c)
	if err != nil {
		return nil, err
	}
	ns, n := c.GetNamespace(), c.GetNameFromParam()
	sd, err := wrapRegistry(api.secretService.Get(ns, n, ""))
	if err != nil {
		return nil, err
	}
	// only edit description by design
	if cfg.Description == sd.Description {
		return hidePwd(sd), nil
	}
	sd.Description = cfg.Description
	sd.UpdateTimestamp = time.Now()
	if err = api.validateRegistryModel(sd); err != nil {
		return nil, err
	}
	res, err := wrapRegistry(api.secretService.Update(ns, sd.ToSecret()))
	if err != nil {
		return nil, err
	}
	return hidePwd(res), nil
}

func (api *API) RefreshRegistryPassword(c *common.Context) (interface{}, error) {
	cfg, err := api.parseAndCheckRegistryModel(c)
	if err != nil {
		return nil, err
	}
	ns, n := c.GetNamespace(), c.GetNameFromParam()
	sd, err := wrapRegistry(api.secretService.Get(ns, n, ""))
	if err != nil {
		return nil, err
	}
	sd.UpdateTimestamp = time.Now()
	sd.Password = cfg.Password
	res, err := wrapRegistry(api.secretService.Update(ns, sd.ToSecret()))
	if err != nil {
		return nil, err
	}
	if err = api.validateRegistryModel(res); err != nil {
		return nil, err
	}
	err = api.updateAppSecret(ns, res.ToSecret())
	if err != nil {
		return nil, err
	}
	return hidePwd(res), nil
}

// DeleteRegistry delete the Registry
func (api *API) DeleteRegistry(c *common.Context) (interface{}, error) {
	ns, n := c.GetNamespace(), c.GetNameFromParam()
	_, err := wrapRegistry(api.secretService.Get(ns, n, ""))
	if err != nil {
		return nil, err
	}
	return api.deleteSecret(ns, n, "registry")
}

// GetAppByRegistry list app
func (api *API) GetAppByRegistry(c *common.Context) (interface{}, error) {
	ns, n := c.GetNamespace(), c.GetNameFromParam()
	res, err := wrapRegistry(api.secretService.Get(ns, n, ""))
	if err != nil {
		return nil, err
	}
	return api.listAppBySecret(ns, res.Name)
}

// parseAndCheckRegistryModel parse and check the config model
func (api *API) parseAndCheckRegistryModel(c *common.Context) (*models.Registry, error) {
	registry := new(models.Registry)
	registry.Name = c.GetNameFromParam()
	err := c.LoadBody(registry)
	if err != nil {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", err.Error()))
	}
	if name := c.GetNameFromParam(); name != "" {
		registry.Name = name
	}
	if registry.Name == "" {
		err = common.Error(common.ErrRequestParamInvalid, common.Field("error", "name is required"))
	}
	return registry, err
}

func hidePwd(r *models.Registry) *models.Registry {
	if r != nil {
		r.Password = ""
	}
	return r
}

func wrapRegistry(s *specV1.Secret, e error) (*models.Registry, error) {
	if s != nil {
		return models.FromSecret(s), e
	}
	return nil, e
}

func wrapRegistryList(s *models.SecretList, e error) (*models.RegistryList, error) {
	if s != nil {
		return models.FromSecretList(s), e
	}
	return nil, e
}

func wrapRegistryListOption(lo *models.ListOptions) *models.ListOptions {
	// TODO 增加type字段代替label标签
	lo.LabelSelector = fmt.Sprintf("%s=%s", specV1.SecretLabel, specV1.SecretRegistry)
	return lo
}

// validateRegistryModel validate the registry model
func (api *API) validateRegistryModel(r *models.Registry) error {
	if r.Address == "" || r.Username == "" || r.Password == "" {
		return common.Error(common.ErrRequestParamInvalid, common.Field("error", "address/username/password is required"))
	}
	return nil
}
